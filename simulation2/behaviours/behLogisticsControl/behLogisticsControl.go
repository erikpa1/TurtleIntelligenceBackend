package behLogisticsControl

import (
	"turtle/simulation2/entities"
)

var BEH_LOGISTICS_CONTROL = entities.SimFunctions{}

func InitBehWorkerPool() {

	entities.BEH_FACTORY.Behaviours["logisticsControl"] = NewBehWorkerPool

	var _init1 entities.FnInit = _BehWorkerInit1
	BEH_LOGISTICS_CONTROL[entities.FN_INIT1] = _init1

	var _step entities.FnStep = _BehWorkerStep
	BEH_LOGISTICS_CONTROL[entities.FN_STEP] = _step
}

func NewBehWorkerPool(entity *entities.SimEntity) {
	control := &BehLogisticsControl{}
	control.Entity = entity
	control.World = entity.World

	entity.Impl = control
	entity.Functions = BEH_LOGISTICS_CONTROL
}

func _BehWorkerInit1(self *entities.SimEntity) {
	control := GetLogisticsControl(self)
	control.Init()
}

func _BehWorkerStep(self *entities.SimEntity) {
	control := GetLogisticsControl(self)
	control.Step()
}
