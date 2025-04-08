package models

type TurtleScene struct {
	Uid         string `json:"uid"`
	Name        string `json:"name"`
	At          int64  `json:"at"`
	CreatedBy   string `json:"created_by" bson:"created_by"`
	Org         string `json:"org"`
	Description string `json:"description"`
	Parent      string `json:"parent"`
	Type        string `json:"type"`
}

func (self *TurtleScene) GetUid() string {
	return self.Uid
}

func (self *TurtleScene) SetUid(uid string) {
	self.Uid = uid
}
