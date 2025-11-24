package main

import (
    "github.com/pdfcpu/pdfcpu/pkg/api"
    pdf "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type WatermarkRequest struct {
    PdfBase64 string `json:"pdf_base64"`
    Text      string `json:"text"`
}

func main() {
    http.HandleFunc("/wm", handleWatermark)
    log.Println("Server started on :8080")
    http.ListenAndServe(":8080", nil)
}

func handleWatermark(w http.ResponseWriter, r *http.Request) {
    var req WatermarkRequest
    json.NewDecoder(r.Body).Decode(&req)

    pdfBytes, _ := base64.StdEncoding.DecodeString(req.PdfBase64)

    reader := bytes.NewReader(pdfBytes)
    var buf bytes.Buffer

    wm := model.DefaultWatermarkConfig()
    wm.Mode = model.WMText
    wm.Diagonal = true
    wm.Scale = 0.9
    wm.Opacity = 0.2
    wm.TextLines = []string{req.Text}

    err := pdfcpu.AddWatermarks(reader, &buf, nil, wm)
    if err != nil {
        w.WriteHeader(500)
        w.Write([]byte(err.Error()))
        return
    }

    encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"pdf_base64":"` + encoded + `"}`))
}
