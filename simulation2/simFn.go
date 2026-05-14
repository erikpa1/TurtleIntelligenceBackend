package simulation2

type SimFunctions map[string]interface{}

func GetSimFunction[T any](entity *SimEntity, fnName string) (T, bool) {
	tmp, casted := entity.Functions[fnName].(T)
	return tmp, casted
}

type TakeActorFunction func(self *SimEntity, actor *SimActor) bool
