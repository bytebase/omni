package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
)

const (
	maxCommentLength = 1000
)

// RestoreOptions names the backed-up table and its backup.
type RestoreOptions struct {
	SourceDatabase string
	SourceTable    string
	BackupDatabase string
	BackupTable    string
	// Table is the source table's metadata.
	Table *metadata.TableMetadata
	// PrimaryKey is the source table's primary key, or nil.
	PrimaryKey *metadata.IndexMetadata
}

// RestoreSQL returns the SQL that puts back the rows a backup holds: an
// INSERT for a DELETE, a MERGE on a unique key the UPDATE leaves alone for an
// UPDATE. node is FindFirstDML's, and originalSQL, the statements the backup
// covers, is quoted in a leading comment.
func RestoreSQL(originalSQL string, node ast.StmtNode, opts RestoreOptions) (string, error) {
	sqlForComment, truncated := truncateString(originalSQL, maxCommentLength)
	if truncated {
		sqlForComment += "..."
	}
	g := &generator{
		backupDatabase:   opts.BackupDatabase,
		backupTable:      opts.BackupTable,
		originalDatabase: opts.SourceDatabase,
		originalTable:    opts.SourceTable,
		pk:               opts.PrimaryKey,
		table:            opts.Table,
		isFirst:          true,
	}
	if err := g.generate(node); err != nil {
		return "", err
	}
	return fmt.Sprintf("/*\nOriginal SQL:\n%s\n*/\n%s", sqlForComment, g.result), nil
}

// NormalizeRestoreSQL trims a statement and its ";" for the restore's comment.
func NormalizeRestoreSQL(statement string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(statement), ";"))
}

