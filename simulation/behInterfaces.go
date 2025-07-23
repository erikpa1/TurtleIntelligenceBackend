package simulation

type ISimBehaviour interface {
	SetWorld(world *SimWorld)
	SetEntity(entity *SimEntity)
	Step()
	Init1()
	Init2()
}

type ActorTakerBehaviour interface {
	TakeActor(actor *SimActor) bool
	CanTakeActor(actor *SimActor) bool
}

type ActorProviderBehaviour interface {
	PopActor() *SimActor
	HasAnyActor() bool
	HasActorOfType(actorType string) bool
	HasActorWithVariable(variable string, value any) bool
}
