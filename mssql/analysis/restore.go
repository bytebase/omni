package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
)

const maxCommentLength = 1000

// RestoreOptions configures RestoreSQL.
type RestoreOptions struct {
	// OriginalDatabase, OriginalSchema and OriginalTable name the table the
	// backup came from.
	OriginalDatabase string
	OriginalSchema   string
	OriginalTable    string
	// BackupDatabase and BackupTable name the backup table, in schema dbo.
	BackupDatabase string
	BackupTable    string
	// Table is the original table's metadata, and PrimaryKey its primary key
	// or nil.
	Table      *metadata.TableMetadata
	PrimaryKey *metadata.IndexMetadata
	// OriginalSQL is the backed-up SQL the restore's comment quotes, cut to
	// 1000 characters.
	OriginalSQL string
}

// BackedUpDML is an UPDATE or DELETE of a script, with the positions a prior
// backup records for it.
type BackedUpDML struct {
	Node ast.Node
	// SQL is the statement's text without its trailing ";".
	SQL           string
	StartPosition Position
	EndPosition   Position
}

// BackedUpDMLs returns the UPDATE and DELETE statements of a parsed script,
// the statements a prior backup covers, in order.
func BackedUpDMLs(stmts []ParsedStatement) []BackedUpDML {
	var dmls []BackedUpDML
	for _, stmt := range stmts {
		if stmt.Node == nil || !isRestorableDML(stmt.Node) {
			continue
		}
		loc := dmlNodeLoc(stmt.Node)
		dmls = append(dmls, BackedUpDML{
			Node:          stmt.Node,
			SQL:           strings.TrimSuffix(sourceFromLoc(stmt.Text, loc), ";"),
			StartPosition: positionFromByteOffset(stmt.StartLine, stmt.Text, loc.Start),
			EndPosition:   positionFromByteOffset(stmt.StartLine, stmt.Text, dmlEndOffset(stmt.Text, loc)),
		})
	}
	return dmls
}

// changesTable reports whether node, an UPDATE or DELETE, changes table in
// schema, resolving its target as a backup does.
func changesTable(node ast.Node, database, schema, table string) bool {
	var target *tableReference
	var err error
	switch n := node.(type) {
	case *ast.UpdateStmt:
		target, _, err = resolveDMLTargetTable(n.Relation, n.FromClause, database)
	case *ast.DeleteStmt:
		target, _, err = resolveDMLTargetTable(n.Relation, n.FromClause, database)
	default:
	}
	return err == nil && target != nil && target.Schema == schema && target.Table == table
}

