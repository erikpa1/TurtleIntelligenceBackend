package simulation

import (
	"fmt"

	"github.com/erikpa1/TurtleIntelligenceBackend/ctrlApp"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/modelsApp"
	"github.com/erikpa1/TurtleIntelligenceBackend/simulation/simInternal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SimWorld struct {
	Uid  primitive.ObjectID
	Name string

	SimBehaviours  map[primitive.ObjectID]ISimBehaviour
	SimActors      map[int64]*SimActor
	SimConnections map[primitive.ObjectID][]ISimBehaviour

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

	tmp.SimBehaviours = make(map[primitive.ObjectID]ISimBehaviour)
	tmp.SimActors = make(map[int64]*SimActor)
	tmp.SimConnections = make(map[primitive.ObjectID][]ISimBehaviour)
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
func (self *SimWorld) LoadEntities(entities []*modelsApp.Entity) {
	for _, entity := range entities {
		simEntity := SimEntity{}
		simEntity.FromEntity(entity)
		simEntity.RuntimeId = self.GetRuntimeId()

		entityType := entity.Type

		constructor, ok := BEH_FACTORY.Behaviours[entityType]

		if ok == false {
			constructor = NewUndefinedBehaviour
		}

		behaviour := constructor()

		behaviour.SetWorld(self)
		behaviour.SetEntity(&simEntity)

		self.SimBehaviours[entity.Uid] = behaviour

	}
}

func (self *SimWorld) LoadConnections(connections []*modelsApp.EntityConnection) {
	for _, connection := range connections {

		array, exists := self.SimConnections[connection.A]

		if exists == false {
			array = []ISimBehaviour{}
			self.SimConnections[connection.A] = array
		}

		simBehaviour, found := self.SimBehaviours[connection.B]

		if found {
			self.SimConnections[connection.A] = append(array, simBehaviour)
		} else {
			lg.LogE("Unable to find entity [%s] in world", connection.B.Hex())
		}
	}
}

func (self *SimWorld) PrepareSimulation() {

	for _, behaviour := range self.SimBehaviours {
		behaviour.Init1()
	}

	for _, behaviour := range self.SimBehaviours {
		behaviour.Init2()
	}

}

func (self *SimWorld) GetConnectionsOf(entity primitive.ObjectID) []ISimBehaviour {

	entities, exists := self.SimConnections[entity]

	if exists {
		return entities
	} else {
		return []ISimBehaviour{}
	}
}

func (self *SimWorld) ClearStates() {
	self.StatesCreatedActors = make(map[int64]*SimActor, 0)
	self.StatesDestroyedActors = make([]int64, 0)
	self.StatesUpdates = make(map[int64]bson.M)
	self.StatesUpcomingEvents = make([]simInternal.SimUpcomingEvent, 0)
}

func (self *SimWorld) Step() {
	lg.LogI(fmt.Sprintf("Step (%d/%d)", self.Stepper.Now, self.Stepper.End))

	for _, behaviour := range self.SimBehaviours {
		behaviour.Step()
	}

}

func (self *SimWorld) UnspawnActor(actor *SimActor) {
	actor.World = nil //Make GC life easier
	self.StatesDestroyedActors = append(self.StatesDestroyedActors, actor.Id)
	delete(self.SimActors, actor.Id)
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

	for _, behaiour := range self.SimBehaviours {
		entity := behaiour.GetEntity()
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
