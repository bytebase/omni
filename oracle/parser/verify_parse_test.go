package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestVerifyCorpus is the multi-dimensional corpus verifier. It reads SQL from
// corpus files and runs each statement through multiple verification layers:
//
//  1. Parse — does the parser accept/reject as expected? (@valid annotation)
//  2. Crash — the parser must never panic on any input
//  3. Loc  — all AST nodes must have valid Loc (Start >= 0, End > Start)
//  4. Error — if @valid: false, error position must be reasonable (>= 0)
//
// Corpus files live in oracle/quality/corpus/*.sql with annotations:
//
//	-- @name: descriptive name
//	-- @valid: true|false       (expected parse result; omit = just check no crash)
//	-- @source: where this SQL comes from
//
// Run:
//
//	go test ./oracle/parser/ -run TestVerifyCorpus -count=1 -v
func TestVerifyCorpus(t *testing.T) {
	corpusDir := filepath.Join("..", "quality", "corpus")
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		corpusDir = filepath.Join("oracle", "quality", "corpus")
		entries, err = os.ReadDir(corpusDir)
		if err != nil {
			t.Fatalf("Cannot read corpus directory: %v", err)
		}
	}

	var corpusFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			corpusFiles = append(corpusFiles, filepath.Join(corpusDir, e.Name()))
		}
	}
	if len(corpusFiles) == 0 {
		t.Skip("No corpus files found in oracle/quality/corpus/")
	}

	var stats verifyStats
	for _, file := range corpusFiles {
		statements := loadCorpusFile(t, file)
		baseName := filepath.Base(file)

		for _, stmt := range statements {
			name := stmt.name
			if name == "" {
				name = truncateSQL(stmt.sql, 60)
			}
			testName := baseName + "/" + name

			t.Run(testName, func(t *testing.T) {
				stats.total++
				result := verifyStatement(t, stmt)
				accumulateStats(&stats, result)
			})
		}
	}

	// Print summary
	t.Logf("\n=== Corpus Verifier Summary ===")
	t.Logf("Total statements:     %d", stats.total)
	t.Logf("Parse OK:             %d", stats.parseOK)
	t.Logf("Parse expected fail:  %d", stats.parseExpectedFail)
	t.Logf("PARSE VIOLATIONS:     %d (should accept but rejects)", stats.parseViolations)
	t.Logf("Loc clean:            %d (all nodes have valid Loc)", stats.locClean)
	t.Logf("LOC VIOLATIONS:       %d (nodes with bad/missing Loc)", stats.locViolations)
	t.Logf("CRASHES:              %d", stats.crashes)
	if stats.parseViolations > 0 || stats.crashes > 0 {
		t.Logf("\n⚠ %d issue(s) require attention", stats.parseViolations+stats.crashes)
	}
}

type verifyStats struct {
	total             int
	parseOK           int
	parseExpectedFail int
	parseViolations   int
	locClean          int
	locViolations     int
	crashes           int
}

type verifyResult struct {
	crashed       bool
	parsed        bool
	parseErr      error
	locViolations []LocViolation
	stmt          corpusStatement
}

// verifyStatement runs all verification layers on a single corpus statement.
func verifyStatement(t *testing.T, stmt corpusStatement) verifyResult {
	t.Helper()
	var result verifyResult
	result.stmt = stmt

	// --- Layer 1: Crash check (always) ---
	var parseResult *ast.List
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CRASH: parser panicked on: %s\npanic: %v",
					truncateSQL(stmt.sql, 200), r)
				result.crashed = true
			}
		}()
		var err error
		parseResult, err = Parse(stmt.sql)
		result.parseErr = err
		result.parsed = err == nil
	}()

	if result.crashed {
		return result
	}

	// --- Layer 2: Parse check (if @valid is annotated) ---
	switch stmt.valid {
	case "true":
		if !result.parsed {
			t.Errorf("PARSE VIOLATION: @valid: true but parser rejects: %v", result.parseErr)
		}
	case "false":
		if result.parsed {
			// Parser accepts invalid SQL — this is informational, not a failure.
			// Parsers are often more lenient than the DB, which is usually OK.
			t.Logf("PARSE LENIENT: @valid: false but parser accepts")
		}
		if !result.parsed && result.parseErr != nil {
			// Verify error has a reasonable position
			if pe, ok := result.parseErr.(*ParseError); ok {
				if pe.Position < 0 {
					t.Errorf("ERROR POSITION: @valid: false, error position = %d (want >= 0)", pe.Position)
				}
			}
		}
	default:
		// No @valid annotation — just check no crash (already done above)
	}

	// --- Layer 3: Loc check (if parse succeeded) ---
	if result.parsed && parseResult != nil {
		violations := checkLocOnResult(parseResult, stmt.sql)
		result.locViolations = violations
		if len(violations) > 0 {
			t.Errorf("LOC VIOLATIONS: %d nodes with invalid Loc:", len(violations))
			for i, v := range violations {
				if i >= 5 {
					t.Errorf("  ... and %d more", len(violations)-5)
					break
				}
				t.Errorf("  %s", v)
			}
		}
	}

	return result
}

