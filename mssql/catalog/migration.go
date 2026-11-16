package catalog

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// MigrationOpType names one emitted chunk of a migration script.
type MigrationOpType string

// OpScript is the whole script. SQL Server's generator interleaves GO batch
// separators with its statements, so the plan keeps the script whole rather
// than re-joining statements.
const OpScript MigrationOpType = "Script"

// MigrationOp is one operation in a migration plan.
type MigrationOp struct {
	Type       MigrationOpType
	SchemaName string
	ObjectName string
	SQL        string
}

// MigrationPlan holds the DDL that transforms one catalog state into another.
type MigrationPlan struct {
	Ops []MigrationOp
}

// SQL returns the migration script.
func (p *MigrationPlan) SQL() string {
	var buf strings.Builder
	for _, op := range p.Ops {
		buf.WriteString(op.SQL)
	}
	return buf.String()
}

// IsEmpty reports whether the plan emits no DDL.
func (p *MigrationPlan) IsEmpty() bool { return p.SQL() == "" }

// GenerateMigration produces the script that transforms from into to, using the
// precomputed diff. It is Bytebase's schema/mssql generator ported verbatim.
func GenerateMigration(from, to *Catalog, diff *SchemaDiff) *MigrationPlan {
	_ = from
	_ = to
	var buf strings.Builder

	var schemasToCreate []string
	for _, schemaDiff := range diff.Schemas {
		if schemaDiff.Action == DiffCreate && strings.ToLower(schemaDiff.SchemaName) != defaultSchema {
			schemasToCreate = append(schemasToCreate, schemaDiff.SchemaName)
		}
	}
	slices.Sort(schemasToCreate)

	// Drop dependent objects first: foreign keys, procedures, functions, views,
	// check constraints, indexes, columns, tables, then schemas.
	for _, tableDiff := range diff.Tables {
		switch tableDiff.Action {
		case DiffAlter:
			for _, fkDiff := range tableDiff.ForeignKeys {
				if fkDiff.Action == DiffDrop {
					fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] DROP CONSTRAINT [%s];\nGO\n", tableDiff.SchemaName, tableDiff.TableName, fkDiff.OldForeignKey.Name)
				}
			}
		case DiffDrop:
			for _, fk := range tableDiff.OldTable.ForeignKeys {
				fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] DROP CONSTRAINT [%s];\n", tableDiff.SchemaName, tableDiff.TableName, fk.Name)
			}
		}
	}

	for _, procDiff := range diff.Procedures {
		if procDiff.Action == DiffDrop {
			fmt.Fprintf(&buf, "DROP PROCEDURE [%s].[%s];\nGO\n", procDiff.SchemaName, procDiff.ProcedureName)
		}
	}
	for _, funcDiff := range diff.Functions {
		if funcDiff.Action == DiffDrop {
			fmt.Fprintf(&buf, "DROP FUNCTION [%s].[%s];\nGO\n", funcDiff.SchemaName, funcDiff.FunctionName)
		}
	}

	dropViewsInOrder(diff, &buf)

	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffAlter {
			continue
		}
		for _, checkDiff := range tableDiff.Checks {
			if checkDiff.Action == DiffDrop {
				fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] DROP CONSTRAINT [%s];\n", tableDiff.SchemaName, tableDiff.TableName, checkDiff.OldCheckConstraint.Name)
			}
		}
		for _, indexDiff := range tableDiff.Indexes {
			if indexDiff.Action != DiffDrop {
				continue
			}
			if indexDiff.OldIndex.IsConstraint {
				fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] DROP CONSTRAINT [%s];\n", tableDiff.SchemaName, tableDiff.TableName, indexDiff.OldIndex.Name)
			} else {
				fmt.Fprintf(&buf, "DROP INDEX [%s] ON [%s].[%s];\n", indexDiff.OldIndex.Name, tableDiff.SchemaName, tableDiff.TableName)
			}
		}

		hasDroppedIndex := false
		for _, indexDiff := range tableDiff.Indexes {
			if indexDiff.Action == DiffDrop {
				hasDroppedIndex = true
				break
			}
		}
		for i, colDiff := range tableDiff.Columns {
			if colDiff.Action != DiffDrop {
				continue
			}
			if i == 0 && hasDroppedIndex {
				buf.WriteString("GO\n")
			}
			if getColumnDefaultValue(colDiff.OldColumn) != "" && colDiff.OldColumn.DefaultConstraintName != "" {
				fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] DROP CONSTRAINT [%s];\n", tableDiff.SchemaName, tableDiff.TableName, colDiff.OldColumn.DefaultConstraintName)
			}
			fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] DROP COLUMN [%s];\n", tableDiff.SchemaName, tableDiff.TableName, colDiff.OldColumn.Name)
		}
	}

	hasTableAlterations := false
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffAlter {
			continue
		}
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action == DiffDrop {
				hasTableAlterations = true
				break
			}
		}
	}

	hasTableDrops := false
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffDrop {
			continue
		}
		if !hasTableDrops && hasTableAlterations {
			buf.WriteString("GO\n")
		}
		hasTableDrops = true
		fmt.Fprintf(&buf, "DROP TABLE [%s].[%s];\n", tableDiff.SchemaName, tableDiff.TableName)
	}

	for _, schemaDiff := range diff.Schemas {
		if schemaDiff.Action != DiffDrop || schemaDiff.OldSchema == nil {
			continue
		}
		for _, table := range schemaDiff.OldSchema.Tables {
			alreadyHandled := false
			for _, tableDiff := range diff.Tables {
				if tableDiff.Action == DiffDrop && tableDiff.SchemaName == schemaDiff.SchemaName && tableDiff.TableName == table.Name {
					alreadyHandled = true
					break
				}
			}
			if !alreadyHandled {
				fmt.Fprintf(&buf, "DROP TABLE [%s].[%s];\n", schemaDiff.SchemaName, table.Name)
			}
		}
	}

	for _, schemaDiff := range diff.Schemas {
		if schemaDiff.Action != DiffDrop || strings.ToLower(schemaDiff.SchemaName) == defaultSchema {
			continue
		}
		fmt.Fprintf(&buf, "DROP SCHEMA [%s];\nGO\n", schemaDiff.SchemaName)
	}

	dropPhaseHasContent := buf.Len() > 0
	createPhaseWillHaveContent := len(schemasToCreate) > 0 ||
		hasCreateOrAlterTables(diff) || hasCreateViews(diff) ||
		hasCreateFunctions(diff) || hasCreateProcedures(diff)
	if dropPhaseHasContent && createPhaseWillHaveContent {
		buf.WriteString("\n")
	}

	for _, schemaName := range schemasToCreate {
		fmt.Fprintf(&buf, "CREATE SCHEMA [%s];\nGO\n", schemaName)
	}
	if len(schemasToCreate) > 0 && (hasCreateOrAlterTables(diff) || hasCreateViews(diff) || hasCreateFunctions(diff) || hasCreateProcedures(diff)) {
		buf.WriteString("\n")
	}

	createTablesInOrder(diff, &buf)

	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffAlter {
			alterTableSQL := generateAlterTable(tableDiff)
			buf.WriteString(alterTableSQL)
			if alterTableSQL != "" {
				buf.WriteString("\n")
			}
		}
	}

	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffCreate || tableDiff.NewTable == nil {
			continue
		}
		for _, fk := range tableDiff.NewTable.ForeignKeys {
			buf.WriteString(generateAddForeignKey(tableDiff.SchemaName, tableDiff.TableName, fk))
			buf.WriteString(";\n")
		}
	}

	hasCommentChanges := false
	for _, tableDiff := range diff.Tables {
		switch tableDiff.Action {
		case DiffCreate:
			if tableDiff.NewTable == nil {
				continue
			}
			if tableDiff.NewTable.Comment != "" {
				buf.WriteString(generateTableCommentSQL("ADD", tableDiff.SchemaName, tableDiff.TableName, tableDiff.NewTable.Comment))
				buf.WriteString(";\n")
				hasCommentChanges = true
			}
			for _, column := range tableDiff.NewTable.Columns {
				if column.Comment != "" {
					buf.WriteString(generateColumnCommentSQL("ADD", tableDiff.SchemaName, tableDiff.TableName, column.Name, column.Comment))
					buf.WriteString(";\n")
					hasCommentChanges = true
				}
			}
		case DiffAlter:
			oldTableComment, newTableComment := "", ""
			if tableDiff.OldTable != nil {
				oldTableComment = tableDiff.OldTable.Comment
			}
			if tableDiff.NewTable != nil {
				newTableComment = tableDiff.NewTable.Comment
			}
			if oldTableComment != newTableComment {
				switch {
				case oldTableComment == "" && newTableComment != "":
					buf.WriteString(generateTableCommentSQL("ADD", tableDiff.SchemaName, tableDiff.TableName, newTableComment))
					buf.WriteString(";\n")
					hasCommentChanges = true
				case oldTableComment != "" && newTableComment == "":
					buf.WriteString(generateTableCommentSQL("DROP", tableDiff.SchemaName, tableDiff.TableName, ""))
					buf.WriteString(";\n")
					hasCommentChanges = true
				default:
					buf.WriteString(generateTableCommentSQL("UPDATE", tableDiff.SchemaName, tableDiff.TableName, newTableComment))
					buf.WriteString(";\n")
					hasCommentChanges = true
				}
			}
			for _, colDiff := range tableDiff.Columns {
				switch colDiff.Action {
				case DiffCreate:
					if colDiff.NewColumn != nil && colDiff.NewColumn.Comment != "" {
						buf.WriteString(generateColumnCommentSQL("ADD", tableDiff.SchemaName, tableDiff.TableName, colDiff.NewColumn.Name, colDiff.NewColumn.Comment))
						buf.WriteString(";\n")
						hasCommentChanges = true
					}
				case DiffAlter:
					oldComment, newComment := "", ""
					if colDiff.OldColumn != nil {
						oldComment = colDiff.OldColumn.Comment
					}
					if colDiff.NewColumn != nil {
						newComment = colDiff.NewColumn.Comment
					}
					if oldComment == newComment {
						continue
					}
					switch {
					case oldComment == "" && newComment != "":
						buf.WriteString(generateColumnCommentSQL("ADD", tableDiff.SchemaName, tableDiff.TableName, colDiff.NewColumn.Name, newComment))
						buf.WriteString(";\n")
						hasCommentChanges = true
					case oldComment != "" && newComment == "":
						buf.WriteString(generateColumnCommentSQL("DROP", tableDiff.SchemaName, tableDiff.TableName, colDiff.NewColumn.Name, ""))
						buf.WriteString(";\n")
						hasCommentChanges = true
					default:
						buf.WriteString(generateColumnCommentSQL("UPDATE", tableDiff.SchemaName, tableDiff.TableName, colDiff.NewColumn.Name, newComment))
						buf.WriteString(";\n")
						hasCommentChanges = true
					}
				}
			}
		}
	}
	if hasCommentChanges {
		buf.WriteString("GO\n")
	}
	if hasCreateOrAlterTables(diff) && (hasCreateViews(diff) || hasCreateFunctions(diff) || hasCreateProcedures(diff)) {
		buf.WriteString("\nGO\n")
	}

	createViewsInOrder(diff, &buf)

	for _, funcDiff := range diff.Functions {
		if funcDiff.Action != DiffCreate {
			continue
		}
		buf.WriteString(funcDiff.NewFunction.Definition)
		if !strings.HasSuffix(strings.TrimSpace(funcDiff.NewFunction.Definition), ";") {
			buf.WriteString(";")
		}
		buf.WriteString("\nGO\n")
		if funcDiff.NewFunction.Comment != "" {
			buf.WriteString(generateFunctionCommentSQL("ADD", funcDiff.SchemaName, funcDiff.FunctionName, funcDiff.NewFunction.Comment))
			buf.WriteString(";\nGO\n")
		}
	}

	for _, procDiff := range diff.Procedures {
		switch procDiff.Action {
		case DiffCreate:
			buf.WriteString(procDiff.NewProcedure.Definition)
			if !strings.HasSuffix(strings.TrimSpace(procDiff.NewProcedure.Definition), ";") {
				buf.WriteString(";")
			}
			buf.WriteString("\nGO\n")
		case DiffAlter:
			fmt.Fprintf(&buf, "DROP PROCEDURE [%s].[%s];\nGO\n", procDiff.SchemaName, procDiff.ProcedureName)
			buf.WriteString(procDiff.NewProcedure.Definition)
			if !strings.HasSuffix(strings.TrimSpace(procDiff.NewProcedure.Definition), ";") {
				buf.WriteString(";")
			}
			buf.WriteString("\nGO\n")
		}
	}

	for _, sequenceDiff := range diff.Sequences {
		if sequenceDiff.Action != DiffCreate {
			continue
		}
		fmt.Fprintf(&buf, "CREATE SEQUENCE [%s].[%s]", sequenceDiff.SchemaName, sequenceDiff.SequenceName)
		if sequenceDiff.NewSequence.DataType != "" {
			buf.WriteString(" AS ")
			buf.WriteString(sequenceDiff.NewSequence.DataType)
		}
		buf.WriteString(";\nGO\n")
		if sequenceDiff.NewSequence.Comment != "" {
			buf.WriteString(generateSequenceCommentSQL("ADD", sequenceDiff.SchemaName, sequenceDiff.SequenceName, sequenceDiff.NewSequence.Comment))
			buf.WriteString(";\nGO\n")
		}
	}

	text := buf.String()
	var ops []MigrationOp
	if text != "" {
		ops = append(ops, MigrationOp{Type: OpScript, SQL: text})
	}
	return &MigrationPlan{Ops: ops}
}

