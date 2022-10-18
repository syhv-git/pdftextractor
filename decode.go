package pdftextractor

import (
	"bufio"
	"bytes"
	"compress/lzw"
	"compress/zlib"
	"image/jpeg"
	"io"
	"log"
	"strings"
)

func decoder(rdr *bufio.Reader, level string) []byte {
	switch {
	case strings.Contains(level, lzwc):
		return decodeLZW(rdr)
	case strings.Contains(level, zip):
		return decodeZLib(rdr)
	default:
		b := &bytes.Buffer{}
		if _, err := rdr.WriteTo(b); err != nil {
			log.Fatal(err.Error())
		}
		return decodeText(b.Bytes())
	}
	return nil
}

func decodeZLib(rdr *bufio.Reader) []byte {
	return decodeText(decodeZLibStream(rdr))
}

func decodeZLibStream(rdr *bufio.Reader) []byte {
	dec, err := zlib.NewReader(rdr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(d io.ReadCloser) {
		if e := d.Close(); e != nil {
			log.Fatal(e.Error())
		}
	}(dec)

	res := &bytes.Buffer{}
	buf := bufio.NewReader(dec)
	if _, err = buf.WriteTo(res); err != nil {
		log.Fatal(err.Error())
	}
	return res.Bytes()
}

func decodeLZW(rdr *bufio.Reader) []byte {
	return decodeText(decodeLZWStream(rdr))
}

func decodeLZWStream(rdr *bufio.Reader) []byte {
	dec := lzw.NewReader(rdr, 1, 8)
	defer func() {
		if err := dec.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()
	ret := &bytes.Buffer{}
	if _, err := ret.ReadFrom(dec); err != nil {
		log.Fatal(err.Error())
	}
	return ret.Bytes()
}

func decodeText(strm []byte) []byte {
	buf := bytes.NewBuffer(strm)
	res := &bytes.Buffer{}

	for {
		b, err := buf.ReadBytes('\n')
		if err != nil {
			return res.Bytes()
		}
		if bytes.Contains(b, text) {
			if _, err = res.Write(decodeString(b)); err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}

func decodeString(b []byte) []byte {
	_, r, _ := bytes.Cut(b, textpre)
	l, _, _ := bytes.Cut(r, text)

	if bytes.Contains(l, []byte{'('}) {
		return decodeRawString(l)
	}
	if bytes.Contains(l, []byte{'<'}) {
		return decodeHexString(l)
	}
	return nil
}

func decodeRawString(b []byte) []byte {
	buf := &bytes.Buffer{}
	l, r, _ := bytes.Cut(b, []byte{'('})
	for {
		l, r, _ = bytes.Cut(r, []byte{')'})

		if _, err := buf.Write(l); err != nil {
			log.Fatal(err.Error())
		}
		if !bytes.Contains(r, []byte{'('}) {
			break
		}
	}
	return buf.Bytes()
}

func decodeHexString(b []byte) []byte {
	buf := &bytes.Buffer{}
	l, r, _ := bytes.Cut(b, []byte{'<'})
	for {
		l, r, _ = bytes.Cut(r, []byte{'>'})

		for len(l) > 3 {
			m := uint64(l[0]) - 48
			n := uint64(l[1]) - 48
			o := uint64(l[2]) - 48
			p := uint64(l[3]) - 48
			rn := (m << 12) | (n << 8) | (o << 4) | p
			buf.WriteRune(rune(rn))
			l = l[4:]
		}

		if !bytes.Contains(r, []byte{'<'}) {
			break
		}
		_, r, _ = bytes.Cut(r, []byte{'<'})
	}
	return buf.Bytes()
}

func decodeJPEG(rdr *bufio.Reader) []byte {
	img, err := jpeg.Decode(rdr)
	if err != nil {
		log.Fatal(err.Error())
	}

	buf := &bytes.Buffer{}
	if err = jpeg.Encode(buf, img, nil); err != nil {
		log.Fatal(err.Error())
	}
	return buf.Bytes()
}
