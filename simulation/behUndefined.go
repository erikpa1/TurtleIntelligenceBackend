package simulation

type UndefinedBehaviour struct {
	World  *SimWorld
	Entity SimEntity
}

func NewUndefinedBehaviour() ISimBehaviour {
	return &UndefinedBehaviour{}
}

// Implementacia IBehaviour
func (self *UndefinedBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *UndefinedBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = *entity
}

func (self *UndefinedBehaviour) Init1() {

}

func (self *UndefinedBehaviour) Init2() {

}

func (self *UndefinedBehaviour) Step() {

}
