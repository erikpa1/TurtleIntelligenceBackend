package simulation2

import (
	"turtle/simulation2/models"
	"turtle/tools"
)

type Worker struct {
	RuntimeId  int64
	Name       string
	ActiveTask *models.WorkerTask
	Aee        tools.AnyEventEmitter
	Actor      *SimActor
}

func NewWorker() *Worker {
	return &Worker{
		Actor: &SimActor{},
	}
}
