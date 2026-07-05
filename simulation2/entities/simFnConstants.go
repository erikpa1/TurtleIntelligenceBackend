package entities

const FN_STEP = "Step"

type FnStep func(self *SimEntity)
type FnInit func(self *SimEntity)

const FN_INIT_BASE = "Init"
const FN_INIT1 = "Init1"
const FN_INIT2 = "Init2"

const FN_TAKE_ACTOR_FUNCTION_NAME = "TakeActor"
const FN_CAN_TAKE_ACTOR_FUNCTION_NAME = "CanTakeActor"
