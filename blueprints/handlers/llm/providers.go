package llm

type StaticMemory struct {
	MemoryText string `json:"memoryText" bson:"memoryText"`
}
