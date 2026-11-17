package catalog

import (
	"regexp"
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
)

// DiffAction describes what happened to an object between two catalog states.
type DiffAction string

const (
	DiffActionCreate DiffAction = "CREATE"
	DiffActionDrop   DiffAction = "DROP"
	DiffActionAlter  DiffAction = "ALTER"
)

type SchemaDiffEntry struct {
	Action     DiffAction
	SchemaName string
	OldSchema  *metadatapb.SchemaMetadata
	NewSchema  *metadatapb.SchemaMetadata
}

type TableDiff struct {
	Action     DiffAction
	SchemaName string
	TableName  string
	OldTable   *metadatapb.TableMetadata
	NewTable   *metadatapb.TableMetadata

	ColumnChanges          []*ColumnDiff
	IndexChanges           []*IndexDiff
	ForeignKeyChanges      []*ForeignKeyDiff
	CheckConstraintChanges []*CheckConstraintDiff
	TriggerChanges         []*TriggerDiff
}

type ColumnDiff struct {
	Action    DiffAction
	OldColumn *metadatapb.ColumnMetadata
	NewColumn *metadatapb.ColumnMetadata
}

type IndexDiff struct {
	Action   DiffAction
	OldIndex *metadatapb.IndexMetadata
	NewIndex *metadatapb.IndexMetadata
}

type ForeignKeyDiff struct {
	Action        DiffAction
	OldForeignKey *metadatapb.ForeignKeyMetadata
	NewForeignKey *metadatapb.ForeignKeyMetadata
}

type CheckConstraintDiff struct {
	Action             DiffAction
	OldCheckConstraint *metadatapb.CheckConstraintMetadata
	NewCheckConstraint *metadatapb.CheckConstraintMetadata
}

type TriggerDiff struct {
	Action      DiffAction
	SchemaName  string
	TableName   string
	TriggerName string
	OldTrigger  *metadatapb.TriggerMetadata
	NewTrigger  *metadatapb.TriggerMetadata
}

type ViewDiff struct {
	Action     DiffAction
	SchemaName string
	ViewName   string
	OldView    *metadatapb.ViewMetadata
	NewView    *metadatapb.ViewMetadata
}

type MaterializedViewDiff struct {
	Action               DiffAction
	SchemaName           string
	MaterializedViewName string
	OldMaterializedView  *metadatapb.MaterializedViewMetadata
	NewMaterializedView  *metadatapb.MaterializedViewMetadata
	IndexChanges         []*IndexDiff
}

type FunctionDiff struct {
	Action       DiffAction
	SchemaName   string
	FunctionName string
	OldFunction  *metadatapb.FunctionMetadata
	NewFunction  *metadatapb.FunctionMetadata
}

type ProcedureDiff struct {
	Action        DiffAction
	SchemaName    string
	ProcedureName string
	OldProcedure  *metadatapb.ProcedureMetadata
	NewProcedure  *metadatapb.ProcedureMetadata
}

type SequenceDiff struct {
	Action       DiffAction
	SchemaName   string
	SequenceName string
	OldSequence  *metadatapb.SequenceMetadata
	NewSequence  *metadatapb.SequenceMetadata
}

// MetadataDiff is the full diff between two metadata snapshots.
type MetadataDiff struct {
	DatabaseName string

	OldDatabase *metadatapb.DatabaseSchemaMetadata
	NewDatabase *metadatapb.DatabaseSchemaMetadata

	SchemaChanges           []*SchemaDiffEntry
	TableChanges            []*TableDiff
	ViewChanges             []*ViewDiff
	MaterializedViewChanges []*MaterializedViewDiff
	FunctionChanges         []*FunctionDiff
	ProcedureChanges        []*ProcedureDiff
	SequenceChanges         []*SequenceDiff

	preAltered   map[*ColumnDiff]bool
	preAdded     map[*ColumnDiff]bool
	alterHandled map[*TableDiff]bool
}

// SchemaDiff is the diff Diff returns.
type SchemaDiff = MetadataDiff

// IsEmpty reports whether the diff changes nothing.
func (d *SchemaDiff) IsEmpty() bool {
	if d == nil {
		return true
	}
	return len(d.SchemaChanges) == 0 &&
		len(d.TableChanges) == 0 &&
		len(d.ViewChanges) == 0 &&
		len(d.MaterializedViewChanges) == 0 &&
		len(d.FunctionChanges) == 0 &&
		len(d.ProcedureChanges) == 0 &&
		len(d.SequenceChanges) == 0
}

