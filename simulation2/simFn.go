package simulation2

import "turtle/core/lgr"

type SimFunctions map[string]interface{}

func GetSimFunction[T any](entity *SimEntity, fnName string) (T, bool) {
	var zero T
	fn, haveFn := entity.Functions[fnName]

	if haveFn == false {
		return zero, false
	}

	tmp, casted := fn.(T)

	if casted == false {
		lgr.Error("Failed to cast [%s][%s][%T]", entity.Type, fnName, fn)
		return zero, false
	}

	return tmp, true
}

type TakeActorFunction func(self *SimEntity, actor *SimActor) bool
