package completion

import (
	"testing"

	"github.com/bytebase/omni/redshift/catalog"
)

func TestRedshiftTopLevelCompletion(t *testing.T) {
	candidates := Complete("", 0, nil)
	for _, kw := range []string{"SELECT", "INSERT", "UPDATE", "DELETE", "MERGE", "CREATE", "ALTER", "DROP", "COPY", "UNLOAD", "SHOW"} {
		if !hasCandidate(candidates, kw, CandidateKeyword) {
			t.Fatalf("missing top-level keyword %q in %#v", kw, candidates)
		}
	}
}

func TestRedshiftCreateTableOptionCompletion(t *testing.T) {
	candidates := Complete("CREATE TABLE t (id INT) ", len("CREATE TABLE t (id INT) "), nil)
	for _, kw := range []string{"DISTSTYLE", "DISTKEY", "SORTKEY", "ENCODE"} {
		if !hasCandidate(candidates, kw, CandidateKeyword) {
			t.Fatalf("missing CREATE TABLE option %q in %#v", kw, candidates)
		}
	}

	candidates = Complete("CREATE TABLE t (id INT ", len("CREATE TABLE t (id INT "), nil)
	for _, kw := range []string{"ENCODE", "DISTKEY", "SORTKEY"} {
		if !hasCandidate(candidates, kw, CandidateKeyword) {
			t.Fatalf("missing CREATE TABLE column option %q in %#v", kw, candidates)
		}
	}
}

func TestRedshiftCopyUnloadOptionCompletion(t *testing.T) {
	copyCandidates := Complete("COPY t FROM 's3://bucket/file' ", len("COPY t FROM 's3://bucket/file' "), nil)
	for _, kw := range []string{"IAM_ROLE", "CREDENTIALS", "FORMAT", "DELIMITER"} {
		if !hasCandidate(copyCandidates, kw, CandidateKeyword) {
			t.Fatalf("missing COPY option %q in %#v", kw, copyCandidates)
		}
	}

	unloadCandidates := Complete("UNLOAD ('SELECT * FROM t') TO 's3://bucket/out' ", len("UNLOAD ('SELECT * FROM t') TO 's3://bucket/out' "), nil)
	for _, kw := range []string{"IAM_ROLE", "MANIFEST", "HEADER", "FORMAT"} {
		if !hasCandidate(unloadCandidates, kw, CandidateKeyword) {
			t.Fatalf("missing UNLOAD option %q in %#v", kw, unloadCandidates)
		}
	}
}

func TestRedshiftShowSubcommandCompletion(t *testing.T) {
	candidates := Complete("SHOW ", len("SHOW "), nil)
	for _, kw := range []string{"DATABASES", "SCHEMAS", "TABLES", "COLUMNS", "GRANTS", "DATASHARES"} {
		if !hasCandidate(candidates, kw, CandidateKeyword) {
			t.Fatalf("missing SHOW subcommand %q in %#v", kw, candidates)
		}
	}
}

func TestRedshiftCompletionVirtualColumns(t *testing.T) {
	cat := catalog.New()
	if _, err := cat.Exec(`CREATE TABLE users (id int, name text);`, nil); err != nil {
		t.Fatal(err)
	}

	cteCandidates := Complete("WITH active(x, y) AS (SELECT id, name FROM users) SELECT  FROM active", len("WITH active(x, y) AS (SELECT id, name FROM users) SELECT "), cat)
	for _, col := range []string{"x", "y"} {
		if !hasCandidate(cteCandidates, col, CandidateColumn) {
			t.Fatalf("missing CTE alias column %q in %#v", col, cteCandidates)
		}
	}

	subqueryCandidates := Complete("SELECT  FROM (SELECT id AS user_id, name FROM users) u", len("SELECT "), cat)
	for _, col := range []string{"user_id", "name"} {
		if !hasCandidate(subqueryCandidates, col, CandidateColumn) {
			t.Fatalf("missing subquery alias column %q in %#v", col, subqueryCandidates)
		}
	}

	exprCandidates := Complete("SELECT  FROM (SELECT 1 + 2) u", len("SELECT "), cat)
	if hasCandidate(exprCandidates, "?column?", CandidateColumn) || hasCandidate(exprCandidates, "1 + 2", CandidateColumn) {
		t.Fatalf("invented unsupported expression column in %#v", exprCandidates)
	}
}

func TestRedshiftCompletionMetadataSchemasAndRelations(t *testing.T) {
	cat := catalog.New()
	if _, err := cat.Exec(`
		CREATE TABLE users (id int);
		CREATE SCHEMA analytics;
		CREATE TABLE analytics.events (id int);
	`, nil); err != nil {
		t.Fatal(err)
	}

	candidates := Complete("SELECT * FROM ", len("SELECT * FROM "), cat)
	if !hasCandidate(candidates, "analytics", CandidateSchema) {
		t.Fatalf("missing schema candidate analytics in %#v", candidates)
	}
	if !hasCandidate(candidates, "users", CandidateTable) {
		t.Fatalf("missing public table candidate users in %#v", candidates)
	}

	qualifiedCandidates := Complete("SELECT * FROM analytics.", len("SELECT * FROM analytics."), cat)
	if !hasCandidate(qualifiedCandidates, "events", CandidateTable) {
		t.Fatalf("missing analytics table candidate events in %#v", qualifiedCandidates)
	}
	if hasCandidate(qualifiedCandidates, "users", CandidateTable) {
		t.Fatalf("unexpected public table candidate users after analytics. in %#v", qualifiedCandidates)
	}
}
