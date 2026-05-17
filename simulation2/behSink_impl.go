package simulation2

import "turtle/core/lgr"

type BehSink struct {
	Count int64
}

func GetBehSink(entity *SimEntity) *BehSink {
	return CastImplementation[BehSink](entity.Impl)
}

func _SinkTakeActor(entity *SimEntity, actor *SimActor) bool {
	//TODO dat sem statistiky
	lgr.Error("Unspawning entity: %v", actor)
	entity.World.UnspawnActor(actor)

	tmp := GetBehSink(entity)
	tmp.Count += 1
	entity.World.UpdateActorState(entity.RuntimeId, "count", tmp.Count)

	return true
}
