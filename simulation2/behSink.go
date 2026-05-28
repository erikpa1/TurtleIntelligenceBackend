package simulation2

var SINK_FUNCITONS = SimFunctions{}

func InitBehSink() {

	BEH_FACTORY.Behaviours["sink"] = NewSinkBehaviour

	var _take FnTakeActor = _SinkTakeActor
	SINK_FUNCITONS[FN_TAKE_ACTOR_FUNCTION_NAME] = _take

}

func NewSinkBehaviour(entity *SimEntity) {
	entity.Functions = SINK_FUNCITONS
	entity.Impl = &BehSink{}
}
