package catalog

import (
	"fmt"
	"slices"
	"strings"

	crparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	plpgsqlparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	metadatapb "github.com/bytebase/omni/metadata"
)

// MetadataMigrationOptions configures MetadataMigration.
type MetadataMigrationOptions struct {
	// SkipSchemas lists schema names left out of the migration, such as a
	// backup schema.
	SkipSchemas []string
}

// MigrationOpType classifies one statement of a migration.
type MigrationOpType string

// The statements a migration holds.
const (
	OpCreateSchema    MigrationOpType = "CreateSchema"
	OpDropSchema      MigrationOpType = "DropSchema"
	OpCreateExtension MigrationOpType = "CreateExtension"
	OpDropExtension   MigrationOpType = "DropExtension"
	OpCreateType      MigrationOpType = "CreateType"
	OpDropType        MigrationOpType = "DropType"
	OpAlterType       MigrationOpType = "AlterType"
	OpCreateSequence  MigrationOpType = "CreateSequence"
	OpDropSequence    MigrationOpType = "DropSequence"
	OpAlterSequence   MigrationOpType = "AlterSequence"
	OpCreateTable     MigrationOpType = "CreateTable"
	OpDropTable       MigrationOpType = "DropTable"
	OpAddColumn       MigrationOpType = "AddColumn"
	OpDropColumn      MigrationOpType = "DropColumn"
	OpAlterColumn     MigrationOpType = "AlterColumn"
	OpAlterTable      MigrationOpType = "AlterTable"
	OpAddConstraint   MigrationOpType = "AddConstraint"
	OpDropConstraint  MigrationOpType = "DropConstraint"
	OpCreateIndex     MigrationOpType = "CreateIndex"
	OpDropIndex       MigrationOpType = "DropIndex"
	OpCreateView      MigrationOpType = "CreateView"
	OpDropView        MigrationOpType = "DropView"
	OpCreateFunction  MigrationOpType = "CreateFunction"
	OpDropFunction    MigrationOpType = "DropFunction"
	OpAlterFunction   MigrationOpType = "AlterFunction"
	OpComment         MigrationOpType = "Comment"
)

// MigrationOp is one statement of a migration. ParentObject names the table
// of a column, constraint or index.
type MigrationOp struct {
	Type         MigrationOpType
	SchemaName   string
	ObjectName   string
	ParentObject string
	SQL          string
}

// MigrationPlan is a migration's statements in the order they run.
type MigrationPlan struct {
	Ops []MigrationOp
}

// SQL returns the statements joined with ";\n".
func (p *MigrationPlan) SQL() string {
	if p == nil || len(p.Ops) == 0 {
		return ""
	}
	parts := make([]string, len(p.Ops))
	for i, op := range p.Ops {
		parts[i] = op.SQL
	}
	return strings.Join(parts, ";\n")
}

// MetadataMigration generates the migration from one synced CockroachDB
// snapshot to another.
func MetadataMigration(oldMeta, newMeta *metadatapb.DatabaseSchemaMetadata, opts MetadataMigrationOptions) (string, error) {
	plan, err := MetadataMigrationPlan(oldMeta, newMeta, opts)
	if err != nil {
		return "", err
	}
	return plan.SQL(), nil
}

// MetadataMigrationPlan is MetadataMigration's plan. It compares the two
// snapshots object by object, each in the form the definitions render, so
// two databases of different names compare equal where their objects do.
//
// The statements run in phases: foreign keys, the views a change rebuilds and
// the indexes and checks that change are dropped first; then removed tables,
// columns, routines, sequences, types and schemas; then new schemas, types and
// sequences, the columns of kept tables, and new tables, views and routines
// in dependency order; then indexes, checks, foreign keys and comments.
func MetadataMigrationPlan(oldMeta, newMeta *metadatapb.DatabaseSchemaMetadata, opts MetadataMigrationOptions) (*MigrationPlan, error) {
	skip := make(map[string]bool)
	for _, name := range opts.SkipSchemas {
		skip[name] = true
	}
	d := &differ{
		from: newSnapshot(oldMeta, skip),
		to:   newSnapshot(newMeta, skip),
	}
	if err := d.diff(); err != nil {
		return nil, err
	}
	var ops []MigrationOp
	for _, phase := range d.phases {
		ops = append(ops, phase...)
	}
	return &MigrationPlan{Ops: ops}, nil
}

// snapshot is one side of a migration, its schemas keyed by name.
type snapshot struct {
	meta     *metadatapb.DatabaseSchemaMetadata
	database string
	schemas  []*metadatapb.SchemaMetadata
	byName   map[string]*metadatapb.SchemaMetadata
}

func newSnapshot(meta *metadatapb.DatabaseSchemaMetadata, skip map[string]bool) *snapshot {
	s := &snapshot{meta: meta, byName: make(map[string]*metadatapb.SchemaMetadata)}
	if meta == nil {
		s.meta = &metadatapb.DatabaseSchemaMetadata{}
		return s
	}
	s.database = meta.Name
	for _, schema := range meta.Schemas {
		if skip[schema.Name] || schema.SkipDump {
			continue
		}
		s.schemas = append(s.schemas, schema)
		s.byName[schema.Name] = schema
	}
	return s
}

func (s *snapshot) table(schema, name string) *metadatapb.TableMetadata {
	for _, t := range s.byName[schema].GetTables() {
		if t.Name == name {
			return t
		}
	}
	return nil
}

// phase indexes: the order the statements run in.
const (
	phaseDetachDefaults = iota
	phaseDropForeignKeys
	phaseDropRoutinesFirst
	phaseDropViews
	phaseDropIndexes
	phaseDropTables
	phaseDropColumns
	phaseDropRoutines
	phaseDropSequences
	phaseDropTypes
	phaseDropSchemas
	phaseCreateSchemas
	phaseCreateTypes
	phaseCreateSequences
	phaseAlterTables
	phaseCreateObjects
	phaseLateColumns
	phaseCreateIndexes
	phaseAddForeignKeys
	phaseComments
	phaseCount
)