func generateCreateTable(schemaName, tableName string, table *metadata.TableMetadata) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "CREATE TABLE [%s].[%s] (\n", schemaName, tableName)
	for i, column := range table.Columns {
		buf.WriteString("  ")
		buf.WriteString(generateColumnDefinition(column))
		if i < len(table.Columns)-1 || hasConstraintsInTable(table) || len(table.CheckConstraints) > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}

	constraintAdded := false
	for _, idx := range table.Indexes {
		if !idx.IsConstraint {
			continue
		}
		if constraintAdded {
			buf.WriteString(",\n")
		}
		fmt.Fprintf(&buf, "  CONSTRAINT [%s]", idx.Name)
		switch {
		case idx.Primary:
			buf.WriteString(" PRIMARY KEY")
		case idx.Unique:
			buf.WriteString(" UNIQUE")
		}
		switch idx.Type {
		case "CLUSTERED":
			buf.WriteString(" CLUSTERED")
		case "NONCLUSTERED":
			buf.WriteString(" NONCLUSTERED")
		}
		buf.WriteString(" (")
		for j, expr := range idx.Expressions {
			if j > 0 {
				buf.WriteString(", ")
			}
			fmt.Fprintf(&buf, "[%s]", expr)
		}
		buf.WriteString(")")
		constraintAdded = true
	}

	if constraintAdded && len(table.CheckConstraints) > 0 {
		buf.WriteString(",\n")
	} else if constraintAdded {
		buf.WriteString("\n")
	}
	for i, check := range table.CheckConstraints {
		fmt.Fprintf(&buf, "  CONSTRAINT [%s] CHECK (%s)", check.Name, check.Expression)
		if i < len(table.CheckConstraints)-1 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}
	buf.WriteString(");")

	for _, idx := range table.Indexes {
		if idx.IsConstraint {
			continue
		}
		buf.WriteString("\n")
		buf.WriteString(generateCreateIndex(schemaName, tableName, idx))
		buf.WriteString(";")
	}
	return buf.String()
}