// checkLocOnResult runs the Loc walker on parse results and returns violations.
func checkLocOnResult(result *ast.List, sql string) []LocViolation {
	var violations []LocViolation
	if result == nil {
		return violations
	}

	// Use the existing reflection walker from loc_walker_test.go
	for i, item := range result.Items {
		path := fmt.Sprintf("Items[%d]", i)
		walkNodeLocs(reflect.ValueOf(item), path, &violations)
	}


	return violations
}

func accumulateStats(stats *verifyStats, r verifyResult) {
	if r.crashed {
		stats.crashes++
		return
	}

	switch r.stmt.valid {
	case "true":
		if r.parsed {
			stats.parseOK++
		} else {
			stats.parseViolations++
		}
	case "false":
		stats.parseExpectedFail++
	default:
		if r.parsed {
			stats.parseOK++
		}
	}

	if r.parsed && len(r.locViolations) == 0 {
		stats.locClean++
	} else if r.parsed && len(r.locViolations) > 0 {
		stats.locViolations++
	}
}

// --- Corpus file parsing ---

// corpusStatement is a single SQL statement from a corpus file.
type corpusStatement struct {
	sql    string
	name   string // from "-- @name: ..."
	valid  string // "true", "false", or "" (unknown)
	source string // from "-- @source: ..."
}

// loadCorpusFile reads a .sql corpus file and extracts individual statements.
//
// Format:
//
//	-- @name: descriptive name
//	-- @valid: true
//	-- @source: Oracle 23ai SQL Reference
//	SELECT * FROM employees
//	WHERE dept_id = 10;
//
//	-- @name: next statement
//	INSERT INTO t VALUES (1);
//
// Statements are separated by blank lines or new "-- @name:" annotations.
// Lines starting with "--" (but not "-- @") are SQL comments and included in the statement.
func loadCorpusFile(t *testing.T, path string) []corpusStatement {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Cannot open corpus file %s: %v", path, err)
	}
	defer f.Close()

	var statements []corpusStatement
	var current corpusStatement
	var lines []string

	flush := func() {
		if sql := buildSQL(lines); sql != "" {
			current.sql = sql
			statements = append(statements, current)
		}
		current = corpusStatement{}
		lines = nil
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Metadata annotations
		if strings.HasPrefix(line, "-- @name:") {
			flush()
			current.name = strings.TrimSpace(strings.TrimPrefix(line, "-- @name:"))
			continue
		}
		if strings.HasPrefix(line, "-- @valid:") {
			current.valid = strings.TrimSpace(strings.TrimPrefix(line, "-- @valid:"))
			continue
		}
		if strings.HasPrefix(line, "-- @source:") {
			current.source = strings.TrimSpace(strings.TrimPrefix(line, "-- @source:"))
			continue
		}

		// Blank line = statement separator
		if strings.TrimSpace(line) == "" {
			if len(lines) > 0 {
				flush()
			}
			continue
		}

		lines = append(lines, line)
	}

	flush()
	return statements
}

// buildSQL joins lines into a single SQL string.
func buildSQL(lines []string) string {
	sql := strings.TrimSpace(strings.Join(lines, "\n"))
	return sql
}

// truncateSQL returns a shortened version of sql for display.
func truncateSQL(sql string, maxLen int) string {
	s := strings.Join(strings.Fields(sql), " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
