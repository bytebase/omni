package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// statement lowers one statement (spec §5), false when the walk does not
// follow it, which ends the walk before it.
func (b *builder) statement(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.GoStmt:
		return b.goStmt(v)
	case *ast.UseStmt:
		return b.use(v)
	case *ast.CreateSchemaStmt:
		return b.createSchema(v)
	case *ast.CreateTableStmt:
		return b.createTable(v)
	case *ast.AlterTableStmt:
		return b.alterTable(v)
	case *ast.TruncateStmt:
		return b.truncate(v)
	case *ast.CreateIndexStmt:
		return b.createIndex(v)
	case *ast.AlterIndexStmt:
		return b.alterIndex(v)
	case *ast.DropStmt:
		return b.drop(v)
	case *ast.CreateViewStmt:
		return b.createModule(v.Loc, v.Name, v.OrAlter, kindView, false)
	case *ast.CreateProcedureStmt:
		return b.createModule(v.Loc, v.Name, v.OrAlter, kindProcedure, v.Number > 1)
	case *ast.CreateFunctionStmt:
		return b.createModule(v.Loc, v.Name, v.OrAlter, kindFunction, false)
	case *ast.CreateTriggerStmt:
		return b.createTrigger(v)
	case *ast.CreateSequenceStmt:
		return b.createObject(v.Name, kindSequence)
	case *ast.CreateSynonymStmt:
		return b.createObject(v.Name, kindSynonym)
	case *ast.AlterSequenceStmt:
		return b.alterSequence(v)
	case *ast.CreateTypeStmt:
		return b.createType(v)
	case *ast.ExecStmt:
		return b.exec(v)
	case *ast.InsertStmt:
		if _, runs := v.Source.(*ast.ExecStmt); runs {
			return false
		}
		return b.dml(v.Relation, v.WithClause, nil)
	case *ast.UpdateStmt:
		return b.dml(v.Relation, v.WithClause, v.FromClause)
	case *ast.DeleteStmt:
		return b.dml(v.Relation, v.WithClause, v.FromClause)
	case *ast.MergeStmt:
		return b.dml(v.Target, v.WithClause, nil)
	case *ast.SelectStmt:
		return b.selectStmt(v)
	case *ast.IfStmt:
		return b.ifStmt(v)
	case *ast.EnableDisableTriggerStmt:
		return b.enableTrigger(v)
	case *ast.SetOptionStmt:
		return b.setOption(v)
	case *ast.SetStmt:
		return strings.HasPrefix(v.Variable, "@")
	case *ast.PrintStmt, *ast.DeclareStmt, *ast.BeginTransStmt, *ast.CommitTransStmt, *ast.SaveTransStmt,
		*ast.GrantStmt, *ast.UpdateStatisticsStmt, *ast.RaiseErrorStmt, *ast.ThrowStmt, *ast.WaitForStmt:
		// None changes a name the walk keeps. SAVE TRAN keeps what the
		// change did; a ROLLBACK, which undoes it, is not here.
		return true
	}
	return false
}

// Messages, one per claim; the oracle maps each to the errors it stands
// for.
func msgExists(written string) string        { return "an object named " + written + " already exists" }
func msgSchemaMissing(written string) string { return "schema " + written + " does not exist" }
func msgMissing(kind, written string) string { return kind + " " + written + " does not exist" }
func msgNotKind(written, kind string) string { return written + " is not a " + kind }
func msgNameTwice(spelling string) string {
	return "the name " + ident(spelling) + " is given twice in the statement"
}

// note builds a check's finding at a range.
func note(rng review.Range, absent, present string) *walkNote {
	return &walkNote{rng: rng, absent: absent, present: present}
}

// goStmt follows a GO, which ends a batch. Bytebase splits only at a GO
// alone on its line, and repeats the batch for a count; omni parses GO
// anywhere a statement may start, so a GO Bytebase would not split at, or
// a count, ends the walk. A temporary table outlives a GO only in
// transaction mode, and the database a USE selected only there too
// (spec §1, verified with go-mssqldb).
func (b *builder) goStmt(v *ast.GoStmt) bool {
	if goCount(b.w.sql, v.Loc) != 1 || v.Count != 1 {
		return false
	}
	b.Session(func() {
		if b.w.db != dbTarget {
			b.w.db = dbUnknown
		}
	})
	b.Forget(walk.Scope{Space: spaceTemp, Deep: true})
	for _, space := range childSpaces {
		b.Forget(walk.Scope{Space: space, Parent: walk.Path(tempSchema), Deep: true})
	}
	return true
}

