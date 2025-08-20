package simulation

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/simulation/simInternal"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func NewSpawnBehaviour() ISimBehaviour {
	return &SpawnBehaviour{}
}

// Implementacia IBehaviour
func (self *SpawnBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *SpawnBehaviour) Init1() {
	self._CalculateNextSpawn()
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

		if actualTime >= self.NextSpawnTime {
			self._Spawn()
		} else {
			lg.LogW(self.NextSpawnTime)
		}
	} else {
		self._TryToAddActorNext()
	}

}

func (self *SpawnBehaviour) __Interface() {
	var _ ISimBehaviour = &SpawnBehaviour{}

	//TODO mozno dorobit providera?
}

func (self *SpawnBehaviour) _Spawn() {
	self.ActiveActor = self.World.SpawnActorWithUid(self.SpawnActor)

	if self.ActiveActor != nil {
		self.ActiveActor.Position = self.Entity.Position
	}

}

func (self *SpawnBehaviour) _CalculateNextSpawn() {
	self.NextSpawnTime = self.World.Stepper.Now + self.SpawnInterval

	self.World.CreateUpcomingEvent(0, simInternal.SimUpcomingEvent{
		simInternal.UPC_EVNT_SPAWN,
		self.NextSpawnTime,
	})
}

func (self *SpawnBehaviour) _TryToAddActorNext() {

	connections, hasConnections := self.World.SimConnections[self.Entity.Uid]

	if hasConnections {
		for _, connection := range connections {

			taker, isITaker := connection.(ActorTakerBehaviour)

			if isITaker {
				taken := taker.TakeActor(self.ActiveActor)

				if taken {
					self.ActiveActor = nil
					self._CalculateNextSpawn()
					return
				}
			}
		}
	}

}
