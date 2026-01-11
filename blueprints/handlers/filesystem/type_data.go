package filesystem

type LoadFileStringNode struct {
	FilePath string `json:"filePath" bson:"filePath"`
	UseWd    bool   `json:"useWd" bson:"useWd"`
}
