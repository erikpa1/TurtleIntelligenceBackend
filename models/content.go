package models

type TurtleContent struct {
	Uid    string `json:"uid"`
	Name   string `json:"name"`
	At     int64  `json:"at"`
	Org    string `json:"org"`
	Scene  string `json:"scene"`
	Parent string `json:"parent"`

	TypeData map[string]string `json:"type_data" bson:"type_data"`
}

func (self *TurtleContent) GetUid() string {
	return self.Uid
}

func (self *TurtleContent) SetUid(uid string) {
	self.Uid = uid
}
