package catalog

import (
	"fmt"
	"regexp"
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
	if view := blockingSchemaBoundView(from, to, diff); view != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: column change under schema-bound view %s; drop and recreate the view manually.\n", view)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	if view := blockingSchemaBoundDependentView(from, to, diff); view != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: replacing view %s under schema-bound dependent view %s; rebuild the dependent view manually.\n", view, view)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	if function := blockingFunctionCheck(from, to, diff); function != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: function %s is used by a CHECK constraint; rebuild the constraint manually.\n", function)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	if column := blockingComputedColumn(diff); column != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: column change under computed column %s; drop and recreate the computed column manually.\n", column)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	var buf strings.Builder
	inbound := inboundForeignKeys(from, diff)
	for _, ref := range inbound {
		fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\nGO\n", quote(ref.schema), quote(ref.table), quote(ref.fk.Name))
	}

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
					fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\nGO\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(fkDiff.OldForeignKey.Name))
				}
			}
		case DiffDrop:
			for _, fk := range tableDiff.OldTable.ForeignKeys {
				fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(fk.Name))
			}
		}
	}

	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffAlter {
			continue
		}
		for _, triggerDiff := range tableDiff.Triggers {
			if triggerDiff.Action == DiffDrop || triggerDiff.Action == DiffAlter {
				fmt.Fprintf(&buf, "DROP TRIGGER %s.%s;\nGO\n", quote(tableDiff.SchemaName), quote(triggerDiff.OldTrigger.Name))
			}
		}
	}

	dropViewsInOrder(diff, &buf)

	for _, procDiff := range diff.Procedures {
		if procDiff.Action == DiffDrop {
			fmt.Fprintf(&buf, "DROP PROCEDURE %s.%s;\nGO\n", quote(procDiff.SchemaName), quote(procDiff.ProcedureName))
		}
	}

	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffAlter {
			continue
		}
		for _, checkDiff := range tableDiff.Checks {
			if checkDiff.Action == DiffDrop {
				fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(checkDiff.OldCheckConstraint.Name))
			}
		}
		for _, indexDiff := range tableDiff.Indexes {
			if indexDiff.Action != DiffDrop {
				continue
			}
			if indexDiff.OldIndex.IsConstraint {
				fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(indexDiff.OldIndex.Name))
			} else {
				fmt.Fprintf(&buf, "DROP INDEX %s ON %s.%s;\n", quote(indexDiff.OldIndex.Name), quote(tableDiff.SchemaName), quote(tableDiff.TableName))
			}
		}

		hasDroppedIndex := false
		for _, indexDiff := range tableDiff.Indexes {
			if indexDiff.Action == DiffDrop {
				hasDroppedIndex = true
				break
			}
		}
		droppedColumns := make([]*ColumnDiffEntry, 0)
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action == DiffDrop {
				droppedColumns = append(droppedColumns, colDiff)
			}
		}
		slices.SortStableFunc(droppedColumns, func(a, b *ColumnDiffEntry) int {
			aComputed := a.OldColumn != nil && isComputedColumnType(a.OldColumn.Type)
			bComputed := b.OldColumn != nil && isComputedColumnType(b.OldColumn.Type)
			if aComputed == bComputed {
				return 0
			}
			if aComputed {
				return -1
			}
			return 1
		})
		for i, colDiff := range droppedColumns {
			_ = i
			if i == 0 && hasDroppedIndex {
				buf.WriteString("GO\n")
			}
			if getColumnDefaultValue(colDiff.OldColumn) != "" && colDiff.OldColumn.DefaultConstraintName != "" {
				fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(colDiff.OldColumn.DefaultConstraintName))
			}
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP COLUMN %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(colDiff.OldColumn.Name))
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
		fmt.Fprintf(&buf, "DROP TABLE %s.%s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName))
	}

	for _, schemaDiff := range diff.Schemas {
		if schemaDiff.Action != DiffDrop || schemaDiff.OldSchema == nil {
			continue
		}
		for _, table := range schemaDiff.OldSchema.Tables {
			if table.GetSkipDump() {
				continue
			}
			alreadyHandled := false
			for _, tableDiff := range diff.Tables {
				if tableDiff.Action == DiffDrop && tableDiff.SchemaName == schemaDiff.SchemaName && tableDiff.TableName == table.Name {
					alreadyHandled = true
					break
				}
			}
			if !alreadyHandled {
				fmt.Fprintf(&buf, "DROP TABLE %s.%s;\n", quote(schemaDiff.SchemaName), quote(table.Name))
			}
		}
	}

	droppedDefaults := make(map[string]bool)
	droppedFunctionDefaults, restoreFunctionDefaults := dropFunctionDefaults(from, diff, &buf)
	for name := range droppedFunctionDefaults {
		droppedDefaults[name] = true
	}
	for _, funcDiff := range orderedFunctionDrops(diff.Functions) {
		if funcDiff.Action != DiffDrop {
			continue
		}
		if functionReplacementUsedByComputed(diff, funcDiff) {
			fmt.Fprintf(&buf, "-- UNSUPPORTED: function %s.%s is used by a computed column; left unchanged.\n", quote(funcDiff.SchemaName), quote(funcDiff.FunctionName))
			continue
		}
		fmt.Fprintf(&buf, "DROP FUNCTION %s.%s;\nGO\n", quote(funcDiff.SchemaName), quote(funcDiff.FunctionName))
	}
	for name := range dropSequenceDefaults(from, diff, &buf) {
		droppedDefaults[name] = true
	}
	for _, sequenceDiff := range diff.Sequences {
		if sequenceDiff.Action == DiffDrop {
			fmt.Fprintf(&buf, "DROP SEQUENCE %s.%s;\nGO\n", quote(sequenceDiff.SchemaName), quote(sequenceDiff.SequenceName))
		}
	}

	for _, synonymDiff := range diff.Synonyms {
		if synonymDiff.Action == DiffDrop {
			fmt.Fprintf(&buf, "DROP SYNONYM %s.%s;\nGO\n", quote(synonymDiff.SchemaName), quote(synonymDiff.SynonymName))
		}
	}
	for _, schemaDiff := range diff.Schemas {
		if schemaDiff.Action != DiffDrop || strings.ToLower(schemaDiff.SchemaName) == defaultSchema {
			continue
		}
		if schemaHasSkipDump(schemaDiff.OldSchema) {
			fmt.Fprintf(&buf, "-- UNSUPPORTED: schema %s still holds objects marked SkipDump; left unchanged.\n", quote(schemaDiff.SchemaName))
			continue
		}
		fmt.Fprintf(&buf, "DROP SCHEMA %s;\nGO\n", quote(schemaDiff.SchemaName))
	}

	dropPhaseHasContent := buf.Len() > 0
	createPhaseWillHaveContent := len(schemasToCreate) > 0 ||
		hasCreateOrAlterTables(diff) || hasCreateViews(diff) ||
		hasCreateFunctions(diff) || hasCreateProcedures(diff)
	if dropPhaseHasContent && createPhaseWillHaveContent {
		buf.WriteString("\n")
	}

	for _, schemaName := range schemasToCreate {
		fmt.Fprintf(&buf, "CREATE SCHEMA %s;\nGO\n", quote(schemaName))
	}
	for _, synonymDiff := range diff.Synonyms {
		if synonymDiff.Action == DiffCreate || synonymDiff.Action == DiffAlter {
			buf.WriteString(generateCreateSynonym(synonymDiff.SchemaName, synonymDiff.NewSynonym))
			buf.WriteString("\nGO\n")
		}
	}
	if len(schemasToCreate) > 0 && (hasCreateOrAlterTables(diff) || hasCreateViews(diff) || hasCreateFunctions(diff) || hasCreateProcedures(diff)) {
		buf.WriteString("\n")
	}

	for _, sequenceDiff := range diff.Sequences {
		switch sequenceDiff.Action {
		case DiffCreate:
			buf.WriteString(generateCreateSequence(sequenceDiff.SchemaName, sequenceDiff.SequenceName, sequenceDiff.NewSequence))
			buf.WriteString("\nGO\n")
			if sequenceDiff.NewSequence.Comment != "" {
				buf.WriteString(generateSequenceCommentSQL("ADD", sequenceDiff.SchemaName, sequenceDiff.SequenceName, sequenceDiff.NewSequence.Comment))
				buf.WriteString(";\nGO\n")
			}
		case DiffAlter:
			buf.WriteString(generateAlterSequence(sequenceDiff))
		}
	}

	tableFunctions := migrationTableFunctions(diff)
	for _, funcDiff := range orderedFunctionDiffs(diff.Functions) {
		if tableFunctions[funcDiff.SchemaName+"."+funcDiff.FunctionName] && !functionReplacementUsedByComputed(diff, funcDiff) {
			writeFunctionMigration(&buf, funcDiff)
		}
	}

	createTablesInOrder(diff, &buf)

	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffAlter {
			alterTableSQL := generateAlterTable(tableDiff, droppedDefaults)
			buf.WriteString(alterTableSQL)
			if alterTableSQL != "" {
				buf.WriteString("\n")
			}
		}
	}

	for _, ref := range inbound {
		if !ref.recreate {
			continue
		}
		buf.WriteString(generateAddForeignKey(ref.schema, ref.table, ref.fk))
		buf.WriteString(";\n")
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

	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffAlter {
			continue
		}
		for _, fkDiff := range tableDiff.ForeignKeys {
			if fkDiff.Action == DiffCreate {
				buf.WriteString(generateAddForeignKey(tableDiff.SchemaName, tableDiff.TableName, fkDiff.NewForeignKey))
				buf.WriteString(";\n")
			}
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

	if hasCreateOrAlterTables(diff) && hasTriggerCreations(diff) {
		buf.WriteString("\nGO\n")
	}
	for _, tableDiff := range diff.Tables {
		switch tableDiff.Action {
		case DiffCreate:
			if tableDiff.NewTable == nil {
				continue
			}
			for _, trigger := range tableDiff.NewTable.Triggers {
				writeTriggerMigration(&buf, trigger)
			}
		case DiffAlter:
			for _, triggerDiff := range tableDiff.Triggers {
				if triggerDiff.Action == DiffCreate || triggerDiff.Action == DiffAlter {
					writeTriggerMigration(&buf, triggerDiff.NewTrigger)
				}
			}
		}
	}

	earlyFunctions := migrationEarlyFunctions(diff)
	for name := range tableFunctions {
		delete(earlyFunctions, name)
	}
	for _, funcDiff := range orderedFunctionDiffs(diff.Functions) {
		if !earlyFunctions[funcDiff.SchemaName+"."+funcDiff.FunctionName] || functionReplacementUsedByComputed(diff, funcDiff) {
			continue
		}
		writeFunctionMigration(&buf, funcDiff)
	}

	createViewsInOrder(diff, &buf)

	for _, funcDiff := range orderedFunctionDiffs(diff.Functions) {
		if funcDiff.Action != DiffCreate {
			continue
		}
		if earlyFunctions[funcDiff.SchemaName+"."+funcDiff.FunctionName] || tableFunctions[funcDiff.SchemaName+"."+funcDiff.FunctionName] || functionReplacementUsedByComputed(diff, funcDiff) {
			continue
		}
		buf.WriteString(creationDefinition(funcDiff.NewFunction.Definition))
		if !strings.HasSuffix(strings.TrimSpace(funcDiff.NewFunction.Definition), ";") {
			buf.WriteString(";")
		}
		buf.WriteString("\nGO\n")
		if funcDiff.NewFunction.Comment != "" {
			buf.WriteString(generateFunctionCommentSQL("ADD", funcDiff.SchemaName, funcDiff.FunctionName, funcDiff.NewFunction.Comment))
			buf.WriteString(";\nGO\n")
		}
	}

	restoredDefault := false
	for _, ref := range restoreFunctionDefaults {
		created := false
		for _, funcDiff := range diff.Functions {
			if funcDiff.Action == DiffCreate && strings.EqualFold(funcDiff.FunctionName, ref.functionName) {
				created = true
			}
		}
		if !created || columnDefaultChanged(diff, ref.schema, ref.table, ref.column) {
			continue
		}
		fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s DEFAULT %s FOR %s;\n",
			quote(ref.schema), quote(ref.table), quote(ref.constraint), ref.expression, quote(ref.column))
		restoredDefault = true
	}
	if restoredDefault {
		buf.WriteString("GO\n")
	}

	for _, procDiff := range diff.Procedures {
		switch procDiff.Action {
		case DiffCreate:
			buf.WriteString(creationDefinition(procDiff.NewProcedure.Definition))
			if !strings.HasSuffix(strings.TrimSpace(procDiff.NewProcedure.Definition), ";") {
				buf.WriteString(";")
			}
			buf.WriteString("\nGO\n")
		case DiffAlter:
			fmt.Fprintf(&buf, "DROP PROCEDURE %s.%s;\nGO\n", quote(procDiff.SchemaName), quote(procDiff.ProcedureName))
			buf.WriteString(creationDefinition(procDiff.NewProcedure.Definition))
			if !strings.HasSuffix(strings.TrimSpace(procDiff.NewProcedure.Definition), ";") {
				buf.WriteString(";")
			}
			buf.WriteString("\nGO\n")
		}
	}

	text := buf.String()
	var ops []MigrationOp
	if text != "" {
		ops = append(ops, MigrationOp{Type: OpScript, SQL: text})
	}
	return &MigrationPlan{Ops: ops}
}

