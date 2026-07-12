package entities

type IActorTaker = SimEntity

func (self *IActorTaker) TakeActor(actor *SimActor) bool {
	takeFn, haveTakeFn := GetSimFunction[FnTakeActor](self, FN_TAKE_ACTOR_FUNCTION_NAME)
	if haveTakeFn {
		tmp := takeFn(self, actor)

		if tmp {
			actor.ParentEntity = self
		}

		return tmp
	} else {
		return false
	}
}
func (self *IActorTaker) CanTakeActor(actor *SimActor) bool {
	takeFn, haveTakeFn := GetSimFunction[FnTakeActor](self, FN_CAN_TAKE_ACTOR_FUNCTION_NAME)
	if haveTakeFn {
		return takeFn(self, actor)
	} else {
		return false
	}
}