// use follows USE, which selects the database for what follows: the
// target for the spelling Bytebase gives it. A missing database fails
// (911), which the walk does not claim: databases are not synced.
func (b *builder) use(v *ast.UseStmt) bool {
	state := b.w.databaseNamed(v.Database)
	b.Session(func() { b.w.db = state })
	return true
}

// setOption follows a SET of session options. After NOEXEC, PARSEONLY,
// FMTONLY, or a SHOWPLAN option is on, statements no longer run
// (verified), and with QUOTED_IDENTIFIER off the server reads double
// quotes otherwise than omni: both end the walk.
func (b *builder) setOption(v *ast.SetOptionStmt) bool {
	value := ""
	if c, ok := v.Value.(*ast.ColumnRef); ok {
		value = strings.ToUpper(c.Column)
	}
	options := v.Options
	if len(options) == 0 {
		options = []string{v.Option}
	}
	for _, o := range options {
		switch o = strings.ToUpper(o); {
		case o == "NOEXEC", o == "PARSEONLY", o == "FMTONLY", strings.HasPrefix(o, "SHOWPLAN"):
			if value != "OFF" {
				return false
			}
		case o == "QUOTED_IDENTIFIER":
			if value != "ON" {
				return false
			}
		}
	}
	return true
}

// first reports whether the statement is its batch's first, as a view,
// procedure, function, trigger, or schema statement must be (111,
// verified); it records the failure when not.
func (b *builder) first(loc ast.Loc, what string) bool {
	if b.w.batch[b.s.index].first || b.optional {
		return true
	}
	b.fail(rangeOf(loc), what+" must be the first statement in its batch")
	return false
}

// moduleVerb returns how a module statement begins: CREATE, CREATE OR
// ALTER, or ALTER, which share an AST node.
func (b *builder) moduleVerb(loc ast.Loc, orAlter bool) string {
	if loc.Start >= 0 && loc.Start < len(b.w.sql) && strings.HasPrefix(strings.ToUpper(b.w.sql[loc.Start:min(len(b.w.sql), loc.Start+5)]), "ALTER") {
		return "ALTER"
	}
	if orAlter {
		return "CREATE OR ALTER"
	}
	return "CREATE"
}

// schemaNeeds adds the check that a creation's schema exists (2760).
func (b *builder) schemaNeeds(r objRef) {
	if r.temp || r.where != placed {
		return
	}
	b.needs(spaceSchema, walk.Path(), r.schema, nil, note(r.rng, msgSchemaMissing(ident(r.schema.spelling)), ""))
}

// lacksObject adds the check that no object holds a reference's name
// (2714), from the relation space too, which the sync lists completely.
func (b *builder) lacksObject(r objRef) {
	if r.where != placed {
		return
	}
	n := note(r.rng, "", msgExists(r.written))
	if r.temp {
		b.lacks(spaceTemp, walk.Path(), r.n, n)
		return
	}
	b.lacks(spaceObject, r.parent, r.n, n)
	b.lacks(spaceRelation, r.parent, r.n, n)
}

// createSchema follows CREATE SCHEMA without elements. A new schema holds
// nothing; one beside another of its class may share its keys, which the
// walk then does not trust.
func (b *builder) createSchema(v *ast.CreateSchemaStmt) bool {
	if !b.first(v.Loc, "CREATE SCHEMA") {
		return true
	}
	if v.Elements != nil && len(v.Elements.Items) > 0 {
		return false
	}
	sn, ok := keyed(v.Name)
	if !ok {
		return false
	}
	rng := b.anchor()
	if b.tooLong(v.Name, false, rng) {
		return true
	}
	if b.otherDatabase(sn) {
		return true
	}
	b.lacks(spaceSchema, walk.Path(), sn, note(rng, "", "schema "+ident(v.Name)+" already exists"))
	fresh, _ := b.classFact(sn.classKey(spaceSchema, walk.Path()))
	b.arrive(spaceSchema, walk.Path(), sn, "")
	if fresh == walk.Absent {
		b.Clear(walk.Scope{Parent: walk.Path(sn.class), Deep: true})
	}
	return true
}

