package dbw

import (
	"go.mongodb.org/mongo-driver/bson"
	"turtle/lg"
	"turtle/models"
)

func QueryOrgEntities(user *models.User) {
	query := bson.M{
		"org": user.Org,
	}

	lg.LogE(query)
}

func QueryUserData(user *models.User) {
	query := bson.M{
		"org":  user.Org,
		"user": user.Uid,
	}

	lg.LogE(query)

}
