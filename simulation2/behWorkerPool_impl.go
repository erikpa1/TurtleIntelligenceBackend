package simulation2

import (
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

func (worker *BehWorkerPool) InitSpawnActors() {

}
