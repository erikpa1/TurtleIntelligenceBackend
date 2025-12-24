package dbw

import (
	"turtle/core/users"
	"turtle/lg"

	"go.mongodb.org/mongo-driver/bson"
)

func QueryOrgEntities(user *users.User) {
	query := bson.M{
		"org": user.Org,
	}

	lg.LogE(query)
}

func QueryUserData(user *users.User) {
	query := bson.M{
		"org":  user.Org,
		"user": user.Uid,
	}

	lg.LogE(query)

}
