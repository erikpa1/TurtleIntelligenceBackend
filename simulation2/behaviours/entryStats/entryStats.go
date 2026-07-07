package entryStats

import "turtle/simulation2/entities"

var ENTRY_STATS_FN = entities.SimFunctions{}

func InitEntryStatistics() {

	entities.BEH_FACTORY.Behaviours["entry_statistics"] = NewEntryStatisticsBeh

	var _takeActor entities.FnTakeActor = _StatisticsTakeActor
	ENTRY_STATS_FN[entities.FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor

}

func NewEntryStatisticsBeh(entity *entities.SimEntity) {
	stats := &BehEntryStatistics{}
	stats.Entity = entity
	stats.World = entity.World

	entity.Impl = stats
	entity.Functions = ENTRY_STATS_FN
}

func _StatisticsTakeActor(entity *entities.SimEntity, actor *entities.SimActor) bool {
	spawn := GetBehEntryStatistics(entity)
	return spawn.TakeActor(actor)
}
