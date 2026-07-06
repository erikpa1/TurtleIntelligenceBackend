package items

import "go.mongodb.org/mongo-driver/bson/primitive"

// Item is a material master record (a product / part) inspired by the SAP
// material master. It is the atomic entity used both for stock keeping in the
// inventory module and as a component/product inside manufacturing BOMs.
type Item struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org         primitive.ObjectID `json:"org" bson:"org"`
	Sku         string             `json:"sku" bson:"sku"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`

	// Category is the high level classification used by manufacturing:
	// raw, semi-finished, finished, trading.
	Category string `json:"category,omitempty" bson:"category,omitempty"`
	// Type is a free form material type (e.g. "Material", "Service").
	Type string `json:"type,omitempty" bson:"type,omitempty"`
	// Uom is the base unit of measure (pcs, kg, m, l ...).
	Uom string `json:"uom,omitempty" bson:"uom,omitempty"`

	UnitPrice float64 `json:"unitPrice" bson:"unitPrice"`
	Currency  string  `json:"currency,omitempty" bson:"currency,omitempty"`

	QtyOnHand    float64 `json:"qtyOnHand" bson:"qtyOnHand"`
	ReorderPoint float64 `json:"reorderPoint" bson:"reorderPoint"`
	Warehouse    string  `json:"warehouse,omitempty" bson:"warehouse,omitempty"`

	// MRP / planning view (SAP material master MRP view). These drive the
	// requirements planning run in the manufacturing/mrp module.
	//
	// ProcurementType is "make" (produced in-house, exploded through its BOM)
	// or "buy" (procured externally). Empty is inferred from the presence of a
	// BOM at planning time.
	ProcurementType string  `json:"procurementType,omitempty" bson:"procurementType,omitempty"`
	LeadTimeDays    float64 `json:"leadTimeDays" bson:"leadTimeDays"`
	SafetyStock     float64 `json:"safetyStock" bson:"safetyStock"`
	// LotSize is a fixed rounding value; 0 means lot-for-lot (order the exact
	// net requirement).
	LotSize    float64 `json:"lotSize" bson:"lotSize"`
	MinLotSize float64 `json:"minLotSize" bson:"minLotSize"`

	Active bool `json:"active" bson:"active"`
}
