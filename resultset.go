package pgkebab

import (
	"encoding/json"
	"fmt"
)

const resultsetFirstRow uint = 0

// Resultset contains rows and result metadata
type Resultset struct {
	pointer     uint
	columns     []string
	columnTypes []string
	records     []Row
	err         error
	debugPrint  bool
}

// Count returns the number of rows in this resultset
func (r Resultset) Count() int {
	return len(r.records)
}

// Bottom returns true if the cursor reached its final row
func (r Resultset) Bottom() bool {
	var (
		l          = len(r.records)
		hasRecords = l > 0
		lastItem   = r.pointer > uint(l)
	)

	return hasRecords && lastItem
}

// Top returns true if the cursor is set on first item
func (r Resultset) Top() bool {
	return len(r.records) > 0 && r.pointer == 0
}

// Next returns true if there's a next item and position the internal cursonr on it
func (r *Resultset) Next() bool {
	if r == nil || len(r.records) == 0 || r.pointer >= uint(len(r.records)) {
		return false
	}

	r.pointer++

	return true
}

// Row returns the current row, set by the internal pointer/position.
func (r *Resultset) Row() (Row, error) {
	if r == nil || len(r.records) == 0 {
		return Row{}, fmt.Errorf("there's no current row")
	}

	if r.Bottom() {
		return Row{}, fmt.Errorf("cursor has reached bottom")
	}

	return r.records[r.pointer-1], nil
}

// Rewind set the internal cursor at first row
func (r *Resultset) Rewind() {
	r.pointer = resultsetFirstRow
}

// Rows returns a row array from the current resultset
func (r *Resultset) Rows() ([]Row, error) {
	if r == nil || len(r.records) == 0 {
		return nil, fmt.Errorf("resultset is empty")
	}

	return r.records, nil
}

// Map returns the resultset as a []map[string]interface{}
// If the resultset is empty, returns nil
func (r Resultset) Map() []map[string]interface{} {
	if len(r.records) < 1 {
		return nil
	}

	m := make([]map[string]interface{}, len(r.records))

	for i, row := range r.records {
		m[i] = make(map[string]interface{})

		for k, v := range row.tuple {
			m[i][k] = v
		}
	}

	return m
}

// JSON returns a json resultset representation
func (r Resultset) JSON(errorIfEmpty bool) (string, error) {
	if len(r.records) < 1 {
		if errorIfEmpty {
			return "[]", fmt.Errorf("empty resultset")
		} else {
			return "[]", nil
		}
	}

	m := r.Map()

	j, err := json.Marshal(m)

	return string(j), err
}
