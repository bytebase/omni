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
	if function := blockingSchemaBoundDrop(from, diff); function != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: schema-bound function %s reads a table being dropped; drop the function first manually.\n", function)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	if function := blockingSchemaBoundViewDrop(from, diff); function != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: schema-bound function %s reads a view being dropped; drop the function first manually.\n", function)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	if column := blockingComputedColumn(diff); column != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: column change under computed column %s; drop and recreate the computed column manually.\n", column)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	rebuiltFunctions, rebuiltViews := rebuiltSchemaBoundDependencies(from, to, diff)
	allFunctionDrops := collectFunctionDrops(diff, rebuiltFunctions)
	allFunctionCreates := collectFunctionCreates(diff, rebuiltFunctions)
	viewsToDrop := collectViewsToDrop(diff, rebuiltViews)
	viewsToCreate := collectViewsToCreate(diff, rebuiltViews)
	earlyFuncDrops, remainingFuncDrops := splitFunctionsBeforeViewDrops(allFunctionDrops, viewsToDrop)
	if msg := blockingSkipDumpSchemaBoundConsumer(from, allFunctionDrops, viewsToDrop, collectDroppedTables(diff)); msg != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: %s; update the skipped object manually.\n", msg)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	if msg := blockingFunctionComputedColumn(from, diff, earlyFuncDrops, remainingFuncDrops); msg != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: function %s; drop and recreate the computed column manually.\n", msg)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	deferredFuncCreates := functionsDeferredPastTables(diff, allFunctionCreates, viewsToCreate)
	if msg := blockingDeferredComputedColumn(diff, allFunctionCreates, deferredFuncCreates); msg != "" {
		text := fmt.Sprintf("-- UNSUPPORTED: computed column %s; apply this change manually.\n", msg)
		return &MigrationPlan{Ops: []MigrationOp{{Type: OpScript, SQL: text}}}
	}
	hasViewsToCreate := hasCreateViews(diff) || len(rebuiltViews) > 0
	hasFunctionsToCreate := hasCreateFunctions(diff) || len(rebuiltFunctions) > 0

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

	droppedDefaults := make(map[string]bool)
	var restoreFunctionDefaults []functionDefaultRef
	droppedChecks := make(map[string]bool)
	var restoreFunctionChecks []functionCheckRef

	if len(earlyFuncDrops) > 0 {
		earlyDefaults, earlyRestoreDefaults := dropFunctionDefaults(from, diff, earlyFuncDrops, droppedDefaults, &buf)
		for name := range earlyDefaults {
			droppedDefaults[name] = true
		}
		restoreFunctionDefaults = append(restoreFunctionDefaults, earlyRestoreDefaults...)
		earlyChecks, earlyRestoreChecks := dropFunctionChecks(from, diff, earlyFuncDrops, droppedChecks, &buf)
		for name := range earlyChecks {
			droppedChecks[name] = true
		}
		restoreFunctionChecks = append(restoreFunctionChecks, earlyRestoreChecks...)
		dropScheduledConstraintsForEarlyFunctions(diff, earlyFuncDrops, droppedDefaults, droppedChecks, &buf)
		dropViewsAndFunctionsInOrder(viewsToDrop, earlyFuncDrops, &buf)
	} else {
		dropViewsInOrder(diff, rebuiltViews, &buf)
	}

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
				if droppedChecks[checkKey(tableDiff.SchemaName, tableDiff.TableName, checkDiff.OldCheckConstraint.Name)] {
					continue
				}
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
			if getColumnDefaultValue(colDiff.OldColumn) != "" && colDiff.OldColumn.DefaultConstraintName != "" && !droppedDefaults[tableDiff.SchemaName+"."+colDiff.OldColumn.DefaultConstraintName] {
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

	droppedTables := collectDroppedTables(diff)
	hasTableDrops := false

	if droppedFunctionsDependOnDroppedTables(remainingFuncDrops, droppedTables) {
		droppedFnDefaults, moreRestoreDefaults := dropFunctionDefaults(from, diff, remainingFuncDrops, droppedDefaults, &buf)
		for name := range droppedFnDefaults {
			droppedDefaults[name] = true
		}
		restoreFunctionDefaults = append(restoreFunctionDefaults, moreRestoreDefaults...)
		droppedFnChecks, moreRestoreChecks := dropFunctionChecks(from, diff, remainingFuncDrops, droppedChecks, &buf)
		for name := range droppedFnChecks {
			droppedChecks[name] = true
		}
		restoreFunctionChecks = append(restoreFunctionChecks, moreRestoreChecks...)
		for _, funcDiff := range diff.Functions {
			if funcDiff.Action == DiffDrop && functionReplacementUsedByComputed(diff, funcDiff) {
				fmt.Fprintf(&buf, "-- UNSUPPORTED: function %s.%s is used by a computed column; left unchanged.\n", quote(funcDiff.SchemaName), quote(funcDiff.FunctionName))
			}
		}
		hasTableDrops = dropTablesAndFunctionsInOrder(droppedTables, remainingFuncDrops, hasTableAlterations, &buf)
	} else {
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

		droppedFnDefaults, moreRestoreDefaults := dropFunctionDefaults(from, diff, remainingFuncDrops, droppedDefaults, &buf)
		for name := range droppedFnDefaults {
			droppedDefaults[name] = true
		}
		restoreFunctionDefaults = append(restoreFunctionDefaults, moreRestoreDefaults...)
		droppedFnChecks, moreRestoreChecks := dropFunctionChecks(from, diff, remainingFuncDrops, droppedChecks, &buf)
		for name := range droppedFnChecks {
			droppedChecks[name] = true
		}
		restoreFunctionChecks = append(restoreFunctionChecks, moreRestoreChecks...)
		for _, funcDiff := range diff.Functions {
			if funcDiff.Action == DiffDrop && functionReplacementUsedByComputed(diff, funcDiff) {
				fmt.Fprintf(&buf, "-- UNSUPPORTED: function %s.%s is used by a computed column; left unchanged.\n", quote(funcDiff.SchemaName), quote(funcDiff.FunctionName))
			}
		}
		for _, funcDiff := range orderedFunctionDrops(remainingFuncDrops) {
			fmt.Fprintf(&buf, "DROP FUNCTION %s.%s;\nGO\n", quote(funcDiff.SchemaName), quote(funcDiff.FunctionName))
		}
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
		hasCreateOrAlterTables(diff) || hasViewsToCreate ||
		hasFunctionsToCreate || hasCreateProcedures(diff)
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
	if len(schemasToCreate) > 0 && (hasCreateOrAlterTables(diff) || hasViewsToCreate || hasFunctionsToCreate || hasCreateProcedures(diff)) {
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

	tableFunctions := migrationTableFunctions(diff, allFunctionCreates)
	// Functions that depend on a view being created/rebuilt or on a table whose
	// columns are being altered cannot be created ahead of tables; keep them in
	// the post-table view/function phase and defer any table DEFAULT or CHECK
	// additions calling them until afterward.
	for id := range deferredFuncCreates {
		delete(tableFunctions, id)
	}
	var deferredTableChecks []functionCheckRef
	deferredTableCheckKeys := make(map[string]bool)
	var deferredTableDefaults []functionDefaultRef
	deferredTableDefaultCols := make(map[string]bool)
	if len(deferredFuncCreates) > 0 {
		for _, tableDiff := range diff.Tables {
			switch tableDiff.Action {
			case DiffCreate:
				if tableDiff.NewTable == nil {
					continue
				}
				for _, col := range tableDiff.NewTable.Columns {
					if col.GetDefault() == "" || !exprReferencesAnyDeferredFunc(col.GetDefault(), tableDiff.SchemaName, allFunctionCreates, deferredFuncCreates) {
						continue
					}
					colKey := checkKey(tableDiff.SchemaName, tableDiff.TableName, col.Name)
					if !deferredTableDefaultCols[colKey] {
						deferredTableDefaultCols[colKey] = true
						constraintName := col.DefaultConstraintName
						if constraintName == "" {
							constraintName = "DF_" + tableDiff.TableName + "_" + col.Name
						}
						deferredTableDefaults = append(deferredTableDefaults, functionDefaultRef{
							schema:     tableDiff.SchemaName,
							table:      tableDiff.TableName,
							column:     col.Name,
							constraint: constraintName,
							expression: col.GetDefault(),
						})
					}
				}
				for _, check := range tableDiff.NewTable.CheckConstraints {
					if check.Expression == "" || !exprReferencesAnyDeferredFunc(check.Expression, tableDiff.SchemaName, allFunctionCreates, deferredFuncCreates) {
						continue
					}
					cKey := checkKey(tableDiff.SchemaName, tableDiff.TableName, check.Name)
					if !deferredTableCheckKeys[cKey] {
						deferredTableCheckKeys[cKey] = true
						deferredTableChecks = append(deferredTableChecks, functionCheckRef{
							schema:     tableDiff.SchemaName,
							table:      tableDiff.TableName,
							constraint: check.Name,
							expression: check.Expression,
						})
					}
				}
			case DiffAlter:
				for _, colDiff := range tableDiff.Columns {
					if (colDiff.Action != DiffCreate && colDiff.Action != DiffAlter) || colDiff.NewColumn == nil {
						continue
					}
					newDefault := getColumnDefaultValue(colDiff.NewColumn)
					if newDefault == "" || !exprReferencesAnyDeferredFunc(newDefault, tableDiff.SchemaName, allFunctionCreates, deferredFuncCreates) {
						continue
					}
					if colDiff.Action == DiffAlter && !columnRequiresRebuild(colDiff) {
						oldDefault := getColumnDefaultValue(colDiff.OldColumn)
						alterColumn := !typesEqual(colDiff.OldColumn.Type, colDiff.NewColumn.Type) ||
							colDiff.OldColumn.Nullable != colDiff.NewColumn.Nullable ||
							colDiff.OldColumn.Collation != colDiff.NewColumn.Collation
						if expressionsEqual(oldDefault, newDefault) && !alterColumn {
							continue
						}
					}
					colKey := checkKey(tableDiff.SchemaName, tableDiff.TableName, colDiff.NewColumn.Name)
					if !deferredTableDefaultCols[colKey] {
						deferredTableDefaultCols[colKey] = true
						constraintName := colDiff.NewColumn.DefaultConstraintName
						if constraintName == "" {
							constraintName = "DF_" + tableDiff.TableName + "_" + colDiff.NewColumn.Name
						}
						deferredTableDefaults = append(deferredTableDefaults, functionDefaultRef{
							schema:     tableDiff.SchemaName,
							table:      tableDiff.TableName,
							column:     colDiff.NewColumn.Name,
							constraint: constraintName,
							expression: newDefault,
						})
					}
				}
				for _, checkDiff := range tableDiff.Checks {
					if checkDiff.Action != DiffCreate || checkDiff.NewCheckConstraint == nil {
						continue
					}
					if !exprReferencesAnyDeferredFunc(checkDiff.NewCheckConstraint.Expression, tableDiff.SchemaName, allFunctionCreates, deferredFuncCreates) {
						continue
					}
					cKey := checkKey(tableDiff.SchemaName, tableDiff.TableName, checkDiff.NewCheckConstraint.Name)
					if !deferredTableCheckKeys[cKey] {
						deferredTableCheckKeys[cKey] = true
						deferredTableChecks = append(deferredTableChecks, functionCheckRef{
							schema:     tableDiff.SchemaName,
							table:      tableDiff.TableName,
							constraint: checkDiff.NewCheckConstraint.Name,
							expression: checkDiff.NewCheckConstraint.Expression,
						})
					}
				}
			}
		}
	}

	if len(tableFunctions) > 0 && createdFunctionsDependOnCreatedTables(diff, allFunctionCreates, tableFunctions) {
		createTablesAndFunctionsInOrder(diff, allFunctionCreates, tableFunctions, deferredTableDefaultCols, deferredTableCheckKeys, &buf)
	} else {
		for _, funcDiff := range orderedFunctionDiffs(allFunctionCreates) {
			if tableFunctions[objectID(funcDiff.SchemaName, funcDiff.FunctionName)] {
				writeFunctionMigration(&buf, funcDiff)
			}
		}
		createTablesInOrder(diff, deferredTableDefaultCols, deferredTableCheckKeys, &buf)
	}

	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffAlter {
			alterTableSQL := generateAlterTable(tableDiff, droppedDefaults, droppedChecks, deferredTableDefaultCols, deferredTableCheckKeys)
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
					if columnRequiresRebuild(colDiff) {
						if newComment != "" {
							buf.WriteString(generateColumnCommentSQL("ADD", tableDiff.SchemaName, tableDiff.TableName, colDiff.NewColumn.Name, newComment))
							buf.WriteString(";\n")
							hasCommentChanges = true
						}
						continue
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
	if (hasCreateOrAlterTables(diff) || hasTableDrops) && (hasViewsToCreate || hasFunctionsToCreate || hasCreateProcedures(diff)) {
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

	var postTableFuncCreates []*FunctionDiffEntry
	for _, funcDiff := range allFunctionCreates {
		if !tableFunctions[objectID(funcDiff.SchemaName, funcDiff.FunctionName)] {
			postTableFuncCreates = append(postTableFuncCreates, funcDiff)
		}
	}
	if createdFunctionsDependOnCreatedViews(postTableFuncCreates, viewsToCreate) {
		createViewsAndFunctionsInOrder(viewsToCreate, postTableFuncCreates, &buf)
	} else {
		earlyFunctions := migrationEarlyFunctions(diff, rebuiltViews, allFunctionCreates)
		for name := range tableFunctions {
			delete(earlyFunctions, name)
		}
		for _, funcDiff := range orderedFunctionDiffs(allFunctionCreates) {
			if !earlyFunctions[objectID(funcDiff.SchemaName, funcDiff.FunctionName)] {
				continue
			}
			writeFunctionMigration(&buf, funcDiff)
		}

		createViewsInOrder(diff, rebuiltViews, &buf)

		for _, funcDiff := range orderedFunctionDiffs(allFunctionCreates) {
			if earlyFunctions[objectID(funcDiff.SchemaName, funcDiff.FunctionName)] || tableFunctions[objectID(funcDiff.SchemaName, funcDiff.FunctionName)] {
				continue
			}
			writeFunctionMigration(&buf, funcDiff)
		}
	}

	restoredConstraint := false
	for _, ref := range restoreFunctionDefaults {
		if deferredTableDefaultCols[checkKey(ref.schema, ref.table, ref.column)] {
			continue
		}
		if !allFunctionsRecreated(allFunctionCreates, ref.functions) || columnDefaultChanged(diff, ref.schema, ref.table, ref.column) {
			continue
		}
		fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s DEFAULT %s FOR %s;\n",
			quote(ref.schema), quote(ref.table), quote(ref.constraint), ref.expression, quote(ref.column))
		restoredConstraint = true
	}
	for _, ref := range restoreFunctionChecks {
		if deferredTableCheckKeys[checkKey(ref.schema, ref.table, ref.constraint)] {
			continue
		}
		if !allFunctionsRecreated(allFunctionCreates, ref.functions) || !checkExistsUnchangedInTarget(to, diff, ref.schema, ref.table, ref.constraint) {
			continue
		}
		fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s CHECK %s;\n",
			quote(ref.schema), quote(ref.table), quote(ref.constraint), checkExpression(ref.expression))
		restoredConstraint = true
	}
	for _, ref := range deferredTableDefaults {
		fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s DEFAULT %s FOR %s;\n",
			quote(ref.schema), quote(ref.table), quote(ref.constraint), ref.expression, quote(ref.column))
		restoredConstraint = true
	}
	for _, ref := range deferredTableChecks {
		fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s CHECK %s;\n",
			quote(ref.schema), quote(ref.table), quote(ref.constraint), checkExpression(ref.expression))
		restoredConstraint = true
	}
	if restoredConstraint {
		buf.WriteString("GO\n")
	}

	for _, procDiff := range diff.Procedures {
		switch procDiff.Action {
		case DiffCreate:
			buf.WriteString(qualifyModuleDefinition(creationDefinition(procDiff.NewProcedure.Definition), procDiff.SchemaName, false))
			if !strings.HasSuffix(strings.TrimSpace(procDiff.NewProcedure.Definition), ";") {
				buf.WriteString(";")
			}
			buf.WriteString("\nGO\n")
		case DiffAlter:
			fmt.Fprintf(&buf, "DROP PROCEDURE %s.%s;\nGO\n", quote(procDiff.SchemaName), quote(procDiff.ProcedureName))
			buf.WriteString(qualifyModuleDefinition(creationDefinition(procDiff.NewProcedure.Definition), procDiff.SchemaName, false))
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

func generateCreateTable(schemaName, tableName string, table *metadata.TableMetadata, deferredDefaultCols, deferredChecks map[string]bool) string {
	var buf strings.Builder
	var activeChecks []*metadata.CheckConstraintMetadata
	for _, check := range table.CheckConstraints {
		if deferredChecks[checkKey(schemaName, tableName, check.Name)] {
			continue
		}
		activeChecks = append(activeChecks, check)
	}
	fmt.Fprintf(&buf, "CREATE TABLE %s.%s (\n", quote(schemaName), quote(tableName))
	for i, column := range table.Columns {
		buf.WriteString("  ")
		buf.WriteString(generateColumnDefinitionFiltered(column, !deferredDefaultCols[checkKey(schemaName, tableName, column.Name)]))
		if i < len(table.Columns)-1 || hasConstraintsInTable(table) || len(activeChecks) > 0 {
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

	if constraintAdded && len(activeChecks) > 0 {
		buf.WriteString(",\n")
	} else if constraintAdded {
		buf.WriteString("\n")
	}
	for i, check := range activeChecks {
		fmt.Fprintf(&buf, "  CONSTRAINT %s CHECK %s", quote(check.Name), checkExpression(check.Expression))
		if i < len(activeChecks)-1 {
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

func generateAlterTable(tableDiff *TableDiffEntry, droppedDefaults, droppedChecks, deferredDefaultCols, deferredChecks map[string]bool) string {
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
		if colDiff.OldColumn.Default != "" && colDiff.OldColumn.DefaultConstraintName != "" && !droppedDefaults[tableDiff.SchemaName+"."+colDiff.OldColumn.DefaultConstraintName] {
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(colDiff.OldColumn.DefaultConstraintName))
		}
		for _, check := range tableDiff.OldTable.CheckConstraints {
			if !strings.Contains(strings.ToLower(check.Expression), strings.ToLower(colDiff.OldColumn.Name)) {
				continue
			}
			if checkChangeScheduled(tableDiff, check.Name) || containsString(rebuiltChecks, check.Name) || droppedChecks[checkKey(tableDiff.SchemaName, tableDiff.TableName, check.Name)] {
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
				checkChangeScheduled(tableDiff, check.Name) || containsString(rebuiltChecks, check.Name) || droppedChecks[checkKey(tableDiff.SchemaName, tableDiff.TableName, check.Name)] {
				continue
			}
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(check.Name))
			rebuiltChecks = append(rebuiltChecks, check.Name)
		}
	}
	for _, colDiff := range tableDiff.Columns {
		if colDiff.Action == DiffCreate || (colDiff.Action == DiffAlter && columnRequiresRebuild(colDiff)) {
			includeDefault := !deferredDefaultCols[checkKey(tableDiff.SchemaName, tableDiff.TableName, colDiff.NewColumn.Name)]
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), generateColumnDefinitionFiltered(colDiff.NewColumn, includeDefault))
		}
	}
	for _, colDiff := range tableDiff.Columns {
		if colDiff.Action == DiffAlter && !columnRequiresRebuild(colDiff) {
			buf.WriteString(generateAlterColumn(tableDiff.SchemaName, tableDiff.TableName, colDiff, droppedDefaults, deferredDefaultCols))
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
			if deferredChecks[checkKey(tableDiff.SchemaName, tableDiff.TableName, checkDiff.NewCheckConstraint.Name)] {
				continue
			}
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s ADD CONSTRAINT %s CHECK %s;\n",
				quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(checkDiff.NewCheckConstraint.Name), checkExpression(checkDiff.NewCheckConstraint.Expression))
		}
	}
	return buf.String()
}

func blockingSchemaBoundViewDrop(from *Catalog, diff *SchemaDiff) string {
	if from == nil {
		return ""
	}
	var droppedViews []viewWithSchema
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffDrop && viewDiff.OldView != nil {
			droppedViews = append(droppedViews, viewWithSchema{schema: viewDiff.SchemaName, view: viewDiff.OldView})
		}
	}
	if len(droppedViews) == 0 {
		return ""
	}
	droppedFunctions := make(map[string]bool)
	for _, funcDiff := range diff.Functions {
		if funcDiff.Action == DiffDrop && !functionReplacementUsedByComputed(diff, funcDiff) {
			droppedFunctions[strings.ToLower(objectID(funcDiff.SchemaName, funcDiffBaseName(funcDiff)))] = true
		}
	}
	for _, schema := range from.Metadata().Schemas {
		for _, function := range schema.Functions {
			if function.GetSkipDump() || !isSchemaBound(function.Definition) {
				continue
			}
			if droppedFunctions[strings.ToLower(objectID(schema.Name, function.Name))] {
				continue
			}
			for _, dropped := range droppedViews {
				if definitionReferencesTable(function.Definition, schema.Name, dropped.schema, dropped.view.Name) {
					return quote(schema.Name) + "." + quote(function.Name)
				}
			}
		}
	}
	return ""
}

func blockingSchemaBoundDrop(from *Catalog, diff *SchemaDiff) string {
	if from == nil {
		return ""
	}
	droppedTables := collectDroppedTables(diff)
	if len(droppedTables) == 0 {
		return ""
	}
	droppedFunctions := make(map[string]bool)
	for _, funcDiff := range diff.Functions {
		if funcDiff.Action == DiffDrop && !functionReplacementUsedByComputed(diff, funcDiff) {
			droppedFunctions[strings.ToLower(objectID(funcDiff.SchemaName, funcDiffBaseName(funcDiff)))] = true
		}
	}
	for _, schema := range from.Metadata().Schemas {
		for _, function := range schema.Functions {
			if function.GetSkipDump() || !isSchemaBound(function.Definition) {
				continue
			}
			// A schema-bound function that is itself being dropped (or replaced)
			// will be scheduled before the table drop.
			if droppedFunctions[strings.ToLower(objectID(schema.Name, function.Name))] {
				continue
			}
			for _, dropped := range droppedTables {
				if definitionReferencesTable(function.Definition, schema.Name, dropped.schema, dropped.table.Name) {
					return quote(schema.Name) + "." + quote(function.Name)
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
	schema     string
	table      string
	column     string
	constraint string
	expression string
	functions  []string
}

type functionCheckRef struct {
	schema     string
	table      string
	constraint string
	expression string
	functions  []string
}

func checkKey(schemaName, tableName, constraintName string) string {
	return strings.ToLower(schemaName + "\x00" + tableName + "\x00" + constraintName)
}

func funcDiffBaseName(funcDiff *FunctionDiffEntry) string {
	if funcDiff.OldFunction != nil && funcDiff.OldFunction.Name != "" {
		return funcDiff.OldFunction.Name
	}
	if funcDiff.NewFunction != nil && funcDiff.NewFunction.Name != "" {
		return funcDiff.NewFunction.Name
	}
	return funcDiff.FunctionName
}

func allFunctionsRecreated(creates []*FunctionDiffEntry, required []string) bool {
	if len(required) == 0 {
		return false
	}
	created := make(map[string]bool, len(creates))
	for _, funcDiff := range creates {
		if funcDiff.Action == DiffCreate {
			created[strings.ToLower(objectID(funcDiff.SchemaName, funcDiffBaseName(funcDiff)))] = true
		}
	}
	for _, req := range required {
		if !created[req] {
			return false
		}
	}
	return true
}

func checkExistsUnchangedInTarget(to *Catalog, diff *SchemaDiff, schemaName, tableName, constraintName string) bool {
	if to == nil {
		return false
	}
	for _, tableDiff := range diff.Tables {
		if !strings.EqualFold(tableDiff.SchemaName, schemaName) || !strings.EqualFold(tableDiff.TableName, tableName) {
			continue
		}
		if tableDiff.Action == DiffDrop {
			return false
		}
		if checkChangeScheduled(tableDiff, constraintName) {
			return false
		}
	}
	table := to.GetTable(schemaName, tableName)
	if table == nil {
		return false
	}
	for _, check := range table.CheckConstraints {
		if strings.EqualFold(check.Name, constraintName) {
			return true
		}
	}
	return false
}

// rebuiltSchemaBoundDependencies finds unchanged schema-bound functions and
// views in from that depend (directly or transitively) on a function being
// dropped or replaced, so they can be dropped before the function drop and
// recreated afterward.
func rebuiltSchemaBoundDependencies(from, to *Catalog, diff *SchemaDiff) ([]functionWithSchema, []viewWithSchema) {
	if from == nil {
		return nil, nil
	}
	affected := make(map[string]functionWithSchema)
	var affectedKeys []string
	addAffected := func(schemaName string, fn *metadata.FunctionMetadata) {
		key := strings.ToLower(objectID(schemaName, fn.Name))
		if _, exists := affected[key]; !exists {
			affected[key] = functionWithSchema{schema: schemaName, function: fn}
			affectedKeys = append(affectedKeys, key)
		}
	}
	for _, funcDiff := range diff.Functions {
		if funcDiff.Action != DiffDrop || functionReplacementUsedByComputed(diff, funcDiff) {
			continue
		}
		fn := funcDiff.OldFunction
		if fn == nil {
			fn = &metadata.FunctionMetadata{Name: funcDiff.FunctionName}
		}
		addAffected(funcDiff.SchemaName, fn)
	}
	// Views already scheduled for DROP or ALTER in diff.Views. Explicitly
	// altered views are seeded into affectedViews so unchanged schema-bound
	// functions reading them are discovered without duplicating their own
	// drop/create entries.
	diffViews := make(map[string]bool)
	var affectedViews []viewWithSchema
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffDrop || viewDiff.Action == DiffAlter {
			diffViews[strings.ToLower(objectID(viewDiff.SchemaName, viewDiff.ViewName))] = true
		}
		if viewDiff.Action == DiffAlter && viewDiff.OldView != nil {
			affectedViews = append(affectedViews, viewWithSchema{schema: viewDiff.SchemaName, view: viewDiff.OldView})
		}
	}
	if len(affected) == 0 && len(affectedViews) == 0 {
		return nil, nil
	}

	// Unified fixpoint over unchanged schema-bound functions and views in from
	// so mixed chains (e.g. f -> v -> g -> v2, including explicitly altered v)
	// are discovered completely.
	var rebuiltFunctions []functionWithSchema
	rebuiltFnSet := make(map[string]bool)
	var rebuiltViews []viewWithSchema
	rebuiltViewSet := make(map[string]bool)
	for {
		added := false
		slices.Sort(affectedKeys)
		for _, schema := range from.Metadata().Schemas {
			for _, fn := range schema.Functions {
				if fn.GetSkipDump() || !isSchemaBound(fn.Definition) {
					continue
				}
				key := strings.ToLower(objectID(schema.Name, fn.Name))
				if _, already := affected[key]; already || rebuiltFnSet[key] {
					continue
				}
				if to != nil && !functionExistsInCatalog(to, schema.Name, fn.Name) {
					continue
				}
				depends := false
				for _, targetKey := range affectedKeys {
					target := affected[targetKey]
					if definitionReferencesFunction(fn.Definition, schema.Name, target.schema, target.function.Name) {
						depends = true
						break
					}
				}
				if !depends {
					for _, vws := range affectedViews {
						if definitionReferencesTable(fn.Definition, schema.Name, vws.schema, vws.view.Name) {
							depends = true
							break
						}
					}
				}
				if depends {
					rebuiltFnSet[key] = true
					rebuiltFunctions = append(rebuiltFunctions, functionWithSchema{schema: schema.Name, function: fn})
					addAffected(schema.Name, fn)
					added = true
				}
			}
			for _, view := range schema.Views {
				if view.GetSkipDump() || !isSchemaBound(view.Definition) {
					continue
				}
				vKey := strings.ToLower(objectID(schema.Name, view.Name))
				if diffViews[vKey] || rebuiltViewSet[vKey] {
					continue
				}
				if to != nil && to.GetView(schema.Name, view.Name) == nil {
					continue
				}
				depends := false
				for _, targetKey := range affectedKeys {
					target := affected[targetKey]
					if expressionReferencesFunction(view.Definition, schema.Name, target.schema, target.function.Name) {
						depends = true
						break
					}
				}
				if !depends && len(rebuiltViewSet) > 0 {
					for _, dep := range viewTableDependencies(view.Definition, schema.Name) {
						if rebuiltViewSet[strings.ToLower(dep)] {
							depends = true
							break
						}
					}
				}
				if depends {
					rebuiltViewSet[vKey] = true
					vws := viewWithSchema{schema: schema.Name, view: view}
					rebuiltViews = append(rebuiltViews, vws)
					affectedViews = append(affectedViews, vws)
					added = true
				}
			}
		}
		if !added {
			break
		}
	}
	return rebuiltFunctions, rebuiltViews
}

// blockingSkipDumpSchemaBoundConsumer checks whether any schema-bound view or
// function marked SkipDump in from depends on a function, view, or table being
// dropped/replaced. Because SkipDump objects are excluded from migration, they
// cannot be rebuilt automatically and would cause SQL Server to reject the drop.
func blockingSkipDumpSchemaBoundConsumer(from *Catalog, functionDrops []*FunctionDiffEntry, viewsToDrop []*ViewDiffEntry, droppedTables []tableWithSchema) string {
	if from == nil || (len(functionDrops) == 0 && len(viewsToDrop) == 0 && len(droppedTables) == 0) {
		return ""
	}
	for _, schema := range from.Metadata().Schemas {
		for _, fn := range schema.Functions {
			if !fn.GetSkipDump() || !isSchemaBound(fn.Definition) {
				continue
			}
			for _, fd := range functionDrops {
				if definitionReferencesFunction(fn.Definition, schema.Name, fd.SchemaName, funcDiffBaseName(fd)) {
					return fmt.Sprintf("schema-bound function %s.%s (marked SkipDump) calls %s.%s",
						quote(schema.Name), quote(fn.Name), quote(fd.SchemaName), quote(fd.FunctionName))
				}
			}
			for _, vd := range viewsToDrop {
				if definitionReferencesTable(fn.Definition, schema.Name, vd.SchemaName, vd.ViewName) {
					return fmt.Sprintf("schema-bound function %s.%s (marked SkipDump) reads view %s.%s",
						quote(schema.Name), quote(fn.Name), quote(vd.SchemaName), quote(vd.ViewName))
				}
			}
			for _, tws := range droppedTables {
				if definitionReferencesTable(fn.Definition, schema.Name, tws.schema, tws.table.Name) {
					return fmt.Sprintf("schema-bound function %s.%s (marked SkipDump) reads table %s.%s",
						quote(schema.Name), quote(fn.Name), quote(tws.schema), quote(tws.table.Name))
				}
			}
		}
		for _, view := range schema.Views {
			if !view.GetSkipDump() || !isSchemaBound(view.Definition) {
				continue
			}
			for _, fd := range functionDrops {
				if expressionReferencesFunction(view.Definition, schema.Name, fd.SchemaName, funcDiffBaseName(fd)) {
					return fmt.Sprintf("schema-bound view %s.%s (marked SkipDump) calls %s.%s",
						quote(schema.Name), quote(view.Name), quote(fd.SchemaName), quote(fd.FunctionName))
				}
			}
			for _, vd := range viewsToDrop {
				if definitionReferencesTable(view.Definition, schema.Name, vd.SchemaName, vd.ViewName) {
					return fmt.Sprintf("schema-bound view %s.%s (marked SkipDump) reads view %s.%s",
						quote(schema.Name), quote(view.Name), quote(vd.SchemaName), quote(vd.ViewName))
				}
			}
			for _, tws := range droppedTables {
				if definitionReferencesTable(view.Definition, schema.Name, tws.schema, tws.table.Name) {
					return fmt.Sprintf("schema-bound view %s.%s (marked SkipDump) reads table %s.%s",
						quote(schema.Name), quote(view.Name), quote(tws.schema), quote(tws.table.Name))
				}
			}
		}
	}
	return ""
}

// blockingDeferredComputedColumn checks whether any created or altered table
// adds or rebuilds a computed column calling a function whose creation must
// wait until after table alterations or view creation.
func blockingDeferredComputedColumn(diff *SchemaDiff, creates []*FunctionDiffEntry, deferredFuncCreates map[string]bool) string {
	if len(deferredFuncCreates) == 0 {
		return ""
	}
	checkColumn := func(schemaName, tableName string, col *metadata.ColumnMetadata) string {
		if col == nil || !isComputedColumnType(col.Type) {
			return ""
		}
		for _, fd := range creates {
			if !deferredFuncCreates[objectID(fd.SchemaName, fd.FunctionName)] {
				continue
			}
			if expressionReferencesFunction(col.Type, schemaName, fd.SchemaName, funcDiffBaseName(fd)) {
				return fmt.Sprintf("%s.%s.%s calls deferred function %s.%s",
					quote(schemaName), quote(tableName), quote(col.Name), quote(fd.SchemaName), quote(fd.FunctionName))
			}
		}
		return ""
	}
	for _, tableDiff := range diff.Tables {
		switch tableDiff.Action {
		case DiffCreate:
			if tableDiff.NewTable == nil {
				continue
			}
			for _, col := range tableDiff.NewTable.Columns {
				if msg := checkColumn(tableDiff.SchemaName, tableDiff.TableName, col); msg != "" {
					return msg
				}
			}
		case DiffAlter:
			for _, colDiff := range tableDiff.Columns {
				if colDiff.Action != DiffCreate && colDiff.Action != DiffAlter {
					continue
				}
				if msg := checkColumn(tableDiff.SchemaName, tableDiff.TableName, colDiff.NewColumn); msg != "" {
					return msg
				}
			}
		}
	}
	return ""
}

// blockingFunctionComputedColumn checks whether any function scheduled for
// dropping (explicitly or via implicit schema-bound rebuild) is still
// referenced by a computed column in from when its DROP FUNCTION runs.
func blockingFunctionComputedColumn(from *Catalog, diff *SchemaDiff, earlyFuncDrops, remainingFuncDrops []*FunctionDiffEntry) string {
	if from == nil || (len(earlyFuncDrops) == 0 && len(remainingFuncDrops) == 0) {
		return ""
	}
	droppedTables := make(map[string]bool)
	droppedColumns := make(map[string]bool)
	for _, tws := range collectDroppedTables(diff) {
		droppedTables[strings.ToLower(objectID(tws.schema, tws.table.Name))] = true
	}
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffAlter {
			continue
		}
		tKey := strings.ToLower(objectID(tableDiff.SchemaName, tableDiff.TableName))
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action == DiffDrop && colDiff.OldColumn != nil {
				droppedColumns[tKey+"\x00"+strings.ToLower(colDiff.OldColumn.Name)] = true
			}
		}
	}
	for _, schema := range from.Metadata().Schemas {
		for _, table := range schema.Tables {
			tKey := strings.ToLower(objectID(schema.Name, table.Name))
			tableDropped := droppedTables[tKey]
			for _, col := range table.Columns {
				if !isComputedColumnType(col.Type) {
					continue
				}
				colDropped := droppedColumns[tKey+"\x00"+strings.ToLower(col.Name)]
				// Early function drops run before DROP COLUMN and DROP TABLE, so
				// any computed column in from still blocks them.
				for _, fd := range earlyFuncDrops {
					if expressionReferencesFunction(col.Type, schema.Name, fd.SchemaName, funcDiffBaseName(fd)) {
						return fmt.Sprintf("%s.%s is used by computed column %s.%s.%s",
							quote(fd.SchemaName), quote(fd.FunctionName), quote(schema.Name), quote(table.Name), quote(col.Name))
					}
				}
				if tableDropped || colDropped {
					continue
				}
				for _, fd := range remainingFuncDrops {
					if expressionReferencesFunction(col.Type, schema.Name, fd.SchemaName, funcDiffBaseName(fd)) {
						return fmt.Sprintf("%s.%s is used by computed column %s.%s.%s",
							quote(fd.SchemaName), quote(fd.FunctionName), quote(schema.Name), quote(table.Name), quote(col.Name))
					}
				}
			}
		}
	}
	// A dropped function left unchanged because a computed column references it
	// also blocks dropping any function it schema-binds to.
	for _, skipped := range diff.Functions {
		if skipped.Action != DiffDrop || skipped.OldFunction == nil || !isSchemaBound(skipped.OldFunction.Definition) {
			continue
		}
		if !functionReplacementUsedByComputed(diff, skipped) {
			continue
		}
		for _, drops := range [][]*FunctionDiffEntry{earlyFuncDrops, remainingFuncDrops} {
			for _, fd := range drops {
				if definitionReferencesFunction(skipped.OldFunction.Definition, skipped.SchemaName, fd.SchemaName, funcDiffBaseName(fd)) {
					return fmt.Sprintf("%s.%s is referenced by schema-bound function %s.%s (used by a computed column)",
						quote(fd.SchemaName), quote(fd.FunctionName), quote(skipped.SchemaName), quote(skipped.FunctionName))
				}
			}
		}
	}
	return ""
}

// dropScheduledConstraintsForEarlyFunctions emits DROP CONSTRAINT before
// dropViewsAndFunctionsInOrder for any scheduled CHECK/DEFAULT drop (on
// altered tables, dropped columns, or dropped tables) that references an
// early-dropped function, recording it so the later drop phase skips it.
func dropScheduledConstraintsForEarlyFunctions(diff *SchemaDiff, earlyFuncDrops []*FunctionDiffEntry, droppedDefaults, droppedChecks map[string]bool, buf *strings.Builder) {
	for _, tableDiff := range diff.Tables {
		switch tableDiff.Action {
		case DiffAlter:
			for _, checkDiff := range tableDiff.Checks {
				if checkDiff.Action != DiffDrop || checkDiff.OldCheckConstraint == nil {
					continue
				}
				cKey := checkKey(tableDiff.SchemaName, tableDiff.TableName, checkDiff.OldCheckConstraint.Name)
				if droppedChecks[cKey] {
					continue
				}
				for _, fd := range earlyFuncDrops {
					if expressionReferencesFunction(checkDiff.OldCheckConstraint.Expression, tableDiff.SchemaName, fd.SchemaName, funcDiffBaseName(fd)) {
						fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(checkDiff.OldCheckConstraint.Name))
						droppedChecks[cKey] = true
						break
					}
				}
			}
			for _, colDiff := range tableDiff.Columns {
				if colDiff.Action != DiffDrop || colDiff.OldColumn == nil {
					continue
				}
				if getColumnDefaultValue(colDiff.OldColumn) == "" || colDiff.OldColumn.DefaultConstraintName == "" {
					continue
				}
				dKey := tableDiff.SchemaName + "." + colDiff.OldColumn.DefaultConstraintName
				if droppedDefaults[dKey] {
					continue
				}
				for _, fd := range earlyFuncDrops {
					if expressionReferencesFunction(colDiff.OldColumn.Default, tableDiff.SchemaName, fd.SchemaName, funcDiffBaseName(fd)) {
						fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(colDiff.OldColumn.DefaultConstraintName))
						droppedDefaults[dKey] = true
						break
					}
				}
			}
		case DiffDrop:
			if tableDiff.OldTable == nil {
				continue
			}
			for _, check := range tableDiff.OldTable.CheckConstraints {
				if check.Name == "" || check.Expression == "" {
					continue
				}
				for _, fd := range earlyFuncDrops {
					if expressionReferencesFunction(check.Expression, tableDiff.SchemaName, fd.SchemaName, funcDiffBaseName(fd)) {
						fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(check.Name))
						break
					}
				}
			}
			for _, col := range tableDiff.OldTable.Columns {
				if col.Default == "" || col.DefaultConstraintName == "" {
					continue
				}
				for _, fd := range earlyFuncDrops {
					if expressionReferencesFunction(col.Default, tableDiff.SchemaName, fd.SchemaName, funcDiffBaseName(fd)) {
						fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tableDiff.SchemaName), quote(tableDiff.TableName), quote(col.DefaultConstraintName))
						break
					}
				}
			}
		}
	}
}

