package review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// statement lowers one statement (spec §5), false when the walk does not
// follow it, which ends the walk before it (spec §6).
func (b *builder) statement(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.CreateTableStmt:
		return b.createTable(v)
	case *ast.AlterTableStmt:
		return b.alterTable(v)
	case *ast.DropStmt:
		return b.drop(v)
	case *ast.TruncateStmt:
		return b.truncate(v)
	case *ast.CreateIndexStmt:
		return b.createIndex(v)
	case *ast.CreateVectorIndexStmt:
		return b.createVectorIndex(v)
	case *ast.AlterIndexStmt:
		return b.alterIndex(v)
	case *ast.CreateViewStmt:
		return b.createView(v)
	case *ast.AlterViewStmt:
		return b.alterView(v)
	case *ast.AlterMaterializedViewStmt:
		return b.alterMaterializedView(v)
	case *ast.CreateSequenceStmt:
		return b.createObject(v.Name, kindSequence, v.IfNotExists, false)
	case *ast.AlterSequenceStmt:
		return b.alterSequence(v)
	case *ast.CreateSynonymStmt:
		if v.Public {
			return b.publicSynonymNamed(v.Name, walk.Present)
		}
		return b.createObject(v.Name, kindSynonym, v.IfNotExists, v.OrReplace)
	case *ast.AlterSynonymStmt:
		if v.Public {
			return true
		}
		if !b.alterObject(v.Name, kindSynonym, v.IfExists, missingSynonym, "private synonym") {
			return false
		}
		b.mayRestoreSynonyms()
		return true
	case *ast.CreateTypeStmt:
		if v.IsBody {
			// A body has a namespace of its own.
			return true
		}
		return b.createObject(v.Name, kindType, v.IfNotExists, v.OrReplace)
	case *ast.AlterTypeStmt:
		if v.CompileTarget == "BODY" {
			return true
		}
		if !b.alterObject(v.Name, kindType, v.IfExists, func(w string) string { return missingKind(kindType, w) }, kindWord(kindType)) {
			return false
		}
		if v.Action != "" && v.Action != "COMPILE" && v.Action != "EDITIONABLE" && v.Action != "NONEDITIONABLE" {
			// Structural type evolution may affect dependent table columns/constraints.
			b.Last = true
		}
		return true
	case *ast.CreatePackageStmt:
		if v.IsBody {
			return true
		}
		return b.createObject(v.Name, kindPackage, v.IfNotExists, v.OrReplace)
	case *ast.AlterPackageStmt:
		if v.CompileTarget == "BODY" {
			return true
		}
		return b.alterObject(v.Name, kindPackage, v.IfExists, func(w string) string { return missingKind(kindPackage, w) }, kindWord(kindPackage))
	case *ast.CreateProcedureStmt:
		return b.createObject(v.Name, kindProcedure, v.IfNotExists, v.OrReplace)
	case *ast.AlterProcedureStmt:
		return b.alterObject(v.Name, kindProcedure, v.IfExists, func(w string) string { return missingKind(kindProcedure, w) }, kindWord(kindProcedure))
	case *ast.CreateFunctionStmt:
		return b.createObject(v.Name, kindFunction, v.IfNotExists, v.OrReplace)
	case *ast.AlterFunctionStmt:
		return b.alterObject(v.Name, kindFunction, v.IfExists, func(w string) string { return missingKind(kindFunction, w) }, kindWord(kindFunction))
	case *ast.CreateTriggerStmt:
		return b.createTrigger(v)
	case *ast.AlterTriggerStmt:
		return b.alterTrigger(v)
	case *ast.RenameStmt:
		return b.rename(v)
	case *ast.AlterSessionStmt:
		return b.alterSession(v)
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
		return b.dml(n)
	case *ast.SelectStmt:
		b.Last = b.callsUserCode(n)
		return true
	case *ast.GrantStmt:
		b.mayRestoreSynonyms()
		return true
	case *ast.CommentStmt:
		return b.comment(v)
	case *ast.RevokeStmt, *ast.LockTableStmt, *ast.AnalyzeStmt,
		*ast.ExplainPlanStmt, *ast.PurgeStmt, *ast.AuditStmt, *ast.NoauditStmt, *ast.CommitStmt,
		*ast.RollbackStmt, *ast.SavepointStmt, *ast.SetTransactionStmt, *ast.SetConstraintsStmt:
		// None changes a name the walk keeps: DDL commits as it runs, so a
		// ROLLBACK undoes none, and PURGE empties only the recycle bin.
		return true
	}
	return false
}

// createObject follows CREATE of a sequence, a private synonym, a type, a
// package, a procedure, or a function: the name must be free in the
// object namespace (00955), and OR REPLACE replaces only an object of its
// own kind (00955 for another, verified). A PL/SQL unit with compilation
// errors is created all the same (spec §4).
func (b *builder) createObject(name *ast.ObjectName, kind string, ifNotExists, orReplace bool) bool {
	r, ok := b.ref(name)
	if !ok {
		return false
	}
	if b.refTooLong(name, r) {
		return true
	}
	b.Order = walk.Lookups
	if b.claimName(r, ifNotExists, orReplace, kind) {
		b.arriveObject(r, kind)
	}
	return true
}

// claimName adds the checks of a CREATE on its name in the object
// namespace, and reports whether the statement creates anything: CREATE
// ... IF NOT EXISTS does nothing when an object of any kind holds the
// name, and may do nothing when one may (spec §4), which makes the
// creation an optional part.
func (b *builder) claimName(r objRef, ifNotExists, orReplace bool, kind string) bool {
	k := r.key(spaceObject)
	f, got := b.look(k)
	switch {
	case ifNotExists:
		switch f {
		case walk.Present:
			return false
		case walk.Unknown:
			b.Next()
			b.Optional()
			// Either way, an object holds the name after it.
			b.settle(k, walk.Present)
		}
	case orReplace:
		if f == walk.Present && got != "" && got != kind {
			b.claim(k, r, kind)
		}
	default:
		b.claim(k, r, kind)
	}
	return true
}

// claim adds the check that the name is free in the object namespace
// (00955), whose finding says the object exists when one of kind holds
// the name, and that the name is taken when another does.
func (b *builder) claim(k walk.Key, r objRef, kind string) {
	b.needNote(k, walk.Absent, nil, &walkNote{rng: r.rng, present: takenName(r.written),
		absentCode: review.ObjectNotFound, presentCode: review.NameTaken, written: r.written, kind: kind})
}

// heldByOther reports whether the walk knows that k, the name's key in
// the space of kind, is absent while an object of another kind holds the
// name in the object namespace, and then adds the check that fails
// there. Oracle refuses such a name as missing (00942, 02289), as it
// refuses one nothing holds; the finding says what the name is instead.
func (b *builder) heldByOther(k walk.Key, r objRef, kind string) bool {
	if f, _ := b.look(k); f != walk.Absent {
		return false
	}
	ok := r.key(spaceObject)
	if f, got := b.look(ok); f != walk.Present || got == "" || got == kind {
		return false
	}
	b.needNote(ok, walk.Present, []string{kind}, &walkNote{rng: r.rng, present: notKind(r.written, kindWord(kind)),
		absentCode: review.ObjectNotFound, presentCode: review.WrongObjectKind, written: r.written, kind: kind})
	return true
}

// constraintNames returns the names a column definition gives its inline
// constraints. The parser keeps the name of a CHECK, key, or reference,
// and not of a NOT NULL or NULL, so the names are read from the tokens.
func (b *builder) constraintNames(c *ast.ColumnDef) []string {
	var out []string
	toks := tokens(b.w.sql, c.Loc)
	for i := 0; i+1 < len(toks); i++ {
		if toks[i].Type != tokenConstraint {
			continue
		}
		if name, ok := tokenName(toks[i+1]); ok {
			out = append(out, name)
		}
	}
	return out
}

// columnKeys returns what a column definition's inline constraints
// declare: whether one is a primary key, and whether one is a key, whose
// index takes a name.
func columnKeys(c *ast.ColumnDef) (primary, key bool) {
	if c.Constraints == nil {
		return false, false
	}
	for _, item := range c.Constraints.Items {
		cc, ok := item.(*ast.ColumnConstraint)
		if !ok {
			continue
		}
		switch cc.Type {
		case ast.CONSTRAINT_PRIMARY:
			primary, key = true, true
		case ast.CONSTRAINT_UNIQUE:
			key = true
		}
	}
	return primary, key
}

// columnDefs lists the column definitions of a list.
func columnDefs(l *ast.List) []*ast.ColumnDef {
	var out []*ast.ColumnDef
	if l == nil {
		return nil
	}
	for _, item := range l.Items {
		if c, ok := item.(*ast.ColumnDef); ok {
			out = append(out, c)
		}
	}
	return out
}

// outOfLine lists the out-of-line constraints of a list.
func outOfLine(l *ast.List) []*ast.TableConstraint {
	var out []*ast.TableConstraint
	if l == nil {
		return nil
	}
	for _, item := range l.Items {
		if c, ok := item.(*ast.TableConstraint); ok {
			out = append(out, c)
		}
	}
	return out
}

