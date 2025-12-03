package loginPenetration

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BruteForceResult struct {
	Success      bool
	Password     string
	Attempts     int64
	Duration     time.Duration
	StatusCode   int
	ResponseBody string
}

type LoginPenetration struct {
	Uid             primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name            string             `json:"name"`
	Url             string             `json:"url"`
	Email           string             `json:"email"`
	IterationsCount int64              `json:"iteratotionsCount" bson:"iterationsCount"`
}