// dropSchema follows DROP SCHEMA [IF EXISTS]. A schema that still holds
// an object, synced or made by the change, cannot be dropped (3729,
// verified); the walk claims that only for an object it certainly sees.
func (b *builder) dropSchema(v *ast.DropStmt) bool {
	if v.Names == nil || len(v.Names.Items) != 1 {
		return false
	}
	tr, ok := v.Names.Items[0].(*ast.TableRef)
	if !ok || tr.Schema != "" || tr.Database != "" || tr.Server != "" {
		return false
	}
	sn, ok := keyed(tr.Object)
	if !ok {
		return false
	}
	if b.otherDatabase(sn) {
		return true
	}
	rng := rangeOf(tr.Loc)
	root := walk.Path()
	s := b.see(spaceSchema, root, sn)
	switch {
	case !v.IfExists:
		b.needs(spaceSchema, root, sn, nil, note(rng, msgSchemaMissing(ident(tr.Object)), ""))
	case s.fact == walk.Unknown:
		b.maybe()
	}
	single := b.schemaKnown(sn)
	if single && s.fact == walk.Present {
		if k, ok := b.heldObject(sn); ok {
			b.Fail([]walk.Key{k}, walk.Present, note(rng, "", "schema "+ident(tr.Object)+" still holds objects"))
		}
	}
	b.leave(spaceSchema, root, sn)
	under := walk.Scope{Parent: walk.Path(sn.class), Deep: true}
	if single {
		b.Clear(under)
	} else {
		b.Forget(under)
	}
	return true
}

// otherDatabase reports whether a schema statement runs outside the
// target, as after USE: in another database it changes nothing the walk
// keeps; in one the walk does not know, the schema and what it holds may
// have changed.
func (b *builder) otherDatabase(sn name) bool {
	switch b.w.db {
	case dbOther:
		return true
	case dbUnknown:
		b.forgetName(spaceSchema, sn)
		b.Forget(walk.Scope{Parent: walk.Path(sn.class), Deep: true})
		return true
	}
	return false
}

// heldObject returns the key of an object the schema certainly holds: a
// synced object, or one the change made there, that the ledger still
// shows.
func (b *builder) heldObject(sn name) (walk.Key, bool) {
	parent := walk.Path(sn.class)
	var candidates []walk.Key
	if objs := b.w.snap.objects[sn.class]; objs != nil && !b.w.target.OtherSpecs {
		for class, names := range objs.classes {
			for _, n := range names {
				candidates = append(candidates, walk.Key{Space: spaceObject, Parent: parent, Name: class + exactSep + n.spelling})
			}
		}
	}
	candidates = append(candidates, b.w.made[parent]...)
	for _, k := range candidates {
		if f, _ := b.l.Lookup(k); f == walk.Present {
			return k, true
		}
	}
	return walk.Key{}, false
}

// ifStmt follows IF ... [ELSE ...] (decision M5): either branch may run,
// so each statement in them is an optional part, with no check and its
// effects weakened to what holds either way. A branch the walk cannot
// lower, or one that switches the database, ends the walk; a view,
// procedure, function, trigger, or schema cannot stand in one (156).
func (b *builder) ifStmt(v *ast.IfStmt) bool {
	saved := b.optional
	b.optional = true
	defer func() { b.optional = saved }()
	return b.branch(v.Then) && b.branch(v.Else)
}

func (b *builder) branch(n ast.Node) bool {
	switch v := n.(type) {
	case nil:
		return true
	case *ast.BeginEndStmt:
		if v == nil || v.Stmts == nil {
			return true
		}
		for _, item := range v.Stmts.Items {
			if !b.branch(item) {
				return false
			}
		}
		return true
	case *ast.IfStmt:
		if v == nil {
			return true
		}
		return b.ifStmt(v)
	case *ast.UseStmt, *ast.GoStmt:
		return false
	}
	if moduleOf(n) != moduleNone {
		return false
	}
	b.Next()
	return b.statement(n)
}

// enableTrigger follows ENABLE/DISABLE TRIGGER on a table or view; a
// trigger the change enables may fire on DML. A database or server
// trigger ends the walk.
func (b *builder) enableTrigger(v *ast.EnableDisableTriggerStmt) bool {
	if v.OnDatabase || v.OnAllServer || v.OnObject == nil {
		return false
	}
	if v.Enable {
		b.Session(func() { b.w.triggers = true })
	}
	return true
}