func appendMigrationDirection(buf *strings.Builder, descending []bool, index int) {
	if index < len(descending) && descending[index] {
		buf.WriteString(" DESC")
	}
}

// generateAlterSequence emits the options ALTER SEQUENCE can change. A type
// change is not alterable in SQL Server, so it is recreated instead.
func generateAlterSequence(diff *SequenceDiffEntry) string {
	oldSequence, newSequence := diff.OldSequence, diff.NewSequence
	var buf strings.Builder
	if !typesEqual(oldSequence.DataType, newSequence.DataType) {
		fmt.Fprintf(&buf, "-- UNSUPPORTED: sequence %s.%s data type change requires recreating it; current position would be lost.\n",
			quote(diff.SchemaName), quote(diff.SequenceName))
		return buf.String()
	}
	fmt.Fprintf(&buf, "ALTER SEQUENCE %s.%s", quote(diff.SchemaName), quote(diff.SequenceName))
	if newSequence.Start != "" && newSequence.Start != oldSequence.Start {
		fmt.Fprintf(&buf, " RESTART WITH %s", newSequence.Start)
	}
	if newSequence.Increment != "" {
		fmt.Fprintf(&buf, " INCREMENT BY %s", newSequence.Increment)
	} else if oldSequence.Increment != "" {
		buf.WriteString(" INCREMENT BY 1")
	}
	if newSequence.MinValue != "" {
		fmt.Fprintf(&buf, " MINVALUE %s", newSequence.MinValue)
	} else if oldSequence.MinValue != "" {
		buf.WriteString(" NO MINVALUE")
	}
	if newSequence.MaxValue != "" {
		fmt.Fprintf(&buf, " MAXVALUE %s", newSequence.MaxValue)
	} else if oldSequence.MaxValue != "" {
		buf.WriteString(" NO MAXVALUE")
	}
	if newSequence.Cycle {
		buf.WriteString(" CYCLE")
	} else {
		buf.WriteString(" NO CYCLE")
	}
	if newSequence.CacheSize != "" {
		if newSequence.CacheSize == "NO CACHE" {
			buf.WriteString(" NO CACHE")
		} else {
			fmt.Fprintf(&buf, " CACHE %s", newSequence.CacheSize)
		}
	}
	buf.WriteString(";\nGO\n")
	return buf.String()
}

