package pgkebab

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// InsertID inserts a new record into given table and returns the last inserted id
// The 3rd param is the optional field name. If not given, the default value "id" will be used
func (l *DBLink) InsertID(table string, pairs map[string]interface{}, idFieldName ...string) (int64, error) {
	if !l.supposedReady {
		return 0, fmt.Errorf("connection not properly initialized")
	}

	if len(pairs) == 0 {
		return 0, errors.New(`pgkebab.InsertID(undefined values)`)
	}

	var (
		fields       []string
		placeholders []string
		parameters   []interface{}
		i            uint
	)

	idField := "id"

	for _, x := range idFieldName {
		idField = x
		break
	}

	for k, v := range pairs {
		fields = append(fields, k)
		i++
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		parameters = append(parameters, v)
	}

	sqlQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s", table, strings.Join(fields, ","), strings.Join(placeholders, ","), idField)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	var lastInsertedID int64

	if err := l.db.QueryRowContext(ctx, sqlQuery, parameters...).Scan(&lastInsertedID); err != nil {
		l.log(`pgkebab.InsertID %s db.QueryRowContext has failed: "%v"`, table, err)

		return 0, err
	}

	return lastInsertedID, nil
}