func generateAlterTable(tableDiff *TableDiffEntry) string {
	var buf strings.Builder
	for _, colDiff := range tableDiff.Columns {
		if colDiff.Action == DiffCreate {
			fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] ADD %s;\n", tableDiff.SchemaName, tableDiff.TableName, generateColumnDefinition(colDiff.NewColumn))
		}
	}
	for _, colDiff := range tableDiff.Columns {
		if colDiff.Action == DiffAlter {
			buf.WriteString(generateAlterColumn(tableDiff.SchemaName, tableDiff.TableName, colDiff))
		}
	}
	for _, indexDiff := range tableDiff.Indexes {
		if indexDiff.Action != DiffCreate {
			continue
		}
		if indexDiff.NewIndex.IsConstraint {
			fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] ADD CONSTRAINT [%s]", tableDiff.SchemaName, tableDiff.TableName, indexDiff.NewIndex.Name)
			switch {
			case indexDiff.NewIndex.Primary:
				buf.WriteString(" PRIMARY KEY")
			case indexDiff.NewIndex.Unique:
				buf.WriteString(" UNIQUE")
			}
			switch indexDiff.NewIndex.Type {
			case "CLUSTERED":
				buf.WriteString(" CLUSTERED")
			case "NONCLUSTERED":
				buf.WriteString(" NONCLUSTERED")
			}
			buf.WriteString(" (")
			for j, expr := range indexDiff.NewIndex.Expressions {
				if j > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(&buf, "[%s]", expr)
			}
			buf.WriteString(");\n")
		} else {
			buf.WriteString(generateCreateIndex(tableDiff.SchemaName, tableDiff.TableName, indexDiff.NewIndex))
			buf.WriteString(";\n")
		}
	}
	for _, checkDiff := range tableDiff.Checks {
		if checkDiff.Action == DiffCreate {
			fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] ADD CONSTRAINT [%s] CHECK (%s);\n",
				tableDiff.SchemaName, tableDiff.TableName, checkDiff.NewCheckConstraint.Name, checkDiff.NewCheckConstraint.Expression)
		}
	}
	for _, fkDiff := range tableDiff.ForeignKeys {
		if fkDiff.Action == DiffCreate {
			buf.WriteString(generateAddForeignKey(tableDiff.SchemaName, tableDiff.TableName, fkDiff.NewForeignKey))
			buf.WriteString(";\n")
		}
	}
	return buf.String()
}

