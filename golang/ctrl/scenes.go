package ctrl

import (
	"go.mongodb.org/mongo-driver/bson"
	"turtle/db"
	"turtle/models"
)

const CT_SCENES = "scenes"

func ListScenes(org, parent string) []*models.TurtleScene {
	return db.QueryEntities[models.TurtleScene](CT_SCENES, bson.M{
		"org":    org,
		"parent": parent,
	})
}

func COUScene(org string, project *models.TurtleScene) {
	db.COUEntity(CT_SCENES, project)
}

func GetScene(org, uid string) *models.TurtleScene {
	return db.QueryEntity[models.TurtleScene](CT_SCENES, bson.M{
		"org": org,
		"uid": uid,
	})
}
