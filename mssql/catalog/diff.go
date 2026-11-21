package catalog

import (
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// DiffAction says what happened to an object between two snapshots.
type DiffAction string

const (
	DiffCreate DiffAction = "CREATE"
	DiffDrop   DiffAction = "DROP"
	DiffAlter  DiffAction = "ALTER"
)

// SchemaDiffEntry describes a schema that was added, dropped, or altered.
type SchemaDiffEntry struct {
	Action     DiffAction
	SchemaName string
	OldSchema  *metadata.SchemaMetadata
	NewSchema  *metadata.SchemaMetadata
}

// TableDiffEntry describes a table change. A table's columns, indexes, foreign
// keys, checks and triggers are compared only for an altered table.
type TableDiffEntry struct {
	Action     DiffAction
	SchemaName string
	TableName  string
	OldTable   *metadata.TableMetadata
	NewTable   *metadata.TableMetadata

	Columns     []*ColumnDiffEntry
	Indexes     []*IndexDiffEntry
	ForeignKeys []*ForeignKeyDiffEntry
	Checks      []*CheckConstraintDiffEntry
	Triggers    []*TriggerDiffEntry
}

// ColumnDiffEntry describes a column change.
type ColumnDiffEntry struct {
	Action    DiffAction
	OldColumn *metadata.ColumnMetadata
	NewColumn *metadata.ColumnMetadata
}

// IndexDiffEntry describes an index or key constraint change.
type IndexDiffEntry struct {
	Action   DiffAction
	OldIndex *metadata.IndexMetadata
	NewIndex *metadata.IndexMetadata
}

// ForeignKeyDiffEntry describes a foreign key change.
type ForeignKeyDiffEntry struct {
	Action        DiffAction
	OldForeignKey *metadata.ForeignKeyMetadata
	NewForeignKey *metadata.ForeignKeyMetadata
}

// CheckConstraintDiffEntry describes a check constraint change.
type CheckConstraintDiffEntry struct {
	Action             DiffAction
	OldCheckConstraint *metadata.CheckConstraintMetadata
	NewCheckConstraint *metadata.CheckConstraintMetadata
}

// TriggerDiffEntry describes a trigger change.
type TriggerDiffEntry struct {
	Action     DiffAction
	OldTrigger *metadata.TriggerMetadata
	NewTrigger *metadata.TriggerMetadata
}

// SynonymDiffEntry describes a synonym change.
type SynonymDiffEntry struct {
	Action      DiffAction
	SchemaName  string
	SynonymName string
	OldSynonym  *metadata.SynonymMetadata
	NewSynonym  *metadata.SynonymMetadata
}

// ViewDiffEntry describes a view change.
type ViewDiffEntry struct {
	Action     DiffAction
	SchemaName string
	ViewName   string
	OldView    *metadata.ViewMetadata
	NewView    *metadata.ViewMetadata
}

// FunctionDiffEntry describes a function change.
type FunctionDiffEntry struct {
	Action       DiffAction
	SchemaName   string
	FunctionName string
	OldFunction  *metadata.FunctionMetadata
	NewFunction  *metadata.FunctionMetadata
}

// ProcedureDiffEntry describes a procedure change.
type ProcedureDiffEntry struct {
	Action        DiffAction
	SchemaName    string
	ProcedureName string
	OldProcedure  *metadata.ProcedureMetadata
	NewProcedure  *metadata.ProcedureMetadata
}

// SequenceDiffEntry describes a sequence change.
type SequenceDiffEntry struct {
	Action       DiffAction
	SchemaName   string
	SequenceName string
	OldSequence  *metadata.SequenceMetadata
	NewSequence  *metadata.SequenceMetadata
}

// SchemaDiff holds every difference between two catalog states.
type SchemaDiff struct {
	DatabaseName string
	Schemas      []*SchemaDiffEntry
	Tables       []*TableDiffEntry
	Views        []*ViewDiffEntry
	Synonyms     []*SynonymDiffEntry
	Functions    []*FunctionDiffEntry
	Procedures   []*ProcedureDiffEntry
	Sequences    []*SequenceDiffEntry
}

// IsEmpty reports whether the two snapshots differ in nothing the migration
// generator models.
func (d *SchemaDiff) IsEmpty() bool {
	return len(d.Schemas) == 0 &&
		len(d.Tables) == 0 &&
		len(d.Views) == 0 &&
		len(d.Synonyms) == 0 &&
		len(d.Functions) == 0 &&
		len(d.Procedures) == 0 &&
		len(d.Sequences) == 0
}

// Diff compares from and to. from is the current state and to is the desired one.
func Diff(from, to *Catalog) *SchemaDiff {
	diff := &SchemaDiff{DatabaseName: to.DatabaseName()}

	for _, name := range sortedSchemaNames(from.schemas) {
		oldSchema := from.schemaByName[name]
		if peerSchema(to, from, name) == nil {
			diff.Schemas = append(diff.Schemas, &SchemaDiffEntry{
				Action:     DiffDrop,
				SchemaName: name,
				OldSchema:  oldSchema,
			})
			if oldSchema != nil {
				addDroppedSchemaObjects(diff, name, oldSchema)
			}
		}
	}

	for _, name := range sortedSchemaNames(to.schemas) {
		newSchema := to.schemaByName[name]
		if newSchema == nil {
			continue
		}
		oldSchema := peerSchema(from, to, name)
		if oldSchema == nil {
			diff.Schemas = append(diff.Schemas, &SchemaDiffEntry{
				Action:     DiffCreate,
				SchemaName: name,
				NewSchema:  newSchema,
			})
			addNewSchemaObjects(diff, name, newSchema)
			continue
		}
		compareSchemaObjects(diff, name, oldSchema, newSchema)
	}

	sortDiff(diff)
	return diff
}

func peerSchema(target, peer *Catalog, name string) *metadata.SchemaMetadata {
	if schema := target.schemaByName[name]; schema != nil {
		return schema
	}
	targetFold := target.schemasByFold[normalize(name)]
	peerFold := peer.schemasByFold[normalize(name)]
	if len(targetFold) == 1 && len(peerFold) == 1 {
		return targetFold[0]
	}
	return nil
}

// addDroppedSchemaObjects turns a removed schema into object-level drops, so
// the dependency-ordered drop phases empty it before DROP SCHEMA runs and
// every foreign key sees its own drop.
func addDroppedSchemaObjects(diff *SchemaDiff, schemaName string, schema *metadata.SchemaMetadata) {
	for _, name := range sortedTableNames(schema) {
		table := tableByName(schema, name)
		if table != nil && !table.GetSkipDump() {
			diff.Tables = append(diff.Tables, &TableDiffEntry{
				Action:     DiffDrop,
				SchemaName: schemaName,
				TableName:  name,
				OldTable:   table,
			})
		}
	}
	for _, name := range sortedViewNames(schema) {
		view := viewByName(schema, name)
		if view != nil && !view.GetSkipDump() {
			diff.Views = append(diff.Views, &ViewDiffEntry{
				Action:     DiffDrop,
				SchemaName: schemaName,
				ViewName:   name,
				OldView:    view,
			})
		}
	}
	for _, function := range schema.GetFunctions() {
		if function.GetSkipDump() {
			continue
		}
		name := function.Signature
		if name == "" {
			name = function.Name
		}
		diff.Functions = append(diff.Functions, &FunctionDiffEntry{
			Action:       DiffDrop,
			SchemaName:   schemaName,
			FunctionName: name,
			OldFunction:  function,
		})
	}
	for _, name := range sortedProcedureNames(schema) {
		procedure := procedureByName(schema, name)
		if procedure != nil && !procedure.SkipDump {
			diff.Procedures = append(diff.Procedures, &ProcedureDiffEntry{
				Action:        DiffDrop,
				SchemaName:    schemaName,
				ProcedureName: name,
				OldProcedure:  procedure,
			})
		}
	}
	for _, synonym := range schema.Synonyms {
		diff.Synonyms = append(diff.Synonyms, &SynonymDiffEntry{
			Action:      DiffDrop,
			SchemaName:  schemaName,
			SynonymName: synonym.Name,
			OldSynonym:  synonym,
		})
	}
	for _, sequence := range schema.GetSequences() {
		if sequence.GetSkipDump() {
			continue
		}
		diff.Sequences = append(diff.Sequences, &SequenceDiffEntry{
			Action:       DiffDrop,
			SchemaName:   schemaName,
			SequenceName: sequence.Name,
			OldSequence:  sequence,
		})
	}
}

func addNewSchemaObjects(diff *SchemaDiff, schemaName string, schema *metadata.SchemaMetadata) {
	for _, name := range sortedTableNames(schema) {
		table := tableByName(schema, name)
		if table != nil && !table.GetSkipDump() {
			diff.Tables = append(diff.Tables, &TableDiffEntry{
				Action:     DiffCreate,
				SchemaName: schemaName,
				TableName:  name,
				NewTable:   table,
			})
		}
	}
	for _, name := range sortedViewNames(schema) {
		view := viewByName(schema, name)
		if view != nil && !view.GetSkipDump() {
			diff.Views = append(diff.Views, &ViewDiffEntry{
				Action:     DiffCreate,
				SchemaName: schemaName,
				ViewName:   name,
				NewView:    view,
			})
		}
	}
	for _, synonym := range schema.Synonyms {
		diff.Synonyms = append(diff.Synonyms, &SynonymDiffEntry{Action: DiffCreate, SchemaName: schemaName, SynonymName: synonym.Name, NewSynonym: synonym})
	}
	for _, function := range schema.GetFunctions() {
		if function.GetSkipDump() {
			continue
		}
		name := function.Signature
		if name == "" {
			name = function.Name
		}
		diff.Functions = append(diff.Functions, &FunctionDiffEntry{
			Action:       DiffCreate,
			SchemaName:   schemaName,
			FunctionName: name,
			NewFunction:  function,
		})
	}
	for _, name := range sortedProcedureNames(schema) {
		procedure := procedureByName(schema, name)
		if procedure != nil && !procedure.SkipDump {
			diff.Procedures = append(diff.Procedures, &ProcedureDiffEntry{
				Action:        DiffCreate,
				SchemaName:    schemaName,
				ProcedureName: name,
				NewProcedure:  procedure,
			})
		}
	}
	for _, sequence := range schema.GetSequences() {
		if sequence.GetSkipDump() {
			continue
		}
		diff.Sequences = append(diff.Sequences, &SequenceDiffEntry{
			Action:       DiffCreate,
			SchemaName:   schemaName,
			SequenceName: sequence.Name,
			NewSequence:  sequence,
		})
	}
}

func compareSchemaObjects(diff *SchemaDiff, schemaName string, oldSchema, newSchema *metadata.SchemaMetadata) {
	for _, name := range sortedTableNames(oldSchema) {
		oldTable := tableByName(oldSchema, name)
		if newTable := tableByName(newSchema, name); newTable == nil {
			if oldTable != nil && !oldTable.GetSkipDump() {
				diff.Tables = append(diff.Tables, &TableDiffEntry{
					Action:     DiffDrop,
					SchemaName: schemaName,
					TableName:  name,
					OldTable:   oldTable,
				})
			}
		}
	}
	for _, name := range sortedTableNames(newSchema) {
		newTable := tableByName(newSchema, name)
		if newTable == nil || newTable.GetSkipDump() {
			continue
		}
		oldTable := tableByName(oldSchema, name)
		if oldTable == nil {
			diff.Tables = append(diff.Tables, &TableDiffEntry{
				Action:     DiffCreate,
				SchemaName: schemaName,
				TableName:  name,
				NewTable:   newTable,
			})
			continue
		}
		if oldTable.GetSkipDump() {
			continue
		}
		if tableDiff := compareTableDetails(schemaName, name, oldTable, newTable); tableDiff != nil {
			diff.Tables = append(diff.Tables, tableDiff)
		}
	}

	compareViews(diff, schemaName, oldSchema, newSchema)
	compareSynonyms(diff, schemaName, oldSchema, newSchema)
	compareFunctions(diff, schemaName, oldSchema, newSchema)
	compareProcedures(diff, schemaName, oldSchema, newSchema)
	compareSequences(diff, schemaName, oldSchema, newSchema)
}

func compareTableDetails(schemaName, tableName string, oldTable, newTable *metadata.TableMetadata) *TableDiffEntry {
	tableDiff := &TableDiffEntry{
		Action:     DiffAlter,
		SchemaName: schemaName,
		TableName:  tableName,
		OldTable:   oldTable,
		NewTable:   newTable,
	}
	hasChanges := false
	if changes := compareColumns(oldTable, newTable); len(changes) > 0 {
		tableDiff.Columns = changes
		hasChanges = true
	}
	if changes := compareIndexes(oldTable, newTable); len(changes) > 0 {
		tableDiff.Indexes = changes
		hasChanges = true
	}
	if changes := compareForeignKeys(oldTable.ForeignKeys, newTable.ForeignKeys); len(changes) > 0 {
		tableDiff.ForeignKeys = changes
		hasChanges = true
	}
	if changes := compareCheckConstraints(oldTable.CheckConstraints, newTable.CheckConstraints); len(changes) > 0 {
		tableDiff.Checks = changes
		hasChanges = true
	}
	if changes := compareTriggers(oldTable.Triggers, newTable.Triggers); len(changes) > 0 {
		tableDiff.Triggers = changes
		hasChanges = true
	}
	if oldTable.Comment != newTable.Comment {
		hasChanges = true
	}
	if !hasChanges {
		return nil
	}
	return tableDiff
}

func compareColumns(oldTable, newTable *metadata.TableMetadata) []*ColumnDiffEntry {
	var changes []*ColumnDiffEntry
	for _, oldColumn := range oldTable.GetColumns() {
		if columnByName(newTable, oldColumn.Name) == nil {
			changes = append(changes, &ColumnDiffEntry{Action: DiffDrop, OldColumn: oldColumn})
		}
	}
	for _, newColumn := range newTable.GetColumns() {
		oldColumn := columnByName(oldTable, newColumn.Name)
		if oldColumn == nil {
			changes = append(changes, &ColumnDiffEntry{Action: DiffCreate, NewColumn: newColumn})
			continue
		}
		if !columnsEqual(oldColumn, newColumn) {
			changes = append(changes, &ColumnDiffEntry{Action: DiffAlter, OldColumn: oldColumn, NewColumn: newColumn})
		}
	}
	return changes
}

func columnsEqual(a, b *metadata.ColumnMetadata) bool {
	if isComputedColumnType(a.Type) || isComputedColumnType(b.Type) {
		if !expressionsEqual(a.Type, b.Type) {
			return false
		}
	} else if !typesEqual(a.Type, b.Type) {
		return false
	}
	if a.Nullable != b.Nullable {
		return false
	}
	if !expressionsEqual(a.Default, b.Default) {
		return false
	}
	if a.Comment != b.Comment || a.CharacterSet != b.CharacterSet {
		return false
	}
	if a.Collation != b.Collation && a.Collation != "" && b.Collation != "" {
		return false
	}
	if a.OnUpdate != b.OnUpdate || a.DefaultOnNull != b.DefaultOnNull {
		return false
	}
	if !generationEqual(a.Generation, b.Generation) {
		return false
	}
	if a.IsIdentity != b.IsIdentity || a.IdentityGeneration != b.IdentityGeneration ||
		a.IdentitySeed != b.IdentitySeed || a.IdentityIncrement != b.IdentityIncrement {
		return false
	}
	if (a.Srid == nil) != (b.Srid == nil) || a.GetSrid() != b.GetSrid() {
		return false
	}
	return a.IsInvisible == b.IsInvisible
}

func typesEqual(a, b string) bool {
	return normalizeType(a) == normalizeType(b)
}

// normalizeType canonicalizes SQL Server built-in type names, synonyms, and
// implicit default parameters (for example datetime2 -> datetime2(7),
// decimal -> decimal(18,0), decimal(10) -> decimal(10,0), float -> float(53),
// real -> real(24), varchar -> varchar(1)).
func normalizeType(value string) string {
	value = strings.Join(strings.Fields(strings.ToLower(value)), "")
	if value == "" {
		return ""
	}
	base, params := value, ""
	if open := strings.IndexByte(value, '('); open >= 0 && strings.HasSuffix(value, ")") {
		base = value[:open]
		params = value[open:]
	}
	if len(base) >= 2 && base[0] == '[' && base[len(base)-1] == ']' && !strings.Contains(base, ".") {
		unbracketed := base[1 : len(base)-1]
		if isBuiltinSQLServerType(unbracketed) {
			base = unbracketed
		}
	}
	switch base {
	case "integer":
		base = "int"
	case "dec":
		base = "decimal"
	case "character":
		base = "char"
	case "charactervarying", "charvarying":
		base = "varchar"
	case "nationalchar", "nationalcharacter":
		base = "nchar"
	case "nationalcharvarying", "nationalcharactervarying":
		base = "nvarchar"
	case "binaryvarying":
		base = "varbinary"
	case "doubleprecision":
		if params == "" {
			return "float(53)"
		}
	}

	switch base {
	case "datetime2", "datetimeoffset", "time":
		if params == "" {
			return base + "(7)"
		}
	case "decimal", "numeric":
		if params == "" {
			return base + "(18,0)"
		}
		inner := params[1 : len(params)-1]
		if !strings.Contains(inner, ",") && inner != "" {
			return base + "(" + inner + ",0)"
		}
	case "float":
		if params == "" {
			return "float(53)"
		}
		inner := params[1 : len(params)-1]
		if n, err := strconv.Atoi(inner); err == nil {
			if n >= 1 && n <= 24 {
				return "real(24)"
			}
			if n >= 25 && n <= 53 {
				return "float(53)"
			}
		}
	case "real":
		if params == "" {
			return "real(24)"
		}
	case "char", "nchar", "varchar", "nvarchar", "binary", "varbinary":
		if params == "" {
			return base + "(1)"
		}
	}
	return base + params
}

func isBuiltinSQLServerType(name string) bool {
	switch name {
	case "bigint", "int", "integer", "smallint", "tinyint", "bit",
		"decimal", "dec", "numeric", "money", "smallmoney",
		"float", "real", "doubleprecision",
		"date", "time", "datetime", "datetime2", "smalldatetime", "datetimeoffset",
		"char", "character", "varchar", "charactervarying", "charvarying", "text",
		"nchar", "nationalchar", "nationalcharacter", "nvarchar", "nationalcharvarying", "nationalcharactervarying", "ntext",
		"binary", "varbinary", "binaryvarying", "image",
		"uniqueidentifier", "xml", "sql_variant", "timestamp", "rowversion",
		"hierarchyid", "geometry", "geography", "sysname":
		return true
	default:
		return false
	}
}

func generationEqual(a, b *metadata.GenerationMetadata) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Type == b.Type && expressionsEqual(a.Expression, b.Expression)
}

