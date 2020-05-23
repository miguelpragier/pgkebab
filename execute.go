package pgkebab

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Exec sends the given sql query to database server. typically update, delete and insert, and return the number of affected rows
func (l *DBLink) Exec(sqlQuery string, args ...interface{}) (int64, error) {
	if !l.supposedReady {
		return 0, fmt.Errorf("connection not properly initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	rs, err := l.db.ExecContext(ctx, sqlQuery, args...)

	if err != nil {
		if l.debugPrint {
			log.Printf("pgkebab.Exec.ExecContext(%s) %v\n", sqlQuery, err)
		}

		return 0, err
	}

	return rs.RowsAffected()
}
