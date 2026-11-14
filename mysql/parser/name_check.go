package parser

import (
	"strings"
	"unicode/utf8"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// The name checks MySQL's parser runs. A database or table name that is empty, ends in
// a space or runs past 64 characters is 1102, 1103 or 1059 from check_db_name and
// check_table_name, and a routine, trigger or event name is 1458 from
// check_routine_name. The statements below call them from a grammar action, so the
// answer does not change with a syntax error after the statement. SELECT 1 FROM `a `.t
// and CREATE TABLE `a `.t (...) reach the same checks only after parsing, when the tree
// is resolved, and keep reading.
//
// Ref: table_ident, sp_name and the DROP, RENAME, VIEW, TRIGGER and GRANT actions in
// MySQL 8.0 sql/sql_yacc.yy; check_db_name, check_table_name and check_routine_name in
// sql/sql_parse.cc and sql/table.cc.

const maxIdentifierChars = 64

// refuseRepeatedTables is 1066, ER_NONUNIQ_TABLE, for a name a statement's table list
// holds twice: DROP TABLE, DROP VIEW, CHECKSUM TABLE and FLUSH TABLES add each name in
// an action (add_table_to_list), so `DROP TABLE t, t` is refused while it parses.
// CHECK, REPAIR, OPTIMIZE and ANALYZE TABLE add theirs after parsing and keep reading.
// Names compare as written: whether t and T, or t and d.t, are one table is the
// server's lower_case_table_names and current database, which the parser can't see.
func refuseRepeatedTables(refs []*nodes.TableRef) error {
	type key struct{ schema, name string }
	seen := make(map[key]bool, len(refs))
	for _, r := range refs {
		k := key{r.Schema, r.Name}
		if seen[k] {
			return &ParseError{Message: "Not unique table/alias: '" + r.Name + "'", Position: r.Loc.Start}
		}
		seen[k] = true
	}
	return nil
}

func identifierNameRefused(name string) bool {
	return name == "" || strings.HasSuffix(name, " ") || utf8.RuneCountInString(name) > maxIdentifierChars
}

func refuseDatabaseName(name string, pos int) error {
	if identifierNameRefused(name) {
		return &ParseError{Message: "Incorrect database name '" + name + "'", Position: pos}
	}
	return nil
}

func refuseRoutineName(name string, pos int) error {
	if identifierNameRefused(name) {
		return &ParseError{Message: "Incorrect routine name '" + name + "'", Position: pos}
	}
	return nil
}

// parseNamedTableRef is parseTableRef where the grammar action holds the names to the
// checks above.
func (p *Parser) parseNamedTableRef() (*nodes.TableRef, error) {
	start := p.pos()
	qualified := p.peekNext().Type == '.'
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	return ref, refuseTableRefName(ref, qualified, start)
}

// refuseTableRefName holds a table reference to check_db_name and check_table_name.
// qualified says whether it wrote a database, which an empty ref.Schema can't tell.
func refuseTableRefName(ref *nodes.TableRef, qualified bool, pos int) error {
	if qualified {
		if err := refuseDatabaseName(ref.Schema, pos); err != nil {
			return err
		}
	}
	if identifierNameRefused(ref.Name) {
		return &ParseError{Message: "Incorrect table name '" + ref.Name + "'", Position: pos}
	}
	return nil
}

// parseNamedRoutineRef reads sp_name into a TableRef, the way routines hold their name.
func (p *Parser) parseNamedRoutineRef() (*nodes.TableRef, error) {
	start := p.pos()
	qualified := p.peekNext().Type == '.'
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if qualified {
		if err := refuseDatabaseName(ref.Schema, start); err != nil {
			return nil, err
		}
	}
	return ref, refuseRoutineName(ref.Name, start)
}

// parseSpName reads sp_name, [db.]name, for the statements that hold the two apart. A
// trigger and an event take a database as a routine does: `CREATE TRIGGER d.tr` and
// `DROP EVENT d.e` both parse.
func (p *Parser) parseSpName() (schema, name string, err error) {
	ref, err := p.parseNamedRoutineRef()
	if err != nil {
		return "", "", err
	}
	return ref.Schema, ref.Name, nil
}
