package behWorkerPool

import (
	"turtle/simulation/simMath"
	"turtle/simulation2/entities"
	"turtle/simulation2/models"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var WORKER_OBJECTID = primitive.ObjectID{}

type Worker struct {
	Name       string
	ActiveTask *models.WorkerTask
	Aee        tools.AnyEventEmitter
	Actor      *entities.SimActor

	WalkSpeed   float32 // km/h
	Destination simMath.Position
	IsWalking   bool
}

func NewWorker(actor *entities.SimActor) *Worker {
	return &Worker{
		Actor:     actor,
		WalkSpeed: simMath.AVG_WALKING_SPEED,
	}
}

func (self *Worker) GoForActor(actor *entities.SimActor) {
	self.WalkTo(actor.Position)
}

// WalkTo sets a destination and starts walking towards it on each Step.
func (self *Worker) WalkTo(dest simMath.Position) {
	self.Destination = dest
	self.IsWalking = true
}

// Step advances the worker by one simulation tick, walking towards its
// destination at walking speed until it arrives.
func (self *Worker) Step() {
	if !self.IsWalking {
		return
	}

	position := self.Actor.Position
	remaining := position.MoveTo(self.Destination, self.WalkSpeed)
	self.Actor.UpdatePosition(position)

	if remaining == 0 {
		self.IsWalking = false
	}
}
