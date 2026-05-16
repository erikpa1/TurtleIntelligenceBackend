package simulation2

type IActorTaker = SimEntity

func (self *IActorTaker) TakeActor(actor *SimActor) bool {
	takeFn, haveTakeFn := GetSimFunction[FnTakeActor](self, FN_TAKE_ACTOR_FUNCTION_NAME)
	if haveTakeFn {
		return takeFn(self, actor)
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
