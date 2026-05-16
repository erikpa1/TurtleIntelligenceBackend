package simulation2

var SPAWN_FUNCTIONS = SimFunctions{}

func InitSpawnBehaviour() {
	var _step FnStep = _SpawnStep
	SPAWN_FUNCTIONS[FN_STEP] = _step

	var _init1 FnInit = _SpawnInit1
	SPAWN_FUNCTIONS[FN_INIT1] = _init1

}

func NewSpawnBehaviour(entity *SimEntity) {
	entity.Functions = SPAWN_FUNCTIONS

	spawnBehaviour := make(SimBehData)
	spawnBehaviour["SpawnInterval"] = entity.TypeData.GetSeconds("spawn_interval", 1)
	spawnBehaviour["SpawnLimit"] = entity.TypeData.GetInt64("spawn_limit", 1)
	spawnBehaviour["SpawnActor"] = entity.TypeData.GetPrimitiveObjectId("actor")

	spawnBehaviour["Actor"] = nil

	entity.BehaviourEntity = spawnBehaviour

}
