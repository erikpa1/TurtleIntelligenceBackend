package behBuffer

import "turtle/simulation2/entities"

var BUFFER_FUNCTIONS = entities.SimFunctions{}

func InitBehBuffer() {
	entities.BEH_FACTORY.Behaviours["buffer"] = NewBufferBehaviour

	var _takeActor entities.FnTakeActor = _BufferTakeEntity
	BUFFER_FUNCTIONS[entities.FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor

	var _step entities.FnStep = _StepBufferBeh
	BUFFER_FUNCTIONS[entities.FN_STEP] = _step
}

func NewBufferBehaviour(entity *entities.SimEntity) {

	buffer := &BehBuffer{}
	buffer.Entity = entity
	buffer.World = entity.World

	buffer.Capacity = entity.TypeData.GetInt64("capacity", 8)
	buffer.InitialActor = entity.TypeData.GetPrimitiveObjectId("initial_actor")
	buffer.InitialCount = entity.TypeData.GetInt64("initial_count", 8)

	entity.Impl = buffer
	entity.Functions = BUFFER_FUNCTIONS
}

func _BufferTakeEntity(self *entities.SimEntity, actor *entities.SimActor) bool {
	buffer := GetBehBuffer(self)
	return buffer.TakeActor(actor)
}

func _StepBufferBeh(self *entities.SimEntity) {
	buffer := GetBehBuffer(self)
	buffer.Step()
}
