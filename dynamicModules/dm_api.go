package dynamicModules

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/auth"
	"turtle/db"
	"turtle/lg"
	"turtle/tools"
)

func InitDefaultEntitiesApi(r *gin.Engine, namespace string, bucket string) {
	totalName := namespace + "/" + bucket

	dbName := namespace + "_" + bucket

	lg.LogOk(fmt.Sprintf("Initializing dynamic module [%s]", totalName))

	_QueryDefaultEntities := func(c *gin.Context) {
		user := auth.GetUserFromContext(c)
		query := tools.QueryHeader[bson.M](c)
		tools.AutoReturn(c, db.QueryEntities[bson.M](dbName, user.FillOrgQuery(query)))
	}

	_COUDefaultEntity := func(c *gin.Context) {
		user := auth.GetUserFromContext(c)
		entity := tools.ObjFromJson[bson.M](c.PostForm("data"))

		uid, _ := primitive.ObjectIDFromHex(entity["uid"].(string))

		if uid.IsZero() {
			db.InsertEntity(dbName, entity)
		} else {
			db.UpdateOneCustom(dbName, bson.M{
				"uid": uid,
				"org": user.Org,
			}, entity)
		}
	}

	_DeleteDefaultEntity := func(c *gin.Context) {
		user := auth.GetUserFromContext(c)
		entityUid := tools.MongoObjectIdFromQuery(c)
		db.DeleteEntity(dbName, bson.M{"_id": entityUid, "org": user.Org})

	}

	r.GET("/api/"+totalName+"s", auth.LoginOrApp, _QueryDefaultEntities)
	r.GET("/api/"+totalName+"s/query", auth.LoginOrApp, _QueryDefaultEntities)
	r.POST("/api/"+totalName, auth.LoginOrApp, _COUDefaultEntity)
	r.DELETE("/api/"+totalName, auth.LoginRequired, _DeleteDefaultEntity)
}
