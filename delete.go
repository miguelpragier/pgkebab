package pgkebab

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// Delete executes "delete" statements against specific table, considering the mandatory "where" criteria
// The routine tries to return affected rows count.
func (l *DBLink) Delete(table string, wherePairs map[string]interface{}) (int64, error) {
	if !l.supposedReady {
		return 0, fmt.Errorf("connection not properly initialized")
	}

	if len(wherePairs) == 0 {
		return 0, errors.New(`pgkebab.Delete("where" criteria not given`)
	}

	var (
		whereList  []string
		parameters []interface{}
		i          uint
	)

	for k, v := range wherePairs {
		i++
		s := fmt.Sprintf("%s=$%d", k, i)
		whereList = append(whereList, s)
		parameters = append(parameters, v)
	}

	whereExpression := strings.Join(whereList, " AND ")

	sqlQuery := fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereExpression)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	rs, err := l.db.ExecContext(ctx, sqlQuery, parameters...)

	if err != nil {
		if l.IsEmptyErr(err) {
			return 0, nil
		}

		if l.debugPrint {
			log.Printf(`pgkebab.Delete "%s" db.Exec has failed: "%v"\n`, table, err)
		}

		return 0, err
	}

	return rs.RowsAffected()
}
