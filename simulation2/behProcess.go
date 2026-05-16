package simulation2

import "turtle/simulation/stats"

var PROCESS_FUNCTIONS = SimFunctions{}

func InitBehProcess() {
	var _takeActor FnTakeActor = _SimProcessTakeEntity
	PROCESS_FUNCTIONS[FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor

	var _step FnStep = _StepProcessBehaviour
	PROCESS_FUNCTIONS[FN_STEP] = _step

}

func NewProcessBehaviour(entity *SimEntity) {
	proces := &BehProcess{}
	proces.Entity = entity
	proces.World = entity.World
	proces.ProcessTime = entity.TypeData.GetString("processTime", "00:10")

	proces.Statistics = stats.NewProcessStats()

	entity.Impl = proces
	entity.Functions = PROCESS_FUNCTIONS
}

func _SimProcessTakeEntity(self *SimEntity, actor *SimActor) bool {
	proces := GetBehProcess(self)
	return proces.TakeActor(actor)
}

func _StepProcessBehaviour(self *SimEntity) {
	proces := GetBehProcess(self)
	proces.Step()
}
