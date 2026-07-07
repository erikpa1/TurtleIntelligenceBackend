package behLogisticsControl

import (
	"turtle/simulation2/entities"
)

type BehLogisticsControl struct {
	World  *entities.SimWorld
	Entity *entities.SimEntity
}

func GetLogisticsControl(entity *entities.SimEntity) *BehLogisticsControl {
	return entity.Impl.(*BehLogisticsControl)
}

func (self *BehLogisticsControl) Step() {

}
