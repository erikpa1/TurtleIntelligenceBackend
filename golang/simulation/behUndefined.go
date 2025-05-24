package simulation

type UndefinedBehaviour struct {
	World *SimWorld
}

func NewUndefinedBehaviour() *UndefinedBehaviour {
	return &UndefinedBehaviour{}
}

// Implementacia IBehaviour
func (self *UndefinedBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *UndefinedBehaviour) Step() {

}
