package ctrlApp

import (
	"bytes"
	"fmt"
	"github.com/ledongthuc/pdf"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/db"
	"turtle/modelsApp"
)

const CT_DOC = "docs"
const CT_DOC_EMBEDDINGS = "docs_embeddings"

func AddDocumentEmbedding(documentUid primitive.ObjectID, embedding [][]float32) {

	emb := modelsApp.DocumentEmbedding{}
	emb.Uid = documentUid
	emb.Embedding = embedding

	db.InsertEntity(CT_DOC_EMBEDDINGS, emb)

}

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
