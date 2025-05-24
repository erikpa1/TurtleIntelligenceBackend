package simulation

type ProcessBehaviour struct {
	World *SimWorld
}

func NewProcessBehaviour() *ProcessBehaviour {
	return &ProcessBehaviour{}
}

// IBehaviour implementation
func (self *ProcessBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *ProcessBehaviour) Step() {

}

// Taker behaviour
func (self *ProcessBehaviour) TakeActor(actor *SimActor) {

}

func (self *ProcessBehaviour) CanTakeActor(actor *SimActor) bool {
	return false
}
