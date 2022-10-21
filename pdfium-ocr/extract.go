package pdfium_ocr

import (
	"bytes"
	render "github.com/brunsgaard/go-pdfium-render"
	"github.com/otiai10/gosseract/v2"
	"image/jpeg"
	"log"
	"os"
)

func ExtractText(src string) []byte {
	render.InitLibrary()
	defer render.DestroyLibrary()
	c := gosseract.NewClient()
	defer c.Close()

	b, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(err.Error())
	}

	doc, err := render.NewDocument(&b)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer doc.Close()

	txt := &bytes.Buffer{}
	for i, j := 1, doc.GetPageCount(); i <= j; i++ {
		buf := &bytes.Buffer{}
		img := doc.RenderPage(i, 0)
		if err = jpeg.Encode(buf, img, nil); err != nil {
			log.Fatal(err.Error())
		}
		if err = c.SetImageFromBytes(buf.Bytes()); err != nil {
			log.Fatal(err.Error())
		}
		tmp, err := c.Text()
		if err != nil {
			log.Fatal(err.Error())
		}
		txt.Write([]byte(tmp))
	}
	return txt.Bytes()
}