// Diff compares two catalog states and returns every difference.
func Diff(from, to *Catalog) *SchemaDiff {
	oldMeta, newMeta := from.Metadata(), to.Metadata()
	if oldMeta == nil || newMeta == nil {
		return &MetadataDiff{}
	}
	diff := &MetadataDiff{
		DatabaseName: newMeta.Name,
		OldDatabase:  oldMeta,
		NewDatabase:  newMeta,
	}

	oldSchemas := schemaByName(oldMeta)
	newSchemas := schemaByName(newMeta)
	matchedOld := map[*metadatapb.SchemaMetadata]bool{}
	for name, newSchema := range newSchemas {
		oldSchema := oldSchemas[name]
		if oldSchema == nil {
			switch {
			case name != "" && name == oldMeta.GetName():
				oldSchema = oldSchemas[""]
			case name == "" && oldMeta.GetName() != "":
				oldSchema = oldSchemas[oldMeta.GetName()]
			}
		}
		if oldSchema == nil {
			diff.SchemaChanges = append(diff.SchemaChanges, &SchemaDiffEntry{Action: DiffActionCreate, SchemaName: name, NewSchema: newSchema})
			addNewSchemaObjects(diff, name, newSchema)
			continue
		}
		matchedOld[oldSchema] = true
		compareSchemaObjects(diff, name, oldSchema, newSchema)
	}
	for name, oldSchema := range oldSchemas {
		if matchedOld[oldSchema] {
			continue
		}
		if name == "" && newMeta.GetName() != "" && newSchemas[newMeta.GetName()] != nil {
			continue
		}
		diff.SchemaChanges = append(diff.SchemaChanges, &SchemaDiffEntry{Action: DiffActionDrop, SchemaName: name, OldSchema: oldSchema})
	}

	sortDiff(diff)
	return diff
}

func schemaByName(meta *metadatapb.DatabaseSchemaMetadata) map[string]*metadatapb.SchemaMetadata {
	out := make(map[string]*metadatapb.SchemaMetadata, len(meta.GetSchemas()))
	for _, schema := range meta.GetSchemas() {
		out[schema.GetName()] = schema
	}
	return out
}

func tableByName(schema *metadatapb.SchemaMetadata) map[string]*metadatapb.TableMetadata {
	out := make(map[string]*metadatapb.TableMetadata, len(schema.GetTables()))
	for _, table := range schema.GetTables() {
		out[table.GetName()] = table
	}
	return out
}

func viewByName(schema *metadatapb.SchemaMetadata) map[string]*metadatapb.ViewMetadata {
	out := make(map[string]*metadatapb.ViewMetadata, len(schema.GetViews()))
	for _, view := range schema.GetViews() {
		out[view.GetName()] = view
	}
	return out
}

func materializedViewByName(schema *metadatapb.SchemaMetadata) map[string]*metadatapb.MaterializedViewMetadata {
	out := make(map[string]*metadatapb.MaterializedViewMetadata, len(schema.GetMaterializedViews()))
	for _, view := range schema.GetMaterializedViews() {
		out[view.GetName()] = view
	}
	return out
}

func procedureByName(schema *metadatapb.SchemaMetadata) map[string]*metadatapb.ProcedureMetadata {
	out := make(map[string]*metadatapb.ProcedureMetadata, len(schema.GetProcedures()))
	for _, proc := range schema.GetProcedures() {
		out[proc.GetName()] = proc
	}
	return out
}

func sequenceByName(schema *metadatapb.SchemaMetadata) map[string]*metadatapb.SequenceMetadata {
	out := make(map[string]*metadatapb.SequenceMetadata, len(schema.GetSequences()))
	for _, seq := range schema.GetSequences() {
		out[seq.GetName()] = seq
	}
	return out
}

func functionKey(fn *metadatapb.FunctionMetadata) string {
	if fn.GetSignature() != "" {
		return fn.GetSignature()
	}
	return fn.GetName()
}

func functionByName(schema *metadatapb.SchemaMetadata) map[string]*metadatapb.FunctionMetadata {
	out := make(map[string]*metadatapb.FunctionMetadata, len(schema.GetFunctions()))
	for _, fn := range schema.GetFunctions() {
		if !fn.GetSkipDump() {
			out[functionKey(fn)] = fn
		}
	}
	return out
}

func addNewSchemaObjects(diff *MetadataDiff, schemaName string, schema *metadatapb.SchemaMetadata) {
	for _, table := range schema.GetTables() {
		if !table.GetSkipDump() {
			diff.TableChanges = append(diff.TableChanges, &TableDiff{Action: DiffActionCreate, SchemaName: schemaName, TableName: table.GetName(), NewTable: table})
		}
	}
	for _, view := range schema.GetViews() {
		if !view.GetSkipDump() {
			diff.ViewChanges = append(diff.ViewChanges, &ViewDiff{Action: DiffActionCreate, SchemaName: schemaName, ViewName: view.GetName(), NewView: view})
		}
	}
	for _, view := range schema.GetMaterializedViews() {
		if !view.GetSkipDump() {
			diff.MaterializedViewChanges = append(diff.MaterializedViewChanges, &MaterializedViewDiff{Action: DiffActionCreate, SchemaName: schemaName, MaterializedViewName: view.GetName(), NewMaterializedView: view})
		}
	}
	for _, fn := range schema.GetFunctions() {
		if !fn.GetSkipDump() {
			diff.FunctionChanges = append(diff.FunctionChanges, &FunctionDiff{Action: DiffActionCreate, SchemaName: schemaName, FunctionName: functionKey(fn), NewFunction: fn})
		}
	}
	for _, proc := range schema.GetProcedures() {
		if !proc.GetSkipDump() {
			diff.ProcedureChanges = append(diff.ProcedureChanges, &ProcedureDiff{Action: DiffActionCreate, SchemaName: schemaName, ProcedureName: proc.GetName(), NewProcedure: proc})
		}
	}
	for _, seq := range schema.GetSequences() {
		if !seq.GetSkipDump() {
			diff.SequenceChanges = append(diff.SequenceChanges, &SequenceDiff{Action: DiffActionCreate, SchemaName: schemaName, SequenceName: seq.GetName(), NewSequence: seq})
		}
	}
}

