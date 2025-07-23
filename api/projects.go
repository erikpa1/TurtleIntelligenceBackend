package api

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/ctrl"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
)

func _ListProjects(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ctrl.ListProjects(user.Org))
}

func _COUProject(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	project := tools.ObjFromJsonPtr[models.TurtleProject](c.PostForm("data"))

	ctrl.COUProject(user, project)
	tools.AutoReturn(c, project.Uid)

}

func initApiProjects(r *gin.Engine) {
	r.GET("/api/projects", auth.LoginOrApp, _ListProjects)
	r.POST("/api/project", auth.LoginRequired, _COUProject)
}
