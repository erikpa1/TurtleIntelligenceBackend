package knowledge

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/llm/llmModels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KnowledgeType int8

const (
	KNOW_TYPE_PLAINTEXT KnowledgeType = 0
	KNOW_TYPE_DOCUMENT  KnowledgeType = 1
	KNOW_TYPE_GUIDANCE  KnowledgeType = 2
)

type Knowledge struct {
	Uid          primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org          primitive.ObjectID `json:"org"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Type         KnowledgeType      `json:"type"`
	HasEmbedding bool               `json:"hasEmbedding" bson:"hasEmbedding"`
	TypeData     bson.M             `json:"typeData" bson:"typeData"`
}

type KnowledgeEmbedding struct {
	Uid       primitive.ObjectID  `json:"uid" bson:"_id,omitempty"`
	Org       primitive.ObjectID  `json:"org"`
	Embedding llmModels.Embedding `json:"embedding"`
}

type KnowledgePlainTextTypeData bson.M

func (self *KnowledgePlainTextTypeData) GetEmbeddableString() string {

	tmp, okTmp := (*self)["text"].(string)

	if okTmp {
		return tmp
	} else {
		return ""
	}

}
