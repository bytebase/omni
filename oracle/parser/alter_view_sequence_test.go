package parser

import (
	"errors"
	"testing"
)

// TestAdminTailTwoOracleAnswers pins CREATE CONTEXT, CREATE RESTORE POINT,
// CREATE and ALTER ROLLBACK SEGMENT, CREATE and ALTER INMEMORY JOIN GROUP,
// ALTER VIEW, and ALTER SEQUENCE forms to Oracle 23ai's answers, probed on
// 23.26.3. The accepted forms fail on the server only past the parser:
// ORA-01031 (privilege), ORA-00942 or ORA-02289 (no such view or sequence),
// ORA-00947.
func TestAdminTailTwoOracleAnswers(t *testing.T) {
	accepts := []string{
		"CREATE CONTEXT a USING a.b INITIALIZED EXTERNALLY",
		`CREATE OR REPLACE CONTEXT "a" USING "b" ACCESSED GLOBALLY`,
		"CREATE CLEAN RESTORE POINT a FOR PLUGGABLE DATABASE b AS OF SCN 1 + 1 GUARANTEE FLASHBACK DATABASE",
		`CREATE RESTORE POINT "a" AS OF TIMESTAMP SYSTIMESTAMP - 1 PRESERVE`,
		"CREATE ROLLBACK SEGMENT a STORAGE (INITIAL 1M) TABLESPACE b STORAGE (NEXT 1M)",
		"CREATE PUBLIC ROLLBACK SEGMENT a TABLESPACE b",
		"ALTER ROLLBACK SEGMENT a STORAGE (MAXEXTENTS UNLIMITED)",
		"ALTER ROLLBACK SEGMENT a SHRINK TO 1M",
		"CREATE INMEMORY JOIN GROUP IF NOT EXISTS a(a.b(c), d.e(f))",
		`CREATE INMEMORY JOIN GROUP a(a.b(c, d), e.f("g", h))`,
		"ALTER INMEMORY JOIN GROUP IF EXISTS a REMOVE(a.b(c), d.e(f))",
		"ALTER VIEW a MODIFY PRIMARY KEY RELY",
		"ALTER VIEW a MODIFY UNIQUE (b) NORELY",
		"ALTER VIEW a MODIFY CONSTRAINT RELY RELY",
		"ALTER VIEW a DROP UNIQUE (b, c)",
		"ALTER VIEW a ANNOTATIONS(ADD b, DROP IF EXISTS c 'x', REPLACE d 'y', ADD OR REPLACE e N'z', ADD IF NOT EXISTS f)",
		"ALTER SEQUENCE a RESTART START WITH 1",
		"ALTER SEQUENCE a NOCACHE INCREMENT BY 2 SESSION",
	}
	for _, sql := range accepts {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}

	rejects := []struct {
		sql string
		pos int
	}{
		{"CREATE IF NOT EXISTS CONTEXT a USING a", 7},                             // ORA-00901
		{"CREATE PUBLIC CONTEXT a USING b", 14},                                   // ORA-00901
		{"CREATE CONTEXT a.b USING a", 16},                                        // ORA-02012
		{"CREATE CONTEXT a USING a@b", 24},                                        // ORA-03048
		{"CREATE CONTEXT a USING a SHARING = NONE", 25},                           // ORA-03049
		{"CREATE CONTEXT a USING a INITIALIZED EXTERNALLY ACCESSED GLOBALLY", 48}, // ORA-03049
		{"CREATE OR REPLACE RESTORE POINT a", 18},                                 // ORA-00922
		{"CREATE PUBLIC RESTORE POINT a", 14},                                     // ORA-00901
		{"CREATE RESTORE POINT SELECT", 21},                                       // ORA-03050
		{"CREATE RESTORE POINT a AS OF SCN 1 FOR PLUGGABLE DATABASE b", 35},       // ORA-03048
		{"CREATE RESTORE POINT a GUARANTEE FLASHBACK DATABASE PRESERVE", 52},      // ORA-03049
		{"CREATE IF NOT EXISTS ROLLBACK SEGMENT a", 7},                            // ORA-00901
		{"CREATE OR REPLACE PUBLIC ROLLBACK SEGMENT a", 25},                       // ORA-00922
		{"CREATE ROLLBACK SEGMENT a TABLESPACE b TABLESPACE c", 39},               // ORA-02215
		{"CREATE ROLLBACK SEGMENT a STORAGE (FOO 1)", 35},                         // ORA-02143
		{"ALTER ROLLBACK SEGMENT a STORAGE (NEXT 1M) ONLINE", 43},                 // ORA-03048
		{"ALTER ROLLBACK SEGMENT a SHRINK TO 1.5", 35},                            // ORA-02017
		{"CREATE OR REPLACE INMEMORY JOIN GROUP a(a.b(c))", 18},                   // ORA-00922
		{"CREATE INMEMORY JOIN GROUP a(a.b(c,))", 35},                             // ORA-00900
		{"CREATE INMEMORY JOIN GROUP a@b(a.b(c))", 28},                            // ORA-00906
		{"ALTER INMEMORY JOIN GROUP a ADD(a.b(c)) REMOVE(d.e(f))", 40},            // ORA-00922
		{"ALTER VIEW a@b COMPILE", 12},                                            // ORA-00922
		{"ALTER VIEW a MODIFY CONSTRAINT c", 32},                                  // ORA-00921
		{"ALTER VIEW a MODIFY a RELY", 20},                                        // ORA-00922
		{"ALTER VIEW a DROP PRIMARY KEY CASCADE", 30},                             // ORA-00922
		{"ALTER VIEW a DROP CONSTRAINT SELECT", 29},                               // ORA-02250
		{"ALTER VIEW a ANNOTATIONS()", 25},                                        // ORA-11551
		{"ALTER VIEW a ANNOTATIONS(ADD IF EXISTS b)", 29},                         // ORA-11543
		{"ALTER VIEW a ANNOTATIONS(DROP IF NOT EXISTS b)", 30},                    // ORA-11544
		{"ALTER VIEW a ANNOTATIONS(ADD b 1)", 31},                                 // ORA-11549
		{"ALTER VIEW a ANNOTATIONS(ADD b 'x' 'y')", 35},                           // ORA-11550
		{"ALTER SEQUENCE a", 16},                                                  // ORA-02286
		{"ALTER SEQUENCE a@b CYCLE", 16},                                          // ORA-02286
		{"ALTER SEQUENCE a RESTART WITH 1", 25},                                   // ORA-03048
		{"ALTER SEQUENCE a CYCLE NOCYCLE", 23},                                    // ORA-02280
		{"ALTER SEQUENCE a INCREMENT BY 1 INCREMENT BY 2", 32},                    // ORA-02284
		{"ALTER SEQUENCE a START WITH 1", 17},                                     // ORA-02283
	}
	for _, tt := range rejects {
		_, err := Parse(tt.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) = %v, want a ParseError", tt.sql, err)
			continue
		}
		if pe.Position != tt.pos {
			t.Errorf("Parse(%q) failed at %d, want %d", tt.sql, pe.Position, tt.pos)
		}
	}
}
