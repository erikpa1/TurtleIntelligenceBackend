package turtleio

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/lg"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

type CreateRoomReq struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (self *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	self.hub.Rooms[req.Id] = &Room{
		Id:      req.Id,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

func (self *Handler) JoinRoom(c *gin.Context) {

	user := auth.GetUserFromContext(c)
	lg.LogE(user.Uid.Hex())

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		lg.LogE("WebSocket upgrade error:", err)
		return
	}

	roomUid := c.Param("roomUid")

	lg.LogE(roomUid)

	cl := &Client{
		Conn:    conn,
		Message: make(chan *Message, 10),
		Id:      user.Uid.Hex(),
		RoomId:  roomUid,
	}

	self.hub.Register <- cl

	go cl.writeMessage()
	cl.readMessage(self.hub)
}
