package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// MustGetInt64 returns the query result's single value as an int64
// If sql.ErrNoRows occurs, it's returned.
// The other routines ( without "Must" preffix ) ignores sql.ErrNoRows
func (l *DBLink) MustGetInt64(sqlQuery string, args ...interface{}) (int64, error) {
	if !l.supposedReady {
		return 0, fmt.Errorf("connection not properly initialized")
	}

	var i sql.NullInt64

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&i); err != nil {
		l.log("pgkebab.GetInt64.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return 0, err
	}

	return i.Int64, nil
}

// GetInt64 returns the query result's single value as an int64
func (l *DBLink) GetInt64(sqlQuery string, args ...interface{}) (int64, error) {
	i, err := l.MustGetInt64(sqlQuery, args...)

	if l.IsEmptyErr(err) {
		err = nil
	}

	return i, err
}
