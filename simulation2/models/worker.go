package models

import "turtle/tools"

type Worker struct {
	RuntimeId  int64
	Name       string
	ActiveTask *WorkerTask
	Aee        tools.AnyEventEmitter
}

func NewWorker() *Worker {
	return &Worker{}
}
