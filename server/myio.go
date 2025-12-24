package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Myio struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         *tools.TimedMutex
}

func NewMyio() *Myio {
	tmp := &Myio{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mu:         tools.NewTimedMutex("MyioMutex", 5),
	}
	tmp.mu.Monitor(3 * time.Second)
	return tmp
}

func (self *Myio) NotifyAll(data any) {
	self.Emit("notifications", data)
}

func (self *Myio) EmitToRoom(room string, channel string, data any) {
	defer tools.Recover("Failed to emit MYIO")

	self.mu.Lock("Emit")
	clientsCopy := copyClientMap(self.clients)
	self.mu.Unlock()
	//passss

	for _, client := range clientsCopy {
		_, ok := client.rooms[room]
		if ok {
			client.emit(channel, data, nil)
		}

	}

}

func (self *Myio) Emit(channel string, data any) {
	defer tools.Recover("Failed to emit MYIO")

	self.mu.Lock("Emit")
	clientsCopy := copyClientMap(self.clients)
	self.mu.Unlock()
	//passss

	for _, client := range clientsCopy {
		client.emit(channel, data, nil)
	}
}

func (self *Myio) EmitSync(channel string, data any) {
	defer tools.Recover("Failed to emit MYIO")

	self.mu.Lock("Emit")
	clientsCopy := copyClientMap(self.clients)
	self.mu.Unlock()
	//passss

	for _, client := range clientsCopy {
		client.emitSync(channel, data, nil)
	}
}

// Message represents the structure of the messages exchanged
type Message struct {
	Event   string      `json:"event"`
	Data    interface{} `json:"data"`
	AckId   *int        `json:"ackId,omitempty"` // Optional ackId
	AckData interface{} `json:"ackData,omitempty"`
}

// Client represents a connected WebSocket client
type Client struct {
	uid   string
	conn  *websocket.Conn
	send  chan []byte
	mu    *tools.TimedMutex
	rooms map[string]any
}

func NewClient(conn *websocket.Conn) *Client {

	tmp := Client{
		uid:   tools.GetUUID4(),
		conn:  conn,
		send:  make(chan []byte, 256),
		mu:    tools.NewTimedMutex("ClientMutex", 5),
		rooms: make(map[string]any),
	}
	return &tmp
}

// Event handlers map
var eventHandlers = map[string]func(*Client, Message){}

func copyClientMap(original map[string]*Client) map[string]*Client {
	// Initialize a new map with the same type
	copied := make(map[string]*Client)

	// Copy each key-value pair from the original map to the new map
	for uid, client := range original {
		copied[uid] = client
	}

	return copied
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (for dev purposes)
	},
}

// HandleWebSocketConnection handles a new WebSocket connection
func HandleWebSocketConnection(c *gin.Context) {
	// WebSocket connection upgrader

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		lg.LogE("WebSocket upgrade error:", err)
		return
	}
	client := NewClient(conn)
	client.mu.Monitor(3 * time.Second)

	MYIO.mu.Lock("HandleWebSocketConnection")
	MYIO.clients[client.uid] = client
	MYIO.mu.Unlock()

	go client.readPump()
	go client.writePump()

	// When a client connects, you can log or send a welcome message
	lg.LogI("Client connected")
}

// readPump handles incoming messages from the client
func (c *Client) readPump() {
	defer func() {
		// This will ensure client is unregistered and resources are freed up
		MYIO.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// Detect if it's a normal close event
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected WebSocket close: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			lg.LogE("Invalid message format:", err)
			continue
		}

		if handler, ok := eventHandlers[msg.Event]; ok {
			handler(c, msg)
		}
	}
}

// writePump sends messages to the client
func (c *Client) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				lg.LogE("Closing connection")
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)

			if err != nil {
				lg.LogI("Closing the connection")
				c.conn.Close()

				MYIO.mu.Lock("writePump")
				delete(MYIO.clients, c.uid)
				MYIO.mu.Unlock()

			}

		}
	}
}

// writePump sends messages to the client
func (c *Client) writePumpSync(message []byte) {
	err := c.conn.WriteMessage(websocket.TextMessage, message)

	if err != nil {
		lg.LogI("Closing the connection")
		c.conn.Close()

		MYIO.mu.Lock("writePump")
		delete(MYIO.clients, c.uid)
		MYIO.mu.Unlock()

	}

}

// emit sends a message to the client
func (c *Client) emitSync(event string, data interface{}, ackId *int) {
	defer func() {
		tools.Recover("Failed to emit to MYIO client")
		c.mu.Unlock()
	}()

	msg := Message{
		Event: event,
		Data:  data,
		AckId: ackId,
	}
	message, err := json.Marshal(msg)

	c.mu.Lock("Client|emit")
	if err != nil {
		lg.LogStackTraceErr(err.Error())
	} else {
		c.writePumpSync(message)
	}

}

// emit sends a message to the client
func (c *Client) emit(event string, data interface{}, ackId *int) {
	defer func() {
		tools.Recover("Failed to emit to MYIO client")
		c.mu.Unlock()
	}()

	msg := Message{
		Event: event,
		Data:  data,
		AckId: ackId,
	}
	message, err := json.Marshal(msg)

	c.mu.Lock("Client|emit")
	if err != nil {
		lg.LogStackTraceErr(err.Error())
	} else {
		c.send <- message
	}

}

// handleAcknowledgment sends an acknowledgment for a received message
func (c *Client) handleAcknowledgment(ackId int, ackData interface{}) {
	msg := Message{
		AckId:   &ackId,
		AckData: ackData,
	}
	message, _ := json.Marshal(msg)
	c.send <- message
}

func _JoinFeed(feed_uid string) {
	//TODO dorobit registraciu na feed
}

func _RegisterClientOnChannel(c *gin.Context) {
	room := c.Query("room")
	who := c.Query("who")

	_ = room
	_ = who

	MYIO.mu.Lock("_RegisterClientOnChannel")

	//client, ok := MYIO.clients[who]{}

	defer MYIO.mu.Unlock()

}

func RunMyioServer(r *gin.Engine) {

	r.POST("/my.io/conn", _RegisterClientOnChannel)

	eventHandlers["message"] = func(c *Client, msg Message) {
		lg.LogE("Received message:", msg.Data)

		if msg.AckId != nil {
			ackData := "Server received your message!"
			c.handleAcknowledgment(*msg.AckId, ackData)
		}

	}

	r.GET("/my.io", HandleWebSocketConnection)
	lg.LogOk("My.io registered")
}

var MYIO = NewMyio()
