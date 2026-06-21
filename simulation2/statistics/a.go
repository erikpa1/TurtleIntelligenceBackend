package statistics

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SimulationStatsPrepare(modelUid, runUid primitive.ObjectID) {

}

func SimulationStateChanged(runUid primitive.ObjectID, state SimulationState) {

}

func SimulationPing(runUid primitive.ObjectID) {

}

func SimulationFailed(runUid primitive.ObjectID, error string) {

}

func WriteSimulationStatistics(runUid primitive.ObjectID, statistics bson.M) {

}
