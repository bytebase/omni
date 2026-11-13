package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequired judges each top-level statement on its own for required
// columns it leaves a table without. A CREATE TABLE with a name lacks every
// entry of required, duplicates once, that no column it defines matches,
// both lower-cased; its findings name the entries as given, sorted
// bytewise. An ALTER TABLE with a name lacks each column a DROP COLUMN
// action drops, lower-cased, that is exactly an entry, so an entry with an
// upper-case letter is never dropped; its findings name the dropped column
// lower-cased, in action order and then drop-list order, once per drop.
// What an action adds or renames is not judged, nor are statements nested in
// a routine body.
//
// Each finding has code ColumnMissing, Name the column, Table the table as
// written, database.schema.name with absent parts left out, and anchors the
// statement. Ported from Bytebase's V1 rule COLUMN_REQUIRED.
func ColumnRequired(stmts []ast.Node, required []string) []v1review.Finding {
	need := make(map[string]bool, len(required))
	for _, column := range required {
		need[column] = true
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Name == nil {
				continue
			}
			present := make(map[string]bool)
			for _, col := range listOf[*ast.ColumnDef](n.Columns) {
				present[strings.ToLower(col.Name)] = true
			}
			var missing []string
			for column := range need {
				if !present[strings.ToLower(column)] {
					missing = append(missing, column)
				}
			}
			slices.Sort(missing)
			start, end := span(n.Loc)
			table := tableRefText(n.Name)
			for _, column := range missing {
				findings = append(findings, v1review.Finding{Code: v1review.ColumnMissing, Name: column, Table: table, Start: start, End: end})
			}
		case *ast.AlterTableStmt:
			if n.Name == nil || n.Actions == nil {
				continue
			}
			start, end := span(n.Loc)
			table := tableRefText(n.Name)
			for _, action := range listOf[*ast.AlterTableAction](n.Actions) {
				if action.Type != ast.ATDropColumn {
					continue
				}
				var dropped []string
				if action.ColName != "" {
					dropped = append(dropped, strings.ToLower(action.ColName))
				}
				for _, ref := range listOf[*ast.TableRef](action.Names) {
					if ref.Object != "" {
						dropped = append(dropped, strings.ToLower(ref.Object))
					}
				}
				for _, column := range dropped {
					if need[column] {
						findings = append(findings, v1review.Finding{Code: v1review.ColumnMissing, Name: column, Table: table, Start: start, End: end})
					}
				}
			}
		}
	}
	return findings
}

// tableRefText is a table reference as written, database.schema.name with
// absent parts left out.
func tableRefText(ref *ast.TableRef) string {
	result := ref.Object
	if ref.Schema != "" {
		result = ref.Schema + "." + result
	}
	if ref.Database != "" {
		result = ref.Database + "." + result
	}
	return result
}

// listOf returns the items of list that are of type T, nil for a nil list.
func listOf[T ast.Node](list *ast.List) []T {
	if list == nil {
		return nil
	}
	var result []T
	for _, item := range list.Items {
		if v, ok := item.(T); ok {
			result = append(result, v)
		}
	}
	return result
}
