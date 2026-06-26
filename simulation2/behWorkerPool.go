package simulation2

import (
	"turtle/simulation2/rvar"
)

var BEH_WORKER_FUNCTTIONS = SimFunctions{}

func InitBehWorkerPool() {

	BEH_FACTORY.Behaviours["workerPool"] = NewSpawnBehaviour

	var _init1 FnInit = _SpawnInit1
	SPAWN_FUNCTIONS[FN_INIT1] = _init1

	var _step FnStep = _SpawnStep
	SPAWN_FUNCTIONS[FN_STEP] = _step
}

func NewBehWorkerPool(entity *SimEntity) {
	pool := &BehWorkerPool{}
	pool.Entity = entity
	pool.World = entity.World

	pool.WorkersCount = rvar.NewRvarr(entity.TypeData.GetString("workers_count", "1"))
	pool.SpawnOnRequest = entity.TypeData.GetBool("spawn_limit", false)
	pool.SpawnOnInit = entity.TypeData.GetBool("spawn_on_init", false)

	entity.Impl = pool
	entity.Functions = BEH_WORKER_FUNCTTIONS
}

func _BehWorkerInit1(self *SimEntity) {
	spawn := GetBehSpawn(self)
	if spawn.SpawnOnInit {
		spawn.Spawn()
	} else {
		spawn._CalculateNextSpawn()
	}

}

func _BehWorkerStepStep(self *SimEntity) {
	spawn := GetBehSpawn(self)
	spawn.Step()
}