// exec follows EXEC of sp_rename with literal arguments and of the
// extended-property procedures, which change no name. Any other EXEC may
// run anything and ends the walk.
func (b *builder) exec(v *ast.ExecStmt) bool {
	if v.Name == nil || v.Name.Server != "" || v.Name.Database != "" || v.Name.Schema != "" && !sameName(v.Name.Schema, "sys") || v.ReturnVar != "" || v.AsLogin != "" || v.AsUser != "" || v.AtServer != "" || v.AtDataSource != "" {
		return false
	}
	switch strings.ToLower(v.Name.Object) {
	case "sp_rename":
		return b.spRename(v)
	case "sp_addextendedproperty", "sp_updateextendedproperty", "sp_dropextendedproperty":
		return true
	}
	return false
}

// dml follows INSERT, UPDATE, DELETE, and MERGE: the target must exist
// (208), which the walk knows only from the change's own facts, since the
// sync does not list synonyms. A CTE, an alias of the FROM list, or a
// table variable is no target to check. The statement ends the walk after
// its checks when a trigger may fire (spec §6): one the sync lists or the
// change created or enabled, or one of another database, which a synonym
// or a qualified name may reach.
func (b *builder) dml(target ast.TableExpr, with *ast.WithClause, from *ast.List) bool {
	tr, ok := target.(*ast.TableRef)
	if !ok {
		if b.w.triggers {
			b.Last = true
		}
		return true
	}
	check := true
	if local, ok := localTarget(tr, with, from); ok {
		// A CTE, or an alias or name of the FROM list: the table is the
		// FROM item's, if any, which the walk reads only for its triggers.
		if local == nil {
			b.Last = true
			return true
		}
		tr, check = local, false
	}
	r, ok := b.ref(tr, false)
	if !ok {
		return false
	}
	reaches := false
	switch r.where {
	case placed:
		if r.temp {
			if check {
				b.needs(spaceTemp, walk.Path(), r.n, nil, note(r.rng, msgMissing("object", r.written), ""))
			}
			break
		}
		if check {
			b.needs(spaceObject, r.parent, r.n, nil, note(r.rng, msgMissing("object", r.written), ""))
		}
		// A synonym, which the sync does not list, may lead to another
		// database.
		reaches = b.see(spaceRelation, r.parent, r.n).fact != walk.Present
	default:
		reaches = !r.temp
	}
	if b.w.triggers || reaches {
		b.Last = true
	}
	return true
}

// localTarget returns the table a DML target means when the statement
// itself names it: nil for a CTE, the FROM item an alias or name of the
// FROM list stands for. ok is false for a target the statement does not
// name.
func localTarget(tr *ast.TableRef, with *ast.WithClause, from *ast.List) (*ast.TableRef, bool) {
	if tr.Schema != "" || tr.Database != "" || tr.Server != "" {
		return nil, false
	}
	if with != nil && with.CTEs != nil {
		for _, item := range with.CTEs.Items {
			if cte, ok := item.(*ast.CommonTableExpr); ok && !differ(cte.Name, tr.Object) {
				return nil, true
			}
		}
	}
	if from == nil || len(from.Items) == 0 {
		return nil, false
	}
	for _, item := range from.Items {
		if t, ok := item.(*ast.TableRef); ok && (t.Alias != "" && !differ(t.Alias, tr.Object) || t.Alias == "" && !differ(t.Object, tr.Object)) {
			return t, true
		}
	}
	// A join or another source may hold it.
	return nil, true
}

// selectStmt follows SELECT, which creates a table with INTO.
func (b *builder) selectStmt(v *ast.SelectStmt) bool {
	into := v
	for into.Larg != nil {
		into = into.Larg
	}
	if into.IntoTable == nil {
		return true
	}
	r, ok := b.ref(into.IntoTable, true)
	if !ok {
		return false
	}
	if b.refTooLong(into.IntoTable, r) {
		return true
	}
	b.schemaNeeds(r)
	b.lacksObject(r)
	b.createdTable(r, func(path string) {
		// SELECT ... INTO adds columns the walk does not read, and no
		// index or constraint.
		b.Forget(walk.Scope{Space: spaceColumn, Parent: path})
	})
	return true
}
