package filesystem

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

type LoadFileStringNode struct {
	FilePath string `json:"filePath" bson:"filePath"`
	UseWd    bool   `json:"useWd" bson:"useWd"`
}
