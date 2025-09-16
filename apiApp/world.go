package apiApp

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/ctrlApp"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/modelsApp"
	"github.com/erikpa1/TurtleIntelligenceBackend/simulation"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func _SaveWorld(c *gin.Context) {

	type SaveWorldResponse struct {
		Uid                primitive.ObjectID      `json:"uid"`
		Modified           []*modelsApp.Entity     `json:"modified"`
		Created            []*modelsApp.Entity     `json:"created"`
		Deleted            []primitive.ObjectID    `json:"deleted"`
		CreatedConnections [][2]primitive.ObjectID `json:"createdConnections"`
		DeletedConnections [][2]primitive.ObjectID `json:"deletedConnections"`
	}

	request := tools.ObjFromJson[SaveWorldResponse](c.PostForm("data"))

	lg.LogOk(request.Created)
	if len(request.Created) > 0 {
		ctrlApp.CreateEntities(request.Created)
	}

	lg.LogOk(request.CreatedConnections)
	if len(request.CreatedConnections) > 0 {
		ctrlApp.CreateConnections(request.Uid, request.CreatedConnections)
	}

	lg.LogOk(request.Modified)
	if len(request.Modified) > 0 {
		ctrlApp.UpdateEntities(request.Modified)
	}

	lg.LogOk(request.DeletedConnections)
	if len(request.DeletedConnections) > 0 {

		for _, conn := range request.DeletedConnections {
			ctrlApp.DeleteConnection(conn[0], conn[1])
		}
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
			"uid":         model.Uid,
			"name":        model.Name,
			"entities":    ctrlApp.ListEntitiesOfWorld(model.Uid),
			"connections": ctrlApp.ListConnectionsOfWorld(model.Uid),
		})
	} else {
		tools.AutoNotFound(c, "notfound")
	}

}

func _PlayWorld(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	tools.AutoReturn(c, simulation.RunSimulation(uid, bson.M{}))
}

func _PauseWorld(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	simulation.PauseSimulation(uid)
}

func _StopWorld(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	simulation.StopSimulation(uid)
}

func _ResumeWorld(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	simulation.ResumeSimulation(uid)
}

func init_api_world(r *gin.Engine) {
	r.GET("/api/w", auth.LoginRequired, _GetWorld)
	r.POST("/api/w/save", auth.LoginRequired, _SaveWorld)
	r.POST("/api/w/simulate", _PlayWorld)
	r.POST("/api/w/stop", _StopWorld)
	r.POST("/api/w/pause", _PauseWorld)
	r.POST("/api/w/resume", _ResumeWorld)

}
