package aps

// ScheduledOperation is one routing operation placed on the timeline for a
// planned production order.
type ScheduledOperation struct {
	OrderRef       string  `json:"orderRef"`
	ItemUid        string  `json:"itemUid"`
	Sku            string  `json:"sku"`
	OperationName  string  `json:"operationName"`
	Sequence       int     `json:"sequence"`
	WorkCenterUid  string  `json:"workCenterUid"`
	WorkCenterName string  `json:"workCenterName"`
	Quantity       float64 `json:"quantity"`
	Start          string  `json:"start"` // RFC3339
	End            string  `json:"end"`   // RFC3339
	DurationHours  float64 `json:"durationHours"`
	DueDate        string  `json:"dueDate"`
	Late           bool    `json:"late"`
}

// WorkCenterLoad summarises how heavily a work center is loaded across the
// scheduling horizon.
type WorkCenterLoad struct {
	WorkCenterUid  string  `json:"workCenterUid"`
	WorkCenterName string  `json:"workCenterName"`
	LoadHours      float64 `json:"loadHours"`
	CapacityHours  float64 `json:"capacityHours"`
	Utilization    float64 `json:"utilization"` // percent
	Operations     int     `json:"operations"`
}

// ApsResult is the finite-capacity schedule.
type ApsResult struct {
	GeneratedAt     string               `json:"generatedAt"`
	HorizonStart    string               `json:"horizonStart"`
	HorizonEnd      string               `json:"horizonEnd"`
	Operations      []ScheduledOperation `json:"operations"`
	WorkCenterLoads []WorkCenterLoad     `json:"workCenterLoads"`
	// Unscheduled lists production orders that could not be scheduled because
	// no routing was found.
	Unscheduled []string `json:"unscheduled"`
}
