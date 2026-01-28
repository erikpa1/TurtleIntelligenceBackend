package simulation

import (
	"math/rand"
	"turtle/lgr"
	"turtle/simulation/simMath"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HumanBehaviour struct {
	World  *SimWorld
	Entity *SimEntity

	SpawnLimit int
	SpawnActor primitive.ObjectID

	Actors []*SimActor

	tmpTargets []ISimBehaviour
}

func NewHumanBehaviour() ISimBehaviour {
	return &HumanBehaviour{}
}

// Implementacia IBehaviour
func (self *HumanBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *HumanBehaviour) Init1() {

	self.Actors = make([]*SimActor, self.SpawnLimit)

	for i := 0; i < self.SpawnLimit; i++ {

		lgr.Error("Spawning actor")

		human := CreateHumanActor()
		human.Uid = primitive.NewObjectID()
		{
			//Position duplication
			myPosition := self.Entity.Position
			human.Position[0] = myPosition[0]
			human.Position[1] = myPosition[1]
			human.Position[2] = myPosition[2]
		}

		self.Actors[i] = human
		self.World.SpwanCustomActor(human)
	}

	self.tmpTargets = make([]ISimBehaviour, len(self.Actors))

	behaviours := self.World.ListBehaviours()

	for i, _ := range self.Actors {
		randomIndex := rand.Intn(len(behaviours))
		self.tmpTargets[i] = behaviours[randomIndex]

	}

}

func (self *HumanBehaviour) Init2() {

}

func (self *HumanBehaviour) GetEntity() *SimEntity {
	return self.Entity
}

func (self *HumanBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = entity

	self.SpawnLimit = entity.TypeData.GetInt("spawn_limit", 1)
	self.SpawnActor = entity.TypeData.GetPrimitiveObjectId("actor")
}

func (self *HumanBehaviour) TakeActor(actor *SimActor) bool {
	self.World.UnspawnActor(actor)
	return true
}

func (self *HumanBehaviour) CanTakeActor(actor *SimActor) bool {
	return true
}

func (self *HumanBehaviour) __Interface() {
	procesB := HumanBehaviour{}
	var _ ISimBehaviour = &procesB

}

func CreateHumanActor() *SimActor {
	tmp := NewSimActor()
	tmp.Name = "Human"
	tmp.Color = "#eb8a02"

	return tmp
}

func (self *HumanBehaviour) Step() {

	behaviours := self.World.ListBehaviours()

	for i, actor := range self.Actors {
		simEntity := self.tmpTargets[i]

		distanceToPoint := actor.Position.MoveTo(simEntity.GetEntity().Position, simMath.AVG_WALKING_SPEED)
		actor.PosChanged()

		if distanceToPoint <= 0.05 {
			randomIndex := rand.Intn(len(behaviours))
			self.tmpTargets[i] = behaviours[randomIndex]
			lgr.Info("Going on the another point")
		}
	}

}
