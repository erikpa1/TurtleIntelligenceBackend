package tools

//Copied on https://github.com/erikpa1/TurtleSim/blob/main/TurtleSimCpp/src/serialization/safejson.cpp

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// SafeJson provides safe methods to parse, access, and modify JSON Data.
type SafeJson struct {
	Data          map[string]interface{}
	OptimizeSpace bool
}

// NewSafeJson initializes a SafeJson instance.
func NewSafeJson() *SafeJson {
	return &SafeJson{
		Data:          make(map[string]interface{}),
		OptimizeSpace: true,
	}
}

// ParseString parses a JSON string into the SafeJson's internal map.
func (s *SafeJson) ParseString(str string) error {
	err := json.Unmarshal([]byte(str), &s.Data)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return err
	}
	return nil
}

// GetString retrieves a string value or a default if the key is not found.
func (s *SafeJson) GetString(key, notFound string) string {
	if val, ok := s.Data[key].(string); ok {
		return val
	}
	return notFound
}

// GetInt retrieves an int value or a default if the key is not found.
func (s *SafeJson) GetInt(key string, notFound int) int {
	if val, ok := s.Data[key].(float64); ok {
		return int(val)
	}
	return notFound
}

// GetInt retrieves an int value or a default if the key is not found.
func (s *SafeJson) GetInt64(key string, notFound int64) int64 {
	if val, ok := s.Data[key].(float64); ok {
		return int64(val)
	}
	return notFound
}

// GetLong retrieves a long (int64) value or a default if the key is not found.
func (s *SafeJson) GetLong(key string, notFound int64) int64 {
	if val, ok := s.Data[key].(float64); ok {
		return int64(val)
	}
	return notFound
}

// GetDouble retrieves a double (float64) value or a default if the key is not found.
func (s *SafeJson) GetDouble(key string, notFound float64) float64 {
	if val, ok := s.Data[key].(float64); ok {
		return val
	}
	return notFound
}

func (s *SafeJson) GetObjectArray(key string) []*SafeJson {
	//Tento kus kodu tu musi byt kvoli subtypovaniu golangovych struktur

	value := reflect.ValueOf(s.Data[key])

	if value.Kind() == reflect.Slice {
		result := make([]*SafeJson, 0, value.Len())

		// Iterate over the slice elements.
		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)

			// Check if the item is a pointer to a map or can be converted to one.
			// item.Kind() == reflect.Ptr && toto doporucil chat gpt a je to blbost
			if item.Elem().Kind() == reflect.Map {
				// Convert item to map[string]interface{}
				mapItem, ok := item.Interface().(*map[string]interface{})
				if !ok {
					// If it's a custom type, try converting to map[string]interface{} with reflection.
					if item.Elem().Type().ConvertibleTo(reflect.TypeOf(map[string]interface{}{})) {
						converted := item.Elem().Convert(reflect.TypeOf(map[string]interface{}{})).Interface()
						convertedMap := converted.(map[string]interface{})

						// Wrap in SafeJson and append to result.
						jsonObj := NewSafeJson()
						jsonObj.Data = convertedMap
						result = append(result, jsonObj)
					}
				} else {
					// If the item is already *map[string]interface{}, use it directly.
					jsonObj := NewSafeJson()
					jsonObj.Data = *mapItem
					result = append(result, jsonObj)
				}
			}
		}
		return result
	}
	return nil
}

func (s *SafeJson) GetFloats64Array(key string, defaultValue []float64) []float64 {
	result := make([]float64, 3)

	value := reflect.ValueOf(s.Data[key])

	if value.Kind() == reflect.Slice {

		result = make([]float64, value.Len())

		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)

			kind := item.Elem().Kind()

			if kind == reflect.Float64 ||
				kind == reflect.Float32 ||
				kind == reflect.Int ||
				kind == reflect.Int32 ||
				kind == reflect.Int16 ||
				kind == reflect.Int64 {
				result[i] = item.Elem().Interface().(float64)
			}
		}
	}

	return result
}

