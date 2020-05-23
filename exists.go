package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Exists returns true if the query returns non null resultset
func (l *DBLink) Exists(sqlQuery string, args ...interface{}) (bool, error) {
	if !l.supposedReady {
		return false, fmt.Errorf("connection not properly initialized")
	}

	var b sql.NullBool

	sqlQuery = fmt.Sprintf("SELECT EXISTS ( %s )", sqlQuery)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	if err := l.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&b); err != nil {
		if l.IsEmptyErr(err) {
			return false, nil
		}

		if l.debugPrint {
			log.Printf("pgkebab.Exists.QueryRowContext().Scan(%s) %v\n", sqlQuery, err)
		}

		return false, err
	}

	return b.Bool, nil
}