type differ struct {
	from, to *snapshot
	phases   [phaseCount][]MigrationOp
	// rebuilt holds the views and materialized views a change drops and
	// creates again, keyed by object ID.
	rebuilt map[string]bool
	// newRoutines holds the names, bare and schema-qualified, of the routines
	// the target adds.
	newRoutines map[string]bool
	// droppedRoutines holds the names, bare and schema-qualified, of the
	// source routines the migration drops: removed, or recreated.
	droppedRoutines map[string]bool
	// recreated holds the routines, by object ID, that are dropped and
	// created again rather than replaced in place.
	recreated map[string]bool
	// r renders the target's objects; its types are both sides'.
	r renderer
	// cascaded holds the foreign keys cascadedForeignKeys returns.
	cascaded map[string]bool
	// droppedSequences holds the schema-qualified names of the sequences
	// the target removes.
	droppedSequences map[string]bool
}

func (d *differ) add(phase int, op MigrationOp) {
	d.phases[phase] = append(d.phases[phase], op)
}

func (d *differ) diff() error {
	d.r = newRenderer(d.to.database, append(slices.Clone(d.to.schemas), d.from.schemas...))
	d.diffSchemas()
	d.diffExtensions()
	if err := d.diffTypes(); err != nil {
		return err
	}
	d.diffSequences()
	d.rebuilt = d.rebuiltViews()
	d.newRoutines = make(map[string]bool)
	for _, schema := range d.to.schemas {
		for _, function := range schema.Functions {
			if findFunction(d.from.byName[schema.Name], metadataFuncIdentity(function)) == nil {
				d.newRoutines[function.Name] = true
				d.newRoutines[schema.Name+"."+function.Name] = true
			}
		}
	}
	d.planRoutines()
	if err := d.diffTables(); err != nil {
		return err
	}
	d.diffViews()
	d.diffRoutines()
	d.createObjects()
	return nil
}

func (d *differ) diffSchemas() {
	for _, schema := range d.from.schemas {
		if d.to.byName[schema.Name] == nil && schema.Name != "public" {
			d.add(phaseDropSchemas, MigrationOp{Type: OpDropSchema, ObjectName: schema.Name,
				SQL: "DROP SCHEMA " + quoteIdent(schema.Name)})
		}
	}
	for _, schema := range d.to.schemas {
		old := d.from.byName[schema.Name]
		if old == nil && schema.Name != "public" {
			d.add(phaseCreateSchemas, MigrationOp{Type: OpCreateSchema, ObjectName: schema.Name,
				SQL: "CREATE SCHEMA " + quoteIdent(schema.Name)})
		}
		if old.GetComment() != schema.Comment {
			d.comment(schema.Name, schema.Name, "", "SCHEMA "+quoteIdent(schema.Name), schema.Comment)
		}
	}
}

func (d *differ) diffExtensions() {
	names := func(meta *metadatapb.DatabaseSchemaMetadata) map[string]*metadatapb.ExtensionMetadata {
		m := make(map[string]*metadatapb.ExtensionMetadata)
		for _, extension := range meta.Extensions {
			m[extension.Name] = extension
		}
		return m
	}
	from, to := names(d.from.meta), names(d.to.meta)
	for _, extension := range d.from.meta.Extensions {
		if to[extension.Name] == nil {
			d.add(phaseDropTypes, MigrationOp{Type: OpDropExtension, ObjectName: extension.Name,
				SQL: "DROP EXTENSION " + quoteIdent(extension.Name)})
		}
	}
	for _, extension := range d.to.meta.Extensions {
		if from[extension.Name] == nil {
			var buf strings.Builder
			writeExtension(&buf, extension)
			d.add(phaseCreateSchemas, MigrationOp{Type: OpCreateExtension, ObjectName: extension.Name, SQL: statement(buf.String())})
		}
	}
}

func (d *differ) diffTypes() error {
	// Types drop dependents first: a composite type before the types its
	// attributes use.
	dropped := make(map[string]bool)
	for _, schema := range d.from.schemas {
		for _, t := range collectTypes([]*metadatapb.SchemaMetadata{schema}, false) {
			old := d.toType(schema.Name, t)
			if old == nil || (t.enum == nil || old.enum == nil) && d.typeDefinition(*old) != d.typeDefinition(t) {
				dropped[getObjectID(schema.Name, t.name())] = true
			}
		}
	}
	all := collectTypes(d.from.schemas, false)
	for i := len(all) - 1; i >= 0; i-- {
		t := all[i]
		if dropped[getObjectID(t.schema, t.name())] {
			d.add(phaseDropTypes, MigrationOp{Type: OpDropType, SchemaName: t.schema, ObjectName: t.name(),
				SQL: fmt.Sprintf("DROP TYPE %s.%s", quoteIdent(t.schema), quoteIdent(t.name()))})
		}
	}
	var created []qualifiedType
	for _, schema := range d.to.schemas {
		for _, t := range collectTypes([]*metadatapb.SchemaMetadata{schema}, false) {
			old := d.fromType(schema.Name, t)
			switch {
			case old == nil:
				created = append(created, t)
			case t.enum != nil && old.enum != nil:
				if err := d.alterEnum(schema.Name, old.enum, t.enum); err != nil {
					return err
				}
			case d.typeDefinition(*old) != d.typeDefinition(t):
				// A composite type has no ALTER that CockroachDB supports:
				// it was dropped above and is created again.
				created = append(created, t)
			}
			oldComment := ""
			if old != nil {
				oldComment = typeComment(*old)
			}
			if typeComment(t) != oldComment {
				d.comment(schema.Name, t.name(), "", fmt.Sprintf("TYPE %s.%s", quoteIdent(schema.Name), quoteIdent(t.name())), typeComment(t))
			}
		}
	}
	for _, t := range orderTypes(created) {
		d.add(phaseCreateTypes, MigrationOp{Type: OpCreateType, SchemaName: t.schema, ObjectName: t.name(), SQL: d.typeDefinition(t)})
	}
	return nil
}

func (d *differ) fromType(schema string, t qualifiedType) *qualifiedType {
	return findType(d.from.byName[schema], t)
}

func (d *differ) toType(schema string, t qualifiedType) *qualifiedType {
	return findType(d.to.byName[schema], t)
}

func findType(schema *metadatapb.SchemaMetadata, t qualifiedType) *qualifiedType {
	if schema == nil {
		return nil
	}
	for _, candidate := range collectTypes([]*metadatapb.SchemaMetadata{schema}, false) {
		if candidate.name() == t.name() {
			return &candidate
		}
	}
	return nil
}

