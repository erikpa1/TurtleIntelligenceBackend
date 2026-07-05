package behaviours

import (
	"turtle/core/lgr"
	"turtle/simulation2/entities"
	"turtle/simulation2/rvar"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

	for self.SpawnOnInit {

		spawnVar := self.WorkersCount.GetInt()

		for i := 0; i < spawnVar; i++ {
			self.SpawnWorker()
		}
	}
}

func (self *BehWorkerPool) SpawnWorker() {
	worker := NewWorker()

	actor := self.World.SpawnActorWithUid(primitive.ObjectID{})
	worker.Actor = actor
	self.WorkersMap[worker.RuntimeId] = worker

}

func (self *BehWorkerPool) Step() {
	lgr.Error("Stepping")
}
