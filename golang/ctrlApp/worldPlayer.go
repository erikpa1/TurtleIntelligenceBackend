package ctrlApp

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/lg"
	"turtle/server"
	"turtle/tools"
)

func PlayWorld(worldUid primitive.ObjectID) {

	go func() {
		tools.Recover("Failed to run simulation")

	}()

	entities := QueryWorldEntities(bson.M{"model": worldUid})

	for _, entity := range entities {
		lg.LogI(entity)
	}

	server.MYIO.Emit("sim", bson.M{
		"spawned":   []bson.M{},
		"unspawned": []int64{},
		"states":    []bson.M{},
	})

}
