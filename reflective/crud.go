package reflective

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateDeleteIdRoute(r *gin.Engine, container string, path string, authMethod gin.HandlerFunc) {
	r.DELETE(path, authMethod, func(c *gin.Context) {
		uid := tools.MongoObjectIdFromQuery(c)
		db.DeleteEntity(container, bson.M{"_id": uid})
	})
}
