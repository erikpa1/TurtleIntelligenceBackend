package ctrlApp

import (
	"bytes"
	"fmt"
	"github.com/ledongthuc/pdf"
	"turtle/modelsApp"
)

func COUDocument(document *modelsApp.Document, documentData []byte) {

	pdf.DebugOn = true

	f, r, err := pdf.Open("./pdf_test.pdf")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		panic(err)
	}
	buf.ReadFrom(b)
	content := buf.String()
	fmt.Println(content)

}

func DeleteDocument() {

}
