package pgkebab

import (
	"database/sql"
)

// MustGetMany returns sql query result as an array of map[string]string
// The difference to GetMany is that, if no record is found, it returns an error
func (l *DBLink) MustGetMany(sqlQuery string, args ...interface{}) (Resultset, error) {
	rs, err := l.GetMany(sqlQuery, args...)

	if err != nil {
		return Resultset{err: err}, err
	}

	if rs.Count() == 0 {
		return Resultset{err: sql.ErrNoRows}, sql.ErrNoRows
	}

	return rs, nil
}