func generateColumnDefinition(column *metadata.ColumnMetadata) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "[%s] ", column.Name)
	isComputedColumn := strings.Contains(column.Type, " AS (") || strings.HasPrefix(column.Type, "AS (")
	if isComputedColumn {
		buf.WriteString(column.Type)
		return buf.String()
	}
	buf.WriteString(column.Type)
	if column.IsIdentity {
		fmt.Fprintf(&buf, " IDENTITY(%d,%d)", column.IdentitySeed, column.IdentityIncrement)
	}
	if column.Nullable {
		buf.WriteString(" NULL")
	} else {
		buf.WriteString(" NOT NULL")
	}
	if column.GetDefault() != "" {
		buf.WriteString(" DEFAULT ")
		buf.WriteString(column.GetDefault())
	}
	return buf.String()
}

func generateAlterColumn(schemaName, tableName string, colDiff *ColumnDiffEntry) string {
	var buf strings.Builder
	if colDiff.OldColumn.Type != colDiff.NewColumn.Type {
		fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] ALTER COLUMN [%s] %s", schemaName, tableName, colDiff.NewColumn.Name, colDiff.NewColumn.Type)
		if colDiff.NewColumn.Nullable {
			buf.WriteString(" NULL")
		} else {
			buf.WriteString(" NOT NULL")
		}
		buf.WriteString(";\n")
	} else if colDiff.OldColumn.Nullable != colDiff.NewColumn.Nullable {
		fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] ALTER COLUMN [%s] %s", schemaName, tableName, colDiff.NewColumn.Name, colDiff.NewColumn.Type)
		if colDiff.NewColumn.Nullable {
			buf.WriteString(" NULL")
		} else {
			buf.WriteString(" NOT NULL")
		}
		buf.WriteString(";\n")
	}
	oldDefault := getColumnDefaultValue(colDiff.OldColumn)
	newDefault := getColumnDefaultValue(colDiff.NewColumn)
	if oldDefault != newDefault {
		if oldDefault != "" && colDiff.OldColumn.DefaultConstraintName != "" {
			fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] DROP CONSTRAINT [%s];\n", schemaName, tableName, colDiff.OldColumn.DefaultConstraintName)
		}
		if newDefault != "" {
			fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] ADD CONSTRAINT [DF_%s_%s] DEFAULT %s FOR [%s];\n",
				schemaName, tableName, tableName, colDiff.NewColumn.Name, newDefault, colDiff.NewColumn.Name)
		}
	}
	return buf.String()
}

