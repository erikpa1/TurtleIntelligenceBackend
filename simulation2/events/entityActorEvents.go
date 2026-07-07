package events

import "turtle/simulation2/entities"

const ACTOR_TAKEN = "ActorTaken"
const ACTOR_PASSED = "ActorPassed"

type ActorTakenStruct struct {
	Entity *entities.SimEntity
	Actor  *entities.SimActor
}