func (s *SafeJson) GetStringsArray(key string) []string {

	result := make([]string, 0)

	value := reflect.ValueOf(s.Data[key])

	if value.Kind() == reflect.Slice {

		result = make([]string, value.Len())

		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)
			if item.Elem().Kind() == reflect.String {
				result[i] = item.Elem().Interface().(string)
			}
		}
	}

	return result
}

func (s *SafeJson) GetFloatsPositionArray(key string, defaultValue [3]float64) [3]float64 {
	result := defaultValue

	value := reflect.ValueOf(s.Data[key])

	if value.Kind() == reflect.Slice {

		len := 3
		if value.Len() < len {
			len = value.Len()
		}

		for i := 0; i < len; i++ {
			item := value.Index(i)

			kind := item.Elem().Kind()

			switch kind {
			case reflect.Float64:
				result[i] = item.Elem().Interface().(float64)
			case reflect.Float32:
				result[i] = float64(item.Elem().Interface().(float32))
			case reflect.Int:
				result[i] = float64(item.Elem().Interface().(int))
			case reflect.Int32:
				result[i] = float64(item.Elem().Interface().(int32))
			case reflect.Int16:
				result[i] = float64(item.Elem().Interface().(int16))
			case reflect.Int64:
				result[i] = float64(item.Elem().Interface().(int64))
			default:
				// handle other cases or unexpected types
				result[i] = 0
			}
		}
	}

	return result
}

// WriteString writes a string value to the JSON.
func (s *SafeJson) WriteString(key, value string) {
	if !s.OptimizeSpace || value != "" {
		s.Data[key] = value
	}
}

// WriteFloat writes a float value to the JSON.
func (s *SafeJson) WriteFloat(key string, value float32) {
	if !s.OptimizeSpace || value != 0.0 {
		s.Data[key] = value
	}
}

// WriteFloat3 writes an array of three float values to the JSON.
func (s *SafeJson) WriteFloat3(key string, x, y, z float32) {
	s.Data[key] = []float32{x, y, z}
}

// WriteDouble writes a double (float64) value to the JSON.
func (s *SafeJson) WriteDouble(key string, value float64) {
	if !s.OptimizeSpace || value != 0 {
		s.Data[key] = value
	}
}

// WriteDouble3 writes an array of three double values to the JSON.
func (s *SafeJson) WriteDouble3(key string, x, y, z float64) {
	s.Data[key] = []float64{x, y, z}
}

// WriteInt writes an int value to the JSON.
func (s *SafeJson) WriteInt(key string, value int) {
	if !s.OptimizeSpace || value != 0 {
		s.Data[key] = value
	}
}

// WriteLong writes a long (int64) value to the JSON.
func (s *SafeJson) WriteLong(key string, value int64) {
	if !s.OptimizeSpace || value != 0 {
		s.Data[key] = value
	}
}

// WriteBoolean writes a boolean value to the JSON.
func (s *SafeJson) WriteBoolean(key string, value bool) {
	if !s.OptimizeSpace || value != false {
		s.Data[key] = value
	}
}

// WriteJson writes another SafeJson object as a nested JSON object.
func (s *SafeJson) WriteJson(key string, another *SafeJson) {
	if another != nil {
		s.Data[key] = another.Data
	}
}

// WriteJsonArray writes an array of SafeJson objects to the JSON.
func (s *SafeJson) WriteJsonArray(key string, array []*SafeJson) {
	if s.OptimizeSpace && len(array) == 0 {
		return
	}
	var jsonArray []map[string]interface{}
	for _, item := range array {
		jsonArray = append(jsonArray, item.Data)
	}
	s.Data[key] = jsonArray
}

// GetKeysLength returns the number of top-level keys in the JSON.
func (s *SafeJson) GetKeysLength() int {
	return len(s.Data)
}

// Dump returns the JSON as a pretty-printed string.
func (s *SafeJson) Dump() string {
	bytes, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
