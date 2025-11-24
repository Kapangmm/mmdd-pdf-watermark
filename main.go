package main

import (
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func main() {
	// ไฟล์ทดสอบแบบง่าย ๆ ก่อน
	in := "input.pdf"
	out := "output.pdf"

	// ข้อความลายน้ำ (ตอนนี้ใช้ข้อความสั้น ๆ ก่อน)
	wmText := "Test Watermark"

	// สร้าง watermark object จากข้อความ
	// params ตัวอย่าง:
	//  - "rot:45"  = หมุน 45 องศา
	//  - "op:0.6"  = opacity 60%
	//  - "sc:1 abs"= scale เต็มหน้าแบบ absolute
	wm, err := pdfcpu.ParseTextWatermark(
		wmText,
		"rot:45, op:0.6, sc:1 abs",
		true,  // onTop:   วางทับด้านบนเนื้อหา
		true,  // update:  อัปเดตแทนที่ watermark เดิมถ้ามี
		pdfcpu.POINTS,
	)
	if err != nil {
		log.Fatal(err)
	}

	// selectedPages = nil  => ใส่ทุกหน้า
	// conf          = nil  => ใช้ default configuration
	if err := api.AddWatermarksFile(in, out, nil, wm, nil); err != nil {
		log.Fatal(err)
	}
}
