package simulation

var BUFFER_FUNCTIONS = SimFunctions{}

func InitBehBuffer() {
	var TakeActor TakeActorFunction = _SimBufferTakeEntity
	BUFFER_FUNCTIONS[TAKE_ACTOR_FUNCTION_NAME] = TakeActor
}

func _SimBufferTakeEntity(self *SimEntity, actor *SimActor) bool {
	return false
}

func NewNewBufferBehaviour() *SimEntity {
	entity := &SimEntity{
		Functions: BUFFER_FUNCTIONS,
	}
	return entity
}