func generateCreateSynonym(schemaName string, synonym *metadata.SynonymMetadata) string {
	return fmt.Sprintf("CREATE SYNONYM %s.%s FOR %s;", quote(schemaName), quote(synonym.Name), synonymTarget(synonym))
}

func generateCreateSequence(schemaName, sequenceName string, sequence *metadata.SequenceMetadata) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "CREATE SEQUENCE %s.%s", quote(schemaName), quote(sequenceName))
	if sequence.DataType != "" {
		fmt.Fprintf(&buf, " AS %s", sequence.DataType)
	}
	if sequence.Start != "" {
		fmt.Fprintf(&buf, " START WITH %s", sequence.Start)
	}
	if sequence.Increment != "" {
		fmt.Fprintf(&buf, " INCREMENT BY %s", sequence.Increment)
	}
	if sequence.MinValue != "" {
		fmt.Fprintf(&buf, " MINVALUE %s", sequence.MinValue)
	}
	if sequence.MaxValue != "" {
		fmt.Fprintf(&buf, " MAXVALUE %s", sequence.MaxValue)
	}
	if sequence.Cycle {
		buf.WriteString(" CYCLE")
	}
	if sequence.CacheSize != "" {
		if sequence.CacheSize == "NO CACHE" {
			buf.WriteString(" NO CACHE")
		} else {
			fmt.Fprintf(&buf, " CACHE %s", sequence.CacheSize)
		}
	}
	buf.WriteString(";")
	return buf.String()
}

// foreignKeyRef is an unchanged foreign key that references a key being
// replaced, so it must be dropped before and recreated after the replacement.
type foreignKeyRef struct {
	schema   string
	table    string
	fk       *metadata.ForeignKeyMetadata
	recreate bool
}

func inboundForeignKeys(from *Catalog, diff *SchemaDiff) []foreignKeyRef {
	if from == nil {
		return nil
	}
	replaced := make(map[string]bool)
	droppedTables := make(map[string]bool)
	explicit := make(map[string]bool)
	for _, tableDiff := range diff.Tables {
		tableID := objectID(tableDiff.SchemaName, tableDiff.TableName)
		switch tableDiff.Action {
		case DiffDrop:
			droppedTables[tableID] = true
		case DiffAlter:
			droppedKey := false
			for _, indexDiff := range tableDiff.Indexes {
				if indexDiff.Action == DiffDrop && (indexDiff.OldIndex.IsConstraint || indexDiff.OldIndex.Unique) {
					droppedKey = true
				}
			}
			for _, colDiff := range tableDiff.Columns {
				if colDiff.Action == DiffAlter && (!typesEqual(colDiff.OldColumn.Type, colDiff.NewColumn.Type) || colDiff.OldColumn.Collation != colDiff.NewColumn.Collation) {
					droppedKey = true
				}
			}
			if droppedKey {
				replaced[tableID] = true
			}
			for _, fkDiff := range tableDiff.ForeignKeys {
				name := fkDiffName(fkDiff)
				explicit[tableID+"."+name] = true
			}
		}
	}
	if len(replaced) == 0 && len(droppedTables) == 0 {
		return nil
	}
	var refs []foreignKeyRef
	for _, schema := range from.Metadata().Schemas {
		for _, table := range schema.Tables {
			tableID := objectID(schema.Name, table.Name)
			if droppedTables[tableID] {
				continue
			}
			for _, fk := range table.ForeignKeys {
				referenced := objectID(fk.ReferencedSchema, fk.ReferencedTable)
				if !replaced[referenced] && !droppedTables[referenced] {
					continue
				}
				if explicit[tableID+"."+fk.Name] {
					continue
				}
				refs = append(refs, foreignKeyRef{
					schema:   schema.Name,
					table:    table.Name,
					fk:       fk,
					recreate: replaced[referenced],
				})
			}
		}
	}
	return refs
}

func generateCreateTable(schemaName, tableName string, table *metadata.TableMetadata) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "CREATE TABLE %s.%s (\n", quote(schemaName), quote(tableName))
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
		fmt.Fprintf(&buf, "  CONSTRAINT %s", quote(idx.Name))
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
			fmt.Fprintf(&buf, "%s", quote(expr))
			appendMigrationDirection(&buf, idx.Descending, j)
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
		fmt.Fprintf(&buf, "  CONSTRAINT %s CHECK %s", quote(check.Name), checkExpression(check.Expression))
		if i < len(table.CheckConstraints)-1 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}
	buf.WriteString(");")

	for _, idx := range orderedIndexes(table.Indexes) {
		if idx.IsConstraint {
			continue
		}
		buf.WriteString("\n")
		buf.WriteString(generateCreateIndex(schemaName, tableName, idx))
		buf.WriteString(";")
	}
	return buf.String()
}

