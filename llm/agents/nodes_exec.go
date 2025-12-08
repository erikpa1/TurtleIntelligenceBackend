package agents

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ExeNodeWithUid(orgUid primitive.ObjectID, uid primitive.ObjectID) {

	node := GetAgentNode(orgUid, uid)

	if node != nil {
		relation := GetRelationOfNode(orgUid, bson.M{"b": uid})

		if relation != nil {
			_ExecuteNodeRecursive(orgUid, relation.B)
		}
	}
}

func _ExecuteNodeRecursive(orgUid primitive.ObjectID, uid primitive.ObjectID) {
	node := GetAgentNode(orgUid, uid)
	lg.LogEson(node)
	//TODO sem treba dorobit aby sa nody citali rekurzivne

}
