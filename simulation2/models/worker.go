package models

import (
	"turtle/simulation/simMath"
	"turtle/tools"
)

type Worker struct {
	RuntimeId  int64
	Name       string
	ActiveTask *WorkerTask
	Aee        tools.AnyEventEmitter
	Position   simMath.Position
}

func NewWorker() *Worker {
	return &Worker{
		Position: simMath.Position{0, 0, 0},
	}
}
