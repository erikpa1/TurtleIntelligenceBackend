package behLogisticsControl

import (
	"turtle/core/lgr"
	"turtle/simulation2/behaviours/behBuffer"
	"turtle/simulation2/behaviours/behWorkerPool"
	"turtle/simulation2/entities"
	"turtle/simulation2/events"
)

type BehLogisticsControl struct {
	World  *entities.SimWorld
	Entity *entities.SimEntity

	WorkerPools map[int64]*behWorkerPool.BehWorkerPool
}

func GetLogisticsControl(entity *entities.SimEntity) *BehLogisticsControl {
	return entity.Impl.(*BehLogisticsControl)
}

func (self *BehLogisticsControl) Init() {

	for _, entity := range self.World.SimEntities {

		switch impl := entity.Impl.(type) {
		case *behBuffer.BehBuffer:
			{
				entity.Aee.On(events.ACTOR_TAKEN, self._GoForActor_FromEvent)
			}
		case *behWorkerPool.BehWorkerPool:
			{
				self.WorkerPools[entity.RuntimeId] = impl
			}
		}
	}

}

func (self *BehLogisticsControl) Step() {

}

func (self *BehLogisticsControl) _GoForActor_FromEvent(args ...interface{}) {
	data, castOk := args[0].(events.ActorTakenStruct)

	if castOk {
		self.SendSomethingForActor(data.Actor)
	} else {
		lgr.Error("Failed to cast to events.ActorTakenStruct")
	}
}

func (self *BehLogisticsControl) SendSomethingForActor(actor *entities.SimActor) {

	for _, workers := range self.WorkerPools {
		tmp := workers.GetFreeWorker()

		if tmp != nil {
			tmp.GoForActor(actor)
		}
	}

}