func collectViewsToDrop(diff *SchemaDiff, rebuiltViews []viewWithSchema) []*ViewDiffEntry {
	var views []*ViewDiffEntry
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffDrop || viewDiff.Action == DiffAlter {
			views = append(views, viewDiff)
		}
	}
	for _, vws := range rebuiltViews {
		views = append(views, &ViewDiffEntry{
			Action:     DiffDrop,
			SchemaName: vws.schema,
			ViewName:   vws.view.Name,
			OldView:    vws.view,
		})
	}
	slices.SortFunc(views, func(a, b *ViewDiffEntry) int {
		return strings.Compare(objectID(a.SchemaName, a.ViewName), objectID(b.SchemaName, b.ViewName))
	})
	return views
}

func collectViewsToCreate(diff *SchemaDiff, rebuiltViews []viewWithSchema) []*ViewDiffEntry {
	var views []*ViewDiffEntry
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffCreate || viewDiff.Action == DiffAlter {
			views = append(views, viewDiff)
		}
	}
	for _, vws := range rebuiltViews {
		views = append(views, &ViewDiffEntry{
			Action:     DiffCreate,
			SchemaName: vws.schema,
			ViewName:   vws.view.Name,
			NewView:    vws.view,
		})
	}
	slices.SortFunc(views, func(a, b *ViewDiffEntry) int {
		return strings.Compare(objectID(a.SchemaName, a.ViewName), objectID(b.SchemaName, b.ViewName))
	})
	return views
}

