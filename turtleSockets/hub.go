/*
	Inspired by tutorial:

https://www.youtube.com/watch?v=760GKM7s_5Y
*/
package turtleSockets

type Room struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Clients map[string]*Client
}

type Hub struct {
	Rooms map[string]*Room `json:"rooms"`
}
