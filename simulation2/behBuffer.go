package simulation2

var BUFFER_FUNCTIONS = SimFunctions{}

func InitBehBuffer() {
	var _takeActor FnTakeActor = _BufferTakeEntity
	BUFFER_FUNCTIONS[FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor

	var _step FnStep = _StepBufferBeh
	BUFFER_FUNCTIONS[FN_STEP] = _step
}

func NewBufferBehaviour(entity *SimEntity) {

	buffer := &BehBuffer{}
	buffer.Entity = entity
	buffer.World = entity.World

	buffer.Capacity = entity.TypeData.GetInt64("capacity", 8)
	buffer.InitialActor = entity.TypeData.GetPrimitiveObjectId("initial_actor")
	buffer.InitialCount = entity.TypeData.GetInt64("initial_count", 8)

	entity.Impl = buffer
	entity.Functions = BUFFER_FUNCTIONS
}

func _BufferTakeEntity(self *SimEntity, actor *SimActor) bool {
	buffer := GetBehBuffer(self)
	return buffer.TakeActor(actor)
}

func _StepBufferBeh(self *SimEntity) {
	buffer := GetBehBuffer(self)
	buffer.Step()
}
