package pgkebab

////noinspection ALL
//const (
//	BehaviourSkipErrors          = 0
//	BehaviourConvertNullsToZero  = 1
//	BehaviourConvertNullsToEmpty = 1
//	BehaviourAbortOnScanError    = 2
//)
//
//// GetIntArray returns first col of query result as []int
//// The first param is the behaviour in case of "null scan" errors
//// The available behaviours are:
//// BehaviourSkipErrors
//// BehaviourConvertNullsToZero
//// BehaviourAbortOnScanError
//func (k *DBLink) GetIntArray(behaviour uint8, sqlQuery string, params ...interface{}) ([]int, error) {
//	k.assureConnection()
//
//	k.lastSQLQuery = sqlQuery
//
//	rows, err := k.db.Query(sqlQuery, params...)
//
//	var a []int
//
//	if err != nil {
//		if k.DebugPrint {
//			log.Printf(`pgkebab.GetIntArray db.Query The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
//		}
//
//		return a, err
//	}
//
//	defer func() { _ = rows.Close() }()
//
//	for rows.Next() {
//		var i int
//
//		if errx := rows.Scan(&i); errx == nil {
//			a = append(a, i)
//		} else {
//			switch behaviour {
//			case BehaviourSkipErrors:
//				continue
//			case BehaviourConvertNullsToZero:
//				if strings.Contains(errx.Error(), "converting driver.Value type <nil>") {
//					a = append(a, 0)
//				}
//			default:
//				return []int{}, errx
//			}
//		}
//	}
//
//	return a, nil
//}
//
