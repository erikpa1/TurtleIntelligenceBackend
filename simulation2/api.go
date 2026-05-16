package simulation2

import (
	"turtle/auth"
	"turtle/core/lgr"
	"turtle/core/serverKit"
	"turtle/ctrlApp"
	"turtle/modelsApp"
	"turtle/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func _SaveWorld(c *gin.Context) {

	type SaveWorldResponse struct {
		Uid                primitive.ObjectID       `json:"uid"`
		Modified           []*modelsApp.WorldEntity `json:"modified"`
		Created            []*modelsApp.WorldEntity `json:"created"`
		Deleted            []primitive.ObjectID     `json:"deleted"`
		CreatedConnections [][2]primitive.ObjectID  `json:"createdConnections"`
		DeletedConnections [][2]primitive.ObjectID  `json:"deletedConnections"`
	}

	request := tools.ObjFromJson[SaveWorldResponse](c.PostForm("data"))

	lgr.OkJson(request.Created)
	if len(request.Created) > 0 {
		ctrlApp.CreateEntities(request.Created)
	}

	lgr.OkJson(request.CreatedConnections)
	if len(request.CreatedConnections) > 0 {
		ctrlApp.CreateConnections(request.Uid, request.CreatedConnections)
	}

	lgr.OkJson(request.Modified)
	if len(request.Modified) > 0 {
		ctrlApp.UpdateEntities(request.Modified)
	}

	lgr.OkJson(request.DeletedConnections)
	if len(request.DeletedConnections) > 0 {

		for _, conn := range request.DeletedConnections {
			ctrlApp.DeleteConnection(conn[0], conn[1])
		}
	}

	lgr.OkJson(request.Deleted)
	if len(request.Deleted) > 0 {
		ctrlApp.DeleteEntities(request.Deleted)
	}

}

func _GetWorld(c *gin.Context) {

	uid, _ := primitive.ObjectIDFromHex(c.Query("uid"))

	model := ctrlApp.GetModel(uid)

	if model != nil {
		serverKit.ReturnOkJson(c, bson.M{
			"uid":         model.Uid,
			"name":        model.Name,
			"entities":    ctrlApp.ListEntitiesOfWorld(model.Uid),
			"connections": ctrlApp.ListConnectionsOfWorld(model.Uid),
		})
	} else {
		serverKit.Return404(c, "notfound")
	}

}

func _PlayWorld(c *gin.Context) {
	uid := serverKit.MongoObjectIdFromQuery(c)
	serverKit.ReturnOkJson(c, RunSimulation(uid, bson.M{}))
}

func _PauseWorld(c *gin.Context) {
	uid := serverKit.MongoObjectIdFromQuery(c)
	PauseSimulation(uid)
}

func _StopWorld(c *gin.Context) {
	uid := serverKit.MongoObjectIdFromQuery(c)
	StopSimulation(uid)
}

func _ResumeWorld(c *gin.Context) {
	uid := serverKit.MongoObjectIdFromQuery(c)
	ResumeSimulation(uid)
}

func InitSimulationApi(r *gin.Engine) {
	r.GET("/api/simulation", auth.LoginRequired, _GetWorld)
	r.POST("/api/simulation/save", auth.LoginRequired, _SaveWorld)
	r.POST("/api/simulation/simulate", _PlayWorld)
	r.POST("/api/simulation/stop", _StopWorld)
	r.POST("/api/simulation/pause", _PauseWorld)
	r.POST("/api/simulation/resume", _ResumeWorld)

}
