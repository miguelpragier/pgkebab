package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// GetCount returns the number of records produced by given query
func (l *DBLink) GetCount(sqlQuery string, args ...interface{}) (int64, error) {
	if !l.supposedReady {
		return 0, fmt.Errorf("connection not properly initialized")
	}

	var i sql.NullInt64

	sqlQuery = fmt.Sprintf("SELECT COUNT(*) FROM ( %s )", sqlQuery)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&i); err != nil {
		if l.IsEmptyErr(err) {
			return 0, nil
		}

		l.log("pgkebab.GetCount.QueryRowContext().Scan(%s) %v", sqlQuery, err)

		return 0, err
	}

	return i.Int64, nil
}
