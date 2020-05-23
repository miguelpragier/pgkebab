package pgkebab

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// MustGetJSONStruct tries to scan query result's single column value into given struct
// If sql.ErrNoRows occurs, it's returned.
// The other routines ( without "Must" preffix ) ignores sql.ErrNoRows
func (l *DBLink) MustGetJSONStruct(target interface{}, sqlQuery string, args ...interface{}) error {
	if !l.supposedReady {
		return fmt.Errorf("connection not properly initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	var j json.RawMessage

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&j); err != nil {
		l.log("pgkebab.GetJSONAsMap.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return err
	}

	if data, err := j.MarshalJSON(); err != nil {
		return err
	} else {
		if err0 := json.Unmarshal(data, &target); err0 != nil {
			return err0
		}
	}

	return nil
}

// GetJSONStruct tries to scan query result's single column value into given struct
func (l *DBLink) GetJSONStruct(target interface{}, sqlQuery string, args ...interface{}) error {
	err := l.MustGetJSONStruct(target, sqlQuery, args...)

	if l.IsEmptyErr(err) {
		err = nil
	}

	return err
}