func compareIndexes(oldTable, newTable *metadata.TableMetadata) []*IndexDiffEntry {
	var changes []*IndexDiffEntry
	oldIndexes := oldTable.GetIndexes()
	newIndexes := newTable.GetIndexes()
	for _, oldIndex := range oldIndexes {
		if indexByName(newTable, oldIndex.Name) == nil {
			changes = append(changes, &IndexDiffEntry{Action: DiffDrop, OldIndex: oldIndex})
		}
	}
	for _, newIndex := range newIndexes {
		oldIndex := indexByName(oldTable, newIndex.Name)
		if oldIndex == nil {
			changes = append(changes, &IndexDiffEntry{Action: DiffCreate, NewIndex: newIndex})
			continue
		}
		if !indexesEqual(oldIndex, newIndex) {
			changes = append(changes,
				&IndexDiffEntry{Action: DiffDrop, OldIndex: oldIndex},
				&IndexDiffEntry{Action: DiffCreate, NewIndex: newIndex},
			)
		}
	}
	return changes
}

func indexesEqual(a, b *metadata.IndexMetadata) bool {
	if !indexTypesEqual(a, b) || a.Unique != b.Unique || a.Primary != b.Primary || a.IsConstraint != b.IsConstraint {
		return false
	}
	if len(a.Expressions) != len(b.Expressions) {
		return false
	}
	for i := range a.Expressions {
		if !expressionsEqual(a.Expressions[i], b.Expressions[i]) {
			return false
		}
	}
	if len(a.KeyLength) != len(b.KeyLength) {
		return false
	}
	for i := range a.KeyLength {
		if a.KeyLength[i] != b.KeyLength[i] {
			return false
		}
	}
	if !descendingArraysEqual(a.Descending, b.Descending) {
		return false
	}
	if a.Visible != b.Visible {
		return false
	}
	if a.Type == "XML" && a.Definition != b.Definition {
		return false
	}
	if !slices.Equal(indexIncludeColumns(a.Definition), indexIncludeColumns(b.Definition)) {
		return false
	}
	if !spatialConfigsEqual(a.SpatialConfig, b.SpatialConfig) {
		return false
	}
	return whereClausesEqual(a.Definition, b.Definition)
}

