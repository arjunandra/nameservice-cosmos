package types

import "strings"

// Querier Types 

type QueryResResolve struct {
	Value string `json:"value"`
}

type QueryResNames []string

// Implement fmt.Stringer

func (r QueryResResolve) String() string {
	return r.Value
}

func (n QueryResNames) String() string {
	return strings.Join(n[:], "\n")
}
