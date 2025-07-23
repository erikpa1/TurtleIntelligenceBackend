package nn

import (
	"github.com/erikpa1/turtle/db"
	"github.com/erikpa1/turtle/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_NN = "nn"
const CT_NN_CONFIGS = "nn_config"
const CT_NN_EXPERIMENTS = "nn_experiments"
const CT_NN_EXPERIMENT_RESULTS = "nn_experiment_result"

func ListNN(user *models.User) []*NeuralNetwork {
	return db.QueryEntities[NeuralNetwork](CT_NN, bson.M{"org": user.Org})
}

func COUNN(user *models.User, nn *NeuralNetwork) {
	nn.Org = user.Org

	if nn.Uid.IsZero() {
		db.InsertEntity(CT_NN, nn)
	} else {
		db.UpdateOneCustom(CT_NN, bson.M{"_id": nn.Uid}, bson.M{"$set": nn})
	}
}

func DeleteNN(user *models.User, nnUid primitive.ObjectID) {

	nnQuery := bson.M{
		"_id": nnUid,
		"org": user.Org,
	}

	if db.EntityExists(CT_NN, nnQuery) {
		db.DeleteEntities(CT_NN_EXPERIMENT_RESULTS, bson.M{"parentNetwork": nnUid})
		db.DeleteEntities(CT_NN_EXPERIMENTS, bson.M{"parentNetwork": nnUid})
		db.DeleteEntities(CT_NN_CONFIGS, bson.M{"parentNetwork": nnUid})
		db.DeleteEntity(CT_NN, nnQuery)
	}

}
