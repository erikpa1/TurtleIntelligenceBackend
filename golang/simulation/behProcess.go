package simulation

type ProcessBehaviour struct {
	World  *SimWorld
	Entity SimEntity

	ActiveActor *SimActor
}

func NewProcessBehaviour() *ProcessBehaviour {
	return &ProcessBehaviour{}
}

// IBehaviour implementation
func (self *ProcessBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *ProcessBehaviour) Init1() {

}

func (self *ProcessBehaviour) Init2() {

}

func (self *ProcessBehaviour) Step() {

}

func (self *ProcessBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = *entity
}

// Taker behaviour
func (self *ProcessBehaviour) TakeActor(actor *SimActor) bool {
	canTake := self.CanTakeActor(actor)

	if canTake {
		self.ActiveActor = actor
		actor.Position = self.Entity.Position
	}

	return canTake
}

func (self *ProcessBehaviour) CanTakeActor(actor *SimActor) bool {
	return self.ActiveActor == nil
}