// RestoreSQL returns the statements that restore, from the backup table, the
// rows nodes, the UPDATE and DELETE statements one backup covers, changed:
// when any is an UPDATE, a MERGE of every column the UPDATEs set on a unique
// key none of them changes, which also inserts deleted rows again; otherwise
// an INSERT of the deleted rows. A table with an IDENTITY column turns
// IDENTITY_INSERT on around them and reseeds after. A comment quotes the
// original SQL. UPDATEs that change every unique key are an error.
func RestoreSQL(nodes []ast.Node, opts RestoreOptions) (string, error) {
	sqlForComment, truncated := truncateString(opts.OriginalSQL, maxCommentLength)
	if truncated {
		sqlForComment += "..."
	}
	sqlForComment = commentText(sqlForComment)
	g := &generator{
		backupDatabase:   opts.BackupDatabase,
		backupTable:      opts.BackupTable,
		originalDatabase: opts.OriginalDatabase,
		originalSchema:   opts.OriginalSchema,
		originalTable:    opts.OriginalTable,
		pk:               opts.PrimaryKey,
		table:            opts.Table,
	}
	result, err := g.generate(nodes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/*\nOriginal SQL:\n%s\n*/\n%s", sqlForComment, result), nil
}

type generator struct {
	backupDatabase   string
	backupTable      string
	originalDatabase string
	originalSchema   string
	originalTable    string
	pk               *metadata.IndexMetadata
	table            *metadata.TableMetadata
}

// hasIdentityColumn checks if the table has any IDENTITY columns
func (g *generator) hasIdentityColumn() bool {
	if g.table == nil {
		return false
	}
	for _, col := range g.table.GetColumns() {
		if col.IsIdentity {
			return true
		}
	}
	return false
}

func (g *generator) restorableColumns() []string {
	var columns []string
	for _, column := range g.table.GetColumns() {
		if column.GetGeneration() != nil {
			continue
		}
		columns = append(columns, column.Name)
	}
	return columns
}

func quoteTSQLColumns(columns []string) string {
	var quotedColumns []string
	for _, column := range columns {
		quotedColumns = append(quotedColumns, fmt.Sprintf("[%s]", column))
	}
	return strings.Join(quotedColumns, ", ")
}

func (g *generator) generate(nodes []ast.Node) (string, error) {
	if len(nodes) == 0 {
		return "", fmt.Errorf("no DML statement to restore")
	}
	var updateColumns []string
	seen := make(map[string]bool)
	hasUpdate := false
	for _, node := range nodes {
		switch n := node.(type) {
		case *ast.DeleteStmt:
		case *ast.UpdateStmt:
			hasUpdate = true
			for _, column := range updateSetColumns(n.SetClause) {
				if name := g.tableColumn(column); !seen[name] {
					seen[name] = true
					updateColumns = append(updateColumns, column)
				}
			}
		default:
			return "", fmt.Errorf("unsupported restore statement %T", node)
		}
	}
	if !hasUpdate {
		return g.generateDelete()
	}
	return g.generateUpdate(updateColumns)
}

// tableColumn returns the table's column that name, a SET column, assigns:
// the one spelled the same, or else the one equal but for case, which is what
// a case-insensitive collation assigns. A name of no column is returned as is.
func (g *generator) tableColumn(name string) string {
	folded := ""
	for _, column := range g.table.GetColumns() {
		if column.Name == name {
			return name
		}
		if folded == "" && strings.EqualFold(column.Name, name) {
			folded = column.Name
		}
	}
	if folded != "" {
		return folded
	}
	return name
}

func (g *generator) generateDelete() (string, error) {
	// Check if the table has IDENTITY columns
	hasIdentity := g.hasIdentityColumn()
	columnList := quoteTSQLColumns(g.restorableColumns())

	if hasIdentity {
		// For tables with IDENTITY columns, we need to enable IDENTITY_INSERT
		// and use explicit column lists
		var buf strings.Builder

		fmt.Fprintf(&buf, "SET IDENTITY_INSERT [%s].[%s].[%s] ON;\n",
			g.originalDatabase, g.originalSchema, g.originalTable)
		fmt.Fprintf(&buf, "INSERT INTO [%s].[%s].[%s] (%s) SELECT %s FROM [%s].[dbo].[%s];\n",
			g.originalDatabase, g.originalSchema, g.originalTable,
			columnList,
			columnList,
			g.backupDatabase, g.backupTable)
		fmt.Fprintf(&buf, "SET IDENTITY_INSERT [%s].[%s].[%s] OFF;\n",
			g.originalDatabase, g.originalSchema, g.originalTable)
		fmt.Fprintf(&buf, "EXEC('DBCC CHECKIDENT (''[%s].[%s].[%s]'', RESEED)');",
			g.originalDatabase, g.originalSchema, g.originalTable)

		return buf.String(), nil
	}

	// Simple INSERT for tables without IDENTITY columns
	return fmt.Sprintf(`INSERT INTO [%s].[%s].[%s] (%s) SELECT %s FROM [%s].[dbo].[%s];`,
		g.originalDatabase, g.originalSchema, g.originalTable,
		columnList,
		columnList,
		g.backupDatabase, g.backupTable), nil
}

func disjoint(a []string, b map[string]bool) bool {
	for _, item := range a {
		if _, ok := b[item]; ok {
			return false
		}
	}
	return true
}

func (g *generator) findDisjointUniqueKey(updateColumns []string) ([]string, error) {
	columnMap := make(map[string]bool)
	for _, column := range updateColumns {
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
		if disjoint(index.Expressions, columnMap) {
			return index.Expressions, nil
		}
	}
	return nil, fmt.Errorf("no disjoint unique key found for %s.%s.%s", g.originalDatabase, g.originalSchema, g.originalTable)
}

func (g *generator) generateUpdate(updateColumns []string) (string, error) {
	uk, err := g.findDisjointUniqueKey(updateColumns)
	if err != nil {
		return "", err
	}

	var buf strings.Builder

	// Check if the table has IDENTITY columns
	hasIdentity := g.hasIdentityColumn()

	if hasIdentity {
		// Enable IDENTITY_INSERT for MERGE statement
		if _, err := fmt.Fprintf(&buf, "SET IDENTITY_INSERT [%s].[%s].[%s] ON;\n",
			g.originalDatabase, g.originalSchema, g.originalTable); err != nil {
			return "", err
		}
	}

	if _, err := fmt.Fprintf(&buf, "MERGE INTO [%s].[%s].[%s] AS t\nUSING [%s].[dbo].[%s] AS b\n  ON", g.originalDatabase, g.originalSchema, g.originalTable, g.backupDatabase, g.backupTable); err != nil {
		return "", err
	}
	for i, column := range uk {
		if i > 0 {
			if _, err := fmt.Fprint(&buf, " AND"); err != nil {
				return "", err
			}
		}
		if _, err := fmt.Fprintf(&buf, " t.[%s] = b.[%s]", column, column); err != nil {
			return "", err
		}
	}
	if _, err := fmt.Fprint(&buf, "\nWHEN MATCHED THEN\n  UPDATE SET"); err != nil {
		return "", err
	}
	for i, field := range updateColumns {
		if i > 0 {
			if _, err := fmt.Fprint(&buf, ","); err != nil {
				return "", err
			}
		}
		if _, err := fmt.Fprintf(&buf, " t.[%s] = b.[%s]", field, field); err != nil {
			return "", err
		}
	}
	if _, err := fmt.Fprint(&buf, "\nWHEN NOT MATCHED THEN\n INSERT ("); err != nil {
		return "", err
	}
	for i, column := range g.restorableColumns() {
		if i > 0 {
			if _, err := fmt.Fprint(&buf, ", "); err != nil {
				return "", err
			}
		}
		if _, err := fmt.Fprintf(&buf, "[%s]", column); err != nil {
			return "", err
		}
	}
	if _, err := fmt.Fprint(&buf, ") VALUES ("); err != nil {
		return "", err
	}
	for i, column := range g.restorableColumns() {
		if i > 0 {
			if _, err := fmt.Fprint(&buf, ", "); err != nil {
				return "", err
			}
		}
		if _, err := fmt.Fprintf(&buf, "b.[%s]", column); err != nil {
			return "", err
		}
	}
	if _, err := fmt.Fprint(&buf, ");"); err != nil {
		return "", err
	}

	// Check if we need to disable IDENTITY_INSERT and reseed
	if hasIdentity {
		if _, err := fmt.Fprintf(&buf, "\n\nSET IDENTITY_INSERT [%s].[%s].[%s] OFF;\nEXEC('DBCC CHECKIDENT (''[%s].[%s].[%s]'', RESEED)');",
			g.originalDatabase, g.originalSchema, g.originalTable,
			g.originalDatabase, g.originalSchema, g.originalTable); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func updateSetColumns(setClause *ast.List) []string {
	if setClause == nil {
		return nil
	}
	var columns []string
	for _, item := range setClause.Items {
		setExpr, ok := item.(*ast.SetExpr)
		if !ok {
			continue
		}
		if setExpr.Column != nil {
			columns = append(columns, setExpr.Column.Column)
			continue
		}
		if setExpr.VarColumn != nil {
			columns = append(columns, setExpr.VarColumn.Column)
			continue
		}
		if setExpr.Variable != "" {
			if column := updateColumnFromDualAssignment(setExpr.Value); column != "" {
				columns = append(columns, column)
			}
		}
	}
	return columns
}

func updateColumnFromDualAssignment(expr ast.ExprNode) string {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op == ast.BinOpEq {
			return updateColumnFromExpr(e.Left)
		}
	case *ast.ParenExpr:
		return updateColumnFromDualAssignment(e.Expr)
	default:
	}
	return ""
}

func updateColumnFromExpr(expr ast.ExprNode) string {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		return e.Column
	case *ast.ParenExpr:
		return updateColumnFromExpr(e.Expr)
	default:
	}
	return ""
}

func isRestorableDML(node ast.Node) bool {
	switch node.(type) {
	case *ast.UpdateStmt, *ast.DeleteStmt:
		return true
	default:
		return false
	}
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
