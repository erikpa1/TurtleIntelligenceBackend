package simInternal

import "github.com/erikpa1/TurtleIntelligenceBackend/tools"

type SimUpcomingEventType int8

const (
	UPC_EVNT_SPAWN    SimUpcomingEventType = 0
	UPC_EVENT_UNSPAWN SimUpcomingEventType = 1
	UPC_EVENT_START   SimUpcomingEventType = 2
	UPC_EVENT_FINISH  SimUpcomingEventType = 3
)

type SimUpcomingEvent struct {
	Id     int64                `json:"id"`
	Type   SimUpcomingEventType `json:"type"`
	Second tools.Seconds        `json:"second"`
}
