package simulation

import "go.mongodb.org/mongo-driver/bson/primitive"

type QueueBehaviour struct {
	World  *SimWorld
	Entity *SimEntity

	Actors []*SimActor

	Capacity     int64
	InitialCount int64
	InitialActor primitive.ObjectID
}

func NewQueueBehaviour() ISimBehaviour {
	return &QueueBehaviour{
		Actors: make([]*SimActor, 0),
	}
}

func (self *QueueBehaviour) __Interface() {
	var _ ISimBehaviour = &QueueBehaviour{}
	var _ ActorTakerBehaviour = &QueueBehaviour{}
	var _ ActorProviderBehaviour = &QueueBehaviour{}
}

// Implementacia IBehaviour

func (self *QueueBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *QueueBehaviour) GetEntity() *SimEntity {
	return self.Entity
}

func (self *QueueBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = entity

	self.Capacity = entity.TypeData.GetInt64("capacity", 8)
	self.InitialActor = entity.TypeData.GetPrimitiveObjectId("initial_actor")
	self.InitialCount = entity.TypeData.GetInt64("initial_count", 8)

}

func (self *QueueBehaviour) Init1() {

}

func (self *QueueBehaviour) Init2() {

}

func (self *QueueBehaviour) Step() {
	self._TryToPassActorsNext()
}

func (self *QueueBehaviour) _TryToPassActorsNext() {

	if len(self.Actors) > 0 {

		lastOne := self.Actors[0]

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

// Entity taker behavour
func (self *QueueBehaviour) TakeActor(actor *SimActor) bool {
	canTake := self.CanTakeActor(actor)

	if canTake {
		self.Actors = append(self.Actors, actor)
	}

	return canTake
}

func (self *QueueBehaviour) CanTakeActor(actor *SimActor) bool {
	return false
}

// Entity provideder behaviour
func (self *QueueBehaviour) PopActor() *SimActor {
	return self.PopFront()
}

func (self *QueueBehaviour) PopFront() *SimActor {
	if len(self.Actors) == 0 {
		return nil // or handle empty slice case
	}

	// Get the last element
	lastActor := self.Actors[0]

	// Remove the last element by reslicing
	self.Actors = self.Actors[1 : len(self.Actors)-1]

	return lastActor
}

func (self *QueueBehaviour) HasAnyActor() bool {
	return false
}

func (self *QueueBehaviour) HasActorOfType(actorType string) bool {
	return false
}

func (self *QueueBehaviour) HasActorWithVariable(variable string, value any) bool {
	return true
}
