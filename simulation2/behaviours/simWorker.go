package behaviours

import (
	"turtle/simulation2/entities"
	"turtle/simulation2/models"
	"turtle/tools"
)

type Worker struct {
	RuntimeId  int64
	Name       string
	ActiveTask *models.WorkerTask
	Aee        tools.AnyEventEmitter
	Actor      *entities.SimActor
}

func NewWorker() *Worker {
	return &Worker{
		Actor: &entities.SimActor{},
	}
}
