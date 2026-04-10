package main

import (
	"turtle/simulation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	worldUid, _ := primitive.ObjectIDFromHex("69d9148ae3786775863e2fcf")
	simulation.RunSimulation(worldUid, bson.M{})

}
