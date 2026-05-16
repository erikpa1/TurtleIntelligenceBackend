package simulation2

import (
	"fmt"

	"turtle/core/lgr"
	"turtle/ctrlApp"
	"turtle/modelsApp"
	"turtle/simulation/simInternal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SimWorld struct {
	Uid  primitive.ObjectID
	Name string

	SimEntities    map[primitive.ObjectID]*SimEntity
	SimActors      map[int64]*SimActor
	SimConnections map[primitive.ObjectID][]*SimEntity

	ActorsDefinitions map[primitive.ObjectID]*modelsApp.Actor

	Stepper  SimStepper
	IsOnline bool

	StatesCreatedActors   map[int64]*SimActor
	StatesDestroyedActors []int64
	StatesUpdates         map[int64]bson.M
	StatesUpcomingEvents  []simInternal.SimUpcomingEvent

	ActorsIds int64
	RuntimeId int64
}

func NewSimWorld() *SimWorld {
	tmp := &SimWorld{}
	tmp.Stepper.End = 100
	tmp.IsOnline = true

	tmp.SimEntities = make(map[primitive.ObjectID]*SimEntity)
	tmp.SimActors = make(map[int64]*SimActor)
	tmp.SimConnections = make(map[primitive.ObjectID][]*SimEntity)
	tmp.ActorsDefinitions = make(map[primitive.ObjectID]*modelsApp.Actor)

	tmp.StatesCreatedActors = make(map[int64]*SimActor, 0)
	tmp.StatesDestroyedActors = make([]int64, 0)
	tmp.StatesUpdates = make(map[int64]bson.M)
	tmp.StatesUpcomingEvents = make([]simInternal.SimUpcomingEvent, 0)

	return tmp
}

func (self *SimWorld) GetRuntimeId() int64 {
	tmp := self.RuntimeId
	self.RuntimeId += 1
	return tmp

}

func (self *SimWorld) LoadEntities(entities []*modelsApp.WorldEntity) {
	for _, entity := range entities {

		simEntity := NewSimEntity()
		simEntity.FromEntity(entity)
		simEntity.RuntimeId = self.GetRuntimeId()
		simEntity.World = self

		entityType := entity.Type

		constructor, ok := BEH_FACTORY.Behaviours[entityType]

		if ok {
			constructor(simEntity)
		} else {
			lgr.Error("No constructor found for entity type %s", entityType)
		}

		self.SimEntities[entity.Uid] = simEntity

	}
}

func (self *SimWorld) LoadConnections(connections []*modelsApp.EntityConnection) {
	for _, connection := range connections {

		array, exists := self.SimConnections[connection.A]

		if exists == false {
			array = []*SimEntity{}
			self.SimConnections[connection.A] = array
		}

		simBehaviour, found := self.SimEntities[connection.B]

		if found {
			self.SimConnections[connection.A] = append(array, simBehaviour)
		} else {
			lgr.Error("Unable to find entity [%s] in world", connection.B.Hex())
		}
	}
}

func (self *SimWorld) PrepareSimulation() {

	stages := []int{1, 2}

	for _, stage := range stages {
		for _, entity := range self.SimEntities {
			fnName := fmt.Sprintf("%s%d", FN_INIT_BASE, stage)
			fnInit, haveInit := GetSimFunction[FnInit](entity, fnName)

			if haveInit {
				fnInit(entity)
			} else {
				//lgr.Error("%s, have no fn %s", entity.type, fnName)
			}
		}
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

func (self *SimWorld) ClearStates() {
	self.StatesCreatedActors = make(map[int64]*SimActor, 0)
	self.StatesDestroyedActors = make([]int64, 0)
	self.StatesUpdates = make(map[int64]bson.M)
	self.StatesUpcomingEvents = make([]simInternal.SimUpcomingEvent, 0)
}

func (self *SimWorld) Step() {
	lgr.Info("Step [%d/%d]", self.Stepper.Now, self.Stepper.End)

	for _, entity := range self.SimEntities {
		fnStep, haveStep := GetSimFunction[FnStep](entity, FN_STEP)

		if haveStep {
			fnStep(entity)
		}
	}

}

func (self *SimWorld) UnspawnActor(actor *SimActor) {
	actor.World = nil //Make GC life easier
	self.StatesDestroyedActors = append(self.StatesDestroyedActors, actor.Id)
	delete(self.SimActors, actor.Id)
}

func (self *SimWorld) SpwanCustomActor(actor *SimActor) {
	actor.Id = self.ActorsIds
	actor.World = self

	self.ActorsIds += 1
	self.StatesCreatedActors[actor.Id] = actor

}

func (self *SimWorld) SpawnActorWithUid(uid primitive.ObjectID) *SimActor {

	definition, exists := self.ActorsDefinitions[uid]

	if exists == false {
		tmp := ctrlApp.GetActor(uid)

		if tmp != nil {
			definition = tmp
			self.ActorsDefinitions[uid] = definition
		} else {
			lgr.ErrorStack("SimActor definition [%s] not found", uid.Hex())
			return nil
		}
	}

	actor := NewSimActor()
	actor.Id = self.ActorsIds
	actor.World = self

	self.ActorsIds += 1
	actor.FromActorDefinition(definition)

	self.StatesCreatedActors[actor.Id] = actor

	return actor

}

func (self *SimWorld) UpdateActorState(key int64, stateKey string, value any) {

	_, inSpawned := self.StatesCreatedActors[key]

	if inSpawned == false {

		stateSetter, bsonExists := self.StatesUpdates[key]

		if bsonExists {
			stateSetter[stateKey] = value

		} else {
			self.StatesUpdates[key] = bson.M{
				stateKey: value}
		}

	}
}

func (self *SimWorld) CreateUpcomingEvent(value simInternal.SimUpcomingEvent) {
	self.StatesUpcomingEvents = append(self.StatesUpcomingEvents, value)
}

func (self *SimWorld) ToJsonInit() bson.M {

	entitiesRuntime := map[string]int64{}

	for _, entity := range self.SimEntities {
		entitiesRuntime[entity.Uid.Hex()] = entity.RuntimeId
	}

	return bson.M{
		"runtimeIds": entitiesRuntime,
	}
}

func (self *SimWorld) ToJsonClient() bson.M {

	return bson.M{
		"second":    self.Stepper.Now,
		"spawned":   self.StatesCreatedActors,
		"unspawned": self.StatesDestroyedActors,
		"states":    self.StatesUpdates,
		"events":    self.StatesUpcomingEvents,
	}
}
