package turtleio

import (
	"github.com/gorilla/websocket"
	"turtle/lg"
)

type Client struct {
	Id      string
	Conn    *websocket.Conn
	Message chan *Message
	RoomId  string
}

type Message struct {
	Content  string `json:"content"`
	RoomId   string `json:"roomId"`
	Username string `json:"username"`
}

func (self *Client) writeMessage() {
	defer func() {
		self.Conn.Close()
	}()

	for {
		message, ok := <-self.Message

		if ok {
			self.Conn.WriteJSON(message)
		} else {
			return
		}
	}

}

func (self *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- self
		self.Conn.Close()
	}()

	for {
		_, m, err := self.Conn.ReadMessage()

		if err != nil {

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				lg.LogE(err.Error())
			} else {
				lg.LogE(err.Error())
			}

			break
		} else {
			msg := &Message{
				Content: string(m),
				RoomId:  self.RoomId,
			}

			hub.Broadcast <- msg
		}
	}

}