func indexTypesEqual(a, b *metadata.IndexMetadata) bool {
	if a.Type == b.Type {
		return true
	}
	// An omitted keyword takes SQL Server's default. A key constraint defaults
	// to CLUSTERED, every other index to NONCLUSTERED.
	effective := func(index *metadata.IndexMetadata) string {
		if index.Type != "" {
			return index.Type
		}
		if index.IsConstraint && index.Primary {
			return "CLUSTERED"
		}
		return "NONCLUSTERED"
	}
	return effective(a) == effective(b)
}

func descendingArraysEqual(a, b []bool) bool {
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
	if len(a) != len(b) {
		return false
	}
	return slices.Equal(a, b)
}

func spatialConfigsEqual(a, b *metadata.SpatialIndexConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Method != b.Method {
		return false
	}
	if !tessellationConfigsEqual(a.Tessellation, b.Tessellation) || !storageConfigsEqual(a.Storage, b.Storage) {
		return false
	}
	if !dimensionalConfigsEqual(a.Dimensional, b.Dimensional) {
		return false
	}
	if len(a.EngineSpecific) != len(b.EngineSpecific) {
		return false
	}
	for key, value := range a.EngineSpecific {
		other, ok := b.EngineSpecific[key]
		if !ok || value != other {
			return false
		}
	}
	return true
}

