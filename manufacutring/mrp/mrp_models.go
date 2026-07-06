package mrp

// PlannedOrder is a supply recommendation produced by the MRP run: either a
// production order (for "make" items, later scheduled by APS) or a purchase
// order (for "buy" items).
type PlannedOrder struct {
	ItemUid     string  `json:"itemUid"`
	Sku         string  `json:"sku"`
	Name        string  `json:"name"`
	OrderType   string  `json:"orderType"` // production | purchase
	Quantity    float64 `json:"quantity"`
	Uom         string  `json:"uom"`
	ReleaseDate string  `json:"releaseDate"`
	DueDate     string  `json:"dueDate"`
	BomLevel    int     `json:"bomLevel"`
}

// RequirementRow is the netting calculation for a single item, exposed so the
// UI can render the classic MRP table.
type RequirementRow struct {
	ItemUid           string  `json:"itemUid"`
	Sku               string  `json:"sku"`
	Name              string  `json:"name"`
	BomLevel          int     `json:"bomLevel"`
	GrossRequirement  float64 `json:"grossRequirement"`
	OnHand            float64 `json:"onHand"`
	SafetyStock       float64 `json:"safetyStock"`
	ScheduledReceipts float64 `json:"scheduledReceipts"`
	NetRequirement    float64 `json:"netRequirement"`
	PlannedOrder      float64 `json:"plannedOrder"`
	ProcurementType   string  `json:"procurementType"`
	RequiredDate      string  `json:"requiredDate"`
}

type MrpException struct {
	Severity string `json:"severity"` // warning | error
	ItemUid  string `json:"itemUid"`
	Sku      string `json:"sku"`
	Message  string `json:"message"`
}

// MrpResult is the full output of a planning run.
type MrpResult struct {
	GeneratedAt   string           `json:"generatedAt"`
	PlannedOrders []PlannedOrder   `json:"plannedOrders"`
	Requirements  []RequirementRow `json:"requirements"`
	Exceptions    []MrpException   `json:"exceptions"`
}
