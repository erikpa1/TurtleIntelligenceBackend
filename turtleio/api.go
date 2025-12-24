package turtleio

import (
	"net/http"

	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (for dev purposes)
	},
}

func InitTurtleSocketsApi(r *gin.Engine) {
	wsHub := NewHub()
	wsHandler := NewHandler(wsHub)

	r.POST("/turtleio/room", auth.LoginRequired, wsHandler.CreateRoom)
	r.GET("/turtleio/join/:roomId", auth.LoginRequired, wsHandler.JoinRoom)

	wsHub.GoRun()

	lg.LogEson(wsHandler)

}
