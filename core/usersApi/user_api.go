package usersApi

import (
	"turtle/lg"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"turtle/auth"
	"turtle/core/pipeline"
	"turtle/core/users"
	"turtle/reflective"
	"turtle/tools"
)

func _COUUser(c *gin.Context) {
	ginPipeline := pipeline.NewGinPipeline(c)

	var user users.User
	if ginPipeline.ShouldBindJSON(&user) {
		users.COUUser(ginPipeline, &user)
	} else {
		lg.LogE("Here")
	}

}

func _ListUsers(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, users.QueryUsers(user, bson.M{}))
}

func _QueryUsers(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryBsonHeader(c)
	tools.AutoReturn(c, users.QueryUsers(user, query))
}

func InitUsersApi(r *gin.Engine) {

	r.GET("/api/users", auth.AdminRequired, _ListUsers)
	r.GET("/api/users/query", auth.AdminRequired, _QueryUsers)
	r.POST("/api/user", auth.AdminRequired, _COUUser)

	reflective.CreateDeleteIdRoute(r, "users", "/users/:id", auth.AdminRequired)
}