// splitFunctionsBeforeViewDrops separates dropped functions that must run
// before or interleaved with view drops (schema-bound functions reading a
// dropped view, plus any dropped functions calling them) from functions that
// stay in the normal post-table function drop phase.
func splitFunctionsBeforeViewDrops(functionDrops []*FunctionDiffEntry, viewsToDrop []*ViewDiffEntry) (early, remaining []*FunctionDiffEntry) {
	if len(functionDrops) == 0 || len(viewsToDrop) == 0 {
		return nil, functionDrops
	}
	earlySet := make(map[string]bool)
	for _, fd := range functionDrops {
		if fd.OldFunction == nil || !isSchemaBound(fd.OldFunction.Definition) {
			continue
		}
		for _, vd := range viewsToDrop {
			if definitionReferencesTable(fd.OldFunction.Definition, fd.SchemaName, vd.SchemaName, vd.ViewName) {
				earlySet[objectID(fd.SchemaName, fd.FunctionName)] = true
				break
			}
		}
	}
	if len(earlySet) == 0 {
		return nil, functionDrops
	}
	// Transitively include any dropped function that calls an early-dropped function.
	for {
		added := false
		for _, fd := range functionDrops {
			id := objectID(fd.SchemaName, fd.FunctionName)
			if earlySet[id] || fd.OldFunction == nil {
				continue
			}
			for _, other := range functionDrops {
				otherID := objectID(other.SchemaName, other.FunctionName)
				if !earlySet[otherID] || otherID == id {
					continue
				}
				if definitionReferencesFunction(fd.OldFunction.Definition, fd.SchemaName, other.SchemaName, funcDiffBaseName(other)) {
					earlySet[id] = true
					added = true
					break
				}
			}
		}
		if !added {
			break
		}
	}
	for _, fd := range functionDrops {
		if earlySet[objectID(fd.SchemaName, fd.FunctionName)] {
			early = append(early, fd)
		} else {
			remaining = append(remaining, fd)
		}
	}
	return early, remaining
}

