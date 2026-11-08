package dbqueries

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

type Queries map[string]string

var requiredQueries = []string{
	"CreateConversations",
	"CreateMessages",
	"FetchMessages",
	"SaveConversation",
	"SaveMessage",
}

func Load() (Queries, error) {
	queries := make(Queries)
	err := fs.WalkDir(sqlFiles, "sql", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		contents, err := sqlFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read embedded query file %s: %w", path, err)
		}
		fileQueries, err := parseQueries(path, string(contents))
		if err != nil {
			return err
		}
		for name, query := range fileQueries {
			if _, exists := queries[name]; exists {
				return fmt.Errorf("duplicate query name %q", name)
			}
			queries[name] = query
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load embedded queries: %w", err)
	}

	for _, name := range requiredQueries {
		if _, exists := queries[name]; !exists {
			return nil, fmt.Errorf("query %q is not defined", name)
		}
	}
	return queries, nil
}

func parseQueries(filename, contents string) (map[string]string, error) {
	queries := make(map[string]string)
	var name string
	var statement strings.Builder

	save := func() error {
		if name == "" {
			if strings.TrimSpace(statement.String()) != "" {
				return fmt.Errorf("SQL without a -- name: marker in %s", filename)
			}
			return nil
		}
		query := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(statement.String()), ";"))
		if query == "" {
			return fmt.Errorf("query %q in %s is empty", name, filename)
		}
		if _, exists := queries[name]; exists {
			return fmt.Errorf("duplicate query name %q in %s", name, filename)
		}
		queries[name] = query
		name = ""
		statement.Reset()
		return nil
	}

	for _, line := range strings.Split(contents, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "-- name:") {
			if err := save(); err != nil {
				return nil, err
			}
			name = strings.TrimSpace(strings.TrimPrefix(trimmed, "-- name:"))
			if name == "" {
				return nil, fmt.Errorf("empty query name in %s", filename)
			}
			continue
		}
		if statement.Len() > 0 {
			statement.WriteByte('\n')
		}
		statement.WriteString(line)
	}
	if err := save(); err != nil {
		return nil, err
	}
	return queries, nil
}
