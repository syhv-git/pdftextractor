package pdftextractor

import (
	"bytes"
	"log"
	"strconv"
	"strings"
)

var (
	includeImages bool
	head          uint64
	version       string
	glyphs        map[string]string
)

var (
	baseenc  = []byte("/BaseEncoding")
	diff     = []byte("/Differences")
	filter   = []byte("/Filter")
	image    = []byte("/Image")
	indirect = []byte(" 0 R")
	length   = []byte("/Length ")
	touni    = []byte("/ToUnicode")
	typ      = []byte("/Type")
)

var (
	edition   = []byte("%PDF-")
	newObj    = []byte(" 0 obj")
	endObj    = []byte("endobj")
	stream    = []byte("stream")
	endStream = []byte("endstream")
)

var (
	objDelim = []byte{'>', '/', '\n'}
	textpre  = []byte("Td ")
	text     = []byte(" Tj")
	zip      = "Flate"
	lzwc     = "LZW"
)

type objDict struct {
	typ     string
	base    string
	filter  string
	refs    refMap
	length  uint64
	content []byte
}

type refMap map[string][]uint64

type objMap map[uint64]*objDict

func (o *objDict) setBaseEncoding(b []byte) {
	_, r, f := bytes.Cut(b, baseenc)
	if !f {
		return
	}
	m := strings.Split(string(r), "/")
	if len(m) < 2 {
		return
	}

	v := m[1]
	if v[len(v)-1] == '\n' {
		v = v[:len(v)-1]
	}
	o.base = v
}

func (o *objDict) setType(b []byte) {
	_, r, f := bytes.Cut(b, typ)
	if !f {
		return
	}
	m := strings.Split(string(r), "/")
	if len(m) < 2 {
		return
	}

	v := m[1]
	if v[len(v)-1] == '\n' {
		v = v[:len(v)-1]
	}
	o.typ = v
}

func (o *objDict) setFilter(b []byte) {
	_, r, f := bytes.Cut(b, filter)
	if !f {
		return
	}
	m := strings.Split(string(r), "/")
	if len(m) < 2 {
		return
	}

	v := m[1]
	if v[len(v)-1] == '\n' {
		v = v[:len(v)-1]
	}
	o.filter = v
}

func (o *objDict) setReferences(b []byte) {
	l, r, f := bytes.Cut(b, indirect)
	if !f {
		return
	}
	name := ""

	for {
		i := bytes.Split(l, []byte{'/'})
		if len(i) > 1 {
			j := bytes.Split(i[len(i)-1], []byte{' '})
			v := j[1]
			if bytes.Contains(v, []byte{'['}) {
				name = string(j[0])
				v = v[1:]
			}

			w, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				log.Fatal(err.Error())
			}
			o.refs[string(j[0])] = append(o.refs[string(j[0])], w)
		} else if name != "" {
			v := bytes.TrimPrefix(i[0], []byte{' '})
			w, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				log.Fatal(err.Error())
			}
			o.refs[name] = append(o.refs[name], w)
		}

		l, r, f = bytes.Cut(r, indirect)
		if !f {
			break
		}
	}
}

func (o *objDict) setLength(b []byte) {
	_, r, f := bytes.Cut(b, length)
	if !f {
		o.length = 0
		return
	}

	buf := &bytes.Buffer{}
	for _, x := range r {
		if bytes.ContainsAny([]byte{x}, string(objDelim)) {
			break
		}
		if err := buf.WriteByte(x); err != nil {
			log.Fatal(err.Error())
		}
	}

	if i := buf.Len(); i < 8 {
		m := buf.Bytes()
		n := make([]byte, 8)
		for h := 0; h < 8-i; h++ {
			n[h] = '\x00'
		}
		for j, k := 0, 8-i; k < 8; j, k = j+1, k+1 {
			n[k] = m[j]
		}
	} else if i > 8 {
		log.Fatal("Error when parsing the length")
	}

	w, err := strconv.ParseUint(string(buf.Bytes()), 10, 64)
	if err != nil {
		log.Fatal(err.Error())
	}

	o.length = w
}
