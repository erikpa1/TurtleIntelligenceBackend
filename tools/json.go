package tools

//Copied on https://github.com/erikpa1/TurtleSim/blob/main/TurtleSimCpp/src/serialization/safejson.cpp

import (
	"encoding/json"
	"fmt"
	"github.com/erikpa1/turtle/lg"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

// SafeJson provides safe methods to parse, access, and modify JSON Data.
type SafeJson struct {
	Data          map[string]interface{}
	OptimizeSpace bool
}

func UnmarshalToSafeJson(str string) *SafeJson {

	data := make(map[string]interface{})

	err := json.Unmarshal([]byte(str), &data)

	if err != nil {
		lg.LogE(err.Error())
	}

	return &SafeJson{
		Data:          data,
		OptimizeSpace: true,
	}
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

// GetString retrieves a string value or a default if the key is not found.
func (s *SafeJson) GetPrimitiveObjectId(key string) primitive.ObjectID {
	if val, ok := s.Data[key].(string); ok {
		primitiv, _ := primitive.ObjectIDFromHex(val)

		return primitiv
	}
	return primitive.ObjectID{}
}

// GetString retrieves a string value or a default if the key is not found.
func (s *SafeJson) GetInterface(key string, notFound interface{}) interface{} {
	if val, ok := s.Data[key]; ok {
		return val
	}
	return notFound
}

// GetString retrieves a string value or a default if the key is not found.
func (s *SafeJson) GetStringMatrix(key string) [][]string {
	data, okData := s.Data[key]

	if okData {
		marshaled, err := json.Marshal(data)

		if err == nil {
			tmp := make([][]string, 0)
			json.Unmarshal(marshaled, &tmp)
			return tmp
		}

	}

	return make([][]string, 0)

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

// GetInt retrieves an int value or a default if the key is not found.
func (s *SafeJson) GetInt8(key string, notFound int8) int8 {
	if val, ok := s.Data[key].(float64); ok {
		return int8(val)
	}
	return notFound
}

// GetInt retrieves an int value or a default if the key is not found.
func (s *SafeJson) GetSeconds(key string, notFound Seconds) Seconds {
	if val, ok := s.Data[key].(float64); ok {
		return Seconds(val)
	}
	return notFound
}

// GetInt retrieves an int value or a default if the key is not found.
func (s *SafeJson) GetFloat64(key string, notFound float64) float64 {
	if val, ok := s.Data[key].(float64); ok {
		return float64(val)
	}
	return notFound
}
func (s *SafeJson) GetBool(key string, notFound bool) bool {
	if val, ok := s.Data[key].(bool); ok {
		return bool(val)
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

func (s *SafeJson) GetSafeJson(key string) *SafeJson {
	if val, ok := s.Data[key].(map[string]interface{}); ok {
		tmp := NewSafeJson()
		tmp.Data = val

		return tmp
	}
	return nil

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

func MarshalWithoutDefaults(v interface{}) ([]byte, error) {
	// Reflect the value and ensure it's a struct
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to a struct")
	}

	// Create a map to hold non-zero values
	nonZeroFields := make(map[string]interface{})

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Check if the field has a zero value
		zeroValue := reflect.Zero(field.Type()).Interface()
		if !reflect.DeepEqual(field.Interface(), zeroValue) {
			nonZeroFields[fieldType.Name] = field.Interface()
		}
	}

	// Marshal the filtered map to JSON
	return json.Marshal(nonZeroFields)
}

func FindFirstJsonString(s string) (string, int, error) {
	// Look for opening brace or bracket
	start := -1
	for i, char := range s {
		if char == '{' || char == '[' {
			start = i
			break
		}
	}

	if start == -1 {
		return "", -1, fmt.Errorf("no JSON start found")
	}

	// Count brackets to find matching closing bracket
	openBraces := 0
	openBrackets := 0
	inString := false
	escaped := false

	for i := start; i < len(s); i++ {
		char := rune(s[i])

		if escaped {
			escaped = false
			continue
		}

		if char == '\\' && inString {
			escaped = true
			continue
		}

		if char == '"' {
			inString = !inString
			continue
		}

		if !inString {
			switch char {
			case '{':
				openBraces++
			case '}':
				openBraces--
			case '[':
				openBrackets++
			case ']':
				openBrackets--
			}

			// Found complete JSON
			if openBraces == 0 && openBrackets == 0 {
				candidate := s[start : i+1]
				// Validate it's actually JSON
				var temp interface{}
				if json.Unmarshal([]byte(candidate), &temp) == nil {
					return candidate, start, nil
				}
			}
		}
	}

	return "", -1, fmt.Errorf("no valid JSON found")
}
