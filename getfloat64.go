package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// MustGetFloat64 returns the query result's single value as a float64
// If sql.ErrNoRows occurs, it's returned.
// The other routines ( without "Must" preffix ) ignores sql.ErrNoRows
func (l *DBLink) MustGetFloat64(sqlQuery string, args ...interface{}) (float64, error) {
	if !l.supposedReady {
		return 0, fmt.Errorf("connection not properly initialized")
	}

	var f sql.NullFloat64

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&f); err != nil {
		l.log("pgkebab.GetFloat64.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return 0, err
	}

	return f.Float64, nil
}

// GetFloat64 returns the query result's single value as a float64
func (l *DBLink) GetFloat64(sqlQuery string, args ...interface{}) (float64, error) {
	x, err := l.MustGetFloat64(sqlQuery, args...)

	if l.IsEmptyErr(err) {
		err = nil
	}

	return x, err
}
