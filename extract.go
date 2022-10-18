package pdftextractor

import (
	"bufio"
	"bytes"
	"github.com/otiai10/gosseract/v2"
	"log"
	"os"
	"strings"
)

// ExtractText extracts the text of the source PDF and optionally extracts text from drawn images.
// Text from images are out of order from extracted text objects
func ExtractText(src string, incl bool) []byte {
	objs := objMap{}
	includeImages = incl
	f, err := os.Open(src)
	if err != nil {
		log.Fatal(err.Error())
	}
	//loadGlyphs()

	rdr := bufio.NewReader(f)
	parsePDF(rdr, objs)
	if err = f.Close(); err != nil {
		log.Fatal(err.Error())
	}

	o := getObjectOrder(objs)
	buf := &bytes.Buffer{}
	for _, x := range o {
		if _, err = buf.Write(objs[x].content); err != nil {
			log.Fatal(err.Error())
		}
	}

	return buf.Bytes()
}

func extractObject(objs objMap, rdr *bufio.Reader, curr uint64) {
	img := false
	for {
		b, err := rdr.ReadBytes('\n')
		if err != nil {
			delete(objs, curr)
			return
		}
		if bytes.Contains(b, endObj) {
			return
		}
		if bytes.Contains(b, typ) {
			objs[curr].setType(b)
			if objs[curr].typ == "Pages" {
				head = curr
			}
		}

		if bytes.Contains(b, filter) {
			objs[curr].setFilter(b)
		}
		if bytes.Contains(b, length) {
			objs[curr].setLength(b)
		}
		if bytes.Contains(b, indirect) {
			objs[curr].setReferences(b)
		}
		if bytes.Contains(b, image) && includeImages {
			img = true
		}
		if bytes.Contains(b, stream) {
			rawrdr := extractStreamContent(rdr)
			if img {
				tmp := extractImage(rawrdr, objs[curr])
				c := gosseract.NewClient()
				if err = c.SetImageFromBytes(tmp); err != nil {
					log.Fatal(err.Error())
				}
				txt, err := c.Text()
				if err != nil {
					log.Fatal(err.Error())
				}
				objs[curr].content = []byte(txt)
				return
			}
			objs[curr].content = decoder(rawrdr, objs[curr].filter)
			return
		}
	}
}

func extractStreamContent(rdr *bufio.Reader) *bufio.Reader {
	buf := &bytes.Buffer{}
	for {
		ln, err := rdr.ReadBytes('\n')
		if err != nil {
			log.Fatal(err.Error())
		}

		if bytes.Contains(ln, endStream) {
			return bufio.NewReader(buf)
		}
		if _, err = buf.Write(ln); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func extractImage(rdr *bufio.Reader, obj *objDict) []byte {
	switch {
	case strings.Contains(obj.filter, "JPEG"):
		return decodeJPEG(rdr)
	case strings.Contains(obj.filter, zip):
		return decodeZLibStream(rdr)
	default:
		return nil
	}
}
