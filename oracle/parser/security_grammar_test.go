package parser

import (
	"errors"
	"testing"
)

// TestSecurityStatementGrammar holds forms of CREATE and ALTER PROFILE,
// CREATE and ALTER AUDIT POLICY, AUDIT and NOAUDIT POLICY, and CREATE and
// ALTER USER that Oracle Free 23.26.3 answered: accepted ones, and refused
// ones at the byte Oracle names.
func TestSecurityStatementGrammar(t *testing.T) {
	accepts := []string{
		"ALTER PROFILE a LIMIT PASSWORD_LIFE_TIME 1/24",
		"ALTER PROFILE a LIMIT PASSWORD_LIFE_TIME (1)",
		"ALTER PROFILE a LIMIT PRIVATE_SGA 1K",
		"ALTER PROFILE a LIMIT PASSWORD_VERIFY_FUNCTION a.b",
		"CREATE MANDATORY PROFILE a LIMIT PASSWORD_VERIFY_FUNCTION a CONTAINER = ALL",
		"CREATE AUDIT POLICY p PRIVILEGES CREATE TABLE, SELECT ANY TABLE",
		"CREATE AUDIT POLICY p ACTIONS SELECT, INSERT ON t, CREATE TABLE",
		"CREATE AUDIT POLICY p ACTIONS ALTER(a) ON t",
		"CREATE AUDIT POLICY p ACTIONS SELECT ON t ACTIONS COMPONENT = DATAPUMP ALL ACTIONS COMPONENT = OLS ALL",
		"CREATE AUDIT POLICY p ACTIONS COMPONENT = DV FACTOR TRUST LEVEL NULL ON f",
		"CREATE AUDIT POLICY p ACTIONS COMPONENT = SQL_FIREWALL ALL ON a, SQL VIOLATION ON b",
		"CREATE AUDIT POLICY p ACTIONS COMPONENT = XS ALL, CREATE USER",
		"CREATE AUDIT POLICY p ROLES CONNECT, RESOURCE",
		"ALTER AUDIT POLICY p ADD PRIVILEGES CREATE TABLE ACTIONS ALL ROLES r",
		"ALTER AUDIT POLICY p DROP ONLY TOPLEVEL",
		"AUDIT POLICY a BY USERS WITH GRANTED ROLES r, s",
		"AUDIT POLICY a BY u WHENEVER NOT SUCCESSFUL",
		"NOAUDIT POLICY a BY USERS WITH GRANTED ROLES r",
		"AUDIT CONTEXT NAMESPACE n ATTRIBUTES a, b BY u, v",
		"CREATE USER a",
		"CREATE USER a IDENTIFIED BY a QUOTA 1 T ON a READ ONLY READ WRITE",
		"CREATE USER a IDENTIFIED BY a ENABLE EDITIONS FOR VIEW",
		"ALTER USER a IDENTIFIED BY 1",
		"ALTER USER a DEFAULT TABLESPACE a TEMPORARY TABLESPACE b QUOTA 1 ON a PROFILE p DEFAULT ROLE ALL PASSWORD EXPIRE ACCOUNT LOCK",
		"ALTER USER a ENABLE EDITIONS FOR VIEW, SYNONYM, SQL TRANSLATION PROFILE FORCE",
		"ALTER USER a SET CONTAINER_DATA = (a, b) FOR a.b",
		"ALTER USER a GRANT CONNECT THROUGH b WITH ROLE ALL EXCEPT c AUTHENTICATION REQUIRED",
		"ALTER USER a GRANT CONNECT THROUGH a WITH ROLES a",
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
		{"ALTER PROFILE a LIMIT SESSIONS_PER_USER 1 LIMIT IDLE_TIME 2", 42},         // ORA-02376
		{"ALTER PROFILE a LIMIT PRIVATE_SGA 1 T", 34},                               // ORA-02377
		{"ALTER PROFILE a LIMIT PASSWORD_VERIFY_FUNCTION SELECT", 47},               // ORA-03050
		{"ALTER PROFILE a CONTAINER = ALL", 16},                                     // ORA-02000
		{"CREATE AUDIT POLICY a ACTIONS a", 30},                                     // ORA-00956
		{"CREATE AUDIT POLICY p ACTIONS SELECT ON t, COMPONENT = DATAPUMP ALL", 43}, // ORA-00956
		{"CREATE AUDIT POLICY p ACTIONS SELECT ON t ACTIONS INSERT ON t", 50},       // ORA-03048
		{"CREATE AUDIT POLICY p ACTIONS COMPONENT = PROTOCOL zzz", 51},              // ORA-00956
		{"CREATE AUDIT POLICY p ROLES SELECT", 28},                                  // ORA-01937
		{"CREATE AUDIT POLICY p ACTIONS ALL WHEN 's'", 42},                          // ORA-02000
		{"ALTER AUDIT a CONDITION DROP", 12},                                        // ORA-00940
		{"ALTER AUDIT POLICY a CONDITION 's' EVALUATE STATEMENT", 44},               // ORA-02000
		{"AUDIT POLICY a BY USERS WITH ROLE a", 29},                                 // ORA-03049
		{"NOAUDIT POLICY a BY a WITH ROLE a", 22},                                   // ORA-03049
		{"AUDIT CONTEXT a ATTRIBUTES a", 14},                                        // ORA-03049
		{"CREATE USER IF NOT a IDENTIFIED EXTERNALLY", 12},                          // ORA-11543
		{"CREATE USER a IDENTIFIED BY a SET CONTAINER_DATA = ALL", 30},              // ORA-00922
		{"ALTER USER a IDENTIFIED BY 'x'", 27},                                      // ORA-00988
		{"ALTER USER a ENABLE EDITIONS FOR VIEW, a", 39},                            // ORA-00922
		{"ALTER USER a SET CONTAINER_DATA ALL", 32},                                 // ORA-02000
		{"ALTER USER a GRANT CONNECT THROUGH b ACCOUNT LOCK", 37},                   // ORA-03049
		{"ALTER USER a REVOKE CONNECT THROUGH b WITH ROLE c", 38},                   // ORA-00922
		{"ALTER USER a GRANT CONNECT THROUGH a WITH ROLE ALL a", 51},                // ORA-28160
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
