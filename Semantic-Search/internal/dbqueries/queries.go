package dbqueries

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

type Queries struct {
	Create CreateQueries
	Fetch  FetchQueries
	Save   SaveQueries
	Update UpdateQueries
	Delete DeleteQueries
}

type CreateQueries struct {
	VectorExtension string
	DocumentsTable  string
	ChunksTable     string
}

type FetchQueries struct {
	SimilarDocuments         string
	ChunkedDocuments         string
	FilteredChunkedDocuments string
	DocumentByID             string
}

type SaveQueries struct {
	Document string
	Chunk    string
}

type UpdateQueries struct {
	Document       string
	IndexingStatus string
}

type DeleteQueries struct {
	Document string
	Chunks   string
}

func Load() (Queries, error) {
	namedQueries := make(map[string]string)
	err := fs.WalkDir(sqlFiles, "sql", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".sql") {
			return nil
		}

		contents, err := sqlFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read embedded query file %s: %w", path, err)
		}
		queries, err := parseQueries(path, string(contents))
		if err != nil {
			return err
		}
		for name, query := range queries {
			if _, exists := namedQueries[name]; exists {
				return fmt.Errorf("duplicate query name %q", name)
			}
			namedQueries[name] = query
		}
		return nil
	})
	if err != nil {
		return Queries{}, fmt.Errorf("load embedded queries: %w", err)
	}

	required := func(name string) (string, error) {
		query, exists := namedQueries[name]
		if !exists {
			return "", fmt.Errorf("query %q is not defined", name)
		}
		return query, nil
	}
	var queries Queries
	var queryErr error
	if queries.Create.VectorExtension, queryErr = required("CreateVectorExtension"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Create.DocumentsTable, queryErr = required("CreateDocumentsTable"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Create.ChunksTable, queryErr = required("CreateChunksTable"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Fetch.SimilarDocuments, queryErr = required("FetchSimilarDocuments"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Fetch.ChunkedDocuments, queryErr = required("FetchChunkedDocuments"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Fetch.FilteredChunkedDocuments, queryErr = required("FetchFilteredChunkedDocuments"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Fetch.DocumentByID, queryErr = required("FetchDocumentByID"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Save.Document, queryErr = required("SaveDocument"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Save.Chunk, queryErr = required("SaveChunk"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Update.Document, queryErr = required("UpdateDocument"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Update.IndexingStatus, queryErr = required("UpdateIndexingStatus"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Delete.Document, queryErr = required("DeleteDocument"); queryErr != nil {
		return Queries{}, queryErr
	}
	if queries.Delete.Chunks, queryErr = required("DeleteChunks"); queryErr != nil {
		return Queries{}, queryErr
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
