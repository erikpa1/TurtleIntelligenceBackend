package tables

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func _QueryTables(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryHeader[bson.M](c)
	serverKit.ReturnOkJson(c, QueryTables(user, query))
}

func _COUTable(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	table := tools.ObjFromJsonPtr[TurtleTable](c.PostForm("data"))
	COUTable(user, table)
}

func _DeleteTable(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	table := serverKit.MongoObjectIdFromQuery(c)
	DeleteTable(user, table)
}

func InitTablesApi(r *gin.Engine) {
	r.GET("/api/tables", auth.LoginOrApp, _QueryTables)
	r.POST("/api/table", auth.LoginOrApp, _COUTable)
	r.DELETE("/api/table", auth.LoginRequired, _DeleteTable)
}
