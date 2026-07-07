package behSpawn

import (
	"turtle/simulation2/entities"
	"turtle/simulation2/rvar"
	"turtle/tools"
)

var SPAWN_FUNCTIONS = entities.SimFunctions{}

func InitBehSpawn() {

	entities.BEH_FACTORY.Behaviours["spawn"] = NewSpawnBehaviour

	var _init1 entities.FnInit = _SpawnInit1
	SPAWN_FUNCTIONS[entities.FN_INIT1] = _init1

	var _step entities.FnStep = _SpawnStep
	SPAWN_FUNCTIONS[entities.FN_STEP] = _step
}

func NewSpawnBehaviour(entity *entities.SimEntity) {
	spawn := &BehSpawn{}
	spawn.Entity = entity
	spawn.World = entity.World

	spawn.SpawnInterval = rvar.NewRvarr(entity.TypeData.GetString("spawn_interval", "1"))
	spawn.SpawnLimit = entity.TypeData.GetInt("spawn_limit", 1)
	spawn.SpawnOnInit = entity.TypeData.GetBool("spawn_on_init", false)

	spawn.SpawnMultiplication = entity.TypeData.GetInt("spawn_multiplication", 1)

	if spawn.SpawnMultiplication <= 0 {
		spawn.SpawnMultiplication = 1
	}

	spawn.SpawnActorUid = entity.TypeData.GetPrimitiveObjectId("actor")
	spawn.NextSpawnTime = tools.MaxSeconds()

	entity.Impl = spawn
	entity.Functions = SPAWN_FUNCTIONS
}

func _SpawnInit1(self *entities.SimEntity) {
	spawn := GetBehSpawn(self)
	if spawn.SpawnOnInit {
		spawn.Spawn()
	} else {
		spawn._CalculateNextSpawn()
	}

}

func _SpawnStep(self *entities.SimEntity) {
	spawn := GetBehSpawn(self)
	spawn.Step()
}
