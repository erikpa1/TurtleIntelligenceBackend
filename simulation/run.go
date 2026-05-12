package simulation

import "turtle/core/lgr"

func Run() {
	//worldUid, _ := primitive.ObjectIDFromHex("69d9148ae3786775863e2fcf")
	//simulation.RunSimulation(worldUid, bson.M{})

	newActor := NewSimActor()

	buffer := NewNewBufferBehaviour()

	takeEntity, haveTakeEntity := GetSimFunction[TakeActorFunction](buffer, TAKE_ACTOR_FUNCTION_NAME)

	if haveTakeEntity {
		takeEntity(buffer, newActor)
	}

	lgr.Info("Here")
}
