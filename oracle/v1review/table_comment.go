package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// TableComment judges, after reading every statement, the comment each
// top-level CREATE TABLE's table is left with by top-level COMMENT ON TABLE
// statements, the last COMMENT for a table winning wherever it stands.
// COMMENT ON COLUMN counts as a comment on the column's table, since the
// AST gives it the table object type. A table is keyed by its qualified
// name (see qualifiedTable), so a COMMENT must name it as the CREATE did
// once currentSchema is applied. With no comment or an empty one, the table
// is TableCommentMissing if cfg.Required; otherwise a comment longer than a
// positive cfg.MaxLength, in bytes, is TableCommentTooLong. Name is the
// qualified name.
//
// There is one finding per CREATE TABLE, in statement order, so a table
// created twice is judged twice. Every finding, too-long ones included,
// anchors the last CREATE TABLE of its table, and Statement is that
// statement's index in stmts. A CREATE TABLE nested in CREATE SCHEMA or a
// PL/SQL block is not judged. A statement may be passed bare or wrapped in
// an *ast.RawStmt. Ported from Bytebase's V1 rule TABLE_COMMENT.
//
// A nil cfg panics at the first created table, as the V1 advisor did when
// the rule had no payload.
func TableComment(stmts []ast.Node, currentSchema string, cfg *v1review.TableCommentConfig) []v1review.Finding {
	var created []string
	anchors := make(map[string]tableAnchor)
	comments := make(map[string]string)
	for i, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			name := qualifiedTable(n.Name, currentSchema)
			created = append(created, name)
			anchors[name] = tableAnchor{i, n.Loc}
		case *ast.CommentStmt:
			if n.ObjectType == ast.OBJECT_TABLE {
				comments[qualifiedTable(n.Object, currentSchema)] = n.Comment
			}
		}
	}
	var findings []v1review.Finding
	for _, name := range created {
		comment, ok := comments[name]
		if !ok || comment == "" {
			if cfg.Required {
				findings = append(findings, anchors[name].finding(v1review.TableCommentMissing, name))
			}
		} else if cfg.MaxLength > 0 && len(comment) > cfg.MaxLength {
			findings = append(findings, anchors[name].finding(v1review.TableCommentTooLong, name))
		}
	}
	return findings
}

// tableAnchor is the statement, and the node in it, a table rule that
// judges statements together took a table's line from.
type tableAnchor struct {
	statement int
	loc       ast.Loc
}

func (a tableAnchor) finding(code v1review.Code, name string) v1review.Finding {
	start, end := span(a.loc)
	return v1review.Finding{Code: code, Name: name, Start: start, End: end, Statement: a.statement}
}

// qualifiedTable returns a table's name as the V1 table rules keyed and
// printed it: schema.name, the schema as written or else currentSchema,
// each part as the parser stores it; the name alone when both schemas are
// empty; "" when the AST names no table.
func qualifiedTable(name *ast.ObjectName, currentSchema string) string {
	if name == nil {
		return ""
	}
	schema := name.Schema
	if schema == "" {
		schema = currentSchema
	}
	if schema == "" {
		return name.Name
	}
	return schema + "." + name.Name
}
