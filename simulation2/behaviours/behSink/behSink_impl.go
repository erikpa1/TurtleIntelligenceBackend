package behSink

import (
	"turtle/core/lgr"
	"turtle/simulation2/entities"
)

type BehSink struct {
	Count int64
}

func GetBehSink(entity *entities.SimEntity) *BehSink {
	return entities.CastImplementation[BehSink](entity.Impl)
}

func _SinkTakeActor(entity *entities.SimEntity, actor *entities.SimActor) bool {
	//TODO dat sem statistiky
	lgr.Error("Unspawning entity: %v", actor)
	entity.World.UnspawnActor(actor)

	tmp := GetBehSink(entity)
	tmp.Count += 1
	entity.World.UpdateActorState(entity.RuntimeId, "count", tmp.Count)

	return true
}
