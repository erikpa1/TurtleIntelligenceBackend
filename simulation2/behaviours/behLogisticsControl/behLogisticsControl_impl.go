package behLogisticsControl

import (
	"sort"
	"strconv"

	"turtle/core/lgr"
	"turtle/simulation2/behaviours/behWorkerPool"
	"turtle/simulation2/entities"
	"turtle/simulation2/events"
)

// typeData keys authored by the frontend (LogisticsControlBehProperties.tsx).
const (
	keyTableMode   = "tableMode"    // "embedded" | "reference"
	keyTable       = "controlTable" // embedded string[][] matrix
	keyTableRef    = "tableRef"     // uid of a world Table entity
	keyDefaultPool = "defaultPool"  // uid of a workerPool entity

	// By-name column mapping for referenced tables (values are column keys).
	keyMapSource      = "mapSource"
	keyMapDestination = "mapDestination"
	keyMapActor       = "mapActor"
	keyMapPool        = "mapPool"
	keyMapPriority    = "mapPriority"

	modeReference = "reference"
)

// Fixed column order for the embedded table (matches COLUMNS on the frontend).
const (
	embSource      = 0
	embDestination = 1
	embActor       = 2
	embPool        = 3
	embPriority    = 4
)

// columnIndex holds the resolved cell index of each mission field (-1 if absent).
type columnIndex struct {
	source      int
	destination int
	actor       int
	pool        int
	priority    int
}

// routeTarget is one delivery a source can trigger, in priority order.
type routeTarget struct {
	Destination *entities.SimEntity
	Pool        *behWorkerPool.BehWorkerPool
	ActorType   string
	Priority    int
}

type BehLogisticsControl struct {
	World  *entities.SimWorld
	Entity *entities.SimEntity

	// routes maps a source entity RuntimeId to its ordered delivery targets.
	routes      map[int64][]routeTarget
	defaultPool *behWorkerPool.BehWorkerPool
}

func GetLogisticsControl(entity *entities.SimEntity) *BehLogisticsControl {
	return entity.Impl.(*BehLogisticsControl)
}

func (self *BehLogisticsControl) Init() {

	self.routes = make(map[int64][]routeTarget)
	self.defaultPool = self.resolvePool(self.Entity.TypeData.GetString(keyDefaultPool, ""))

	rows, cols := self.loadRows()
	if len(rows) == 0 {
		lgr.Error("LogisticsControl [%s]: no mission rows to route", self.Entity.Name)
		return
	}

	for rowIdx, row := range rows {

		source := self.entityAt(row, cols.source)
		if source == nil {
			lgr.Error("LogisticsControl row %d: source entity not found", rowIdx)
			continue
		}

		destination := self.entityAt(row, cols.destination)
		if destination == nil {
			lgr.Error("LogisticsControl row %d: destination entity not found", rowIdx)
			continue
		}

		pool := self.resolvePool(cellAt(row, cols.pool))
		if pool == nil {
			pool = self.defaultPool
		}
		if pool == nil {
			lgr.Error("LogisticsControl row %d: no worker pool and no default pool", rowIdx)
			continue
		}

		priority := 0
		if parsed, err := strconv.Atoi(cellAt(row, cols.priority)); err == nil {
			priority = parsed
		}

		// Subscribe to a source's actor-taken events only once.
		if _, seen := self.routes[source.RuntimeId]; !seen {
			source.Aee.On(events.ACTOR_TAKEN, self._GoForActor_FromEvent)
		}

		self.routes[source.RuntimeId] = append(self.routes[source.RuntimeId], routeTarget{
			Destination: destination,
			Pool:        pool,
			ActorType:   cellAt(row, cols.actor),
			Priority:    priority,
		})
	}

	// Order each source's targets by ascending priority.
	for id := range self.routes {
		targets := self.routes[id]
		sort.SliceStable(targets, func(i, j int) bool {
			return targets[i].Priority < targets[j].Priority
		})
		self.routes[id] = targets
	}
}

func (self *BehLogisticsControl) Step() {

}

func (self *BehLogisticsControl) _GoForActor_FromEvent(args ...interface{}) {

	data, castOk := args[0].(events.ActorTakenStruct)
	if !castOk {
		lgr.Error("LogisticsControl: failed to cast to events.ActorTakenStruct")
		return
	}

	self.DispatchMission(data.Entity, data.Actor)
}

// DispatchMission sends a free worker on a pickup→deliver mission for `item`,
// trying the source's pools in priority order until one has a free worker.
func (self *BehLogisticsControl) DispatchMission(source *entities.SimEntity, item *entities.SimActor) {

	if source == nil || item == nil {
		return
	}

	targets, ok := self.routes[source.RuntimeId]
	if !ok {
		return
	}

	for _, target := range targets {
		worker := target.Pool.GetFreeWorker()
		if worker != nil {
			worker.StartMission(item, target.Destination.Position, target.ActorType)
			return
		}
	}
}
