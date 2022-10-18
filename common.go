package pdftextractor

import (
	"bufio"
	"bytes"
	"log"
	"strconv"
	"strings"
)

//func loadGlyphs() {
//	b, err := os.ReadFile("glyphs.txt")
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	wg := &sync.WaitGroup{}
//	wg.Add(4280)
//	buf := bytes.NewBuffer(b)
//	for {
//		l, err := buf.ReadBytes('\n')
//		if err != nil {
//			log.Fatal(err.Error())
//		}
//		if len(l) < 2 {
//			break
//		}
//		if bytes.Contains(l, []byte{'#'}) {
//			continue
//		}
//		go func(i []byte) {
//			defer wg.Done()
//			if err != nil {
//				log.Fatal(err.Error())
//			}
//			v := bytes.Split(i, []byte{';'})
//			k := string(v[0])
//			glyphs[k] = string(v[1])
//		}(l)
//	}
//	wg.Wait()
//}

func parsePDF(rdr *bufio.Reader, objs objMap) {
	for {
		b, err := rdr.ReadBytes('\n')
		if err != nil {
			break
		}

		if bytes.Contains(b, edition) {
			_, r, f := bytes.Cut(b, edition)
			if !f {
				log.Fatal("Unable to detect PDF version")
			}
			if r[len(r)-1] == '\n' {
				r = r[:len(r)-1]
			}
			version = string(r)
		}

		if bytes.Contains(b, newObj) {
			id := getObjID(b)
			if id == 0 {
				continue
			}
			objs[id] = &objDict{refs: refMap{}}
			current := id
			extractObject(objs, rdr, current)
		}
	}
}

func getObjID(b []byte) uint64 {
	args := strings.Split(string(b), " ")
	if len(args) != 3 {
		return 0
	}

	r, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		log.Fatal(err.Error())
	}
	return r
}

func getObjectOrder(objs objMap) (order []uint64) {
	pgs := objs[head].refs["Kids"]
	for _, x := range pgs {
		c := objs[x].refs["Contents"]
		order = append(order, c...)
	}
	return
}
