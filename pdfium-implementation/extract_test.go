package pdfium_implementation

import (
	"fmt"
	"testing"
)

func TestExtractText(t *testing.T) {
	src := "sample.pdf"
	//b, _ := os.ReadFile(src)
	//os.WriteFile("output.txt", b, 0666)
	b := ExtractText(src)
	if len(b) < 1 {
		t.Error("Error when extracting text")
	}
	fmt.Println(b)
}
