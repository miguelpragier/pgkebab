package pgkebab

import (
	"database/sql"
)

// GetOne returns sql query result as a Row
func (l *DBLink) GetOne(sqlQuery string, args ...interface{}) (Row, error) {
	rs, err := l.MustGetMany(sqlQuery, args...)

	if err != nil {
		return Row{}, err
	}

	if rs.Next() {
		return rs.Row()
	}

	return Row{}, sql.ErrNoRows
}
