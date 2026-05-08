//go:build oracle_ref

package parser

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

func TestOracleReference(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "reference_oracle.tsv"))
	ctx, db := openOracleReferenceDB(t)
	runID := fmt.Sprintf("R%d", time.Now().UnixNano())

	for _, row := range rows {
		if oracleReferenceRowSkipped(row.Fields["expect"]) {
			continue
		}
		t.Run(row.Key, func(t *testing.T) {
			sqlText := oracleReferenceSQL(row.Fields["sql"], runID)
			oracleErr := oracleParseOnly(ctx, db, sqlText)
			_, omniErr := Parse(sqlText)
			oracleAccepts := oracleErr == nil
			omniAccepts := omniErr == nil
			if oracleAccepts != omniAccepts {
				t.Fatalf("reference mismatch: oracle_accepts=%v oracle_err=%v omni_accepts=%v omni_err=%v sql=%q",
					oracleAccepts, oracleErr, omniAccepts, omniErr, sqlText)
			}
		})
	}
}

func oracleReferenceRowSkipped(expect string) bool {
	return expect == "catalog" || expect == "unsafe"
}

func TestOracleVReservedWordsKeywordAudit(t *testing.T) {
	ctx, db := openOracleReservedWordsDB(t)

	rows, err := db.QueryContext(ctx, `
SELECT keyword, reserved, res_type, res_attr, res_semi
FROM v$reserved_words
WHERE duplicate = 'N' OR duplicate IS NULL`)
	if err != nil {
		t.Fatalf("query v$reserved_words: %v", err)
	}
	defer rows.Close()

	manifestRows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "oracle_keywords.tsv"))
	manifest := make(map[string]coverageRow, len(manifestRows))
	for _, row := range manifestRows {
		manifest[row.Fields["word"]] = row
	}

	var checked int
	var missing []string
	for rows.Next() {
		var keyword, reserved, resType, resAttr, resSemi sql.NullString
		if err := rows.Scan(&keyword, &reserved, &resType, &resAttr, &resSemi); err != nil {
			t.Fatalf("scan v$reserved_words: %v", err)
		}
		if !keyword.Valid || keyword.String == "" {
			continue
		}
		if !isOracleReferenceWordKeyword(keyword.String) {
			continue
		}
		if reserved.String == "Y" || resType.String == "Y" || resAttr.String == "Y" || resSemi.String == "Y" {
			checked++
			if _, ok := manifest[keyword.String]; !ok {
				missing = append(missing, keyword.String)
			}
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate v$reserved_words: %v", err)
	}
	if len(missing) > 0 {
		t.Fatalf("Oracle V$RESERVED_WORDS reserved entries missing from local manifest: %v", missing)
	}
	t.Logf("Oracle V$RESERVED_WORDS reserved/context audit checked=%d", checked)
}

func openOracleReferenceDB(t *testing.T) (context.Context, *sql.DB) {
	t.Helper()
	if dsn := os.Getenv("ORACLE_PARSER_REF_DSN"); dsn != "" {
		db, err := sql.Open("oracle", dsn)
		if err != nil {
			t.Fatalf("open Oracle reference connection: %v", err)
		}
		t.Cleanup(func() { _ = db.Close() })

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		t.Cleanup(cancel)
		if err := db.PingContext(ctx); err != nil {
			t.Fatalf("ping Oracle reference connection: %v", err)
		}
		return ctx, db
	}
	if os.Getenv("ORACLE_PARSER_REF_CONTAINER") == "1" {
		oracle := startOracleDB(t)
		return oracle.ctx, oracle.db
	}
	if os.Getenv("ORACLE_PARSER_REF_STRICT") == "1" {
		t.Fatal("strict Oracle reference mode requires ORACLE_PARSER_REF_DSN or ORACLE_PARSER_REF_CONTAINER=1")
	}
	t.Skip("ORACLE_PARSER_REF_DSN is not set; set ORACLE_PARSER_REF_CONTAINER=1 to run against Oracle Free")
	return nil, nil
}

func openOracleReservedWordsDB(t *testing.T) (context.Context, *sql.DB) {
	t.Helper()
	if dsn := os.Getenv("ORACLE_PARSER_REF_ADMIN_DSN"); dsn != "" {
		db, err := sql.Open("oracle", dsn)
		if err != nil {
			t.Fatalf("open Oracle admin reference connection: %v", err)
		}
		t.Cleanup(func() { _ = db.Close() })

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		t.Cleanup(cancel)
		if err := db.PingContext(ctx); err != nil {
			t.Fatalf("ping Oracle admin reference connection: %v", err)
		}
		return ctx, db
	}
	if os.Getenv("ORACLE_PARSER_REF_CONTAINER") == "1" {
		oracle := startOracleDB(t)
		return oracle.ctx, oracle.adminDB
	}
	return openOracleReferenceDB(t)
}

func oracleReferenceSQL(sqlText, runID string) string {
	return strings.ReplaceAll(sqlText, "{run}", runID)
}

func isOracleReferenceWordKeyword(keyword string) bool {
	for i := 0; i < len(keyword); i++ {
		ch := keyword[i]
		if ch >= 'A' && ch <= 'Z' {
			continue
		}
		if i > 0 && ch >= '0' && ch <= '9' {
			continue
		}
		if ch == '_' || ch == '$' || ch == '#' {
			continue
		}
		return false
	}
	return true
}

func oracleParseOnly(ctx context.Context, db *sql.DB, sqlText string) error {
	const block = `
DECLARE
  c INTEGER;
BEGIN
  c := DBMS_SQL.OPEN_CURSOR;
  BEGIN
    DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
    DBMS_SQL.CLOSE_CURSOR(c);
  EXCEPTION
    WHEN OTHERS THEN
      IF DBMS_SQL.IS_OPEN(c) THEN
        DBMS_SQL.CLOSE_CURSOR(c);
      END IF;
      RAISE;
  END;
END;`
	_, err := db.ExecContext(ctx, block, sqlText)
	return err
}
