package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB v8.5.0 stored-program grammar matrix. CREATE PROCEDURE bodies are parsed
// via the compound-statement grammar (parseCompoundStmtOrStmt), not a raw-text
// scanner. Every row below was probed on a pingcap/tidb:v8.5.0 container:
// accepted forms report ERROR 8108 (parse-OK, plan-unsupported), rejected forms
// report ERROR 1064 at parse time. TiDB's grammar is a strict subset of MySQL's:
// it has no LOOP, DECLARE...CONDITION, SIGNAL/RESIGNAL, and no CREATE
// FUNCTION/TRIGGER/EVENT at all. Verified on TiDB v8.5.0.

func TestRoutineCompoundBodyAccepted(t *testing.T) {
	cases := []string{
		`CREATE PROCEDURE p() BEGIN SELECT 1; END`,
		`CREATE PROCEDURE p() BEGIN IF 1>0 THEN SELECT 1; END IF; END`,
		`CREATE PROCEDURE p() BEGIN CASE 1 WHEN 1 THEN SELECT 1; END CASE; END`,
		`CREATE PROCEDURE p() BEGIN CASE WHEN 1>0 THEN SELECT 1; ELSE SELECT 2; END CASE; END`,
		`CREATE PROCEDURE p() BEGIN DECLARE x INT DEFAULT 3; WHILE x>0 DO SET x=x-1; END WHILE; END`,
		`CREATE PROCEDURE p() BEGIN DECLARE x INT DEFAULT 0; REPEAT SET x=x+1; UNTIL x>5 END REPEAT; END`,
		`CREATE PROCEDURE p() lbl: BEGIN SELECT 1; END lbl`,
		`CREATE PROCEDURE p() BEGIN END`,
		`CREATE PROCEDURE p() BEGIN IF EXISTS (SELECT 1 FROM t) THEN SELECT 1; END IF; END`,
		`CREATE PROCEDURE p() BEGIN IF (1>0) THEN SELECT 1; END IF; END`,
		`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR SQLEXCEPTION SET @x=1; SELECT 1; END`,
		`CREATE PROCEDURE p() BEGIN DECLARE cur CURSOR FOR SELECT a FROM t; OPEN cur; CLOSE cur; END`,
		`CREATE PROCEDURE p() w: WHILE 1=0 DO IF 1>0 THEN LEAVE w; END IF; END WHILE w`,
		`CREATE PROCEDURE p() BEGIN IF 1=1 THEN WHILE 1=0 DO BEGIN DECLARE y INT; SET y=1; END; END WHILE; END IF; END`,
		// Comments inside the body are handled by the lexer, not a text scanner.
		`CREATE PROCEDURE p() BEGIN IF /*!50000 EXISTS */ (SELECT 1 FROM t) THEN SELECT 1; END IF; END`,
		`CREATE PROCEDURE p() BEGIN SELECT IF /*c*/ (1,2,3); END`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

func TestRoutineGrammarRejected(t *testing.T) {
	// Forms MySQL accepts but TiDB v8.5.0 rejects (1064). omni must reject them
	// to match TiDB, not inherit MySQL's grammar.
	cases := []string{
		`CREATE PROCEDURE p() BEGIN LOOP SELECT 1; END LOOP; END`,
		`CREATE PROCEDURE p() BEGIN l: LOOP LEAVE l; END LOOP l; END`,
		`CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR SQLSTATE '45000'; SELECT 1; END`,
		`CREATE PROCEDURE p() BEGIN SIGNAL SQLSTATE '45000'; END`,
		`CREATE PROCEDURE p() BEGIN RESIGNAL; END`,
		`CREATE FUNCTION f() RETURNS INT BEGIN RETURN 1; END`,
		`CREATE FUNCTION f RETURNS INTEGER SONAME 'x.so'`,
		`CREATE AGGREGATE FUNCTION agg() RETURNS INT SONAME 'x.so'`,
		`CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @x=1`,
		`CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1`,
		`ALTER EVENT e ON SCHEDULE EVERY 2 DAY`,
		`SIGNAL SQLSTATE '45000'`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but TiDB v8.5.0 rejects it (1064)", sql)
			}
		})
	}
}

// A CREATE PROCEDURE body is parsed into an AST node (Body) and the raw bytes
// preserved (BodyText) for catalog re-emission.
func TestRoutineBodyParsedAsNode(t *testing.T) {
	const sql = `CREATE PROCEDURE p() BEGIN IF 1>0 THEN SELECT 1; END IF; END`
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	proc, ok := list.Items[0].(*nodes.CreateFunctionStmt)
	if !ok {
		t.Fatalf("statement is %T, want *nodes.CreateFunctionStmt", list.Items[0])
	}
	if proc.Body == nil {
		t.Errorf("Body is nil, want a parsed compound-statement node")
	}
	if proc.BodyText == "" {
		t.Errorf("BodyText is empty, want the raw body bytes")
	}
}

// The body parse must consume exactly the routine body — a statement following
// the routine is parsed separately (the grammar handles the boundary; no text
// scanner over/under-consumes, including across comment boundaries).
func TestRoutineBodyDoesNotSwallowTrailingStatement(t *testing.T) {
	cases := []string{
		`CREATE PROCEDURE p() BEGIN IF 1>0 THEN SELECT 1; END IF; END; SELECT 2`,
		`CREATE PROCEDURE p() BEGIN /*!50003 SET @x = 1 */; END; SELECT 2`,
		"DELIMITER //\nCREATE PROCEDURE p() BEGIN SELECT 1; END//\nSELECT 2//",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			list, err := Parse(sql)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
			if len(list.Items) != 2 {
				t.Fatalf("Parse(%q) produced %d statements, want 2 (procedure + trailing)", sql, len(list.Items))
			}
			if _, ok := list.Items[0].(*nodes.CreateFunctionStmt); !ok {
				t.Errorf("first statement is %T, want *nodes.CreateFunctionStmt", list.Items[0])
			}
		})
	}
}