func getColumnDefaultValue(column *metadata.ColumnMetadata) string {
	if column == nil {
		return ""
	}
	return column.GetDefault()
}

func generateCreateIndex(schemaName, tableName string, index *metadata.IndexMetadata) string {
	var buf strings.Builder
	buf.WriteString("CREATE")
	switch strings.ToUpper(index.Type) {
	case "CLUSTERED COLUMNSTORE":
		fmt.Fprintf(&buf, " CLUSTERED COLUMNSTORE INDEX [%s] ON [%s].[%s]", index.Name, schemaName, tableName)
		return buf.String()
	case "NONCLUSTERED COLUMNSTORE":
		fmt.Fprintf(&buf, " NONCLUSTERED COLUMNSTORE INDEX [%s] ON [%s].[%s]", index.Name, schemaName, tableName)
		if len(index.Expressions) > 0 {
			buf.WriteString(" (")
			for i, expr := range index.Expressions {
				if i > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(&buf, "[%s]", expr)
			}
			buf.WriteString(")")
		}
		return buf.String()
	case "SPATIAL":
		return generateSpatialIndexDDL(index, schemaName, tableName)
	default:
		if index.Unique {
			buf.WriteString(" UNIQUE")
		}
		switch index.Type {
		case "CLUSTERED":
			buf.WriteString(" CLUSTERED")
		case "NONCLUSTERED":
			buf.WriteString(" NONCLUSTERED")
		}
		fmt.Fprintf(&buf, " INDEX [%s] ON [%s].[%s] (", index.Name, schemaName, tableName)
		for i, expr := range index.Expressions {
			if i > 0 {
				buf.WriteString(", ")
			}
			fmt.Fprintf(&buf, "[%s]", expr)
		}
		buf.WriteString(")")
		return buf.String()
	}
}