func (d *differ) typeDefinition(t qualifiedType) string {
	var buf strings.Builder
	if t.enum != nil {
		writeEnum(&buf, t.schema, t.enum)
	} else {
		d.r.writeCompositeType(&buf, t.schema, t.composite)
	}
	return buf.String()
}

func typeComment(t qualifiedType) string {
	if t.enum != nil {
		return t.enum.Comment
	}
	return t.composite.Comment
}

// alterEnum adds and drops an enum's values in place. Values both sides hold
// must keep their order; CockroachDB cannot reorder them.
func (d *differ) alterEnum(schema string, old, new *metadatapb.EnumTypeMetadata) error {
	name := fmt.Sprintf("%s.%s", quoteIdent(schema), quoteIdent(new.Name))
	var kept []string
	for _, value := range old.Values {
		if slices.Contains(new.Values, value) {
			kept = append(kept, value)
		} else {
			d.add(phaseCreateTypes, MigrationOp{Type: OpAlterType, SchemaName: schema, ObjectName: new.Name,
				SQL: fmt.Sprintf("ALTER TYPE %s DROP VALUE '%s'", name, escapeSingleQuote(value))})
		}
	}
	var common []string
	for _, value := range new.Values {
		if slices.Contains(old.Values, value) {
			common = append(common, value)
		}
	}
	if !slices.Equal(kept, common) {
		return fmt.Errorf("enum %s.%s reorders its values, which CockroachDB cannot do in place", schema, new.Name)
	}
	// Each value is placed after the one before it, which is kept or was
	// added just before; leading values go before the first kept value.
	for i, value := range new.Values {
		if slices.Contains(old.Values, value) {
			continue
		}
		position := ""
		if i > 0 {
			position = fmt.Sprintf(" AFTER '%s'", escapeSingleQuote(new.Values[i-1]))
		} else if len(common) > 0 {
			position = fmt.Sprintf(" BEFORE '%s'", escapeSingleQuote(common[0]))
		}
		d.add(phaseCreateTypes, MigrationOp{Type: OpAlterType, SchemaName: schema, ObjectName: new.Name,
			SQL: fmt.Sprintf("ALTER TYPE %s ADD VALUE '%s'%s", name, escapeSingleQuote(value), position)})
	}
	return nil
}

func (d *differ) diffSequences() {
	d.droppedSequences = make(map[string]bool)
	for _, schema := range d.from.schemas {
		for _, sequence := range schema.Sequences {
			if findSequence(d.to.byName[schema.Name], sequence.Name) == nil {
				d.droppedSequences[schema.Name+"."+sequence.Name] = true
				d.add(phaseDropSequences, MigrationOp{Type: OpDropSequence, SchemaName: schema.Name, ObjectName: sequence.Name,
					SQL: fmt.Sprintf("DROP SEQUENCE %s.%s", quoteIdent(schema.Name), quoteIdent(sequence.Name))})
			}
		}
	}
	for _, schema := range d.to.schemas {
		for _, sequence := range schema.Sequences {
			old := findSequence(d.from.byName[schema.Name], sequence.Name)
			if old == nil {
				var buf strings.Builder
				writeCreateSequence(&buf, schema.Name, sequence)
				d.add(phaseCreateSequences, MigrationOp{Type: OpCreateSequence, SchemaName: schema.Name, ObjectName: sequence.Name, SQL: statement(buf.String())})
				continue
			}
			var clauses []string
			for _, c := range [][3]string{
				{"AS", old.DataType, sequence.DataType},
				{"INCREMENT BY", old.Increment, sequence.Increment},
				{"MINVALUE", old.MinValue, sequence.MinValue},
				{"MAXVALUE", old.MaxValue, sequence.MaxValue},
				{"START WITH", old.Start, sequence.Start},
				{"CACHE", old.CacheSize, sequence.CacheSize},
			} {
				if c[1] != c[2] && c[2] != "" {
					clauses = append(clauses, c[0]+" "+c[2])
				}
			}
			if old.Cycle != sequence.Cycle {
				if sequence.Cycle {
					clauses = append(clauses, "CYCLE")
				} else {
					clauses = append(clauses, "NO CYCLE")
				}
			}
			if len(clauses) > 0 {
				d.add(phaseCreateSequences, MigrationOp{Type: OpAlterSequence, SchemaName: schema.Name, ObjectName: sequence.Name,
					SQL: fmt.Sprintf("ALTER SEQUENCE %s.%s %s", quoteIdent(schema.Name), quoteIdent(sequence.Name), strings.Join(clauses, " "))})
			}
		}
	}
}

func findSequence(schema *metadatapb.SchemaMetadata, name string) *metadatapb.SequenceMetadata {
	for _, sequence := range schema.GetSequences() {
		if sequence.Name == name {
			return sequence
		}
	}
	return nil
}

// rebuiltViews returns the source views and materialized views a change drops
// and creates again: those whose definition changes, those reading a column
// that is dropped or changes type, and those reading a view rebuilt or
// dropped.
func (d *differ) rebuiltViews() map[string]bool {
	changedColumns := make(map[string]bool) // schema\x00table\x00column
	dropped := make(map[string]bool)        // relation IDs dropped
	for _, schema := range d.from.schemas {
		for _, table := range schema.Tables {
			target := d.to.table(schema.Name, table.Name)
			if target == nil {
				dropped[getObjectID(schema.Name, table.Name)] = true
				continue
			}
			for _, column := range table.Columns {
				newColumn := findColumn(target, column.Name)
				if newColumn == nil || columnType(newColumn.Type) != columnType(column.Type) || newColumn.Collation != column.Collation {
					changedColumns[getObjectID(schema.Name, table.Name)+objectIDSeparator+column.Name] = true
				}
			}
		}
	}
	type relation struct {
		id   string
		deps []*metadatapb.DependencyColumn
		same bool
	}
	var relations []relation
	for _, schema := range d.from.schemas {
		for _, view := range schema.Views {
			target := findView(d.to.byName[schema.Name], view.Name)
			same := target != nil && relationQuery(target.Definition, d.to.database) == relationQuery(view.Definition, d.from.database) && sameColumnNames(view.Columns, target.Columns)
			relations = append(relations, relation{getObjectID(schema.Name, view.Name), view.DependencyColumns, same})
			if target == nil {
				dropped[getObjectID(schema.Name, view.Name)] = true
			}
		}
		for _, view := range schema.MaterializedViews {
			target := findMaterializedView(d.to.byName[schema.Name], view.Name)
			same := target != nil && relationQuery(target.Definition, d.to.database) == relationQuery(view.Definition, d.from.database)
			relations = append(relations, relation{getObjectID(schema.Name, view.Name), view.DependencyColumns, same})
			if target == nil {
				dropped[getObjectID(schema.Name, view.Name)] = true
			}
		}
	}
	rebuilt := make(map[string]bool)
	for changed := true; changed; {
		changed = false
		for _, r := range relations {
			if rebuilt[r.id] || dropped[r.id] {
				continue
			}
			need := !r.same
			for _, dep := range r.deps {
				relationID := getObjectID(dep.Schema, dep.Table)
				if dropped[relationID] || rebuilt[relationID] || changedColumns[relationID+objectIDSeparator+dep.Column] {
					need = true
				}
			}
			if need {
				rebuilt[r.id] = true
				changed = true
			}
		}
	}
	return rebuilt
}

