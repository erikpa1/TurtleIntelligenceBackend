package themes

import (
	"turtle/auth"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _ListThemes(c *gin.Context) {
	tools.AutoReturn(c, ListThemes())
}

func _GetTheme(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, GetTheme(user, uid))
}

func _GetDefaultTheme(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, GetDefaultTheme(user))
}

func _COUTheme(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	theme := tools.ObjFromJsonPtr[Theme](c.PostForm("data"))
	COUTheme(user, theme)
}

func _ImportTheme(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	theme := tools.ObjFromJsonPtr[Theme](c.PostForm("data"))
	ImportTheme(user, theme)
}

func _DeleteTheme(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	DeleteTheme(user, uid)

}

func InitThemesApi(r *gin.Engine) {
	r.GET("/api/themes", auth.LoginRequired, _ListThemes)
	r.GET("/api/theme", auth.LoginRequired, _GetTheme)
	r.GET("/api/theme/default", auth.LoginRequired, _GetDefaultTheme)

	r.POST("/api/theme", auth.AdminRequired, _COUTheme)
	r.POST("/api/theme/import", auth.AdminRequired, _ImportTheme)
	r.DELETE("/api/theme", auth.AdminRequired, _DeleteTheme)

	//Copilo part
	r.POST("/api/themes/copilot", auth.AdminRequired, _ChatCopilot)
	r.POST("/api/themes/copilot/examples", auth.AdminRequired, _ListCopilotExamples)
}
