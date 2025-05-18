package simulation

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/ctrlApp"
	"turtle/lg"
	"turtle/server"
	"turtle/tools"
)

func RunSimulation(modelUid primitive.ObjectID, simParams bson.M) {

	entities := ctrlApp.QueryWorldEntities(bson.M{"model": modelUid})

	for _, entity := range entities {
		lg.LogI(entity)
	}

	world := NewSimWorld()
	world.IsOnline = true

	go func() {
		tools.Recover("Failed to run simulation")

		world.RunSimulation()

		server.MYIO.Emit("sim", bson.M{
			"spawned":   []bson.M{},
			"unspawned": []int64{},
			"states":    []bson.M{},
		})

	}()
}