func dropViewsAndFunctionsInOrder(viewsToDrop []*ViewDiffEntry, funcDrops []*FunctionDiffEntry, buf *strings.Builder) {
	g := newGraph()
	viewMap := make(map[string]*ViewDiffEntry, len(viewsToDrop))
	viewNodeByFold := make(map[string]string, len(viewsToDrop))
	funcMap := make(map[string]*FunctionDiffEntry, len(funcDrops))
	for _, vd := range viewsToDrop {
		obj := objectID(vd.SchemaName, vd.ViewName)
		node := "V\x00" + obj
		g.addNode(node)
		viewMap[node] = vd
		viewNodeByFold[strings.ToLower(obj)] = node
	}
	for _, fd := range funcDrops {
		node := "F\x00" + objectID(fd.SchemaName, fd.FunctionName)
		g.addNode(node)
		funcMap[node] = fd
	}
	// View depends on View -> dependent view drops first.
	// View calls Function -> view drops before function.
	for _, vd := range viewsToDrop {
		if vd.OldView == nil || vd.OldView.Definition == "" {
			continue
		}
		vNode := "V\x00" + objectID(vd.SchemaName, vd.ViewName)
		for _, dep := range viewTableDependencies(vd.OldView.Definition, vd.SchemaName) {
			if targetNode, ok := viewNodeByFold[strings.ToLower(dep)]; ok && targetNode != vNode {
				g.addEdge(vNode, targetNode)
			}
		}
		for _, fd := range funcDrops {
			if expressionReferencesFunction(vd.OldView.Definition, vd.SchemaName, fd.SchemaName, funcDiffBaseName(fd)) {
				g.addEdge(vNode, "F\x00"+objectID(fd.SchemaName, fd.FunctionName))
			}
		}
	}
	// Schema-bound Function reads View -> function drops before view.
	// Function calls Function -> caller drops before callee.
	for _, fd := range funcDrops {
		if fd.OldFunction == nil {
			continue
		}
		fNode := "F\x00" + objectID(fd.SchemaName, fd.FunctionName)
		if isSchemaBound(fd.OldFunction.Definition) {
			for _, vd := range viewsToDrop {
				if definitionReferencesTable(fd.OldFunction.Definition, fd.SchemaName, vd.SchemaName, vd.ViewName) {
					g.addEdge(fNode, "V\x00"+objectID(vd.SchemaName, vd.ViewName))
				}
			}
		}
		for _, other := range funcDrops {
			otherNode := "F\x00" + objectID(other.SchemaName, other.FunctionName)
			if otherNode == fNode {
				continue
			}
			if definitionReferencesFunction(fd.OldFunction.Definition, fd.SchemaName, other.SchemaName, funcDiffBaseName(other)) {
				g.addEdge(fNode, otherNode)
			}
		}
	}
	order, err := g.topologicalSort()
	if err != nil {
		for _, fd := range orderedFunctionDrops(funcDrops) {
			fmt.Fprintf(buf, "DROP FUNCTION %s.%s;\nGO\n", quote(fd.SchemaName), quote(fd.FunctionName))
		}
		for _, vd := range viewsToDrop {
			fmt.Fprintf(buf, "DROP VIEW %s.%s;\nGO\n", quote(vd.SchemaName), quote(vd.ViewName))
		}
		return
	}
	for _, node := range order {
		if vd, ok := viewMap[node]; ok {
			fmt.Fprintf(buf, "DROP VIEW %s.%s;\nGO\n", quote(vd.SchemaName), quote(vd.ViewName))
		} else if fd, ok := funcMap[node]; ok {
			fmt.Fprintf(buf, "DROP FUNCTION %s.%s;\nGO\n", quote(fd.SchemaName), quote(fd.FunctionName))
		}
	}
}

