package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Uid       primitive.ObjectID
	Email     string
	Firstname string
	Surname   string
	Password  string
	Type      string
	Org       primitive.ObjectID
}

func NewUser() *User {
	tmp := User{}
	tmp.Firstname = ""
	tmp.Password = ""
	tmp.Surname = ""
	tmp.Type = "admin"
	return &tmp
}

func (self *User) FromAnotherUserNoPass(another *User) {
	self.Firstname = another.Firstname
	self.Surname = another.Surname
	self.Type = another.Type
	self.Uid = another.Uid
}
