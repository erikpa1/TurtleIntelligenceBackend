package auth

import (
	"fmt"
	"github.com/erikpa1/turtle/credentials"
	"github.com/erikpa1/turtle/ctrl"
	"github.com/erikpa1/turtle/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthContext struct {
	GinC *gin.Context
	User *models.User
}

func GetUserFromContext(c *gin.Context) *models.User {
	tmp := models.NewSuperAdmin()

	return tmp
}

func GetUserFromCookies(cookiestring string) *models.User {
	return &models.User{}
}

func LoginRequiredWithUser(fn func(ctx *AuthContext)) gin.HandlerFunc {

	return func(c *gin.Context) {
		str_token, err := c.Cookie("")

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"why": "no auth token"})
		} else {
			_, err := jwt.Parse(str_token, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(credentials.AuthInfinityJwtSecret()), nil
			})

			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"why": "Wrong auth token"})
				return
			}

			ctx := AuthContext{}
			ctx.GinC = c
			ctx.User = models.NewUser()
			fn(&ctx)
		}

	}
}

func IsLoggedInInfinity(c *gin.Context) (*models.User, string) {
	tokenName := credentials.AuthInfinityJwtKey()

	token, err := c.Cookie(tokenName)
	if err != nil {
		token = "" // If the cookie is not found, use an empty string
	}

	// Check the token and retrieve the user
	user, _ := ctrl.CheckInfinityAuth(token)

	if user != nil {
		tokenName := credentials.AuthInfinityJwtKey()

		var expireTime, _ = strconv.Atoi(credentials.AuthInfinityTokenExpire())

		c.SetCookie(tokenName, token, expireTime, "/", "", false, true)

		return user, token
	}

	c.SetCookie(tokenName, "", -1, "/", "", false, true)

	return nil, token
}

type IdpImplementation interface {
	RedirectToIdp(c *gin.Context)
}

func LoginRequired(c *gin.Context) {
	if credentials.AuthProvider() == credentials.AUTH_PROVIDER_NONE {
		c.Next()

		return
	}

	if credentials.AuthProvider() == credentials.AUTH_PROVIDER_INFINITY || credentials.AuthProvider() == credentials.AUTH_PROVIDER_ALL {
		user, _ := IsLoggedInInfinity(c)

		if user != nil {
			c.Next()

			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"why": "Unauthorized"})
}

func AdminRequired(c *gin.Context) {
	c.Next()
}

func LoginOrApp(c *gin.Context) {
	c.Next()
}