func tableHasColumnCreateOrAlter(tableDiff *TableDiffEntry) bool {
	if tableDiff.Action != DiffAlter {
		return false
	}
	for _, colDiff := range tableDiff.Columns {
		if colDiff.Action == DiffCreate || colDiff.Action == DiffAlter {
			return true
		}
	}
	return false
}

// functionsDeferredPastTables returns the objectIDs of functions in funcCreates
// that must be created after table creation/alteration because they are
// schema-bound to a view in viewsToCreate, read a table whose columns are being
// created or altered, or transitively call such a function.
func functionsDeferredPastTables(diff *SchemaDiff, funcCreates []*FunctionDiffEntry, viewsToCreate []*ViewDiffEntry) map[string]bool {
	deps := make(map[string]bool)
	if len(funcCreates) == 0 {
		return deps
	}
	var alteredTables []*TableDiffEntry
	for _, tableDiff := range diff.Tables {
		if tableHasColumnCreateOrAlter(tableDiff) {
			alteredTables = append(alteredTables, tableDiff)
		}
	}
	for _, fd := range funcCreates {
		if fd.NewFunction == nil {
			continue
		}
		if isSchemaBound(fd.NewFunction.Definition) {
			for _, vd := range viewsToCreate {
				if definitionReferencesTable(fd.NewFunction.Definition, fd.SchemaName, vd.SchemaName, vd.ViewName) {
					deps[objectID(fd.SchemaName, fd.FunctionName)] = true
					break
				}
			}
		}
		if deps[objectID(fd.SchemaName, fd.FunctionName)] {
			continue
		}
		for _, td := range alteredTables {
			if definitionReferencesTable(fd.NewFunction.Definition, fd.SchemaName, td.SchemaName, td.TableName) {
				deps[objectID(fd.SchemaName, fd.FunctionName)] = true
				break
			}
		}
	}
	if len(deps) == 0 {
		return deps
	}
	for {
		added := false
		for _, fd := range funcCreates {
			id := objectID(fd.SchemaName, fd.FunctionName)
			if deps[id] || fd.NewFunction == nil {
				continue
			}
			for _, other := range funcCreates {
				otherID := objectID(other.SchemaName, other.FunctionName)
				if !deps[otherID] || otherID == id {
					continue
				}
				if definitionReferencesFunction(fd.NewFunction.Definition, fd.SchemaName, other.SchemaName, funcDiffBaseName(other)) {
					deps[id] = true
					added = true
					break
				}
			}
		}
		if !added {
			break
		}
	}
	return deps
}

func createdFunctionsDependOnCreatedViews(funcCreates []*FunctionDiffEntry, viewsToCreate []*ViewDiffEntry) bool {
	if len(funcCreates) == 0 || len(viewsToCreate) == 0 {
		return false
	}
	for _, fd := range funcCreates {
		if fd.NewFunction == nil || !isSchemaBound(fd.NewFunction.Definition) {
			continue
		}
		for _, vd := range viewsToCreate {
			if definitionReferencesTable(fd.NewFunction.Definition, fd.SchemaName, vd.SchemaName, vd.ViewName) {
				return true
			}
		}
	}
	return false
}

func createViewsAndFunctionsInOrder(viewsToCreate []*ViewDiffEntry, funcCreates []*FunctionDiffEntry, buf *strings.Builder) {
	g := newGraph()
	viewMap := make(map[string]*ViewDiffEntry, len(viewsToCreate))
	viewNodeByFold := make(map[string]string, len(viewsToCreate))
	funcMap := make(map[string]*FunctionDiffEntry, len(funcCreates))
	for _, vd := range viewsToCreate {
		obj := objectID(vd.SchemaName, vd.ViewName)
		node := "V\x00" + obj
		g.addNode(node)
		viewMap[node] = vd
		viewNodeByFold[strings.ToLower(obj)] = node
	}
	for _, fd := range funcCreates {
		node := "F\x00" + objectID(fd.SchemaName, fd.FunctionName)
		g.addNode(node)
		funcMap[node] = fd
	}
	// Prerequisite View -> dependent View; Function -> View calling it.
	for _, vd := range viewsToCreate {
		if vd.NewView == nil || vd.NewView.Definition == "" {
			continue
		}
		vNode := "V\x00" + objectID(vd.SchemaName, vd.ViewName)
		for _, dep := range viewTableDependencies(vd.NewView.Definition, vd.SchemaName) {
			if targetNode, ok := viewNodeByFold[strings.ToLower(dep)]; ok && targetNode != vNode {
				g.addEdge(targetNode, vNode)
			}
		}
		for _, fd := range funcCreates {
			if expressionReferencesFunction(vd.NewView.Definition, vd.SchemaName, fd.SchemaName, funcDiffBaseName(fd)) {
				g.addEdge("F\x00"+objectID(fd.SchemaName, fd.FunctionName), vNode)
			}
		}
	}
	// View -> schema-bound Function reading it; Callee Function -> Caller Function.
	for _, fd := range funcCreates {
		if fd.NewFunction == nil {
			continue
		}
		fNode := "F\x00" + objectID(fd.SchemaName, fd.FunctionName)
		if isSchemaBound(fd.NewFunction.Definition) {
			for _, vd := range viewsToCreate {
				if definitionReferencesTable(fd.NewFunction.Definition, fd.SchemaName, vd.SchemaName, vd.ViewName) {
					g.addEdge("V\x00"+objectID(vd.SchemaName, vd.ViewName), fNode)
				}
			}
		}
		for _, other := range funcCreates {
			otherNode := "F\x00" + objectID(other.SchemaName, other.FunctionName)
			if otherNode == fNode {
				continue
			}
			if definitionReferencesFunction(fd.NewFunction.Definition, fd.SchemaName, other.SchemaName, funcDiffBaseName(other)) {
				g.addEdge(otherNode, fNode)
			}
		}
	}
	writeViewEntry := func(viewDiff *ViewDiffEntry) {
		buf.WriteString(qualifyModuleDefinition(creationDefinition(viewDiff.NewView.Definition), viewDiff.SchemaName, false))
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
	order, err := g.topologicalSort()
	if err != nil {
		for _, fd := range orderedFunctionDiffs(funcCreates) {
			writeFunctionMigration(buf, fd)
		}
		for _, vd := range viewsToCreate {
			writeViewEntry(vd)
		}
		return
	}
	for _, node := range order {
		if vd, ok := viewMap[node]; ok {
			writeViewEntry(vd)
		} else if fd, ok := funcMap[node]; ok {
			writeFunctionMigration(buf, fd)
		}
	}
}

func functionExistsInCatalog(c *Catalog, schemaName, fnName string) bool {
	schema := c.GetSchema(schemaName)
	if schema == nil {
		return false
	}
	for _, fn := range schema.Functions {
		if strings.EqualFold(fn.Name, fnName) {
			return true
		}
	}
	return false
}

func collectFunctionDrops(diff *SchemaDiff, rebuiltFunctions []functionWithSchema) []*FunctionDiffEntry {
	var drops []*FunctionDiffEntry
	for _, funcDiff := range diff.Functions {
		if funcDiff.Action != DiffDrop || functionReplacementUsedByComputed(diff, funcDiff) {
			continue
		}
		drops = append(drops, funcDiff)
	}
	for _, fws := range rebuiltFunctions {
		name := fws.function.Signature
		if name == "" {
			name = fws.function.Name
		}
		drops = append(drops, &FunctionDiffEntry{
			Action:       DiffDrop,
			SchemaName:   fws.schema,
			FunctionName: name,
			OldFunction:  fws.function,
		})
	}
	return drops
}

func collectFunctionCreates(diff *SchemaDiff, rebuiltFunctions []functionWithSchema) []*FunctionDiffEntry {
	var creates []*FunctionDiffEntry
	for _, funcDiff := range diff.Functions {
		if funcDiff.Action != DiffCreate || functionReplacementUsedByComputed(diff, funcDiff) {
			continue
		}
		creates = append(creates, funcDiff)
	}
	for _, fws := range rebuiltFunctions {
		name := fws.function.Signature
		if name == "" {
			name = fws.function.Name
		}
		creates = append(creates, &FunctionDiffEntry{
			Action:       DiffCreate,
			SchemaName:   fws.schema,
			FunctionName: name,
			NewFunction:  fws.function,
		})
	}
	return creates
}

func collectDroppedTables(diff *SchemaDiff) []tableWithSchema {
	var dropped []tableWithSchema
	seen := make(map[string]bool)
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffDrop || tableDiff.OldTable == nil {
			continue
		}
		key := strings.ToLower(objectID(tableDiff.SchemaName, tableDiff.TableName))
		if !seen[key] {
			seen[key] = true
			dropped = append(dropped, tableWithSchema{schema: tableDiff.SchemaName, table: tableDiff.OldTable})
		}
	}
	for _, schemaDiff := range diff.Schemas {
		if schemaDiff.Action != DiffDrop || schemaDiff.OldSchema == nil {
			continue
		}
		for _, table := range schemaDiff.OldSchema.Tables {
			if table.GetSkipDump() {
				continue
			}
			key := strings.ToLower(objectID(schemaDiff.SchemaName, table.Name))
			if !seen[key] {
				seen[key] = true
				dropped = append(dropped, tableWithSchema{schema: schemaDiff.SchemaName, table: table})
			}
		}
	}
	return dropped
}

