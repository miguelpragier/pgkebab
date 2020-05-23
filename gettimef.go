package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// MustGetTimef returns the query result's single value as a formatted time.Time, considering the given format string
// If sql.ErrNoRows occurs, it's returned.
// The other routines ( without "Must" preffix ) ignores sql.ErrNoRows
func (l *DBLink) MustGetTimef(format, sqlQuery string, args ...interface{}) (string, error) {
	if !l.supposedReady {
		return "", fmt.Errorf("connection not properly initialized")
	}

	var t sql.NullTime

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&t); err != nil {
		l.log("pgkebab.GetTime.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return "", err
	}

	return t.Time.Format(format), nil
}

// GetTimef returns the query result's single value as a formatted time.Time, considering the given format string
func (l *DBLink) GetTimef(format, sqlQuery string, args ...interface{}) (string, error) {
	s, err := l.MustGetTimef(format, sqlQuery, args...)

	if l.IsEmptyErr(err) {
		err = nil
	}

	return s, err
}
