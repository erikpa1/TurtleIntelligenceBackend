package simulation2

var BUFFER_FUNCTIONS = SimFunctions{}

func InitBehBuffer() {
	var _takeActor FnTakeActor = _BufferTakeEntity
	BUFFER_FUNCTIONS[FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor

	var _step FnStep = _StepBufferBeh
	PROCESS_FUNCTIONS[FN_STEP] = _step
}

func NewBufferBehaviour(entity *SimEntity) {

	//TODO sem naimplementovat premenne

	entity.Functions = BUFFER_FUNCTIONS
}

func _BufferTakeEntity(self *SimEntity, actor *SimActor) bool {
	proces := GetBehProcess(self)
	return proces.TakeActor(actor)
}

func _StepBufferBeh(self *SimEntity) {
	buffer := GetBehBuffer(self)
	buffer.Step()
}
