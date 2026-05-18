package simulation2

import (
	"turtle/core/lgr"
	"turtle/simulation/simInternal"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BehSpawn struct {
	World  *SimWorld
	Entity *SimEntity

	SpawnInterval       tools.Seconds
	SpawnLimit          int
	SpawnMultiplication int
	SpawnActorUid       primitive.ObjectID

	ActiveActor   *SimActor
	SpawnedCount  int
	NextSpawnTime tools.Seconds
}

func GetBehSpawn(entity *SimEntity) *BehSpawn {
	return entity.Impl.(*BehSpawn)
}

func (self *BehSpawn) Step() {
	actualTime := self.World.Stepper.Now

	if self.ActiveActor == nil {
		if actualTime >= self.NextSpawnTime {
			self.Spawn()
		}
	} else {
		self._TryToPassActor()
	}
}

/*
Spawn tries to spawn entity, it have limit of max spawn entity
*/
func (self *BehSpawn) Spawn() {
	if self.CanSpawn() {
		self.ForceSpawn()
	}
}

/*
ForceSpawn spawns entity without limit
*/
func (self *BehSpawn) ForceSpawn() {

	for i := 1; i <= self.SpawnMultiplication; i++ {

		//CLaude

		self.ActiveActor = self.World.SpawnActorWithUid(self.SpawnActorUid)
		lgr.Ok("Spawned actor %v", self.ActiveActor)

		if self.ActiveActor != nil {
			self.SpawnedCount++
			self.ActiveActor.Position = self.Entity.Position
			self.NextSpawnTime = tools.MaxSeconds()
			self._TryToPassActor()
		}

	}

}

func (self *BehSpawn) _TryToPassActor() {

	connections, hasConnections := self.World.SimConnections[self.Entity.Uid]

	if hasConnections {
		for _, connection := range connections {
			taken := connection.TakeActor(self.ActiveActor)
			if taken {
				self.ActiveActor = nil
				self._CalculateNextSpawn()
				return
			}
		}
	}
}

func (self *BehSpawn) _CalculateNextSpawn() {
	self.NextSpawnTime = self.World.Stepper.Now + self.SpawnInterval

	self.World.CreateUpcomingEvent(simInternal.SimUpcomingEvent{
		Id:     self.Entity.RuntimeId,
		Type:   simInternal.UPC_EVNT_SPAWN,
		Second: self.NextSpawnTime,
	})

	lgr.Info(self.Entity.FormatInfo("Next spawn event at: (%d/%d)", self.NextSpawnTime, self.World.Stepper.End))
}

func (self *BehSpawn) CanSpawn() bool {
	if self.SpawnLimit < 0 {
		return true
	} else {
		return self.SpawnedCount < self.SpawnLimit
	}
}
