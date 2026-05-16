package simulation2

import "turtle/core/lgr"

func _SinkTakeActor(entity *SimEntity, actor *SimActor) bool {
	//TODO dat sem statistiky
	lgr.Error("Unspawning entity: %v", actor)
	entity.World.UnspawnActor(actor)
	return true
}