// FindFirstDML returns the first UPDATE or DELETE in a script, or nil. When
// the whole script does not parse, it tries each run of lines, so the DML of
// a statement the backup sliced out of a larger one is still found.
func FindFirstDML(statement string) (ast.StmtNode, error) {
	node, err := findFirstOracleDMLOnce(statement)
	if err == nil {
		return node, nil
	}
	if !strings.Contains(statement, "\n") {
		return nil, fmt.Errorf("failed to parse statement: %w", err)
	}

	lines := strings.Split(statement, "\n")
	for start := range lines {
		var candidate strings.Builder
		for end := start; end < len(lines); end++ {
			line := strings.TrimSpace(lines[end])
			if line == "" && candidate.Len() == 0 {
				continue
			}
			if candidate.Len() > 0 {
				candidate.WriteByte('\n')
			}
			candidate.WriteString(lines[end])
			node, candidateErr := findFirstOracleDMLOnce(candidate.String())
			if candidateErr != nil {
				continue
			}
			if node != nil {
				return node, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to parse statement: %w", err)
}

func findFirstOracleDMLOnce(statement string) (ast.StmtNode, error) {
	list, err := Parse(statement)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	for _, item := range list.Items {
		raw, ok := item.(*ast.RawStmt)
		if !ok || raw.Stmt == nil {
			continue
		}
		switch raw.Stmt.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			return raw.Stmt, nil
		default:
		}
	}
	return nil, nil
}

type generator struct {
	backupDatabase   string
	backupTable      string
	originalDatabase string
	originalTable    string
	pk               *metadata.IndexMetadata
	table            *metadata.TableMetadata

	isFirst bool
	result  string
}

func (g *generator) generate(node ast.StmtNode) error {
	if !g.isFirst {
		return nil
	}
	g.isFirst = false

	switch n := node.(type) {
	case *ast.DeleteStmt:
		g.generateDelete()
	case *ast.UpdateStmt:
		return g.generateUpdate(n)
	default:
		return fmt.Errorf("unexpected statement type: %T", node)
	}
	return nil
}

func (g *generator) generateDelete() {
	columnList := quoteOracleColumns(g.restorableColumns())
	g.result = fmt.Sprintf(`INSERT INTO "%s"."%s" (%s) SELECT %s FROM "%s"."%s";`, g.originalDatabase, g.originalTable, columnList, columnList, g.backupDatabase, g.backupTable)
}

func (g *generator) restorableColumns() []string {
	var columns []string
	for _, column := range g.table.GetColumns() {
		// The backup is copied with SELECT t.*, which has no invisible columns.
		if column.GetGeneration() != nil || column.IsInvisible {
			continue
		}
		columns = append(columns, column.Name)
	}
	return columns
}

func quoteOracleColumns(columns []string) string {
	var quotedColumns []string
	for _, column := range columns {
		quotedColumns = append(quotedColumns, fmt.Sprintf(`"%s"`, column))
	}
	return strings.Join(quotedColumns, ", ")
}

func disjoint(a []string, b map[string]bool) bool {
	for _, item := range a {
		if _, ok := b[item]; ok {
			return false
		}
	}
	return true
}

func (g *generator) findDisjointUniqueKey(columns []string) ([]string, error) {
	columnMap := make(map[string]bool)
	for _, column := range columns {
		columnMap[column] = true
	}
	if g.pk != nil {
		if disjoint(g.pk.Expressions, columnMap) {
			return g.pk.Expressions, nil
		}
	}
	for _, index := range g.table.Indexes {
		if index.Primary {
			continue
		}
		if !index.Unique {
			continue
		}
		if strings.Contains(index.Type, "FUNCTION-BASED") {
			continue
		}
		if disjoint(index.Expressions, columnMap) {
			return index.Expressions, nil
		}
	}
	return nil, fmt.Errorf("no disjoint unique key found for %s.%s", g.originalDatabase, g.originalTable)
}

func (g *generator) generateUpdate(stmt *ast.UpdateStmt) error {
	// SET VALUE (t) = ... replaces whole rows, whose columns the restore
	// cannot list.
	if stmt.SetClauses != nil {
		for _, item := range stmt.SetClauses.Items {
			if clause, ok := item.(*ast.SetClause); ok && clause.ValueOf != "" {
				return fmt.Errorf("cannot restore SET VALUE (%s), which replaces whole rows", clause.ValueOf)
			}
		}
	}
	updateColumns := extractUpdateColumns(stmt)
	uk, err := g.findDisjointUniqueKey(updateColumns)
	if err != nil {
		return err
	}

	var buf strings.Builder
	if _, err := fmt.Fprintf(&buf, "MERGE INTO \"%s\".\"%s\" t\nUSING \"%s\".\"%s\" b\n  ON(", g.originalDatabase, g.originalTable, g.backupDatabase, g.backupTable); err != nil {
		return err
	}
	for i, column := range uk {
		if i > 0 {
			if _, err := fmt.Fprint(&buf, " AND"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(&buf, " t.\"%s\" = b.\"%s\"", column, column); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(&buf, ")\nWHEN MATCHED THEN\n  UPDATE SET"); err != nil {
		return err
	}
	for i, field := range updateColumns {
		if i > 0 {
			if _, err := fmt.Fprint(&buf, ","); err != nil {
				return err
			}
		}
		// The field is written by user and no need to escape.
		if _, err := fmt.Fprintf(&buf, " t.\"%s\" = b.\"%s\"", field, field); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(&buf, "\nWHEN NOT MATCHED THEN\n INSERT ("); err != nil {
		return err
	}
	for i, column := range g.restorableColumns() {
		if i > 0 {
			if _, err := fmt.Fprint(&buf, ", "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(&buf, "\"%s\"", column); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(&buf, ") VALUES ("); err != nil {
		return err
	}
	for i, column := range g.restorableColumns() {
		if i > 0 {
			if _, err := fmt.Fprint(&buf, ", "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(&buf, "b.\"%s\"", column); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(&buf, ");"); err != nil {
		return err
	}
	g.result = buf.String()
	return nil
}

func extractUpdateColumns(stmt *ast.UpdateStmt) []string {
	var result []string
	if stmt == nil || stmt.SetClauses == nil {
		return result
	}
	for _, item := range stmt.SetClauses.Items {
		clause, ok := item.(*ast.SetClause)
		if !ok {
			continue
		}
		if clause.Column != nil {
			result = append(result, clause.Column.Column)
		}
		for _, column := range columnRefList(clause.Columns) {
			result = append(result, column.Column)
		}
	}
	return result
}

func columnRefList(list *ast.List) []*ast.ColumnRef {
	if list == nil {
		return nil
	}
	result := make([]*ast.ColumnRef, 0, len(list.Items))
	for _, item := range list.Items {
		if column, ok := item.(*ast.ColumnRef); ok {
			result = append(result, column)
		}
	}
	return result
}
