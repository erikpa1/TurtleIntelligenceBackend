package stats

import "github.com/erikpa1/TurtleIntelligenceBackend/tools"

type ProcessStats struct {
	IdleTime    tools.Seconds
	ProcessTime tools.Seconds
	BlockedTime tools.Seconds
}

func NewProcessStats() *ProcessStats {
	return &ProcessStats{}
}
