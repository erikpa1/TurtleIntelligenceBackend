package simulation2

import (
	"turtle/simulation/stats"
	"turtle/tools"
)

var DELAY_FUNCTIONS = SimFunctions{}

func InitBehDelay() {

	BEH_FACTORY.Behaviours["delay"] = NewDelayBehaviour

	var _takeActor FnTakeActor = _SimDelayTakeEntity
	DELAY_FUNCTIONS[FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor

	var _step FnStep = _StepDelayBehaviour
	DELAY_FUNCTIONS[FN_STEP] = _step

}

func NewDelayBehaviour(entity *SimEntity) {
	delay := &BehDelay{}
	delay.Entity = entity
	delay.World = entity.World
	delay.Actors = make(map[tools.Seconds][]*SimActor)

	delay.DelayTime = entity.TypeData.GetString("delayTime", "00:10")

	delay.Statistics = stats.NewProcessStats()

	entity.Impl = delay
	entity.Functions = DELAY_FUNCTIONS
}

func _SimDelayTakeEntity(self *SimEntity, actor *SimActor) bool {
	delay := GetBehDelay(self)
	return delay.TakeActor(actor)
}

func _StepDelayBehaviour(self *SimEntity) {
	delay := GetBehDelay(self)
	delay.Step()
}
