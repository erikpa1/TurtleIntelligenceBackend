package ctrl

import (
	"github.com/erikpa1/turtle/db"
	"github.com/erikpa1/turtle/models"
	"github.com/erikpa1/turtle/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_PROJECTS = "projects"

func ListProjects(org primitive.ObjectID) []*models.TurtleProject {
	return db.QueryEntities[models.TurtleProject](CT_PROJECTS, bson.M{"org": org})
}

func COUProject(user *models.User, project *models.TurtleProject) {
	if project.Uid.IsZero() {
		project.Org = user.Org
		project.CreatedBy = user.Uid
		db.InsertEntity(CT_PROJECTS, project)
	} else {

		project.UpdatedAt = tools.GetTimeNowMillis()

		db.UpdateOneCustom(CT_PROJECTS,
			bson.M{
				"_id": project.Uid,
				"org": user.Org,
			},
			bson.M{"$set": project})
	}
}