// sameColumnNames reports whether two views name the same output columns,
// where both state them.
func sameColumnNames(a, b []*metadatapb.ColumnMetadata) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}
	return slices.EqualFunc(a, b, func(x, y *metadatapb.ColumnMetadata) bool { return x.Name == y.Name })
}

func findColumn(table *metadatapb.TableMetadata, name string) *metadatapb.ColumnMetadata {
	for _, column := range table.GetColumns() {
		if column.Name == name {
			return column
		}
	}
	return nil
}

func findView(schema *metadatapb.SchemaMetadata, name string) *metadatapb.ViewMetadata {
	for _, view := range schema.GetViews() {
		if view.Name == name {
			return view
		}
	}
	return nil
}

func findMaterializedView(schema *metadatapb.SchemaMetadata, name string) *metadatapb.MaterializedViewMetadata {
	for _, view := range schema.GetMaterializedViews() {
		if view.Name == name {
			return view
		}
	}
	return nil
}

// cascadedForeignKeys returns the source foreign keys, by schema, table and
// name, that reference the key columns of a unique index the migration drops.
// The index drops with CASCADE, which takes such a key with it, so the key is
// dropped and added again like a changed one.
func (d *differ) cascadedForeignKeys() map[string]bool {
	dropped := make(map[string]bool) // schema\x00table\x00sorted key columns
	for _, schema := range d.from.schemas {
		for _, table := range schema.Tables {
			target := d.to.table(schema.Name, table.Name)
			if target == nil {
				continue
			}
			for _, index := range table.Indexes {
				if !index.Unique || index.Primary {
					continue
				}
				var newIndex *metadatapb.IndexMetadata
				for _, candidate := range target.Indexes {
					if candidate.Name == index.Name {
						newIndex = candidate
					}
				}
				if newIndex == nil || newIndex.Primary || !sameIndex(schema.Name, table.Name, newIndex, index) || indexReads(index, retypedColumns(table, target)) {
					dropped[getObjectID(schema.Name, table.Name)+objectIDSeparator+columnSet(index.Expressions)] = true
				}
			}
		}
	}
	cascaded := make(map[string]bool)
	for _, schema := range d.from.schemas {
		for _, table := range schema.Tables {
			for _, fk := range table.ForeignKeys {
				if dropped[getObjectID(fk.ReferencedSchema, fk.ReferencedTable)+objectIDSeparator+columnSet(fk.ReferencedColumns)] {
					cascaded[getObjectID(schema.Name, table.Name)+objectIDSeparator+fk.Name] = true
				}
			}
		}
	}
	return cascaded
}

func columnSet(columns []string) string {
	sorted := slices.Clone(columns)
	slices.Sort(sorted)
	return strings.Join(sorted, objectIDSeparator)
}

func (d *differ) diffTables() error {
	d.cascaded = d.cascadedForeignKeys()
	for _, schema := range d.from.schemas {
		for _, table := range schema.Tables {
			target := d.to.table(schema.Name, table.Name)
			for _, fk := range table.ForeignKeys {
				var newFK *metadatapb.ForeignKeyMetadata
				if target != nil {
					newFK = findForeignKey(target, fk.Name)
				}
				if newFK == nil || foreignKeyClause(newFK) != foreignKeyClause(fk) || d.cascaded[getObjectID(schema.Name, table.Name)+objectIDSeparator+fk.Name] {
					d.add(phaseDropForeignKeys, MigrationOp{Type: OpDropConstraint, SchemaName: schema.Name, ObjectName: fk.Name, ParentObject: table.Name,
						SQL: fmt.Sprintf("ALTER TABLE %s.%s DROP CONSTRAINT %s", quoteIdent(schema.Name), quoteIdent(table.Name), quoteIdent(fk.Name))})
				}
			}
			if target == nil {
				d.add(phaseDropTables, MigrationOp{Type: OpDropTable, SchemaName: schema.Name, ObjectName: table.Name,
					SQL: fmt.Sprintf("DROP TABLE %s.%s", quoteIdent(schema.Name), quoteIdent(table.Name))})
			}
		}
	}
	for _, schema := range d.to.schemas {
		for _, table := range schema.Tables {
			old := d.from.table(schema.Name, table.Name)
			if old == nil {
				// New tables are created with the other new objects.
				for _, fk := range table.ForeignKeys {
					d.addForeignKey(schema.Name, table.Name, fk)
				}
				continue
			}
			if err := d.alterTable(schema.Name, old, table); err != nil {
				return err
			}
		}
	}
	return nil
}

func findForeignKey(table *metadatapb.TableMetadata, name string) *metadatapb.ForeignKeyMetadata {
	for _, fk := range table.ForeignKeys {
		if fk.Name == name {
			return fk
		}
	}
	return nil
}

func foreignKeyClause(fk *metadatapb.ForeignKeyMetadata) string {
	var buf strings.Builder
	writeForeignKey(&buf, "", "", fk)
	return buf.String()
}

