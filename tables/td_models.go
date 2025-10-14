package tables

import "go.mongodb.org/mongo-driver/bson/primitive"

type TableData struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org       primitive.ObjectID `json:"org"`
	Name      string             `json:"name"`
	Theme     TableTheme         `json:"theme"`
	Schema    TableDataSchema    `json:"schema"`
	BatchSize int                `json:"batchSize" bson:"batchSize"`
}

type TableDataSchema struct {
	Uid   primitive.ObjectID  `json:"uid" bson:"_id,omitempty"`
	Org   primitive.ObjectID  `json:"org"`
	Items TableDataSchemaItem `json:"items"`
}

type TableDataSchemaItemType int8

const (
	TD_SCH_T_STRING TableDataSchemaItemType = 0
	TD_SCH_T_F64    TableDataSchemaItemType = 1
	TD_SCH_T_F32    TableDataSchemaItemType = 2
)

type TableDataSchemaItem struct {
	Name string
	Type TableDataSchemaItemType
}

type TableTheme struct {
	HeaderColor string `json:"headerColor" bson:"headerColor"`
}
