![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg) ![GitHub](https://img.shields.io/badge/goDoc-Yes!-blue.svg) 
[![Go Report Card](https://goreportcard.com/badge/github.com/miguelpragier/pgkebab?update)](https://goreportcard.com/report/github.com/miguelpragier/pgkebab) 
![Go Version](https://img.shields.io/badge/GO%20version-%3E%3D1.13-blue)

# PGKebab
### GOLang PostgreSQL Helper Over [PQ](https://github.com/lib/pq/)
##### Makes PostgreSQL handling as easy and simple as GOlang
    Replace heavy ORMs and dense routines with simple SQL queries
---
![PGKebab](./etc/img/pgkebab.png "PGKebab")
<br>
<br>
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/miguelpragier/pgkebab )](https://goreportcard.com/report/github.com/miguelpragier/pgkebab) -->

##### Simple Sample
```golang
package main

import (
	"fmt"
	"github.com/miguelpragier/pgkebab"
	"log"
)

func main() {
	const (
		connectionTimeout                  = 10
		executiontionTimeout               = 10
		connectionMaxAttempts              = 5
		connectionMaxMinutesRetrying       = 5
		secondsBetweenReconnectionAttempts = 10
		debugLogPrint                      = true
	)

	var (
		cs         = pgkebab.ConnStringEnvVar("{YOURAPPCONNECTIONSTRING}")
		opts       = pgkebab.Options(cs, connectionTimeout, executiontionTimeout, connectionMaxAttempts, connectionMaxMinutesRetrying, secondsBetweenReconnectionAttempts, debugLogPrint)
		customerID = 1
	)

	db, errcnx := pgkebab.NewConnected(opts)

	if errcnx != nil {
		log.Fatal(errcnx)
	}

	if row, err := db.GetOne("SELECT name, status_id FROM customers WHERE id=$1", customerID); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("the customer", row.String("name"), "has status", row.Int64("status_id"))
	}

	if n, err := db.GetCount("customers"); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("table customer counts", n, "rows")
	}
}
```
##### Dependencies:
[pq - Pure Go Postgres driver for database/sql](https://github.com/lib/pq)
<br>

---
<!-- ![Requires.io](https://img.shields.io/requires/github.com/lib/pq) -->
