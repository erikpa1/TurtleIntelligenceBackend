package simulation2

import (
	"fmt"
	"turtle/core/lgr"
	"turtle/simulation/simMath"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"turtle/modelsApp"
	"turtle/tools"
)

type SimEntity struct {
	RuntimeId int64
	Uid       primitive.ObjectID
	Name      string
	Type      string
	Position  simMath.Position
	TypeData  tools.SafeJson
	World     *SimWorld
	Functions map[string]any
	Data      SimBehData
	Impl      any
}

func (self *SimEntity) FromEntity(def *modelsApp.WorldEntity) {
	self.Uid = def.Uid
	self.Name = def.Name
	self.Type = def.Type
	self.Position = def.Position
	self.TypeData = def.TypeData

}

func NewSimEntity() *SimEntity {
	return &SimEntity{
		Functions: make(map[string]any),
		TypeData:  tools.NewSafeJson(),
	}
}

func (self *SimEntity) FormatInfo(format string, args ...interface{}) string {
	return fmt.Sprintf("[%s] %s", self.Type, fmt.Sprintf(format, args...))
}

func CastImplementation[T any](imp any) *T {
	tmp, ok := imp.(*T)
	if ok {
		return tmp
	}

	var fullType T
	lgr.Error("Failed to cast [%T] to [%T]", imp, fullType)
	return nil
}