func generateAddForeignKey(schemaName, tableName string, fk *metadata.ForeignKeyMetadata) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "ALTER TABLE [%s].[%s] ADD CONSTRAINT [%s] FOREIGN KEY (", schemaName, tableName, fk.Name)
	for i, col := range fk.Columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "[%s]", col)
	}
	buf.WriteString(") REFERENCES ")
	if fk.ReferencedSchema != "" {
		fmt.Fprintf(&buf, "[%s].", fk.ReferencedSchema)
	}
	fmt.Fprintf(&buf, "[%s] (", fk.ReferencedTable)
	for i, col := range fk.ReferencedColumns {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "[%s]", col)
	}
	buf.WriteString(")")
	if fk.OnDelete != "" {
		buf.WriteString(" ON DELETE ")
		buf.WriteString(fk.OnDelete)
	}
	if fk.OnUpdate != "" {
		buf.WriteString(" ON UPDATE ")
		buf.WriteString(fk.OnUpdate)
	}
	return buf.String()
}

func hasConstraintsInTable(table *metadata.TableMetadata) bool {
	for _, idx := range table.Indexes {
		if idx.IsConstraint {
			return true
		}
	}
	return false
}

func generateExtendedPropertySQL(action, objectType, schemaName, objectName, comment string) string {
	var buf strings.Builder
	switch action {
	case "ADD":
		fmt.Fprintf(&buf, "EXEC sp_addextendedproperty 'MS_Description', '%s', 'SCHEMA', '%s', '%s', '%s'", strings.ReplaceAll(comment, "'", "''"), schemaName, objectType, objectName)
	case "UPDATE":
		fmt.Fprintf(&buf, "EXEC sp_updateextendedproperty 'MS_Description', '%s', 'SCHEMA', '%s', '%s', '%s'", strings.ReplaceAll(comment, "'", "''"), schemaName, objectType, objectName)
	case "DROP":
		fmt.Fprintf(&buf, "EXEC sp_dropextendedproperty 'MS_Description', 'SCHEMA', '%s', '%s', '%s'", schemaName, objectType, objectName)
	}
	return buf.String()
}

func generateViewCommentSQL(action, schemaName, viewName, comment string) string {
	return generateExtendedPropertySQL(action, "VIEW", schemaName, viewName, comment)
}

func generateFunctionCommentSQL(action, schemaName, functionName, comment string) string {
	return generateExtendedPropertySQL(action, "FUNCTION", schemaName, functionName, comment)
}

func generateSequenceCommentSQL(action, schemaName, sequenceName, comment string) string {
	return generateExtendedPropertySQL(action, "SEQUENCE", schemaName, sequenceName, comment)
}

func generateTableCommentSQL(action, schemaName, tableName, comment string) string {
	return generateExtendedPropertySQL(action, "TABLE", schemaName, tableName, comment)
}

func generateColumnCommentSQL(action, schemaName, tableName, columnName, comment string) string {
	var buf strings.Builder
	switch action {
	case "ADD":
		fmt.Fprintf(&buf, "EXEC sp_addextendedproperty 'MS_Description', '%s', 'SCHEMA', '%s', 'TABLE', '%s', 'COLUMN', '%s'", strings.ReplaceAll(comment, "'", "''"), schemaName, tableName, columnName)
	case "UPDATE":
		fmt.Fprintf(&buf, "EXEC sp_updateextendedproperty 'MS_Description', '%s', 'SCHEMA', '%s', 'TABLE', '%s', 'COLUMN', '%s'", strings.ReplaceAll(comment, "'", "''"), schemaName, tableName, columnName)
	case "DROP":
		fmt.Fprintf(&buf, "EXEC sp_dropextendedproperty 'MS_Description', 'SCHEMA', '%s', 'TABLE', '%s', 'COLUMN', '%s'", schemaName, tableName, columnName)
	}
	return buf.String()
}

func dropViewsInOrder(diff *SchemaDiff, buf *strings.Builder) {
	graph := newGraph()
	viewMap := make(map[string]*ViewDiffEntry)
	var viewsToProcess []*ViewDiffEntry
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffDrop || viewDiff.Action == DiffAlter {
			viewsToProcess = append(viewsToProcess, viewDiff)
		}
	}
	slices.SortFunc(viewsToProcess, func(a, b *ViewDiffEntry) int {
		return strings.Compare(objectID(a.SchemaName, a.ViewName), objectID(b.SchemaName, b.ViewName))
	})
	for _, viewDiff := range viewsToProcess {
		viewID := objectID(viewDiff.SchemaName, viewDiff.ViewName)
		graph.addNode(viewID)
		viewMap[viewID] = viewDiff
	}
	for _, viewDiff := range viewsToProcess {
		viewID := objectID(viewDiff.SchemaName, viewDiff.ViewName)
		if viewDiff.OldView == nil || viewDiff.OldView.Definition == "" {
			continue
		}
		deps, err := getViewDependencies(viewDiff.OldView.Definition, viewDiff.SchemaName)
		if err != nil {
			continue
		}
		for _, dep := range deps {
			if _, ok := viewMap[dep]; ok {
				graph.addEdge(viewID, dep)
			}
		}
	}
	ordered, err := graph.topologicalSort()
	if err != nil {
		for _, viewDiff := range viewsToProcess {
			fmt.Fprintf(buf, "DROP VIEW [%s].[%s];\nGO\n", viewDiff.SchemaName, viewDiff.ViewName)
		}
		return
	}
	for _, viewID := range ordered {
		if viewDiff, ok := viewMap[viewID]; ok {
			fmt.Fprintf(buf, "DROP VIEW [%s].[%s];\nGO\n", viewDiff.SchemaName, viewDiff.ViewName)
		}
	}
}

