package pgkebab

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

func rowsClose(rows *sql.Rows) {
	_ = rows.Close()
}

func decimalAsFloat64(databaseType string, value interface{}) (float64, bool) {
	if databaseType == "NUMERIC" || databaseType == "DECIMAL" {
		if x, ok := value.([]uint8); ok {
			f, _ := strconv.ParseFloat(string(x), 64)

			return f, true
		}
	}

	return 0, false
}

// GetMany returns sql query result as an array of map[string]string
func (l *DBLink) GetMany(sqlQuery string, args ...interface{}) (Resultset, error) {
	if !l.supposedReady {
		return Resultset{}, fmt.Errorf("connection not properly initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.executionTimeoutSeconds)*time.Second)

	defer cancel()

	rows, err := l.db.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		if l.debugPrint {
			log.Printf(`pgkebab.GetMany.db.Query("%s") has failed: "%v"\n`, sqlQuery, err)
		}

		return Resultset{err: err}, err
	}

	defer rowsClose(rows)

	cols, err0 := rows.Columns()

	if err0 != nil {
		if l.debugPrint {
			log.Printf(`pgkebab.GetMany.rows.Columns("%s") has failed: "%v"\n`, sqlQuery, err0)
		}

		return Resultset{err: err0}, err0
	}

	rs := Resultset{
		pointer:    resultsetFirstRow,
		columns:    cols,
		debugPrint: l.debugPrint,
	}

	if colTypes, err8 := rows.ColumnTypes(); err8 != nil {
		if l.debugPrint {
			log.Printf(`pgkebab.GetMany.rows.ColumnTypes("%s") has failed: "%v"\n`, sqlQuery, err8)
		}
	} else {
		for _, x := range colTypes {
			rs.columnTypes = append(rs.columnTypes, x.DatabaseTypeName())
		}
	}

	howManyCols := len(cols)

	// Byte array for store values
	rs.records = make([]Row, 0)

	for rows.Next() {
		// Once Rows.Scan() expects []interface{} ( slice of interfaces ), we need one of these as a temporary buffer
		fakeDest := make([]interface{}, howManyCols)

		// Here we set the real []interface's address for each interface instance
		realDest := make([]interface{}, howManyCols)

		for i := range realDest {
			fakeDest[i] = &realDest[i]
		}

		if er2 := rows.Scan(fakeDest...); er2 != nil {
			if l.debugPrint {
				log.Printf(`pgkebab.GetMany.rows.Scan("%s") has failed: "%v"\n`, sqlQuery, er2)
			}

			return Resultset{}, er2
		}

		rw := Row{tuple: make(map[string]interface{})}

		for i, v := range realDest {
			col := cols[i]

			// Try to intercept and handle float64
			// To know further: https://stackoverflow.com/questions/31946344/why-does-go-treat-a-postgresql-numeric-decimal-columns-as-uint8
			if len(rs.columnTypes) >= i {
				if f, ok := decimalAsFloat64(rs.columnTypes[i], v); ok {
					rw.tuple[col] = f
				} else {
					rw.tuple[col] = v
				}
			}
		}

		rs.records = append(rs.records, rw)
	}

	return rs, nil
}
