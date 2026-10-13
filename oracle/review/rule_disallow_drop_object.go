package review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
	"github.com/bytebase/omni/review"
)

// objectKinds names the object kinds of a DROP statement.
var objectKinds = map[ast.ObjectType]string{
	ast.OBJECT_TABLE:                      "table",
	ast.OBJECT_INDEX:                      "index",
	ast.OBJECT_VIEW:                       "view",
	ast.OBJECT_MATERIALIZED_VIEW:          "materialized view",
	ast.OBJECT_SEQUENCE:                   "sequence",
	ast.OBJECT_SYNONYM:                    "synonym",
	ast.OBJECT_PACKAGE:                    "package",
	ast.OBJECT_PACKAGE_BODY:               "package body",
	ast.OBJECT_PROCEDURE:                  "procedure",
	ast.OBJECT_FUNCTION:                   "function",
	ast.OBJECT_TRIGGER:                    "trigger",
	ast.OBJECT_TYPE:                       "type",
	ast.OBJECT_TYPE_BODY:                  "type body",
	ast.OBJECT_DATABASE_LINK:              "database link",
	ast.OBJECT_USER:                       "user",
	ast.OBJECT_ROLE:                       "role",
	ast.OBJECT_TABLESPACE:                 "tablespace",
	ast.OBJECT_PROFILE:                    "profile",
	ast.OBJECT_DIRECTORY:                  "directory",
	ast.OBJECT_CONTEXT:                    "context",
	ast.OBJECT_CLUSTER:                    "cluster",
	ast.OBJECT_DIMENSION:                  "dimension",
	ast.OBJECT_FLASHBACK_ARCHIVE:          "flashback archive",
	ast.OBJECT_JAVA:                       "Java object",
	ast.OBJECT_LIBRARY:                    "library",
	ast.OBJECT_MATERIALIZED_VIEW_LOG:      "materialized view log on",
	ast.OBJECT_DATABASE:                   "database",
	ast.OBJECT_DISKGROUP:                  "disk group",
	ast.OBJECT_PLUGGABLE_DATABASE:         "pluggable database",
	ast.OBJECT_AUDIT_POLICY:               "audit policy",
	ast.OBJECT_ANALYTIC_VIEW:              "analytic view",
	ast.OBJECT_ATTRIBUTE_DIMENSION:        "attribute dimension",
	ast.OBJECT_HIERARCHY:                  "hierarchy",
	ast.OBJECT_DOMAIN:                     "domain",
	ast.OBJECT_INDEXTYPE:                  "indextype",
	ast.OBJECT_OPERATOR:                   "operator",
	ast.OBJECT_LOCKDOWN_PROFILE:           "lockdown profile",
	ast.OBJECT_OUTLINE:                    "outline",
	ast.OBJECT_MATERIALIZED_ZONEMAP:       "materialized zonemap",
	ast.OBJECT_INMEMORY_JOIN_GROUP:        "inmemory join group",
	ast.OBJECT_JSON_DUALITY_VIEW:          "JSON relational duality view",
	ast.OBJECT_ROLLBACK_SEGMENT:           "rollback segment",
	ast.OBJECT_EDITION:                    "edition",
	ast.OBJECT_TABLESPACE_SET:             "tablespace set",
	ast.OBJECT_MLE_ENV:                    "MLE environment",
	ast.OBJECT_MLE_MODULE:                 "MLE module",
	ast.OBJECT_PROPERTY_GRAPH:             "property graph",
	ast.OBJECT_VECTOR_INDEX:               "vector index",
	ast.OBJECT_RESTORE_POINT:              "restore point",
	ast.OBJECT_LOGICAL_PARTITION_TRACKING: "logical partition tracking",
	ast.OBJECT_PMEM_FILESTORE:             "PMEM filestore",
}

