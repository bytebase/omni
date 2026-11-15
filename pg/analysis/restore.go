package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
)

// RestoreOptions configures RestoreSQL.
type RestoreOptions struct {
	// OriginalSchema and OriginalTable name the table the backup came from.
	OriginalSchema string
	OriginalTable  string
	// BackupSchema and BackupTable name the backup table.
	BackupSchema string
	BackupTable  string
	// Table is the original table's metadata, and PrimaryKey its primary key
	// or nil.
	Table      *metadata.TableMetadata
	PrimaryKey *metadata.IndexMetadata
	// OriginalSQL is the backed-up statement as the restore's comment quotes it.
	OriginalSQL string
	// PrependStatements, when set, runs before the restore.
	PrependStatements string
}

// BackedUpDML is an UPDATE or DELETE of a script, with its node's byte range.
type BackedUpDML struct {
	Node  ast.Node
	Text  string
	Start int
	End   int
}

// BackedUpDMLs returns the UPDATE and DELETE statements of a parsed script,
// the statements a prior backup covers, in order.
func BackedUpDMLs(stmts []omnipg.Statement) []BackedUpDML {
	var dmls []BackedUpDML
	for _, stmt := range stmts {
		if stmt.Empty() {
			continue
		}
		switch stmt.AST.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			nodeLoc := ast.NodeLoc(stmt.AST)
			dmls = append(dmls, BackedUpDML{Node: stmt.AST, Text: stmt.Text, Start: nodeLoc.Start, End: nodeLoc.End})
		default:
		}
	}
	return dmls
}

// RestoreSQL returns the statement that restores, from its backup table, the
// rows node changed: an INSERT of the deleted rows for a DELETE, and for an
// UPDATE an upsert on a unique key the UPDATE leaves alone. A comment quotes
// the original statement. An UPDATE whose every unique key it changes is an
// error.
func RestoreSQL(node ast.Node, opts RestoreOptions) (string, error) {
	sqlForComment := opts.OriginalSQL
	prependStatements := opts.PrependStatements
	backupSchema := opts.BackupSchema
	backupTable := opts.BackupTable
	originalSchema := opts.OriginalSchema
	originalTable := opts.OriginalTable
	quotedColumnList := quotePGColumns(restorableColumns(opts.Table))

	var result string
	switch n := node.(type) {
	case *ast.DeleteStmt:
		result = fmt.Sprintf(`INSERT INTO %s.%s (%s) SELECT %s FROM %s.%s;`, quoteIdentifier(originalSchema), quoteIdentifier(originalTable), quotedColumnList, quotedColumnList, quoteIdentifier(backupSchema), quoteIdentifier(backupTable))
	case *ast.UpdateStmt:
		fields := extractSetFieldNames(n)
		conflictTarget, err := findDisjointUniqueKey(opts.Table, opts.PrimaryKey, fields)
		if err != nil {
			return "", err
		}

		var buf strings.Builder
		if _, err := fmt.Fprintf(&buf, `INSERT INTO %s.%s (%s) SELECT %s FROM %s.%s ON CONFLICT %s DO UPDATE SET `, quoteIdentifier(originalSchema), quoteIdentifier(originalTable), quotedColumnList, quotedColumnList, quoteIdentifier(backupSchema), quoteIdentifier(backupTable), conflictTarget); err != nil {
			return "", fmt.Errorf("failed to generate update statement: %w", err)
		}
		for i, field := range fields {
			if i > 0 {
				if _, err := fmt.Fprint(&buf, ", "); err != nil {
					return "", fmt.Errorf("failed to generate update statement: %w", err)
				}
			}
			if _, err := fmt.Fprintf(&buf, `%s = EXCLUDED.%s`, quoteIdentifier(field), quoteIdentifier(field)); err != nil {
				return "", fmt.Errorf("failed to generate update statement: %w", err)
			}
		}
		if _, err := fmt.Fprint(&buf, `;`); err != nil {
			return "", fmt.Errorf("failed to generate update statement: %w", err)
		}
		result = buf.String()
	default:
		return "", fmt.Errorf("unexpected statement type: %T", node)
	}

	sqlForComment = commentText(sqlForComment)
	if len(prependStatements) > 0 {
		// Ensure prependStatements ends with semicolon
		if !strings.HasSuffix(prependStatements, ";") {
			prependStatements += ";"
		}
		return fmt.Sprintf("%s\n/*\nOriginal SQL:\n%s\n*/\n%s", prependStatements, sqlForComment, result), nil
	}
	return fmt.Sprintf("/*\nOriginal SQL:\n%s\n*/\n%s", sqlForComment, result), nil
}

