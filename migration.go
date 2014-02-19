package main

import (
	"fmt"
	"io/ioutil"
	path "path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Migration struct {
	Version int
	Name    string
	SQL     string
}

func LoadMigrations(dir string) ([]*Migration, error) {
	var (
		err        error
		dir2       string
		matches    []string
		migrations []*Migration
	)

	dir2, err = path.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("error while reading migrations in %s: %q", dir, err)
	}

	matches, err = path.Glob(path.Join(dir2, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf("error while reading migrations in %s: %q", dir2, err)
	}

	for _, file_path := range matches {
		migration, err := LoadMigration(file_path)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	sort.Sort(migrationSorter(migrations))
	return migrations, nil
}

func LoadMigration(file_path string) (*Migration, error) {
	var (
		file_name = path.Base(file_path)
		version   int
		name      string
		data      []byte
		sql       string
		err       error
	)

	version, name, err = split_file_name(file_name)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadFile(file_path)
	if err != nil {
		return nil, fmt.Errorf("migration %d: %s", version, err)
	}

	sql = string(data)

	return &Migration{version, name, sql}, nil
}

func split_file_name(s string) (int, string, error) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("Invalid migration name: %q", s)
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("Invalid migration name: %q", s)
	}

	name := parts[1]
	name = strings.TrimSuffix(name, ".sql")
	name = strings.Replace(name, "--", "-", -1)
	name = strings.Replace(name, "-", " ", -1)

	return version, name, nil
}

type migrationSorter []*Migration

func (l migrationSorter) Len() int           { return len(l) }
func (l migrationSorter) Less(i, j int) bool { return l[i].Version < l[j].Version }
func (l migrationSorter) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