func (d *differ) addForeignKey(schema, table string, fk *metadatapb.ForeignKeyMetadata) {
	var buf strings.Builder
	writeForeignKey(&buf, schema, table, fk)
	d.add(phaseAddForeignKeys, MigrationOp{Type: OpAddConstraint, SchemaName: schema, ObjectName: fk.Name, ParentObject: table, SQL: statement(buf.String())})
	if fk.Comment != "" {
		d.constraintComment(schema, table, fk.Name, fk.Comment)
	}
}

// alterTable changes a kept table's columns, key, checks, indexes, foreign
// keys and comments from old to new.
func (d *differ) alterTable(schema string, old, new *metadatapb.TableMetadata) error {
	name := fmt.Sprintf("%s.%s", quoteIdent(schema), quoteIdent(new.Name))
	oldColumns, oldPrimary, oldChecks := tableShape(old)
	newColumns, newPrimary, newChecks := tableShape(new)
	alter := func(phase int, typ MigrationOpType, object, sql string) {
		d.add(phase, MigrationOp{Type: typ, SchemaName: schema, ObjectName: object, ParentObject: new.Name, SQL: "ALTER TABLE " + name + " " + sql})
	}

	// Checks and indexes that go or change are dropped before the columns
	// they read.
	for _, check := range oldChecks {
		if c := findCheck(newChecks, check.Name); c == nil || c.Expression != check.Expression {
			alter(phaseDropIndexes, OpDropConstraint, check.Name, "DROP CONSTRAINT "+quoteIdent(check.Name))
		}
	}
	d.diffIndexes(schema, new.Name, old.Indexes, new.Indexes, retypedColumns(old, new))

	oldKey, newKey := "", ""
	if oldPrimary != nil {
		oldKey = primaryKeyClause(oldPrimary)
	}
	if newPrimary != nil {
		newKey = primaryKeyClause(newPrimary)
	}
	// A column of a primary key that changes is dropped once the new key is
	// in place.
	var afterKey []string
	for _, column := range oldColumns {
		if findColumnIn(newColumns, column.Name) != nil {
			continue
		}
		if oldKey != newKey && slices.Contains(oldPrimary.GetExpressions(), column.Name) {
			afterKey = append(afterKey, column.Name)
			continue
		}
		alter(phaseDropColumns, OpDropColumn, column.Name, "DROP COLUMN "+quoteIdent(column.Name))
	}
	for _, column := range newColumns {
		oldColumn := findColumnIn(oldColumns, column.Name)
		// A default that calls a routine the target adds waits for it.
		defaultPhase := phaseAlterTables
		for _, name := range calledFunctions(column.Default) {
			if d.newRoutines[name] {
				defaultPhase = phaseLateColumns
			}
		}
		if oldColumn == nil {
			var buf strings.Builder
			d.r.writeColumnDefinition(&buf, column)
			alter(defaultPhase, OpAddColumn, column.Name, "ADD COLUMN "+buf.String())
			continue
		}
		col := "ALTER COLUMN " + quoteIdent(column.Name) + " "
		// A default calling a routine the migration drops is detached first.
		detached := false
		for _, name := range calledFunctions(oldColumn.Default) {
			detached = detached || d.droppedRoutines[name]
		}
		for _, name := range usedSequences(oldColumn.Default) {
			detached = detached || d.droppedSequences[name]
		}
		if detached {
			alter(phaseDetachDefaults, OpAlterColumn, column.Name, col+"DROP DEFAULT")
		}
		if columnType(oldColumn.Type) != columnType(column.Type) || oldColumn.Collation != column.Collation {
			typ := d.r.typeName(column.Type)
			if column.Collation != "" {
				typ += " COLLATE " + quoteIdent(column.Collation)
			}
			alter(phaseAlterTables, OpAlterColumn, column.Name, col+"SET DATA TYPE "+typ)
		}
		if generationText(oldColumn) != generationText(column) {
			return fmt.Errorf("column %s.%s.%s changes its computed expression, which CockroachDB cannot alter in place", schema, new.Name, column.Name)
		}
		if detached && column.Default != "" {
			alter(defaultPhase, OpAlterColumn, column.Name, col+"SET DEFAULT "+defaultExpression(column.Default))
		} else if !detached && oldColumn.Default != column.Default {
			if column.Default == "" {
				alter(phaseAlterTables, OpAlterColumn, column.Name, col+"DROP DEFAULT")
			} else {
				alter(defaultPhase, OpAlterColumn, column.Name, col+"SET DEFAULT "+defaultExpression(column.Default))
			}
		}
		if oldColumn.Nullable != column.Nullable {
			if column.Nullable {
				alter(phaseAlterTables, OpAlterColumn, column.Name, col+"DROP NOT NULL")
			} else {
				alter(phaseAlterTables, OpAlterColumn, column.Name, col+"SET NOT NULL")
			}
		}
		if oldColumn.IsInvisible != column.IsInvisible {
			if column.IsInvisible {
				alter(phaseAlterTables, OpAlterColumn, column.Name, col+"SET NOT VISIBLE")
			} else {
				alter(phaseAlterTables, OpAlterColumn, column.Name, col+"SET VISIBLE")
			}
		}
	}

	if oldKey != newKey {
		switch {
		case newPrimary == nil:
			return fmt.Errorf("table %s.%s drops its primary key, which CockroachDB cannot do", schema, new.Name)
		case oldPrimary == nil:
			alter(phaseAlterTables, OpAddConstraint, newPrimary.Name, "ADD "+newKey)
		default:
			// Dropping and adding the key in one statement keeps CockroachDB
			// from leaving the old key behind as a unique index.
			alter(phaseAlterTables, OpAddConstraint, newPrimary.Name, "DROP CONSTRAINT "+quoteIdent(oldPrimary.Name)+", ADD "+newKey)
		}
	}
	for _, name := range afterKey {
		alter(phaseAlterTables, OpDropColumn, name, "DROP COLUMN "+quoteIdent(name))
	}

	for _, check := range newChecks {
		if c := findCheck(oldChecks, check.Name); c == nil || c.Expression != check.Expression {
			alter(phaseCreateIndexes, OpAddConstraint, check.Name, fmt.Sprintf("ADD CONSTRAINT %s CHECK %s", quoteIdent(check.Name), check.Expression))
		}
	}
	for _, fk := range new.ForeignKeys {
		if oldFK := findForeignKey(old, fk.Name); oldFK == nil || foreignKeyClause(oldFK) != foreignKeyClause(fk) || d.cascaded[getObjectID(schema, new.Name)+objectIDSeparator+fk.Name] {
			d.addForeignKey(schema, new.Name, fk)
		} else if oldFK.Comment != fk.Comment {
			d.constraintComment(schema, new.Name, fk.Name, fk.Comment)
		}
	}

	if old.Comment != new.Comment {
		d.comment(schema, new.Name, "", "TABLE "+name, new.Comment)
	}
	for _, column := range newColumns {
		oldComment := ""
		if oldColumn := findColumnIn(oldColumns, column.Name); oldColumn != nil {
			oldComment = oldColumn.Comment
		}
		if oldComment != column.Comment {
			d.comment(schema, column.Name, new.Name, "COLUMN "+name+"."+quoteIdent(column.Name), column.Comment)
		}
	}
	for _, check := range newChecks {
		oldComment := ""
		if c := findCheck(oldChecks, check.Name); c != nil && c.Expression == check.Expression {
			oldComment = c.Comment
		}
		if oldComment != check.Comment {
			d.constraintComment(schema, new.Name, check.Name, check.Comment)
		}
	}
	return nil
}

