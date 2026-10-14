package review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowRename reports, from the SQL alone, ALTER TABLE ... RENAME
// TO, which renames a table, and ALTER TABLE ... RENAME COLUMN, each
// anchoring its clause. Renaming an index, a constraint, a partition, or a
// LOB segment is not in scope.
func checkDisallowRename(s *statement, r *reporter) {
	v, ok := s.node.(*ast.AlterTableStmt)
	if !ok || v.Actions == nil {
		return
	}
	for _, item := range v.Actions.Items {
		cmd, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch {
		case cmd.Action == ast.AT_RENAME && cmd.NewName != "":
			r.report(review.DisallowRename, s.index, rangeOf(cmd.Loc), "renames table "+objectName(v.Name)+" to "+ident(cmd.NewName))
		case cmd.Action == ast.AT_RENAME_COLUMN:
			r.report(review.DisallowRename, s.index, rangeOf(cmd.Loc), "renames column "+ident(cmd.ColumnName)+" of "+objectName(v.Name)+" to "+ident(cmd.NewName))
		}
	}
}

// hasRename reports whether the change holds a RENAME statement, the one
// rename checkRenames classifies.
func hasRename(stmts []statement) bool {
	for i := range stmts {
		if _, ok := stmts[i].node.(*ast.RenameStmt); ok {
			return true
		}
	}
	return false
}

// checkRenames reports a RENAME old TO new that renames a table or a
// view. RENAME also renames a sequence or a private synonym, which is not
// in scope, so what the old name means decides: the object an earlier
// statement of the change made it, or renamed to it, when one did, and
// otherwise the synced schema's. A name the change dropped, made as
// another kind, or made with IF NOT EXISTS, and a name neither the change
// nor the snapshot places, is not judged.
//
// oracle: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/RENAME.html
func checkRenames(stmts []statement, x schemaIndex, r *reporter) {
	schemas := schemasAt(stmts, x.db())
	for pos := range stmts {
		v, ok := stmts[pos].node.(*ast.RenameStmt)
		if !ok || v.OldName == nil || v.NewName == nil {
			continue
		}
		k := keyIn(v.OldName, schemas[pos])
		if kind := relationKindAt(stmts, schemas, pos, k, x); kind != "" {
			r.report(review.DisallowRename, stmts[pos].index, rangeOf(stmts[pos].loc), "renames "+kind+" "+objectName(v.OldName)+" to "+objectName(v.NewName))
		}
	}
}

// match is how a statement's name relates to the name looked for.
type match int

const (
	noMatch match = iota
	maybeMatch
	sureMatch
)

func matchKey(ref *ast.ObjectName, current string, k tableKey) match {
	if ref == nil || ref.DBLink != "" {
		return noMatch
	}
	o := keyIn(ref, current)
	switch {
	case o.is(k):
		return sureMatch
	case o.mayBe(k):
		return maybeMatch
	}
	return noMatch
}

// relationKindAt returns "table" or "view" for what the key names when
// the statement at pos runs, "" for anything else or not known. It reads
// the change backwards for the last statement that made, dropped, or
// renamed an object of the name, and falls back to the synced schema.
func relationKindAt(stmts []statement, schemas []string, pos int, k tableKey, x schemaIndex) string {
	// found decides at a statement naming the key: a sure match gives the
	// kind, a match that is only possible gives nothing certain.
	found := func(m match, kind string) (string, bool) {
		switch m {
		case sureMatch:
			return kind, true
		case maybeMatch:
			return "", true
		}
		return "", false
	}
	for i := pos - 1; i >= 0; i-- {
		cur := schemas[i]
		switch v := stmts[i].node.(type) {
		case *ast.CreateTableStmt:
			kind := "table"
			if v.IfNotExists {
				kind = ""
			}
			if kind, ok := found(matchKey(v.Name, cur, k), kind); ok {
				return kind
			}
		case *ast.CreateSchemaStmt:
			if v.Stmts == nil {
				continue
			}
			for j := len(v.Stmts.Items) - 1; j >= 0; j-- {
				switch e := v.Stmts.Items[j].(type) {
				case *ast.CreateTableStmt:
					if kind, ok := found(matchKey(e.Name, v.SchemaName, k), "table"); ok {
						return kind
					}
				case *ast.CreateViewStmt:
					if kind, ok := found(matchKey(e.Name, v.SchemaName, k), "view"); ok {
						return kind
					}
				}
			}
		case *ast.CreateViewStmt:
			// RENAME refuses a materialized view.
			kind := "view"
			if v.Materialized || v.IfNotExists {
				kind = ""
			}
			if kind, ok := found(matchKey(v.Name, cur, k), kind); ok {
				return kind
			}
		case *ast.CreateJsonDualityViewStmt:
			kind := "view"
			if v.IfNotExists {
				kind = ""
			}
			if kind, ok := found(matchKey(v.Name, cur, k), kind); ok {
				return kind
			}
		case *ast.CreateSequenceStmt:
			if kind, ok := found(matchKey(v.Name, cur, k), ""); ok {
				return kind
			}
		case *ast.CreateSynonymStmt:
			if !v.Public {
				if kind, ok := found(matchKey(v.Name, cur, k), ""); ok {
					return kind
				}
			}
		case *ast.DropStmt:
			// An index, a cluster, a trigger, and the like have names of
			// their own namespaces, which a table's or view's does not share.
			if v.Names == nil || !sharesTableNamespace(v.ObjectType) {
				continue
			}
			for _, item := range v.Names.Items {
				if name, ok := item.(*ast.ObjectName); ok {
					if kind, ok := found(matchKey(name, cur, k), ""); ok {
						return kind
					}
				}
			}
		case *ast.RenameStmt:
			if kind, ok := found(matchKey(v.OldName, cur, k), ""); ok {
				return kind
			}
			switch matchKey(v.NewName, cur, k) {
			case sureMatch:
				return relationKindAt(stmts, schemas, i, keyIn(v.OldName, cur), x)
			case maybeMatch:
				return ""
			}
		case *ast.AlterTableStmt:
			if v.Actions == nil {
				continue
			}
			for _, item := range v.Actions.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok || cmd.Action != ast.AT_RENAME || cmd.NewName == "" {
					continue
				}
				if kind, ok := found(matchKey(v.Name, cur, k), ""); ok {
					return kind
				}
				to := &ast.ObjectName{Schema: v.Name.Schema, Name: cmd.NewName}
				if kind, ok := found(matchKey(to, cur, k), "table"); ok {
					return kind
				}
			}
		case *ast.FlashbackTableStmt:
			if !v.ToBeforeDrop {
				continue
			}
			for _, t := range v.Tables {
				name := t
				if v.Rename != "" {
					name = &ast.ObjectName{Schema: t.Schema, Name: v.Rename}
				}
				if kind, ok := found(matchKey(name, cur, k), "table"); ok {
					return kind
				}
			}
		}
	}
	if !x.own(k) {
		return ""
	}
	return x.relationKind(k)
}

// sharesTableNamespace reports whether a DROP's object kind names objects
// of the namespace tables and views share in a schema.
//
// oracle: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Database-Object-Names-and-Qualifiers.html (Schema Object Namespaces)
func sharesTableNamespace(t ast.ObjectType) bool {
	switch t {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATERIALIZED_VIEW, ast.OBJECT_SEQUENCE,
		ast.OBJECT_SYNONYM, ast.OBJECT_PACKAGE, ast.OBJECT_PROCEDURE, ast.OBJECT_FUNCTION, ast.OBJECT_TYPE:
		return true
	}
	return false
}
