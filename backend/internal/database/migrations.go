package database

import (
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
)

type Migration struct {
	Version int
	Name    string
	Path    string
	SQL     string
}

func LoadMigrations(files fs.FS, dir string) ([]Migration, error) {
	entries, err := fs.ReadDir(files, dir)
	if err != nil {
		return nil, err
	}

	migrations := make([]Migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") {
			return nil, fmt.Errorf("invalid migration file %q: expected .sql", name)
		}

		parts := strings.SplitN(strings.TrimSuffix(name, ".sql"), "_", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid migration file %q: expected version_name.sql", name)
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid migration file %q: version must be numeric", name)
		}

		path := dir + "/" + name
		content, err := fs.ReadFile(files, path)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    parts[1],
			Path:    path,
			SQL:     string(content),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}
