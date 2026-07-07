package behBuffer

import (
	"turtle/core/lgr"
	"turtle/simulation2/entities"
	"turtle/simulation2/events"
)

func (self *BehBuffer) EmitActorTaken(actor *entities.SimActor) {

	lgr.Error("Emitting %s", events.ACTOR_TAKEN)

	self.Entity.Aee.Emit(events.ACTOR_TAKEN, events.ActorTakenStruct{
		Actor:  actor,
		Entity: self.Entity,
	})
}

func (self *BehBuffer) EmitActorLeft(actor *entities.SimActor) {
	self.Entity.Aee.Emit(events.ACTOR_PASSED, events.ActorTakenStruct{
		Actor:  actor,
		Entity: self.Entity,
	})
}
