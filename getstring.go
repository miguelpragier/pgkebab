package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// MustGetString returns the query result's single value as a string
// If sql.ErrNoRows occurs, it's returned.
// The other routines ( without "Must" preffix ) ignores sql.ErrNoRows
func (l *DBLink) MustGetString(sqlQuery string, args ...interface{}) (string, error) {
	if !l.supposedReady {
		return "", fmt.Errorf("connection not properly initialized")
	}

	var s sql.NullString

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&s); err != nil {
		l.log("pgkebab.GetString.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return "", err
	}

	return s.String, nil
}

// GetString returns the query result's single value as a string
func (l *DBLink) GetString(sqlQuery string, args ...interface{}) (string, error) {
	s, err := l.MustGetString(sqlQuery, args...)

	if l.IsEmptyErr(err) {
		err = nil
	}

	return s, err
}
