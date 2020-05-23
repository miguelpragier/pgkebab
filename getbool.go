package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// MustGetBool returns the query result's single value as a boolean
// If sql.ErrNoRows occurs, it's returned.
// The other routines ( without "Must" preffix ) ignores sql.ErrNoRows
func (l *DBLink) MustGetBool(sqlQuery string, args ...interface{}) (bool, error) {
	if !l.supposedReady {
		return false, fmt.Errorf("connection not properly initialized")
	}

	var b sql.NullBool

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&b); err != nil {
		l.log("pgkebab.GetBool.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return false, err
	}

	return b.Bool, nil
}

// GetBool returns the query result's single value as a boolean
func (l *DBLink) GetBool(sqlQuery string, args ...interface{}) (bool, error) {
	b, err := l.MustGetBool(sqlQuery, args...)

	if l.IsEmptyErr(err) {
		err = nil
	}

	return b, err
}

// Bool returns the query result's single value as a boolean.
// in case of any error, it returns false
func (l *DBLink) Bool(sqlQuery string, args ...interface{}) bool {
	b, _ := l.MustGetBool(sqlQuery, args...)

	return b
}
