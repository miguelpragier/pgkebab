package pgkebab

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// Update executes "update" sql queries against database
// At least one whereParameters is mandatory.
// Important: WHERE only accepts "AND" where criteria.
// If you need to send more complex operations/statements, use Execute()
func (l *DBLink) Update(table string, updatePairs map[string]interface{}, wherePairs map[string]interface{}) (int64, error) {
	if !l.supposedReady {
		return 0, fmt.Errorf("connection not properly initialized")
	}

	if len(updatePairs) < 1 {
		return 0, errors.New("there are no update-pairs to operate")
	}

	if len(wherePairs) < 1 {
		return 0, errors.New("there are no where-pairs to operate")
	}

	var (
		fields      []string
		whereFields []string
		i           uint
		values      []interface{}
	)

	for l, v := range updatePairs {
		i++
		fields = append(fields, fmt.Sprintf("%s=$%d", l, i))
		values = append(values, v)
	}

	for l, v := range wherePairs {
		i++
		whereFields = append(whereFields, fmt.Sprintf("%s=$%d", l, i))
		values = append(values, v)
	}

	sqlQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(fields, ","),
		strings.Join(whereFields, " AND "),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	rs, err := l.db.ExecContext(ctx, sqlQuery, values...)

	if err != nil {
		if l.debugPrint {
			log.Printf(`pgkebab.Update %s db.ExecContext has failed with "%v"\n`, table, err)
		}

		return 0, err
	}

	return rs.RowsAffected()
}