func droppedFunctionsDependOnDroppedTables(functionDrops []*FunctionDiffEntry, droppedTables []tableWithSchema) bool {
	if len(functionDrops) == 0 || len(droppedTables) == 0 {
		return false
	}
	for _, funcDiff := range functionDrops {
		if funcDiff.OldFunction == nil || !isSchemaBound(funcDiff.OldFunction.Definition) {
			continue
		}
		for _, tws := range droppedTables {
			if definitionReferencesTable(funcDiff.OldFunction.Definition, funcDiff.SchemaName, tws.schema, tws.table.Name) {
				return true
			}
		}
	}
	return false
}

// dropTablesAndFunctionsInOrder drops tables and functions in dependency order
// when at least one dropped schema-bound function reads a dropped table.
// Returns true if a DROP TABLE was emitted without a trailing GO.
func dropTablesAndFunctionsInOrder(droppedTables []tableWithSchema, functionDrops []*FunctionDiffEntry, hasTableAlterations bool, buf *strings.Builder) bool {
	g := newGraph()
	tableMap := make(map[string]tableWithSchema, len(droppedTables))
	funcMap := make(map[string]*FunctionDiffEntry, len(functionDrops))
	for _, tws := range droppedTables {
		node := "T\x00" + objectID(tws.schema, tws.table.Name)
		g.addNode(node)
		tableMap[node] = tws
	}
	for _, funcDiff := range functionDrops {
		node := "F\x00" + objectID(funcDiff.SchemaName, funcDiff.FunctionName)
		g.addNode(node)
		funcMap[node] = funcDiff
	}
	// Function reads Table (WITH SCHEMABINDING) -> Function drops before Table.
	for _, funcDiff := range functionDrops {
		if funcDiff.OldFunction == nil || !isSchemaBound(funcDiff.OldFunction.Definition) {
			continue
		}
		fNode := "F\x00" + objectID(funcDiff.SchemaName, funcDiff.FunctionName)
		for _, tws := range droppedTables {
			if definitionReferencesTable(funcDiff.OldFunction.Definition, funcDiff.SchemaName, tws.schema, tws.table.Name) {
				g.addEdge(fNode, "T\x00"+objectID(tws.schema, tws.table.Name))
			}
		}
	}
	// Table references Function -> Table drops before Function.
	for _, tws := range droppedTables {
		tNode := "T\x00" + objectID(tws.schema, tws.table.Name)
		for _, funcDiff := range functionDrops {
			if tableReferencesFunction(tws.schema, tws.table, funcDiff.SchemaName, funcDiffBaseName(funcDiff)) {
				g.addEdge(tNode, "F\x00"+objectID(funcDiff.SchemaName, funcDiff.FunctionName))
			}
		}
	}
	// Function calls Function -> caller drops before callee.
	for _, funcDiff := range functionDrops {
		if funcDiff.OldFunction == nil {
			continue
		}
		fNode := "F\x00" + objectID(funcDiff.SchemaName, funcDiff.FunctionName)
		for _, other := range functionDrops {
			otherNode := "F\x00" + objectID(other.SchemaName, other.FunctionName)
			if otherNode == fNode {
				continue
			}
			if definitionReferencesFunction(funcDiff.OldFunction.Definition, funcDiff.SchemaName, other.SchemaName, funcDiffBaseName(other)) {
				g.addEdge(fNode, otherNode)
			}
		}
	}

	needAlterGo := hasTableAlterations
	endedWithTableDrop := false
	order, err := g.topologicalSort()
	if err == nil {
		for _, node := range order {
			if tws, ok := tableMap[node]; ok {
				if needAlterGo {
					buf.WriteString("GO\n")
					needAlterGo = false
				}
				fmt.Fprintf(buf, "DROP TABLE %s.%s;\n", quote(tws.schema), quote(tws.table.Name))
				endedWithTableDrop = true
			} else if funcDiff, ok := funcMap[node]; ok {
				needAlterGo = false
				fmt.Fprintf(buf, "DROP FUNCTION %s.%s;\nGO\n", quote(funcDiff.SchemaName), quote(funcDiff.FunctionName))
				endedWithTableDrop = false
			}
		}
		return endedWithTableDrop
	}

	// Cycle between dropped table constraint/computed column and schema-bound
	// function reading that table: drop independent tables first, strip
	// function-referencing constraints/computed columns from bound tables, drop
	// functions, then drop the bound tables.
	var independentTables, boundTables []tableWithSchema
	for _, tws := range droppedTables {
		bound := false
		for _, funcDiff := range functionDrops {
			if funcDiff.OldFunction != nil && isSchemaBound(funcDiff.OldFunction.Definition) &&
				definitionReferencesTable(funcDiff.OldFunction.Definition, funcDiff.SchemaName, tws.schema, tws.table.Name) {
				bound = true
				break
			}
		}
		if bound {
			boundTables = append(boundTables, tws)
		} else {
			independentTables = append(independentTables, tws)
		}
	}
	for _, tws := range independentTables {
		if needAlterGo {
			buf.WriteString("GO\n")
			needAlterGo = false
		}
		fmt.Fprintf(buf, "DROP TABLE %s.%s;\n", quote(tws.schema), quote(tws.table.Name))
	}
	for _, tws := range boundTables {
		for _, col := range tws.table.Columns {
			if col.Default == "" || col.DefaultConstraintName == "" {
				continue
			}
			for _, funcDiff := range functionDrops {
				if expressionReferencesFunction(col.Default, tws.schema, funcDiff.SchemaName, funcDiffBaseName(funcDiff)) {
					fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tws.schema), quote(tws.table.Name), quote(col.DefaultConstraintName))
					break
				}
			}
		}
		var droppedCompCols map[string]bool
		for _, col := range tws.table.Columns {
			if !isComputedColumnType(col.Type) {
				continue
			}
			for _, funcDiff := range functionDrops {
				if expressionReferencesFunction(col.Type, tws.schema, funcDiff.SchemaName, funcDiffBaseName(funcDiff)) {
					if droppedCompCols == nil {
						droppedCompCols = make(map[string]bool)
					}
					droppedCompCols[strings.ToLower(col.Name)] = true
					break
				}
			}
		}
		if len(droppedCompCols) > 0 {
			// The table itself is dropped immediately after the functions, so
			// dropping all its checks and indexes avoids leaving any CHECK,
			// INCLUDE, or filtered-index dependency on a computed column.
			for _, check := range tws.table.CheckConstraints {
				if check.Name != "" {
					fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tws.schema), quote(tws.table.Name), quote(check.Name))
				}
			}
			droppedAnyIndex := false
			emitIndexDrop := func(idx *metadata.IndexMetadata) {
				if idx.Name == "" {
					return
				}
				if idx.IsConstraint {
					fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tws.schema), quote(tws.table.Name), quote(idx.Name))
				} else {
					fmt.Fprintf(buf, "DROP INDEX %s ON %s.%s;\n", quote(idx.Name), quote(tws.schema), quote(tws.table.Name))
				}
				droppedAnyIndex = true
			}
			for _, idx := range tws.table.Indexes {
				if idx.Type == "XML" && strings.HasPrefix(idx.Definition, "USING XML INDEX ") {
					emitIndexDrop(idx)
				}
			}
			for _, idx := range tws.table.Indexes {
				if idx.Primary || (idx.Type == "XML" && strings.HasPrefix(idx.Definition, "USING XML INDEX ")) {
					continue
				}
				emitIndexDrop(idx)
			}
			for _, idx := range tws.table.Indexes {
				if idx.Primary {
					emitIndexDrop(idx)
				}
			}
			if droppedAnyIndex {
				buf.WriteString("GO\n")
			}
			for _, col := range tws.table.Columns {
				if droppedCompCols[strings.ToLower(col.Name)] {
					fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP COLUMN %s;\n", quote(tws.schema), quote(tws.table.Name), quote(col.Name))
				}
			}
		} else {
			for _, check := range tws.table.CheckConstraints {
				if check.Name == "" {
					continue
				}
				for _, funcDiff := range functionDrops {
					if expressionReferencesFunction(check.Expression, tws.schema, funcDiff.SchemaName, funcDiffBaseName(funcDiff)) {
						fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(tws.schema), quote(tws.table.Name), quote(check.Name))
						break
					}
				}
			}
		}
	}
	for _, funcDiff := range orderedFunctionDrops(functionDrops) {
		fmt.Fprintf(buf, "DROP FUNCTION %s.%s;\nGO\n", quote(funcDiff.SchemaName), quote(funcDiff.FunctionName))
	}
	for _, tws := range boundTables {
		fmt.Fprintf(buf, "DROP TABLE %s.%s;\n", quote(tws.schema), quote(tws.table.Name))
		endedWithTableDrop = true
	}
	return endedWithTableDrop
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

