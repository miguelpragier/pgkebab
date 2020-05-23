# PGKebab
GOLang PostgreSQL Helper Operators Over [PQ](https://github.com/lib/pq/)
---
![PGKebab](./etc/img/pgkebab.png "PGKebab")
<br>
<br>
[![Go Report Card](https://goreportcard.com/badge/github.com/miguelpragier/pgkebab)](https://goreportcard.com/report/github.com/miguelpragier/pgkebab)

<br>
##### Simple Sample
```golang
package main

import (
	"fmt"
    "os"
	"github.com/miguelpragier/pgkebab"
)

func main(){
    var (
        db *pgkebab.DBLink
    	cs = pgkebab.ConnStringDirect(os.Getenv("{YOURAPPCONNECTIONSTRING}"))
    	opts = pgkebab.Options(cs, 10, 10, 5, 5, 10, true)
    )

	if _db, err := pgkebab.NewConnected(opts); err != nil {
		t.Fatal(err)
	} else {
		db = _db
	}

    row, err := db.GetOne("SELECT name, status_id FROM customer WHERE id=$1",customerID)

    if err != nil {
        handleError(err)
        return
    }
    
    fmt.Println("The customer %s has status %d\n", row.String("name"), row.Int64("status_id"))
    
    if n, err := db.Count("customer"); err == nil {
        fmt.Println("The table customer got %d rows\n", n)
    }
}
```
##### Dependencies:
[pq - Pure Go Postgres driver for database/sql](https://github.com/lib/pq)
<br>

---
![Requires.io](https://img.shields.io/requires/github.com/lib/pq)