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
	ast.OBJECT_MATERIALIZED_VIEW_LOG:      "materialized view log",
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
// clause, and ALTER TYPE ... DROP ATTRIBUTE, each attribute anchoring its
// name. Dropping a constraint belongs to DisallowDropConstraint. SET
// UNUSED, which hides a column for a later DROP UNUSED COLUMNS to remove,
// drops nothing itself. A DROP runs inside PL/SQL only as the text of an
// EXECUTE IMMEDIATE, which is not read. Each finding has the code
// review.DropCode gives the kind its message names.
func checkDisallowDropObject(sql string, s *statement, r *reporter) {
	report := func(loc ast.Loc, kind, what string) {
		r.report(review.DropCode(kind), s.index, rangeOf(loc), "Drops "+what+".")
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
		report(v.Loc, objectKinds[v.ObjectType], dropWhat(v.ObjectType, names...))
	case *ast.AdminDDLStmt:
		if strings.EqualFold(v.Action, "DROP") {
			var names []string
			if v.Name != nil {
				names = append(names, objectName(v.Name))
			}
			report(v.Loc, objectKinds[v.ObjectType], dropWhat(v.ObjectType, names...))
		}
	case *ast.DropTablespaceStmt:
		kind := "tablespace"
		if v.IsSet {
			kind = "tablespace set"
		}
		report(v.Loc, kind, kind+" "+review.Code(objectName(v.Name)))
	case *ast.DropAuditPolicyStmt:
		report(v.Loc, "audit policy", "audit policy "+review.Code(ident(v.Name)))
	case *ast.AlterTableStmt:
		if v.Actions == nil {
			return
		}
		table := review.Code(objectName(v.Name))
		for _, item := range v.Actions.Items {
			cmd, ok := item.(*ast.AlterTableCmd)
			if !ok {
				continue
			}
			switch cmd.Action {
			case ast.AT_DROP_COLUMN:
				names := columnList(sql, cmd)
				report(cmd.Loc, "column", review.Plural("column", len(names))+" "+review.Names(names...)+" of "+table)
			case ast.AT_DROP_UNUSED_COLUMNS:
				if strings.EqualFold(cmd.Subtype, "COLUMNS CONTINUE") {
					report(cmd.Loc, "columns", "the columns an interrupted drop left in "+table)
				} else {
					report(cmd.Loc, "unused columns", "the unused columns of "+table)
				}
			case ast.AT_DROP_PARTITION:
				report(cmd.Loc, "partition", partitionWhat("partition", sql, cmd)+" of "+table)
			case ast.AT_DROP_SUBPARTITION:
				report(cmd.Loc, "subpartition", partitionWhat("subpartition", sql, cmd)+" of "+table)
			case ast.AT_DROP_PERIOD:
				report(cmd.Loc, "period", "period "+review.Code(ident(cmd.ColumnName))+" of "+table)
			}
		}
	case *ast.AlterTypeStmt:
		if v.Action != "DROP_ATTRIBUTE" {
			return
		}
		for _, a := range v.Attributes {
			report(a.Loc, "attribute", "attribute "+review.Code(ident(a.Name))+" of type "+review.Code(objectName(v.Name)))
		}
	}
}

// dropWhat writes, as Markdown, what a DROP statement of the kind drops:
// the kind and the names, "the" kind when it names nothing, and the names
// alone when the kind is not known. A materialized view log has no name of
// its own; the DROP names its table.
func dropWhat(t ast.ObjectType, names ...string) string {
	kind, known := objectKinds[t]
	switch {
	case t == ast.OBJECT_MATERIALIZED_VIEW_LOG && len(names) > 0:
		return "the materialized view log on " + review.Names(names...)
	case !known && len(names) > 0:
		return review.Names(names...)
	case !known:
		return "an object"
	case len(names) == 0:
		return "the " + kind
	}
	return review.Plural(kind, len(names)) + " " + review.Names(names...)
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

// partitionWhat writes, as Markdown, the partitions or subpartitions a
// clause names, as "partition `P1`" or "partitions `P1` and `FOR (10)`":
// each name, or the FOR clause that names one by its key, which the parser
// does not keep beyond the first name.
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
	if len(items) == 0 {
		return "a " + kind
	}
	return review.Plural(kind, len(items)) + " " + review.Names(items...)
}
