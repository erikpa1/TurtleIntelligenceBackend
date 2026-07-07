package behWorkerPool

import (
	"turtle/core/lgr"
	"turtle/simulation2/entities"
	"turtle/simulation2/rvar"
)

type BehWorkerPool struct {
	World  *entities.SimWorld
	Entity *entities.SimEntity

	SpawnOnInit    bool
	SpawnOnRequest bool //Spawns entity when no free actor required
	WorkersCount   *rvar.Rvar

	WorkersMap map[int64]*Worker
}

func GetWorkerPool(entity *entities.SimEntity) *BehWorkerPool {
	return entity.Impl.(*BehWorkerPool)
}

func (self *BehWorkerPool) InitSpawnWorkers() {

	if self.SpawnOnInit {

		spawnVar := self.WorkersCount.GetInt()

		for i := 0; i < spawnVar; i++ {
			self.SpawnWorker()
		}
	}
}

func (self *BehWorkerPool) SpawnWorker() {
	worker := NewWorker(self.World.SpawnActorWithUid(WORKER_OBJECTID))
	self.WorkersMap[worker.Actor.Id] = worker

	worker.Actor.Color = "#ffffff"
	worker.Actor.UpdatePosition(self.Entity.Position)

	lgr.ErrorJson(self.Entity.Position)

}

func (self *BehWorkerPool) GetFreeWorker() *Worker {
	for _, worker := range self.WorkersMap {
		return worker
	}

	return nil
}

func (self *BehWorkerPool) Step() {

	for _, worker := range self.WorkersMap {
		worker.Step()
	}
}
