package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemCharsetAllowlist judges the character sets top-level CREATE and
// ALTER DATABASE and CREATE and ALTER TABLE statements name against the
// allowed list, compared without regard to letter case. Statements nested in
// a routine body are not judged.
//
// A CREATE DATABASE has only its first CHARACTER SET or CHARSET option
// judged, an ALTER DATABASE every one. A CREATE TABLE has its first
// CHARACTER SET or CHARSET table option judged, then the character set of
// each column's type. An ALTER TABLE has judged the column types its ADD,
// MODIFY and CHANGE COLUMN actions give and each of its table options of
// those names. A CONVERT TO CHARACTER SET is not judged: the V1 advisor
// read its character set from an option the parser does not fill.
//
// Name is the character set lower-cased. Every finding anchors the whole
// statement. Ported from Bytebase's V1 rule SYSTEM_CHARSET_ALLOWLIST.
func SystemCharsetAllowlist(stmts []ast.Node, allowed []string) []v1review.Finding {
	c := &allowlist{code: v1review.CharsetNotAllowed, allowed: make(map[string]bool, len(allowed))}
	for _, charset := range allowed {
		c.allowed[strings.ToLower(charset)] = true
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateDatabaseStmt:
			for _, opt := range n.Options {
				if isCharsetOption(opt.Name, false) {
					c.check(n.Loc, opt.Value)
					break
				}
			}
		case *ast.AlterDatabaseStmt:
			for _, opt := range n.Options {
				if isCharsetOption(opt.Name, false) {
					c.check(n.Loc, opt.Value)
				}
			}
		case *ast.CreateTableStmt:
			for _, opt := range n.Options {
				if isCharsetOption(opt.Name, true) {
					c.check(n.Loc, opt.Value)
					break
				}
			}
			for _, col := range n.Columns {
				if col != nil && col.TypeName != nil {
					c.check(n.Loc, col.TypeName.Charset)
				}
			}
		case *ast.AlterTableStmt:
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
					for _, col := range commandColumns(cmd) {
						if col != nil && col.TypeName != nil {
							c.check(n.Loc, col.TypeName.Charset)
						}
					}
				case ast.ATTableOption:
					if cmd.Option != nil && isCharsetOption(cmd.Option.Name, true) {
						c.check(n.Loc, cmd.Option.Value)
					}
				case ast.ATConvertCharset:
					if cmd.Option != nil && isCharsetOption(cmd.Option.Name, false) {
						c.check(n.Loc, cmd.Option.Value)
					}
				}
			}
		}
	}
	return c.findings
}

// isCharsetOption reports whether an option name is one the V1 advisor read
// a character set from; the DEFAULT CHARSET spelling counted only where a
// table option is read.
func isCharsetOption(name string, tableOption bool) bool {
	return strings.EqualFold(name, "CHARACTER SET") || strings.EqualFold(name, "CHARSET") ||
		tableOption && strings.EqualFold(name, "DEFAULT CHARSET")
}

// allowlist holds a charset or collation allow-list, lower-cased, and the
// findings against it.
type allowlist struct {
	code     v1review.Code
	allowed  map[string]bool
	findings []v1review.Finding
}

// check reports a name that is written and not allowed, anchored at loc.
func (c *allowlist) check(loc ast.Loc, name string) {
	name = strings.ToLower(name)
	if name == "" || c.allowed[name] {
		return
	}
	start, end := span(loc)
	c.findings = append(c.findings, v1review.Finding{Code: c.code, Name: name, Start: start, End: end})
}
