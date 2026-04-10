package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// TestVerifyCorpus is the multi-dimensional corpus verifier. It reads SQL from
// corpus files and runs each statement through multiple verification layers:
//
//  1. Parse — does the parser accept/reject as expected? (@valid annotation)
//  2. Crash — the parser must never panic on any input
//  3. Loc  — all AST nodes must have valid Loc (Start >= 0, End > Start)
//  4. Error — if @valid: false, error position must be reasonable (>= 0)
//
// Corpus files live in mysql/quality/corpus/*.sql with annotations:
//
//	-- @name: descriptive name
//	-- @valid: true|false       (expected parse result; omit = just check no crash)
//	-- @source: where this SQL comes from
//
// Run:
//
//	go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
func TestVerifyCorpus(t *testing.T) {
	corpusDir := filepath.Join("..", "quality", "corpus")
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		corpusDir = filepath.Join("mysql", "quality", "corpus")
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
		t.Skip("No corpus files found in mysql/quality/corpus/")
	}

	var stats corpusVerifyStats
	for _, file := range corpusFiles {
		statements := loadCorpusStatements(t, file)
		baseName := filepath.Base(file)

		for _, stmt := range statements {
			name := stmt.name
			if name == "" {
				name = truncateCorpusSQL(stmt.sql, 60)
			}
			testName := baseName + "/" + name

			t.Run(testName, func(t *testing.T) {
				stats.total++
				result := verifyCorpusStatement(t, stmt)
				accumulateCorpusStats(&stats, result)
			})
		}
	}

	// Print summary
	t.Logf("\n=== Corpus Verifier Summary ===")
	t.Logf("Total statements:     %d", stats.total)
	t.Logf("Parse OK:             %d", stats.parseOK)
	t.Logf("Parse expected fail:  %d", stats.parseExpectedFail)
	t.Logf("PARSE VIOLATIONS:     %d (should accept but rejects)", stats.parseViolations)
	t.Logf("PARSE LENIENT:        %d (should reject but accepts)", stats.parseLenient)
	t.Logf("Loc clean:            %d (all nodes have valid Loc)", stats.locClean)
	t.Logf("LOC VIOLATIONS:       %d (nodes with bad/missing Loc)", stats.locViolations)
	t.Logf("CRASHES:              %d", stats.crashes)
	if stats.parseViolations > 0 || stats.crashes > 0 {
		t.Logf("\n%d issue(s) require attention", stats.parseViolations+stats.crashes)
	}
}

type corpusVerifyStats struct {
	total             int
	parseOK           int
	parseExpectedFail int
	parseViolations   int
	parseLenient      int
	locClean          int
	locViolations     int
	crashes           int
}

type corpusVerifyResult struct {
	crashed       bool
	parsed        bool
	parseErr      error
	locViolations []LocViolation
	stmt          corpusStmt
}

// verifyCorpusStatement runs all verification layers on a single corpus statement.
func verifyCorpusStatement(t *testing.T, stmt corpusStmt) corpusVerifyResult {
	t.Helper()
	var result corpusVerifyResult
	result.stmt = stmt

	// --- Layer 1: Crash check (always) ---
	var parseResult *ast.List
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CRASH: parser panicked on: %s\npanic: %v",
					truncateCorpusSQL(stmt.sql, 200), r)
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
		violations := checkCorpusLoc(parseResult)
		result.locViolations = violations
		if len(violations) > 0 {
			t.Logf("LOC VIOLATIONS: %d nodes with invalid Loc:", len(violations))
			for i, v := range violations {
				if i >= 5 {
					t.Logf("  ... and %d more", len(violations)-5)
					break
				}
				t.Logf("  %s", v)
			}
		}
	}

	return result
}

// checkCorpusLoc runs the Loc walker on parse results and returns violations.
func checkCorpusLoc(result *ast.List) []LocViolation {
	var violations []LocViolation
	if result == nil {
		return violations
	}
	for i, item := range result.Items {
		path := fmt.Sprintf("Items[%d]", i)
		walkNodeLocs(reflect.ValueOf(item), path, &violations)
	}
	return violations
}

func accumulateCorpusStats(stats *corpusVerifyStats, r corpusVerifyResult) {
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
		if r.parsed {
			stats.parseLenient++
		}
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

// corpusStmt is a single SQL statement from a corpus file.
type corpusStmt struct {
	sql    string
	name   string // from "-- @name: ..."
	valid  string // "true", "false", or "" (unknown)
	source string // from "-- @source: ..."
}

// loadCorpusStatements reads a .sql corpus file and extracts individual statements.
func loadCorpusStatements(t *testing.T, path string) []corpusStmt {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Cannot open corpus file %s: %v", path, err)
	}
	defer f.Close()

	var statements []corpusStmt
	var current corpusStmt
	var lines []string

	flush := func() {
		if sql := buildCorpusSQL(lines); sql != "" {
			current.sql = sql
			statements = append(statements, current)
		}
		current = corpusStmt{}
		lines = nil
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

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

func buildCorpusSQL(lines []string) string {
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func truncateCorpusSQL(sql string, maxLen int) string {
	s := strings.Join(strings.Fields(sql), " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
