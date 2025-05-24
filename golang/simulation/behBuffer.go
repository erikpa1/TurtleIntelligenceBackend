package simulation

type BufferBehaviour struct {
	World  *SimWorld
	Entity SimEntity
}

func NewBufferBehaviour() *BufferBehaviour {
	return &BufferBehaviour{}
}

// Implementacia IBehaviour

func (self *BufferBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *BufferBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = *entity
}

func (self *BufferBehaviour) Init1() {

}

func (self *BufferBehaviour) Init2() {

}

func (self *BufferBehaviour) Step() {

}

//Entity taker behavour

func (self *BufferBehaviour) TakeActor(actor *SimActor) {

}

func (self *BufferBehaviour) CanTakeActor(actor *SimActor) bool {
	return false
}

// Entity provideder behaviour
func (self *BufferBehaviour) PopActor() *SimActor {
	return nil
}

func (self *BufferBehaviour) HasAnyActor() bool {
	return false
}

func (self *BufferBehaviour) HasActorOfType(actorType string) bool {
	return false
}

func (self *BufferBehaviour) HasActorWithVariable(variable string, value any) bool {
	return true
}
