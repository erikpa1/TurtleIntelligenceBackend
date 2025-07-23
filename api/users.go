package api

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/ctrl"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
)

func _TryLogin(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	if ctrl.UserExists(login, password) {
		lg.LogE("User don't exists")
	}

	//TODO toto naimpelmentovat
	tmp := models.NewUser()
	tmp.Email = "poseidon@turtle.sk"
	tmp.Firstname = "Poseidon"
	tmp.Surname = "The God"
	tmp.Type = 3

	tools.AutoReturn(c, tmp)

}

func initUsersApi(r *gin.Engine) {
	r.POST("/api/login", _TryLogin)
}
