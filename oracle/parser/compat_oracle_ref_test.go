//go:build oracle_ref

package parser

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOracleCompatibilityReferenceReport(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "compat_oracle.tsv"))
	ctx, db := openOracleReferenceDB(t)
	runID := fmt.Sprintf("R%d", time.Now().UnixNano())
	runOracleCompatibilityReferenceReport(t, ctx, db, rows, runID, oracleCompatLaneStandard)
}

func TestOracleCompatibilityPrivilegedReferenceReport(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "compat_oracle.tsv"))
	ctx, db := openOraclePrivilegedReferenceDB(t)
	runID := fmt.Sprintf("R%d", time.Now().UnixNano())
	runOracleCompatibilityReferenceReport(t, ctx, db, rows, runID, oracleCompatLanePrivileged)
}

func runOracleCompatibilityReferenceReport(t *testing.T, ctx context.Context, db *sql.DB, rows []coverageRow, runID string, lane string) {
	t.Helper()
	outcomes := make(map[string]int)
	byFamily := make(map[string]int)
	for _, row := range rows {
		expect := row.Fields["expect"]
		if lane == oracleCompatLaneStandard && (expect == "unsafe" || expect == "version_dependent") {
			outcomes[expect]++
			continue
		}
		if !oracleCompatLaneIncludesRow(lane, row) {
			continue
		}

		sqlText := oracleCompatSQL(row.Fields["sql"], runID)
		setupSQL := oracleCompatSQL(row.Fields["setup"], runID)
		cleanupSQL := oracleCompatSQL(row.Fields["cleanup"], runID)
		if lane == oracleCompatLanePrivileged {
			privilegedSetup, privilegedCleanup := oracleCompatPrivilegedFixture(row, runID)
			setupSQL = joinOracleCompatSQL(setupSQL, privilegedSetup)
			cleanupSQL = joinOracleCompatSQL(cleanupSQL, privilegedCleanup)
		}
		_ = execOracleCompatSQLList(ctx, db, cleanupSQL)
		if err := execOracleCompatSQLList(ctx, db, setupSQL); err != nil {
			t.Fatalf("%s: setup failed: %v setup=%q", row.Fields["id"], err, setupSQL)
		}
		oracleErr := oracleParseOnly(ctx, db, sqlText)
		if err := execOracleCompatSQLList(ctx, db, cleanupSQL); err != nil {
			t.Logf("%s: cleanup failed: %v cleanup=%q", row.Fields["id"], err, cleanupSQL)
		}
		_, omniErr := Parse(sqlText)
		oracleAccepts := oracleErr == nil
		omniAccepts := omniErr == nil

		outcome := oracleCompatibilityOutcome(expect, row.Fields["reference_class"], oracleAccepts, omniAccepts)
		outcomes[outcome]++
		byFamily[row.Fields["family"]+"/"+outcome]++
		if outcome != "match_accept" && outcome != "match_reject" {
			t.Logf("%s: %s oracle_accepts=%v oracle_err=%v omni_accepts=%v omni_err=%v sql=%q",
				row.Fields["id"], outcome, oracleAccepts, oracleErr, omniAccepts, omniErr, sqlText)
		}
	}

	if lane == oracleCompatLanePrivileged {
		t.Logf("Oracle privileged compatibility reference report: %s", sortedCountString(outcomes))
		t.Logf("Oracle privileged compatibility reference report by family: %s", sortedCountString(byFamily))
		return
	}
	t.Logf("Oracle compatibility reference report: %s", sortedCountString(outcomes))
	t.Logf("Oracle compatibility reference report by family: %s", sortedCountString(byFamily))
}

func joinOracleCompatSQL(parts ...string) string {
	var cleaned []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || strings.EqualFold(part, "none") {
			continue
		}
		cleaned = append(cleaned, part)
	}
	if len(cleaned) == 0 {
		return ""
	}
	return strings.Join(cleaned, ";")
}

func execOracleCompatSQLList(ctx context.Context, db *sql.DB, sqlText string) error {
	sqlText = strings.TrimSpace(sqlText)
	if sqlText == "" || strings.EqualFold(sqlText, "none") {
		return nil
	}
	var firstErr error
	for _, stmt := range strings.Split(sqlText, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.ExecContext(ctx, stmt); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func oracleCompatibilityOutcome(expect string, referenceClass string, oracleAccepts bool, omniAccepts bool) string {
	if referenceClass != "" && referenceClass != "syntax" && !oracleAccepts && omniAccepts {
		return "oracle_" + referenceClass
	}
	switch {
	case expect == "accept" && oracleAccepts && omniAccepts:
		return "match_accept"
	case expect == "reject" && !oracleAccepts && !omniAccepts:
		return "match_reject"
	case oracleAccepts && !omniAccepts:
		return "omni_too_strict"
	case !oracleAccepts && omniAccepts:
		return "omni_too_lenient"
	case expect == "accept" && !oracleAccepts && !omniAccepts:
		return "oracle_unexpected_reject"
	case expect == "reject" && oracleAccepts && omniAccepts:
		return "oracle_unexpected_accept"
	default:
		return "mismatch_other"
	}
}