func restorableColumns(table *metadata.TableMetadata) []string {
	var columns []string
	for _, column := range table.GetColumns() {
		if column.GetGeneration() != nil {
			continue
		}
		columns = append(columns, column.Name)
	}
	return columns
}

func quotePGColumns(columns []string) string {
	var quotedColumns []string
	for _, column := range columns {
		quotedColumns = append(quotedColumns, quoteIdentifier(column))
	}
	return strings.Join(quotedColumns, ", ")
}

// extractSetFieldNames extracts the target column names from an UPDATE SET clause.
func extractSetFieldNames(stmt *ast.UpdateStmt) []string {
	if stmt.TargetList == nil {
		return nil
	}
	var fields []string
	for _, item := range stmt.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Name == "" {
			continue
		}
		// Take the last name component (the actual column name).
		// For "SET test.c1 = 1", Name="test", Indirection=["c1"] → want "c1".
		// For "SET c1 = 1", Name="c1", Indirection=nil → want "c1".
		// For "SET schema.col[1] = val", Indirection=["col", integer] → want "col".
		if rt.Indirection != nil && rt.Indirection.Len() > 0 {
			found := false
			for j := rt.Indirection.Len() - 1; j >= 0; j-- {
				if s, ok := rt.Indirection.Items[j].(*ast.String); ok {
					fields = append(fields, s.Str)
					found = true
					break
				}
			}
			if found {
				continue
			}
		}
		fields = append(fields, rt.Name)
	}
	return fields
}

func disjoint(a []string, b map[string]bool) bool {
	for _, item := range a {
		if _, ok := b[item]; ok {
			return false
		}
	}
	return true
}

// findDisjointUniqueKey returns the ON CONFLICT target of the first unique
// key whose columns the UPDATE leaves alone: the primary key or a unique
// constraint by name, or a unique index, which is no constraint, by its
// columns. A unique index on an expression or with a predicate is passed
// over: inference by columns cannot name it, and whether the UPDATE changes
// an expression's value is not known.
func findDisjointUniqueKey(table *metadata.TableMetadata, pk *metadata.IndexMetadata, fields []string) (string, error) {
	columnMap := make(map[string]bool)
	for _, field := range fields {
		columnMap[field] = true
	}
	if pk != nil {
		if columns, ok := keyColumns(table, pk); ok && disjoint(columns, columnMap) {
			return "ON CONSTRAINT " + quoteIdentifier(pk.Name), nil
		}
	}
	for _, index := range table.GetIndexes() {
		if index.Primary || !index.Unique {
			continue
		}
		columns, ok := keyColumns(table, index)
		if !ok || !disjoint(columns, columnMap) {
			continue
		}
		if index.IsConstraint {
			return "ON CONSTRAINT " + quoteIdentifier(index.Name), nil
		}
		if index.Definition == "" || strings.Contains(index.Definition, " WHERE ") {
			continue
		}
		return "(" + quotePGColumns(columns) + ")", nil
	}

	return "", fmt.Errorf("no disjoint unique key found for table")
}

// keyColumns returns the columns of key's parts, or false when a part is an
// expression. A part is recorded as pg_get_indexdef writes it: a column's
// name, quoted when it needs to be.
func keyColumns(table *metadata.TableMetadata, key *metadata.IndexMetadata) ([]string, bool) {
	var columns []string
	for _, part := range key.Expressions {
		found := false
		for _, column := range table.GetColumns() {
			if part == column.Name || part == quoteIdentifier(column.Name) {
				columns = append(columns, column.Name)
				found = true
				break
			}
		}
		if !found {
			return nil, false
		}
	}
	return columns, true
}

// PrependStatements returns the first SET role or SET search_path of a parsed
// script, without its trailing ';', which a restore runs first; or "".
func PrependStatements(stmts []omnipg.Statement) string {
	for _, stmt := range stmts {
		if stmt.Empty() {
			continue
		}
		vs, ok := stmt.AST.(*ast.VariableSetStmt)
		if !ok {
			continue
		}
		name := strings.ToLower(vs.Name)
		if name == "role" || name == "search_path" {
			text := strings.TrimRight(strings.TrimSpace(stmt.Text), ";")
			return text
		}
	}

	return ""
}

// commentText returns s for the inside of a block comment: a space goes
// between every adjacent "/*" and "*/", since the engine nests block comments
// and either would move where the comment ends.
func commentText(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		b.WriteByte(s[i])
		if i+1 < len(s) && (s[i] == '/' && s[i+1] == '*' || s[i] == '*' && s[i+1] == '/') {
			b.WriteByte(' ')
		}
	}
	return b.String()
}
