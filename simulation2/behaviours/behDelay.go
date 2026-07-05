package behaviours

import (
	"turtle/simulation/stats"
	"turtle/simulation2/entities"
	"turtle/simulation2/rvar"
	"turtle/tools"
)

var DELAY_FUNCTIONS = entities.SimFunctions{}

func InitBehDelay() {

	entities.BEH_FACTORY.Behaviours["delay"] = NewDelayBehaviour

	var _takeActor entities.FnTakeActor = _SimDelayTakeEntity
	DELAY_FUNCTIONS[entities.FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor

	var _step entities.FnStep = _StepDelayBehaviour
	DELAY_FUNCTIONS[entities.FN_STEP] = _step

}

func NewDelayBehaviour(entity *entities.SimEntity) {
	delay := &BehDelay{}
	delay.Entity = entity
	delay.World = entity.World
	delay.Actors = make(map[tools.Seconds][]*entities.SimActor)

	delay.DelayTime = rvar.NewRvarr(entity.TypeData.GetString("delayTime", "00:10"))

	delay.Statistics = stats.NewProcessStats()

	entity.Impl = delay
	entity.Functions = DELAY_FUNCTIONS
}

func _SimDelayTakeEntity(self *entities.SimEntity, actor *entities.SimActor) bool {
	delay := GetBehDelay(self)
	return delay.TakeActor(actor)
}

func _StepDelayBehaviour(self *entities.SimEntity) {
	delay := GetBehDelay(self)
	delay.Step()
}
