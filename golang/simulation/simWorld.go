package simulation

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/ctrlApp"
	"turtle/lg"
	"turtle/modelsApp"
	"turtle/server"
	"turtle/tools"
)

type SimWorld struct {
	Uid  primitive.ObjectID
	Name string

	SimEntities    map[primitive.ObjectID]*SimEntity
	SimActors      map[primitive.ObjectID]*SimActor
	SimConnections map[primitive.ObjectID][]*SimEntity

	ActorsDefinitions map[primitive.ObjectID]*modelsApp.Actor

	Stepper    SimStepper
	IsOnline   bool
	IdsCounter int64
}

func NewSimWorld() *SimWorld {
	tmp := &SimWorld{}
	tmp.Stepper.End = 100
	tmp.IsOnline = true
	return tmp
}

func (self *SimWorld) LoadEntities(entities []*modelsApp.Entity) {
	for _, entity := range entities {
		simEntity := SimEntity{}
		simEntity.FromEntity(entity)

		self.SimEntities[entity.Uid] = &simEntity

	}
}

func (self *SimWorld) LoadConnections(connections []*modelsApp.EntityConnection) {
	for _, connection := range connections {

		array, exists := self.SimConnections[connection.A]

		if exists == false {
			array = []*SimEntity{}
			self.SimConnections[connection.A] = array
		}

		simEntity, found := self.SimEntities[connection.B]

		if found {
			self.SimConnections[connection.A] = append(array, simEntity)
		} else {
			lg.LogE("Unable to find entity [%s] in world", connection.B.Hex())
		}
	}
}

func (self *SimWorld) RunSimulation() {

	var second tools.Seconds = 0

	for second = 0; second < self.Stepper.End; second++ {
		self.Step()
		self.Stepper.Step()

		server.MYIO.EmitSync("simstep", bson.M{
			"second": second,
		})
	}

}

func (self *SimWorld) GetConnectionsOf(entity primitive.ObjectID) []*SimEntity {

	entities, exists := self.SimConnections[entity]

	if exists {
		return entities
	} else {
		return []*SimEntity{}
	}

}

func (self *SimWorld) Step() {

	lg.LogI(fmt.Sprintf("Step (%d/%d)", self.Stepper.Now, self.Stepper.End))

}

func (self *SimWorld) SpawnActorWithUid(uid primitive.ObjectID) *SimActor {

	definition, exists := self.ActorsDefinitions[uid]

	if exists == false {
		tmp := ctrlApp.GetActor(uid)

		if tmp != nil {
			definition = tmp
			self.ActorsDefinitions[uid] = definition
		} else {
			lg.LogE("SimActor definition [%s] not found", uid.Hex())
			return nil
		}
	}

	actor := SimActor{}
	actor.Id = self.IdsCounter
	self.IdsCounter += 1
	actor.FromActorDefinition(definition)

	return &actor

}
