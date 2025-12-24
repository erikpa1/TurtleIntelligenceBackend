/*
	Inspired by tutorial:

https://www.youtube.com/watch?v=760GKM7s_5Y
*/
package turtleio

type Room struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Clients map[string]*Client
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (self *Hub) GoRun() {
	go self.Run()
}

func (self *Hub) Run() {
	for {
		select {
		//Registration of client into room
		case cl := <-self.Register:
			room, roomOk := self.Rooms[cl.RoomId]

			if roomOk {
				_, clientOk := room.Clients[cl.Id]

				if clientOk == false {
					room.Clients[cl.Id] = cl
				}
			}

			room.Clients[cl.Id] = cl
		//Unregister of client from room
		case cl := <-self.Unregister:
			if _, roomExists := self.Rooms[cl.RoomId]; roomExists {
				if _, clientExists := self.Rooms[cl.RoomId].Clients[cl.Id]; clientExists {
					delete(self.Rooms[cl.RoomId].Clients, cl.Id)
					close(cl.Message)
				}
			}
		//Braodcasting
		case m := <-self.Broadcast:
			if _, roomExists := self.Rooms[m.RoomId]; roomExists {
				for _, cl := range self.Rooms[m.RoomId].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
