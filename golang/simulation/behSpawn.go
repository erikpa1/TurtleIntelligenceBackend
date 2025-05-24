package simulation

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpawnBehaviour struct {
	World      *SimWorld
	SpawnLimit int64 `mapstructure:"spawn_limit"`
}

func NewSpawnBehaviour() *SpawnBehaviour {
	return &SpawnBehaviour{}
}

func (self *SpawnBehaviour) SpawnActorOfType(actorUid primitive.ObjectID) {
	self.World.SpawnActorWithUid(actorUid)
}

// Implementacia IBehaviour
func (self *SpawnBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *SpawnBehaviour) Step() {

}