// createsIndex reports whether a constraint's text creates an index of a
// name of its own: USING INDEX (CREATE INDEX ...).
func (b *builder) createsIndex(loc ast.Loc) bool {
	for _, tok := range tokens(b.w.sql, loc) {
		if tok.Type == tokenCreate {
			return true
		}
	}
	return false
}

// tableDecl is what a CREATE TABLE, or an ADD or MODIFY of ALTER TABLE,
// declares beside its columns: the constraint names, whether one is a
// primary key, the key names, whose indexes take those names or others,
// and whether a key creates an index of a name of its own.
type tableDecl struct {
	names     []string
	rngs      []review.Range
	primaries int
	keys      []string
	indexes   bool
}

func (b *builder) declareColumn(d *tableDecl, c *ast.ColumnDef) {
	names := b.constraintNames(c)
	for _, n := range names {
		d.names = append(d.names, n)
		d.rngs = append(d.rngs, rangeOf(c.Loc))
	}
	primary, key := columnKeys(c)
	if primary {
		d.primaries++
	}
	if key {
		d.keys = append(d.keys, names...)
		d.indexes = d.indexes || b.createsIndex(c.Loc)
	}
}

func (b *builder) declareConstraint(d *tableDecl, c *ast.TableConstraint) {
	if c.Name != "" {
		d.names = append(d.names, c.Name)
		d.rngs = append(d.rngs, rangeOf(c.Loc))
	}
	switch c.Type {
	case ast.CONSTRAINT_PRIMARY:
		d.primaries++
		fallthrough
	case ast.CONSTRAINT_UNIQUE:
		if c.Name != "" {
			d.keys = append(d.keys, c.Name)
		}
		d.indexes = d.indexes || b.createsIndex(c.Loc)
	}
}

// declTextFails fails a declaration its text alone refutes: a constraint
// name twice (02264), or two primary keys (02260) (verified).
func (b *builder) declTextFails(d *tableDecl, r objRef) bool {
	seen := make(map[string]bool)
	for i, n := range d.names {
		if seen[n] {
			b.fail(d.rngs[i], review.DuplicateName, twice("declares", "constraint", n))
			return true
		}
		seen[n] = true
	}
	if d.primaries > 1 {
		b.fail(r.rng, review.NameTaken, takenPrimaryKey(r.written))
		return true
	}
	return false
}

// declChecks adds the checks a declaration needs on a table: each
// constraint name free in the schema (02264), and with a primary key, none
// on the table (02260).
func (b *builder) declChecks(d *tableDecl, r objRef) {
	for i, n := range d.names {
		b.need(r.schemaKey(spaceConstraint, n), walk.Absent, nil, d.rngs[i], "", takenConstraint(n))
	}
	if d.primaries > 0 {
		b.need(r.childKey(spacePrimary, ""), walk.Absent, nil, r.rng, "", takenPrimaryKey(r.written))
	}
}

// declEffects records what a declaration makes: its constraints, on the
// table and in the schema, the primary key, and the index names of its
// keys, which are Unknown: a key's index takes the constraint's name, or
// another, or none (spec §3).
func (b *builder) declEffects(d *tableDecl, r objRef) {
	for _, n := range d.names {
		b.Set(r.schemaKey(spaceConstraint, n), walk.Present, r.path())
		b.Set(r.childKey(spaceTableConstraint, n), walk.Present, "")
	}
	if d.primaries > 0 {
		b.Set(r.childKey(spacePrimary, ""), walk.Present, "")
	}
	for _, n := range d.keys {
		b.Set(r.schemaKey(spaceIndex, n), walk.Unknown, "")
	}
	if d.indexes {
		b.MayAdd(walk.Scope{Space: spaceIndex, Parent: walk.Path(r.schema)})
	}
}

// createTable follows CREATE [GLOBAL TEMPORARY] TABLE [IF NOT EXISTS],
// with a column list or AS SELECT. A private temporary table, a nested
// table, an external table, or a table with a supplemental log group,
// which is a constraint, ends the walk (spec §6), as does an object or
// XMLType table, whose columns come from its type, and a table with a
// period, whose start and end columns Oracle adds.
func (b *builder) createTable(v *ast.CreateTableStmt) bool {
	if v.Private || v.OrReplace {
		return false
	}
	stopAfterCheck := v.Sharded || v.Duplicated || v.Parent != nil || v.Organization == "EXTERNAL" || v.OfType != nil
	if v.Options != nil {
		for _, item := range v.Options.Items {
			if o, ok := item.(*ast.DDLOption); ok {
				if o.Key == "SUPPLEMENTAL LOG" {
					return false
				}
				if o.Key == "NESTED" || o.Key == "ORGANIZATION EXTERNAL" || o.Key == "PERIOD FOR" {
					stopAfterCheck = true
				}
			}
		}
	}
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	if stopAfterCheck {
		if b.refTooLong(v.Name, r) {
			return true
		}
		b.Order = walk.Lookups
		b.claimName(r, v.IfNotExists, false, kindTable)
		b.Last = true
		return true
	}
	columns := columnDefs(v.Columns)
	var d tableDecl
	for _, c := range columns {
		b.declareColumn(&d, c)
	}
	for _, c := range outOfLine(v.Constraints) {
		b.declareConstraint(&d, c)
	}
	if b.refTooLong(v.Name, r) || b.namesTooLong(columns, &d) || b.columnsTwice(columns) {
		return true
	}
	// One object: every check holds against the state before the
	// statement, and the creation is a part of its own, one IF NOT EXISTS
	// may skip. A column named twice fails it even when the table exists
	// (verified).
	b.Order = walk.Lookups
	if !b.claimName(r, v.IfNotExists, false, kindTable) {
		return true
	}
	if !v.IfNotExists {
		if v.AsQuery == nil && v.OfType == nil {
			colSet := make(map[string]bool, len(columns))
			for _, c := range columns {
				if columnMissingDatatype(c) {
					b.fail(rangeOf(c.Loc), review.UnsupportedForm, "Column "+review.Code(ident(c.Name))+" must specify a datatype.")
					return true
				}
				colSet[c.Name] = true
			}
			if len(columns) > 0 {
				for _, tc := range outOfLine(v.Constraints) {
					if tc.Columns == nil {
						continue
					}
					for _, item := range tc.Columns.Items {
						if s, ok := item.(*ast.String); ok && s.Str != "" && !colSet[s.Str] && !generated(spaceColumn, s.Str) {
							b.fail(rangeOf(tc.Loc), review.ObjectNotFound, missingIn("Column", s.Str, r.written))
							return true
						}
					}
				}
			}
		} else if v.AsQuery != nil && len(columns) == 0 {
			if b.checkSelectAliases(v.AsQuery, r) {
				return true
			}
		}
	}
	if b.declTextFails(&d, r) {
		return true
	}
	for i, n := range d.names {
		b.need(r.schemaKey(spaceConstraint, n), walk.Absent, nil, d.rngs[i], "", takenConstraint(n))
	}
	if v.AsQuery != nil && b.callsUserCode(v.AsQuery) {
		b.Last = true
	}
	b.arriveObject(r, kindTable)
	path := walk.Scope{Parent: r.path(), Deep: true}
	b.Clear(path)
	if v.AsQuery != nil {
		// The query gives the columns.
		b.Forget(path)
	}
	for _, c := range columns {
		b.Set(r.childKey(spaceColumn, c.Name), walk.Present, "")
	}
	b.declEffects(&d, r)
	return true
}

// columnMissingDatatype reports whether a column definition in CREATE TABLE
// (without AS SELECT or OF type) or ALTER TABLE ADD omits its datatype where
// Oracle requires one (ORA-02263): virtual columns (AS (expr)), DOMAIN
// columns, and columns with an inline REFERENCES constraint inherit their type.
func columnMissingDatatype(c *ast.ColumnDef) bool {
	if c == nil || c.TypeName != nil || c.Virtual != nil || c.Domain != nil {
		return false
	}
	if c.Constraints != nil {
		for _, item := range c.Constraints.Items {
			if cc, ok := item.(*ast.ColumnConstraint); ok && cc.Type == ast.CONSTRAINT_FOREIGN {
				return false
			}
		}
	}
	return true
}

// namesTooLong checks the names a table declares.
func (b *builder) namesTooLong(columns []*ast.ColumnDef, d *tableDecl) bool {
	for _, c := range columns {
		if b.tooLong(c.Name, rangeOf(c.Loc)) {
			return true
		}
	}
	for i, n := range d.names {
		if b.tooLong(n, d.rngs[i]) {
			return true
		}
	}
	return false
}

// columnsTwice fails a column list that names a column twice (00957).
func (b *builder) columnsTwice(columns []*ast.ColumnDef) bool {
	seen := make(map[string]bool)
	for _, c := range columns {
		if seen[c.Name] {
			b.fail(rangeOf(c.Loc), review.DuplicateName, twice("declares", "column", c.Name))
			return true
		}
		seen[c.Name] = true
	}
	return false
}