func generateAlterTable(tableDiff *TableDiffEntry, droppedDefaults map[string]bool) string {
	var buf strings.Builder
	// A computed column cannot be changed with ALTER COLUMN. Drop the indexes
	// that use it, drop the column, add the new expression, and rebuild them.
	var rebuiltIndexes []string
	var rebuiltChecks []string
	for _, colDiff := range tableDiff.Columns {
		if colDiff.Action != DiffAlter || !columnRequiresRebuild(colDiff) {
			continue
		}
		for _, index := range tableDiff.OldTable.Indexes {
			if !indexUsesColumn(index, colDiff.OldColumn.Name) {
				continue
			}
			if indexChangeScheduled(tableDiff, index.Name) || containsString(rebuiltIndexes, index.Name) {
				continue
			}
			if index.IsConstraint {
				fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(index.Name))
			} else {
				fmt.Fprintf(&buf, "DROP INDEX %s ON %s.%s;\n", quote(index.Name), quote(tableDiff.SchemaName), quote(tableDiff.TableName))
			}
			rebuiltIndexes = append(rebuiltIndexes, index.Name)
		}
		if colDiff.OldColumn.Default != "" && colDiff.OldColumn.DefaultConstraintName != "" {
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(colDiff.OldColumn.DefaultConstraintName))
		}
		for _, check := range tableDiff.OldTable.CheckConstraints {
			if !strings.Contains(strings.ToLower(check.Expression), strings.ToLower(colDiff.OldColumn.Name)) {
				continue
			}
			if checkChangeScheduled(tableDiff, check.Name) || containsString(rebuiltChecks, check.Name) {
				continue
			}
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(check.Name))
			rebuiltChecks = append(rebuiltChecks, check.Name)
		}
		fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP COLUMN %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(colDiff.OldColumn.Name))
	}
	for _, colDiff := range tableDiff.Columns {
		if colDiff.Action != DiffAlter || columnRequiresRebuild(colDiff) ||
			(typesEqual(colDiff.OldColumn.Type, colDiff.NewColumn.Type) && colDiff.OldColumn.Collation == colDiff.NewColumn.Collation) {
			continue
		}
		for _, index := range tableDiff.OldTable.Indexes {
			if !indexUsesColumn(index, colDiff.OldColumn.Name) || indexChangeScheduled(tableDiff, index.Name) || containsString(rebuiltIndexes, index.Name) {
				continue
			}
			if index.IsConstraint {
				fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(index.Name))
			} else {
				fmt.Fprintf(&buf, "DROP INDEX %s ON %s.%s;\n", quote(index.Name), quote(tableDiff.SchemaName), quote(tableDiff.TableName))
			}
			rebuiltIndexes = append(rebuiltIndexes, index.Name)
		}
		for _, check := range tableDiff.OldTable.CheckConstraints {
			if !strings.Contains(strings.ToLower(check.Expression), strings.ToLower(colDiff.OldColumn.Name)) ||
				checkChangeScheduled(tableDiff, check.Name) || containsString(rebuiltChecks, check.Name) {
				continue
			}
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(check.Name))
			rebuiltChecks = append(rebuiltChecks, check.Name)
		}
	}
	for _, colDiff := range tableDiff.Columns {
		if colDiff.Action == DiffCreate || (colDiff.Action == DiffAlter && columnRequiresRebuild(colDiff)) {
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), generateColumnDefinition(colDiff.NewColumn))
		}
	}
	for _, colDiff := range tableDiff.Columns {
		if colDiff.Action == DiffAlter && !columnRequiresRebuild(colDiff) {
			buf.WriteString(generateAlterColumn(tableDiff.SchemaName, tableDiff.TableName, colDiff, droppedDefaults))
		}
	}
	for _, name := range rebuiltChecks {
		for _, check := range tableDiff.NewTable.CheckConstraints {
			if check.Name == name && !checkChangeScheduled(tableDiff, name) {
				fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s CHECK %s;\n",
					quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(check.Name), checkExpression(check.Expression))
			}
		}
	}
	for _, name := range rebuiltIndexes {
		for _, index := range tableDiff.NewTable.Indexes {
			if index.Name != name || indexChangeScheduled(tableDiff, name) {
				continue
			}
			if index.IsConstraint {
				fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(index.Name))
				switch {
				case index.Primary:
					buf.WriteString(" PRIMARY KEY")
				case index.Unique:
					buf.WriteString(" UNIQUE")
				}
				switch index.Type {
				case "CLUSTERED":
					buf.WriteString(" CLUSTERED")
				case "NONCLUSTERED":
					buf.WriteString(" NONCLUSTERED")
				}
				buf.WriteString(" (")
				for j, expr := range index.Expressions {
					if j > 0 {
						buf.WriteString(", ")
					}
					buf.WriteString(quote(expr))
					appendMigrationDirection(&buf, index.Descending, j)
				}
				buf.WriteString(");\n")
			} else {
				buf.WriteString(generateCreateIndex(tableDiff.SchemaName, tableDiff.TableName, index))
				buf.WriteString(";\n")
			}
		}
	}
	var secondaryXML []*metadata.IndexMetadata
	for _, indexDiff := range tableDiff.Indexes {
		if indexDiff.Action != DiffCreate {
			continue
		}
		if indexDiff.NewIndex.Type == "XML" && strings.HasPrefix(indexDiff.NewIndex.Definition, "USING XML INDEX ") {
			secondaryXML = append(secondaryXML, indexDiff.NewIndex)
			continue
		}
		if indexDiff.NewIndex.IsConstraint {
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(indexDiff.NewIndex.Name))
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
				fmt.Fprintf(&buf, "%s", quote(expr))
				appendMigrationDirection(&buf, indexDiff.NewIndex.Descending, j)
			}
			buf.WriteString(");\n")
		} else {
			buf.WriteString(generateCreateIndex(tableDiff.SchemaName, tableDiff.TableName, indexDiff.NewIndex))
			buf.WriteString(";\n")
		}
	}
	for _, index := range secondaryXML {
		buf.WriteString(generateCreateIndex(tableDiff.SchemaName, tableDiff.TableName, index))
		buf.WriteString(";\n")
	}
	for _, checkDiff := range tableDiff.Checks {
		if checkDiff.Action == DiffCreate {
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s CHECK %s;\n",
				quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(checkDiff.NewCheckConstraint.Name), checkExpression(checkDiff.NewCheckConstraint.Expression))
		}
	}
	return buf.String()
}

func blockingFunctionCheck(from, to *Catalog, diff *SchemaDiff) string {
	replaced := make(map[string]bool)
	for _, funcDiff := range diff.Functions {
		if funcDiff.Action == DiffDrop || funcDiff.Action == DiffAlter {
			replaced[strings.ToLower(funcDiff.FunctionName)] = true
		}
	}
	if len(replaced) == 0 {
		return ""
	}
	for _, catalog := range []*Catalog{from, to} {
		if catalog == nil {
			continue
		}
		for _, schema := range catalog.Metadata().Schemas {
			for _, function := range schema.Functions {
				if !strings.Contains(strings.ToUpper(function.Definition), "SCHEMABINDING") || replaced[strings.ToLower(function.Name)] {
					continue
				}
				for name := range replaced {
					if strings.Contains(strings.ToLower(function.Definition), name) {
						return quote(schema.Name) + "." + quote(function.Name)
					}
				}
			}
		}
	}
	for _, tableDiff := range diff.Tables {
		table := tableDiff.OldTable
		if table == nil {
			table = tableDiff.NewTable
		}
		if table == nil {
			continue
		}
		for _, check := range table.CheckConstraints {
			for name := range replaced {
				if strings.Contains(strings.ToLower(check.Expression), name) {
					return quote(tableDiff.SchemaName) + "." + quote(tableDiff.TableName) + "." + quote(check.Name)
				}
			}
		}
	}
	return ""
}

