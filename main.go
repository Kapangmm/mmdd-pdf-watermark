package main

import (
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func main() {
	in := "input.pdf"
	out := "output.pdf"

	conf := pdfcpu.NewDefaultConfiguration()
	err := api.AddWatermarksFile(in, out, nil, "Test Watermark", conf)
	if err != nil {
		log.Fatal(err)
	}
}
