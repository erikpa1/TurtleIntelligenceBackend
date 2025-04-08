package api

import (
	"github.com/gin-gonic/gin"
	"turtle/ctrl"
	"turtle/lg"
	"turtle/models"
	"turtle/tools"
)

func _TryLogin(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	if ctrl.UserExists(login, password) {
		lg.LogE("User don't exists")
	}

	//TODO toto naimpelmentovat
	tmp := models.NewUser()
	tmp.Uid = "poseidon"
	tmp.Email = "poseidon@turtle.sk"
	tmp.Firstname = "Poseidon"
	tmp.Surname = "The God"
	tmp.Type = "superadmin"
	tmp.Org = "olymp"

	tools.AutoReturn(c, tmp)

}

func initUsersApi(r *gin.Engine) {
	r.POST("/api/login", _TryLogin)
}
