package entities

import (
	"fmt"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SimBehData tools.SafeJson

// GetString returns the value at key as a string. Returns empty string if key doesn't exist or value isn't a string.
func (d SimBehData) GetString(key string) string {
	if v, ok := d[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetStringOk returns the value at key as a string and a bool indicating success.
func (d SimBehData) GetStringOk(key string) (string, bool) {
	if v, ok := d[key]; ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}

// GetInt returns the value at key as an int. Returns 0 if key doesn't exist or value isn't an int.
func (d SimBehData) GetInt(key string) int {
	if v, ok := d[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int32:
			return int(n)
		case int64:
			return int(n)
		case float32:
			return int(n)
		case float64:
			return int(n)
		}
	}
	return 0
}

func (d SimBehData) SetSeconds(key string, data tools.Seconds) tools.Seconds {
	d[key] = data
	return data
}
func (d SimBehData) GetSeconds(key string) tools.Seconds {
	a := tools.SafeJson(d)
	return a.GetSeconds(key, 0)
}

// GetInt64 returns the value at key as an int64.
func (d SimBehData) GetInt64(key string) int64 {
	if v, ok := d[key]; ok {
		switch n := v.(type) {
		case int:
			return int64(n)
		case int32:
			return int64(n)
		case int64:
			return n
		case float32:
			return int64(n)
		case float64:
			return int64(n)
		}
	}
	return 0
}

// GetFloat returns the value at key as a float64.
func (d SimBehData) GetFloat(key string) float64 {
	if v, ok := d[key]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case float32:
			return float64(n)
		case int:
			return float64(n)
		case int32:
			return float64(n)
		case int64:
			return float64(n)
		}
	}
	return 0
}

// GetBool returns the value at key as a bool.
func (d SimBehData) GetBool(key string) bool {
	if v, ok := d[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// GetBool returns the value at key as a bool.
func (d SimBehData) GetPrimitiveObjectId(key string) primitive.ObjectID {
	if v, ok := d[key]; ok {
		if b, ok := v.(primitive.ObjectID); ok {
			return b
		}
	}
	return primitive.NilObjectID
}

// GetActor returns the value at key as a *SimActor. Returns nil if key doesn't exist or value isn't a *SimActor.
func (d SimBehData) GetActor(key string) *SimActor {
	if v, ok := d[key]; ok {
		if a, ok := v.(*SimActor); ok {
			return a
		}
	}
	return nil
}

// GetActorOk returns the value at key as a *SimActor and a bool indicating success.
func (d SimBehData) GetActorOk(key string) (*SimActor, bool) {
	if v, ok := d[key]; ok {
		if a, ok := v.(*SimActor); ok {
			return a, true
		}
	}
	return nil, false
}

// GetEntity returns the value at key as a *SimEntity.
func (d SimBehData) GetEntity(key string) *SimEntity {
	if v, ok := d[key]; ok {
		if e, ok := v.(*SimEntity); ok {
			return e
		}
	}
	return nil
}

// GetEntityOk returns the value at key as a *SimEntity and a bool indicating success.
func (d SimBehData) GetEntityOk(key string) (*SimEntity, bool) {
	if v, ok := d[key]; ok {
		if e, ok := v.(*SimEntity); ok {
			return e, true
		}
	}
	return nil, false
}

// GetActors returns the value at key as a []*SimActor.
func (d SimBehData) GetActors(key string) []*SimActor {
	if v, ok := d[key]; ok {
		if a, ok := v.([]*SimActor); ok {
			return a
		}
	}
	return nil
}

// GetActorsOk returns the value at key as a []*SimActor and a bool indicating success.
func (d SimBehData) GetActorsOk(key string) ([]*SimActor, bool) {
	if v, ok := d[key]; ok {
		if a, ok := v.([]*SimActor); ok {
			return a, true
		}
	}
	return nil, false
}

// GetEntities returns the value at key as a []*SimEntity.
func (d SimBehData) GetEntities(key string) []*SimEntity {
	if v, ok := d[key]; ok {
		if e, ok := v.([]*SimEntity); ok {
			return e
		}
	}
	return nil
}

// GetEntitiesOk returns the value at key as a []*SimEntity and a bool indicating success.
func (d SimBehData) GetEntitiesOk(key string) ([]*SimEntity, bool) {
	if v, ok := d[key]; ok {
		if e, ok := v.([]*SimEntity); ok {
			return e, true
		}
	}
	return nil, false
}

// GetStrings returns the value at key as a []string.
func (d SimBehData) GetStrings(key string) []string {
	if v, ok := d[key]; ok {
		if s, ok := v.([]string); ok {
			return s
		}
	}
	return nil
}

// GetMap returns the value at key as a SimBehData (nested dictionary).
func (d SimBehData) GetMap(key string) SimBehData {
	if v, ok := d[key]; ok {
		if m, ok := v.(SimBehData); ok {
			return m
		}
		if m, ok := v.(map[string]interface{}); ok {
			return SimBehData(m)
		}
	}
	return nil
}

// Get returns the raw value and a bool indicating whether the key exists.
func (d SimBehData) Get(key string) (interface{}, bool) {
	v, ok := d[key]
	return v, ok
}

// Has reports whether the key exists in the dictionary.
func (d SimBehData) Has(key string) bool {
	_, ok := d[key]
	return ok
}

// Set stores a value under the given key. Returns d for chaining.
func (d SimBehData) SetActor(key string, actor *SimActor) *SimActor {
	d[key] = actor
	return actor
}

// Set stores a value under the given key. Returns d for chaining.
func (d SimBehData) Set(key string, value interface{}) SimBehData {
	d[key] = value
	return d
}

// String provides a debug representation.
func (d SimBehData) String() string {
	return fmt.Sprintf("SimBehData%v", map[string]interface{}(d))
}
