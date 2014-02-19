package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func open_postgres(u string) (*Adapter, error) {
	ds, err := pq.ParseURL(u)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("postgres", ds)
	if err != nil {
		return nil, fmt.Errorf("error while connectiong to %q: %s", u, err)
	}

	return &Adapter{
		db,
		`CREATE TABLE mdb_versions (version INT, name VARCHAR(256));`,
		`SELECT version FROM mdb_versions;`,
		`INSERT INTO mdb_versions (version, name) VALUES ($1, $2);`,
	}, nil
}
