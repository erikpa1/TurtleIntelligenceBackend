package llm

type OllamaNode struct {
	OllamaUrl string `json:"ollamaUrl" bson:"ollamaUrl"`
	ModelName string `json:"modelName" bson:"modelName"`
}
