package simulation2

import (
	"turtle/core/lgr"
	"turtle/simulation2/behaviours/behBuffer"
	"turtle/simulation2/behaviours/behProcess"
	"turtle/simulation2/behaviours/behSpawn"
	"turtle/simulation2/entities"
)

func Run() {
	//worldUid, _ := primitive.ObjectIDFromHex("69d9148ae3786775863e2fcf")
	//simulation.RunSimulation(worldUid, bson.M{})

	newActor := entities.NewSimActor()

	spawn := entities.NewSimEntity()
	behSpawn.NewSpawnBehaviour(spawn)

	buffer := entities.NewSimEntity()
	behBuffer.NewBufferBehaviour(buffer)

	takeEntity, haveTakeEntity := entities.GetSimFunction[entities.FnTakeActor](buffer, entities.FN_TAKE_ACTOR_FUNCTION_NAME)

	if haveTakeEntity {
		takeEntity(buffer, newActor)
	}

	process := entities.NewSimEntity()
	behProcess.NewProcessBehaviour(process)

	takeEntity, haveTakeEntity = entities.GetSimFunction[entities.FnTakeActor](process, entities.FN_TAKE_ACTOR_FUNCTION_NAME)

	if haveTakeEntity {
		takeEntity(buffer, newActor)
	}

	lgr.Info("Ending test run")
}
