package simulation2

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BehBuffer struct {
	World  *SimWorld
	Entity *SimEntity

	Actors []*SimActor

	Capacity     int64
	InitialCount int64
	InitialActor primitive.ObjectID
}

func GetBehBuffer(entity *SimEntity) *BehBuffer {
	return entity.Impl.(*BehBuffer)
}

func (self *BehBuffer) Step() {
	self._TryToPassActorsNext()
}

func (self *BehBuffer) _TryToPassActorsNext() {
	for _, conn := range self.World.GetConnectionsOf(self.Entity.Uid) {
		for len(self.Actors) > 0 {
			lastOne := self.Actors[len(self.Actors)-1]
			if !conn.TakeActor(lastOne) {
				break
			}
			self.PopActor()
		}
		if len(self.Actors) == 0 {
			return
		}
	}
}

//Entity taker behavour

func (self *BehBuffer) TakeActor(actor *SimActor) bool {
	canTake := self.CanTakeActor(actor)

	if canTake {
		self.Actors = append(self.Actors, actor)
		actor.UpdatePosition(self.Entity.Position.RandomizeXZ(1))
	}

	return canTake
}

func (self *BehBuffer) CanTakeActor(actor *SimActor) bool {
	if self.Capacity == -1 {
		return true
	}
	return len(self.Actors) < int(self.Capacity)
}

// Entity provideder behaviour
func (self *BehBuffer) PopActor() *SimActor {
	return self.PopBack()
}

func (self *BehBuffer) PopBack() *SimActor {
	if len(self.Actors) == 0 {
		return nil // or handle empty slice case
	}

	// Get the last element
	lastActor := self.Actors[len(self.Actors)-1]

	// Remove the last element by reslicing
	self.Actors = self.Actors[:len(self.Actors)-1]

	return lastActor
}
