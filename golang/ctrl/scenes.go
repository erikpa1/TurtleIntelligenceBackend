package ctrl

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/db"
	"turtle/models"
	"turtle/tools"
)

const CT_SCENES = "scenes"

func ListScenes(org primitive.ObjectID, parent primitive.ObjectID) []*models.TurtleScene {
	return db.QueryEntities[models.TurtleScene](CT_SCENES, bson.M{
		"org":    org,
		"parent": parent,
	})
}

func COUScene(user *models.User, scene *models.TurtleScene) {

	if scene.Uid.IsZero() {
		scene.Org = user.Org
		scene.CreatedBy = user.Org
		scene.UpdatedBy = user.Org
		scene.CreatedAt = tools.GetTimeNowMillis()

		db.InsertEntity(CT_SCENES, scene)
	} else {
		scene.UpdatedAt = tools.GetTimeNowMillis()
		scene.UpdatedBy = user.Org

		db.UpdateOneCustom(CT_SCENES,
			bson.M{
				"_id": scene.Uid,
				"org": scene.Org,
			},
			bson.M{"$set": scene})
	}

}

func GetScene(org primitive.ObjectID, uid primitive.ObjectID) *models.TurtleScene {
	return db.QueryEntity[models.TurtleScene](CT_SCENES, bson.M{
		"org": org,
		"_id": uid,
	})
}
