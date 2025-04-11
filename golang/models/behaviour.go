package models

type Behaviour struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	TypeData string `json:"type_data"`
}

type IBehaviour interface{}
