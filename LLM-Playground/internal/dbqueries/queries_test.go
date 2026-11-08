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

	for name, prefix := range map[string]string{
		"CreateConversations": "CREATE TABLE",
		"FetchMessages":       "SELECT",
		"SaveConversation":    "INSERT",
		"SaveMessage":         "INSERT",
	} {
		if !strings.HasPrefix(queries[name], prefix) {
			t.Errorf("query %q does not start with %q: %q", name, prefix, queries[name])
		}
	}
}

func TestParseQueriesRequiresNames(t *testing.T) {
	if _, err := parseQueries("test.sql", "SELECT 1;"); err == nil {
		t.Fatal("expected unnamed SQL to be rejected")
	}
}