func dropFunctionDefaults(from *Catalog, diff *SchemaDiff, functionDrops []*FunctionDiffEntry, alreadyDropped map[string]bool, buf *strings.Builder) (map[string]bool, []functionDefaultRef) {
	removed := make(map[string]bool)
	var restore []functionDefaultRef
	if from == nil || len(functionDrops) == 0 {
		return removed, restore
	}
	droppedTables := make(map[string]bool)
	droppedColumns := make(map[string]bool)
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffDrop {
			droppedTables[strings.ToLower(objectID(tableDiff.SchemaName, tableDiff.TableName))] = true
		}
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action == DiffDrop && colDiff.OldColumn != nil {
				droppedColumns[strings.ToLower(objectID(tableDiff.SchemaName, tableDiff.TableName)+"\x00"+colDiff.OldColumn.Name)] = true
			}
		}
	}
	for _, schema := range from.Metadata().Schemas {
		for _, table := range schema.Tables {
			tKey := strings.ToLower(objectID(schema.Name, table.Name))
			if droppedTables[tKey] {
				continue
			}
			for _, column := range table.Columns {
				if column.Default == "" || column.DefaultConstraintName == "" || droppedColumns[tKey+"\x00"+strings.ToLower(column.Name)] {
					continue
				}
				if alreadyDropped[schema.Name+"."+column.DefaultConstraintName] {
					continue
				}
				var matched []string
				for _, funcDiff := range functionDrops {
					if expressionReferencesFunction(column.Default, schema.Name, funcDiff.SchemaName, funcDiffBaseName(funcDiff)) {
						matched = append(matched, strings.ToLower(objectID(funcDiff.SchemaName, funcDiffBaseName(funcDiff))))
					}
				}
				if len(matched) > 0 {
					fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(schema.Name), quote(table.Name), quote(column.DefaultConstraintName))
					removed[schema.Name+"."+column.DefaultConstraintName] = true
					restore = append(restore, functionDefaultRef{
						schema: schema.Name, table: table.Name, column: column.Name,
						constraint: column.DefaultConstraintName, expression: column.Default,
						functions: matched,
					})
				}
			}
		}
	}
	return removed, restore
}

