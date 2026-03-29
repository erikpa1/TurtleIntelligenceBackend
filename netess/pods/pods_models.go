package pods

import (
	"turtle/core/entities"
)

type NetessPod struct {
	entities.EntityMinimal `json:",inline" bson:",inline"`
	RestConfig             NetessRestConfig `json:"restConfig" bson:"restConfig"`
}

type NetessRestConfig struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}
