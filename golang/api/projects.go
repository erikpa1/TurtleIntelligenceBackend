package api

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/ctrl"
	"turtle/lg"
	"turtle/models"
	"turtle/tools"
)

func _ListProjects(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ctrl.ListProjects(user.Org))
}

func _COUProject(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	lg.LogOk(c.PostForm("data"))

	project := tools.ObjFromJsonPtr[models.TurtleProject](c.PostForm("data"))

	if project.Uid == "" {
		project.Uid = tools.GetUUID4()
		project.At = tools.GetTimeNowMillis()
	}

	project.Org = user.Org
	project.CreatedBy = user.Org

	ctrl.COUProject(user.Org, project)

	tools.AutoReturn(c, project.Uid)

}

func initApiProjects(r *gin.Engine) {
	r.GET("/api/projects", auth.LoginOrApp, _ListProjects)
	r.POST("/api/project", auth.LoginRequired, _COUProject)
}