func dropFunctionChecks(from *Catalog, diff *SchemaDiff, functionDrops []*FunctionDiffEntry, alreadyDropped map[string]bool, buf *strings.Builder) (map[string]bool, []functionCheckRef) {
	removed := make(map[string]bool)
	var restore []functionCheckRef
	if from == nil || len(functionDrops) == 0 {
		return removed, restore
	}
	droppedTables := make(map[string]bool)
	droppedChecks := make(map[string]bool)
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffDrop {
			droppedTables[strings.ToLower(objectID(tableDiff.SchemaName, tableDiff.TableName))] = true
			continue
		}
		if tableDiff.Action != DiffAlter {
			continue
		}
		for _, checkDiff := range tableDiff.Checks {
			if checkDiff.Action == DiffDrop && checkDiff.OldCheckConstraint != nil {
				droppedChecks[checkKey(tableDiff.SchemaName, tableDiff.TableName, checkDiff.OldCheckConstraint.Name)] = true
			}
		}
	}
	for _, schema := range from.Metadata().Schemas {
		for _, table := range schema.Tables {
			tKey := strings.ToLower(objectID(schema.Name, table.Name))
			if droppedTables[tKey] {
				continue
			}
			for _, check := range table.CheckConstraints {
				if check.Expression == "" || check.Name == "" {
					continue
				}
				cKey := checkKey(schema.Name, table.Name, check.Name)
				if droppedChecks[cKey] || alreadyDropped[cKey] {
					continue
				}
				var matched []string
				for _, funcDiff := range functionDrops {
					if expressionReferencesFunction(check.Expression, schema.Name, funcDiff.SchemaName, funcDiffBaseName(funcDiff)) {
						matched = append(matched, strings.ToLower(objectID(funcDiff.SchemaName, funcDiffBaseName(funcDiff))))
					}
				}
				if len(matched) > 0 {
					fmt.Fprintf(buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(schema.Name), quote(table.Name), quote(check.Name))
					removed[cKey] = true
					restore = append(restore, functionCheckRef{
						schema: schema.Name, table: table.Name,
						constraint: check.Name, expression: check.Expression,
						functions: matched,
					})
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
	if !functionUsedByComputed(diff, funcDiff.SchemaName, funcDiffBaseName(funcDiff)) {
		return false
	}
	for _, other := range diff.Functions {
		if other.Action == DiffDrop && strings.EqualFold(other.SchemaName, funcDiff.SchemaName) && strings.EqualFold(other.FunctionName, funcDiff.FunctionName) {
			return true
		}
	}
	return false
}

func functionUsedByComputed(diff *SchemaDiff, fnSchema, fnName string) bool {
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action != DiffAlter || tableDiff.OldTable == nil {
			continue
		}
		droppedCols := make(map[string]bool)
		for _, colDiff := range tableDiff.Columns {
			if colDiff.Action == DiffDrop && colDiff.OldColumn != nil {
				droppedCols[strings.ToLower(colDiff.OldColumn.Name)] = true
			}
		}
		for _, column := range tableDiff.OldTable.Columns {
			if droppedCols[strings.ToLower(column.Name)] {
				continue
			}
			if isComputedColumnType(column.Type) && expressionReferencesFunction(column.Type, tableDiff.SchemaName, fnSchema, fnName) {
				return true
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
		id := objectID(function.SchemaName, function.FunctionName)
		graph.addNode(id)
		byName[id] = function
	}
	for _, function := range functions {
		if function.Action != DiffDrop || function.OldFunction == nil {
			continue
		}
		id := objectID(function.SchemaName, function.FunctionName)
		for _, other := range functions {
			if other.Action != DiffDrop || other == function {
				continue
			}
			if definitionReferencesFunction(function.OldFunction.Definition, function.SchemaName, other.SchemaName, funcDiffBaseName(other)) {
				graph.addEdge(id, objectID(other.SchemaName, other.FunctionName))
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
		id := objectID(function.SchemaName, function.FunctionName)
		graph.addNode(id)
		byName[id] = function
	}
	for _, function := range functions {
		if function.Action != DiffCreate || function.NewFunction == nil {
			continue
		}
		id := objectID(function.SchemaName, function.FunctionName)
		for _, other := range functions {
			if other.Action != DiffCreate || other == function {
				continue
			}
			if definitionReferencesFunction(function.NewFunction.Definition, function.SchemaName, other.SchemaName, funcDiffBaseName(other)) {
				graph.addEdge(objectID(other.SchemaName, other.FunctionName), id)
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

func expandMigrationFunctionDeps(functions []*FunctionDiffEntry, seeds map[string]bool) map[string]bool {
	if len(seeds) == 0 {
		return seeds
	}
	byID := make(map[string]*FunctionDiffEntry, len(functions))
	for _, fn := range functions {
		if fn.Action == DiffCreate {
			byID[objectID(fn.SchemaName, fn.FunctionName)] = fn
		}
	}
	result := make(map[string]bool, len(seeds))
	var queue []string
	for id := range seeds {
		result[id] = true
		queue = append(queue, id)
	}
	slices.Sort(queue)
	for len(queue) > 0 {
		curID := queue[0]
		queue = queue[1:]
		cur := byID[curID]
		if cur == nil || cur.NewFunction == nil {
			continue
		}
		for _, other := range functions {
			if other.Action != DiffCreate {
				continue
			}
			otherID := objectID(other.SchemaName, other.FunctionName)
			if result[otherID] || otherID == curID {
				continue
			}
			if definitionReferencesFunction(cur.NewFunction.Definition, cur.SchemaName, other.SchemaName, funcDiffBaseName(other)) {
				result[otherID] = true
				queue = append(queue, otherID)
			}
		}
	}
	return result
}

func migrationTableFunctions(diff *SchemaDiff, creates []*FunctionDiffEntry) map[string]bool {
	type exprCtx struct {
		schema string
		expr   string
	}
	var expressions []exprCtx
	for _, tableDiff := range diff.Tables {
		if tableDiff.NewTable == nil {
			continue
		}
		for _, column := range tableDiff.NewTable.Columns {
			if column.Default != "" {
				expressions = append(expressions, exprCtx{schema: tableDiff.SchemaName, expr: column.Default})
			}
			if isComputedColumnType(column.Type) {
				expressions = append(expressions, exprCtx{schema: tableDiff.SchemaName, expr: column.Type})
			}
		}
		for _, check := range tableDiff.NewTable.CheckConstraints {
			if check.Expression != "" {
				expressions = append(expressions, exprCtx{schema: tableDiff.SchemaName, expr: check.Expression})
			}
		}
	}
	referenced := make(map[string]bool)
	for _, funcDiff := range creates {
		if funcDiff.Action != DiffCreate {
			continue
		}
		name := funcDiffBaseName(funcDiff)
		for _, item := range expressions {
			if expressionReferencesFunction(item.expr, item.schema, funcDiff.SchemaName, name) {
				referenced[objectID(funcDiff.SchemaName, funcDiff.FunctionName)] = true
				break
			}
		}
	}
	return expandMigrationFunctionDeps(creates, referenced)
}

func createdFunctionsDependOnCreatedTables(diff *SchemaDiff, creates []*FunctionDiffEntry, tableFunctions map[string]bool) bool {
	for _, funcDiff := range creates {
		if funcDiff.NewFunction == nil || !tableFunctions[objectID(funcDiff.SchemaName, funcDiff.FunctionName)] {
			continue
		}
		if !isSchemaBound(funcDiff.NewFunction.Definition) {
			continue
		}
		for _, tableDiff := range diff.Tables {
			if tableDiff.Action != DiffCreate || tableDiff.NewTable == nil {
				continue
			}
			if definitionReferencesTable(funcDiff.NewFunction.Definition, funcDiff.SchemaName, tableDiff.SchemaName, tableDiff.TableName) {
				return true
			}
		}
	}
	return false
}

func exprReferencesAnyDeferredFunc(expr, ownerSchema string, creates []*FunctionDiffEntry, deferredFuncCreates map[string]bool) bool {
	for _, fd := range creates {
		if !deferredFuncCreates[objectID(fd.SchemaName, fd.FunctionName)] {
			continue
		}
		if expressionReferencesFunction(expr, ownerSchema, fd.SchemaName, funcDiffBaseName(fd)) {
			return true
		}
	}
	return false
}

func createTablesAndFunctionsInOrder(diff *SchemaDiff, creates []*FunctionDiffEntry, tableFunctions, deferredDefaultCols, deferredChecks map[string]bool, buf *strings.Builder) {
	var tablesToCreate []*TableDiffEntry
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffCreate {
			tablesToCreate = append(tablesToCreate, tableDiff)
		}
	}
	var funcsToCreate []*FunctionDiffEntry
	for _, funcDiff := range creates {
		if tableFunctions[objectID(funcDiff.SchemaName, funcDiff.FunctionName)] {
			funcsToCreate = append(funcsToCreate, funcDiff)
		}
	}
	g := newGraph()
	tableMap := make(map[string]*TableDiffEntry, len(tablesToCreate))
	tableNodeByObj := make(map[string]string, len(tablesToCreate))
	funcMap := make(map[string]*FunctionDiffEntry, len(funcsToCreate))
	for _, td := range tablesToCreate {
		obj := objectID(td.SchemaName, td.TableName)
		node := "T\x00" + obj
		g.addNode(node)
		tableMap[node] = td
		tableNodeByObj[obj] = node
	}
	for _, fd := range funcsToCreate {
		node := "F\x00" + objectID(fd.SchemaName, fd.FunctionName)
		g.addNode(node)
		funcMap[node] = fd
	}
	for _, td := range tablesToCreate {
		tNode := "T\x00" + objectID(td.SchemaName, td.TableName)
		if td.NewTable == nil {
			continue
		}
		for _, fk := range td.NewTable.ForeignKeys {
			if refNode, ok := tableNodeByObj[objectID(fk.ReferencedSchema, fk.ReferencedTable)]; ok && refNode != tNode {
				g.addEdge(refNode, tNode)
			}
		}
		for _, fd := range funcsToCreate {
			if tableReferencesFunction(td.SchemaName, td.NewTable, fd.SchemaName, funcDiffBaseName(fd)) {
				g.addEdge("F\x00"+objectID(fd.SchemaName, fd.FunctionName), tNode)
			}
		}
	}
	for _, fd := range funcsToCreate {
		if fd.NewFunction == nil {
			continue
		}
		fNode := "F\x00" + objectID(fd.SchemaName, fd.FunctionName)
		for _, other := range funcsToCreate {
			otherNode := "F\x00" + objectID(other.SchemaName, other.FunctionName)
			if otherNode == fNode {
				continue
			}
			if definitionReferencesFunction(fd.NewFunction.Definition, fd.SchemaName, other.SchemaName, funcDiffBaseName(other)) {
				g.addEdge(otherNode, fNode)
			}
		}
		if isSchemaBound(fd.NewFunction.Definition) {
			for _, td := range tablesToCreate {
				if definitionReferencesTable(fd.NewFunction.Definition, fd.SchemaName, td.SchemaName, td.TableName) {
					g.addEdge("T\x00"+objectID(td.SchemaName, td.TableName), fNode)
				}
			}
		}
	}
	order, err := g.topologicalSort()
	if err != nil {
		for _, fd := range orderedFunctionDiffs(funcsToCreate) {
			writeFunctionMigration(buf, fd)
		}
		createTablesInOrder(diff, deferredDefaultCols, deferredChecks, buf)
		return
	}
	needGoBeforeFn := false
	for _, node := range order {
		if td, ok := tableMap[node]; ok {
			buf.WriteString(generateCreateTable(td.SchemaName, td.TableName, td.NewTable, deferredDefaultCols, deferredChecks))
			buf.WriteString("\n")
			needGoBeforeFn = true
		} else if fd, ok := funcMap[node]; ok {
			if needGoBeforeFn {
				buf.WriteString("GO\n")
			}
			writeFunctionMigration(buf, fd)
			needGoBeforeFn = false
		}
	}
}

func migrationEarlyFunctions(diff *SchemaDiff, rebuiltViews []viewWithSchema, creates []*FunctionDiffEntry) map[string]bool {
	type viewCtx struct {
		schema string
		def    string
	}
	var definitions []viewCtx
	for _, viewDiff := range diff.Views {
		if viewDiff.NewView != nil {
			definitions = append(definitions, viewCtx{schema: viewDiff.SchemaName, def: viewDiff.NewView.Definition})
		}
	}
	for _, vws := range rebuiltViews {
		definitions = append(definitions, viewCtx{schema: vws.schema, def: vws.view.Definition})
	}
	referenced := make(map[string]bool)
	for _, funcDiff := range creates {
		if funcDiff.Action != DiffCreate {
			continue
		}
		name := funcDiffBaseName(funcDiff)
		for _, item := range definitions {
			if expressionReferencesFunction(item.def, item.schema, funcDiff.SchemaName, name) {
				referenced[objectID(funcDiff.SchemaName, funcDiff.FunctionName)] = true
				break
			}
		}
	}
	return expandMigrationFunctionDeps(creates, referenced)
}

func writeFunctionMigration(buf *strings.Builder, funcDiff *FunctionDiffEntry) {
	if funcDiff.NewFunction == nil {
		return
	}
	buf.WriteString(qualifyModuleDefinition(creationDefinition(funcDiff.NewFunction.Definition), funcDiff.SchemaName, false))
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

func generateColumnDefinitionFiltered(column *metadata.ColumnMetadata, includeDefault bool) string {
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
	if includeDefault && column.GetDefault() != "" {
		if column.DefaultConstraintName != "" {
			fmt.Fprintf(&buf, " CONSTRAINT %s DEFAULT %s", quote(column.DefaultConstraintName), column.GetDefault())
		} else {
			buf.WriteString(" DEFAULT ")
			buf.WriteString(column.GetDefault())
		}
	}
	return buf.String()
}

func generateAlterColumn(schemaName, tableName string, colDiff *ColumnDiffEntry, droppedDefaults, deferredDefaultCols map[string]bool) string {
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
	if !expressionsEqual(oldDefault, newDefault) || (droppedDefaultForAlter && newDefault != "") {
		if !droppedDefaultForAlter && oldDefault != "" && colDiff.OldColumn.DefaultConstraintName != "" && !droppedDefaults[schemaName+"."+colDiff.OldColumn.DefaultConstraintName] {
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s DROP CONSTRAINT %s;\n", quote(schemaName), quote(tableName), quote(colDiff.OldColumn.DefaultConstraintName))
		}
		if newDefault != "" && !deferredDefaultCols[checkKey(schemaName, tableName, colDiff.NewColumn.Name)] {
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
	if index.Definition == "SELECTIVE" || strings.HasPrefix(index.Definition, "SELECTIVE ") {
		rest := strings.TrimSpace(strings.TrimPrefix(index.Definition, "SELECTIVE"))
		if rest != "" {
			return fmt.Sprintf("CREATE SELECTIVE XML INDEX %s ON %s.%s (%s) %s",
				quote(index.Name), quote(schemaName), quote(tableName), column, rest)
		}
		return fmt.Sprintf("CREATE SELECTIVE XML INDEX %s ON %s.%s (%s)",
			quote(index.Name), quote(schemaName), quote(tableName), column)
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

func dropViewsInOrder(diff *SchemaDiff, rebuiltViews []viewWithSchema, buf *strings.Builder) {
	graph := newGraph()
	viewMap := make(map[string]*ViewDiffEntry)
	viewIDByFold := make(map[string]string)
	var viewsToProcess []*ViewDiffEntry
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffDrop || viewDiff.Action == DiffAlter {
			viewsToProcess = append(viewsToProcess, viewDiff)
		}
	}
	for _, vws := range rebuiltViews {
		viewsToProcess = append(viewsToProcess, &ViewDiffEntry{
			Action:     DiffDrop,
			SchemaName: vws.schema,
			ViewName:   vws.view.Name,
			OldView:    vws.view,
		})
	}
	slices.SortFunc(viewsToProcess, func(a, b *ViewDiffEntry) int {
		return strings.Compare(objectID(a.SchemaName, a.ViewName), objectID(b.SchemaName, b.ViewName))
	})
	for _, viewDiff := range viewsToProcess {
		viewID := objectID(viewDiff.SchemaName, viewDiff.ViewName)
		graph.addNode(viewID)
		viewMap[viewID] = viewDiff
		viewIDByFold[strings.ToLower(viewID)] = viewID
	}
	for _, viewDiff := range viewsToProcess {
		viewID := objectID(viewDiff.SchemaName, viewDiff.ViewName)
		if viewDiff.OldView == nil || viewDiff.OldView.Definition == "" {
			continue
		}
		for _, dep := range viewTableDependencies(viewDiff.OldView.Definition, viewDiff.SchemaName) {
			if targetID, ok := viewIDByFold[strings.ToLower(dep)]; ok && targetID != viewID {
				graph.addEdge(viewID, targetID)
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

func createTablesInOrder(diff *SchemaDiff, deferredDefaultCols, deferredChecks map[string]bool, buf *strings.Builder) {
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
			buf.WriteString(generateCreateTable(tableDiff.SchemaName, tableDiff.TableName, tableDiff.NewTable, deferredDefaultCols, deferredChecks))
			buf.WriteString("\n")
		}
		return
	}
	for _, tableID := range ordered {
		if tableDiff, ok := tableMap[tableID]; ok {
			buf.WriteString(generateCreateTable(tableDiff.SchemaName, tableDiff.TableName, tableDiff.NewTable, deferredDefaultCols, deferredChecks))
			buf.WriteString("\n")
		}
	}
}

func createViewsInOrder(diff *SchemaDiff, rebuiltViews []viewWithSchema, buf *strings.Builder) {
	graph := newGraph()
	viewMap := make(map[string]*ViewDiffEntry)
	viewIDByFold := make(map[string]string)
	var viewsToProcess []*ViewDiffEntry
	for _, viewDiff := range diff.Views {
		if viewDiff.Action == DiffCreate || viewDiff.Action == DiffAlter {
			viewsToProcess = append(viewsToProcess, viewDiff)
		}
	}
	for _, vws := range rebuiltViews {
		viewsToProcess = append(viewsToProcess, &ViewDiffEntry{
			Action:     DiffCreate,
			SchemaName: vws.schema,
			ViewName:   vws.view.Name,
			NewView:    vws.view,
		})
	}
	slices.SortFunc(viewsToProcess, func(a, b *ViewDiffEntry) int {
		return strings.Compare(objectID(a.SchemaName, a.ViewName), objectID(b.SchemaName, b.ViewName))
	})
	for _, viewDiff := range viewsToProcess {
		viewID := objectID(viewDiff.SchemaName, viewDiff.ViewName)
		graph.addNode(viewID)
		viewMap[viewID] = viewDiff
		viewIDByFold[strings.ToLower(viewID)] = viewID
	}
	for _, viewDiff := range viewsToProcess {
		viewID := objectID(viewDiff.SchemaName, viewDiff.ViewName)
		if viewDiff.NewView == nil || viewDiff.NewView.Definition == "" {
			continue
		}
		for _, dep := range viewTableDependencies(viewDiff.NewView.Definition, viewDiff.SchemaName) {
			if targetID, ok := viewIDByFold[strings.ToLower(dep)]; ok && targetID != viewID {
				graph.addEdge(targetID, viewID)
			}
		}
	}
	writeView := func(viewDiff *ViewDiffEntry) {
		buf.WriteString(qualifyModuleDefinition(creationDefinition(viewDiff.NewView.Definition), viewDiff.SchemaName, false))
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

func hasTableDrops(diff *SchemaDiff) bool {
	for _, tableDiff := range diff.Tables {
		if tableDiff.Action == DiffDrop {
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
