package themes

import (
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CT_THEMES = "themes"

func ListThemes() []bson.M {
	opts := options.FindOptions{}
	opts.Projection = bson.M{"name": 1}
	opts.Projection = bson.M{"_id": 1}

	return db.QueryEntitiesAsCopy[bson.M](CT_THEMES, bson.M{}, &opts)
}

func GetTheme(uid string) bson.M {
	return bson.M{}
}