func tessellationConfigsEqual(a, b *metadata.TessellationConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Scheme != b.Scheme || a.CellsPerObject != b.CellsPerObject {
		return false
	}
	if (a.BoundingBox == nil) != (b.BoundingBox == nil) {
		return false
	}
	if a.BoundingBox != nil {
		if a.BoundingBox.Xmin != b.BoundingBox.Xmin || a.BoundingBox.Ymin != b.BoundingBox.Ymin ||
			a.BoundingBox.Xmax != b.BoundingBox.Xmax || a.BoundingBox.Ymax != b.BoundingBox.Ymax {
			return false
		}
	}
	if len(a.GridLevels) != len(b.GridLevels) {
		return false
	}
	for i := range a.GridLevels {
		if a.GridLevels[i].Level != b.GridLevels[i].Level || a.GridLevels[i].Density != b.GridLevels[i].Density {
			return false
		}
	}
	return true
}

func storageConfigsEqual(a, b *metadata.StorageConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Fillfactor == b.Fillfactor &&
		a.Buffering == b.Buffering &&
		a.Tablespace == b.Tablespace &&
		a.WorkTablespace == b.WorkTablespace &&
		a.SdoLevel == b.SdoLevel &&
		a.CommitInterval == b.CommitInterval &&
		a.PadIndex == b.PadIndex &&
		a.SortInTempdb == b.SortInTempdb &&
		a.DropExisting == b.DropExisting &&
		a.Online == b.Online &&
		a.AllowRowLocks == b.AllowRowLocks &&
		a.AllowPageLocks == b.AllowPageLocks &&
		a.Maxdop == b.Maxdop &&
		a.DataCompression == b.DataCompression
}

func dimensionalConfigsEqual(a, b *metadata.DimensionalConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Dimensions == b.Dimensions &&
		a.DataType == b.DataType &&
		a.OperatorClass == b.OperatorClass &&
		a.LayerGtype == b.LayerGtype &&
		a.ParallelBuild == b.ParallelBuild
}

func whereClausesEqual(a, b string) bool {
	return indexWhereClause(a) == indexWhereClause(b)
}

