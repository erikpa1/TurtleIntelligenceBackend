package behLogisticsControl

import (
	"turtle/core/lgr"
	"turtle/simulation2/behaviours/behWorkerPool"
	"turtle/simulation2/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// loadRows returns the mission matrix and the column index of each field,
// resolved from either the embedded table or a referenced world Table.
func (self *BehLogisticsControl) loadRows() ([][]string, columnIndex) {

	td := self.Entity.TypeData

	if td.GetString(keyTableMode, "") == modeReference {
		return self.loadReferencedRows()
	}

	return td.GetStringMatrix(keyTable), columnIndex{
		source:      embSource,
		destination: embDestination,
		actor:       embActor,
		pool:        embPool,
		priority:    embPriority,
	}
}

// loadReferencedRows reads a world Table entity's grid and resolves the field
// column indices from the by-name mapping stored on this control.
func (self *BehLogisticsControl) loadReferencedRows() ([][]string, columnIndex) {

	td := self.Entity.TypeData
	refUid := td.GetString(keyTableRef, "")

	table := self.entityByUid(refUid)
	if table == nil {
		lgr.Error("LogisticsControl: referenced table %q not found", refUid)
		return nil, columnIndex{}
	}

	rows := table.TypeData.GetStringMatrix("rows")

	// Map the referenced table's column keys to their positions.
	keyToIdx := map[string]int{}
	for i, col := range table.TypeData.GetObjectArray("columns") {
		keyToIdx[col.GetString("key", "")] = i
	}

	idxOf := func(mappingKey string) int {
		if idx, ok := keyToIdx[td.GetString(mappingKey, "")]; ok {
			return idx
		}
		return -1
	}

	return rows, columnIndex{
		source:      idxOf(keyMapSource),
		destination: idxOf(keyMapDestination),
		actor:       idxOf(keyMapActor),
		pool:        idxOf(keyMapPool),
		priority:    idxOf(keyMapPriority),
	}
}

func cellAt(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func (self *BehLogisticsControl) entityAt(row []string, idx int) *entities.SimEntity {
	return self.entityByUid(cellAt(row, idx))
}

func (self *BehLogisticsControl) entityByUid(hex string) *entities.SimEntity {
	if hex == "" {
		return nil
	}
	uid, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		lgr.Error("LogisticsControl: invalid uid %q: %v", hex, err)
		return nil
	}
	return self.World.SimEntities[uid]
}

func (self *BehLogisticsControl) resolvePool(hex string) *behWorkerPool.BehWorkerPool {
	entity := self.entityByUid(hex)
	if entity == nil {
		return nil
	}
	pool, ok := entity.Impl.(*behWorkerPool.BehWorkerPool)
	if !ok {
		lgr.Error("LogisticsControl: entity %q is not a worker pool (got %T)", entity.Name, entity.Impl)
		return nil
	}
	return pool
}