func blockingComputedColumn(diff *SchemaDiff) string {
	for _, tableDiff := range diff.Tables {
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action != DiffAlter || typesEqual(colDiff.OldColumn.Type, colDiff.NewColumn.Type) {
				continue
			}
			for _, column := range tableDiff.OldTable.Columns {
				if isComputedColumnType(column.Type) && strings.Contains(strings.ToLower(column.Type), strings.ToLower(colDiff.OldColumn.Name)) {
					return quote(tableDiff.SchemaName) + "." + quote(tableDiff.TableName) + "." + quote(column.Name)
				}
			}
		}
	}
	return ""
}

func blockingSchemaBoundDependentView(from, to *Catalog, diff *SchemaDiff) string {
	changed := make(map[string]bool)
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffAlter || viewDiff.Action == DiffDrop {
			changed[viewDiff.SchemaName+"."+viewDiff.ViewName] = true
		}
	}
	if len(changed) == 0 {
		return ""
	}
	for _, catalog := range []*Catalog{from, to} {
		if catalog == nil {
			continue
		}
		for _, schema := range catalog.Metadata().Schemas {
			for _, view := range schema.Views {
				if !strings.Contains(strings.ToUpper(view.Definition), "SCHEMABINDING") {
					continue
				}
				for name := range changed {
					parts := strings.SplitN(name, ".", 2)
					if len(parts) == 2 && !strings.EqualFold(view.Name, parts[1]) && strings.Contains(strings.ToLower(view.Definition), strings.ToLower(parts[1])) {
						return quote(schema.Name) + "." + quote(view.Name)
					}
				}
			}
		}
	}
	return ""
}

func blockingSchemaBoundView(from, to *Catalog, diff *SchemaDiff) string {
	for _, tableDiff := range diff.Tables {
		changed := false
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action == DiffAlter && !typesEqual(colDiff.OldColumn.Type, colDiff.NewColumn.Type) {
				changed = true
			}
		}
		if !changed {
			continue
		}
		for _, catalog := range []*Catalog{from, to} {
			if catalog == nil {
				continue
			}
			for _, schema := range catalog.Metadata().Schemas {
				for _, view := range schema.Views {
					if strings.Contains(strings.ToUpper(view.Definition), "SCHEMABINDING") &&
						strings.Contains(strings.ToLower(view.Definition), strings.ToLower(tableDiff.TableName)) {
						return quote(schema.Name) + "." + quote(view.Name)
					}
				}
			}
		}
	}
	return ""
}

func dropSequenceDefaults(from *Catalog, diff *SchemaDiff, buf *strings.Builder) map[string]bool {
	removed := make(map[string]bool)
	if from == nil {
		return removed
	}
	droppedSequences := make(map[string]bool)
	for _, sequenceDiff := range diff.Sequences {
		if sequenceDiff.Action == DiffDrop {
			droppedSequences[strings.ToLower(sequenceDiff.SchemaName+"."+sequenceDiff.SequenceName)] = true
		}
	}
	droppedTables := make(map[string]bool)
	droppedColumns := make(map[string]bool)
	for _, tableDiff := range diff.Tables {
		tableID := tableDiff.SchemaName + "." + tableDiff.TableName
		if tableDiff.Action == DiffDrop {
			droppedTables[tableID] = true
		}
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action == DiffDrop && colDiff.OldColumn != nil {
				droppedColumns[tableID+"."+colDiff.OldColumn.Name] = true
			}
		}
	}
	for _, schema := range from.Metadata().Schemas {
		for _, table := range schema.Tables {
			tableID := schema.Name + "." + table.Name
			if droppedTables[tableID] {
				continue
			}
			for _, column := range table.Columns {
				if column.Default == "" || column.DefaultConstraintName == "" || droppedColumns[tableID+"."+column.Name] {
					continue
				}
				if !strings.Contains(strings.ToLower(column.Default), "next value for") {
					continue
				}
				for name := range droppedSequences {
					parts := strings.SplitN(name, ".", 2)
					if len(parts) != 2 {
						continue
					}
					if defaultReferencesSequence(column.Default, parts[0], parts[1]) {
						fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(schema.Name), quote(table.Name), quote(column.DefaultConstraintName))
						removed[schema.Name+"."+column.DefaultConstraintName] = true
						break
					}
				}
			}
		}
	}
	return removed
}

type functionDefaultRef struct {
	schema       string
	table        string
	column       string
	constraint   string
	expression   string
	functionName string
}

func columnDefaultChanged(diff *SchemaDiff, schemaName, tableName, columnName string) bool {
	for _, tableDiff := range diff.Tables {
		if !strings.EqualFold(tableDiff.SchemaName, schemaName) || !strings.EqualFold(tableDiff.TableName, tableName) {
			continue
		}
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action == DiffDrop && colDiff.OldColumn != nil && strings.EqualFold(colDiff.OldColumn.Name, columnName) {
				return true
			}
			if colDiff.Action == DiffAlter && colDiff.NewColumn != nil && strings.EqualFold(colDiff.NewColumn.Name, columnName) {
				return !expressionsEqual(colDiff.OldColumn.Default, colDiff.NewColumn.Default)
			}
		}
	}
	return false
}

func dropFunctionDefaults(from *Catalog, diff *SchemaDiff, buf *strings.Builder) (map[string]bool, []functionDefaultRef) {
	removed := make(map[string]bool)
	var restore []functionDefaultRef
	if from == nil {
		return removed, restore
	}
	droppedFunctions := make(map[string]bool)
	for _, funcDiff := range diff.Functions {
		if funcDiff.Action == DiffDrop {
			droppedFunctions[strings.ToLower(funcDiff.FunctionName)] = true
		}
	}
	droppedTables := make(map[string]bool)
	droppedColumns := make(map[string]bool)
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffDrop {
			droppedTables[tableDiff.SchemaName+"."+tableDiff.TableName] = true
		}
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action == DiffDrop && colDiff.OldColumn != nil {
				droppedColumns[tableDiff.SchemaName+"."+tableDiff.TableName+"."+colDiff.OldColumn.Name] = true
			}
		}
	}
	for _, schema := range from.Metadata().Schemas {
		for _, table := range schema.Tables {
			if droppedTables[schema.Name+"."+table.Name] {
				continue
			}
			for _, column := range table.Columns {
				if column.Default == "" || column.DefaultConstraintName == "" || droppedColumns[schema.Name+"."+table.Name+"."+column.Name] {
					continue
				}
				for name := range droppedFunctions {
					if strings.Contains(strings.ToLower(column.Default), name) {
						fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(schema.Name), quote(table.Name), quote(column.DefaultConstraintName))
						removed[schema.Name+"."+column.DefaultConstraintName] = true
						restore = append(restore, functionDefaultRef{
							schema: schema.Name, table: table.Name, column: column.Name,
							constraint: column.DefaultConstraintName, expression: column.Default,
							functionName: name,
						})
						break
					}
				}
			}
		}
	}
	return removed, restore
}

