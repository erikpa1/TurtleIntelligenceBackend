package themes

import (
	"go.mongodb.org/mongo-driver/bson"
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
	Title              string `json:"title,omitempty" bson:"title,omitempty"`
	Favicon            string `json:"favicon,omitempty" bson:"favicon,omitempty"`
	PrimaryLogo        string `json:"primaryLogo,omitempty" bson:"primaryLogo,omitempty"`
	PrimaryLogoSizeX   string `json:"primaryLogoSizeX,omitempty" bson:"primaryLogoSizeX,omitempty"`
	PrimaryLogoSizeY   string `json:"primaryLogoSizeY,omitempty" bson:"primaryLogoSizeY,omitempty"`
	TopBarHeightBig    string `json:"topBarHeightBig,omitempty" bson:"topBarHeightBig,omitempty"`
	BigPadding         string `json:"bigPadding,omitempty" bson:"bigPadding,omitempty"`
	HeadingFontColor   string `json:"headingFontColor,omitempty" bson:"headingFontColor,omitempty"`
	IconPrimaryColor   string `json:"iconPrimaryColor,omitempty" bson:"iconPrimaryColor,omitempty"`
	IconSecondaryColor string `json:"iconSecondaryColor,omitempty" bson:"iconSecondaryColor,omitempty"`
	BorderColor        string `json:"borderColor,omitempty" bson:"borderColor,omitempty"`
	BorderHoverColor   string `json:"borderHoverColor,omitempty" bson:"borderHoverColor,omitempty"`
	PrimaryColor       string `json:"primaryColor,omitempty" bson:"primaryColor,omitempty"`
}

var ThemeLightProjection = bson.M{
	"_id":     1,
	"name":    1,
	"default": 1,
	"color":   1,
}
