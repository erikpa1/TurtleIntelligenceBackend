package agents

type LLMAgentData struct {
	SystemPrompt string `json:"systemPrompt" bson:"systemPrompt"`
}

type OllamaNode struct {
	OllamaUrl string `json:"ollamaUrl" bson:"ollamaUrl"`
	ModelName string `json:"modelName" bson:"modelName"`
}

type WriteToFileNode struct {
	ParentFolder string `json:"parentFolder" bson:"parentFolder"`
	FileName     string `json:"fileName" bson:"fileName"`
	OpenFolder   bool   `json:"openFolder" bson:"openFolder"`
	UseWd        bool   `json:"useWd" bson:"useWd"`
}

func (self *WriteToFileNode) GetFileName() string {
	if self.FileName == "" {
		return "output.txt"
	}
	return self.FileName
}
