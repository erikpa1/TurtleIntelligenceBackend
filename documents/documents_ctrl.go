package documents

import (
	"bytes"
	"context"
	"fmt"
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/llm/llmCtrl"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VSearchResult struct {
	Similarity float32   `json:"similarity"`
	Doc        *Document `json:"doc"`
}

func ListVSearchDocuments(c context.Context, user *models.User, searchQuery string, threshold float32) ([]*VSearchResult, error) {

	resultList := make([]*VSearchResult, 0)

	queryEmbedding, err := llmCtrl.CreateStringEmbedding(c, searchQuery)

	if err != nil {
		return nil, err
	}

	for _, docEmbedding := range db.QueryEntities[DocumentEmbedding](CT_DOC_EMBEDDINGS,
		bson.M{"org": user.Org}) {

		similarity := docEmbedding.DescEmbedding.GetSimilarity(queryEmbedding)

		if similarity > threshold {
			resultList = append(resultList, &VSearchResult{
				Similarity: similarity,
				Doc:        GetDocument(user, docEmbedding.Uid),
			})
		}
		lg.LogI(similarity)

	}

	return resultList, nil
}

func GetDocument(user *models.User, docUid primitive.ObjectID) *Document {
	return db.QueryEntity[Document](CT_DOC, bson.M{
		"_id": docUid,
		"org": user.Org,
	})
}

func ListDocuments(user *models.User) []*Document {

	return db.QueryEntities[Document](CT_DOC, bson.M{
		"org": user.Org,
	})

}

func ListDocumentExtracts(user *models.User) []*DocumentExtraction {
	return db.QueryEntities[DocumentExtraction](CT_DOC_EXTRACT, user.FillOrgQuery(nil))

}

func DeleteDocument(user *models.User, documentUid primitive.ObjectID) {

	var docToDelete *Document

	if user.IsAdmin() {
		docToDelete = db.QueryEntity[Document](CT_DOC, bson.M{
			"_id": documentUid,
			"org": user.Org,
		})
	} else {
		docToDelete = db.QueryEntity[Document](CT_DOC, bson.M{
			"_id":  documentUid,
			"org":  user.Org,
			"user": user.Uid,
		})
	}

	if docToDelete != nil {

		db.DeleteEntities(CT_DOC_EMBEDDINGS, user.FillOrgQuery(bson.M{
			"_id": documentUid,
		}))

		db.DeleteEntities(CT_DOC_EXTRACT, user.FillOrgQuery(bson.M{
			"_id": documentUid,
		}))

		db.SC.DeleteFileNew(fmt.Sprintf("documents/%s.%s", documentUid, docToDelete.Extension))

		db.DeleteEntities(CT_DOC, user.FillOrgQuery(bson.M{
			"_id": documentUid,
		}))

	}

}

type InsertDocumentParams struct {
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	LLMDescription   bool               `json:"llmDescription"`
	CreateEmbedding  bool               `json:"createEmbedding"`
	DescriptionModel primitive.ObjectID `json:"descriptionModel"`
}

func UpdateDocument(user *models.User, document *Document) {
	document.Org = user.Org
	db.UpdateOneCustom(CT_DOC,
		bson.M{
			"_id": document.Uid,
			"org": user.Org,
		}, bson.M{"$set": document})

}
func CreateAndUploadDocument(c *gin.Context, user *models.User, uploadParams *InsertDocumentParams, documentData []byte) {

	document := &Document{}
	document.Uid = primitive.NewObjectID()
	document.Name = uploadParams.Name
	document.Description = uploadParams.Description
	document.Org = user.Org
	document.Extension = "pdf"

	db.InsertEntity(CT_DOC, document)

	fileName := fmt.Sprintf("%s.%s", document.Uid.Hex(), document.Extension)

	db.SC.UploadFile("documents", fileName, documentData)

	if uploadParams.CreateEmbedding || uploadParams.LLMDescription {

		pdfText, extractError := ExtractPdfTextInMemory(documentData)

		extraction := DocumentExtraction{}
		extraction.Uid = document.Uid
		extraction.Extraction = pdfText
		extraction.Org = user.Org

		db.InsertEntity(CT_DOC_EXTRACT, extraction)

		if extractError == nil {

			if uploadParams.LLMDescription {
				lg.LogI("Going to create LLM description")
				descText := llmCtrl.AskModelForDescription(c, user, uploadParams.DescriptionModel, pdfText, 100)
				lg.LogOk("LLM desc", descText)

				db.UpdateOneCustom(CT_DOC,
					bson.M{"_id": document.Uid},
					bson.M{"$set": bson.M{"description": descText}},
				)

				document.Description = descText

				lg.LogOk("Uploaded document description")

			}

			if uploadParams.CreateEmbedding {
				embedding, embError := llmCtrl.CreateStringEmbedding(context.Background(), pdfText)
				descEmbedding, _ := llmCtrl.CreateStringEmbedding(context.Background(), document.Description)

				if embError == nil {

					document.HasEmbedding = true

					db.UpdateOneCustom(CT_DOC,
						bson.M{"_id": document.Uid},
						bson.M{"$set": bson.M{"hasEmbedding": true}},
					)

					AddDocumentEmbedding(user.Org,
						document.Uid,
						embedding,
						descEmbedding,
					)
				} else {
					lg.LogE(embError.Error())
				}

			}

		} else {
			lg.LogE(extractError.Error())
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

func AddDocumentEmbedding(org primitive.ObjectID,
	documentUid primitive.ObjectID,
	embedding [][]float32,
	descEmbedding [][]float32,
) {
	emb := DocumentEmbedding{}
	emb.Uid = documentUid
	emb.Embedding = embedding
	emb.DescEmbedding = descEmbedding
	emb.Org = org
	db.InsertEntity(CT_DOC_EMBEDDINGS, emb)

}
