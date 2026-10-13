package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// renamedObject is an sp_rename of an object the SQL alone cannot
// classify: it renames a table or a view only if the object is one, which
// the rule looks up in the target's synced schema.
type renamedObject struct {
	// statement is the index in Result.Statements, and position the
	// statement's place in the parsed list.
	statement int
	position  int
	loc       ast.Loc
	name      *ast.TableRef
	newName   string
}

// checkDisallowRename reports the rename of a table or view, by RENAME
// OBJECT, and the rename of a column, by RENAME OBJECT ... COLUMN or
// sp_rename with the COLUMN type, at top level or inside a routine body.
// An sp_rename of an object, with no type or OBJECT, is returned, for the
// target's synced schema to tell whether the object is a table or a view.
// Renaming an index, statistics, a database, or a type is not in scope.
//
// sql: https://learn.microsoft.com/sql/relational-databases/system-stored-procedures/sp-rename-transact-sql
func checkDisallowRename(pos int, s *statement, r *reporter) []renamedObject {
	var open []renamedObject
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.RenameStmt:
			switch {
			case v.ColumnName != "":
				r.report(review.DisallowRename, s.index, rangeOf(v.Loc), "renames column "+ident(v.ColumnName)+" of "+object(v.Name)+" to "+ident(v.NewColumnName))
			case !strings.EqualFold(v.ObjectType, "DATABASE"):
				r.report(review.DisallowRename, s.index, rangeOf(v.Loc), "renames table "+object(v.Name)+" to "+ident(v.NewName))
			}
		case *ast.ExecStmt:
			name, newName, kind, ok := spRename(v)
			if !ok {
				return true
			}
			switch strings.ToUpper(kind) {
			case "COLUMN":
				if ref, column := splitColumn(name); ref != nil {
					r.report(review.DisallowRename, s.index, rangeOf(v.Loc), "renames column "+ident(column)+" of "+object(ref)+" to "+ident(newName))
				}
			case "", "OBJECT":
				// sp_rename called through another database renames an
				// object there, which the synced schema does not hold.
				if v.Name.Database != "" || v.Name.Server != "" {
					break
				}
				if ref := parseObjectName(name); ref != nil {
					open = append(open, renamedObject{s.index, pos, v.Loc, ref, newName})
				}
			}
		}
		return true
	})
	return open
}

// spRename reads a call of sp_rename: the object name, the new name, and
// the object type, from positional or named arguments; ok is false for
// another call, or one whose names are not string literals.
func spRename(v *ast.ExecStmt) (name, newName, kind string, ok bool) {
	if v.Name == nil || !sameName(v.Name.Object, "sp_rename") || v.Name.Schema != "" && !sameName(v.Name.Schema, "sys") || v.Args == nil {
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

// parseObjectName reads a multipart object name as sp_rename takes it:
// parts separated by dots, each bare or in brackets or double quotes. It
// returns nil for a name it cannot read.
func parseObjectName(s string) *ast.TableRef {
	parts := splitParts(s)
	ref := &ast.TableRef{}
	switch len(parts) {
	case 1:
		ref.Object = parts[0]
	case 2:
		ref.Schema, ref.Object = parts[0], parts[1]
	case 3:
		ref.Database, ref.Schema, ref.Object = parts[0], parts[1], parts[2]
	default:
		return nil
	}
	return ref
}

// splitColumn reads table.column, or a longer name ending in the column.
func splitColumn(s string) (*ast.TableRef, string) {
	parts := splitParts(s)
	if len(parts) < 2 || len(parts) > 4 {
		return nil, ""
	}
	ref := parseObjectName(strings.Join(quoteParts(parts[:len(parts)-1]), "."))
	return ref, parts[len(parts)-1]
}

func quoteParts(parts []string) []string {
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = "[" + strings.ReplaceAll(p, "]", "]]") + "]"
	}
	return out
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
		}
		p := strings.TrimSpace(part.String())
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
