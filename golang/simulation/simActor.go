package simulation

import (
	"fmt"
	"turtle/modelsApp"
)

type SimActor struct {
	Name string
	Id   int64
}

func (self *SimActor) FromActorDefinition(def *modelsApp.Actor) {
	def.Name = fmt.Sprintf("%s_%d", def.Name, self.Id)
}
