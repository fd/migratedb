package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Adapter struct {
	*sqlx.DB
	CREATE_VERSION_TABLE string
	SELECT_VERSIONS      string
	INSERT_VERSION       string
}

var adapters = map[string]func(string) (*Adapter, error){
	"postgres": open_postgres,
}

func Open(u string) (*Adapter, error) {
	parts := strings.SplitN(u, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid connection url: %q", u)
	}

	f := adapters[parts[0]]
	if f == nil {
		return nil, fmt.Errorf("Unsupported database type: %q", parts[0])
	}

	return f(u)
}