// drop follows DROP of one object.
func (b *builder) drop(v *ast.DropStmt) bool {
	if v.Names == nil || len(v.Names.Items) != 1 {
		return false
	}
	name, ok := v.Names.Items[0].(*ast.ObjectName)
	if !ok {
		return false
	}
	switch v.ObjectType {
	case ast.OBJECT_PACKAGE_BODY, ast.OBJECT_TYPE_BODY:
		return true
	case ast.OBJECT_SYNONYM:
		if b.publicSynonym(v.Loc) {
			return b.publicSynonymNamed(name, walk.Absent)
		}
	}
	r, ok := b.ref(name)
	if !ok {
		return false
	}
	if b.refTooLong(name, r) {
		return true
	}
	b.Order = walk.Lookups
	switch v.ObjectType {
	case ast.OBJECT_TABLE:
		return b.dropTable(v, r)
	case ast.OBJECT_VIEW:
		if !v.IfExists && b.heldByOther(r.key(spaceView), r, kindView) {
			return true
		}
		if !b.dropsKind(r.key(spaceView), nil, v.IfExists, r.rng, missingView(r.written), "") {
			return true
		}
		b.leaveObject(r)
		b.dropOwned(r)
		if v.Cascade {
			b.cascadeConstraints()
		}
	case ast.OBJECT_MATERIALIZED_VIEW:
		k := r.key(spaceTable)
		if !v.IfExists {
			if b.heldByOther(k, r, kindMView) {
				return true
			}
			if f, kind := b.look(k); f == walk.Present && kind == kindTable {
				b.need(k, walk.Present, []string{kindMView}, r.rng, "", wrongKind(r.written, kindTable, kindWord(kindMView)))
				return true
			}
		}
		if !b.exists(k, []string{kindMView}, v.IfExists, r.rng, missingKind(kindMView, r.written), wrongKind(r.written, kindTable, kindWord(kindMView))) {
			return true
		}
		b.Last = true
	case ast.OBJECT_SEQUENCE:
		if !v.IfExists && b.heldByOther(r.key(spaceSequence), r, kindSequence) {
			return true
		}
		if !b.dropsKind(r.key(spaceSequence), nil, v.IfExists, r.rng, missingSequence(r.written), "") {
			return true
		}
		b.leaveObject(r)
	case ast.OBJECT_SYNONYM:
		if !b.exists(r.key(spaceObject), []string{kindSynonym}, v.IfExists, r.rng, missingSynonym(r.written), notKind(r.written, "private synonym")) {
			return true
		}
		b.leaveObject(r)
	case ast.OBJECT_PACKAGE, ast.OBJECT_PROCEDURE, ast.OBJECT_FUNCTION, ast.OBJECT_TYPE:
		kind := map[ast.ObjectType]string{
			ast.OBJECT_PACKAGE:   kindPackage,
			ast.OBJECT_PROCEDURE: kindProcedure,
			ast.OBJECT_FUNCTION:  kindFunction,
			ast.OBJECT_TYPE:      kindType,
		}[v.ObjectType]
		if !b.exists(r.key(spaceObject), []string{kind}, v.IfExists, r.rng, missingKind(kind, r.written), notKind(r.written, kindWord(kind))) {
			return true
		}
		b.leaveObject(r)
		if v.ObjectType == ast.OBJECT_TYPE && v.Force {
			b.forgetTypeDependents()
		}
	case ast.OBJECT_INDEX:
		b.Last = b.domainIndex(r)
		if !b.dropsKind(r.key(spaceIndex), nil, v.IfExists, r.rng, missingIndex(r.written), "") {
			return true
		}
		b.Set(r.key(spaceIndex), walk.Absent, "")
	case ast.OBJECT_TRIGGER:
		if !b.dropsKind(r.key(spaceTrigger), nil, v.IfExists, r.rng, missingTrigger(r.written), "") {
			return true
		}
		b.Set(r.key(spaceTrigger), walk.Absent, "")
	default:
		return false
	}
	return true
}

// publicSynonym reports whether a DROP SYNONYM drops a public one, which
// the parser does not record.
func (b *builder) publicSynonym(loc ast.Loc) bool {
	for _, tok := range tokens(b.w.sql, loc) {
		if tok.Type == tokenPublic {
			return true
		}
	}
	return false
}

// publicSynonymNamed records what CREATE or DROP PUBLIC SYNONYM leaves of
// the synonym, which the walk claims nothing about.
func (b *builder) publicSynonymNamed(name *ast.ObjectName, f walk.Fact) bool {
	if name == nil || name.Name == "" || name.Schema != "" || name.DBLink != "" {
		return false
	}
	kind := ""
	if f == walk.Present {
		b.mayRestoreSynonyms()
		kind = kindSynonym
	}
	b.Set(publicKey(name.Name), f, kind)
	return true
}

// exists adds the check that a DROP's or ALTER's object exists, of one of
// kinds when kinds is not nil, and reports whether the statement does
// anything. With IF EXISTS it checks nothing (spec §4): it does nothing on
// a missing name or one of another kind, and may do nothing on a name the
// walk does not know, which makes what it does an optional part.
func (b *builder) exists(k walk.Key, kinds []string, ifExists bool, rng review.Range, absent, present string) bool {
	if !ifExists {
		b.need(k, walk.Present, kinds, rng, absent, present)
		return true
	}
	f, kind := b.look(k)
	switch {
	case f == walk.Absent, f == walk.Present && kinds != nil && kind != "" && !contains(kinds, kind):
		return false
	case f == walk.Unknown, f == walk.Present && kind == "" && kinds != nil:
		b.Next()
		b.Optional()
	}
	return true
}

// dropsKind is exists for a DROP ... IF EXISTS whose key is in its kind's
// own space, as a table's in the table space: whether or not it drops
// anything, the key is Absent after it.
func (b *builder) dropsKind(k walk.Key, kinds []string, ifExists bool, rng review.Range, absent, present string) bool {
	if !b.exists(k, kinds, ifExists, rng, absent, present) {
		return false
	}
	if ifExists {
		b.settle(k, walk.Absent)
	}
	return true
}

// dropTable follows DROP TABLE [IF EXISTS] [CASCADE CONSTRAINTS] [PURGE]:
// the table must exist (00942; a view, a synonym, and any other kind too,
// verified), and a materialized view is refused even with IF EXISTS
// (12083). It takes the table's columns, constraints, indexes, and
// triggers, whose names are free at once, the recycle bin or not
// (verified).
func (b *builder) dropTable(v *ast.DropStmt, r objRef) bool {
	k := r.key(spaceTable)
	if !v.IfExists && b.heldByOther(k, r, kindTable) {
		return true
	}
	mview := wrongKind(r.written, kindMView, kindWord(kindTable))
	if v.IfExists {
		if f, kind := b.look(k); f == walk.Present && kind == kindMView {
			b.need(k, walk.Present, []string{kindTable}, r.rng, "", mview)
			return true
		}
	}
	if !b.dropsKind(k, []string{kindTable}, v.IfExists, r.rng, missingTable(r.written), mview) {
		return true
	}
	b.Last = b.domain(r)
	b.forgetStorage(r)
	b.forgetLogs(r)
	b.leaveObject(r)
	b.Clear(walk.Scope{Parent: r.path(), Deep: true})
	b.dropOwned(r)
	if v.Cascade {
		b.cascadeConstraints()
	}
	return true
}

// truncate follows TRUNCATE TABLE, which follows no synonym: the table
// must exist (00942), and a view (01702), a sequence (02201), or a type or
// routine (04044) is refused with its own error (verified).
func (b *builder) truncate(v *ast.TruncateStmt) bool {
	if v.Cluster {
		return true
	}
	r, ok := b.ref(v.Table)
	if !ok {
		return false
	}
	b.Last = b.domain(r)
	if f, _ := b.look(r.key(spaceTable)); f != walk.Absent {
		return true
	}
	b.wrongTarget(r, true)
	return true
}

// wrongTarget adds the check that a statement's table, which the walk
// knows is no table, is not a view, a sequence, or a type or routine,
// each of which the statement refuses with an error of its own (01702,
// 02201, 04044); a statement that follows no synonym (TRUNCATE) refuses
// any other name as missing (00942).
func (b *builder) wrongTarget(r objRef, followsNone bool) {
	ok := r.key(spaceObject)
	f, kind := b.look(ok)
	var msg string
	switch kind {
	case kindView, kindSequence, kindType, kindProcedure, kindFunction, kindPackage:
		msg = wrongKind(r.written, kind, kindWord(kindTable))
	}
	switch {
	case f == walk.Present && msg != "":
		b.need(ok, walk.Present, []string{kindTable, kindMView}, r.rng, "", msg)
	case followsNone:
		b.need(r.key(spaceTable), walk.Present, nil, r.rng, missingTable(r.written), "")
	}
}

