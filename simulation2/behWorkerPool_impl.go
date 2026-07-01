package simulation2

import (
	"turtle/core/lgr"
	"turtle/simulation2/models"
	"turtle/simulation2/rvar"
)

type BehWorkerPool struct {
	World  *SimWorld
	Entity *SimEntity

	SpawnOnInit    bool
	SpawnOnRequest bool //Spawns entity when no free actor required
	WorkersCount   *rvar.Rvar

	WorkersMap map[int64]*models.Worker
}

func GetWorkerPool(entity *SimEntity) *BehWorkerPool {
	return entity.Impl.(*BehWorkerPool)
}

func (self *BehWorkerPool) InitSpawnWorkers() {

	for self.SpawnOnInit {

		spawnVar := self.WorkersCount.GetInt()

		for i := 0; i < spawnVar; i++ {
			self.SpawnWorker()
		}
	}
}

func (self *BehWorkerPool) SpawnWorker() {
	worker := models.NewWorker()
	worker.Position = self.Entity.Position
	worker.RuntimeId = self.World.GetRuntimeId()
}

func (self *BehWorkerPool) Step() {
	lgr.Error("Stepping")
}
