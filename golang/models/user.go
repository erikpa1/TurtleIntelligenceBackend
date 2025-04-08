package models

type User struct {
	Uid       string `json:"uid" bson:"uid"`
	Email     string `json:"email" bson:"email"`
	Firstname string `json:"firstname" bson:"firstname"`
	Surname   string `json:"surname" bson:"surname"`
	Password  string `json:"password" bson:"password"`
	Type      string `json:"type" bson:"type"`
	Org       string `json:"org" bson:"org"`
}

func NewUser() *User {
	tmp := User{}
	tmp.Firstname = ""
	tmp.Password = ""
	tmp.Surname = ""
	tmp.Type = "admin"
	tmp.Uid = "inmemoryuser"
	return &tmp
}

// Interface Implementetion
func (self *User) GetUid() string {
	return self.Uid
}

func (self *User) SetUid(uid string) {
	self.Uid = uid
}

func (self *User) FromAnotherUserNoPass(another *User) {
	self.Firstname = another.Firstname
	self.Surname = another.Surname
	self.Type = another.Type
	self.Uid = another.Uid
}