// createIndex follows CREATE [UNIQUE | BITMAP] INDEX [IF NOT EXISTS]. The
// index name must be free among the schema's indexes (00955). The table
// follows synonyms, so a missing one claims nothing; a view, a sequence,
// or a type or routine of the name is refused (01702, 02201, 04044). An
// index on an expression that may call code of a schema ends the walk
// after its checks.
func (b *builder) createIndex(v *ast.CreateIndexStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	if b.refTooLong(v.Name, r) {
		return true
	}
	b.Order = walk.Lookups
	kind := ""
	if v.Cluster != nil {
		cl, ok := b.ref(v.Cluster)
		if !ok {
			return false
		}
		if f, got := b.look(cl.key(spaceObject)); f == walk.Present && got != "" {
			b.needNote(cl.key(spaceObject), walk.Present, []string{"CLUSTER"}, &walkNote{
				rng: cl.rng, present: wrongKind(cl.written, got, "cluster"),
				absentCode: review.ObjectNotFound, presentCode: review.WrongObjectKind,
				written: cl.written, kind: "CLUSTER",
			})
			return true
		}
	} else {
		t, ok := b.ref(v.Table)
		if !ok {
			return false
		}
		switch f, got := b.look(t.key(spaceTable)); f {
		case walk.Present:
			kind = t.path()
			if got == kindTable && v.FromTables == nil && v.Columns != nil {
				seenCols := make(map[string]bool)
				for _, item := range v.Columns.Items {
					ic, ok := item.(*ast.IndexColumn)
					if !ok {
						continue
					}
					cr, ok := ic.Expr.(*ast.ColumnRef)
					if !ok || cr.Schema != "" {
						continue
					}
					if cr.Table != "" && cr.Table != t.name && (v.Alias == "" || cr.Table != v.Alias) {
						continue
					}
					if cr.Column == "" {
						continue
					}
					rng := rangeOf(cr.Loc)
					if rng == (review.Range{}) {
						rng = r.rng
					}
					if seenCols[cr.Column] {
						b.fail(rng, review.DuplicateName, twice("lists", "column", cr.Column))
						return true
					}
					seenCols[cr.Column] = true
					b.columnNeeds(t, cr.Column, rng)
				}
			}
		case walk.Absent:
			b.wrongTarget(t, false)
		}
		// An indextype's code runs to build a domain index, and can make
		// tables of its own.
		b.Last = b.callsUserCode(v) || v.IndexType != nil
	}
	k := r.key(spaceIndex)
	if v.IfNotExists {
		switch f, _ := b.look(k); f {
		case walk.Present:
			return true
		case walk.Unknown:
			b.Next()
			b.Optional()
			b.settle(k, walk.Present)
		}
	} else {
		b.need(k, walk.Absent, nil, r.rng, "", takenName(r.written))
	}
	b.Set(k, walk.Present, kind)
	return true
}

// createVectorIndex follows CREATE VECTOR INDEX [IF NOT EXISTS]: the index
// name must be free in the schema's index namespace (00955), and the table
// cannot be a view, sequence, routine, or type.
func (b *builder) createVectorIndex(v *ast.CreateVectorIndexStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	if b.refTooLong(v.Name, r) {
		return true
	}
	b.Order = walk.Lookups
	kind := ""
	t, ok := b.ref(v.TableName)
	if !ok {
		return false
	}
	switch f, got := b.look(t.key(spaceTable)); f {
	case walk.Present:
		kind = t.path()
		if got == kindMView {
			b.need(t.key(spaceTable), walk.Present, []string{kindTable}, t.rng, "", wrongKind(t.written, kindMView, kindWord(kindTable)))
			return true
		}
		if got == kindTable {
			if v.Column != "" {
				b.columnNeeds(t, v.Column, t.rng)
			}
			for _, col := range v.IncludeColumns {
				if col != "" {
					b.columnNeeds(t, col, t.rng)
				}
			}
		}
	case walk.Absent:
		if objFact, objKind := b.look(t.key(spaceObject)); objFact == walk.Present && objKind != "" && objKind != kindSynonym {
			b.wrongTarget(t, false)
		} else if !generated(spaceTable, t.name) {
			keys := []walk.Key{t.key(spaceTable), t.key(spaceSynonym)}
			if v.TableName.Schema == "" {
				keys = append(keys, publicKey(t.name))
			}
			allAbsent := true
			for _, sk := range keys {
				if sf, _ := b.look(sk); sf != walk.Absent {
					allAbsent = false
					break
				}
			}
			if allAbsent {
				b.Fail(keys, walk.Absent, noteOf(t.rng, review.ObjectNotFound, missingTable(t.written)))
				return true
			}
		}
	}
	k := r.key(spaceIndex)
	if v.IfNotExists {
		switch f, _ := b.look(k); f {
		case walk.Present:
			return true
		case walk.Unknown:
			b.Next()
			b.Optional()
			b.settle(k, walk.Present)
		}
	} else {
		b.need(k, walk.Absent, nil, r.rng, "", takenName(r.written))
	}
	// A vector index creates auxiliary internal tables and indexes.
	b.Last = true
	b.Set(k, walk.Present, kind)
	return true
}

// alterObject follows ALTER of a procedure, function, package, type, or
// private synonym: the object must exist of that kind in the object namespace.
func (b *builder) alterObject(name *ast.ObjectName, kind string, ifExists bool, missing func(string) string, needed string) bool {
	r, ok := b.ref(name)
	if !ok {
		return false
	}
	if b.refTooLong(name, r) {
		return true
	}
	b.Order = walk.Lookups
	b.exists(r.key(spaceObject), []string{kind}, ifExists, r.rng, missing(r.written), notKind(r.written, needed))
	return true
}

// alterIndex follows ALTER INDEX: the index must exist (01418); RENAME TO
// needs the new name free (00955, the same name too, verified).
func (b *builder) alterIndex(v *ast.AlterIndexStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	b.Order = walk.Lookups
	b.Last = b.domainIndex(r)
	k := r.key(spaceIndex)
	if !b.exists(k, nil, v.IfExists, r.rng, missingIndex(r.written), "") {
		return true
	}
	if v.Action != "RENAME" || v.Partition != "" || v.Subpartition != "" {
		return true
	}
	if strings.Contains(v.NewName, ".") {
		return false
	}
	to := refIn(r, v.NewName)
	if b.tooLong(v.NewName, r.rng) {
		return true
	}
	tk := to.key(spaceIndex)
	b.need(tk, walk.Absent, nil, r.rng, "", takenName(to.written))
	_, kind := b.look(k)
	b.Set(k, walk.Absent, "")
	b.Set(tk, walk.Present, kind)
	return true
}

// createView follows CREATE [OR REPLACE] [FORCE] [EDITIONING] VIEW [IF
// NOT EXISTS]. A view's constraint names must be free (02264). OR REPLACE
// of a view drops its INSTEAD OF triggers and the constraints it no longer
// declares (verified). A materialized view ends the walk.
func (b *builder) createView(v *ast.CreateViewStmt) bool {
	if v.Materialized {
		r, ok := b.ref(v.Name)
		if !ok {
			return false
		}
		if b.refTooLong(v.Name, r) {
			return true
		}
		b.Order = walk.Lookups
		if v.OnPrebuilt {
			// CREATE MATERIALIZED VIEW ... ON PREBUILT TABLE requires an existing
			// table of the same name (ORA-12059 when missing, ORA-12006 on an
			// existing materialized view, ORA-01702/02201/04044/12081 on other kinds).
			if !v.IfNotExists {
				k := r.key(spaceTable)
				if b.heldByOther(k, r, kindTable) {
					return true
				}
				if f, kind := b.look(k); f == walk.Present && kind == kindMView {
					b.need(k, walk.Present, []string{kindTable}, r.rng, "", wrongKind(r.written, kindMView, kindWord(kindTable)))
					return true
				}
				b.need(k, walk.Present, []string{kindTable}, r.rng, missingTable(r.written), wrongKind(r.written, kindMView, kindWord(kindTable)))
			}
			b.Last = true
			return true
		}
		if b.claimName(r, v.IfNotExists, false, kindMView) && !v.IfNotExists {
			b.checkMviewColumnRefs(v, r)
		}
		b.Last = true
		return true
	}
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	var d tableDecl
	if v.ConstraintName != "" {
		d.names = append(d.names, v.ConstraintName)
		d.rngs = append(d.rngs, r.rng)
	}
	for _, c := range outOfLine(v.Constraints) {
		b.declareConstraint(&d, c)
	}
	if b.refTooLong(v.Name, r) || b.namesTooLong(nil, &d) {
		return true
	}
	b.Order = walk.Lookups
	if !b.claimName(r, v.IfNotExists, v.OrReplace, kindView) {
		return true
	}
	if !v.IfNotExists && !v.Force && v.Columns == nil && !b.isObjectView(v) {
		b.checkSelectAliases(v.Query, r)
	}
	for i, n := range d.names {
		k := r.schemaKey(spaceConstraint, n)
		if f, kind := b.look(k); v.OrReplace && f == walk.Present && kind == r.path() {
			// The view's own constraint, which it declares again.
			continue
		}
		b.need(k, walk.Absent, nil, d.rngs[i], "", takenConstraint(n))
	}
	if v.OrReplace {
		b.dropOwned(r)
	}
	b.arriveObject(r, kindView)
	for _, n := range d.names {
		b.Set(r.schemaKey(spaceConstraint, n), walk.Present, r.path())
	}
	return true
}

