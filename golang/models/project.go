package models

type TurtleProject struct {
	Uid         string `json:"uid"`
	Name        string `json:"name"`
	At          int64  `json:"at"`
	CreatedBy   string `json:"created_by" bson:"created_by"`
	Org         string `json:"org"`
	Description string `json:"description"`
}

func (self *TurtleProject) GetUid() string {
	return self.Uid
}

func (self *TurtleProject) SetUid(uid string) {
	self.Uid = uid
}
