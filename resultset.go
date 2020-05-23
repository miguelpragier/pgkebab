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

// JSON returns a json resultset representation
func (r Resultset) JSON() (string, error) {
	if len(r.records) < 1 {
		return "[]", nil
	}

	str, err := json.Marshal(r.records)

	return string(str), err
}
