package themes

import (
	"turtle/core/users"
	"turtle/db"
	"turtle/lg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CT_THEMES = "themes"

func ListThemes() []ThemeLight {
	opts := options.FindOptions{}
	opts.Projection = ThemeLightProjection

	return db.QueryEntitiesAsCopy[ThemeLight](CT_THEMES, bson.M{}, &opts)
}

func GetTheme(user *users.User, uid primitive.ObjectID) *Theme {
	return db.GetByIdAndOrg[Theme](CT_THEMES, uid, user.Org)
}

func ImportTheme(user *users.User, theme *Theme) {
	theme.Org = user.Org
	theme.Default = false

	if db.EntityExists(CT_THEMES, bson.M{"_id": theme.Uid}) {
		db.SetById(CT_THEMES, theme.Uid, theme)
		lg.LogE("Here")
	} else {
		lg.LogOk("Tu som")
		db.InsertEntity(CT_THEMES, theme)
	}
}

func COUTheme(user *users.User, theme *Theme) {
	theme.Org = user.Org

	if theme.Default {
		db.SetEntitiesWhere(CT_THEMES,
			bson.M{"org": user.Org},
			bson.M{"default": false})
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
