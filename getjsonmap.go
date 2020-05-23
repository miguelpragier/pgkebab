package pgkebab

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// MustGetJSONMap returns the query result's single column value as map
// If sql.ErrNoRows occurs, it's returned.
// The other routines ( without "Must" preffix ) ignores sql.ErrNoRows
func (l *DBLink) MustGetJSONMap(sqlQuery string, args ...interface{}) (map[string]interface{}, error) {
	if !l.supposedReady {
		return nil, fmt.Errorf("connection not properly initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	j := json.RawMessage{}

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&j); err != nil {
		l.log("pgkebab.GetJSONAsMap.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return nil, err
	}

	m := make(map[string]interface{})

	if data, err := j.MarshalJSON(); err != nil {
		return nil, err
	} else {
		if err0 := json.Unmarshal(data, &m); err0 != nil {
			return nil, err0
		}
	}

	return m, nil
}

// GetJSONMap returns the query result's single column value as map
func (l *DBLink) GetJSONMap(sqlQuery string, args ...interface{}) (map[string]interface{}, error) {
	m, err := l.MustGetJSONMap(sqlQuery, args...)

	if l.IsEmptyErr(err) {
		err = nil
	}

	return m, err
}
