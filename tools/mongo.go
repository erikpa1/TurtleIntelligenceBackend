package tools

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"turtle/lg"
)

func RecastBson[T any](obj any) *T {
	// Marshal the bson.M back to BSON bytes
	data, err := bson.Marshal(obj)
	if err != nil {
		lg.LogE(fmt.Sprintf("failed to unmarshal TypeData: %w", err))
	}

	target := new(T)
	// Unmarshal into the target type
	if err := bson.Unmarshal(data, target); err != nil {
		lg.LogE(fmt.Sprintf("failed to unmarshal TypeData: %w", err))
	}

	return target
}