func compareSchemaObjects(diff *MetadataDiff, schemaName string, oldSchema, newSchema *metadatapb.SchemaMetadata) {
	oldTables, newTables := tableByName(oldSchema), tableByName(newSchema)
	for name, oldTable := range oldTables {
		if _, ok := newTables[name]; !ok && !oldTable.GetSkipDump() {
			diff.TableChanges = append(diff.TableChanges, &TableDiff{Action: DiffActionDrop, SchemaName: schemaName, TableName: name, OldTable: oldTable})
		}
	}
	for name, newTable := range newTables {
		if newTable.GetSkipDump() {
			continue
		}
		oldTable, ok := oldTables[name]
		if !ok {
			diff.TableChanges = append(diff.TableChanges, &TableDiff{Action: DiffActionCreate, SchemaName: schemaName, TableName: name, NewTable: newTable})
			continue
		}
		if oldTable.GetSkipDump() {
			continue
		}
		if tableDiff := compareTable(schemaName, name, oldTable, newTable); tableDiff != nil {
			diff.TableChanges = append(diff.TableChanges, tableDiff)
		}
	}

	oldViews, newViews := viewByName(oldSchema), viewByName(newSchema)
	for name, oldView := range oldViews {
		if _, ok := newViews[name]; !ok && !oldView.GetSkipDump() {
			diff.ViewChanges = append(diff.ViewChanges, &ViewDiff{Action: DiffActionDrop, SchemaName: schemaName, ViewName: name, OldView: oldView})
		}
	}
	for name, newView := range newViews {
		if newView.GetSkipDump() {
			continue
		}
		oldView, ok := oldViews[name]
		if !ok {
			diff.ViewChanges = append(diff.ViewChanges, &ViewDiff{Action: DiffActionCreate, SchemaName: schemaName, ViewName: name, NewView: newView})
		} else if !oldView.GetSkipDump() && viewChanged(oldView, newView) {
			diff.ViewChanges = append(diff.ViewChanges, &ViewDiff{Action: DiffActionAlter, SchemaName: schemaName, ViewName: name, OldView: oldView, NewView: newView})
		}
	}

	oldMVs, newMVs := materializedViewByName(oldSchema), materializedViewByName(newSchema)
	for name, oldView := range oldMVs {
		if _, ok := newMVs[name]; !ok && !oldView.GetSkipDump() {
			diff.MaterializedViewChanges = append(diff.MaterializedViewChanges, &MaterializedViewDiff{Action: DiffActionDrop, SchemaName: schemaName, MaterializedViewName: name, OldMaterializedView: oldView})
		}
	}
	for name, newView := range newMVs {
		if newView.GetSkipDump() {
			continue
		}
		oldView, ok := oldMVs[name]
		if !ok {
			diff.MaterializedViewChanges = append(diff.MaterializedViewChanges, &MaterializedViewDiff{Action: DiffActionCreate, SchemaName: schemaName, MaterializedViewName: name, NewMaterializedView: newView})
		} else if !oldView.GetSkipDump() && materializedViewChanged(oldView, newView) {
			diff.MaterializedViewChanges = append(diff.MaterializedViewChanges, &MaterializedViewDiff{Action: DiffActionAlter, SchemaName: schemaName, MaterializedViewName: name, OldMaterializedView: oldView, NewMaterializedView: newView})
		}
	}

	oldFuncs, newFuncs := functionByName(oldSchema), functionByName(newSchema)
	for key, oldFn := range oldFuncs {
		if _, ok := newFuncs[key]; !ok {
			diff.FunctionChanges = append(diff.FunctionChanges, &FunctionDiff{Action: DiffActionDrop, SchemaName: schemaName, FunctionName: key, OldFunction: oldFn})
		}
	}
	for key, newFn := range newFuncs {
		oldFn, ok := oldFuncs[key]
		if !ok {
			diff.FunctionChanges = append(diff.FunctionChanges, &FunctionDiff{Action: DiffActionCreate, SchemaName: schemaName, FunctionName: key, NewFunction: newFn})
		} else if !functionEqual(oldFn, newFn) {
			diff.FunctionChanges = append(diff.FunctionChanges, &FunctionDiff{Action: DiffActionAlter, SchemaName: schemaName, FunctionName: key, OldFunction: oldFn, NewFunction: newFn})
		}
	}

	oldProcs, newProcs := procedureByName(oldSchema), procedureByName(newSchema)
	for name, oldProc := range oldProcs {
		if _, ok := newProcs[name]; !ok && !oldProc.GetSkipDump() {
			diff.ProcedureChanges = append(diff.ProcedureChanges, &ProcedureDiff{Action: DiffActionDrop, SchemaName: schemaName, ProcedureName: name, OldProcedure: oldProc})
		}
	}
	for name, newProc := range newProcs {
		if newProc.GetSkipDump() {
			continue
		}
		oldProc, ok := oldProcs[name]
		if !ok {
			diff.ProcedureChanges = append(diff.ProcedureChanges, &ProcedureDiff{Action: DiffActionCreate, SchemaName: schemaName, ProcedureName: name, NewProcedure: newProc})
		} else if !oldProc.GetSkipDump() && oldProc.GetDefinition() != newProc.GetDefinition() {
			diff.ProcedureChanges = append(diff.ProcedureChanges, &ProcedureDiff{Action: DiffActionAlter, SchemaName: schemaName, ProcedureName: name, OldProcedure: oldProc, NewProcedure: newProc})
		}
	}

	oldSeqs, newSeqs := sequenceByName(oldSchema), sequenceByName(newSchema)
	for name, oldSeq := range oldSeqs {
		if _, ok := newSeqs[name]; !ok && !oldSeq.GetSkipDump() {
			diff.SequenceChanges = append(diff.SequenceChanges, &SequenceDiff{Action: DiffActionDrop, SchemaName: schemaName, SequenceName: name, OldSequence: oldSeq})
		}
	}
	for name, newSeq := range newSeqs {
		if newSeq.GetSkipDump() {
			continue
		}
		oldSeq, ok := oldSeqs[name]
		if !ok {
			diff.SequenceChanges = append(diff.SequenceChanges, &SequenceDiff{Action: DiffActionCreate, SchemaName: schemaName, SequenceName: name, NewSequence: newSeq})
		} else if !oldSeq.GetSkipDump() && !sequenceEqual(oldSeq, newSeq) {
			diff.SequenceChanges = append(diff.SequenceChanges, &SequenceDiff{Action: DiffActionAlter, SchemaName: schemaName, SequenceName: name, OldSequence: oldSeq, NewSequence: newSeq})
		}
	}
}

