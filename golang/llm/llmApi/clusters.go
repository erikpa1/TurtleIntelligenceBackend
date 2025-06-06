package llmApi

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/llm/llmCtrl"
	"turtle/llm/llmModels"
	"turtle/tools"
)

func _ListLLMClusters(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, llmCtrl.ListLLMClusters(user.Org))
}

func _ListLlmModels(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, llmCtrl.ListLLMModels(user.Org))
}

func _DeleteLLMCluster(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	llmCtrl.DeleteLLMCluster(user, uid)
}

func _COULLMCluster(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	data := tools.ObjFromJson[llmModels.LLMCluster](c.PostForm("data"))
	llmCtrl.COULLMCluster(user, &data)
}

func _GetLLMCluster(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	tools.AutoReturn(c, llmCtrl.GetLLMCluster(user, uid))
}

func _COULlmModel(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	data := tools.ObjFromJson[llmModels.LlmModel](c.PostForm("data"))
	llmCtrl.COULLMModel(user, &data)
}

func _DeleteLlmModel(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	llmCtrl.DeleteLLMModel(user, uid)
}

func _GetLLMModel(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	tools.AutoReturn(c, llmCtrl.GetLLMModel(user, uid))
}

func InitClustersOnly(r *gin.Engine) {
	r.GET("/api/llm/clusters", auth.LoginRequired, _ListLLMClusters)
	r.GET("/api/llm/cluster", auth.LoginRequired, _GetLLMCluster)
	r.POST("/api/llm/cluster", auth.LoginRequired, _COULLMCluster)
	r.DELETE("/api/llm/cluster", auth.LoginRequired, _DeleteLLMCluster)
}

func InitModelsOnly(r *gin.Engine) {

	r.GET("/api/llm/models", auth.LoginRequired, _ListLlmModels)
	r.GET("/api/llm/model", auth.LoginRequired, _GetLLMModel)
	r.POST("/api/llm/model", auth.LoginRequired, _COULlmModel)
	r.DELETE("/api/llm/model", auth.LoginRequired, _DeleteLlmModel)
}

func InitLlmAndClusters(r *gin.Engine) {
	InitClustersOnly(r)
	InitModelsOnly(r)
}
