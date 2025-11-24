package main

import (
    "log"

    "github.com/pdfcpu/pdfcpu/pkg/api"
    "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
    pdfcpu "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func main() {

    in := "input.pdf"
    out := "output.pdf"

    // 1) สร้าง config
    conf := model.NewDefaultConfiguration()

    // 2) ตั้งข้อความ watermark
    text := "© 2025 Myanmar Daily Digest – LINE User: XXXXXXXXXX…"

    // 3) text watermark options
    // rot = หมุนองศา
    // sc  = scale ขนาด watermark
    // op  = opacity ความโปร่ง (0.0 – 1.0)
    details := "rot:45, sc:0.8, op:0.25, fillc:#666666"

    // onTop = false → ให้ watermark อยู่ใต้ layer บางอย่าง แต่ยังทับเนื้อหา
    onTop := false

    // 4) สร้าง watermark object
    wm, err := pdfcpu.ParseTextWatermarkDetails(text, details, onTop, conf.Unit)
    if err != nil {
        log.Fatal(err)
    }

    // 5) selectedPages = nil → ทุกหน้า
    var pages []string = nil

    // 6) apply watermark
    err = api.AddWatermarksFile(in, out, pages, wm, conf)
    if err != nil {
        log.Fatal(err)
    }

}