func compareTable(schemaName, tableName string, oldTable, newTable *metadatapb.TableMetadata) *TableDiff {
	tableDiff := &TableDiff{Action: DiffActionAlter, SchemaName: schemaName, TableName: tableName, OldTable: oldTable, NewTable: newTable}
	changed := false

	oldColumns := make(map[string]*metadatapb.ColumnMetadata, len(oldTable.GetColumns()))
	for _, column := range oldTable.GetColumns() {
		oldColumns[column.GetName()] = column
	}
	newColumns := make(map[string]*metadatapb.ColumnMetadata, len(newTable.GetColumns()))
	for _, column := range newTable.GetColumns() {
		newColumns[column.GetName()] = column
	}
	for name, oldColumn := range oldColumns {
		if _, ok := newColumns[name]; !ok {
			tableDiff.ColumnChanges = append(tableDiff.ColumnChanges, &ColumnDiff{Action: DiffActionDrop, OldColumn: oldColumn})
			changed = true
		}
	}
	for name, newColumn := range newColumns {
		oldColumn, ok := oldColumns[name]
		if !ok {
			tableDiff.ColumnChanges = append(tableDiff.ColumnChanges, &ColumnDiff{Action: DiffActionCreate, NewColumn: newColumn})
			changed = true
		} else if !columnEqual(oldColumn, newColumn) {
			tableDiff.ColumnChanges = append(tableDiff.ColumnChanges, &ColumnDiff{Action: DiffActionAlter, OldColumn: oldColumn, NewColumn: newColumn})
			changed = true
		}
	}

	tableDiff.IndexChanges = compareIndexes(oldTable.GetIndexes(), newTable.GetIndexes())
	tableDiff.ForeignKeyChanges = compareForeignKeys(oldTable.GetForeignKeys(), newTable.GetForeignKeys())
	tableDiff.CheckConstraintChanges = compareCheckConstraints(oldTable.GetCheckConstraints(), newTable.GetCheckConstraints())
	tableDiff.TriggerChanges = compareTriggers(oldTable.GetTriggers(), newTable.GetTriggers())
	if len(tableDiff.IndexChanges) > 0 || len(tableDiff.ForeignKeyChanges) > 0 || len(tableDiff.CheckConstraintChanges) > 0 || len(tableDiff.TriggerChanges) > 0 || oldTable.GetComment() != newTable.GetComment() {
		changed = true
	}

	if !changed {
		return nil
	}
	return tableDiff
}

