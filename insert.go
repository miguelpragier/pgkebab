package pgkebab

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Insert executes an "insert" sql query against given table
func (l *DBLink) Insert(table string, pairs map[string]interface{}) error {
	if !l.supposedReady {
		return fmt.Errorf("connection not properly initialized")
	}

	if len(pairs) == 0 {
		return errors.New(`pgkebab.Insert(undefined values)`)
	}

	var (
		fields       []string
		placeholders []string
		parameters   []interface{}
		i            uint
	)

	for k, v := range pairs {
		fields = append(fields, k)
		i++
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		parameters = append(parameters, v)
	}

	sqlQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(fields, ","), strings.Join(placeholders, ","))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	_, err := l.db.ExecContext(ctx, sqlQuery, parameters...)

	if err != nil {
		l.log(`pgkebab.Insert %s db.Exec has failed: "%v"`, table, err)

		return err
	}

	return nil
}