func defaultReferencesSequence(expression, schemaName, sequenceName string) bool {
	matches := sequenceReference.FindAllStringSubmatch(expression, -1)
	for _, match := range matches {
		parts := strings.Split(match[1], ".")
		for i := range parts {
			parts[i] = strings.Trim(parts[i], "[]")
		}
		if len(parts) == 0 || !strings.EqualFold(parts[len(parts)-1], sequenceName) {
			continue
		}
		if len(parts) == 1 || strings.EqualFold(parts[len(parts)-2], schemaName) {
			return true
		}
	}
	return false
}

var sequenceReference = regexp.MustCompile(`(?i)NEXT\s+VALUE\s+FOR\s+((?:\[[^\]]+\]|[A-Za-z0-9_]+)(?:\.(?:\[[^\]]+\]|[A-Za-z0-9_]+))*)`)

func schemaHasSkipDump(schema *metadata.SchemaMetadata) bool {
	if schema == nil {
		return false
	}
	for _, table := range schema.Tables {
		if table.GetSkipDump() {
			return true
		}
	}
	for _, view := range schema.Views {
		if view.GetSkipDump() {
			return true
		}
	}
	for _, function := range schema.Functions {
		if function.GetSkipDump() {
			return true
		}
	}
	for _, procedure := range schema.Procedures {
		if procedure.SkipDump {
			return true
		}
	}
	for _, sequence := range schema.Sequences {
		if sequence.GetSkipDump() {
			return true
		}
	}
	return false
}

func functionReplacementUsedByComputed(diff *SchemaDiff, funcDiff *FunctionDiffEntry) bool {
	if !functionUsedByComputed(diff, funcDiff.FunctionName) {
		return false
	}
	for _, other := range diff.Functions {
		if other.Action == DiffDrop && other.FunctionName == funcDiff.FunctionName {
			return true
		}
	}
	return false
}

func functionUsedByComputed(diff *SchemaDiff, functionName string) bool {
	for _, tableDiff := range diff.Tables {
		tables := []*metadata.TableMetadata{tableDiff.OldTable, tableDiff.NewTable}
		for _, table := range tables {
			if table == nil {
				continue
			}
			for _, column := range table.Columns {
				if isComputedColumnType(column.Type) && strings.Contains(strings.ToLower(column.Type), strings.ToLower(functionName)) {
					return true
				}
			}
		}
	}
	return false
}

func orderedFunctionDrops(functions []*FunctionDiffEntry) []*FunctionDiffEntry {
	graph := newGraph()
	byName := make(map[string]*FunctionDiffEntry)
	for _, function := range functions {
		if function.Action != DiffDrop {
			continue
		}
		id := function.SchemaName + "." + function.FunctionName
		graph.addNode(id)
		byName[id] = function
	}
	for _, function := range functions {
		if function.Action != DiffDrop || function.OldFunction == nil {
			continue
		}
		id := function.SchemaName + "." + function.FunctionName
		for _, other := range functions {
			if other.Action != DiffDrop || other == function {
				continue
			}
			name := other.FunctionName
			if strings.Contains(strings.ToLower(function.OldFunction.Definition), strings.ToLower(name)) {
				graph.addEdge(id, other.SchemaName+"."+other.FunctionName)
			}
		}
	}
	order, err := graph.topologicalSort()
	if err != nil {
		return functions
	}
	result := make([]*FunctionDiffEntry, 0, len(functions))
	for _, id := range order {
		if function := byName[id]; function != nil {
			result = append(result, function)
		}
	}
	return result
}

func orderedFunctionDiffs(functions []*FunctionDiffEntry) []*FunctionDiffEntry {
	graph := newGraph()
	byName := make(map[string]*FunctionDiffEntry)
	for _, function := range functions {
		if function.Action != DiffCreate {
			continue
		}
		id := function.SchemaName + "." + function.FunctionName
		graph.addNode(id)
		byName[id] = function
	}
	for _, function := range functions {
		if function.Action != DiffCreate || function.NewFunction == nil {
			continue
		}
		id := function.SchemaName + "." + function.FunctionName
		for _, other := range functions {
			if other.Action != DiffCreate || other == function {
				continue
			}
			name := other.FunctionName
			if other.NewFunction != nil && other.NewFunction.Name != "" {
				name = other.NewFunction.Name
			}
			if strings.Contains(strings.ToLower(function.NewFunction.Definition), strings.ToLower(name)) {
				graph.addEdge(other.SchemaName+"."+other.FunctionName, id)
			}
		}
	}
	order, err := graph.topologicalSort()
	if err != nil {
		return functions
	}
	result := make([]*FunctionDiffEntry, 0, len(functions))
	for _, id := range order {
		if function := byName[id]; function != nil {
			result = append(result, function)
		}
	}
	return result
}

func migrationTableFunctions(diff *SchemaDiff) map[string]bool {
	var expressions []string
	for _, tableDiff := range diff.Tables {
		if tableDiff.NewTable == nil {
			continue
		}
		for _, column := range tableDiff.NewTable.Columns {
			expressions = append(expressions, column.Default, column.Type)
		}
		for _, check := range tableDiff.NewTable.CheckConstraints {
			expressions = append(expressions, check.Expression)
		}
	}
	referenced := make(map[string]bool)
	for _, funcDiff := range orderedFunctionDiffs(diff.Functions) {
		if funcDiff.Action != DiffCreate && funcDiff.Action != DiffAlter {
			continue
		}
		name := funcDiff.FunctionName
		if funcDiff.NewFunction != nil && funcDiff.NewFunction.Name != "" {
			name = funcDiff.NewFunction.Name
		}
		for _, expression := range expressions {
			if strings.Contains(strings.ToLower(expression), strings.ToLower(name)) {
				referenced[funcDiff.SchemaName+"."+funcDiff.FunctionName] = true
				break
			}
		}
	}
	return referenced
}

func migrationEarlyFunctions(diff *SchemaDiff) map[string]bool {
	var definitions []string
	for _, viewDiff := range diff.Views {
		if viewDiff.NewView != nil {
			definitions = append(definitions, viewDiff.NewView.Definition)
		}
	}
	referenced := make(map[string]bool)
	for _, funcDiff := range orderedFunctionDiffs(diff.Functions) {
		if funcDiff.Action != DiffCreate && funcDiff.Action != DiffAlter {
			continue
		}
		name := funcDiff.FunctionName
		if funcDiff.NewFunction != nil && funcDiff.NewFunction.Name != "" {
			name = funcDiff.NewFunction.Name
		}
		for _, definition := range definitions {
			if strings.Contains(strings.ToLower(definition), strings.ToLower(name)) {
				referenced[funcDiff.SchemaName+"."+funcDiff.FunctionName] = true
				break
			}
		}
	}
	return referenced
}

func writeFunctionMigration(buf *strings.Builder, funcDiff *FunctionDiffEntry) {
	if funcDiff.NewFunction == nil {
		return
	}
	buf.WriteString(creationDefinition(funcDiff.NewFunction.Definition))
	if !strings.HasSuffix(strings.TrimSpace(funcDiff.NewFunction.Definition), ";") {
		buf.WriteString(";")
	}
	buf.WriteString("\nGO\n")
	if funcDiff.NewFunction.Comment != "" {
		buf.WriteString(generateFunctionCommentSQL("ADD", funcDiff.SchemaName, funcDiff.FunctionName, funcDiff.NewFunction.Comment))
		buf.WriteString(";\nGO\n")
	}
}

