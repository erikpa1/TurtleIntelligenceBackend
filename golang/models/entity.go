package models

type Entity struct {
	Id       int       `json:"id" gorm:"primaryKey;autoIncrement;columnt:id"`
	Name     string    `json:"name" gorm:"index;columnt:name"`
	Position []float32 `json:"position" gorm:"columnt:position"`

	Behaviours map[int]*Behaviour `json:"behaviours" gorm:"-"`
}

func NewEntity() *Entity {
	return &Entity{}
}