func columnEqual(a, b *metadatapb.ColumnMetadata) bool {
	return a.GetType() == b.GetType() &&
		a.GetNullable() == b.GetNullable() &&
		a.GetDefault() == b.GetDefault() &&
		a.GetDefaultOnNull() == b.GetDefaultOnNull() &&
		a.GetOnUpdate() == b.GetOnUpdate() &&
		a.GetComment() == b.GetComment() &&
		a.GetCharacterSet() == b.GetCharacterSet() &&
		a.GetCollation() == b.GetCollation() &&
		a.GetIsInvisible() == b.GetIsInvisible() &&
		a.GetIsIdentity() == b.GetIsIdentity() &&
		a.GetIdentityGeneration() == b.GetIdentityGeneration() &&
		a.GetIdentitySeed() == b.GetIdentitySeed() &&
		a.GetIdentityIncrement() == b.GetIdentityIncrement() &&
		generationEqual(a.GetGeneration(), b.GetGeneration())
}

func generationEqual(a, b *metadatapb.GenerationMetadata) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.GetType() == b.GetType() && normalizeExpression(a.GetExpression()) == normalizeExpression(b.GetExpression())
}

func compareIndexes(oldIndexes, newIndexes []*metadatapb.IndexMetadata) []*IndexDiff {
	oldMap := make(map[string]*metadatapb.IndexMetadata, len(oldIndexes))
	for _, index := range oldIndexes {
		oldMap[index.GetName()] = index
	}
	newMap := make(map[string]*metadatapb.IndexMetadata, len(newIndexes))
	for _, index := range newIndexes {
		newMap[index.GetName()] = index
	}
	var changes []*IndexDiff
	for name, oldIndex := range oldMap {
		if _, ok := newMap[name]; !ok {
			changes = append(changes, &IndexDiff{Action: DiffActionDrop, OldIndex: oldIndex})
		}
	}
	for name, newIndex := range newMap {
		oldIndex, ok := oldMap[name]
		if !ok {
			changes = append(changes, &IndexDiff{Action: DiffActionCreate, NewIndex: newIndex})
		} else if !indexEqual(oldIndex, newIndex) {
			changes = append(changes, &IndexDiff{Action: DiffActionDrop, OldIndex: oldIndex})
			changes = append(changes, &IndexDiff{Action: DiffActionCreate, NewIndex: newIndex})
		}
	}
	return changes
}

func indexEqual(a, b *metadatapb.IndexMetadata) bool {
	if a.GetIsConstraint() != b.GetIsConstraint() {
		return false
	}
	if a.GetIsConstraint() && b.GetIsConstraint() {
		if a.GetUnique() != b.GetUnique() || a.GetPrimary() != b.GetPrimary() {
			return false
		}
		if len(a.GetExpressions()) != len(b.GetExpressions()) {
			return false
		}
		for i, expression := range a.GetExpressions() {
			if strings.Trim(expression, `"`) != strings.Trim(b.GetExpressions()[i], `"`) {
				return false
			}
		}
		return descendingEqual(a.GetDescending(), b.GetDescending())
	}
	if a.GetType() != b.GetType() || a.GetUnique() != b.GetUnique() || a.GetPrimary() != b.GetPrimary() || a.GetVisible() != b.GetVisible() {
		return false
	}
	if len(a.GetExpressions()) != len(b.GetExpressions()) {
		return false
	}
	for i, expression := range a.GetExpressions() {
		if !keyExpressionsEqual(expression, b.GetExpressions()[i]) {
			return false
		}
	}
	return descendingEqual(a.GetDescending(), b.GetDescending())
}

func descendingEqual(a, b []bool) bool {
	allFalse := func(values []bool) bool {
		for _, value := range values {
			if value {
				return false
			}
		}
		return true
	}
	if allFalse(a) && allFalse(b) {
		return true
	}
	return slices.Equal(a, b)
}

func compareForeignKeys(oldFKs, newFKs []*metadatapb.ForeignKeyMetadata) []*ForeignKeyDiff {
	oldMap := make(map[string]*metadatapb.ForeignKeyMetadata, len(oldFKs))
	for _, fk := range oldFKs {
		oldMap[fk.GetName()] = fk
	}
	newMap := make(map[string]*metadatapb.ForeignKeyMetadata, len(newFKs))
	for _, fk := range newFKs {
		newMap[fk.GetName()] = fk
	}
	var changes []*ForeignKeyDiff
	for name, oldFK := range oldMap {
		if _, ok := newMap[name]; !ok {
			changes = append(changes, &ForeignKeyDiff{Action: DiffActionDrop, OldForeignKey: oldFK})
		}
	}
	for name, newFK := range newMap {
		oldFK, ok := oldMap[name]
		if !ok {
			changes = append(changes, &ForeignKeyDiff{Action: DiffActionCreate, NewForeignKey: newFK})
		} else if !foreignKeyEqual(oldFK, newFK) {
			changes = append(changes, &ForeignKeyDiff{Action: DiffActionDrop, OldForeignKey: oldFK})
			changes = append(changes, &ForeignKeyDiff{Action: DiffActionCreate, NewForeignKey: newFK})
		}
	}
	return changes
}

// keyExpressionsEqual compares two index expressions. A bare identifier,
// quoted or not, is the same column only under its exact name: the sync quotes
// a case-sensitive name, and a quoted name may hold significant spaces, so it
// is not normalized. Anything else compares as normalized text.
func keyExpressionsEqual(a, b string) bool {
	if isKeyIdentifier(a) && isKeyIdentifier(b) {
		return unquoteIdentifier(a) == unquoteIdentifier(b)
	}
	return normalizeExpression(a) == normalizeExpression(b)
}

