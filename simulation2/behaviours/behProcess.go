package behaviours

import (
	"turtle/simulation/stats"
	"turtle/simulation2/entities"
	"turtle/simulation2/rvar"
)

var PROCESS_FUNCTIONS = entities.SimFunctions{}

func InitBehProcess() {

	entities.BEH_FACTORY.Behaviours["process"] = NewProcessBehaviour

	var _takeActor entities.FnTakeActor = _SimProcessTakeEntity
	PROCESS_FUNCTIONS[entities.FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor

	var _step entities.FnStep = _StepProcessBehaviour
	PROCESS_FUNCTIONS[entities.FN_STEP] = _step

}

func NewProcessBehaviour(entity *entities.SimEntity) {
	proces := &BehProcess{}
	proces.Entity = entity
	proces.World = entity.World
	proces.ProcessTime = rvar.NewRvarr(entity.TypeData.GetString("processTime", "00:10"))

	proces.Statistics = stats.NewProcessStats()

	entity.Impl = proces
	entity.Functions = PROCESS_FUNCTIONS
}

func _SimProcessTakeEntity(self *entities.SimEntity, actor *entities.SimActor) bool {
	proces := GetBehProcess(self)
	return proces.TakeActor(actor)
}

func _StepProcessBehaviour(self *entities.SimEntity) {
	proces := GetBehProcess(self)
	proces.Step()
}
