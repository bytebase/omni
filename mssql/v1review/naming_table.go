// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for SQL Server. The contract, and the split between
// these functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTable judges every table name the statements introduce against
// cfg: the table a CREATE TABLE (plain, CTAS, or AS CLONE OF) or a
// SELECT ... INTO creates, and the new name a RENAME OBJECT or an
// sp_rename gives, at top level or inside a routine body. The name is
// judged as written, without its database or schema qualifier, a
// temporary table's # prefix included; CREATE EXTERNAL and REMOTE TABLE
// are not judged. The rule reads an sp_rename with no type or OBJECT, not
// through another database, the old name readable; the SQL alone cannot
// tell such a rename's table from a view, so the new name is judged
// either way. NameMismatch and NameTooLong are the codes. Ported from
// Bytebase's V1 rule NAMING_TABLE.
func NamingTable(stmts []ast.Node, cfg v1review.NamingConfig) ([]v1review.Finding, error) {
	pattern, err := cfg.Compile()
	if err != nil {
		return nil, err
	}
	if pattern == nil && cfg.MaxLength <= 0 {
		return nil, nil
	}
	naming := &tableNaming{pattern: pattern, max: cfg.MaxLength}
	for _, stmt := range stmts {
		if stmt == nil {
			continue
		}
		ast.Inspect(stmt, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.CreateTableStmt:
				if v.Name != nil {
					naming.check(v.Name.Loc, v.Name.Object)
				}
			case *ast.CreateTableAsSelectStmt:
				if v.Name != nil {
					naming.check(v.Name.Loc, v.Name.Object)
				}
			case *ast.CreateTableCloneStmt:
				if v.Name != nil {
					naming.check(v.Name.Loc, v.Name.Object)
				}
			case *ast.SelectStmt:
				if v.IntoTable != nil {
					naming.check(v.IntoTable.Loc, v.IntoTable.Object)
				}
			case *ast.RenameStmt:
				if v.ColumnName == "" && !strings.EqualFold(v.ObjectType, "DATABASE") {
					naming.check(v.Loc, v.NewName)
				}
			case *ast.ExecStmt:
				name, newName, kind, ok := spRename(v)
				if !ok {
					return true
				}
				switch strings.ToUpper(kind) {
				case "", "OBJECT":
					// sp_rename called through another database renames
					// an object there.
					if v.Name.Database != "" || v.Name.Server != "" {
						return true
					}
					if readableObjectName(name) {
						naming.check(v.Loc, newName)
					}
				}
			}
			return true
		})
	}
	return naming.findings, nil
}

// tableNaming applies NamingTable's compiled config.
type tableNaming struct {
	pattern  *regexp.Regexp // nil requires nothing
	max      int            // 0 means no limit
	findings []v1review.Finding
}

func (c *tableNaming) check(loc ast.Loc, name string) {
	start, end := span(loc)
	if c.pattern != nil && !c.pattern.MatchString(name) {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Start: start, End: end})
	}
	if c.max > 0 && utf8.RuneCountInString(name) > c.max {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameTooLong, Name: name, Start: start, End: end})
	}
}

// span returns a Loc as a finding's byte range, 0,0 when the AST kept none.
func span(loc ast.Loc) (int, int) {
	if loc.Start < 0 || loc.End <= loc.Start {
		return 0, 0
	}
	return loc.Start, loc.End
}

// spRename reads a call of sp_rename: the object name, the new name, and
// the object type, from positional or named arguments; ok is false for
// another call, or one whose names are not string literals. Inlined from
// the V2 mssql review package's helper of the same name, since V1 and V2
// share no code.
//
// sql: https://learn.microsoft.com/sql/relational-databases/system-stored-procedures/sp-rename-transact-sql
func spRename(v *ast.ExecStmt) (name, newName, kind string, ok bool) {
	if v.Name == nil || !strings.EqualFold(v.Name.Object, "sp_rename") || v.Name.Schema != "" && !strings.EqualFold(v.Name.Schema, "sys") || v.Args == nil {
		return "", "", "", false
	}
	params := []string{"@objname", "@newname", "@objtype"}
	values := make(map[string]string)
	for i, item := range v.Args.Items {
		arg, isArg := item.(*ast.ExecArg)
		if !isArg {
			return "", "", "", false
		}
		key := strings.ToLower(arg.Name)
		if key == "" {
			if i >= len(params) {
				return "", "", "", false
			}
			key = params[i]
		}
		lit, isLit := arg.Value.(*ast.Literal)
		if !isLit || lit.Type != ast.LitString {
			return "", "", "", false
		}
		values[key] = lit.Str
	}
	name, newName = values["@objname"], values["@newname"]
	if name == "" || newName == "" {
		return "", "", "", false
	}
	return name, newName, values["@objtype"], true
}

// readableObjectName reports whether sp_rename's object name is one the
// rule can read: one to three dot-separated parts, each bare or in
// brackets or double quotes.
func readableObjectName(s string) bool {
	parts := splitParts(s)
	return len(parts) >= 1 && len(parts) <= 3
}

// splitParts splits a dotted name, honoring brackets and double quotes;
// nil when a part is empty or unterminated.
func splitParts(s string) []string {
	var parts []string
	i := 0
	for {
		var part strings.Builder
		switch {
		case i < len(s) && s[i] == '[':
			i++
			for {
				if i >= len(s) {
					return nil
				}
				if s[i] == ']' {
					if i+1 < len(s) && s[i+1] == ']' {
						part.WriteByte(']')
						i += 2
						continue
					}
					i++
					break
				}
				part.WriteByte(s[i])
				i++
			}
		case i < len(s) && s[i] == '"':
			i++
			for {
				if i >= len(s) {
					return nil
				}
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						part.WriteByte('"')
						i += 2
						continue
					}
					i++
					break
				}
				part.WriteByte(s[i])
				i++
			}
		default:
			for i < len(s) && s[i] != '.' {
				part.WriteByte(s[i])
				i++
			}
			// Spaces around a bare name are not part of it; inside
			// delimiters they are.
			bare := strings.TrimSpace(part.String())
			part.Reset()
			part.WriteString(bare)
		}
		p := part.String()
		if p == "" {
			return nil
		}
		parts = append(parts, p)
		if i >= len(s) {
			return parts
		}
		if s[i] != '.' {
			return nil
		}
		i++
	}
}