func isKeyIdentifier(expression string) bool {
	if strings.HasPrefix(expression, `"`) && strings.HasSuffix(expression, `"`) && len(expression) > 1 {
		return !strings.Contains(expression[1:len(expression)-1], `"`)
	}
	return isSimpleIdentifier(expression)
}

func unquoteIdentifiers(identifiers []string) []string {
	unquoted := make([]string, len(identifiers))
	for i, identifier := range identifiers {
		unquoted[i] = unquoteIdentifier(identifier)
	}
	return unquoted
}

func foreignKeyEqual(a, b *metadatapb.ForeignKeyMetadata) bool {
	return a.GetReferencedSchema() == b.GetReferencedSchema() &&
		a.GetReferencedTable() == b.GetReferencedTable() &&
		a.GetOnDelete() == b.GetOnDelete() &&
		a.GetOnUpdate() == b.GetOnUpdate() &&
		a.GetMatchType() == b.GetMatchType() &&
		slices.Equal(unquoteIdentifiers(a.GetColumns()), unquoteIdentifiers(b.GetColumns())) &&
		slices.Equal(unquoteIdentifiers(a.GetReferencedColumns()), unquoteIdentifiers(b.GetReferencedColumns()))
}

func compareCheckConstraints(oldChecks, newChecks []*metadatapb.CheckConstraintMetadata) []*CheckConstraintDiff {
	oldMap := make(map[string]*metadatapb.CheckConstraintMetadata, len(oldChecks))
	for _, check := range oldChecks {
		oldMap[check.GetName()] = check
	}
	newMap := make(map[string]*metadatapb.CheckConstraintMetadata, len(newChecks))
	for _, check := range newChecks {
		newMap[check.GetName()] = check
	}
	var changes []*CheckConstraintDiff
	for name, oldCheck := range oldMap {
		if _, ok := newMap[name]; !ok {
			changes = append(changes, &CheckConstraintDiff{Action: DiffActionDrop, OldCheckConstraint: oldCheck})
		}
	}
	for name, newCheck := range newMap {
		oldCheck, ok := oldMap[name]
		if !ok {
			changes = append(changes, &CheckConstraintDiff{Action: DiffActionCreate, NewCheckConstraint: newCheck})
		} else if normalizeExpression(oldCheck.GetExpression()) != normalizeExpression(newCheck.GetExpression()) {
			changes = append(changes, &CheckConstraintDiff{Action: DiffActionDrop, OldCheckConstraint: oldCheck})
			changes = append(changes, &CheckConstraintDiff{Action: DiffActionCreate, NewCheckConstraint: newCheck})
		}
	}
	return changes
}

func compareTriggers(oldTriggers, newTriggers []*metadatapb.TriggerMetadata) []*TriggerDiff {
	oldMap := make(map[string]*metadatapb.TriggerMetadata, len(oldTriggers))
	for _, trigger := range oldTriggers {
		oldMap[trigger.GetName()] = trigger
	}
	newMap := make(map[string]*metadatapb.TriggerMetadata, len(newTriggers))
	for _, trigger := range newTriggers {
		newMap[trigger.GetName()] = trigger
	}
	var changes []*TriggerDiff
	for name, oldTrigger := range oldMap {
		if _, ok := newMap[name]; !ok {
			changes = append(changes, &TriggerDiff{Action: DiffActionDrop, TriggerName: name, OldTrigger: oldTrigger})
		}
	}
	for name, newTrigger := range newMap {
		oldTrigger, ok := oldMap[name]
		if !ok {
			changes = append(changes, &TriggerDiff{Action: DiffActionCreate, TriggerName: name, NewTrigger: newTrigger})
		} else if !triggerEqual(oldTrigger, newTrigger) {
			changes = append(changes, &TriggerDiff{Action: DiffActionDrop, TriggerName: name, OldTrigger: oldTrigger})
			changes = append(changes, &TriggerDiff{Action: DiffActionCreate, TriggerName: name, NewTrigger: newTrigger})
		}
	}
	return changes
}

func triggerEqual(a, b *metadatapb.TriggerMetadata) bool {
	return a.GetName() == b.GetName() && a.GetEvent() == b.GetEvent() && a.GetTiming() == b.GetTiming() &&
		normalizeTriggerBody(a.GetBody()) == normalizeTriggerBody(b.GetBody())
}

var (
	triggerQualifiedName = regexp.MustCompile(`(?i)(TRIGGER\s+|\sON\s+)[A-Z][A-Z0-9_$#]*\.`)
	triggerCreateReplace = regexp.MustCompile(`(?i)CREATE\s+OR\s+REPLACE\s+TRIGGER`)
)

