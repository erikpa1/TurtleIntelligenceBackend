package simulation2

import (
	"turtle/core/lgr"
	"turtle/simulation2/simInternal"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpawnBehaviour = SimEntity

func (self SpawnBehaviour) GetSpawnLimit() int {
	return self.BehaviourEntity.GetInt("SpawnLimit")
}

func (self SpawnBehaviour) GetSpawnInterval() tools.Seconds {
	return self.BehaviourEntity.GetSeconds("SpawnLimit")
}

func (self SpawnBehaviour) GetNextSpawnTime() tools.Seconds {
	return self.BehaviourEntity.GetSeconds("NextSpawnTime")
}

func (self SpawnBehaviour) GetActorUidToSpawn() primitive.ObjectID {
	return self.BehaviourEntity.GetPrimitiveObjectId("SpawnActor")
}

func (self SpawnBehaviour) GetActiveActor() *SimActor {
	return self.BehaviourEntity.GetActor("ActiveActor")
}

func (self SpawnBehaviour) Step() {
	actualTime := self.World.Stepper.Now

	if self.GetActiveActor() == nil {

		if actualTime >= self.GetNextSpawnTime() {
			self.Spawn()
		} else {
			lgr.Info("%d", self.GetNextSpawnTime())
		}
	} else {
		self._TryToAddActorNext()
	}

}

func (self SpawnBehaviour) SetActiveActor(actor *SimActor) *SimActor {
	return self.BehaviourEntity.SetActor("ActiveActor", actor)
}

func (self SpawnBehaviour) Spawn() {
	actor := self.SetActiveActor(self.World.SpawnActorWithUid(self.GetActorUidToSpawn()))

	if actor != nil {
		actor.Position = self.Position
	}
}

func (self SpawnBehaviour) _TryToAddActorNext() {

	connections, hasConnections := self.World.SimConnections[self.Uid]

	if hasConnections {
		for _, connection := range connections {

			taken := connection.TakeActor(self.GetActiveActor())

			if taken {
				self.SetActiveActor(nil)
				self._CalculateNextSpawn()
				return
			}
		}
	}

}

func (self SpawnBehaviour) _CalculateNextSpawn() {
	nextSpawnTime := self.BehaviourEntity.SetSeconds("NextSpawnTime", self.World.Stepper.Now+self.GetSpawnInterval())

	self.World.CreateUpcomingEvent(simInternal.SimUpcomingEvent{
		Id:     self.RuntimeId,
		Type:   simInternal.UPC_EVNT_SPAWN,
		Second: nextSpawnTime,
	})

	lgr.Info("Next spawn event at, %d", nextSpawnTime)
}
