package apiApp

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/auth"
	"turtle/ctrlApp"
	"turtle/lg"
	"turtle/modelsApp"
	"turtle/tools"
)

func _SaveWorld(c *gin.Context) {
	//TODO

	type SaveWorldResponse struct {
		Uid      string               `json:"uid"`
		Modified []*modelsApp.Entity  `json:"modified"`
		Created  []*modelsApp.Entity  `json:"created"`
		Deleted  []primitive.ObjectID `json:"deleted"`
	}

	request := tools.ObjFromJson[SaveWorldResponse](c.PostForm("data"))

	lg.LogOk(request.Created)
	if len(request.Created) > 0 {
		ctrlApp.CreateEntities(request.Created)
	}

	lg.LogOk(request.Modified)
	if len(request.Modified) > 0 {
		ctrlApp.UpdateEntities(request.Modified)
	}

	lg.LogOk(request.Deleted)
	if len(request.Deleted) > 0 {
		ctrlApp.DeleteEntities(request.Deleted)
	}

}

func _GetWorld(c *gin.Context) {

	uid, _ := primitive.ObjectIDFromHex(c.Query("uid"))

	model := ctrlApp.GetModel(uid)

	if model != nil {
		tools.AutoReturn(c, bson.M{
			"uid":      model.Uid,
			"name":     model.Name,
			"entities": ctrlApp.ListEntitiesOfWorld(model.Uid),
		})
	} else {
		tools.AutoNotFound(c, "notfound")
	}

}

func _PlayWorld(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	ctrlApp.PlayWorld(uid)
}

func init_api_world(r *gin.Engine) {
	r.GET("/api/w", auth.LoginRequired, _GetWorld)
	r.POST("/api/w/save", auth.LoginRequired, _SaveWorld)
	r.POST("/api/w/simulate", _PlayWorld)

}
