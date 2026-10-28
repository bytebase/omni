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
	case *ast.AlterIndexStmt:
		return b.alterIndex(v)
	case *ast.CreateViewStmt:
		return b.createView(v)
	case *ast.AlterViewStmt:
		return b.alterView(v)
	case *ast.CreateSequenceStmt:
		return b.createObject(v.Name, kindSequence, v.IfNotExists, false)
	case *ast.AlterSequenceStmt:
		return b.alterSequence(v)
	case *ast.CreateSynonymStmt:
		if v.Public {
			// The walk never relies on a public synonym's absence.
			return true
		}
		return b.createObject(v.Name, kindSynonym, false, v.OrReplace)
	case *ast.CreateTypeStmt:
		if v.IsBody {
			// A body has a namespace of its own.
			return true
		}
		return b.createObject(v.Name, kindType, v.IfNotExists, v.OrReplace)
	case *ast.CreatePackageStmt:
		if v.IsBody {
			return true
		}
		return b.createObject(v.Name, kindPackage, v.IfNotExists, v.OrReplace)
	case *ast.CreateProcedureStmt:
		return b.createObject(v.Name, kindProcedure, v.IfNotExists, v.OrReplace)
	case *ast.CreateFunctionStmt:
		return b.createObject(v.Name, kindFunction, v.IfNotExists, v.OrReplace)
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
	case *ast.CommentStmt, *ast.GrantStmt, *ast.RevokeStmt, *ast.LockTableStmt, *ast.AnalyzeStmt,
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
			b.need(k, walk.Absent, nil, r.rng, "", takenName(r.written))
		}
	default:
		b.need(k, walk.Absent, nil, r.rng, "", takenName(r.written))
	}
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
			b.fail(d.rngs[i], takenConstraint(n))
			return true
		}
		seen[n] = true
	}
	if d.primaries > 1 {
		b.fail(r.rng, "Table "+r.written+" can have only one primary key")
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
		b.need(r.childKey(spacePrimary, ""), walk.Absent, nil, r.rng, "", "Table "+r.written+" can have only one primary key")
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
// which is a constraint, ends the walk (spec §6).
func (b *builder) createTable(v *ast.CreateTableStmt) bool {
	if v.Private || v.OrReplace || v.Sharded || v.Duplicated || v.Parent != nil || v.Organization == "EXTERNAL" {
		return false
	}
	if v.Options != nil {
		for _, item := range v.Options.Items {
			if o, ok := item.(*ast.DDLOption); ok && (o.Key == "NESTED" || o.Key == "ORGANIZATION EXTERNAL" || o.Key == "SUPPLEMENTAL LOG") {
				return false
			}
		}
	}
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	columns := columnDefs(v.Columns)
	var d tableDecl
	for _, c := range columns {
		b.declareColumn(&d, c)
	}
	for _, c := range outOfLine(v.Constraints) {
		b.declareConstraint(&d, c)
	}
	if b.refTooLong(v.Name, r) || b.namesTooLong(columns, &d) || b.columnsTwice(columns, r) {
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
func (b *builder) columnsTwice(columns []*ast.ColumnDef, r objRef) bool {
	seen := make(map[string]bool)
	for _, c := range columns {
		if seen[c.Name] {
			b.fail(rangeOf(c.Loc), "Column "+ident(c.Name)+" is named twice in table "+r.written)
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
			return true
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
		if !b.dropsKind(r.key(spaceView), nil, v.IfExists, r.rng, missingTable(r.written), "") {
			return true
		}
		b.leaveObject(r)
		b.dropOwned(r)
		if v.Cascade {
			b.cascadeConstraints()
		}
	case ast.OBJECT_SEQUENCE:
		if !b.dropsKind(r.key(spaceSequence), nil, v.IfExists, r.rng, missingSequence(r.written), "") {
			return true
		}
		b.leaveObject(r)
	case ast.OBJECT_SYNONYM:
		msg := "Private synonym " + r.written + " to be dropped does not exist"
		if !b.exists(r.key(spaceObject), []string{kindSynonym}, v.IfExists, r.rng, msg, msg) {
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
		msg := missingObject(r.written)
		if !b.exists(r.key(spaceObject), []string{kind}, v.IfExists, r.rng, msg, msg) {
			return true
		}
		b.leaveObject(r)
	case ast.OBJECT_INDEX:
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
	mview := "Must use DROP MATERIALIZED VIEW to drop " + r.written
	if v.IfExists {
		if f, kind := b.look(k); f == walk.Present && kind == kindMView {
			b.need(k, walk.Present, []string{kindTable}, r.rng, "", mview)
			return true
		}
	}
	if !b.dropsKind(k, []string{kindTable}, v.IfExists, r.rng, missingTable(r.written), mview) {
		return true
	}
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
	case kindView:
		msg = "A view is not appropriate here: " + r.written
	case kindSequence:
		msg = "A sequence is not allowed here: " + r.written
	case kindType, kindProcedure, kindFunction, kindPackage:
		msg = "A procedure, function, package, or type is not allowed here: " + r.written
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
	if v.Cluster == nil {
		t, ok := b.ref(v.Table)
		if !ok {
			return false
		}
		switch f, _ := b.look(t.key(spaceTable)); f {
		case walk.Present:
			kind = t.path()
		case walk.Absent:
			b.wrongTarget(t, false)
		}
		b.Last = b.callsUserCode(v)
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

// alterIndex follows ALTER INDEX: the index must exist (01418); RENAME TO
// needs the new name free (00955, the same name too, verified).
func (b *builder) alterIndex(v *ast.AlterIndexStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	b.Order = walk.Lookups
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
		return false
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

// alterView follows ALTER VIEW: the view must exist (00942). A constraint
// it adds or drops leaves its name Unknown.
func (b *builder) alterView(v *ast.AlterViewStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	b.Order = walk.Lookups
	if !b.exists(r.key(spaceView), nil, v.IfExists, r.rng, missingTable(r.written), "") {
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

// alterSequence follows ALTER SEQUENCE, which follows no synonym (02289).
func (b *builder) alterSequence(v *ast.AlterSequenceStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
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
	if v.Table == nil || v.Table.Loc.Start <= 0 && v.Table.Loc.End <= 0 {
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
			b.need(k, walk.Absent, nil, r.rng, "", "Trigger "+r.written+" already exists on another table, cannot replace it")
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
		b.need(k, walk.Absent, nil, r.rng, "", "Trigger "+r.written+" already exists")
	}
	b.Set(k, walk.Present, on)
	b.Session(func() { b.w.triggers = true })
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
	case kindMView:
		present = "Cannot rename materialized view " + from.written
	case kindType, kindProcedure, kindFunction, kindPackage:
		present = "RENAME cannot rename " + article(kind) + ": " + from.written
	}
	b.need(fk, walk.Present, []string{kindTable, kindView, kindSequence, kindSynonym}, from.rng, missingObject(from.written), present)
	b.need(to.key(spaceObject), walk.Absent, nil, to.rng, "", takenName(to.written))
	if from.name == to.name {
		return true
	}
	switch {
	case f != walk.Present || kind == "":
		// An object of a kind the walk does not know.
		b.leaveObject(from)
		b.forgetObject(to)
		b.Set(to.key(spaceObject), walk.Present, "")
		b.Move(from.path(), to.path())
		b.dropOwned(from)
	case kind == kindTable:
		b.moveTable(from, to)
	default:
		b.leaveObject(from)
		b.arriveObject(to, kind)
		b.dropOwned(from)
	}
	return true
}

// moveTable records a table renamed: what it holds goes with it, and its
// indexes, constraint names, and triggers keep their names, of a kind the
// walk no longer has (spec §5).
func (b *builder) moveTable(from, to objRef) {
	b.leaveObject(from)
	b.arriveObject(to, kindTable)
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

// dml follows INSERT, UPDATE, DELETE, and MERGE, which name no object the
// walk checks: a target follows synonyms. One that may run code that
// changes a schema ends the walk after its checks (spec §6): it may fire a
// trigger, when the snapshot or the change has one, or a target is not a
// table the walk knows in its schema, which a public synonym may make a
// table of another schema; or it may call a function of a schema.
func (b *builder) dml(n ast.Node) bool {
	targets, ok := dmlTargets(n)
	fires := !ok || b.w.triggers
	for _, t := range targets {
		r, ok := b.ref(t)
		if !ok {
			return false
		}
		if f, kind := b.look(r.key(spaceTable)); f != walk.Present || kind != kindTable {
			fires = true
		}
	}
	b.Last = fires || b.callsUserCode(n)
	return true
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
