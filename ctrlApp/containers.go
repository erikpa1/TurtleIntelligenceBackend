package ctrlApp

import (
	"go.mongodb.org/mongo-driver/bson"
	"turtle/db"
	"turtle/modelsApp"
)

const CT_CONTAINERS = "containers"

func QueryContainers(query bson.M) []*modelsApp.Container {
	return db.QueryEntities[modelsApp.Container](CT_CONTAINERS, query)
}

func CreateContainer(ct *modelsApp.Container) {

}
