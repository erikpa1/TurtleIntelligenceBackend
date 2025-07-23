package flows

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CT_FLOWS = "flows"
const CT_FLOWS_EXECUTIONS = "flow_execution"

func DeleteFlow(user *models.User, flowUid primitive.ObjectID) {
	db.DeleteEntities(CT_FLOWS_EXECUTIONS, bson.M{"flow": flowUid})

	db.DeleteEntity(CT_FLOWS,
		bson.M{
			"_id": flowUid,
			"org": user.Org,
		},
	)

}

func GetFlow(user *models.User, flowUid primitive.ObjectID) *Flow {
	return db.QueryEntity[Flow](CT_FLOWS,
		bson.M{
			"_id": flowUid,
			"org": user.Org,
		})
}

func ListFlows(user *models.User) []*Flow {
	opts := options.FindOptions{}
	opts.Projection = bson.M{"_id": 1, "name": 1}

	return db.QueryEntities[Flow](CT_FLOWS, bson.M{"org": user.Org}, &opts)
}

func COUFlow(user *models.User, flow *Flow) {
	flow.Org = user.Org

	if flow.Uid.IsZero() {
		db.InsertEntity(CT_FLOWS,
			flow,
		)
	} else {
		db.UpdateOneCustom(CT_FLOWS, bson.M{
			"_id": flow.Uid,
			"org": user.Org,
		}, bson.M{"$set": flow})
	}

}

func CallFlow(user *models.User, flowUid primitive.ObjectID) (int, string) {
	flowExe := FlowExecution{}
	flowExe.Uid = primitive.NewObjectID()
	flowExe.Flow = flowUid
	flowExe.At = tools.GetNow()
	flowExe.Org = user.Org
	flowExe.Status = 0

	db.InsertEntity(CT_FLOWS_EXECUTIONS, flowExe)

	flow := GetFlow(user, flowUid)

	updateStatus := 404
	callError := ""

	if flow != nil {

	}

	db.UpdateOneCustom(CT_FLOWS_EXECUTIONS, bson.M{
		"_id": flowExe.Uid,
		"org": user.Org,
	},
		bson.M{
			"$set": bson.M{
				"status":    updateStatus,
				"callstack": callError,
			},
		},
	)

	return updateStatus, callError

}