func generationText(column *metadatapb.ColumnMetadata) string {
	if column.Generation == nil {
		return ""
	}
	return column.Generation.Type.String() + ":" + column.Generation.Expression
}

func findColumnIn(columns []*metadatapb.ColumnMetadata, name string) *metadatapb.ColumnMetadata {
	for _, column := range columns {
		if column.Name == name {
			return column
		}
	}
	return nil
}

func findCheck(checks []*metadatapb.CheckConstraintMetadata, name string) *metadatapb.CheckConstraintMetadata {
	for _, check := range checks {
		if check.Name == name {
			return check
		}
	}
	return nil
}

// retypedColumns returns the columns of a kept table whose type changes.
// CockroachDB cannot rewrite a column an index reads, so such an index is
// dropped before the change and created again after it.
func retypedColumns(old, new *metadatapb.TableMetadata) map[string]bool {
	retyped := make(map[string]bool)
	for _, column := range old.Columns {
		newColumn := findColumn(new, column.Name)
		if newColumn != nil && (columnType(newColumn.Type) != columnType(column.Type) || newColumn.Collation != column.Collation) {
			retyped[column.Name] = true
		}
	}
	return retyped
}

// indexReads reports whether an index reads any of columns: as a key, a
// stored column, or in an expression or predicate. A word of its definition
// that names a column counts, which may rebuild an index needlessly but never
// leaves one in a column's way.
func indexReads(index *metadatapb.IndexMetadata, columns map[string]bool) bool {
	if len(columns) == 0 || index == nil {
		return false
	}
	for _, expression := range index.Expressions {
		if columns[expression] {
			return true
		}
	}
	tokens, ok := crparser.Tokens(index.Definition)
	if !ok {
		return false
	}
	for _, token := range tokens {
		if columns[token.Str] {
			return true
		}
	}
	return false
}

// sameIndex reports whether two indexes of schema.table are the same index.
func sameIndex(schema, table string, a, b *metadatapb.IndexMetadata) bool {
	return indexStatement(schema, table, a) == indexStatement(schema, table, b)
}

// diffIndexes drops the secondary indexes of a kept table or materialized view
// that go or change, creates those that come or change, and sets index
// comments. A primary key changes through ALTER TABLE.
func (d *differ) diffIndexes(schema, table string, oldIndexes, newIndexes []*metadatapb.IndexMetadata, retyped map[string]bool) {
	find := func(indexes []*metadatapb.IndexMetadata, name string) *metadatapb.IndexMetadata {
		for _, index := range indexes {
			if index.Name == name {
				return index
			}
		}
		return nil
	}
	for _, index := range oldIndexes {
		if index.Primary {
			continue
		}
		newIndex := find(newIndexes, index.Name)
		if newIndex == nil || newIndex.Primary || !sameIndex(schema, table, newIndex, index) || indexReads(index, retyped) {
			sql := fmt.Sprintf("DROP INDEX %s.%s@%s", quoteIdent(schema), quoteIdent(table), quoteIdent(index.Name))
			if index.Unique {
				// A unique index may back a UNIQUE constraint, which
				// CockroachDB drops only with CASCADE.
				sql += " CASCADE"
			}
			d.add(phaseDropIndexes, MigrationOp{Type: OpDropIndex, SchemaName: schema, ObjectName: index.Name, ParentObject: table, SQL: sql})
		}
	}
	for _, index := range newIndexes {
		oldIndex := find(oldIndexes, index.Name)
		if !index.Primary && (oldIndex == nil || oldIndex.Primary || !sameIndex(schema, table, oldIndex, index) || indexReads(oldIndex, retyped)) {
			d.add(phaseCreateIndexes, MigrationOp{Type: OpCreateIndex, SchemaName: schema, ObjectName: index.Name, ParentObject: table, SQL: indexStatement(schema, table, index)})
			oldIndex = nil
		}
		oldComment := ""
		if oldIndex != nil {
			oldComment = oldIndex.Comment
		}
		if oldComment != index.Comment {
			d.comment(schema, index.Name, table, fmt.Sprintf("INDEX %s.%s@%s", quoteIdent(schema), quoteIdent(table), quoteIdent(index.Name)), index.Comment)
		}
	}
}

// diffViews drops the views and materialized views that go or are rebuilt,
// dependents first, and alters the indexes of materialized views kept.
func (d *differ) diffViews() {
	var drops []MigrationOp
	for _, object := range orderObjects(d.from.schemas) {
		switch {
		case object.view != nil:
			id := getObjectID(object.schema, object.view.Name)
			if d.rebuilt[id] || findView(d.to.byName[object.schema], object.view.Name) == nil {
				drops = append(drops, MigrationOp{Type: OpDropView, SchemaName: object.schema, ObjectName: object.view.Name,
					SQL: fmt.Sprintf("DROP VIEW %s.%s", quoteIdent(object.schema), quoteIdent(object.view.Name))})
			}
		case object.materializedView != nil:
			id := getObjectID(object.schema, object.materializedView.Name)
			target := findMaterializedView(d.to.byName[object.schema], object.materializedView.Name)
			if d.rebuilt[id] || target == nil {
				drops = append(drops, MigrationOp{Type: OpDropView, SchemaName: object.schema, ObjectName: object.materializedView.Name,
					SQL: fmt.Sprintf("DROP MATERIALIZED VIEW %s.%s", quoteIdent(object.schema), quoteIdent(object.materializedView.Name))})
				continue
			}
			d.diffIndexes(object.schema, target.Name, materializedViewIndexes(object.materializedView), materializedViewIndexes(target), nil)
		}
	}
	slices.Reverse(drops)
	d.phases[phaseDropViews] = append(d.phases[phaseDropViews], drops...)
}

