package dbqueries

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	queries, err := Load()
	if err != nil {
		t.Fatalf("load queries: %v", err)
	}

	if !strings.HasPrefix(queries.Create.DocumentsTable, "CREATE TABLE") {
		t.Errorf("unexpected create documents query: %q", queries.Create.DocumentsTable)
	}
	if !strings.HasPrefix(queries.Fetch.SimilarDocuments, "SELECT") {
		t.Errorf("unexpected fetch similar documents query: %q", queries.Fetch.SimilarDocuments)
	}
	if !strings.HasPrefix(queries.Save.Document, "INSERT") {
		t.Errorf("unexpected save document query: %q", queries.Save.Document)
	}
	if !strings.HasPrefix(queries.Update.Document, "UPDATE") {
		t.Errorf("unexpected update document query: %q", queries.Update.Document)
	}
	if !strings.HasPrefix(queries.Delete.Document, "DELETE") {
		t.Errorf("unexpected delete document query: %q", queries.Delete.Document)
	}
}

func TestParseQueriesRequiresNames(t *testing.T) {
	if _, err := parseQueries("test.sql", "SELECT 1;"); err == nil {
		t.Fatal("expected unnamed SQL to be rejected")
	}
}
