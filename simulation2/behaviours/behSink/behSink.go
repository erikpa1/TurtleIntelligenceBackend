package behSink

import "turtle/simulation2/entities"

var SINK_FUNCITONS = entities.SimFunctions{}

func InitBehSink() {

	entities.BEH_FACTORY.Behaviours["sink"] = NewSinkBehaviour

	var _take entities.FnTakeActor = _SinkTakeActor
	SINK_FUNCITONS[entities.FN_TAKE_ACTOR_FUNCTION_NAME] = _take

}

func NewSinkBehaviour(entity *entities.SimEntity) {
	entity.Functions = SINK_FUNCITONS
	entity.Impl = &BehSink{}
}
