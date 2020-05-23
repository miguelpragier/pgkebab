package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// MustGetTime returns the query result's single value as a time.Time
// If sql.ErrNoRows occurs, it's returned.
// The other routines ( without "Must" preffix ) ignores sql.ErrNoRows
func (l *DBLink) MustGetTime(sqlQuery string, args ...interface{}) (time.Time, error) {
	if !l.supposedReady {
		return time.Time{}, fmt.Errorf("connection not properly initialized")
	}

	var t sql.NullTime

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&t); err != nil {
		l.log("pgkebab.GetTime.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return time.Time{}, err
	}

	return t.Time, nil
}

// GetTime returns the query result's single value as a time.Time
func (l *DBLink) GetTime(sqlQuery string, args ...interface{}) (time.Time, error) {
	t, err := l.MustGetTime(sqlQuery, args...)

	if l.IsEmptyErr(err) {
		err = nil
	}

	return t, err
}
