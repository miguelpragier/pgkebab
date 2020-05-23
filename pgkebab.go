package pgkebab

import (
	// PostgreSql Driver
	_ "github.com/lib/pq"
)

type csMethod uint

const (
	connectionTimeoutSecondsDefault             uint = 10
	executionTimeoutSecondsDefault              uint = 10
	timeBetweenConnectionAttemptsSecondsDefault uint = 30
	//maxOpenConnectionsDefault                   uint = 10

	// csMethodDirectParam Connection String sent directly as a string param
	csMethodDirectParam csMethod = 1
	// csMethodEnvVariable Connection String can be found in an environment variable
	csMethodEnvVariable csMethod = 2
)
