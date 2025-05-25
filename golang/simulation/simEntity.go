package simulation

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/modelsApp"
	"turtle/tools"
)

type SimEntity struct {
	Uid      primitive.ObjectID
	Name     string
	Type     string
	TypeData *tools.SafeJson
	Wold     *SimWorld
}

func (self *SimEntity) FromEntity(def *modelsApp.Entity) {
	self.Uid = def.Uid
	self.Name = def.Name
	self.Type = def.Type
	self.TypeData = tools.NewSafeJson()
	self.TypeData.Data = def.TypeData

}
