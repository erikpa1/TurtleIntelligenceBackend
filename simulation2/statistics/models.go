package statistics

type SimulationState int8

const (
	SIM_STATE_LOADING  SimulationState = 0
	SIM_STATE_LOADED   SimulationState = 1
	SIM_STATE_STARTED  SimulationState = 2
	SIM_STATE_PAUSED   SimulationState = 3
	SIM_STATE_FAILED   SimulationState = 4
	SIM_STATE_STOPPED  SimulationState = 5
	SIM_STATE_FINISHED SimulationState = 6
)
