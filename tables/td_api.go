package tables

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/auth"
	"turtle/tools"
)

func _QueryTableData(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryHeader[bson.M](c)
	namespace := c.GetHeader("namespace")
	tools.AutoReturn(c, QueryTableData(user, namespace, query))
}

func CreateGinRouting(r *gin.Engine, nameSpace string) {
	r.GET(fmt.Sprintf("/api/%s/tables/query", nameSpace), auth.LoginOrApp, _QueryTableData)
}
