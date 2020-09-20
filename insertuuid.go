package pgkebab

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// InsertUUID inserts a new record into given table and returns the last inserted uuid
// The 3rd param is the optional field name. If not given, the default value "id" will be used
func (l *DBLink) InsertUUID(table string, pairs map[string]interface{}, idFieldName ...string) (string, error) {
	if !l.supposedReady {
		return "", fmt.Errorf("connection not properly initialized")
	}

	if len(pairs) == 0 {
		return "", errors.New(`pgkebab.InsertUUID(undefined values)`)
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

	var lastInsertedUUID string

	if err := l.db.QueryRowContext(ctx, sqlQuery, parameters...).Scan(&lastInsertedUUID); err != nil {
		l.log(`pgkebab.InsertUUID %s db.QueryRowContext has failed: "%v"`, table, err)

		return "", err
	}

	return lastInsertedUUID, nil
}
