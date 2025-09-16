package simulation

import "go.mongodb.org/mongo-driver/bson/primitive"

type BufferBehaviour struct {
	World  *SimWorld
	Entity *SimEntity

	Actors []*SimActor

	Capacity     int64
	InitialCount int64
	InitialActor primitive.ObjectID
}

func NewBufferBehaviour() ISimBehaviour {
	return &BufferBehaviour{
		Actors: make([]*SimActor, 0),
	}
}

func (self *BufferBehaviour) __Interface() {
	var _ ISimBehaviour = &BufferBehaviour{}
	var _ ActorTakerBehaviour = &BufferBehaviour{}
	var _ ActorProviderBehaviour = &BufferBehaviour{}
}

// Implementacia IBehaviour

func (self *BufferBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *BufferBehaviour) GetEntity() *SimEntity {
	return self.Entity
}

func (self *BufferBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = entity

	self.Capacity = entity.TypeData.GetInt64("capacity", 8)
	self.InitialActor = entity.TypeData.GetPrimitiveObjectId("initial_actor")
	self.InitialCount = entity.TypeData.GetInt64("initial_count", 8)

}

func (self *BufferBehaviour) Init1() {

}

func (self *BufferBehaviour) Init2() {

}

func (self *BufferBehaviour) Step() {
	self._TryToPassActorsNext()
}

func (self *BufferBehaviour) _TryToPassActorsNext() {

	if len(self.Actors) > 0 {

		lastOne := self.Actors[len(self.Actors)-1]

		for _, conn := range self.World.GetConnectionsOf(self.Entity.Uid) {

			iActorTaker, ok := conn.(ActorTakerBehaviour)

			if ok {

				canTake := iActorTaker.CanTakeActor(lastOne)

				self.PopActor()

				if canTake {
					iActorTaker.TakeActor(lastOne)
					lastOne = nil
				}
			}
		}
	}
}

//Entity taker behavour

func (self *BufferBehaviour) TakeActor(actor *SimActor) bool {
	canTake := self.CanTakeActor(actor)

	if canTake {
		self.Actors = append(self.Actors, actor)
	}

	return canTake
}

func (self *BufferBehaviour) CanTakeActor(actor *SimActor) bool {
	return false
}

// Entity provideder behaviour
func (self *BufferBehaviour) PopActor() *SimActor {
	return self.PopBack()
}

func (self *BufferBehaviour) PopBack() *SimActor {
	if len(self.Actors) == 0 {
		return nil // or handle empty slice case
	}

	// Get the last element
	lastActor := self.Actors[len(self.Actors)-1]

	// Remove the last element by reslicing
	self.Actors = self.Actors[:len(self.Actors)-1]

	return lastActor
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