func hasTriggerCreations(diff *SchemaDiff) bool {
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffCreate && tableDiff.NewTable != nil && len(tableDiff.NewTable.Triggers) > 0 {
			return true
		}
		if tableDiff.Action == DiffAlter {
			for _, triggerDiff := range tableDiff.Triggers {
				if triggerDiff.Action == DiffCreate || triggerDiff.Action == DiffAlter {
					return true
				}
			}
		}
	}
	return false
}

func containsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func writeTriggerMigration(buf *strings.Builder, trigger *metadata.TriggerMetadata) {
	if trigger == nil || trigger.Body == "" {
		return
	}
	buf.WriteString(strings.TrimRight(strings.TrimSpace(creationDefinition(trigger.Body)), ";"))
	buf.WriteString("\nGO\n")
}

func columnRequiresRebuild(diff *ColumnDiffEntry) bool {
	if diff.OldColumn == nil || diff.NewColumn == nil {
		return false
	}
	return computedColumnChanged(diff)
}

func computedColumnChanged(diff *ColumnDiffEntry) bool {
	if diff.OldColumn == nil || diff.NewColumn == nil {
		return false
	}
	computed := isComputedColumnType(diff.OldColumn.Type) || isComputedColumnType(diff.NewColumn.Type)
	if !computed {
		return false
	}
	// ALTER COLUMN cannot change a computed expression or its explicit
	// nullability, so either difference takes the drop/add path.
	return diff.OldColumn.Type != diff.NewColumn.Type || diff.OldColumn.Nullable != diff.NewColumn.Nullable
}

func isComputedColumnType(columnType string) bool {
	return strings.HasPrefix(strings.TrimSpace(columnType), "AS ")
}

func indexUsesColumn(index *metadata.IndexMetadata, name string) bool {
	for _, expression := range index.Expressions {
		if strings.EqualFold(expression, name) {
			return true
		}
	}
	for _, column := range indexIncludeColumns(index.Definition) {
		if strings.EqualFold(column, name) {
			return true
		}
	}
	if predicate := indexWhereClause(index.Definition); predicate != "" {
		if strings.Contains(strings.ToLower(predicate), strings.ToLower(name)) {
			return true
		}
	}
	return false
}

func checkChangeScheduled(table *TableDiffEntry, name string) bool {
	for _, checkDiff := range table.Checks {
		if checkDiffName(checkDiff) == name {
			return true
		}
	}
	return false
}

func indexChangeScheduled(table *TableDiffEntry, name string) bool {
	for _, indexDiff := range table.Indexes {
		if indexDiffName(indexDiff) == name {
			return true
		}
	}
	return false
}

func generateColumnDefinition(column *metadata.ColumnMetadata) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "%s ", quote(column.Name))
	isComputedColumn := strings.Contains(column.Type, " AS (") || strings.HasPrefix(column.Type, "AS (")
	if isComputedColumn {
		buf.WriteString(column.Type)
		if !column.Nullable && strings.Contains(strings.ToUpper(column.Type), "PERSISTED") {
			buf.WriteString(" NOT NULL")
		}
		return buf.String()
	}
	buf.WriteString(column.Type)
	if column.Collation != "" {
		fmt.Fprintf(&buf, " COLLATE %s", column.Collation)
	}
	if column.IsIdentity {
		fmt.Fprintf(&buf, " IDENTITY(%d,%d)", column.IdentitySeed, column.IdentityIncrement)
	}
	if column.Nullable {
		buf.WriteString(" NULL")
	} else {
		buf.WriteString(" NOT NULL")
	}
	if column.GetDefault() != "" {
		if column.DefaultConstraintName != "" {
			fmt.Fprintf(&buf, " CONSTRAINT %s DEFAULT %s", quote(column.DefaultConstraintName), column.GetDefault())
		} else {
			buf.WriteString(" DEFAULT ")
			buf.WriteString(column.GetDefault())
		}
	}
	return buf.String()
}

func generateAlterColumn(schemaName, tableName string, colDiff *ColumnDiffEntry, droppedDefaults map[string]bool) string {
	var buf strings.Builder
	oldDefault := getColumnDefaultValue(colDiff.OldColumn)
	newDefault := getColumnDefaultValue(colDiff.NewColumn)
	alterColumn := !typesEqual(colDiff.OldColumn.Type, colDiff.NewColumn.Type) ||
		colDiff.OldColumn.Nullable != colDiff.NewColumn.Nullable ||
		colDiff.OldColumn.Collation != colDiff.NewColumn.Collation
	droppedDefaultForAlter := false
	if alterColumn {
		// SQL Server refuses ALTER COLUMN while a default is bound to the
		// column, even when that default is unchanged.
		if oldDefault != "" && colDiff.OldColumn.DefaultConstraintName != "" && !droppedDefaults[schemaName+"."+colDiff.OldColumn.DefaultConstraintName] {
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(schemaName), quote(tableName), quote(colDiff.OldColumn.DefaultConstraintName))
			droppedDefaultForAlter = true
		}
		fmt.Fprintf(&buf, "ALTER TABLE %s.%s ALTER COLUMN %s %s", quote(schemaName), quote(tableName), quote(colDiff.NewColumn.Name), colDiff.NewColumn.Type)
		if colDiff.NewColumn.Collation != "" {
			fmt.Fprintf(&buf, " COLLATE %s", colDiff.NewColumn.Collation)
		}
		if colDiff.NewColumn.Nullable {
			buf.WriteString(" NULL")
		} else {
			buf.WriteString(" NOT NULL")
		}
		buf.WriteString(";\n")
	}
	if colDiff.OldColumn.IsIdentity != colDiff.NewColumn.IsIdentity ||
		colDiff.OldColumn.IdentitySeed != colDiff.NewColumn.IdentitySeed ||
		colDiff.OldColumn.IdentityIncrement != colDiff.NewColumn.IdentityIncrement {
		fmt.Fprintf(&buf, "-- UNSUPPORTED: identity property change on %s.%s.%s requires a table rebuild; left unchanged.\n",
			quote(schemaName), quote(tableName), quote(colDiff.NewColumn.Name))
	}
	if oldDefault != newDefault || (droppedDefaultForAlter && newDefault != "") {
		if !droppedDefaultForAlter && oldDefault != "" && colDiff.OldColumn.DefaultConstraintName != "" && !droppedDefaults[schemaName+"."+colDiff.OldColumn.DefaultConstraintName] {
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(schemaName), quote(tableName), quote(colDiff.OldColumn.DefaultConstraintName))
		}
		if newDefault != "" {
			constraintName := colDiff.NewColumn.DefaultConstraintName
			if constraintName == "" {
				constraintName = "DF_" + tableName + "_" + colDiff.NewColumn.Name
			}
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s DEFAULT %s FOR %s;\n",
				quote(schemaName), quote(tableName), quote(constraintName), newDefault, quote(colDiff.NewColumn.Name))
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
	case "XML":
		return generateXMLIndex(schemaName, tableName, index)
	case "CLUSTERED COLUMNSTORE":
		fmt.Fprintf(&buf, " CLUSTERED COLUMNSTORE INDEX %s ON %s.%s", quote(index.Name), quote(schemaName), quote(tableName))
		return buf.String()
	case "NONCLUSTERED COLUMNSTORE":
		fmt.Fprintf(&buf, " NONCLUSTERED COLUMNSTORE INDEX %s ON %s.%s", quote(index.Name), quote(schemaName), quote(tableName))
		if len(index.Expressions) > 0 {
			buf.WriteString(" (")
			for i, expr := range index.Expressions {
				if i > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(&buf, "%s", quote(expr))
			}
			buf.WriteString(")")
		}
		if where := indexWhereClause(index.Definition); where != "" {
			fmt.Fprintf(&buf, " WHERE %s", where)
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
		fmt.Fprintf(&buf, " INDEX %s ON %s.%s (", quote(index.Name), quote(schemaName), quote(tableName))
		for i, expr := range index.Expressions {
			if i > 0 {
				buf.WriteString(", ")
			}
			fmt.Fprintf(&buf, "%s", quote(expr))
			appendMigrationDirection(&buf, index.Descending, i)
		}
		buf.WriteString(")")
		if include := indexIncludeColumns(index.Definition); len(include) > 0 {
			buf.WriteString(" INCLUDE (")
			for i, column := range include {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(quote(column))
			}
			buf.WriteString(")")
		}
		if where := indexWhereClause(index.Definition); where != "" {
			fmt.Fprintf(&buf, " WHERE %s", where)
		}
		return buf.String()
	}
}