func createTablesInOrder(diff *SchemaDiff, buf *strings.Builder) {
	graph := newGraph()
	tableMap := make(map[string]*TableDiffEntry)
	var tablesToCreate []*TableDiffEntry
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffCreate {
			tablesToCreate = append(tablesToCreate, tableDiff)
		}
	}
	slices.SortFunc(tablesToCreate, func(a, b *TableDiffEntry) int {
		return strings.Compare(objectID(a.SchemaName, a.TableName), objectID(b.SchemaName, b.TableName))
	})
	for _, tableDiff := range tablesToCreate {
		tableID := objectID(tableDiff.SchemaName, tableDiff.TableName)
		graph.addNode(tableID)
		tableMap[tableID] = tableDiff
	}
	for _, tableDiff := range tablesToCreate {
		tableID := objectID(tableDiff.SchemaName, tableDiff.TableName)
		if tableDiff.NewTable == nil {
			continue
		}
		for _, fk := range tableDiff.NewTable.ForeignKeys {
			refID := objectID(fk.ReferencedSchema, fk.ReferencedTable)
			if _, ok := tableMap[refID]; ok {
				graph.addEdge(refID, tableID)
			}
		}
	}
	ordered, err := graph.topologicalSort()
	if err != nil {
		for _, tableDiff := range tablesToCreate {
			buf.WriteString(generateCreateTable(tableDiff.SchemaName, tableDiff.TableName, tableDiff.NewTable))
			buf.WriteString("\n")
		}
		return
	}
	for _, tableID := range ordered {
		if tableDiff, ok := tableMap[tableID]; ok {
			buf.WriteString(generateCreateTable(tableDiff.SchemaName, tableDiff.TableName, tableDiff.NewTable))
			buf.WriteString("\n")
		}
	}
}

func createViewsInOrder(diff *SchemaDiff, buf *strings.Builder) {
	graph := newGraph()
	viewMap := make(map[string]*ViewDiffEntry)
	var viewsToProcess []*ViewDiffEntry
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffCreate || viewDiff.Action == DiffAlter {
			viewsToProcess = append(viewsToProcess, viewDiff)
		}
	}
	slices.SortFunc(viewsToProcess, func(a, b *ViewDiffEntry) int {
		return strings.Compare(objectID(a.SchemaName, a.ViewName), objectID(b.SchemaName, b.ViewName))
	})
	for _, viewDiff := range viewsToProcess {
		viewID := objectID(viewDiff.SchemaName, viewDiff.ViewName)
		graph.addNode(viewID)
		viewMap[viewID] = viewDiff
	}
	for _, viewDiff := range viewsToProcess {
		viewID := objectID(viewDiff.SchemaName, viewDiff.ViewName)
		if viewDiff.NewView == nil || viewDiff.NewView.Definition == "" {
			continue
		}
		deps, err := getViewDependencies(viewDiff.NewView.Definition, viewDiff.SchemaName)
		if err != nil {
			continue
		}
		for _, dep := range deps {
			if _, ok := viewMap[dep]; ok {
				graph.addEdge(dep, viewID)
			}
		}
	}
	writeView := func(viewDiff *ViewDiffEntry) {
		buf.WriteString(viewDiff.NewView.Definition)
		if !strings.HasSuffix(strings.TrimSpace(viewDiff.NewView.Definition), ";") {
			buf.WriteString(";")
		}
		buf.WriteString("\nGO\n")
		if viewDiff.NewView.Comment != "" {
			buf.WriteString(generateViewCommentSQL("ADD", viewDiff.SchemaName, viewDiff.ViewName, viewDiff.NewView.Comment))
			buf.WriteString(";\nGO\n")
		}
	}
	ordered, err := graph.topologicalSort()
	if err != nil {
		for _, viewDiff := range viewsToProcess {
			writeView(viewDiff)
		}
		return
	}
	for _, viewID := range ordered {
		if viewDiff, ok := viewMap[viewID]; ok {
			writeView(viewDiff)
		}
	}
}

