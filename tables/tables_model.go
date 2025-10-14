package tables

import (
	"fmt"

	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TurtleTable struct {
	Uid  primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org  primitive.ObjectID `json:"org"`
	Name string             `json:"name"`

	CreatedAt tools.Milliseconds `json:"createdAt" bson:"createdAt"`
	UpdatedAt tools.Milliseconds `json:"updatedAt" bson:"updatedAt"`

	CreatedBy primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	UpdatedBy primitive.ObjectID `json:"updatedBy" bson:"updatedBy"`

	Headers          []string `json:"headers"`
	ValueTypes       []string `json:"valueTypes" bson:"valueTypes"`
	DefaultValues    []any    `json:"defaultValues" bson:"defaultValues"`
	HasDatabaseTable bool     `json:"hasDatabaseTable" bson:"hasDatabaseTable"`
}

func (self *TurtleTable) GetTableMongoName() string {
	return fmt.Sprintf("z_table_%s", self.Uid.Hex())
}
