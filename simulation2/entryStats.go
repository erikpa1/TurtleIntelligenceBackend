package simulation2

var ENTRY_STATS_FN = SimFunctions{}

func InitEntryStatistics() {

	BEH_FACTORY.Behaviours["entry_statistics"] = NewEntryStatisticsBeh

	var _takeActor FnTakeActor = _StatisticsTakeActor
	ENTRY_STATS_FN[FN_TAKE_ACTOR_FUNCTION_NAME] = _takeActor

}

func NewEntryStatisticsBeh(entity *SimEntity) {
	stats := &BehEntryStatistics{}
	stats.Entity = entity
	stats.World = entity.World

	entity.Impl = stats
	entity.Functions = ENTRY_STATS_FN
}

func _StatisticsTakeActor(entity *SimEntity, actor *SimActor) bool {
	spawn := GetBehEntryStatistics(entity)
	return spawn.TakeActor(actor)
}