// isObjectView reports whether a CREATE VIEW declares an object or XMLType
// view (CREATE VIEW v OF ...), whose columns come from the type.
func (b *builder) isObjectView(v *ast.CreateViewStmt) bool {
	if v.Name == nil || v.Name.Loc.End <= 0 || v.Name.Loc.End >= v.Loc.End {
		return false
	}
	toks := tokens(b.w.sql, ast.Loc{Start: v.Name.Loc.End, End: v.Loc.End})
	return len(toks) > 0 && toks[0].Type != tokenQuoted && strings.EqualFold(toks[0].Str, "OF")
}

// checkSelectAliases fails a CREATE VIEW (without FORCE or a column alias
// list) or CREATE TABLE ... AS SELECT (without column aliases) when its
// leftmost SELECT list has an unaliased non-column expression (ORA-00998).
func (b *builder) checkSelectAliases(query ast.StmtNode, r objRef) bool {
	sel, ok := query.(*ast.SelectStmt)
	if !ok || sel == nil {
		return false
	}
	for sel.Larg != nil {
		sel = sel.Larg
	}
	if sel.TargetList == nil {
		return false
	}
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Name != "" || rt.Expr == nil {
			continue
		}
		expr := rt.Expr
		for {
			p, ok := expr.(*ast.ParenExpr)
			if !ok || p == nil {
				break
			}
			expr = p.Expr
		}
		switch expr.(type) {
		case *ast.ColumnRef, *ast.Star, *ast.PseudoColumn:
			continue
		}
		rng := rangeOf(rt.Loc)
		if rng == (review.Range{}) {
			rng = r.rng
		}
		b.fail(rng, review.UnsupportedForm, "Select expression in "+review.Code(r.written)+" must have a column alias.")
		return true
	}
	return false
}

// mvKnownColumns returns the known column names of a CREATE MATERIALIZED VIEW
// when complete (from its explicit column alias list or a non-star SELECT list).
func mvKnownColumns(v *ast.CreateViewStmt) (map[string]bool, bool) {
	cols := make(map[string]bool)
	if v.Columns != nil && len(v.Columns.Items) > 0 {
		for _, item := range v.Columns.Items {
			if s, ok := item.(*ast.String); ok && s.Str != "" {
				cols[s.Str] = true
			}
		}
		return cols, len(cols) > 0
	}
	sel, ok := v.Query.(*ast.SelectStmt)
	if !ok || sel == nil {
		return nil, false
	}
	for sel.Larg != nil {
		sel = sel.Larg
	}
	if sel.TargetList == nil || len(sel.TargetList.Items) == 0 {
		return nil, false
	}
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			return nil, false
		}
		if _, isStar := rt.Expr.(*ast.Star); isStar {
			return nil, false
		}
		if cr, isCR := rt.Expr.(*ast.ColumnRef); isCR && cr.Column == "*" {
			return nil, false
		}
		if rt.Name != "" {
			cols[rt.Name] = true
		} else if cr, isCR := rt.Expr.(*ast.ColumnRef); isCR && cr.Column != "" {
			cols[cr.Column] = true
		}
	}
	return cols, true
}

