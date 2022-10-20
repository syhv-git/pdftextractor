# PDFtExtractor
___
This package performs a text extraction on PDF files, with a focus on performance. The function `ExtractText` takes two parameters; the path to the PDF file, and a boolean value defining whether to extract text from drawn images as well. It returns a byte slice of the text contents.

The current package has some issues when handling the PDF text objects, and does not properly decode the PDF glyphs

## Usage
To use this package in your project, run the following command in your module:
```
go get -u github.com/syhv-git/pdftextractor
```
You will also need to download `tesseract-ocr` and `libtesseract-dev` (`tesseract-ocr-dev` for apk). You can add more language data to `/usr/share/tesseract-ocr/$VERSION/tessdata/`

> There may be dependency issues with the Gosseract package. This will require the linux mint package from the same developer

## Roadmap
* Decode PDF string objects and extract the raw text
* Test with PDFs containing images
* Test with an Image based PDF file
* Test interoperability with other PDF versions
* Optimize the codebase
