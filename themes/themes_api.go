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

}

func _COUTheme(c *gin.Context) {
}

func _DeleteTheme(c *gin.Context) {

}

func InitThemesApi(r *gin.Engine) {
	r.GET("/api/themes", auth.AdminRequired, _ListThemes)
	r.GET("/api/theme", auth.AdminRequired, _ListThemes)
	r.GET("/api/theme/default", auth.AdminRequired, _ListThemes)

	r.POST("/api/theme", auth.AdminRequired, _COUTheme)
	r.DELETE("/api/theme", auth.AdminRequired, _DeleteTheme)

	r.POST("/api/themes/copilot", auth.AdminRequired, _ChatCopilot)
}