// planRoutines decides which source routines the migration drops: those the
// target removes, and those it drops and creates again. A routine whose
// definition changes, or whose body reads a relation the migration drops or
// rebuilds, is recreated, so the relation can go while nothing reads it. A
// routine a kept table's default or check still calls cannot be dropped, so it
// is replaced in place instead.
func (d *differ) planRoutines() {
	d.droppedRoutines = make(map[string]bool)
	d.recreated = make(map[string]bool)
	pinned := make(map[string]bool)
	for _, schema := range d.to.schemas {
		for _, table := range schema.Tables {
			if d.from.table(schema.Name, table.Name) == nil {
				continue
			}
			for _, name := range tableCalls(table) {
				pinned[name] = true
			}
		}
	}
	gone := make(map[string]bool)
	for id := range d.rebuilt {
		gone[id] = true
	}
	for _, schema := range d.from.schemas {
		target := d.to.byName[schema.Name]
		for _, table := range schema.Tables {
			if d.to.table(schema.Name, table.Name) == nil {
				gone[getObjectID(schema.Name, table.Name)] = true
			}
		}
		for _, view := range schema.Views {
			if findView(target, view.Name) == nil {
				gone[getObjectID(schema.Name, view.Name)] = true
			}
		}
		for _, view := range schema.MaterializedViews {
			if findMaterializedView(target, view.Name) == nil {
				gone[getObjectID(schema.Name, view.Name)] = true
			}
		}
	}
	for _, schema := range d.from.schemas {
		for _, function := range schema.Functions {
			names := []string{function.Name, schema.Name + "." + function.Name}
			newFunction := findFunction(d.to.byName[schema.Name], metadataFuncIdentity(function))
			if newFunction == nil {
				d.droppedRoutines[names[0]], d.droppedRoutines[names[1]] = true, true
				continue
			}
			if pinned[names[0]] || pinned[names[1]] {
				continue
			}
			recreate := routineDefinition(function.Definition, d.from.database) != routineDefinition(newFunction.Definition, d.to.database)
			for _, relation := range routineRelations(function.Definition) {
				recreate = recreate || gone[relation]
			}
			if recreate {
				d.recreated[getObjectID(schema.Name, metadataFuncIdentity(function))] = true
				d.droppedRoutines[names[0]], d.droppedRoutines[names[1]] = true, true
			}
		}
	}
}

// tableCalls returns the routines a table's defaults and checks call.
func tableCalls(table *metadatapb.TableMetadata) []string {
	var names []string
	for _, column := range table.Columns {
		names = append(names, calledFunctions(column.Default)...)
	}
	for _, check := range table.CheckConstraints {
		names = append(names, calledFunctions(check.Expression)...)
	}
	return names
}

// diffRoutines drops the routines planRoutines chose. A routine that a
// dropped table's or column's default, or any check, calls is dropped after
// the tables, columns and checks go; the others first, before the relations
// their bodies read. A kept column's default is detached from a dropped
// routine before either.
func (d *differ) diffRoutines() {
	late := make(map[string]bool)
	for _, schema := range d.from.schemas {
		for _, table := range schema.Tables {
			target := d.to.table(schema.Name, table.Name)
			for _, column := range table.Columns {
				if target != nil && findColumn(target, column.Name) != nil {
					continue
				}
				for _, name := range calledFunctions(column.Default) {
					late[name] = true
				}
			}
			for _, check := range table.CheckConstraints {
				for _, name := range calledFunctions(check.Expression) {
					late[name] = true
				}
			}
		}
	}
	for _, schema := range d.from.schemas {
		for _, function := range schema.Functions {
			recreated := d.recreated[getObjectID(schema.Name, metadataFuncIdentity(function))]
			if !recreated && findFunction(d.to.byName[schema.Name], metadataFuncIdentity(function)) != nil {
				continue
			}
			kind := "FUNCTION"
			if isProcedure(function.Definition) {
				kind = "PROCEDURE"
			}
			phase := phaseDropRoutinesFirst
			if late[function.Name] || late[schema.Name+"."+function.Name] {
				phase = phaseDropRoutines
			}
			typ := OpDropFunction
			if recreated {
				typ = OpAlterFunction
			}
			d.add(phase, MigrationOp{Type: typ, SchemaName: schema.Name, ObjectName: function.Name,
				SQL: fmt.Sprintf("DROP %s %s.%s", kind, quoteIdent(schema.Name), routineSignature(function))})
		}
	}
}

// routineRelations returns the relations a routine's body names, by object ID.
func routineRelations(definition string) []string {
	stmt, err := parseOne(strings.TrimSuffix(strings.TrimSpace(definition), ";"))
	if err != nil {
		return nil
	}
	routine, ok := stmt.(*tree.CreateRoutine)
	if !ok {
		return nil
	}
	var relations []string
	ctx := tree.NewFmtCtx(tree.FmtSimple, tree.FmtReformatTableNames(func(ctx *tree.FmtCtx, tn *tree.TableName) {
		relations = append(relations, getObjectID(schemaOf(tn), string(tn.ObjectName)))
		ctx.WithReformatTableNames(nil, func() { ctx.FormatNode(tn) })
	}))
	defer ctx.Close()
	language := ""
	for _, option := range routine.Options {
		if l, ok := option.(tree.RoutineLanguage); ok {
			language = strings.ToLower(string(l))
		}
	}
	for _, option := range routine.Options {
		body, ok := option.(tree.RoutineBodyStr)
		if !ok {
			continue
		}
		switch language {
		case "sql":
			if stmts, err := crparser.Parse(string(body)); err == nil {
				for _, s := range stmts {
					ctx.FormatNode(s.AST)
				}
			}
		case "plpgsql":
			if block, err := plpgsqlparser.Parse(string(body)); err == nil {
				ctx.FormatNode(block.AST)
			}
		}
	}
	return relations
}

