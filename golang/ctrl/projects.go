package ctrl

import (
	"go.mongodb.org/mongo-driver/bson"
	"turtle/db"
	"turtle/models"
)

const CT_PROJECTS = "projects"

func ListProjects(org string) []*models.TurtleProject {
	return db.QueryEntities[models.TurtleProject](CT_PROJECTS, bson.M{"org": org})
}

func COUProject(org string, project *models.TurtleProject) {
	db.COUEntity(CT_PROJECTS, project)
}
