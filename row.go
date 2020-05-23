package pgkebab

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Row struct {
	tuple map[string]interface{}
}

// Ready returns true if the tuple contains at least one filled column
func (r Row) Ready() bool {
	return len(r.tuple) > 0
}

// Columns returns an string array filled with column names
func (r Row) Columns() []string {
	if len(r.tuple) < 1 {
		return []string{}
	}

	var a []string

	for x := range r.tuple {
		a = append(a, x)
	}

	return a
}

// has returns true if the tuple is initialized and contains the given key
func (r Row) has(key string) bool {
	if len(r.tuple) > 0 {
		if _, ok := r.tuple[key]; ok {
			return true
		}
	}

	return false
}

// String returns the specified field as string
func (r Row) String(key string) string {
	if !r.has(key) {
		return ""
	}

	switch x := r.tuple[key].(type) {
	case string:
		return x
	case int64:
		return strconv.Itoa(int(x))
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case time.Time:
		return x.Format(time.RFC3339)
	case bool:
		return fmt.Sprintf("%t", x)
	case []byte:
		return string(x)
	}

	return ""
}

// Int64 returns the specified field as Int64
func (r Row) Int64(key string) int64 {
	if !r.has(key) {
		return 0
	}

	switch x := r.tuple[key].(type) {
	case int64:
		return x
	case string:
		{
			i, _ := strconv.ParseInt(x, 10, 64)

			return i
		}
	case float64:
		return int64(x)
	case []byte:
		{
			if len(x) == 1 {
				return int64(x[0])
			}

			return 0
		}
	}

	return 0
}

// Int returns the specified field as Int
func (r Row) Int(key string) int {
	if !r.has(key) {
		return 0
	}

	switch x := r.tuple[key].(type) {
	case int64:
		return int(x)
	case string:
		{
			i, _ := strconv.ParseInt(x, 10, 64)

			return int(i)
		}
	case float64:
		return int(x)
	case []byte:
		{
			if len(x) == 1 {
				return int(x[0])
			}

			return 0
		}
	}

	return 0
}

// Float64 returns the specified field as float64
// Check for conversion details: https://stackoverflow.com/questions/31946344/why-does-go-treat-a-postgresql-numeric-decimal-columns-as-uint8
func (r Row) Float64(key string) float64 {
	if !r.has(key) {
		return 0
	}

	switch x := r.tuple[key].(type) {
	case float64:
		return x
	case string:
		{
			f, _ := strconv.ParseFloat(x, 64)

			return f
		}
	case []uint8:
		{
			f, _ := strconv.ParseFloat(string(x), 64)

			return f
		}
	case int64:
		return float64(x)
	}

	return 0
}

// Bool returns the specified field as boolean
func (r Row) Bool(key string) bool {
	if !r.has(key) {
		return false
	}

	switch x := r.tuple[key].(type) {
	case bool:
		return x
	case string:
		{
			b, _ := strconv.ParseBool(x)

			return b
		}
	}

	return false
}

// Time returns the specified field as time.Time
func (r Row) Time(key string) time.Time {
	if !r.has(key) {
		return time.Time{}
	}

	switch x := r.tuple[key].(type) {
	case time.Time:
		return x
	case string:
		{
			b, _ := time.Parse(time.RFC3339, x)

			return b
		}
	}

	return time.Time{}
}

// Timef returns the specified field as time.Time
// it tries to interpret the date/time value according the given format
func (r Row) Timef(key, format string) time.Time {
	if !r.has(key) {
		return time.Time{}
	}

	switch x := r.tuple[key].(type) {
	case time.Time:
		return x
	case string:
		{
			b, _ := time.Parse(format, x)

			return b
		}
	}

	return time.Time{}
}

// JSON returns the column content serialized as JSON string
func (r Row) JSON() (string, error) {
	if len(r.tuple) < 1 {
		return "", fmt.Errorf("empty row")
	}

	bs, err := json.Marshal(r.tuple)

	return string(bs), err
}

// JSONMap returns the column content desserialized and a bool indicating if the map has some content
func (r Row) JSONMap(key string) (map[string]interface{}, bool) {
	if !r.has(key) {
		return nil, false
	}

	m := make(map[string]interface{})

	switch x := r.tuple[key].(type) {
	case []uint8:
		if err := json.Unmarshal(x, &m); err != nil {
			return m, false
		}
	case string:
		if err := json.Unmarshal([]byte(x), &m); err != nil {
			return m, false
		}
	}

	return m, len(m) > 0
}

// JSONStruct returns the column content desserialized to given target struct
func (r Row) JSONStruct(key string, target interface{}) error {
	if !r.has(key) {
		return fmt.Errorf("unknow column %s", key)
	}

	switch x := r.tuple[key].(type) {
	case []uint8:
		if err := json.Unmarshal(x, &target); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(x), &target); err != nil {
			return err
		}
	}

	return nil
}