func normalizeTriggerBody(body string) string {
	headerEnd := triggerHeaderEnd(body)
	header := triggerCreateReplace.ReplaceAllString(normalizeExpression(body[:headerEnd]), "CREATE TRIGGER")
	header = replaceFirstMatch(header, triggerQualifiedName)
	header = replaceFirstMatch(header, triggerQualifiedName)
	return header + normalizeExpression(body[headerEnd:])
}

func triggerHeaderEnd(body string) int {
	upper := strings.ToUpper(body)
	var quote byte
	for i := 0; i < len(body); i++ {
		c := body[i]
		if quote != 0 {
			if c == quote {
				quote = 0
			}
			continue
		}
		if c == '\'' {
			quote = c
			continue
		}
		if i > 0 && isIdentifierByte(body[i-1]) {
			continue
		}
		for _, keyword := range []string{"BEGIN", "DECLARE"} {
			if strings.HasPrefix(upper[i:], keyword) && (i+len(keyword) == len(body) || !isIdentifierByte(body[i+len(keyword)])) {
				return i
			}
		}
	}
	return len(body)
}

func replaceFirstMatch(text string, pattern *regexp.Regexp) string {
	match := pattern.FindStringSubmatchIndex(text)
	if match == nil || len(match) < 4 {
		return text
	}
	return text[:match[0]] + text[match[2]:match[3]] + text[match[1]:]
}

func viewChanged(a, b *metadatapb.ViewMetadata) bool {
	return viewNeedsRecreate(a, b) || a.GetComment() != b.GetComment()
}

func viewNeedsRecreate(a, b *metadatapb.ViewMetadata) bool {
	return a.GetDefinition() != b.GetDefinition() || normalizeTriggers(a.GetTriggers()) != normalizeTriggers(b.GetTriggers())
}

func materializedViewChanged(a, b *metadatapb.MaterializedViewMetadata) bool {
	return materializedViewNeedsRecreate(a, b) || a.GetComment() != b.GetComment()
}

func materializedViewNeedsRecreate(a, b *metadatapb.MaterializedViewMetadata) bool {
	return a.GetDefinition() != b.GetDefinition() || normalizeTriggers(a.GetTriggers()) != normalizeTriggers(b.GetTriggers()) || len(compareIndexes(materializedViewIndexes(a.GetIndexes()), materializedViewIndexes(b.GetIndexes()))) > 0
}

func normalizeTriggers(triggers []*metadatapb.TriggerMetadata) string {
	var parts []string
	for _, trigger := range triggers {
		parts = append(parts, trigger.GetName()+"\x00"+normalizeTriggerBody(trigger.GetBody()))
	}
	return strings.Join(parts, "\x01")
}

func functionEqual(a, b *metadatapb.FunctionMetadata) bool {
	return a.GetDefinition() == b.GetDefinition() && a.GetComment() == b.GetComment()
}

func sequenceEqual(a, b *metadatapb.SequenceMetadata) bool {
	return a.GetDataType() == b.GetDataType() &&
		a.GetStart() == b.GetStart() &&
		a.GetMinValue() == b.GetMinValue() &&
		a.GetMaxValue() == b.GetMaxValue() &&
		a.GetIncrement() == b.GetIncrement() &&
		a.GetCycle() == b.GetCycle() &&
		a.GetCacheSize() == b.GetCacheSize() &&
		a.GetOwnerTable() == b.GetOwnerTable() &&
		a.GetOwnerColumn() == b.GetOwnerColumn() &&
		a.GetComment() == b.GetComment()
}

// normalizeExpression makes whitespace-only differences outside quoted text
// irrelevant. Whitespace inside a string literal or a quoted identifier is part
// of the value and stays.
func normalizeExpression(expression string) string {
	if strings.Contains(expression, "--") {
		return expression
	}
	var out strings.Builder
	var quote byte
	pendingSpace := false
	for i := 0; i < len(expression); i++ {
		if end, ok := alternativeQuoteEnd(expression, i); ok {
			out.WriteString(expression[i:end])
			i = end - 1
			continue
		}
		c := expression[i]
		if quote != 0 {
			out.WriteByte(c)
			if c == quote {
				if i+1 < len(expression) && expression[i+1] == quote {
					out.WriteByte(expression[i+1])
					i++
				} else {
					quote = 0
				}
			}
			continue
		}
		if c == '\'' || c == '"' {
			if pendingSpace && out.Len() > 0 {
				out.WriteByte(' ')
			}
			pendingSpace = false
			quote = c
			out.WriteByte(c)
			continue
		}
		switch c {
		case ' ', '\t', '\n', '\r', '\f':
			pendingSpace = true
		default:
			if pendingSpace && out.Len() > 0 {
				out.WriteByte(' ')
			}
			pendingSpace = false
			out.WriteByte(c)
		}
	}
	return out.String()
}