func compareForeignKeys(oldFKs, newFKs []*metadata.ForeignKeyMetadata) []*ForeignKeyDiffEntry {
	var changes []*ForeignKeyDiffEntry
	oldMap := make(map[string]*metadata.ForeignKeyMetadata)
	for _, fk := range oldFKs {
		oldMap[fk.Name] = fk
	}
	newMap := make(map[string]*metadata.ForeignKeyMetadata)
	for _, fk := range newFKs {
		newMap[fk.Name] = fk
	}
	for name, oldFK := range oldMap {
		if _, ok := newMap[name]; !ok {
			changes = append(changes, &ForeignKeyDiffEntry{Action: DiffDrop, OldForeignKey: oldFK})
		}
	}
	for name, newFK := range newMap {
		oldFK, ok := oldMap[name]
		if !ok {
			changes = append(changes, &ForeignKeyDiffEntry{Action: DiffCreate, NewForeignKey: newFK})
			continue
		}
		if !foreignKeysEqual(oldFK, newFK) {
			changes = append(changes,
				&ForeignKeyDiffEntry{Action: DiffDrop, OldForeignKey: oldFK},
				&ForeignKeyDiffEntry{Action: DiffCreate, NewForeignKey: newFK},
			)
		}
	}
	return changes
}

func foreignKeysEqual(a, b *metadata.ForeignKeyMetadata) bool {
	return a.ReferencedSchema == b.ReferencedSchema &&
		a.ReferencedTable == b.ReferencedTable &&
		a.OnDelete == b.OnDelete &&
		a.OnUpdate == b.OnUpdate &&
		a.MatchType == b.MatchType &&
		slices.Equal(a.Columns, b.Columns) &&
		slices.Equal(a.ReferencedColumns, b.ReferencedColumns)
}

func compareCheckConstraints(oldChecks, newChecks []*metadata.CheckConstraintMetadata) []*CheckConstraintDiffEntry {
	var changes []*CheckConstraintDiffEntry
	oldMap := make(map[string]*metadata.CheckConstraintMetadata)
	for _, check := range oldChecks {
		oldMap[check.Name] = check
	}
	newMap := make(map[string]*metadata.CheckConstraintMetadata)
	for _, check := range newChecks {
		newMap[check.Name] = check
	}
	for name, oldCheck := range oldMap {
		if _, ok := newMap[name]; !ok {
			changes = append(changes, &CheckConstraintDiffEntry{Action: DiffDrop, OldCheckConstraint: oldCheck})
		}
	}
	for name, newCheck := range newMap {
		oldCheck, ok := oldMap[name]
		if !ok {
			changes = append(changes, &CheckConstraintDiffEntry{Action: DiffCreate, NewCheckConstraint: newCheck})
			continue
		}
		if !checkExpressionsEqual(oldCheck.Expression, newCheck.Expression) {
			changes = append(changes,
				&CheckConstraintDiffEntry{Action: DiffDrop, OldCheckConstraint: oldCheck},
				&CheckConstraintDiffEntry{Action: DiffCreate, NewCheckConstraint: newCheck},
			)
		}
	}
	return changes
}

func compareTriggers(oldTriggers, newTriggers []*metadata.TriggerMetadata) []*TriggerDiffEntry {
	var changes []*TriggerDiffEntry
	oldMap := make(map[string]*metadata.TriggerMetadata)
	for _, trigger := range oldTriggers {
		if trigger.GetSkipDump() {
			continue
		}
		oldMap[trigger.Name] = trigger
	}
	newMap := make(map[string]*metadata.TriggerMetadata)
	for _, trigger := range newTriggers {
		if trigger.GetSkipDump() {
			continue
		}
		newMap[trigger.Name] = trigger
	}
	for name, oldTrigger := range oldMap {
		if _, ok := newMap[name]; !ok {
			changes = append(changes, &TriggerDiffEntry{Action: DiffDrop, OldTrigger: oldTrigger})
		}
	}
	for name, newTrigger := range newMap {
		oldTrigger, ok := oldMap[name]
		if !ok {
			changes = append(changes, &TriggerDiffEntry{Action: DiffCreate, NewTrigger: newTrigger})
			continue
		}
		if !triggersEqual(oldTrigger, newTrigger) {
			changes = append(changes,
				&TriggerDiffEntry{Action: DiffDrop, OldTrigger: oldTrigger},
				&TriggerDiffEntry{Action: DiffCreate, NewTrigger: newTrigger},
			)
		}
	}
	return changes
}

func triggersEqual(a, b *metadata.TriggerMetadata) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Name != b.Name || a.Event != b.Event || a.Timing != b.Timing {
		return false
	}
	if a.Body == b.Body {
		return true
	}
	return normalizeTriggerBody(a.Body) == normalizeTriggerBody(b.Body)
}

// normalizeTriggerBody is Bytebase's whitespace- and case-insensitive trigger
// body comparison.
func normalizeTriggerBody(body string) string {
	return normalizeExpression(body)
}

func compareViews(diff *SchemaDiff, schemaName string, oldSchema, newSchema *metadata.SchemaMetadata) {
	for _, name := range sortedViewNames(oldSchema) {
		oldView := viewByName(oldSchema, name)
		if viewByName(newSchema, name) == nil {
			if oldView != nil && !oldView.GetSkipDump() {
				diff.Views = append(diff.Views, &ViewDiffEntry{Action: DiffDrop, SchemaName: schemaName, ViewName: name, OldView: oldView})
			}
		}
	}
	for _, name := range sortedViewNames(newSchema) {
		newView := viewByName(newSchema, name)
		if newView == nil || newView.GetSkipDump() {
			continue
		}
		oldView := viewByName(oldSchema, name)
		if oldView == nil {
			diff.Views = append(diff.Views, &ViewDiffEntry{Action: DiffCreate, SchemaName: schemaName, ViewName: name, NewView: newView})
			continue
		}
		if !oldView.GetSkipDump() && (oldView.Definition != newView.Definition || oldView.Comment != newView.Comment || !triggerListsEqual(oldView.Triggers, newView.Triggers)) {
			diff.Views = append(diff.Views, &ViewDiffEntry{Action: DiffAlter, SchemaName: schemaName, ViewName: name, OldView: oldView, NewView: newView})
		}
	}
}

