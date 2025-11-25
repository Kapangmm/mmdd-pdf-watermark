// Remark: v2025-11-24.r1 | Minimal HTTP watermark server using pdfcpu v0.11.1
// Created: 2025-11-24

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

const (
	// ถ้า client ไม่ส่ง ?text= มา จะใช้ข้อความนี้เป็น watermark
	defaultWatermarkText = "MMDD"

	// description ของ watermark ตาม syntax ของ pdfcpu
	// - pos:br       = วางมุมขวาล่าง (bottom-right)
	// - op:0.5       = opacity 50%
	// - rot:0        = ไม่หมุน
	defaultWatermarkDesc = "pos:br, op:0.5, rot:0"
)

func main() {
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/watermark", watermarkHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[mmdd-pdf-watermark] starting server on :%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func watermarkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("use POST"))
		return
	}

	// จำกัดขนาดไฟล์สูงสุด ~20MB
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
	defer r.Body.Close()

	pdfData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(pdfData) == 0 {
		http.Error(w, "empty body, send PDF bytes as body", http.StatusBadRequest)
		return
	}

	// รับข้อความ watermark จาก query string: ?text=...
	text := r.URL.Query().Get("text")
	if text == "" {
		text = defaultWatermarkText
	}

	// สามารถ override description ได้ผ่าน ?desc=...
	desc := r.URL.Query().Get("desc")
	if desc == "" {
		desc = defaultWatermarkDesc
	}

	// เขียน PDF เข้า temp file (ฝั่ง Railway ใช้ /tmp ได้)
	inFile, err := os.CreateTemp("", "mmdd-in-*.pdf")
	if err != nil {
		http.Error(w, "failed to create temp input file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(inFile.Name())

	if _, err := io.Copy(inFile, bytes.NewReader(pdfData)); err != nil {
		inFile.Close()
		http.Error(w, "failed to write temp input file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := inFile.Close(); err != nil {
		http.Error(w, "failed to close temp input file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	outFile, err := os.CreateTemp("", "mmdd-out-*.pdf")
	if err != nil {
		http.Error(w, "failed to create temp output file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	outPath := outFile.Name()
	outFile.Close()
	defer os.Remove(outPath)

	// config พื้นฐานของ pdfcpu
	conf := model.NewDefaultConfiguration()

	// แปลง text + desc เป็น *model.Watermark
	wm, err := pdfcpu.ParseTextWatermarkDetails(text, desc, false, conf.Unit)
	if err != nil {
		http.Error(w, "failed to parse watermark: "+err.Error(), http.StatusBadRequest)
		return
	}

    // ถ้า selectedPages เป็น nil = ใส่ watermark ทุกหน้า
    var selectedPages []string // nil slice

    // ใช้ API ระดับไฟล์โดยตรง
    if err := api.AddWatermarksFile(inFile.Name(), outPath, selectedPages, wm, conf); err != nil {
        http.Error(w, "failed to apply watermark: "+err.Error(), http.StatusInternalServerError)
        return
    }

	result, err := os.ReadFile(outPath)
	if err != nil {
		http.Error(w, "failed to read output file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`inline; filename="watermarked-%d.pdf"`, time.Now().Unix()))
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(result); err != nil {
		log.Printf("write response error: %v", err)
	}
}
