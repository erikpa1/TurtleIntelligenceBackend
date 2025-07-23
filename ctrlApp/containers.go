package ctrlApp

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/modelsApp"
	"go.mongodb.org/mongo-driver/bson"
)

const CT_CONTAINERS = "containers"

func QueryContainers(query bson.M) []*modelsApp.Container {
	return db.QueryEntities[modelsApp.Container](CT_CONTAINERS, query)
}

func CreateContainer(ct *modelsApp.Container) {

}
