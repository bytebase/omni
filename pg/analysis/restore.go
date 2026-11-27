package analysis

import (
	"fmt"
	"slices"
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

// changesTable reports whether node, an UPDATE or DELETE, changes table in
// schema. An unqualified table matches by name.
func changesTable(node ast.Node, schema, table string) bool {
	var relation *ast.RangeVar
	switch n := node.(type) {
	case *ast.UpdateStmt:
		relation = n.Relation
	case *ast.DeleteStmt:
		relation = n.Relation
	default:
	}
	return relation != nil && relation.Relname == table && (relation.Schemaname == "" || relation.Schemaname == schema)
}

// RestoreSQL returns the statement that restores, from its backup table, the
// rows nodes, the UPDATE and DELETE statements one backup covers, changed:
// when any is an UPDATE, an upsert of every column the UPDATEs set on a
// unique key none of them changes, which also inserts deleted rows again;
// otherwise an INSERT of the deleted rows. A comment quotes the original
// statements. UPDATEs that change every unique key are an error.
func RestoreSQL(nodes []ast.Node, opts RestoreOptions) (string, error) {
	if len(nodes) == 0 {
		return "", fmt.Errorf("no DML statement to restore")
	}
	var fields []string
	hasUpdate := false
	for _, node := range nodes {
		switch n := node.(type) {
		case *ast.DeleteStmt:
		case *ast.UpdateStmt:
			hasUpdate = true
			for _, field := range extractSetFieldNames(n) {
				if !slices.Contains(fields, field) {
					fields = append(fields, field)
				}
			}
		default:
			return "", fmt.Errorf("unexpected statement type: %T", node)
		}
	}
	sqlForComment := opts.OriginalSQL
	prependStatements := opts.PrependStatements
	backupSchema := opts.BackupSchema
	backupTable := opts.BackupTable
	originalSchema := opts.OriginalSchema
	originalTable := opts.OriginalTable
	quotedColumnList := quotePGColumns(restorableColumns(opts.Table))
	// PostgreSQL takes a value for a GENERATED ALWAYS identity column only
	// with OVERRIDING SYSTEM VALUE.
	overriding := ""
	if len(alwaysIdentityColumns(opts.Table)) > 0 {
		overriding = " OVERRIDING SYSTEM VALUE"
	}

	var result string
	if !hasUpdate {
		result = fmt.Sprintf(`INSERT INTO %s.%s (%s)%s SELECT %s FROM %s.%s;`, quoteIdentifier(originalSchema), quoteIdentifier(originalTable), quotedColumnList, overriding, quotedColumnList, quoteIdentifier(backupSchema), quoteIdentifier(backupTable))
	} else {
		// ON CONFLICT DO UPDATE can set such a column only to DEFAULT.
		for _, column := range alwaysIdentityColumns(opts.Table) {
			if slices.Contains(fields, column) {
				return "", fmt.Errorf("cannot restore an UPDATE of %q, a GENERATED ALWAYS identity column", column)
			}
		}
		conflictTarget, err := findDisjointUniqueKey(opts.Table, opts.PrimaryKey, fields)
		if err != nil {
			return "", err
		}

		var buf strings.Builder
		if _, err := fmt.Fprintf(&buf, `INSERT INTO %s.%s (%s)%s SELECT %s FROM %s.%s ON CONFLICT %s DO UPDATE SET `, quoteIdentifier(originalSchema), quoteIdentifier(originalTable), quotedColumnList, overriding, quotedColumnList, quoteIdentifier(backupSchema), quoteIdentifier(backupTable), conflictTarget); err != nil {
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

func alwaysIdentityColumns(table *metadata.TableMetadata) []string {
	var columns []string
	for _, column := range table.GetColumns() {
		if column.GetGeneration() == nil && column.GetIdentityGeneration() == metadata.ColumnMetadata_ALWAYS {
			columns = append(columns, column.Name)
		}
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

// extractSetFieldNames returns the columns an UPDATE's SET clause assigns,
// each once. PostgreSQL does not allow a SET target qualified by its table, so
// in "SET fn.first = 1" and "SET arr[1] = 1" the column is the first name, fn
// or arr; the rest picks a field or element of it.
func extractSetFieldNames(stmt *ast.UpdateStmt) []string {
	if stmt.TargetList == nil {
		return nil
	}
	var fields []string
	for _, item := range stmt.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Name == "" || slices.Contains(fields, rt.Name) {
			continue
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