// checkDisallowDropObject reports every DROP: the DROP statement of any
// object kind, IF EXISTS or not, and the ALTER TABLE clauses that drop a
// part of a table: DROP COLUMN, DROP UNUSED COLUMNS and COLUMNS CONTINUE,
// DROP PARTITION, DROP SUBPARTITION, and DROP PERIOD, each anchoring its
// clause. Dropping a constraint belongs to DisallowDropConstraint. SET
// UNUSED, which hides a column for a later DROP UNUSED COLUMNS to remove,
// drops nothing itself. A DROP runs inside PL/SQL only as the text of an
// EXECUTE IMMEDIATE, which is not read.
func checkDisallowDropObject(sql string, s *statement, r *reporter) {
	report := func(loc ast.Loc, what string) {
		r.report(review.DisallowDropObject, s.index, rangeOf(loc), "drops "+what)
	}
	kind := func(t ast.ObjectType) string {
		if k, ok := objectKinds[t]; ok {
			return k
		}
		return "object"
	}
	switch v := s.node.(type) {
	case *ast.DropStmt:
		var names []string
		if v.Names != nil {
			for _, item := range v.Names.Items {
				if name, ok := item.(*ast.ObjectName); ok {
					names = append(names, objectName(name))
				}
			}
		}
		what := kind(v.ObjectType)
		if len(names) > 0 {
			what += " " + strings.Join(names, ", ")
		}
		report(v.Loc, what)
	case *ast.AdminDDLStmt:
		if strings.EqualFold(v.Action, "DROP") {
			what := kind(v.ObjectType)
			if v.Name != nil {
				what += " " + objectName(v.Name)
			}
			report(v.Loc, what)
		}
	case *ast.DropTablespaceStmt:
		what := "tablespace "
		if v.IsSet {
			what = "tablespace set "
		}
		report(v.Loc, what+objectName(v.Name))
	case *ast.DropAuditPolicyStmt:
		report(v.Loc, "audit policy "+ident(v.Name))
	case *ast.AlterTableStmt:
		if v.Actions == nil {
			return
		}
		table := objectName(v.Name)
		for _, item := range v.Actions.Items {
			cmd, ok := item.(*ast.AlterTableCmd)
			if !ok {
				continue
			}
			switch cmd.Action {
			case ast.AT_DROP_COLUMN:
				names := columnList(sql, cmd)
				if len(names) == 1 {
					report(cmd.Loc, "column "+names[0]+" of "+table)
				} else {
					report(cmd.Loc, "columns "+strings.Join(names, ", ")+" of "+table)
				}
			case ast.AT_DROP_UNUSED_COLUMNS:
				if strings.EqualFold(cmd.Subtype, "COLUMNS CONTINUE") {
					report(cmd.Loc, "the columns an interrupted drop left of "+table)
				} else {
					report(cmd.Loc, "the unused columns of "+table)
				}
			case ast.AT_DROP_PARTITION:
				report(cmd.Loc, partitionWhat("partition", sql, cmd)+" of "+table)
			case ast.AT_DROP_SUBPARTITION:
				report(cmd.Loc, partitionWhat("subpartition", sql, cmd)+" of "+table)
			case ast.AT_DROP_PERIOD:
				report(cmd.Loc, "period "+ident(cmd.ColumnName)+" of "+table)
			}
		}
	}
}

// columnList returns the columns a DROP COLUMN clause names: the one of
// DROP COLUMN, or each in the parentheses of DROP (a, b), which the parser
// does not keep beyond the first.
func columnList(sql string, cmd *ast.AlterTableCmd) []string {
	toks := tokens(sql, cmd.Loc)
	if len(toks) < 2 || toks[1].Type != '(' {
		return []string{ident(cmd.ColumnName)}
	}
	var names []string
	for _, tok := range toks[2:] {
		if tok.Type == ')' {
			break
		}
		if tok.Type != ',' {
			names = append(names, ident(tok.Str))
		}
	}
	if len(names) == 0 {
		return []string{ident(cmd.ColumnName)}
	}
	return names
}

var keywordFor = parser.Tokenize("FOR")[0].Type

// partitionWhat writes the partitions or subpartitions a clause names, as
// "partition P1" or "partitions P1, FOR (10)": each name, or the FOR
// clause that names one by its key, which the parser does not keep beyond
// the first name.
func partitionWhat(kind, sql string, cmd *ast.AlterTableCmd) string {
	toks := tokens(sql, cmd.Loc)
	// The clause starts DROP or TRUNCATE, then the kind's keyword.
	var items []string
	i := 2
	for i < len(toks) {
		tok := toks[i]
		switch {
		case tok.Type == keywordFor && i+1 < len(toks) && toks[i+1].Type == '(':
			depth, j := 0, i+1
			for ; j < len(toks); j++ {
				if toks[j].Type == '(' {
					depth++
				} else if toks[j].Type == ')' {
					if depth--; depth == 0 {
						break
					}
				}
			}
			if j >= len(toks) {
				j = len(toks) - 1
			}
			items = append(items, "FOR "+strings.Join(strings.Fields(sql[toks[i+1].Loc:toks[j].End]), " "))
			i = j + 1
		case tok.Type == '(' || tok.Type == ')' || tok.Type == ',' || tok.Type == ';':
			i = len(toks)
		default:
			items = append(items, ident(tok.Str))
			i++
		}
		if i >= len(toks) || toks[i].Type != ',' {
			break
		}
		i++
	}
	switch len(items) {
	case 0:
		return "a " + kind
	case 1:
		return kind + " " + items[0]
	}
	return kind + "s " + strings.Join(items, ", ")
}
