package simulation2

import "turtle/core/lgr"

var BUFFER_FUNCTIONS = SimFunctions{}

func InitBehBuffer() {
	var _takeActor FnTakeActor = _SimBufferTakeEntity
	BUFFER_FUNCTIONS[FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor
}

func _SimBufferTakeEntity(self *SimEntity, actor *SimActor) bool {
	lgr.Error("Buffer taking entity")
	return false
}

func NewBufferBehaviour(entity *SimEntity) {
	entity.Functions = BUFFER_FUNCTIONS
}