// checkMviewColumnRefs checks that column names referenced in CREATE
// MATERIALIZED VIEW storage or partitioning clauses (NESTED TABLE, VARRAY,
// XMLTYPE, JSON, PARTITION BY) exist in the view's known column list
// (ORA-00904, ORA-14030).
func (b *builder) checkMviewColumnRefs(v *ast.CreateViewStmt, r objRef) {
	if v.Name == nil || v.Name.Loc.End <= 0 || v.Query == nil {
		return
	}
	qStart := v.Loc.End
	if sel, ok := v.Query.(*ast.SelectStmt); ok && sel.Loc.Start > v.Name.Loc.End {
		qStart = sel.Loc.Start
	}
	if qStart <= v.Name.Loc.End {
		return
	}
	cols, complete := mvKnownColumns(v)
	if !complete {
		return
	}
	toks := tokens(b.w.sql, ast.Loc{Start: v.Name.Loc.End, End: qStart})
	if len(toks) > 0 && toks[0].Type != tokenQuoted && strings.EqualFold(toks[0].Str, "OF") {
		// Object materialized view: columns come from the object type.
		return
	}
	checkCol := func(name string) bool {
		if name == "" || strings.EqualFold(name, "COLUMN_VALUE") || generated(spaceColumn, name) {
			return false
		}
		if !cols[name] {
			b.fail(r.rng, review.ObjectNotFound, missingIn("Column", name, r.written))
			return true
		}
		return false
	}
	isWord := func(i int, word string) bool {
		return i < len(toks) && toks[i].Type != tokenQuoted && strings.EqualFold(toks[i].Str, word)
	}
	depth := 0
	for i := 0; i < len(toks); i++ {
		if toks[i].Type == tokenOpenParen {
			depth++
			continue
		}
		if toks[i].Type == tokenCloseParen {
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		switch {
		case isWord(i, "NESTED") && isWord(i+1, "TABLE"):
			if i+2 < len(toks) {
				if name, ok := tokenName(toks[i+2]); ok && checkCol(name) {
					return
				}
			}
		case isWord(i, "VARRAY"):
			if i+1 < len(toks) {
				if name, ok := tokenName(toks[i+1]); ok && checkCol(name) {
					return
				}
			}
		case isWord(i, "XMLTYPE"):
			j := i + 1
			if isWord(j, "COLUMN") {
				j++
			}
			if j < len(toks) {
				if name, ok := tokenName(toks[j]); ok && checkCol(name) {
					return
				}
			}
		case isWord(i, "JSON") && i+2 < len(toks) && toks[i+1].Type == tokenOpenParen:
			if name, ok := tokenName(toks[i+2]); ok && checkCol(name) {
				return
			}
		case isWord(i, "PARTITION") && isWord(i+1, "BY"):
			// PARTITION BY {RANGE|LIST|HASH|SYSTEM} (col, ...)
			// Skip PARTITION BY REFERENCE (constraint_name).
			if isWord(i+2, "REFERENCE") {
				continue
			}
			if isWord(i+2, "RANGE") || isWord(i+2, "LIST") || isWord(i+2, "HASH") {
				j := i + 3
				if j < len(toks) && toks[j].Type == tokenOpenParen {
					for j++; j < len(toks) && toks[j].Type != tokenCloseParen; j++ {
						if toks[j].Type == tokenComma {
							continue
						}
						if name, ok := tokenName(toks[j]); ok && checkCol(name) {
							return
						}
					}
				}
			}
		}
	}
}

// comment follows COMMENT ON TABLE, COMMENT ON COLUMN, and COMMENT ON
// MATERIALIZED VIEW. TABLE and MATERIALIZED VIEW follow private and public
// synonyms; TABLE additionally refuses materialized views (ORA-12098).
// COLUMN checks the target object and, when it is a table with complete
// synced columns, that the column exists (ORA-00904). Other comment
// categories (EDITION, INDEXTYPE, OPERATOR, AUDIT POLICY, MINING MODEL,
// PROPERTY GRAPH) target unsynced object types and are left unchecked.
func (b *builder) comment(v *ast.CommentStmt) bool {
	toks := tokens(b.w.sql, v.Loc)
	if len(toks) < 3 {
		return true
	}
	cat := strings.ToUpper(toks[2].Str)
	switch cat {
	case "TABLE":
		r, ok := b.ref(v.Object)
		if !ok {
			return false
		}
		b.Order = walk.Lookups
		if f, kind := b.look(r.key(spaceTable)); f == walk.Present && kind == kindMView {
			b.need(r.key(spaceTable), walk.Present, []string{kindTable}, r.rng, "", wrongKind(r.written, kindMView, "table or view"))
			return true
		}
		b.missingTarget(v.Object, r)
	case "MATERIALIZED":
		r, ok := b.ref(v.Object)
		if !ok {
			return false
		}
		b.Order = walk.Lookups
		b.missingTarget(v.Object, r)
	case "COLUMN":
		r, ok := b.ref(v.Object)
		if !ok {
			return false
		}
		b.Order = walk.Lookups
		b.missingTarget(v.Object, r)
		if v.Column != "" {
			if f, kind := b.look(r.key(spaceTable)); f == walk.Present && kind == kindTable {
				b.columnNeeds(r, v.Column, r.rng)
			}
		}
	}
	return true
}

// alterView follows ALTER VIEW: the view must exist (00942). A constraint
// it adds or drops leaves its name Unknown.
func (b *builder) alterView(v *ast.AlterViewStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	b.Order = walk.Lookups
	if !v.IfExists && b.heldByOther(r.key(spaceView), r, kindView) {
		return true
	}
	if !b.exists(r.key(spaceView), nil, v.IfExists, r.rng, missingView(r.written), "") {
		return true
	}
	switch v.Action {
	case "ADD_CONSTRAINT":
		if v.Constraint != nil && v.Constraint.Name != "" {
			b.Set(r.schemaKey(spaceConstraint, v.Constraint.Name), walk.Unknown, "")
		}
	case "DROP_CONSTRAINT":
		b.Set(r.schemaKey(spaceConstraint, v.ConstraintName), walk.Unknown, "")
	}
	return true
}

// alterMaterializedView follows ALTER MATERIALIZED VIEW [IF EXISTS]:
// non-partition actions and mview-only partitioning clauses require an
// existing materialized view (ORA-12003). Table-like partition maintenance
// clauses (ADD/MODIFY/MOVE/SPLIT/MERGE/COALESCE/RENAME/DROP/TRUNCATE/EXCHANGE
// PARTITION/SUBPARTITION) delegate to the container table in Oracle 23 and
// accept a partitioned table (or synonym to one), while still failing on a
// missing name or a non-table/non-synonym object (view, sequence, routine,
// type).
func (b *builder) alterMaterializedView(v *ast.AlterMaterializedViewStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	if b.refTooLong(v.Name, r) {
		return true
	}
	b.Order = walk.Lookups
	k := r.key(spaceTable)
	if v.Action == "PARTITION" && !b.alterMviewRequiresMView(v) {
		if !v.IfExists {
			if f, _ := b.look(k); f == walk.Absent {
				if objFact, objKind := b.look(r.key(spaceObject)); objFact == walk.Present {
					if objKind == kindSynonym {
						b.Last = true
						return true
					}
					if objKind != "" && objKind != kindTable && objKind != kindMView {
						b.heldByOther(k, r, kindTable)
						return true
					}
				}
				if !generated(spaceTable, r.name) {
					keys := []walk.Key{k, r.key(spaceSynonym)}
					if v.Name.Schema == "" {
						keys = append(keys, publicKey(r.name))
					}
					allAbsent := true
					for _, sk := range keys {
						if sf, _ := b.look(sk); sf != walk.Absent {
							allAbsent = false
							break
						}
					}
					if allAbsent {
						b.Fail(keys, walk.Absent, noteOf(r.rng, review.ObjectNotFound, missingTable(r.written)))
						return true
					}
				}
			}
		}
		b.Last = true
		return true
	}
	if !v.IfExists {
		if b.heldByOther(k, r, kindMView) {
			return true
		}
		if f, kind := b.look(k); f == walk.Present && kind == kindTable {
			b.need(k, walk.Present, []string{kindMView}, r.rng, "", wrongKind(r.written, kindTable, kindWord(kindMView)))
			return true
		}
	}
	b.exists(k, []string{kindMView}, v.IfExists, r.rng, missingKind(kindMView, r.written), wrongKind(r.written, kindTable, kindWord(kindMView)))
	b.Last = true
	return true
}

// alterMviewRequiresMView reports whether an ALTER MATERIALIZED VIEW
// PARTITION action requires a materialized view specifically (ORA-12003 on
// plain tables and synonyms) rather than delegating to a table container.
func (b *builder) alterMviewRequiresMView(v *ast.AlterMaterializedViewStmt) bool {
	if v.Name == nil || v.Name.Loc.End <= 0 || v.Name.Loc.End >= v.Loc.End {
		return false
	}
	toks := tokens(b.w.sql, ast.Loc{Start: v.Name.Loc.End, End: v.Loc.End})
	if len(toks) < 2 {
		return false
	}
	w0 := strings.ToUpper(toks[0].Str)
	w1 := strings.ToUpper(toks[1].Str)
	switch w0 {
	case "SET":
		return w1 == "AUTOMATIC" || w1 == "INTERVAL" || w1 == "PARTITIONING" || w1 == "SUBPARTITION"
	case "ADD":
		return w1 == "PARTITIONING"
	case "MERGE":
		return w1 == "SUBPARTITIONS"
	case "MODIFY":
		return w1 == "DEFAULT"
	}
	return false
}

// alterSequence follows ALTER SEQUENCE, which follows no synonym (02289).
func (b *builder) alterSequence(v *ast.AlterSequenceStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	if !v.IfExists && b.heldByOther(r.key(spaceSequence), r, kindSequence) {
		return true
	}
	b.exists(r.key(spaceSequence), nil, v.IfExists, r.rng, missingSequence(r.written), "")
	return true
}

// createTrigger follows CREATE [OR REPLACE] TRIGGER of a DML or INSTEAD OF
// trigger: the name must be free among the schema's triggers (04081), and
// OR REPLACE refuses another table's trigger (04095) (verified). The table
// follows synonyms, so a missing one claims nothing. A trigger on SCHEMA
// or DATABASE ends the walk.
func (b *builder) createTrigger(v *ast.CreateTriggerStmt) bool {
	if !dmlTrigger(v) {
		return false
	}
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	t, ok := b.ref(v.Table)
	if !ok {
		return false
	}
	if b.refTooLong(v.Name, r) {
		return true
	}
	b.Order = walk.Lookups
	on := ""
	for _, space := range []string{spaceTable, spaceView} {
		if f, _ := b.look(t.key(space)); f == walk.Present {
			on = t.path()
		}
	}
	k := r.key(spaceTrigger)
	f, kind := b.look(k)
	switch {
	case v.OrReplace:
		if f == walk.Present && kind != "" && on != "" && kind != on {
			b.need(k, walk.Absent, nil, r.rng, "", "Trigger "+review.Code(r.written)+" already exists on another table.")
		}
	case v.IfNotExists:
		switch f {
		case walk.Present:
			return true
		case walk.Unknown:
			b.Next()
			b.Optional()
			b.settle(k, walk.Present)
		}
	default:
		b.need(k, walk.Absent, nil, r.rng, "", "Trigger "+review.Code(r.written)+" already exists.")
	}
	b.Set(k, walk.Present, on)
	b.Session(func() { b.w.triggers = true })
	return true
}

// dmlTrigger reports whether a CREATE TRIGGER makes a DML or INSTEAD OF
// trigger on a table or view: one whose every event is INSERT, UPDATE, or
// DELETE. A DDL or database event makes a trigger on SCHEMA, of the
// current user or a named one (ON shop.SCHEMA), or on DATABASE, which the
// parser keeps as a table named so; and a DML event on shop.SCHEMA names a
// table (00942 for a missing one, verified).
func dmlTrigger(v *ast.CreateTriggerStmt) bool {
	if v.Table == nil || v.Table.Loc.Start <= 0 && v.Table.Loc.End <= 0 || v.Events == nil || len(v.Events.Items) == 0 {
		return false
	}
	for _, e := range v.Events.Items {
		if _, ok := e.(*ast.Integer); !ok {
			return false
		}
	}
	return true
}

// alterTrigger follows ALTER TRIGGER: the trigger must exist (04080);
// RENAME TO needs the new name free (00955, the same name too, verified).
func (b *builder) alterTrigger(v *ast.AlterTriggerStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	b.Order = walk.Lookups
	k := r.key(spaceTrigger)
	if !b.exists(k, nil, v.IfExists, r.rng, missingTrigger(r.written), "") {
		return true
	}
	if v.Action != "RENAME" {
		return true
	}
	if strings.Contains(v.NewName, ".") {
		return false
	}
	to := refIn(r, v.NewName)
	if b.tooLong(v.NewName, r.rng) {
		return true
	}
	tk := to.key(spaceTrigger)
	b.need(tk, walk.Absent, nil, r.rng, "", takenName(to.written))
	_, kind := b.look(k)
	b.Set(k, walk.Absent, "")
	b.Set(tk, walk.Present, kind)
	return true
}

// rename follows RENAME a TO b, which follows no synonym: a must exist
// (04043) as a table, a view, a sequence, or a private synonym (03001 for
// another kind, 32318 for a materialized view), and b must be free (00955,
// a itself too) (verified). Run by a user who does not own the schema, it
// fails with 03001 whatever it names, which the walk does not claim; the
// claims it makes hold either way.
func (b *builder) rename(v *ast.RenameStmt) bool {
	from, ok1 := b.ref(v.OldName)
	to, ok2 := b.ref(v.NewName)
	if !ok1 || !ok2 || v.OldName.Schema != "" || v.NewName.Schema != "" {
		return false
	}
	if b.refTooLong(v.NewName, to) {
		return true
	}
	b.Order = walk.Lookups
	fk := from.key(spaceObject)
	f, kind := b.look(fk)
	present := ""
	switch kind {
	case kindMView, kindType, kindProcedure, kindFunction, kindPackage:
		present = wrongKind(from.written, kind, renamable)
	}
	b.need(fk, walk.Present, []string{kindTable, kindView, kindSequence, kindSynonym}, from.rng, missingName(from.written), present)
	b.need(to.key(spaceObject), walk.Absent, nil, to.rng, "", takenName(to.written))
	if from.name == to.name {
		return true
	}
	switch {
	case f != walk.Present || kind == "":
		// An object of a kind the walk does not know.
		b.mayRestoreSynonyms()
		b.leaveObject(from)
		b.forgetObject(to)
		b.Set(to.key(spaceObject), walk.Present, "")
		b.Move(from.path(), to.path())
		b.dropOwned(from)
	case kind == kindTable:
		b.Last = b.domain(from)
		b.moveTable(from, to)
	default:
		// to arrives first: its arrival weakens every synonym's absence,
		// and from's, which it certainly lacks, is set after.
		b.arriveObject(to, kind)
		b.leaveObject(from)
		b.dropOwned(from)
	}
	return true
}

// moveTable records a table renamed: what it holds goes with it, and its
// indexes, constraint names, and triggers keep their names, of a kind the
// walk no longer has (spec §5).
func (b *builder) moveTable(from, to objRef) {
	b.arriveObject(to, kindTable)
	b.leaveObject(from)
	b.Move(from.path(), to.path())
	b.dropOwned(from)
}

// sessionParams are the session parameters whose change leaves the names
// the walk keeps as they are, by name or by prefix.
var sessionParams = []string{
	"NLS_", "PLSQL_", "OPTIMIZER_", "TIME_ZONE", "DDL_LOCK_TIMEOUT", "RECYCLEBIN", "STATISTICS_LEVEL",
	"SKIP_UNUSABLE_INDEXES", "SQL_TRACE", "TRACEFILE_IDENTIFIER", "CURSOR_SHARING", "COMMIT_WAIT",
	"COMMIT_LOGGING", "RESULT_CACHE_MODE", "QUERY_REWRITE_ENABLED", "QUERY_REWRITE_INTEGRITY", "ISOLATION_LEVEL",
	"CONSTRAINTS", "PLSCOPE_SETTINGS", "WORKAREA_SIZE_POLICY", "SORT_AREA_SIZE", "HASH_AREA_SIZE",
	"DB_FILE_MULTIBLOCK_READ_COUNT", "PARALLEL_DEGREE_POLICY", "ERROR_ON_OVERLAP_TIME", "GLOBAL_NAMES",
	"EVENTS", "MAX_DUMP_FILE_SIZE", "TIMED_STATISTICS", "DEFAULT_COLLATION", "WEBSERVICE",
}

// alterSession follows ALTER SESSION: SET CURRENT_SCHEMA moves
// unqualified names, and the parameters of sessionParams change nothing
// the walk keeps; any other parameter, such as EDITION or CONTAINER, ends
// the walk. The other actions, such as ENABLE PARALLEL DML, change no
// name.
func (b *builder) alterSession(v *ast.AlterSessionStmt) bool {
	if v.Action != "SET" || v.SetParams == nil {
		return true
	}
	for _, item := range v.SetParams.Items {
		p, ok := item.(*ast.SetParam)
		if !ok {
			return false
		}
		name := strings.ToUpper(p.Name)
		if name == "CURRENT_SCHEMA" {
			schema := ""
			if c, ok := p.Value.(*ast.ColumnRef); ok && c.Table == "" && c.Schema == "" && !outsideASCII(c.Column) {
				schema = c.Column
			}
			b.Session(func() { b.w.current = schema })
			continue
		}
		known := false
		for _, s := range sessionParams {
			known = known || name == s || strings.HasSuffix(s, "_") && strings.HasPrefix(name, s)
		}
		if !known {
			return false
		}
	}
	return true
}

// dml follows INSERT, UPDATE, DELETE, and MERGE: each target must name
// something the statement can write (missingTarget). One that may run code
// that changes a schema ends the walk after its checks (spec §6): it may
// fire a trigger, when the snapshot or the change has one, or a target is
// not a table the walk knows in its schema, which a synonym may make a
// table of another schema; or it may call a function of a schema.
func (b *builder) dml(n ast.Node) bool {
	b.Order = walk.Lookups
	targets, ok := dmlTargets(n)
	fires := !ok || b.w.triggers
	for _, t := range targets {
		r, ok := b.ref(t)
		if !ok {
			return false
		}
		b.missingTarget(t, r)
		if f, kind := b.look(r.key(spaceTable)); f != walk.Present || kind != kindTable || b.domain(r) {
			fires = true
		}
	}
	switch v := n.(type) {
	case *ast.InsertStmt:
		b.dmlInsert(v)
	case *ast.UpdateStmt:
		b.dmlUpdate(v)
	case *ast.DeleteStmt:
		b.dmlDelete(v)
	case *ast.MergeStmt:
		b.dmlMerge(v)
	}
	b.Last = fires || b.callsUserCode(n)
	return true
}

// dmlInsert checks target columns, explicit column-vs-value counts, and
// LOG ERRORS targets on INSERT statements.
func (b *builder) dmlInsert(v *ast.InsertStmt) {
	if v.MultiTable != nil {
		for _, item := range v.MultiTable.Items {
			into, ok := item.(*ast.InsertIntoClause)
			if !ok || into.Table == nil {
				continue
			}
			r, ok := b.ref(into.Table)
			if !ok {
				continue
			}
			if b.dmlCheckColumnList(r, into.Columns, "") {
				return
			}
			if into.Columns != nil && into.Values != nil && len(into.Columns.Items) > 0 && len(into.Values.Items) > 0 && len(into.Columns.Items) != len(into.Values.Items) {
				b.fail(rangeOf(into.Loc), review.UnsupportedForm, "INSERT column count does not match VALUES count.")
				return
			}
		}
		return
	}
	if v.Table == nil {
		return
	}
	r, ok := b.ref(v.Table)
	if !ok {
		return
	}
	if b.dmlCheckColumnList(r, v.Columns, aliasName(v.Alias)) {
		return
	}
	if b.dmlCheckSetClauses(r, v.SetClauses, aliasName(v.Alias)) {
		return
	}
	if v.Columns != nil && len(v.Columns.Items) > 0 {
		if v.Values != nil && len(v.Values.Items) > 0 && v.ValuesRecord == nil && v.Select == nil {
			if len(v.Columns.Items) != len(v.Values.Items) {
				b.fail(b.anchor(), review.UnsupportedForm, "INSERT column count does not match VALUES count.")
				return
			}
		} else if v.Select != nil && !v.ByName {
			if nSel, ok := selectColumnCount(v.Select); ok && nSel != len(v.Columns.Items) {
				b.fail(b.anchor(), review.UnsupportedForm, "INSERT column count does not match SELECT list count.")
				return
			}
		}
	}
	b.dmlErrorLog(v.ErrorLog, r)
}

// dmlUpdate checks READ ONLY inline view targets, target SET columns,
// unaliased self-join ambiguity (UPDATE t SET col = ... FROM t), and
// LOG ERRORS targets.
func (b *builder) dmlUpdate(v *ast.UpdateStmt) {
	if b.dmlCheckReadOnlyTarget(v.Target) {
		return
	}
	t := v.Table
	alias := aliasName(v.Alias)
	if tr, ok := v.Target.(*ast.TableRef); ok {
		if t == nil || t.Name == "" {
			t = tr.Name
		}
		if alias == "" {
			alias = aliasName(tr.Alias)
		}
	}
	if t == nil || t.Name == "" {
		return
	}
	r, ok := b.ref(t)
	if !ok {
		return
	}
	if b.dmlCheckSetClauses(r, v.SetClauses, alias) {
		return
	}
	if alias == "" && v.FromClause != nil && b.fromHasUnaliasedTable(v.FromClause, r) {
		if b.setClausesHaveUnqualified(v.SetClauses) {
			b.fail(b.anchor(), review.UnsupportedForm, "Ambiguous unqualified column in UPDATE with unaliased self-join on "+review.Code(r.written)+".")
			return
		}
	}
	b.dmlErrorLog(v.ErrorLog, r)
}

// dmlDelete checks READ ONLY inline view targets and LOG ERRORS targets.
func (b *builder) dmlDelete(v *ast.DeleteStmt) {
	if b.dmlCheckReadOnlyTarget(v.Target) {
		return
	}
	t := v.Table
	if tr, ok := v.Target.(*ast.TableRef); ok && (t == nil || t.Name == "") {
		t = tr.Name
	}
	if t == nil || t.Name == "" {
		return
	}
	r, ok := b.ref(t)
	if !ok {
		return
	}
	b.dmlErrorLog(v.ErrorLog, r)
}

// dmlMerge checks the USING source table, target UPDATE/INSERT columns,
// unaliased self-join ambiguity (MERGE INTO t USING t ... UPDATE SET col = ...),
// and LOG ERRORS targets.
func (b *builder) dmlMerge(v *ast.MergeStmt) {
	if v.Target == nil {
		return
	}
	r, ok := b.ref(v.Target)
	if !ok {
		return
	}
	targetAlias := aliasName(v.TargetAlias)
	sourceSameUnaliased := false
	if sr, ok := v.Source.(*ast.TableRef); ok && sr.Name != nil && sr.Dblink == "" {
		if sref, sok := b.ref(sr.Name); sok {
			b.missingTarget(sr.Name, sref)
			srcAlias := aliasName(v.SourceAlias)
			if srcAlias == "" {
				srcAlias = aliasName(sr.Alias)
			}
			if targetAlias == "" && srcAlias == "" && sref.schema == r.schema && sref.name == r.name {
				sourceSameUnaliased = true
			}
		}
	}
	if v.Clauses != nil {
		for _, item := range v.Clauses.Items {
			mc, ok := item.(*ast.MergeClause)
			if !ok {
				continue
			}
			if b.dmlCheckSetClauses(r, mc.UpdateSet, targetAlias) {
				return
			}
			if b.dmlCheckColumnList(r, mc.InsertCols, targetAlias) {
				return
			}
			if mc.InsertCols != nil && mc.InsertVals != nil && len(mc.InsertCols.Items) > 0 && len(mc.InsertVals.Items) > 0 && len(mc.InsertCols.Items) != len(mc.InsertVals.Items) {
				b.fail(rangeOf(mc.Loc), review.UnsupportedForm, "MERGE INSERT column count does not match VALUES count.")
				return
			}
		}
	}
	if sourceSameUnaliased {
		b.fail(b.anchor(), review.UnsupportedForm, "Ambiguous unaliased self-join in MERGE on "+review.Code(r.written)+".")
		return
	}
	b.dmlErrorLog(v.ErrorLog, r)
}

// dmlCheckReadOnlyTarget fails UPDATE or DELETE whose direct target is an
// inline view declared WITH READ ONLY.
func (b *builder) dmlCheckReadOnlyTarget(target ast.TableExpr) bool {
	sub, ok := target.(*ast.SubqueryRef)
	if !ok || sub.Restriction == nil || !sub.Restriction.ReadOnly {
		return false
	}
	rng := rangeOf(sub.Loc)
	if rng == (review.Range{}) {
		rng = b.anchor()
	}
	b.fail(rng, review.UnsupportedForm, "Cannot modify a subquery declared WITH READ ONLY.")
	return true
}

// dmlCheckColumnList verifies each target column in an INSERT column list
// has no duplicate and, when r is a known table, exists in r (ORA-00957,
// ORA-00904).
func (b *builder) dmlCheckColumnList(r objRef, cols *ast.List, alias string) bool {
	if cols == nil || len(cols.Items) == 0 {
		return false
	}
	f, kind := b.look(r.key(spaceTable))
	checkTable := f == walk.Present && kind == kindTable
	seen := make(map[string]bool, len(cols.Items))
	for _, item := range cols.Items {
		cr, ok := item.(*ast.ColumnRef)
		if !ok || cr.Column == "" || cr.Schema != "" {
			continue
		}
		if cr.Table != "" && cr.Table != r.name && (alias == "" || cr.Table != alias) {
			continue
		}
		rng := rangeOf(cr.Loc)
		if rng == (review.Range{}) {
			rng = r.rng
		}
		if seen[cr.Column] {
			b.fail(rng, review.DuplicateName, twice("names", "column", cr.Column))
			return true
		}
		seen[cr.Column] = true
		if checkTable {
			b.columnNeeds(r, cr.Column, rng)
		}
	}
	return false
}

// dmlCheckSetClauses verifies target columns on the LHS of SET clauses when
// r is a known table (ORA-00904).
func (b *builder) dmlCheckSetClauses(r objRef, setClauses *ast.List, alias string) bool {
	if setClauses == nil || len(setClauses.Items) == 0 {
		return false
	}
	f, kind := b.look(r.key(spaceTable))
	if f != walk.Present || kind != kindTable {
		return false
	}
	checkRef := func(cr *ast.ColumnRef) {
		if cr == nil || cr.Column == "" || cr.Schema != "" {
			return
		}
		if cr.Table != "" && cr.Table != r.name && (alias == "" || cr.Table != alias) {
			return
		}
		rng := rangeOf(cr.Loc)
		if rng == (review.Range{}) {
			rng = r.rng
		}
		b.columnNeeds(r, cr.Column, rng)
	}
	for _, item := range setClauses.Items {
		sc, ok := item.(*ast.SetClause)
		if !ok {
			continue
		}
		checkRef(sc.Column)
		if sc.Columns != nil {
			for _, ci := range sc.Columns.Items {
				if cr, ok := ci.(*ast.ColumnRef); ok {
					checkRef(cr)
				}
			}
		}
	}
	return false
}

// dmlErrorLog checks that a LOG ERRORS [INTO err_table] error-logging table
// exists. Without INTO, Oracle defaults to ERR$_<first 25 bytes of target>
// in the current schema.
func (b *builder) dmlErrorLog(el *ast.ErrorLogClause, target objRef) {
	if el == nil {
		return
	}
	if el.Into != nil && el.Into.Name != "" {
		if er, ok := b.ref(el.Into); ok {
			b.missingTarget(el.Into, er)
		}
		return
	}
	if b.w.current == "" || target.schema != b.w.current || target.name == "" || generated(spaceTable, target.name) {
		return
	}
	base := target.name
	if len(base) > 25 {
		base = base[:25]
	}
	errName := "ERR$_" + base
	rng := rangeOf(el.Loc)
	if rng == (review.Range{}) {
		rng = target.rng
	}
	er := objRef{schema: b.w.current, name: errName, written: ident(errName), rng: rng}
	b.missingTarget(&ast.ObjectName{Name: errName, Loc: el.Loc}, er)
}

// selectColumnCount returns the exact output column count of a SELECT
// statement when its leftmost target list contains no star expression.
func selectColumnCount(sel *ast.SelectStmt) (int, bool) {
	if sel == nil {
		return 0, false
	}
	for sel.Larg != nil {
		sel = sel.Larg
	}
	if sel.TargetList == nil || len(sel.TargetList.Items) == 0 {
		return 0, false
	}
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			return 0, false
		}
		if _, isStar := rt.Expr.(*ast.Star); isStar {
			return 0, false
		}
		if cr, isCR := rt.Expr.(*ast.ColumnRef); isCR && cr.Column == "*" {
			return 0, false
		}
	}
	return len(sel.TargetList.Items), true
}

