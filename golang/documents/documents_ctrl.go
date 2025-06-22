package documents

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/db"
	"turtle/lg"
	"turtle/llm/llmCtrl"
	"turtle/models"
)

const CT_DOC = "docs"
const CT_DOC_EXTRACT = "docs_extracts"
const CT_DOC_EMBEDDINGS = "docs_embeddings"

func ListDocument(user *models.User) []*Document {

	return db.QueryEntities[Document](CT_DOC, bson.M{
		"org": user.Org,
	})

}

type InsertDocumentParams struct {
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	LLMDescription   bool               `json:"llmDescription"`
	CreateEmbedding  bool               `json:"createEmbedding"`
	DescriptionModel primitive.ObjectID `json:"descriptionModel"`
}

func CreateAndUploadDocument(c *gin.Context, user *models.User, uploadParams *InsertDocumentParams, documentData []byte) {

	document := &Document{}
	document.Uid = primitive.NewObjectID()
	document.Name = uploadParams.Name
	document.Description = uploadParams.Description

	db.InsertEntity(CT_DOC, document)

	fileName := fmt.Sprintf("%s.%s", document.Uid.Hex(), document.Extension)

	db.SC.UploadFile("documents", fileName, documentData)

	if uploadParams.CreateEmbedding || uploadParams.LLMDescription {
		pdfText, err := ExtractPdfTextInMemory(documentData)

		db.InsertEntity(CT_DOC_EXTRACT, bson.M{
			"_id":     document.Uid,
			"extract": pdfText,
		})

		if err != nil {

			if uploadParams.CreateEmbedding {
				embedding, embError := llmCtrl.CreateStringEmbedding(context.Background(), pdfText)

				if err == nil {

					document.HasEmbedding = true

					db.UpdateOneCustom("documents",
						bson.M{"_id": document.Uid},
						bson.M{"hasEmbedding": true},
					)

					AddDocumentEmbedding(document.Uid, embedding)
				} else {
					lg.LogE(embError.Error())
				}

			}

			if uploadParams.LLMDescription {
				text := llmCtrl.AskModelForDescription(c, user, uploadParams.DescriptionModel, pdfText)

				db.UpdateOneCustom("documents",
					bson.M{"_id": document.Uid},
					bson.M{"descriptions": text},
				)

				lg.LogOk("Uploaded document description")

			}

		}
	}

}

func ExtractPdfTextInMemory(data []byte) (string, error) {
	// Create a bytes.Reader which implements io.ReaderAt
	reader := bytes.NewReader(data)

	// Now you can use it with your PDF function
	pdfReader, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		lg.LogE(err.Error())
		return "", err
	}

	var buf bytes.Buffer

	b, err := pdfReader.GetPlainText()

	if err != nil {
		lg.LogE(err.Error())
		return "", err
	}
	buf.ReadFrom(b)
	content := buf.String()
	return content, nil
}

func CreateDocFileEmbedding(filePath string) error {
	pdf.DebugOn = true

	f, r, err := pdf.Open(filePath)

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
