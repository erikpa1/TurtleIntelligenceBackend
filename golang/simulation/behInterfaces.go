package simulation

type ISimBehaviour interface {
	SetWorld(world *SimWorld)
	Step()
}

type ActorTakerBehaviour interface {
	TakeActor(actor *SimActor)
	CanTakeActor(actor *SimActor) bool
}

type ActorProviderBehaviour interface {
	PopActor() *SimActor
	HasAnyActor() bool
	HasActorOfType(actorType string) bool
	HasActorWithVariable(variable string, value any) bool
}