// fromHasUnaliasedTable reports whether a FROM list contains an unaliased
// TableRef pointing to the same (schema, name) as target.
func (b *builder) fromHasUnaliasedTable(from *ast.List, target objRef) bool {
	var check func(e ast.TableExpr) bool
	check = func(e ast.TableExpr) bool {
		switch v := e.(type) {
		case *ast.TableRef:
			if v.Alias != nil || v.Name == nil || v.Dblink != "" {
				return false
			}
			r, ok := b.ref(v.Name)
			return ok && r.schema == target.schema && r.name == target.name
		case *ast.JoinClause:
			return check(v.Left) || check(v.Right)
		}
		return false
	}
	for _, item := range from.Items {
		if te, ok := item.(ast.TableExpr); ok && check(te) {
			return true
		}
	}
	return false
}

// setClausesHaveUnqualified reports whether any SET clause targets an
// unqualified column name.
func (b *builder) setClausesHaveUnqualified(setClauses *ast.List) bool {
	if setClauses == nil {
		return false
	}
	for _, item := range setClauses.Items {
		sc, ok := item.(*ast.SetClause)
		if !ok {
			continue
		}
		if sc.Column != nil && sc.Column.Table == "" && sc.Column.Schema == "" && sc.Column.Column != "" {
			return true
		}
		if sc.Columns != nil {
			for _, ci := range sc.Columns.Items {
				if cr, ok := ci.(*ast.ColumnRef); ok && cr.Table == "" && cr.Schema == "" && cr.Column != "" {
					return true
				}
			}
		}
	}
	return false
}