func sortDiff(diff *MetadataDiff) {
	compareNameAction := func(schemaA, nameA string, actionA DiffAction, schemaB, nameB string, actionB DiffAction) int {
		if schemaA != schemaB {
			return strings.Compare(schemaA, schemaB)
		}
		if nameA != nameB {
			return strings.Compare(nameA, nameB)
		}
		return strings.Compare(string(actionA), string(actionB))
	}
	slices.SortFunc(diff.SchemaChanges, func(a, b *SchemaDiffEntry) int {
		return compareNameAction(a.SchemaName, "", a.Action, b.SchemaName, "", b.Action)
	})
	slices.SortFunc(diff.TableChanges, func(a, b *TableDiff) int {
		return compareNameAction(a.SchemaName, a.TableName, a.Action, b.SchemaName, b.TableName, b.Action)
	})
	for _, table := range diff.TableChanges {
		slices.SortFunc(table.ColumnChanges, func(a, b *ColumnDiff) int {
			aPos, aName := columnPosition(a)
			bPos, bName := columnPosition(b)
			if aPos > 0 && bPos > 0 && aPos != bPos {
				return int(aPos - bPos)
			}
			if aName != bName {
				return strings.Compare(aName, bName)
			}
			return strings.Compare(string(a.Action), string(b.Action))
		})
		slices.SortFunc(table.IndexChanges, func(a, b *IndexDiff) int {
			return compareNameAction("", diffIndexName(a), a.Action, "", diffIndexName(b), b.Action)
		})
		slices.SortFunc(table.ForeignKeyChanges, func(a, b *ForeignKeyDiff) int {
			return compareNameAction("", diffForeignKeyName(a), a.Action, "", diffForeignKeyName(b), b.Action)
		})
		slices.SortFunc(table.CheckConstraintChanges, func(a, b *CheckConstraintDiff) int {
			return compareNameAction("", diffCheckName(a), a.Action, "", diffCheckName(b), b.Action)
		})
		slices.SortFunc(table.TriggerChanges, func(a, b *TriggerDiff) int {
			return compareNameAction("", a.TriggerName, a.Action, "", b.TriggerName, b.Action)
		})
	}
	slices.SortFunc(diff.ViewChanges, func(a, b *ViewDiff) int {
		return compareNameAction(a.SchemaName, a.ViewName, a.Action, b.SchemaName, b.ViewName, b.Action)
	})
	slices.SortFunc(diff.MaterializedViewChanges, func(a, b *MaterializedViewDiff) int {
		return compareNameAction(a.SchemaName, a.MaterializedViewName, a.Action, b.SchemaName, b.MaterializedViewName, b.Action)
	})
	slices.SortFunc(diff.FunctionChanges, func(a, b *FunctionDiff) int {
		return compareNameAction(a.SchemaName, a.FunctionName, a.Action, b.SchemaName, b.FunctionName, b.Action)
	})
	slices.SortFunc(diff.ProcedureChanges, func(a, b *ProcedureDiff) int {
		return compareNameAction(a.SchemaName, a.ProcedureName, a.Action, b.SchemaName, b.ProcedureName, b.Action)
	})
	slices.SortFunc(diff.SequenceChanges, func(a, b *SequenceDiff) int {
		return compareNameAction(a.SchemaName, a.SequenceName, a.Action, b.SchemaName, b.SequenceName, b.Action)
	})
}

func columnPosition(diff *ColumnDiff) (int32, string) {
	if diff.NewColumn != nil {
		return diff.NewColumn.GetPosition(), diff.NewColumn.GetName()
	}
	if diff.OldColumn != nil {
		return diff.OldColumn.GetPosition(), diff.OldColumn.GetName()
	}
	return 0, ""
}

func diffColumnName(diff *ColumnDiff) string {
	if diff.NewColumn != nil {
		return diff.NewColumn.GetName()
	}
	return diff.OldColumn.GetName()
}

func diffIndexName(diff *IndexDiff) string {
	if diff.NewIndex != nil {
		return diff.NewIndex.GetName()
	}
	return diff.OldIndex.GetName()
}

func diffForeignKeyName(diff *ForeignKeyDiff) string {
	if diff.NewForeignKey != nil {
		return diff.NewForeignKey.GetName()
	}
	return diff.OldForeignKey.GetName()
}

func diffCheckName(diff *CheckConstraintDiff) string {
	if diff.NewCheckConstraint != nil {
		return diff.NewCheckConstraint.GetName()
	}
	return diff.OldCheckConstraint.GetName()
}

// alternativeQuoteEnd returns the end of Oracle's q'...' or nq'...' literal at
// start, where the delimiter is the character after the quote and a bracket
// delimiter closes with its pair.
func alternativeQuoteEnd(expression string, start int) (int, bool) {
	i := start
	if i+1 < len(expression) && (expression[i] == 'n' || expression[i] == 'N') {
		i++
	}
	if i+1 >= len(expression) || (expression[i] != 'q' && expression[i] != 'Q') || expression[i+1] != '\'' {
		return 0, false
	}
	if i+2 >= len(expression) {
		return 0, false
	}
	open := expression[i+2]
	close := open
	switch open {
	case '[':
		close = ']'
	case '{':
		close = '}'
	case '(':
		close = ')'
	case '<':
		close = '>'
	}
	for j := i + 3; j+1 < len(expression); j++ {
		if expression[j] == close && expression[j+1] == '\'' {
			return j + 2, true
		}
	}
	return len(expression), true
}
