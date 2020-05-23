package pgkebab

//// GetStringArray returns first col of query result as []string
//// The first param is the behaviour in case of "null scan" errors
//// The available behaviours are:
//// BehaviourSkipErrors
//// BehaviourConvertNullsToEmpty
//// BehaviourAbortOnScanError
//func (k *DBLink) GetStringArray(behaviour uint8, sqlQuery string, params ...interface{}) ([]string, error) {
//	k.assureConnection()
//
//	k.lastSQLQuery = sqlQuery
//
//	rows, err := k.db.Query(sqlQuery, params...)
//
//	var a []string
//
//	if err != nil {
//		if k.DebugPrint {
//			log.Printf(`pgkebab.GetStringArray db.Query The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
//		}
//
//		return a, err
//	}
//
//	defer func() { _ = rows.Close() }()
//
//	for rows.Next() {
//		var s string
//
//		if errx := rows.Scan(&s); errx == nil {
//			a = append(a, s)
//		} else {
//			switch behaviour {
//			case BehaviourSkipErrors:
//				continue
//			case BehaviourConvertNullsToEmpty:
//				if strings.Contains(errx.Error(), "converting driver.Value type <nil>") {
//					a = append(a, ``)
//				}
//			default:
//				return []string{}, errx
//			}
//		}
//	}
//
//	return a, nil
//}
//