func hasCreateOrAlterTables(diff *SchemaDiff) bool {
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffCreate || tableDiff.Action == DiffAlter {
			return true
		}
	}
	return false
}

func hasCreateViews(diff *SchemaDiff) bool {
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffCreate || viewDiff.Action == DiffAlter {
			return true
		}
	}
	return false
}

func hasCreateFunctions(diff *SchemaDiff) bool {
	for _, funcDiff := range diff.Functions {
		if funcDiff.Action == DiffCreate {
			return true
		}
	}
	return false
}

func hasCreateProcedures(diff *SchemaDiff) bool {
	for _, procDiff := range diff.Procedures {
		if procDiff.Action == DiffCreate || procDiff.Action == DiffAlter {
			return true
		}
	}
	return false
}

func generateSpatialIndexDDL(index *metadata.IndexMetadata, schemaName, tableName string) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "CREATE SPATIAL INDEX [%s] ON [%s].[%s] (", index.Name, schemaName, tableName)
	for i, expr := range index.Expressions {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "[%s]", expr)
	}
	buf.WriteString(")")
	if index.SpatialConfig == nil || index.SpatialConfig.Tessellation == nil {
		buf.WriteString("\n/* WARNING: Spatial index configuration may be incomplete. ")
		buf.WriteString("Please verify BOUNDING_BOX, GRIDS, and other parameters manually. */")
		return buf.String()
	}
	tessellation := index.SpatialConfig.Tessellation
	if tessellation.Scheme != "" {
		buf.WriteString("\nUSING ")
		buf.WriteString(tessellation.Scheme)
	}
	var withParams []string
	withParams = append(withParams, buildMigrationTessellationParams(tessellation)...)
	if index.SpatialConfig.Storage != nil {
		withParams = append(withParams, buildMigrationStorageParams(index.SpatialConfig.Storage)...)
	}
	if len(withParams) > 0 {
		buf.WriteString("\nWITH (\n    ")
		buf.WriteString(strings.Join(withParams, ",\n    "))
		buf.WriteString("\n)")
	}
	return buf.String()
}

func buildMigrationTessellationParams(tessellation *metadata.TessellationConfig) []string {
	var params []string
	if tessellation.BoundingBox != nil {
		bbox := tessellation.BoundingBox
		params = append(params, fmt.Sprintf("BOUNDING_BOX = (%g, %g, %g, %g)", bbox.Xmin, bbox.Ymin, bbox.Xmax, bbox.Ymax))
	}
	if len(tessellation.GridLevels) > 0 {
		var gridParts []string
		for _, level := range tessellation.GridLevels {
			gridParts = append(gridParts, fmt.Sprintf("LEVEL_%d = %s", level.Level, level.Density))
		}
		if len(gridParts) > 0 {
			params = append(params, fmt.Sprintf("GRIDS = (%s)", strings.Join(gridParts, ", ")))
		}
	}
	if tessellation.CellsPerObject > 0 {
		params = append(params, fmt.Sprintf("CELLS_PER_OBJECT = %d", tessellation.CellsPerObject))
	}
	return params
}

func buildMigrationStorageParams(storage *metadata.StorageConfig) []string {
	var params []string
	if storage.PadIndex {
		params = append(params, "PAD_INDEX = ON")
	}
	if storage.Fillfactor > 0 {
		params = append(params, fmt.Sprintf("FILLFACTOR = %d", storage.Fillfactor))
	}
	if storage.SortInTempdb != "" {
		params = append(params, fmt.Sprintf("SORT_IN_TEMPDB = %s", storage.SortInTempdb))
	}
	if storage.DropExisting {
		params = append(params, "DROP_EXISTING = ON")
	}
	if storage.Online {
		params = append(params, "ONLINE = ON")
	}
	if storage.AllowRowLocks {
		params = append(params, "ALLOW_ROW_LOCKS = ON")
	}
	if storage.AllowPageLocks {
		params = append(params, "ALLOW_PAGE_LOCKS = ON")
	}
	if storage.Maxdop > 0 {
		params = append(params, fmt.Sprintf("MAXDOP = %d", storage.Maxdop))
	}
	if storage.DataCompression != "" && storage.DataCompression != "NONE" {
		params = append(params, fmt.Sprintf("DATA_COMPRESSION = %s", storage.DataCompression))
	}
	return params
}
