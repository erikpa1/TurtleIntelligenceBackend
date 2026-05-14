package simulation2

const FN_STEP = "Step"

type FnStep func(self *SimEntity)
type FnInit func(self *SimEntity)

const FN_INIT = "Init"

const FN_TAKE_ACTOR_FUNCTION_NAME = "TakeActor"
