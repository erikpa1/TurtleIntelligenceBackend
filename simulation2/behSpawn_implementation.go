package simulation2

import (
	"turtle/core/lgr"
	"turtle/simulation/simInternal"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpawnBehaviour = SimEntity

const spwn_NextSpawnTime = "_nextSpawnTime"
const spwn_SpawnedCount = "_spawnedCount"

func (self SpawnBehaviour) GetSpawnLimit() int {
	return self.BehaviourEntity.GetInt("SpawnLimit")
}

func (self SpawnBehaviour) IncSpawnCount() int {
	tmp := self.BehaviourEntity.GetInt(spwn_SpawnedCount) + 1
	self.BehaviourEntity.Set(spwn_SpawnedCount, tmp)
	return tmp
}

func (self SpawnBehaviour) GetSpawnInterval() tools.Seconds {
	return self.BehaviourEntity.GetSeconds("SpawnLimit")
}

func (self SpawnBehaviour) GetNextSpawnTime() tools.Seconds {
	tmp := self.BehaviourEntity.GetSeconds(spwn_NextSpawnTime)
	return tmp
}

func (self SpawnBehaviour) GetActorUidToSpawn() primitive.ObjectID {
	return self.BehaviourEntity.GetPrimitiveObjectId("SpawnActor")
}

func (self SpawnBehaviour) GetActiveActor() *SimActor {
	return self.BehaviourEntity.GetActor("ActiveActor")
}

func (self SpawnBehaviour) Step() {
	actualTime := self.World.Stepper.Now

	actor := self.GetActiveActor()

	if actor == nil {
		if actualTime >= self.GetNextSpawnTime() {
			self.Spawn()
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

	lgr.Error("Spawned actor %v", actor)

	if actor != nil {
		self.IncSpawnCount()
		actor.Position = self.Position
		self._TryToAddActorNext()
	}

}

func (self SpawnBehaviour) _TryToAddActorNext() {

	lgr.Error("Heeere")

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
	nextSpawnTime := self.BehaviourEntity.SetSeconds(spwn_NextSpawnTime, self.World.Stepper.Now+self.GetSpawnInterval())

	self.World.CreateUpcomingEvent(simInternal.SimUpcomingEvent{
		Id:     self.RuntimeId,
		Type:   simInternal.UPC_EVNT_SPAWN,
		Second: nextSpawnTime,
	})

	lgr.Info(self.FormatInfo("Next spawn event at: (%d/%d)", nextSpawnTime, self.World.Stepper.End))

}

//Implementaiton

func _SpawnInit1(self *SpawnBehaviour) {
	self._CalculateNextSpawn()
}

func _SpawnStep(self *SpawnBehaviour) {
	self.Step()
}

func (self *SpawnBehaviour) _Spawn() {
	actor := self.SetActiveActor(self.World.SpawnActorWithUid(self.GetActorUidToSpawn()))

	if actor != nil {
		actor.Position = self.Position
	}

}
