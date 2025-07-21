package tools

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/lg"
)

func VecFromJStr[T any](data string) []T {

	if data == "" {
		return []T{}
	} else {
		var result []T
		err := json.Unmarshal([]byte(data), &result)

		if err != nil {
			lg.LogStackTraceErr(err, "JSON: ", data)
		}

		return result
	}

}

func ObjFromJsonPtr[T any](data string) *T {
	var result = new(T)
	err := json.Unmarshal([]byte(data), &result)

	if err != nil {
		lg.LogStackTraceErr(err, "JSON: ", data)
		return nil
	}

	return result
}

func QueryHeader[T any](c *gin.Context) T {
	var result T

	headerValue := c.GetHeader("query")
	err := json.Unmarshal([]byte(headerValue), &result)

	if err != nil {
		lg.LogStackTraceErr(err, "JSON: ", headerValue)
		return result
	}

	return result
}

func ObjFromBsonPtr[T any](data string) *T {

	dataDict := ObjFromJson[bson.M](data)

	// Convert map to BSON bytes
	bsonData, err := bson.Marshal(dataDict)
	if err != nil {
		lg.LogStackTraceErr(err)
	}

	var result T
	err = bson.Unmarshal(bsonData, &result)
	if err != nil {
		lg.LogStackTraceErr(err)
	}

	return &result
}

func ObjFromJson[T any](data string) T {
	var result T
	err := json.Unmarshal([]byte(data), &result)

	if err != nil {
		lg.LogStackTraceErr(err, "JSON: ", data)
		return result
	}

	return result
}

func SafeJsonFromJson(data string) *SafeJson {
	result := new(SafeJson)

	var tmp map[string]interface{}

	err := json.Unmarshal([]byte(data), &tmp)

	if err != nil {
		lg.LogStackTraceErr(err, "JSON: ", data)
		return nil
	}
	result.Data = tmp
	return result
}