// missingTarget adds the check that a DML target names something the
// statement can write. Oracle looks an unqualified name up in the current
// schema, private synonyms included, and then among the public synonyms,
// and a qualified one in its schema alone (verified). A name of a
// sequence, routine, or type is refused (02201, 04044); a name nothing
// holds, as the tables, views, and synonyms show, fails with 00942. A
// type or another object the sync does not list fails all the same, with
// an error of its own (04044, 11563, 42429, verified).
func (b *builder) missingTarget(t *ast.ObjectName, r objRef) {
	if generated(spaceTable, r.name) {
		return
	}
	ok := r.key(spaceObject)
	if f, kind := b.look(ok); f == walk.Present {
		msg := ""
		switch kind {
		case kindSequence, kindType, kindProcedure, kindFunction, kindPackage:
			msg = wrongKind(r.written, kind, "table or view")
		}
		if msg != "" {
			b.need(ok, walk.Present, []string{kindTable, kindMView, kindView, kindSynonym}, r.rng, missingName(r.written), msg)
		}
		return
	}
	keys := []walk.Key{r.key(spaceTable), r.key(spaceView), r.key(spaceSynonym)}
	if t.Schema == "" {
		keys = append(keys, publicKey(r.name))
	}
	for _, k := range keys {
		if f, _ := b.look(k); f != walk.Absent {
			return
		}
	}
	b.Fail(keys, walk.Absent, noteOf(r.rng, review.ObjectNotFound, missingName(r.written)))
}

// dmlTargets lists the tables a DML statement writes, false when one is
// no table the statement names, such as an inline view.
func dmlTargets(n ast.Node) ([]*ast.ObjectName, bool) {
	var out []*ast.ObjectName
	add := func(t *ast.ObjectName, e ast.TableExpr) bool {
		if t == nil || t.Name == "" {
			ref, ok := e.(*ast.TableRef)
			if !ok || ref.Name == nil {
				return false
			}
			t = ref.Name
		}
		out = append(out, t)
		return true
	}
	switch v := n.(type) {
	case *ast.InsertStmt:
		if v.MultiTable == nil {
			return out, add(v.Table, nil)
		}
		for _, item := range v.MultiTable.Items {
			into, ok := item.(*ast.InsertIntoClause)
			if !ok || !add(into.Table, nil) {
				return out, false
			}
		}
		return out, true
	case *ast.UpdateStmt:
		return out, add(v.Table, v.Target)
	case *ast.DeleteStmt:
		return out, add(v.Table, v.Target)
	case *ast.MergeStmt:
		return out, add(v.Target, nil)
	}
	return nil, false
}