func triggerListsEqual(a, b []*metadata.TriggerMetadata) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !triggersEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func compareSynonyms(diff *SchemaDiff, schemaName string, oldSchema, newSchema *metadata.SchemaMetadata) {
	oldMap := make(map[string]*metadata.SynonymMetadata)
	for _, synonym := range oldSchema.Synonyms {
		oldMap[synonym.Name] = synonym
	}
	newMap := make(map[string]*metadata.SynonymMetadata)
	for _, synonym := range newSchema.Synonyms {
		newMap[synonym.Name] = synonym
	}
	for name, oldSynonym := range oldMap {
		if _, ok := newMap[name]; !ok {
			diff.Synonyms = append(diff.Synonyms, &SynonymDiffEntry{Action: DiffDrop, SchemaName: schemaName, SynonymName: name, OldSynonym: oldSynonym})
		}
	}
	for name, newSynonym := range newMap {
		oldSynonym, ok := oldMap[name]
		if !ok {
			diff.Synonyms = append(diff.Synonyms, &SynonymDiffEntry{Action: DiffCreate, SchemaName: schemaName, SynonymName: name, NewSynonym: newSynonym})
			continue
		}
		if oldSynonym.TargetServer != newSynonym.TargetServer ||
			oldSynonym.TargetDatabase != newSynonym.TargetDatabase ||
			oldSynonym.TargetSchema != newSynonym.TargetSchema ||
			oldSynonym.TargetName != newSynonym.TargetName {
			diff.Synonyms = append(diff.Synonyms,
				&SynonymDiffEntry{Action: DiffDrop, SchemaName: schemaName, SynonymName: name, OldSynonym: oldSynonym},
				&SynonymDiffEntry{Action: DiffCreate, SchemaName: schemaName, SynonymName: name, NewSynonym: newSynonym},
			)
		}
	}
}

func compareFunctions(diff *SchemaDiff, schemaName string, oldSchema, newSchema *metadata.SchemaMetadata) {
	oldFunctions := functionsBySignature(oldSchema)
	newFunctions := functionsBySignature(newSchema)
	for signature, oldFunction := range oldFunctions {
		if _, ok := newFunctions[signature]; !ok {
			diff.Functions = append(diff.Functions, &FunctionDiffEntry{
				Action:       DiffDrop,
				SchemaName:   schemaName,
				FunctionName: signature,
				OldFunction:  oldFunction,
			})
		}
	}
	for signature, newFunction := range newFunctions {
		oldFunction, ok := oldFunctions[signature]
		if !ok {
			diff.Functions = append(diff.Functions, &FunctionDiffEntry{
				Action:       DiffCreate,
				SchemaName:   schemaName,
				FunctionName: signature,
				NewFunction:  newFunction,
			})
			continue
		}
		if functionsEqual(oldFunction, newFunction) {
			continue
		}
		diff.Functions = append(diff.Functions,
			&FunctionDiffEntry{Action: DiffDrop, SchemaName: schemaName, FunctionName: signature, OldFunction: oldFunction},
			&FunctionDiffEntry{Action: DiffCreate, SchemaName: schemaName, FunctionName: signature, NewFunction: newFunction},
		)
	}
}

func functionsBySignature(schema *metadata.SchemaMetadata) map[string]*metadata.FunctionMetadata {
	result := make(map[string]*metadata.FunctionMetadata)
	for _, function := range schema.GetFunctions() {
		if function.GetSkipDump() {
			continue
		}
		signature := function.Signature
		if signature == "" {
			signature = function.Name
		}
		result[signature] = function
	}
	return result
}

func functionsEqual(a, b *metadata.FunctionMetadata) bool {
	if a.Definition != b.Definition && normalizeFunctionDefinition(a.Definition) != normalizeFunctionDefinition(b.Definition) {
		return false
	}
	return a.CharacterSetClient == b.CharacterSetClient &&
		a.CollationConnection == b.CollationConnection &&
		a.DatabaseCollation == b.DatabaseCollation &&
		a.SqlMode == b.SqlMode &&
		a.Comment == b.Comment
}

func compareProcedures(diff *SchemaDiff, schemaName string, oldSchema, newSchema *metadata.SchemaMetadata) {
	for _, name := range sortedProcedureNames(oldSchema) {
		oldProcedure := procedureByName(oldSchema, name)
		if procedureByName(newSchema, name) == nil {
			if oldProcedure != nil && !oldProcedure.SkipDump {
				diff.Procedures = append(diff.Procedures, &ProcedureDiffEntry{
					Action:        DiffDrop,
					SchemaName:    schemaName,
					ProcedureName: name,
					OldProcedure:  oldProcedure,
				})
			}
		}
	}
	for _, name := range sortedProcedureNames(newSchema) {
		newProcedure := procedureByName(newSchema, name)
		if newProcedure == nil || newProcedure.SkipDump {
			continue
		}
		oldProcedure := procedureByName(oldSchema, name)
		if oldProcedure == nil {
			diff.Procedures = append(diff.Procedures, &ProcedureDiffEntry{
				Action:        DiffCreate,
				SchemaName:    schemaName,
				ProcedureName: name,
				NewProcedure:  newProcedure,
			})
			continue
		}
		if !oldProcedure.SkipDump && !procedureDefinitionsEqual(oldProcedure.Definition, newProcedure.Definition) {
			diff.Procedures = append(diff.Procedures, &ProcedureDiffEntry{
				Action:        DiffAlter,
				SchemaName:    schemaName,
				ProcedureName: name,
				OldProcedure:  oldProcedure,
				NewProcedure:  newProcedure,
			})
		}
	}
}

func compareSequences(diff *SchemaDiff, schemaName string, oldSchema, newSchema *metadata.SchemaMetadata) {
	oldSequences := sequencesByName(oldSchema)
	newSequences := sequencesByName(newSchema)
	for name, oldSequence := range oldSequences {
		if _, ok := newSequences[name]; !ok {
			diff.Sequences = append(diff.Sequences, &SequenceDiffEntry{
				Action:       DiffDrop,
				SchemaName:   schemaName,
				SequenceName: name,
				OldSequence:  oldSequence,
			})
		}
	}
	for name, newSequence := range newSequences {
		oldSequence, ok := oldSequences[name]
		if !ok {
			diff.Sequences = append(diff.Sequences, &SequenceDiffEntry{
				Action:       DiffCreate,
				SchemaName:   schemaName,
				SequenceName: name,
				NewSequence:  newSequence,
			})
			continue
		}
		if !sequencesEqual(oldSequence, newSequence) {
			diff.Sequences = append(diff.Sequences, &SequenceDiffEntry{
				Action:       DiffAlter,
				SchemaName:   schemaName,
				SequenceName: name,
				OldSequence:  oldSequence,
				NewSequence:  newSequence,
			})
		}
	}
}

