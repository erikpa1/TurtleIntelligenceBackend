package simulation

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SinkBehaviour struct {
	World  *SimWorld
	Entity *SimEntity

	NextSpawnTime tools.Seconds
	SpawnInterval tools.Seconds
	SpawnLimit    int64

	SpawnActor  primitive.ObjectID
	ActiveActor *SimActor
}

func NewSinkBehaviour() ISimBehaviour {
	return &SinkBehaviour{}
}

// Implementacia IBehaviour
func (self *SinkBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *SinkBehaviour) Init1() {

}

func (self *SinkBehaviour) Init2() {

}

func (self *SinkBehaviour) GetEntity() *SimEntity {
	return self.Entity
}

func (self *SinkBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = entity
}

func (self *SinkBehaviour) Step() {

}

func (self *SinkBehaviour) TakeActor(actor *SimActor) bool {
	self.World.UnspawnActor(actor)
	return true
}

func (self *SinkBehaviour) CanTakeActor(actor *SimActor) bool {
	return true
}
