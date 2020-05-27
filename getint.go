package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// MustGetInt returns the query result's single value as an int
// If sql.ErrNoRows occurs, it's returned.
// The other routines ( without "Must" preffix ) ignores sql.ErrNoRows
func (l *DBLink) MustGetInt(sqlQuery string, args ...interface{}) (int, error) {
	if !l.supposedReady {
		return 0, fmt.Errorf("connection not properly initialized")
	}

	var i sql.NullInt64

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&i); err != nil {
		l.log("pgkebab.GetInt.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return 0, err
	}

	return strconv.Atoi(fmt.Sprintf("%d", i.Int64))
}

// GetInt returns the query result's single value as an int64
func (l *DBLink) GetInt(sqlQuery string, args ...interface{}) (int, error) {
	i, err := l.MustGetInt(sqlQuery, args...)

	if l.IsEmptyErr(err) {
		err = nil
	}

	return i, err
}