func sequencesEqual(a, b *metadata.SequenceMetadata) bool {
	return typesEqual(a.DataType, b.DataType) &&
		a.Start == b.Start &&
		a.MinValue == b.MinValue &&
		a.MaxValue == b.MaxValue &&
		a.Increment == b.Increment &&
		a.Cycle == b.Cycle &&
		a.CacheSize == b.CacheSize &&
		a.OwnerTable == b.OwnerTable &&
		a.OwnerColumn == b.OwnerColumn &&
		a.Comment == b.Comment
}

// checkExpressionsEqual compares two CHECK expressions ignoring outer
// parenthesization, whitespace, and non-literal case.
func checkExpressionsEqual(a, b string) bool {
	return expressionsEqual(a, b)
}

// normalizeFunctionDefinition is Bytebase's function comparison normalization,
// ported because Bytebase applies it to every engine, SQL Server included.
func normalizeFunctionDefinition(definition string) string {
	return normalizeExpression(definition)
}

func procedureDefinitionsEqual(a, b string) bool {
	if a == b {
		return true
	}
	return normalizeFunctionDefinition(a) == normalizeFunctionDefinition(b)
}

// stripOuterParens removes redundant balanced outer parentheses around an
// expression without touching grouping parentheses inside binary expressions
// such as "(1 + 2) * 3".
func stripOuterParens(expression string) string {
	expression = strings.TrimSpace(expression)
	for strings.HasPrefix(expression, "(") && strings.HasSuffix(expression, ")") && balancedOuterParens(expression) {
		expression = strings.TrimSpace(expression[1 : len(expression)-1])
	}
	return expression
}

func expressionsEqual(a, b string) bool {
	if a == b {
		return true
	}
	if strings.TrimSpace(a) == "" && strings.TrimSpace(b) == "" {
		return true
	}
	if strings.TrimSpace(a) == "" || strings.TrimSpace(b) == "" {
		return false
	}
	return normalizeExpression(stripOuterParens(a)) == normalizeExpression(stripOuterParens(b))
}

// normalizeExpression collapses whitespace and folds case outside string
// literals only, so two defaults that differ inside quotes stay different.
func normalizeExpression(expression string) string {
	var out strings.Builder
	inLiteral := false
	space := false
	for i := 0; i < len(expression); i++ {
		ch := expression[i]
		if inLiteral {
			out.WriteByte(ch)
			if ch == '\'' {
				if i+1 < len(expression) && expression[i+1] == '\'' {
					out.WriteByte('\'')
					i++
					continue
				}
				inLiteral = false
			}
			continue
		}
		if ch == '-' && i+1 < len(expression) && expression[i+1] == '-' {
			for i < len(expression) && expression[i] != '\n' {
				out.WriteByte(expression[i])
				i++
			}
			if i < len(expression) {
				out.WriteByte('\n')
			}
			space = false
			continue
		}
		if ch == '/' && i+1 < len(expression) && expression[i+1] == '*' {
			for i < len(expression) {
				out.WriteByte(expression[i])
				if expression[i] == '*' && i+1 < len(expression) && expression[i+1] == '/' {
					out.WriteByte('/')
					i++
					break
				}
				i++
			}
			space = false
			continue
		}
		switch {
		case ch == '[':
			if space {
				out.WriteByte(' ')
				space = false
			}
			for i < len(expression) {
				out.WriteByte(expression[i])
				if expression[i] == ']' {
					if i+1 < len(expression) && expression[i+1] == ']' {
						out.WriteByte(']')
						i += 2
						continue
					}
					i++
					break
				}
				i++
			}
			i--
		case ch == '"':
			if space {
				out.WriteByte(' ')
				space = false
			}
			out.WriteByte(ch)
			for i++; i < len(expression); i++ {
				out.WriteByte(expression[i])
				if expression[i] == '"' {
					if i+1 < len(expression) && expression[i+1] == '"' {
						out.WriteByte('"')
						i++
						continue
					}
					break
				}
			}
		case ch == '\'':
			if space {
				out.WriteByte(' ')
				space = false
			}
			inLiteral = true
			out.WriteByte(ch)
		case ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n':
			space = true
		default:
			if space {
				out.WriteByte(' ')
				space = false
			}
			if ch >= 'a' && ch <= 'z' {
				ch -= 'a' - 'A'
			}
			out.WriteByte(ch)
		}
	}
	return strings.TrimSpace(out.String())
}

func sortDiff(diff *SchemaDiff) {
	slices.SortFunc(diff.Schemas, func(a, b *SchemaDiffEntry) int {
		if a.SchemaName != b.SchemaName {
			return strings.Compare(a.SchemaName, b.SchemaName)
		}
		return strings.Compare(string(a.Action), string(b.Action))
	})
	slices.SortFunc(diff.Tables, func(a, b *TableDiffEntry) int {
		if a.SchemaName != b.SchemaName {
			return strings.Compare(a.SchemaName, b.SchemaName)
		}
		if a.TableName != b.TableName {
			return strings.Compare(a.TableName, b.TableName)
		}
		return strings.Compare(string(a.Action), string(b.Action))
	})
	slices.SortFunc(diff.Synonyms, func(a, b *SynonymDiffEntry) int {
		if a.SchemaName != b.SchemaName {
			return strings.Compare(a.SchemaName, b.SchemaName)
		}
		if a.SynonymName != b.SynonymName {
			return strings.Compare(a.SynonymName, b.SynonymName)
		}
		return strings.Compare(string(a.Action), string(b.Action))
	})
	slices.SortFunc(diff.Views, func(a, b *ViewDiffEntry) int {
		if a.SchemaName != b.SchemaName {
			return strings.Compare(a.SchemaName, b.SchemaName)
		}
		if a.ViewName != b.ViewName {
			return strings.Compare(a.ViewName, b.ViewName)
		}
		return strings.Compare(string(a.Action), string(b.Action))
	})
	slices.SortFunc(diff.Functions, func(a, b *FunctionDiffEntry) int {
		if a.SchemaName != b.SchemaName {
			return strings.Compare(a.SchemaName, b.SchemaName)
		}
		if a.FunctionName != b.FunctionName {
			return strings.Compare(a.FunctionName, b.FunctionName)
		}
		return strings.Compare(string(a.Action), string(b.Action))
	})
	slices.SortFunc(diff.Procedures, func(a, b *ProcedureDiffEntry) int {
		if a.SchemaName != b.SchemaName {
			return strings.Compare(a.SchemaName, b.SchemaName)
		}
		if a.ProcedureName != b.ProcedureName {
			return strings.Compare(a.ProcedureName, b.ProcedureName)
		}
		return strings.Compare(string(a.Action), string(b.Action))
	})
	slices.SortFunc(diff.Sequences, func(a, b *SequenceDiffEntry) int {
		if a.SchemaName != b.SchemaName {
			return strings.Compare(a.SchemaName, b.SchemaName)
		}
		if a.SequenceName != b.SequenceName {
			return strings.Compare(a.SequenceName, b.SequenceName)
		}
		return strings.Compare(string(a.Action), string(b.Action))
	})
	for _, table := range diff.Tables {
		sortTableDiff(table)
	}
}

