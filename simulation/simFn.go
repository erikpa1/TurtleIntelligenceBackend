package simulation

type SimFunctions map[string]interface{}

func GetSimFunction[T any](entity *SimEntity, fnName string) (T, bool) {
	tmp, casted := entity.Functions[fnName].(T)
	return tmp, casted
}

const TAKE_ACTOR_FUNCTION_NAME = "TakeActor"

type TakeActorFunction func(self *SimEntity, actor *SimActor) bool

func InitSimFunctions() bool {
	InitBehBuffer()

	return true
}

var SimInit = InitSimFunctions()
