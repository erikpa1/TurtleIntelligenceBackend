package bom

import "go.mongodb.org/mongo-driver/bson/primitive"

// BomComponent is a single line of a bill of materials: an inventory item that
// is consumed (in a given quantity) to produce the BOM's product. References to
// items are stored as hex uid strings so a component can reference an item
// without coupling the BOM document to the item lifecycle.
type BomComponent struct {
	ItemUid  string  `json:"itemUid" bson:"itemUid"`
	Sku      string  `json:"sku,omitempty" bson:"sku,omitempty"`
	Name     string  `json:"name,omitempty" bson:"name,omitempty"`
	Quantity float64 `json:"quantity" bson:"quantity"`
	Uom      string  `json:"uom,omitempty" bson:"uom,omitempty"`
	ScrapPct float64 `json:"scrapPct" bson:"scrapPct"`
	Position int     `json:"position" bson:"position"`
}

// Bom is a manufacturing bill of materials header describing how many
// components are needed to build a given quantity of a finished product.
type Bom struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org         primitive.ObjectID `json:"org" bson:"org"`
	Code        string             `json:"code" bson:"code"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`

	// ProductUid / ProductSku identify the inventory item this BOM produces.
	ProductUid string `json:"productUid,omitempty" bson:"productUid,omitempty"`
	ProductSku string `json:"productSku,omitempty" bson:"productSku,omitempty"`

	// BaseQuantity is the number of product units produced by consuming the
	// component quantities below (SAP "base quantity").
	BaseQuantity float64 `json:"baseQuantity" bson:"baseQuantity"`
	Uom          string  `json:"uom,omitempty" bson:"uom,omitempty"`

	Version string `json:"version,omitempty" bson:"version,omitempty"`
	// Status is one of: draft, active, obsolete.
	Status string `json:"status,omitempty" bson:"status,omitempty"`

	Components []BomComponent `json:"components" bson:"components"`
}
