package database

import (
	"testing"
	"testing/fstest"
)

func TestLoadMigrationsSortsByVersion(t *testing.T) {
	files := fstest.MapFS{
		"migrations/000002_second.sql": {Data: []byte("select 2;")},
		"migrations/000001_first.sql":  {Data: []byte("select 1;")},
	}

	migrations, err := LoadMigrations(files, "migrations")
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}

	if len(migrations) != 2 {
		t.Fatalf("expected 2 migrations, got %d", len(migrations))
	}
	if migrations[0].Version != 1 || migrations[1].Version != 2 {
		t.Fatalf("migrations not sorted by version: %#v", migrations)
	}
}

func TestLoadMigrationsRejectsInvalidName(t *testing.T) {
	files := fstest.MapFS{
		"migrations/bad.sql": {Data: []byte("select 1;")},
	}

	if _, err := LoadMigrations(files, "migrations"); err == nil {
		t.Fatal("expected invalid migration filename to fail")
	}
}
