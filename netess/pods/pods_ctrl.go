package pods

import (
	"turtle/core/users"
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_PODS = "netess_pods"

func ListPods(user *users.User) []*NetessPod {
	return db.QueryEntities[NetessPod](CT_PODS, bson.M{})
}

func GetPod(user *users.User, podUid primitive.ObjectID) *NetessPod {
	return db.GetByIdAndOrg[NetessPod](CT_PODS, podUid, user.Org)
}

func COUPod(user *users.User, pod *NetessPod) {

	user.Org = user.Org

	if pod.Uid.IsZero() {
		db.InsertEntity(CT_PODS, pod)
	} else {
		db.SetByOrgAndId(CT_PODS, pod.Uid, user.Org, pod)
	}
}

func DeletePod(user *users.User, uid primitive.ObjectID) {
	db.DeleteByIdAndOrg(CT_PODS, uid, user.Org)
}
