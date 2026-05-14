package simulation2

import "turtle/core/lgr"

var PROCESS_FUNCTIONS = SimFunctions{}

func InitBehProcess() {
	var _takeActor TakeActorFunction = _SimProcessTakeEntity
	PROCESS_FUNCTIONS[FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor
}

func _SimProcessTakeEntity(self *SimEntity, actor *SimActor) bool {
	lgr.Error("Process taking actor")
	return false
}

func NewProcessBehaviour(entity *SimEntity) {
	entity.Functions = PROCESS_FUNCTIONS

}
