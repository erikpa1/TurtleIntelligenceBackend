package themes

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThemeLight struct {
	Uid     primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name    string             `json:"name"`
	Default bool               `json:"default"`
	Org     primitive.ObjectID `json:"org"`
	Color   string             `json:"color"`
}

type Theme struct {
	ThemeLight `bson:",inline"`
	ThemeData  `bson:",inline"`
}

type ThemeData struct {
	TopBarHeightBig    string `json:"topBarHeightBig,omitempty" bson:"topBarHeightBig,omitempty"`
	BigPadding         string `json:"bigPadding,omitempty" bson:"bigPadding,omitempty"`
	HeadingFontColor   string `json:"headingFontColor,omitempty" bson:"headingFontColor,omitempty"`
	IconPrimaryColor   string `json:"iconPrimaryColor,omitempty" bson:"iconPrimaryColor,omitempty"`
	IconSecondaryColor string `json:"iconSecondaryColor,omitempty" bson:"iconSecondaryColor,omitempty"`
	BorderColor        string `json:"borderColor,omitempty" bson:"borderColor,omitempty"`
	BorderHoverColor   string `json:"borderHoverColor,omitempty" bson:"borderHoverColor,omitempty"`
}
