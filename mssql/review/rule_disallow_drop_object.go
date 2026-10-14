package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// dropKinds names the object kinds of a DROP statement.
var dropKinds = map[ast.DropObjectType]string{
	ast.DropTable:               "table",
	ast.DropView:                "view",
	ast.DropIndex:               "index",
	ast.DropProcedure:           "procedure",
	ast.DropFunction:            "function",
	ast.DropDatabase:            "database",
	ast.DropSchema:              "schema",
	ast.DropTrigger:             "trigger",
	ast.DropType:                "type",
	ast.DropSequence:            "sequence",
	ast.DropSynonym:             "synonym",
	ast.DropAssembly:            "assembly",
	ast.DropPartitionFunction:   "partition function",
	ast.DropPartitionScheme:     "partition scheme",
	ast.DropStatistics:          "statistics",
	ast.DropDefault:             "default",
	ast.DropRule:                "rule",
	ast.DropXmlSchemaCollection: "XML schema collection",
	ast.DropFulltextIndex:       "fulltext index",
	ast.DropFulltextCatalog:     "fulltext catalog",
	ast.DropMaterializedView:    "materialized view",
}

// checkDisallowDropObject reports every DROP: the DROP statement of any
// object kind, the DROP of a principal, a key, a Service Broker object, a
// security policy, an aggregate, a stoplist, a search property list, or a
// federation, and ALTER TABLE ... DROP COLUMN, which anchors the column.
// IF EXISTS does not exempt a statement, and neither does a routine body
// or a control-of-flow block around it. Dropping a constraint belongs to
// DisallowDropConstraint.
func checkDisallowDropObject(s *statement, r *reporter) {
	report := func(loc ast.Loc, kind string, names ...string) {
		msg := "drops " + kind
		if len(names) > 0 {
			msg += " " + strings.Join(names, ", ")
		}
		r.report(review.DisallowDropObject, s.index, rangeOf(loc), msg)
	}
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.DropStmt:
			kind, ok := dropKinds[v.ObjectType]
			if !ok {
				kind = "object"
			}
			var names []string
			if v.Names != nil {
				for i, item := range v.Names.Items {
					ref, ok := item.(*ast.TableRef)
					if !ok {
						continue
					}
					name := object(ref)
					// An index name is unique only on its table.
					if v.OnTables != nil && i < len(v.OnTables.Items) {
						if on, ok := v.OnTables.Items[i].(*ast.TableRef); ok && on != nil {
							name += " on " + object(on)
						}
					}
					names = append(names, name)
				}
			}
			report(v.Loc, kind, names...)
		case *ast.SecurityStmt:
			if strings.EqualFold(v.Action, "DROP") {
				// A principal, endpoint, audit, or pool has a one-part
				// name, which may hold a dot.
				render := dotted
				if onePartKinds[strings.ToUpper(v.ObjectType)] {
					render = ident
				}
				names := []string{render(v.Name)}
				// DROP EVENT NOTIFICATION n1, n2 keeps the rest apart.
				if v.Options != nil {
					for _, item := range v.Options.Items {
						if o, ok := item.(*ast.EventNotificationOption); ok {
							for _, extra := range o.ExtraNames {
								names = append(names, render(extra))
							}
						}
					}
				}
				report(v.Loc, strings.ToLower(v.ObjectType), names...)
			}
		case *ast.ServiceBrokerStmt:
			if strings.EqualFold(v.Action, "DROP") {
				// Only a queue's name has a schema; the others are one
				// part, often a URI with dots in it.
				name := ident(v.Name)
				if strings.EqualFold(v.ObjectType, "QUEUE") {
					name = dotted(v.Name)
				}
				report(v.Loc, strings.ToLower(v.ObjectType), name)
			}
		case *ast.SecurityKeyStmt:
			if strings.EqualFold(v.Action, "DROP") {
				names := []string{}
				if v.Name != "" {
					names = append(names, ident(v.Name))
				}
				report(v.Loc, strings.ToLower(v.ObjectType), names...)
			}
		case *ast.SecurityPolicyStmt:
			if strings.EqualFold(v.Action, "DROP") {
				report(v.Loc, "security policy", object(v.Name))
			}
		case *ast.SensitivityClassificationStmt:
			if strings.EqualFold(v.Action, "DROP") {
				var names []string
				if v.Columns != nil {
					for _, item := range v.Columns.Items {
						if ref, ok := item.(*ast.TableRef); ok {
							names = append(names, object(ref))
						}
					}
				}
				report(v.Loc, "sensitivity classification of", names...)
			}
		case *ast.SignatureStmt:
			if strings.EqualFold(v.Action, "DROP") {
				kind := "signature of"
				if v.IsCounter {
					kind = "counter signature of"
				}
				report(v.Loc, kind, object(v.ModuleName))
			}
		case *ast.DropStatisticsStmt:
			var names []string
			if v.Names != nil {
				for _, item := range v.Names.Items {
					if name, ok := item.(*ast.String); ok {
						names = append(names, escapeLines(name.Str))
					}
				}
			}
			report(v.Loc, "statistics", names...)
		case *ast.DropAggregateStmt:
			report(v.Loc, "aggregate", object(v.Name))
		case *ast.DropFulltextStoplistStmt:
			report(v.Loc, "fulltext stoplist", ident(v.Name))
		case *ast.DropSearchPropertyListStmt:
			report(v.Loc, "search property list", ident(v.Name))
		case *ast.DropFederationStmt:
			report(v.Loc, "federation", ident(v.Name))
		case *ast.AlterTableStmt:
			if v.Actions == nil {
				return true
			}
			for _, item := range v.Actions.Items {
				if a, ok := item.(*ast.AlterTableAction); ok && a.Type == ast.ATDropColumn {
					report(a.Loc, "column", ident(a.ColName)+" of "+object(v.Name))
				}
			}
		}
		return true
	})
}

// onePartKinds are the security object kinds whose names have no schema
// or database part.
var onePartKinds = map[string]bool{
	"APPLICATION ROLE":             true,
	"AVAILABILITY GROUP":           true,
	"DATABASE AUDIT SPECIFICATION": true,
	"ENDPOINT":                     true,
	"EVENT NOTIFICATION":           true,
	"EVENT SESSION":                true,
	"EXTERNAL DATA SOURCE":         true,
	"EXTERNAL FILE FORMAT":         true,
	"EXTERNAL LANGUAGE":            true,
	"EXTERNAL LIBRARY":             true,
	"EXTERNAL MODEL":               true,
	"EXTERNAL STREAM":              true,
	"EXTERNAL STREAMING JOB":       true,
	"EXTERNAL RESOURCE POOL":       true,
	"LOGIN":                        true,
	"RESOURCE POOL":                true,
	"ROLE":                         true,
	"SERVER AUDIT":                 true,
	"SERVER AUDIT SPECIFICATION":   true,
	"SERVER ROLE":                  true,
	"USER":                         true,
	"WORKLOAD CLASSIFIER":          true,
	"WORKLOAD GROUP":               true,
}
