package behTable

import (
	"turtle/simulation2/entities"
)

// behTable is the backend counterpart of the placeable "Table" data object
// (Plant Simulation TableFile analogue). It holds no simulation behaviour of its
// own — other behaviours read its grid straight off the entity's TypeData via
// SafeJson.GetStringMatrix("rows"). Registering a constructor keeps LoadEntities
// from logging "no constructor" for every table in the model.

var BEH_TABLE = entities.SimFunctions{}

func InitBehTable() {
	entities.BEH_FACTORY.Behaviours["table"] = NewBehTable
}

func NewBehTable(entity *entities.SimEntity) {
	table := &BehTable{}
	table.Entity = entity
	table.World = entity.World

	entity.Impl = table
	entity.Functions = BEH_TABLE
}

type BehTable struct {
	World  *entities.SimWorld
	Entity *entities.SimEntity
}

func GetTable(entity *entities.SimEntity) *BehTable {
	return entity.Impl.(*BehTable)
}

// Rows returns the table's cell matrix as authored on the frontend.
func (self *BehTable) Rows() [][]string {
	return self.Entity.TypeData.GetStringMatrix("rows")
}
