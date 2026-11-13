package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemCharsetAllowlist judges the character sets top-level CREATE and
// ALTER DATABASE and CREATE and ALTER TABLE statements name against the
// allowed list, compared without regard to letter case. Statements nested in
// a routine body are not judged.
//
// A CREATE or ALTER DATABASE has only its first CHARACTER SET or CHARSET
// option judged. A CREATE TABLE reports at most one character set: its first
// CHARACTER SET or CHARSET table option when that is not allowed, otherwise
// the first column whose type names one that is not. An ALTER TABLE reports
// at most one too, the last not allowed among its CHARACTER SET and CHARSET
// table options, its CONVERT TO CHARACTER SET, and the column types of its
// MODIFY and CHANGE COLUMN actions and of the first such column each ADD
// COLUMN gives.
//
// Name is the character set lower-cased. A finding from a CREATE TABLE
// column anchors that column's definition; every other finding anchors the
// whole statement. Ported from Bytebase's V1 rule SYSTEM_CHARSET_ALLOWLIST.
func SystemCharsetAllowlist(stmts []ast.Node, allowed []string) []v1review.Finding {
	return judgeAllowlist(stmts, allowed, allowlistRule{
		code:    v1review.CharsetNotAllowed,
		options: []string{"CHARACTER SET", "CHARSET"},
		column: func(col *ast.ColumnDef) string {
			if col == nil || col.TypeName == nil {
				return ""
			}
			return col.TypeName.Charset
		},
		convert: func(cmd *ast.AlterTableCmd) string { return cmd.Name },
	})
}

// allowlistRule says where a charset or collation allow-list rule reads its
// names: which option names, a column's, and a CONVERT TO's.
type allowlistRule struct {
	code    v1review.Code
	options []string
	column  func(*ast.ColumnDef) string
	convert func(*ast.AlterTableCmd) string
}

// judgeAllowlist judges stmts as SystemCharsetAllowlist and
// SystemCollationAllowlist describe, reading names as rule says.
func judgeAllowlist(stmts []ast.Node, allowed []string, rule allowlistRule) []v1review.Finding {
	allow := make(map[string]bool, len(allowed))
	for _, name := range allowed {
		allow[strings.ToLower(name)] = true
	}
	// notAllowed returns name lower-cased when it is written and not allowed,
	// "" otherwise.
	notAllowed := func(name string) string {
		name = strings.ToLower(name)
		if name == "" || allow[name] {
			return ""
		}
		return name
	}
	var findings []v1review.Finding
	report := func(loc ast.Loc, name string) {
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: rule.code, Name: name, Start: start, End: end})
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateDatabaseStmt:
			if name := notAllowed(firstDatabaseOption(n.Options, rule.options)); name != "" {
				report(n.Loc, name)
			}
		case *ast.AlterDatabaseStmt:
			if name := notAllowed(firstDatabaseOption(n.Options, rule.options)); name != "" {
				report(n.Loc, name)
			}
		case *ast.CreateTableStmt:
			if name := notAllowed(firstTableOption(n.Options, rule.options)); name != "" {
				report(n.Loc, name)
				continue
			}
			for _, col := range n.Columns {
				if col == nil {
					continue
				}
				if name := notAllowed(rule.column(col)); name != "" {
					report(col.Loc, name)
					break
				}
			}
		case *ast.AlterTableStmt:
			last := ""
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				name := ""
				switch cmd.Type {
				case ast.ATTableOption:
					if cmd.Option != nil && isOptionNamed(cmd.Option.Name, rule.options) {
						name = notAllowed(cmd.Option.Value)
					}
				case ast.ATConvertCharset:
					name = notAllowed(rule.convert(cmd))
				case ast.ATAddColumn:
					for _, col := range addedColumns(cmd) {
						if col == nil {
							continue
						}
						if name = notAllowed(rule.column(col)); name != "" {
							break
						}
					}
				case ast.ATChangeColumn, ast.ATModifyColumn:
					if cmd.Column != nil {
						name = notAllowed(rule.column(cmd.Column))
					}
				}
				if name != "" {
					last = name
				}
			}
			if last != "" {
				report(n.Loc, last)
			}
		}
	}
	return findings
}

// firstDatabaseOption returns the value of the first option with one of the
// names, "" when there is none.
func firstDatabaseOption(opts []*ast.DatabaseOption, names []string) string {
	for _, opt := range opts {
		if opt != nil && isOptionNamed(opt.Name, names) {
			return opt.Value
		}
	}
	return ""
}

// firstTableOption returns the value of the first option with one of the
// names, "" when there is none.
func firstTableOption(opts []*ast.TableOption, names []string) string {
	for _, opt := range opts {
		if opt != nil && isOptionNamed(opt.Name, names) {
			return opt.Value
		}
	}
	return ""
}

// isOptionNamed reports whether an option name is one of names, compared
// without regard to letter case.
func isOptionNamed(name string, names []string) bool {
	for _, n := range names {
		if strings.EqualFold(name, n) {
			return true
		}
	}
	return false
}
