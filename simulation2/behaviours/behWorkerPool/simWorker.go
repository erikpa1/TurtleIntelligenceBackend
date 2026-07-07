package behWorkerPool

import (
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
}

func NewWorker(actor *entities.SimActor) *Worker {
	return &Worker{
		Actor: actor,
	}
}
