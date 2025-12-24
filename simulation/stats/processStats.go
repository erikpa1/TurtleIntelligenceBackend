package stats

import "turtle/tools"

type ProcessStats struct {
	IdleTime    tools.Seconds
	ProcessTime tools.Seconds
	BlockedTime tools.Seconds
}

func NewProcessStats() *ProcessStats {
	return &ProcessStats{}
}
