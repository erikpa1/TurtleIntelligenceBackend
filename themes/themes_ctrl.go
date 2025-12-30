package themes

import (
	"turtle/core/users"
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CT_THEMES = "themes"

func ListThemes() []ThemeLight {
	opts := options.FindOptions{}
	opts.Projection = bson.M{
		"_id":     1,
		"name":    1,
		"default": 1,
	}

	return db.QueryEntitiesAsCopy[ThemeLight](CT_THEMES, bson.M{}, &opts)
}

func GetTheme(user *users.User, uid primitive.ObjectID) *Theme {
	return db.GetByIdAndOrg[Theme](CT_THEMES, uid, user.Org)
}

func COUTheme(user *users.User, theme *Theme) {
	theme.Org = user.Org

	if theme.Default {
		db.UpdateEntitiesWhere(CT_THEMES,
			bson.M{"org": user.Org},
			bson.M{"$set": bson.M{"default": false}})
	}

	if theme.Uid.IsZero() {
		db.InsertEntity(CT_THEMES, theme)
	} else {
		db.SetByOrgAndId(CT_THEMES, theme.Uid, theme.Org, theme)
	}
}

func GetDefaultTheme(user *users.User) *Theme {
	return db.QueryEntity[Theme](CT_THEMES, user.FillOrgQuery(bson.M{"default": true}))
}

func DeleteTheme(user *users.User, uid primitive.ObjectID) {
	db.DeleteEntity(CT_THEMES, user.FillOrgQuery(bson.M{"_id": uid}))
}
