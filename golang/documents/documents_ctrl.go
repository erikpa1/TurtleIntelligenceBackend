package documents

import (
	"bytes"
	"fmt"
	"github.com/ledongthuc/pdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/db"
	"turtle/lg"
	"turtle/models"
	"turtle/vfs"
)

const CT_DOC = "docs"
const CT_DOC_EMBEDDINGS = "docs_embeddings"

func ListDocument(user *models.User) []*Document {
	//TODO vytvorit safe databazu, ktora pracuje aj s uzivatelom

	return db.QueryEntities[Document](CT_DOC, bson.M{
		"org": user.Org,
	})

}

func InsertDocument(document *Document, documentData []byte) {

	document.Uid = primitive.ObjectID{}

	db.InsertEntity(CT_DOC, document)

	fileName := fmt.Sprintf("%s.%s", document.Uid.Hex(), document.Extension)

	db.SC.UploadFile("documents", fileName, documentData)

	filePath := vfs.GetFilePathFromWD("documents", fileName)

	ExtractPdfTextInMemory(documentData)

	if document.HasEmbedding {
		if CreateDocFileEmbedding(filePath) != nil {
			db.UpdateOneCustom("documents",
				bson.M{"_id": document.Uid},
				bson.M{"hasEmbedding": false},
			)
		}
	}

}

func ExtractPdfTextInMemory(data []byte) string {
	// Create a bytes.Reader which implements io.ReaderAt
	reader := bytes.NewReader(data)

	// Now you can use it with your PDF function
	pdfReader, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		lg.LogE(err.Error())
		return ""
	}

	var buf bytes.Buffer

	b, err := pdfReader.GetPlainText()

	if err != nil {
		lg.LogE(err.Error())
		return ""
	}
	buf.ReadFrom(b)
	content := buf.String()
	lg.LogE(content)
	// Work with your PDF...

	return ""
}

func CreateDocFileEmbedding(filePath string) error {
	pdf.DebugOn = true

	f, r, err := pdf.Open("./pdf_test.pdf")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return err
	}
	buf.ReadFrom(b)
	content := buf.String()
	fmt.Println(content)

	return nil

}

func AddDocumentEmbedding(documentUid primitive.ObjectID, embedding [][]float32) {

	emb := DocumentEmbedding{}
	emb.Uid = documentUid
	emb.Embedding = embedding

	db.InsertEntity(CT_DOC_EMBEDDINGS, emb)

}

func DeleteDocument() {

}
