package ctrl

import (
	"errors"
	"fmt"
	"turtle/core/users"
	"turtle/credentials"
	"turtle/db"
	"turtle/lgr"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/scrypt"
)

var CT_USERS = "users"

func GetUsers() []*users.User {
	opts := options.FindOptions{}
	opts.Projection = bson.M{"password": 0}
	return db.QueryEntities[users.User](CT_USERS, bson.M{}, &opts)
}

func GetUser(uid primitive.ObjectID) *users.User {
	return db.QueryEntity[users.User](CT_USERS, bson.M{"uid": uid})
}

func CheckInfinityAuth(token string) (*users.User, error) {
	if token == "" {
		return nil, nil // No token provided
	}

	// Parse the token without verification to extract the algorithm
	unverifiedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, errors.New("failed to parse token")
	}

	// Extract the signing algorithm from the token header
	header, ok := unverifiedToken.Header["alg"].(string)
	if !ok {
		return nil, errors.New("invalid token header")
	}

	// Verify and decode the token
	decoded, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Ensure the algorithm matches the header's algorithm
		if token.Method.Alg() != header {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(credentials.AuthInfinityJwtSecret()), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Extract claims from the decoded token
	claims, ok := decoded.Claims.(jwt.MapClaims)
	if !ok || !decoded.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Retrieve the UID from the claims
	uid, ok := claims["uid"].(string)
	if !ok || uid == "" {
		return nil, errors.New("uid not found in token")
	}

	userId, _ := primitive.ObjectIDFromHex(uid)

	// Get the user by UID
	user := GetUser(userId)
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func GetUserByEmail(email string) *users.User {
	return db.QueryEntity[users.User](CT_USERS, bson.M{"email": email})
}

func EncryptPassword(password string) string {
	// Define the salt (in your case, password + "5zg")
	salt := []byte(password + "5zg")

	// Scrypt parameters: N=16384, r=8, p=1
	const N = 16384
	const r = 8
	const p = 1
	keyLen := 32 // Output key length

	// Derive the key using scrypt
	hashedPassword, err := scrypt.Key([]byte(password), salt, N, r, p, keyLen)
	if err != nil {
		lgr.Error("failed to generate scrypt key: %v", err)
	}

	// Convert the result to a hex string
	return fmt.Sprintf("%x", hashedPassword)
}

func UserExists(email string, password string) bool {
	hash := EncryptPassword(password)
	entity := db.QueryEntity[users.User](CT_USERS, bson.M{"email": email, "password": hash})
	return entity != nil
}

func COUUser(user *users.User) {

	//TODO toto je nejake divne
	from_dbuser := GetUser(user.Uid)

	if from_dbuser != nil {
		from_dbuser.FromAnotherUserNoPass(user)

		if from_dbuser.Uid.IsZero() {
			db.InsertEntity(CT_USERS, from_dbuser)
		} else {
			db.UpdateOneCustom(CT_USERS, bson.M{"_id": user.Uid}, bson.M{"$set": from_dbuser})
		}

	} else {
		user.Password = EncryptPassword(user.Password)
		db.InsertEntity(CT_USERS, user)
	}
}

func DeleteUser(uid string) {
	db.DeleteEntity(CT_USERS, bson.M{"uid": uid})
}

func ChangePassword(session string, password string) {
	//TODO
}

func CheckSession(session string) {
	//TODO
}
