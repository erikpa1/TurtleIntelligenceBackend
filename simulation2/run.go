package simulation2

import "turtle/core/lgr"

func Run() {
	//worldUid, _ := primitive.ObjectIDFromHex("69d9148ae3786775863e2fcf")
	//simulation.RunSimulation(worldUid, bson.M{})

	newActor := NewSimActor()

	spawn := NewSimEntity()
	NewSpawnBehaviour(spawn)

	buffer := NewSimEntity()
	NewBufferBehaviour(buffer)

	takeEntity, haveTakeEntity := GetSimFunction[TakeActorFunction](buffer, FN_TAKE_ACTOR_FUNCTION_NAME)

	if haveTakeEntity {
		takeEntity(buffer, newActor)
	}

	process := NewSimEntity()
	NewProcessBehaviour(process)

	takeEntity, haveTakeEntity = GetSimFunction[TakeActorFunction](process, FN_TAKE_ACTOR_FUNCTION_NAME)

	if haveTakeEntity {
		takeEntity(buffer, newActor)
	}

	lgr.Info("Ending test run")
}
