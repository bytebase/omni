package parser

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestSplitSegmentsCompileInOracle executes every segment Split returns, the
// way Bytebase runs a script, and checks that each PL/SQL unit compiled
// VALID. The SQL layer accepts a unit whose text is cut wrong and leaves it
// INVALID with PLS-00103 without raising an error: a call spec needs its
// ';', while a trigger's CALL routine and a wrapped unit must not have one.
//
// CREATE JAVA SOURCE and MLE code are not covered here: the
// gvenzl/oracle-free image has no Java VM (ORA-29538) and no MLE (ORA-00439).
// TestSplitEmbeddedSource pins their boundaries.
func TestSplitSegmentsCompileInOracle(t *testing.T) {
	ctx, db := openOracleReferenceDB(t)
	run := fmt.Sprintf("S%d", time.Now().UnixNano())

	var wrapped string
	wrapSource := "CREATE OR REPLACE PROCEDURE omni_wrap_" + run + " AS BEGIN NULL; END;"
	if err := db.QueryRowContext(ctx, "SELECT DBMS_DDL.WRAP(:1) FROM dual", wrapSource).Scan(&wrapped); err != nil {
		t.Fatalf("DBMS_DDL.WRAP: %v", err)
	}

	script := strings.ReplaceAll(`CREATE OR REPLACE PROCEDURE omni_target_{run} AS BEGIN NULL; END;
/
  CREATE OR REPLACE EDITIONABLE FUNCTION omni_jfn_{run} (p IN BLOB, q IN VARCHAR2) RETURN VARCHAR2
AS LANGUAGE JAVA
NAME 'Signer.check(oracle.sql.BLOB, java.lang.String) return java.lang.String';
/
/

CREATE OR REPLACE PROCEDURE omni_jproc_{run} (a IN VARCHAR2) AS LANGUAGE JAVA NAME 'X.y(java.lang.String)';
CREATE OR REPLACE PACKAGE omni_pk_{run} AS
  FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int';
  PROCEDURE p;
END;
/
CREATE OR REPLACE PACKAGE BODY omni_pk_{run} AS
  FUNCTION g RETURN NUMBER AS LANGUAGE JAVA NAME 'X.g() return int';
  PROCEDURE p IS BEGIN NULL; END;
END;
/
CREATE OR REPLACE TRIGGER omni_trg_{run} BEFORE INSERT ON t FOR EACH ROW CALL omni_target_{run}
/
`, "{run}", run) + strings.TrimRight(wrapped, "\n") + "\n/\n"

	for _, seg := range Split(script) {
		if seg.Kind == SegmentSQLPlusCommand || seg.Empty() {
			continue
		}
		if _, err := db.ExecContext(ctx, seg.Text); err != nil {
			t.Fatalf("exec segment %q: %v", seg.Text, err)
		}
	}

	rows, err := db.QueryContext(ctx, `
SELECT o.object_name, o.object_type, o.status,
       (SELECT LISTAGG(e.text, ' / ') WITHIN GROUP (ORDER BY e.sequence)
          FROM user_errors e WHERE e.name = o.object_name AND e.type = o.object_type)
  FROM user_objects o
 WHERE o.object_name LIKE '%\_' || :1 ESCAPE '\'
 ORDER BY o.object_name, o.object_type`, strings.ToUpper(run))
	if err != nil {
		t.Fatalf("query user_objects: %v", err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var name, kind, status string
		var errs *string
		if err := rows.Scan(&name, &kind, &status, &errs); err != nil {
			t.Fatalf("scan user_objects: %v", err)
		}
		got = append(got, kind+" "+name)
		if status != "VALID" {
			msg := ""
			if errs != nil {
				msg = *errs
			}
			t.Errorf("%s %s is %s: %s", kind, name, status, msg)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate user_objects: %v", err)
	}
	if len(got) != 7 {
		t.Fatalf("created objects = %q, want 7", got)
	}
}

// TestSplitExternalCallSpecInOracle executes Split's segments for a script
// whose EXTERNAL call spec opens with a clause other than NAME or LIBRARY. The
// library does not exist, so Oracle compiles the procedure with a semantic
// error only; a segment cut before its ';' would add PLS-00103, and a segment
// that swallowed the next statement would leave that statement uncreated.
func TestSplitExternalCallSpecInOracle(t *testing.T) {
	ctx, db := openOracleReferenceDB(t)
	run := fmt.Sprintf("X%d", time.Now().UnixNano())
	script := strings.ReplaceAll(`CREATE OR REPLACE PROCEDURE omni_ext_{run} AS EXTERNAL LANGUAGE C NAME "c_p" LIBRARY omni_nolib_{run};
CREATE OR REPLACE PROCEDURE omni_after_{run} AS BEGIN NULL; END;
`, "{run}", run)

	var segments int
	for _, seg := range Split(script) {
		if seg.Kind == SegmentSQLPlusCommand || seg.Empty() {
			continue
		}
		segments++
		if _, err := db.ExecContext(ctx, seg.Text); err != nil && !strings.Contains(err.Error(), "ORA-24344") {
			t.Fatalf("exec segment %q: %v", seg.Text, err)
		}
	}
	if segments != 2 {
		t.Fatalf("segments = %d, want 2", segments)
	}

	ext := strings.ToUpper("omni_ext_" + run)
	var syntaxErrors int
	if err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM user_errors WHERE name = :1 AND text LIKE 'PLS-00103%'`, ext).Scan(&syntaxErrors); err != nil {
		t.Fatalf("query user_errors: %v", err)
	}
	if syntaxErrors != 0 {
		t.Fatalf("%s compiled with PLS-00103: its segment lost the call spec's ';'", ext)
	}
	var status string
	if err := db.QueryRowContext(ctx,
		`SELECT status FROM user_objects WHERE object_name = :1 AND object_type = 'PROCEDURE'`,
		strings.ToUpper("omni_after_"+run)).Scan(&status); err != nil {
		t.Fatalf("procedure after the call spec was not created: %v", err)
	}
	if status != "VALID" {
		t.Fatalf("procedure after the call spec is %s", status)
	}
}