func generateXMLIndex(schemaName, tableName string, index *metadata.IndexMetadata) string {
	column := ""
	if len(index.Expressions) > 0 {
		column = quote(index.Expressions[0])
	}
	if strings.HasPrefix(index.Definition, "USING XML INDEX ") {
		rest := strings.TrimPrefix(index.Definition, "USING XML INDEX ")
		parent, secondaryType := rest, "PATH"
		if before, after, ok := strings.Cut(rest, " FOR "); ok {
			parent, secondaryType = before, after
		}
		return fmt.Sprintf("CREATE XML INDEX %s ON %s.%s (%s) USING XML INDEX %s FOR %s",
			quote(index.Name), quote(schemaName), quote(tableName), column, quote(parent), secondaryType)
	}
	return fmt.Sprintf("CREATE PRIMARY XML INDEX %s ON %s.%s (%s)",
		quote(index.Name), quote(schemaName), quote(tableName), column)
}

func generateAddForeignKey(schemaName, tableName string, fk *metadata.ForeignKeyMetadata) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s FOREIGN KEY (", quote(schemaName), quote(tableName), quote(fk.Name))
	for i, col := range fk.Columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%s", quote(col))
	}
	buf.WriteString(") REFERENCES ")
	if fk.ReferencedSchema != "" {
		fmt.Fprintf(&buf, "%s.", quote(fk.ReferencedSchema))
	}
	buf.WriteString(quote(fk.ReferencedTable))
	if len(fk.ReferencedColumns) > 0 {
		buf.WriteString(" (")
		for i, col := range fk.ReferencedColumns {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(quote(col))
		}
		buf.WriteString(")")
	}
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

func literal(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func generateExtendedPropertySQL(action, objectType, schemaName, objectName, comment string) string {
	var buf strings.Builder
	switch action {
	case "ADD":
		fmt.Fprintf(&buf, "EXEC sp_addextendedproperty 'MS_Description', N'%s', 'SCHEMA', '%s', '%s', '%s'", strings.ReplaceAll(comment, "'", "''"), literal(schemaName), objectType, literal(objectName))
	case "UPDATE":
		fmt.Fprintf(&buf, "EXEC sp_updateextendedproperty 'MS_Description', N'%s', 'SCHEMA', '%s', '%s', '%s'", strings.ReplaceAll(comment, "'", "''"), literal(schemaName), objectType, literal(objectName))
	case "DROP":
		fmt.Fprintf(&buf, "EXEC sp_dropextendedproperty 'MS_Description', 'SCHEMA', '%s', '%s', '%s'", literal(schemaName), objectType, literal(objectName))
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
		fmt.Fprintf(&buf, "EXEC sp_addextendedproperty 'MS_Description', N'%s', 'SCHEMA', '%s', 'TABLE', '%s', 'COLUMN', '%s'", strings.ReplaceAll(comment, "'", "''"), literal(schemaName), literal(tableName), literal(columnName))
	case "UPDATE":
		fmt.Fprintf(&buf, "EXEC sp_updateextendedproperty 'MS_Description', N'%s', 'SCHEMA', '%s', 'TABLE', '%s', 'COLUMN', '%s'", strings.ReplaceAll(comment, "'", "''"), literal(schemaName), literal(tableName), literal(columnName))
	case "DROP":
		fmt.Fprintf(&buf, "EXEC sp_dropextendedproperty 'MS_Description', 'SCHEMA', '%s', 'TABLE', '%s', 'COLUMN', '%s'", literal(schemaName), literal(tableName), literal(columnName))
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
			fmt.Fprintf(buf, "DROP VIEW %s.%s;\nGO\n", quote(viewDiff.SchemaName), quote(viewDiff.ViewName))
		}
		return
	}
	for _, viewID := range ordered {
		if viewDiff, ok := viewMap[viewID]; ok {
			fmt.Fprintf(buf, "DROP VIEW %s.%s;\nGO\n", quote(viewDiff.SchemaName), quote(viewDiff.ViewName))
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
			if refID == tableID {
				continue
			}
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
		buf.WriteString(creationDefinition(viewDiff.NewView.Definition))
		if !strings.HasSuffix(strings.TrimSpace(viewDiff.NewView.Definition), ";") {
			buf.WriteString(";")
		}
		buf.WriteString("\nGO\n")
		if viewDiff.NewView.Comment != "" {
			buf.WriteString(generateViewCommentSQL("ADD", viewDiff.SchemaName, viewDiff.ViewName, viewDiff.NewView.Comment))
			buf.WriteString(";\nGO\n")
		}
		for _, trigger := range viewDiff.NewView.Triggers {
			writeTriggerMigration(buf, trigger)
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
	fmt.Fprintf(&buf, "CREATE SPATIAL INDEX %s ON %s.%s (", quote(index.Name), quote(schemaName), quote(tableName))
	for i, expr := range index.Expressions {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%s", quote(expr))
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
	if !storage.AllowRowLocks {
		params = append(params, "ALLOW_ROW_LOCKS = OFF")
	}
	if !storage.AllowPageLocks {
		params = append(params, "ALLOW_PAGE_LOCKS = OFF")
	}
	if storage.Maxdop > 0 {
		params = append(params, fmt.Sprintf("MAXDOP = %d", storage.Maxdop))
	}
	if storage.DataCompression != "" && storage.DataCompression != "NONE" {
		params = append(params, fmt.Sprintf("DATA_COMPRESSION = %s", storage.DataCompression))
	}
	return params
}
