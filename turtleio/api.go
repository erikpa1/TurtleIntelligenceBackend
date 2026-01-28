package turtleio

import (
	"net/http"

	"turtle/auth"
	"turtle/lgr"

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

	lgr.ErrorJson(wsHandler)

}
