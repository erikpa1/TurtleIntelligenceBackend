package simulation

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/lg"
	"turtle/tools"
)

type SpawnBehaviour struct {
	World  *SimWorld
	Entity SimEntity

	NextSpawnTime tools.Seconds
	SpawnInterval tools.Seconds
	SpawnLimit    int64

	SpawnActor  primitive.ObjectID
	ActiveActor *SimActor
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

func (self *SpawnBehaviour) Init1() {

}

func (self *SpawnBehaviour) Init2() {

}

func (self *SpawnBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = *entity
	self.SpawnInterval = entity.TypeData.GetSeconds("spawn_interval", 1)
	self.SpawnLimit = entity.TypeData.GetInt64("spawn_limit", 1)
	self.SpawnActor = entity.TypeData.GetPrimitiveObjectId("actor")
}

func (self *SpawnBehaviour) Step() {
	actualTime := self.World.Stepper.Now

	if self.ActiveActor == nil {

		lg.LogW(self.NextSpawnTime)

		if actualTime >= self.NextSpawnTime {
			self._Spawn()
		}
	} else {
		self._TryToAddActorNext()
	}

}

func (self *SpawnBehaviour) _Spawn() {
	self.ActiveActor = self.World.SpawnActorWithUid(self.SpawnActor)
}

func (self *SpawnBehaviour) _CalculateNextSpawn() {
	self.NextSpawnTime = self.World.Stepper.Now + self.SpawnInterval
}

func (self *SpawnBehaviour) _TryToAddActorNext() {

	connections, hasConnections := self.World.SimConnections[self.Entity.Uid]

	if hasConnections {
		for _, connection := range connections {

			taker, isITaker := connection.(ActorTakerBehaviour)

			if isITaker {
				taken := taker.TakeActor(self.ActiveActor)

				if taken {
					lg.LogE("Passing actor next")
					return
				}
			}
		}
	}

}
