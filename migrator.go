package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type Migrator struct {
	Migrations []*Migration
	DB         *Adapter
	Log        *log.Logger
}

func (m *Migrator) Run() {
	var (
		tx          *sqlx.Tx
		err         error
		success     bool
		versions    []struct{ Version int }
		version_map = map[int]bool{}
	)

	err = m.ensure_versions_table_exists()
	if err != nil {
		m.Log.Printf("error while initializing migrations table: %s", err)
		return
	}

	tx, err = m.DB.Beginx()
	if err != nil {
		m.Log.Printf("error while opening the transaction: %s", err)
		return
	}

	defer func() {
		if success {
			err = tx.Commit()
			if err != nil {
				m.Log.Printf("error while commiting the transaction: %s", err)
			}
		} else {
			err = tx.Rollback()
			if err != nil {
				m.Log.Printf("error while aborting the transaction: %s", err)
			}
		}
	}()

	err = tx.Select(&versions, m.DB.SELECT_VERSIONS)
	if err != nil {
		m.Log.Printf("error while loading migrations table: %s", err)
		return
	}

	for _, v := range versions {
		version_map[v.Version] = true
	}

	for _, migration := range m.Migrations {
		if version_map[migration.Version] {
			m.Log.Printf("%014d: %s [DONE]", migration.Version, migration.Name)
			continue
		}

		_, err := tx.Exec(migration.SQL)
		if err != nil {
			m.Log.Printf("%014d: %s [ERROR]:\n  %s", migration.Version, migration.Name, err)
			return
		}

		_, err = tx.Exec(m.DB.INSERT_VERSION, migration.Version, migration.Name)
		if err != nil {
			m.Log.Printf("%014d: %s [ERROR]:\n  %s", migration.Version, migration.Name, err)
			return
		}

		m.Log.Printf("%014d: %s [DONE]", migration.Version, migration.Name)
	}

	success = true
}

func (m *Migrator) ensure_versions_table_exists() error {
	var (
		exists   bool
		versions []struct{ Version int }
	)

	m.tx(func(tx *sqlx.Tx) error {
		err := tx.Select(&versions, m.DB.SELECT_VERSIONS)
		if err == nil {
			exists = true
		}
		return fmt.Errorf("always err")
	})

	if exists {
		return nil
	}

	return m.tx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(m.DB.CREATE_VERSION_TABLE)
		return err
	})
}

func (m *Migrator) tx(f func(*sqlx.Tx) error) (err error) {
	var (
		tx      *sqlx.Tx
		success bool
	)

	tx, err = m.DB.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if success {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	err = f(tx)
	if err != nil {
		return err
	}

	success = true
	return nil
}
