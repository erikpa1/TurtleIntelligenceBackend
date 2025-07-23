package dbw

import (
	"github.com/erikpa1/turtle/lg"
	"github.com/erikpa1/turtle/models"
	"go.mongodb.org/mongo-driver/bson"
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