// schemaOf returns a table name's schema, public when it names none.
func schemaOf(name *tree.TableName) string {
	if name.ExplicitSchema {
		return string(name.SchemaName)
	}
	return "public"
}

// routineSignature is a routine's synced signature with its name quoted.
func routineSignature(function *metadatapb.FunctionMetadata) string {
	if rest, ok := strings.CutPrefix(function.Signature, function.Name); ok && strings.HasPrefix(rest, "(") {
		return quoteIdent(function.Name) + rest
	}
	return quoteIdent(function.Name) + "()"
}

func findFunction(schema *metadatapb.SchemaMetadata, identity string) *metadatapb.FunctionMetadata {
	for _, function := range schema.GetFunctions() {
		if metadataFuncIdentity(function) == identity {
			return function
		}
	}
	return nil
}

// createObjects creates the new tables, the new and rebuilt views and
// materialized views, and the new routines, each after what it depends on,
// and replaces the routines whose definition changes.
func (d *differ) createObjects() {
	r := d.r
	var schemas []*metadatapb.SchemaMetadata
	for _, schema := range d.to.schemas {
		old := d.from.byName[schema.Name]
		subset := &metadatapb.SchemaMetadata{Name: schema.Name}
		for _, table := range schema.Tables {
			if d.from.table(schema.Name, table.Name) == nil {
				subset.Tables = append(subset.Tables, table)
			}
		}
		for _, view := range schema.Views {
			if findView(old, view.Name) == nil || d.rebuilt[getObjectID(schema.Name, view.Name)] {
				subset.Views = append(subset.Views, view)
			}
		}
		for _, view := range schema.MaterializedViews {
			if findMaterializedView(old, view.Name) == nil || d.rebuilt[getObjectID(schema.Name, view.Name)] {
				subset.MaterializedViews = append(subset.MaterializedViews, view)
			}
		}
		for _, function := range schema.Functions {
			oldFunction := findFunction(old, metadataFuncIdentity(function))
			if oldFunction == nil || d.recreated[getObjectID(schema.Name, metadataFuncIdentity(function))] || routineDefinition(oldFunction.Definition, d.from.database) != routineDefinition(function.Definition, d.to.database) {
				subset.Functions = append(subset.Functions, function)
			}
		}
		schemas = append(schemas, subset)
	}
	for _, object := range orderObjects(schemas) {
		var buf strings.Builder
		switch {
		case object.function != nil:
			typ := OpCreateFunction
			definition := r.routine(object.function)
			switch {
			case d.recreated[getObjectID(object.schema, metadataFuncIdentity(object.function))]:
				typ = OpAlterFunction
			case findFunction(d.from.byName[object.schema], metadataFuncIdentity(object.function)) != nil:
				typ = OpAlterFunction
				definition = replaceRoutine(definition)
			}
			d.add(phaseCreateObjects, MigrationOp{Type: typ, SchemaName: object.schema, ObjectName: object.function.Name, SQL: definition})
		case object.table != nil:
			r.writeTable(&buf, object.schema, object.table)
			d.addScript(phaseCreateObjects, OpCreateTable, object.schema, object.table.Name, buf.String())
		case object.view != nil:
			r.writeView(&buf, object.schema, object.view)
			d.addScript(phaseCreateObjects, OpCreateView, object.schema, object.view.Name, buf.String())
		case object.materializedView != nil:
			r.writeMaterializedView(&buf, object.schema, object.materializedView)
			d.addScript(phaseCreateObjects, OpCreateView, object.schema, object.materializedView.Name, buf.String())
		}
	}
}

func (r renderer) routine(function *metadatapb.FunctionMetadata) string {
	return routineDefinition(function.Definition, r.database)
}

// replaceRoutine turns a CREATE FUNCTION or CREATE PROCEDURE into CREATE OR
// REPLACE, on CockroachDB's tree.
func replaceRoutine(definition string) string {
	stmt, err := parseOne(definition)
	if err != nil {
		return definition
	}
	routine, ok := stmt.(*tree.CreateRoutine)
	if !ok {
		return definition
	}
	routine.Replace = true
	return tree.AsString(routine)
}

// addScript adds the statements a definition writer wrote, one op each, split
// by CockroachDB's parser so a quoted semicolon stays inside its statement.
// The first is the object's own; the rest (comments, indexes) are typed by
// their leading keywords.
func (d *differ) addScript(phase int, typ MigrationOpType, schema, object, script string) {
	stmts, err := crparser.Parse(script)
	if err != nil {
		// The renderer wrote the script; keep it whole rather than guess.
		d.add(phase, MigrationOp{Type: typ, SchemaName: schema, ObjectName: object, SQL: statement(script)})
		return
	}
	for i, parsed := range stmts {
		stmt := statement(parsed.SQL)
		if stmt == "" {
			continue
		}
		op := MigrationOp{Type: typ, SchemaName: schema, ObjectName: object, SQL: stmt}
		if i > 0 {
			switch {
			case strings.HasPrefix(stmt, "COMMENT ON"):
				op.Type = OpComment
				d.add(phaseComments, op)
				continue
			case strings.HasPrefix(stmt, "CREATE"):
				op.Type = OpCreateIndex
				op.ParentObject = object
			default:
				op.Type = OpAlterTable
			}
		}
		d.add(phase, op)
	}
}

func (d *differ) comment(schema, object, parent, target, comment string) {
	text := "NULL"
	if comment != "" {
		text = "'" + escapeSingleQuote(comment) + "'"
	}
	d.add(phaseComments, MigrationOp{Type: OpComment, SchemaName: schema, ObjectName: object, ParentObject: parent,
		SQL: fmt.Sprintf("COMMENT ON %s IS %s", target, text)})
}

func (d *differ) constraintComment(schema, table, name, comment string) {
	d.comment(schema, name, table, fmt.Sprintf("CONSTRAINT %s ON %s.%s", quoteIdent(name), quoteIdent(schema), quoteIdent(table)), comment)
}

// statement trims a rendered statement of its terminator and spacing.
func statement(s string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), ";"))
}
