package behaviours

import (
	"turtle/simulation2/entities"
	"turtle/simulation2/rvar"
)

var BEH_WORKER_FUNCTTIONS = entities.SimFunctions{}

func InitBehWorkerPool() {

	entities.BEH_FACTORY.Behaviours["workerPool"] = NewBehWorkerPool

	var _init1 entities.FnInit = _BehWorkerInit1
	BEH_WORKER_FUNCTTIONS[entities.FN_INIT1] = _init1

	var _step entities.FnStep = _BehWorkerStep
	BEH_WORKER_FUNCTTIONS[entities.FN_STEP] = _step
}

func NewBehWorkerPool(entity *entities.SimEntity) {
	pool := &BehWorkerPool{}
	pool.Entity = entity
	pool.World = entity.World
	pool.WorkersMap = make(map[int64]*Worker)

	pool.WorkersCount = rvar.NewRvarr(entity.TypeData.GetString("workers_count", "1"))
	pool.SpawnOnRequest = entity.TypeData.GetBool("spawn_limit", false)
	pool.SpawnOnInit = entity.TypeData.GetBool("spawn_on_init", false)

	entity.Impl = pool
	entity.Functions = BEH_WORKER_FUNCTTIONS
}

func _BehWorkerInit1(self *entities.SimEntity) {
	pool := GetWorkerPool(self)
	pool.InitSpawnWorkers()

}

func _BehWorkerStep(self *entities.SimEntity) {
	pool := GetWorkerPool(self)
	pool.Step()
}
