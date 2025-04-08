package tools

import (
	"fmt"
	"strconv"
	"turtle/lg"
)

func ToIArray[T any](data []T) []interface{} {
	interfaceArray := make([]interface{}, len(data))

	for i, element := range data {
		interfaceArray[i] = element
	}
	return interfaceArray
}

func ToAnotherArray[SRC any, DEST any](data []SRC) ([]DEST, error) {
	interfaceArray := make([]DEST, len(data))

	for i, element := range data {
		destElement, ok := any(element).(DEST)
		if !ok {
			return nil, fmt.Errorf("element at index %d cannot be converted to type", i)
		}
		interfaceArray[i] = destElement
	}
	return interfaceArray, nil
}

func MineValueFromDictList[MAP_T comparable, DATA_T any](data []map[MAP_T]interface{}, key MAP_T) []DATA_T {
	result := make([]DATA_T, len(data))

	for i, element := range data {
		result[i] = element[key].(DATA_T)
	}

	return result

}

func F64fromString(value string) float64 {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		lg.LogE(f)
		return 0
	}
	return f
}

func F32fromString(value string) float32 {
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		lg.LogE(f)
		return 0
	}
	return float32(f)
}
func Int64fromString(value string) int64 {
	f, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		lg.LogE(f)
		return 0
	}
	return f
}