func sortTableDiff(table *TableDiffEntry) {
	slices.SortFunc(table.Columns, func(a, b *ColumnDiffEntry) int {
		aPos, aName := columnPositionName(a)
		bPos, bName := columnPositionName(b)
		if aPos > 0 && bPos > 0 && aPos != bPos {
			return int(aPos - bPos)
		}
		if aName != bName {
			return strings.Compare(aName, bName)
		}
		return strings.Compare(string(a.Action), string(b.Action))
	})
	slices.SortFunc(table.Indexes, func(a, b *IndexDiffEntry) int {
		return compareNamed(string(a.Action), indexDiffName(a), string(b.Action), indexDiffName(b))
	})
	slices.SortFunc(table.ForeignKeys, func(a, b *ForeignKeyDiffEntry) int {
		return compareNamed(string(a.Action), fkDiffName(a), string(b.Action), fkDiffName(b))
	})
	slices.SortFunc(table.Checks, func(a, b *CheckConstraintDiffEntry) int {
		return compareNamed(string(a.Action), checkDiffName(a), string(b.Action), checkDiffName(b))
	})
	slices.SortFunc(table.Triggers, func(a, b *TriggerDiffEntry) int {
		return compareNamed(string(a.Action), triggerDiffName(a), string(b.Action), triggerDiffName(b))
	})
}

func compareNamed(aAction, aName, bAction, bName string) int {
	if aName != bName {
		return strings.Compare(aName, bName)
	}
	return strings.Compare(aAction, bAction)
}

func columnPositionName(diff *ColumnDiffEntry) (int32, string) {
	if diff.NewColumn != nil {
		return diff.NewColumn.Position, diff.NewColumn.Name
	}
	if diff.OldColumn != nil {
		return diff.OldColumn.Position, diff.OldColumn.Name
	}
	return 0, ""
}

func indexDiffName(diff *IndexDiffEntry) string {
	if diff.NewIndex != nil {
		return diff.NewIndex.Name
	}
	if diff.OldIndex != nil {
		return diff.OldIndex.Name
	}
	return ""
}

func fkDiffName(diff *ForeignKeyDiffEntry) string {
	if diff.NewForeignKey != nil {
		return diff.NewForeignKey.Name
	}
	if diff.OldForeignKey != nil {
		return diff.OldForeignKey.Name
	}
	return ""
}

func checkDiffName(diff *CheckConstraintDiffEntry) string {
	if diff.NewCheckConstraint != nil {
		return diff.NewCheckConstraint.Name
	}
	if diff.OldCheckConstraint != nil {
		return diff.OldCheckConstraint.Name
	}
	return ""
}

func triggerDiffName(diff *TriggerDiffEntry) string {
	if diff.NewTrigger != nil {
		return diff.NewTrigger.Name
	}
	if diff.OldTrigger != nil {
		return diff.OldTrigger.Name
	}
	return ""
}

func sortedSchemaNames(schemas []*metadata.SchemaMetadata) []string {
	names := make([]string, 0, len(schemas))
	for _, schema := range schemas {
		names = append(names, schema.Name)
	}
	slices.Sort(names)
	return names
}

func sortedTableNames(schema *metadata.SchemaMetadata) []string {
	names := make([]string, 0, len(schema.Tables))
	for _, table := range schema.Tables {
		names = append(names, table.Name)
	}
	slices.Sort(names)
	return names
}

func sortedViewNames(schema *metadata.SchemaMetadata) []string {
	names := make([]string, 0, len(schema.Views))
	for _, view := range schema.Views {
		names = append(names, view.Name)
	}
	slices.Sort(names)
	return names
}

func sortedProcedureNames(schema *metadata.SchemaMetadata) []string {
	names := make([]string, 0, len(schema.Procedures))
	for _, procedure := range schema.Procedures {
		names = append(names, procedure.Name)
	}
	slices.Sort(names)
	return names
}

func tableByName(schema *metadata.SchemaMetadata, name string) *metadata.TableMetadata {
	for _, table := range schema.Tables {
		if table.Name == name {
			return table
		}
	}
	for _, table := range schema.Tables {
		if strings.EqualFold(table.Name, name) {
			return table
		}
	}
	return nil
}

func viewByName(schema *metadata.SchemaMetadata, name string) *metadata.ViewMetadata {
	for _, view := range schema.Views {
		if view.Name == name {
			return view
		}
	}
	for _, view := range schema.Views {
		if strings.EqualFold(view.Name, name) {
			return view
		}
	}
	return nil
}

func procedureByName(schema *metadata.SchemaMetadata, name string) *metadata.ProcedureMetadata {
	for _, procedure := range schema.Procedures {
		if strings.EqualFold(procedure.Name, name) {
			return procedure
		}
	}
	return nil
}

func columnByName(table *metadata.TableMetadata, name string) *metadata.ColumnMetadata {
	for _, column := range table.Columns {
		if column.Name == name {
			return column
		}
	}
	for _, column := range table.Columns {
		if strings.EqualFold(column.Name, name) {
			return column
		}
	}
	return nil
}

func indexByName(table *metadata.TableMetadata, name string) *metadata.IndexMetadata {
	for _, index := range table.Indexes {
		if strings.EqualFold(index.Name, name) {
			return index
		}
	}
	return nil
}

func sequencesByName(schema *metadata.SchemaMetadata) map[string]*metadata.SequenceMetadata {
	result := make(map[string]*metadata.SequenceMetadata)
	for _, sequence := range schema.Sequences {
		if sequence.GetSkipDump() {
			continue
		}
		result[sequence.Name] = sequence
	}
	return result
}
