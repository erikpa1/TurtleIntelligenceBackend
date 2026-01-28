package api

import (
	"turtle/core/users"
	"turtle/ctrl"
	"turtle/lgr"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _TryLogin(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	if ctrl.UserExists(login, password) {
		lgr.Error("User don't exists")
	}

	//TODO toto naimpelmentovat
	tmp := users.NewUser()
	tmp.Email = "poseidon@turtle.sk"
	tmp.Firstname = "Poseidon"
	tmp.Surname = "The God"
	tmp.Type = 3

	tools.AutoReturn(c, tmp)

}

func initUsersApi(r *gin.Engine) {
	r.POST("/api/login", _TryLogin)
}
