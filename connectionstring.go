package pgkebab

import (
	"fmt"
	"os"
	"strings"
)

// ConnectionString has directives to obtain database connection string
// To create a new ConnectionString, use ConnStringDirect(), ConnStringEnvVar()
type ConnectionString struct {
	method csMethod
	key    string
	value  string
}

func (s ConnectionString) validate() error {
	if s.method == csMethodDirectParam && s.value == "" {
		return fmt.Errorf(`connection string declared as "direct param", but value field empty`)
	}

	if s.method == csMethodEnvVariable {
		if s.key == "" {
			return fmt.Errorf(`connection string declared as "environment variable", but key field is empty`)
		}

		val := os.Getenv(s.key)

		if val == "" {
			return fmt.Errorf(`connection string declared as "environment variable", but the given key doesn't exist or is empty`)
		}

		s.value = val
	}

	return nil
}

func (s ConnectionString) get() string {
	return s.value
}

// ConnStringDirect sets connection string with given value
func ConnStringDirect(val string) *ConnectionString {
	return &ConnectionString{
		method: csMethodDirectParam,
		value:  strings.TrimSpace(val),
	}
}

// ConnStringEnvVar configures the engine to obtain connection string from given environment variable
func ConnStringEnvVar(key string) *ConnectionString {
	return &ConnectionString{
		method: csMethodEnvVariable,
		key:    strings.TrimSpace(key),
		value:  os.Getenv(key),
	}
}
