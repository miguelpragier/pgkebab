package pgkebab

import "fmt"

// Pairs is a helper for creating query pairs
// The first pairs item SHOULD be string, or the function panics
// The function panics for "less than 2 items", or "odd quantity of items".
func Pairs(params ...interface{}) map[string]interface{} {
	if len(params) < 2 {
		panic("tried to build query pairs with less than 2 items")
	}

	if len(params)%2 != 0 {
		panic("tried to build query pairs with odd parameters quantity")
	}

	m := make(map[string]interface{})

	for i := 0; i < len(params); i += 2 {
		key, ok := params[i].(string)

		if !ok {
			panic(fmt.Errorf("every pairs first item should be string. given item %d isn't", i))
		}

		value := params[i+1]
		m[key] = value
	}

	return m
}
