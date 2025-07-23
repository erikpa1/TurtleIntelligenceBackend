package nn

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NeuralNetwork struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org       primitive.ObjectID `json:"org"`
	CreatedAt tools.Milliseconds `json:"createdAt" bson:"createdAt"`
	Name      string             `json:"name"`
}

type NeuralNetworkConfig struct {
	Uid           primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org           primitive.ObjectID `json:"org"`
	ParentNetwork primitive.ObjectID `json:"parentNetwork" bson:"parentNetwork"`
	Name          string             `json:"name"`
}

type NeuralNetworkExperiment struct {
	Uid           primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org           primitive.ObjectID `json:"org"`
	ParentNetwork primitive.ObjectID `json:"parentNetwork" bson:"parentNetwork"`
	Config        primitive.ObjectID `json:"config"`
	Dataset       primitive.ObjectID `json:"dataset"`
	Name          string             `json:"name"`
}

type NeuralNetworkExperimentResult struct {
	Uid              primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org              primitive.ObjectID `json:"org"`
	ParentNetwork    primitive.ObjectID `json:"parentNetwork" bson:"parentNetwork"`
	ParentExperiment primitive.ObjectID `json:"parentNetwork" bson:"parentNetwork"`
	Status           float32            `json:"status"`
}
