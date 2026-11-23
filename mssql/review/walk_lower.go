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
	case *ast.AlterSchemaStmt:
		return b.alterSchema(v)
	case *ast.CreateTableStmt:
		return b.createTable(v)
	case *ast.CreateTableAsSelectStmt:
		// Azure Synapse's CTAS: SQL Server has none (102, 173).
		b.fail(b.anchor(), review.UnsupportedForm, "SQL Server doesn't accept `CREATE TABLE AS SELECT`.")
		return true
	case *ast.AlterTableStmt:
		return b.alterTable(v)
	case *ast.TruncateStmt:
		return b.truncate(v)
	case *ast.CreateIndexStmt:
		return b.createIndex(v)
	case *ast.CreateXmlIndexStmt:
		return b.createSpecializedIndex(v.Name, v.Table, v.XmlColumn, v.UsingIndex, v.Options, kindXMLIndex, true, false)
	case *ast.CreateSelectiveXmlIndexStmt:
		return b.createSpecializedIndex(v.Name, v.Table, v.XmlColumn, "", v.Options, kindXMLIndex, true, false)
	case *ast.CreateSpatialIndexStmt:
		return b.createSpecializedIndex(v.Name, v.Table, v.SpatialColumn, "", v.Options, kindIndex, false, true)
	case *ast.CreateStatisticsStmt:
		return b.createStatistics(v)
	case *ast.UpdateStatisticsStmt:
		return b.updateStatistics(v)
	case *ast.AlterIndexStmt:
		return b.alterIndex(v)
	case *ast.DropStmt:
		return b.drop(v)
	case *ast.CreateViewStmt:
		return b.createView(v)
	case *ast.CreateProcedureStmt:
		if !b.createModule(v.Loc, v.Name, v.OrAlter, kindProcedure, v.Number > 1) {
			return false
		}
		if v.ExternalName == "" {
			defer b.enterModule(b.moduleSchema(v.Name, v.Loc, v.OrAlter))()
			b.moduleBody(v.Body, paramNames(v.Params), typedParams(v.Params), nil, "", "")
		}
		return true
	case *ast.CreateFunctionStmt:
		return b.createFunction(v)
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
	case *ast.SecurityStmt:
		if strings.EqualFold(v.Action, "CREATE") && (strings.EqualFold(v.ObjectType, "DEFAULT") || strings.EqualFold(v.ObjectType, "RULE")) {
			return b.createLegacyObject(v)
		}
		return false
	case *ast.ExecStmt:
		return b.exec(v)
	case *ast.InsertStmt:
		return b.insertStmt(v)
	case *ast.UpdateStmt:
		return b.updateStmt(v)
	case *ast.DeleteStmt:
		return b.deleteStmt(v)
	case *ast.MergeStmt:
		return b.dml(v.Target, v.WithClause, nil, true)
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
	case *ast.DeclareStmt:
		return b.declareStmt(v)
	case *ast.PrintStmt, *ast.BeginTransStmt, *ast.CommitTransStmt, *ast.SaveTransStmt,
		*ast.GrantStmt, *ast.RaiseErrorStmt, *ast.ThrowStmt, *ast.WaitForStmt:
		// None changes a name the walk keeps. SAVE TRAN keeps what the
		// change did; a ROLLBACK, which undoes it, is not here.
		return true
	}
	return false
}

// Messages, one shape per finding code
// (docs/specs/2026-11-02-review-finding-messages.md, Appendix E); the
// oracle maps each to the errors it stands for. A name comes as object or
// ident writes it, and a kind as kindWords names it.

// kindWords are the walk's kinds as a message names them.
var kindWords = map[string]string{
	kindTable:     "table",
	kindView:      "view",
	kindSequence:  "sequence",
	kindProcedure: "procedure",
	kindFunction:  "function",
	kindTrigger:   "trigger",
	kindSynonym:   "synonym",
	kindDefault:   "default",
	kindRule:      "rule",
	kindIndex:     "index",
	kindXMLIndex:  "XML index",
	kindPK:        "primary key",
	kindUQ:        "unique constraint",
	kindFK:        "foreign key",
	kindCK:        "check constraint",
	kindDF:        "default constraint",
}

// msgMissing says that an object of a kind doesn't exist (WT101), with
// no kind word for a kind the walk does not know.
func msgMissing(kind, written string) string {
	if kind == "" {
		return review.Code(written) + " doesn't exist."
	}
	return upperFirst(kind) + " " + review.Code(written) + " doesn't exist."
}

// msgMissingIn says that a table's column, index, or constraint doesn't
// exist (WT101).
func msgMissingIn(kind, written, table string) string {
	return upperFirst(kind) + " " + review.Code(written) + " doesn't exist in " + review.Code(table) + "."
}

// msgExists says that an object of a kind holds the name already (WT201).
func msgExists(kind, written string) string {
	return upperFirst(kind) + " " + review.Code(written) + " already exists."
}

// msgExistsIn says that a table's column or index holds the name already
// (WT201).
func msgExistsIn(kind, written, table string) string {
	return upperFirst(kind) + " " + review.Code(written) + " already exists in " + review.Code(table) + "."
}

// msgTwice says that the statement uses a name twice (WT401): verb is
// what it does with the name, and kind is "" for two objects of different
// kinds.
func msgTwice(verb, kind, written string) string {
	if kind != "" {
		kind += " "
	}
	return "This statement " + verb + " " + kind + review.Code(written) + " twice."
}

// taken writes what the check that the name the statement gives an
// object of a kind is free says when it is not (WT201): the kind's text
// when the object found is of that kind, else the name's alone.
func taken(kind, written string) func(string) string {
	return func(found string) string {
		if w := kindWords[kind]; w != "" && found == kind {
			return msgExists(w, written)
		}
		return "The name " + review.Code(written) + " is already taken."
	}
}

// notKind writes what the check that a name is of one of kinds, named
// need, says when it is of another (WT301): the kind found, when the
// ledger tells it.
func notKind(written, need string, kinds ...string) func(string) string {
	return func(found string) string {
		if w := kindWords[found]; w != "" && !contains(kinds, found) {
			return review.Code(written) + " is " + article(w) + ", not " + article(need) + "."
		}
		return review.Code(written) + " isn't " + article(need) + "."
	}
}

// text writes one text whatever the kind found.
func text(s string) func(string) string {
	return func(string) string { return s }
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// article puts "a" or "an" before a kind word.
func article(word string) string {
	if strings.HasPrefix(word, "index") || strings.HasPrefix(word, "XML") {
		return "an " + word
	}
	return "a " + word
}

// note builds a check's finding at a range.
func note(rng review.Range, absent string, present func(string) string) *walkNote {
	return &walkNote{rng: rng, absent: absent, present: present}
}

// as fixes the note's code, for a Fail or a check whose violation is
// another condition than its facts name.
func (n *walkNote) as(code review.FindingCode) *walkNote {
	n.fixed = code
	return n
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
		b.w.declaredVars = nil
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
	b.fail(rangeOf(loc), review.NotFirstInBatch, review.Code(what)+" must be the first statement of its batch.")
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
	b.needs(spaceSchema, walk.Path(), r.schema, nil, note(r.rng, msgMissing("schema", ident(r.schema.spelling)), nil))
}

// lacksObject adds the check that no object holds the name a reference
// gives an object of a kind (2714), from the relation space too, which
// the sync lists completely.
func (b *builder) lacksObject(r objRef, kind string) {
	if r.where != placed {
		return
	}
	n := note(r.rng, "", taken(kind, r.written))
	if r.temp {
		b.lacks(spaceTemp, walk.Path(), r.n, n)
		return
	}
	b.lacks(spaceObject, r.parent, r.n, n)
	b.lacks(spaceRelation, r.parent, r.n, n)
}

// createSchema follows CREATE SCHEMA. A new schema without elements holds
// nothing; one beside another of its class may share its keys, which the
// walk then does not trust. A CREATE SCHEMA with inline elements still
// checks that the schema name is free and then ends the walk.
func (b *builder) createSchema(v *ast.CreateSchemaStmt) bool {
	if !b.first(v.Loc, "CREATE SCHEMA") {
		return true
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
	b.lacks(spaceSchema, walk.Path(), sn, note(rng, "", text(msgExists("schema", ident(v.Name)))))
	if v.Elements != nil && len(v.Elements.Items) > 0 {
		if b.see(spaceSchema, walk.Path(), sn).fact != walk.Present {
			for _, item := range v.Elements.Items {
				if cv, ok := item.(*ast.CreateViewStmt); ok {
					if starRng, bad := selectStarWithoutFrom(cv.Query); bad {
						b.fail(starRng, review.UnsupportedForm, "`SELECT *` requires a `FROM` clause.")
						break
					}
				}
			}
		}
		b.Last = true
		return true
	}
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
		b.needs(spaceSchema, root, sn, nil, note(rng, msgMissing("schema", ident(tr.Object)), nil))
	case s.fact == walk.Unknown:
		b.maybe()
	}
	single := b.schemaKnown(sn)
	if single && s.fact == walk.Present {
		if k, ok := b.heldObject(sn); ok {
			b.Fail([]walk.Key{k}, walk.Present, note(rng, "", text("Schema "+review.Code(ident(tr.Object))+" still holds objects.")).as(review.SchemaNotEmpty))
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

// enableTrigger follows ENABLE/DISABLE TRIGGER on a table or view: its
// target must exist as a table or view (1088), each named trigger must be
// in the target's schema and not of another object kind (1088, 2103), and
// a trigger the change enables may fire on DML. A database or server
// trigger ends the walk.
func (b *builder) enableTrigger(v *ast.EnableDisableTriggerStmt) bool {
	if v.OnDatabase || v.OnAllServer || v.OnObject == nil {
		return false
	}
	r, ok := b.ref(v.OnObject, false)
	if !ok {
		return false
	}
	if b.refTooLong(v.OnObject, r) {
		return true
	}
	if r.where != elsewhere {
		b.tableNeeds(r, []string{kindTable, kindView}, msgMissing("table", r.written), notKind(r.written, "table or view", kindTable, kindView))
	}
	if r.where == placed && !r.temp && !v.TriggerAll && v.Triggers != nil {
		for _, item := range v.Triggers.Items {
			s, ok := item.(*ast.String)
			if !ok {
				continue
			}
			trRef := parseObjectName(s.Str)
			if trRef == nil {
				continue
			}
			written := object(trRef)
			if trRef.Schema != "" && differ(trRef.Schema, r.schema.spelling) && b.see(spaceRelation, r.parent, r.n).fact != walk.Absent {
				b.fail(b.anchor(), review.TriggerSchemaMismatch, "Trigger "+review.Code(written)+" must be in the schema of "+review.Code(r.written)+".")
				return true
			}
			tn, ok := keyed(trRef.Object)
			if !ok {
				continue
			}
			present := notKind(written, "trigger", kindTrigger)
			b.ofKind(spaceObject, r.parent, tn, []string{kindTrigger}, note(b.anchor(), "", present))
			b.ofKind(spaceRelation, r.parent, tn, []string{kindTrigger}, note(b.anchor(), "", present))
			if b.see(spaceObject, r.parent, tn).fact == walk.Absent {
				b.needs(spaceObject, r.parent, tn, []string{kindTrigger}, note(b.anchor(), msgMissing("trigger", written), present))
			}
		}
	}
	if v.Enable {
		b.Session(func() { b.w.triggers = true })
	}
	return true
}

// exec follows EXEC of sp_rename with literal arguments and of the
// extended-property procedures, which change no name. Any other EXEC
// checks that a named target isn't an existing non-executable object
// (table, view, sequence, trigger, or constraint, 2809) before ending the
// walk.
func (b *builder) exec(v *ast.ExecStmt) bool {
	if v.Name == nil {
		return false
	}
	if v.Name.Server == "" && v.Name.Database == "" && (v.Name.Schema == "" || sameName(v.Name.Schema, "sys")) && v.ReturnVar == "" && v.AsLogin == "" && v.AsUser == "" && v.AtServer == "" && v.AtDataSource == "" {
		switch strings.ToLower(v.Name.Object) {
		case "sp_rename":
			return b.spRename(v)
		case "sp_addextendedproperty", "sp_updateextendedproperty", "sp_dropextendedproperty":
			return true
		}
	}
	if b.optional {
		return false
	}
	b.execTarget(v.Name)
	b.Last = true
	return true
}

// execTarget checks that an EXEC target is not certainly a non-executable
// object (table, view, sequence, trigger, or constraint, 2809).
func (b *builder) execTarget(tr *ast.TableRef) {
	if tr == nil {
		return
	}
	r, ok := b.ref(tr, false)
	if !ok || r.where != placed || r.temp {
		return
	}
	allowed := []string{kindProcedure, kindFunction, kindSynonym}
	present := notKind(r.written, "procedure", allowed...)
	b.ofKind(spaceObject, r.parent, r.n, allowed, note(r.rng, "", present))
	b.ofKind(spaceRelation, r.parent, r.n, allowed, note(r.rng, "", present))
}

// declareStmt records batch-local @variables declared by DECLARE.
func (b *builder) declareStmt(v *ast.DeclareStmt) bool {
	if v.Variables == nil {
		return true
	}
	var names []string
	for _, item := range v.Variables.Items {
		if vd, ok := item.(*ast.VariableDecl); ok && vd.Name != "" {
			names = append(names, vd.Name)
		}
	}
	if len(names) > 0 {
		b.Session(func() {
			b.w.declaredVars = append(b.w.declaredVars, names...)
		})
	}
	return true
}

// createView follows CREATE / CREATE OR ALTER / ALTER VIEW and rejects
// SELECT * without a FROM clause (263).
func (b *builder) createView(v *ast.CreateViewStmt) bool {
	if !b.createModule(v.Loc, v.Name, v.OrAlter, kindView, false) {
		return false
	}
	defer b.enterModule(b.moduleSchema(v.Name, v.Loc, v.OrAlter))()
	b.viewQuery(v.Query, v.Columns, "view", nil, v.SchemaBinding || hasOption(v.Options, "SCHEMABINDING"))
	return true
}

// createFunction follows CREATE / CREATE OR ALTER / ALTER FUNCTION and
// rejects inline TVFs with SELECT * without FROM (263) and scalar or
// multi-statement functions whose body executes a result-returning SELECT
// (444).
func (b *builder) createFunction(v *ast.CreateFunctionStmt) bool {
	if !b.createModule(v.Loc, v.Name, v.OrAlter, kindFunction, false) {
		return false
	}
	if v.ExternalName != "" || v.Body == nil {
		return true
	}
	defer b.enterModule(b.moduleSchema(v.Name, v.Loc, v.OrAlter))()
	if v.ReturnsTable != nil && v.ReturnsTable.Variable == "" && v.ReturnsTable.Columns == nil {
		if ret, ok := v.Body.(*ast.ReturnStmt); ok {
			if sub, ok := ret.Value.(*ast.SubqueryExpr); ok {
				b.viewQuery(sub.Query, nil, "function", paramNames(v.Params), hasOption(v.Options, "SCHEMABINDING"))
			}
		}
		return true
	}
	if selRng, bad := functionBodyReturnsData(v.Body); bad {
		b.fail(selRng, review.UnsupportedForm, "SELECT statements inside a function cannot return data to a client.")
		return true
	}
	if rt := v.ReturnsTable; rt != nil && rt.Variable != "" {
		if rng, bad := returnsValue(v.Body); bad {
			b.fail(rng, review.UnsupportedForm, "SQL Server doesn't accept `RETURN` of a value in a table-valued function.")
			return true
		}
	}
	if b.functionBody(v.Body) {
		return true
	}
	if rt := v.ReturnsTable; rt != nil && rt.Columns != nil {
		var names []string
		for _, item := range rt.Columns.Items {
			if cd, ok := item.(*ast.ColumnDef); ok {
				names = append(names, cd.Name)
			}
		}
		if n := twice(names); n != "" {
			b.fail(rangeOf(rt.Loc), review.DuplicateName, msgTwice("declares", "column", ident(n)))
			return true
		}
	}
	if hasOption(v.Options, "SCHEMABINDING") && b.schemaBound(newBinder(b, v.Body, false), v.Body, "function", false) {
		return true
	}
	params := paramNames(v.Params)
	var tableVars []*ast.VariableDecl
	if rt := v.ReturnsTable; rt != nil && rt.Variable != "" {
		params = append(params, rt.Variable)
		tableVars = append(tableVars, &ast.VariableDecl{Name: rt.Variable, IsTable: true, TableDef: rt.Columns})
	}
	b.moduleBody(v.Body, params, typedParams(v.Params), tableVars, "", "")
	return true
}

// paramNames lists a module's parameters.
func paramNames(params *ast.List) []string {
	var out []string
	if params == nil {
		return nil
	}
	for _, item := range params.Items {
		if p, ok := item.(*ast.ParamDef); ok {
			out = append(out, p.Name)
		}
	}
	return out
}

// selectStarWithoutFrom reports whether a SELECT statement (or any branch
// of a set operation) selects * without a FROM clause (263).
func selectStarWithoutFrom(n ast.Node) (review.Range, bool) {
	s, ok := n.(*ast.SelectStmt)
	if !ok || s == nil {
		return review.Range{}, false
	}
	if s.Larg != nil || s.Rarg != nil {
		if rng, bad := selectStarWithoutFrom(s.Larg); bad {
			return rng, true
		}
		return selectStarWithoutFrom(s.Rarg)
	}
	if s.FromClause != nil && len(s.FromClause.Items) > 0 {
		return review.Range{}, false
	}
	if s.TargetList == nil {
		return review.Range{}, false
	}
	for _, item := range s.TargetList.Items {
		if rt, ok := item.(*ast.ResTarget); ok {
			if st, ok := rt.Val.(*ast.StarExpr); ok {
				return rangeOf(st.Loc), true
			}
		} else if st, ok := item.(*ast.StarExpr); ok {
			return rangeOf(st.Loc), true
		}
	}
	return review.Range{}, false
}

// functionBodyReturnsData reports whether a scalar or multi-statement
// function body contains a SELECT statement that returns a result set to
// the client (444), rather than assigning to local variables.
func functionBodyReturnsData(n ast.Node) (review.Range, bool) {
	switch v := n.(type) {
	case nil:
		return review.Range{}, false
	case *ast.BeginEndStmt:
		if v.Stmts == nil {
			return review.Range{}, false
		}
		for _, item := range v.Stmts.Items {
			if rng, bad := functionBodyReturnsData(item); bad {
				return rng, true
			}
		}
	case *ast.IfStmt:
		if rng, bad := functionBodyReturnsData(v.Then); bad {
			return rng, true
		}
		return functionBodyReturnsData(v.Else)
	case *ast.WhileStmt:
		return functionBodyReturnsData(v.Body)
	case *ast.SelectStmt:
		if v.IntoTable != nil {
			return review.Range{}, false
		}
		if v.Op != ast.SetOpNone {
			return rangeOf(v.Loc), true
		}
		if v.TargetList != nil && len(v.TargetList.Items) > 0 {
			for _, item := range v.TargetList.Items {
				if !isSelectVarAssign(item) {
					return rangeOf(v.Loc), true
				}
			}
		}
	}
	return review.Range{}, false
}

// returnsValue finds a RETURN of a value in a module body (178 in a
// table-valued function).
func returnsValue(n ast.Node) (review.Range, bool) {
	switch v := n.(type) {
	case *ast.BeginEndStmt:
		if v.Stmts == nil {
			return review.Range{}, false
		}
		for _, item := range v.Stmts.Items {
			if rng, bad := returnsValue(item); bad {
				return rng, true
			}
		}
	case *ast.IfStmt:
		if rng, bad := returnsValue(v.Then); bad {
			return rng, true
		}
		return returnsValue(v.Else)
	case *ast.WhileStmt:
		return returnsValue(v.Body)
	case *ast.ReturnStmt:
		if v.Value != nil {
			return rangeOf(v.Loc), true
		}
	}
	return review.Range{}, false
}

func isSelectVarAssign(item ast.Node) bool {
	if _, ok := item.(*ast.SelectAssign); ok {
		return true
	}
	if rt, ok := item.(*ast.ResTarget); ok {
		if _, ok := rt.Val.(*ast.SelectAssign); ok {
			return true
		}
	}
	return false
}

// insertStmt follows INSERT: checks the target relation, empty column list
// parsed by SQL Server as a TVF call (215), missing target columns on a
// known table (207), $IDENTITY with DEFAULT VALUES (339/207), EXEC target
// kind (2809), and SELECT * without FROM (263).
func (b *builder) insertStmt(v *ast.InsertStmt) bool {
	if !b.dml(v.Relation, v.WithClause, nil, false) {
		return false
	}
	if b.insertCounts(v) || b.setWidths(v) {
		return true
	}
	if tr, ok := v.Relation.(*ast.TableRef); ok {
		if _, isLocal := localTarget(tr, v.WithClause, nil); !isLocal {
			if r, ok := b.ref(tr, false); ok && r.where == placed {
				if v.Cols != nil && len(v.Cols.Items) == 0 && !r.temp {
					present := notKind(r.written, "function", kindFunction)
					b.ofKind(spaceObject, r.parent, r.n, []string{kindFunction}, note(r.rng, "", present))
					b.ofKind(spaceRelation, r.parent, r.n, []string{kindFunction}, note(r.rng, "", present))
				} else if v.Cols != nil && len(v.Cols.Items) > 0 {
					known := b.knownTable(r)
					path := r.tablePath()
					for _, item := range v.Cols.Items {
						s, ok := item.(*ast.String)
						if !ok {
							continue
						}
						if v.DefaultValues && strings.EqualFold(s.Str, "$IDENTITY") && !b.dmlTargetFails(r) && !b.hasDelimitedIdent("$IDENTITY") {
							b.fail(b.anchor(), review.UnsupportedForm, "Cannot insert `DEFAULT VALUES` into `$IDENTITY`.")
							return true
						}
						if isPseudoCol(s.Str) || !known {
							continue
						}
						if cn, ok := keyed(s.Str); ok {
							b.needs(spaceColumn, path, cn, nil, note(b.anchor(), msgMissingIn("column", ident(s.Str), r.written), nil))
						}
					}
				}
			}
		}
	}
	switch src := v.Source.(type) {
	case *ast.ExecStmt:
		b.execTarget(src.Name)
		b.Last = true
	case *ast.SelectStmt:
		if starRng, bad := selectStarWithoutFrom(src); bad {
			b.fail(starRng, review.UnsupportedForm, "`SELECT *` requires a `FROM` clause.")
		}
	}
	return true
}

// updateStmt follows UPDATE: checks the target relation, undeclared local
// @variables in SET and WHERE CURRENT OF (137), unmodifiable $IDENTITY
// targets (8102/207), and missing target columns when there is no FROM
// clause (207).
func (b *builder) updateStmt(v *ast.UpdateStmt) bool {
	if !b.dml(v.Relation, v.WithClause, v.FromClause, true) {
		return false
	}
	noFrom := v.FromClause == nil || len(v.FromClause.Items) == 0
	var r objRef
	var hasRef, known bool
	if noFrom {
		if tr, ok := v.Relation.(*ast.TableRef); ok {
			if _, isLocal := localTarget(tr, v.WithClause, nil); !isLocal {
				if ref, ok := b.ref(tr, false); ok {
					r, hasRef, known = ref, true, b.knownTable(ref)
				}
			}
		}
	}
	if v.SetClause != nil {
		for _, item := range v.SetClause.Items {
			se, ok := item.(*ast.SetExpr)
			if !ok {
				continue
			}
			if se.Variable != "" && !b.varDeclared(se.Variable) {
				b.fail(rangeOf(se.Loc), review.ObjectNotFound, "Variable "+review.Code(se.Variable)+" is not declared.")
				return true
			}
			if se.Column != nil && se.Column.Table == "" && se.Column.Schema == "" {
				if strings.EqualFold(se.Column.Column, "$IDENTITY") && !b.isDelimitedAt(se.Column.Loc) && (!hasRef || !b.dmlTargetFails(r)) {
					b.fail(rangeOf(se.Column.Loc), review.UnsupportedForm, "Cannot update `$IDENTITY` column.")
					return true
				}
				if known && !isPseudoCol(se.Column.Column) {
					if cn, ok := keyed(se.Column.Column); ok {
						b.needs(spaceColumn, r.tablePath(), cn, nil, note(rangeOf(se.Column.Loc), msgMissingIn("column", ident(se.Column.Column), r.written), nil))
					}
				}
			}
		}
	}
	b.checkCurrentOfVar(v.WhereClause, v.Relation, v.WithClause, v.FromClause)
	return true
}

// deleteStmt follows DELETE: checks the target relation and WHERE CURRENT
// OF @cursor_variable (137).
func (b *builder) deleteStmt(v *ast.DeleteStmt) bool {
	if !b.dml(v.Relation, v.WithClause, v.FromClause, true) {
		return false
	}
	b.checkCurrentOfVar(v.WhereClause, v.Relation, v.WithClause, v.FromClause)
	return true
}

func (b *builder) checkCurrentOfVar(where ast.ExprNode, target ast.TableExpr, with *ast.WithClause, from *ast.List) {
	co, ok := where.(*ast.CurrentOfExpr)
	if !ok || !strings.HasPrefix(co.CursorName, "@") {
		return
	}
	if tr, ok := target.(*ast.TableRef); ok {
		if _, isLocal := localTarget(tr, with, from); !isLocal {
			if r, ok := b.ref(tr, false); ok && b.dmlTargetFails(r) {
				return
			}
		}
	}
	if !b.varDeclared(co.CursorName) {
		b.fail(rangeOf(co.Loc), review.ObjectNotFound, "Variable "+review.Code(co.CursorName)+" is not declared.")
	}
}

// dmlTargetFails reports whether b.dml already added a failing check on r,
// so a subsequent text-only b.fail does not shadow the target's 208 finding.
func (b *builder) dmlTargetFails(r objRef) bool {
	if r.where != placed {
		return false
	}
	if r.temp {
		return b.see(spaceTemp, walk.Path(), r.n).fact == walk.Absent
	}
	s := b.see(spaceRelation, r.parent, r.n)
	if s.fact == walk.Absent {
		return true
	}
	return s.fact == walk.Present && s.kind != "" && s.kind != kindTable && s.kind != kindView && s.kind != kindSynonym
}

func (b *builder) hasDelimitedIdent(name string) bool {
	if b.s.loc.Start < 0 || b.s.loc.End > len(b.w.sql) || b.s.loc.End < b.s.loc.Start {
		return false
	}
	upper := strings.ToUpper(b.w.sql[b.s.loc.Start:b.s.loc.End])
	u := strings.ToUpper(name)
	return strings.Contains(upper, "["+u+"]") || strings.Contains(upper, `"`+u+`"`)
}

func (b *builder) isDelimitedAt(loc ast.Loc) bool {
	if loc.Start < 0 || loc.Start >= len(b.w.sql) {
		return false
	}
	c := b.w.sql[loc.Start]
	return c == '[' || c == '"'
}

func isPseudoCol(name string) bool {
	return strings.HasPrefix(name, "$") || strings.EqualFold(name, "IDENTITYCOL") || strings.EqualFold(name, "ROWGUIDCOL")
}

// dml follows INSERT, UPDATE, DELETE, and MERGE: the target must be a
// table, view, sequence, or synonym, which the relation space lists. On a
// name it does not hold, DML fails: with 208, or with an error of its own
// on a queue or a table-valued function (verified). A CTE, an alias of
// the FROM list, or a table variable is no target to check. The statement
// ends the walk after its checks when a trigger may fire (spec §6): one
// the sync lists or the change created or enabled, or one of another
// database, which a synonym or a qualified name may reach. An UPDATE, DELETE, or MERGE of a table
// (cascades) ends it always: a cascaded action fires the triggers of the
// tables it reaches, and those of a schema the sync excludes are not
// listed (verified).
func (b *builder) dml(target ast.TableExpr, with *ast.WithClause, from *ast.List, cascades bool) bool {
	switch tv := target.(type) {
	case *ast.TableVarRef:
		if !b.varDeclared(tv.Name) {
			b.fail(rangeOf(tv.Loc), review.ObjectNotFound, "Table variable "+review.Code(tv.Name)+" is not declared.")
			return true
		}
		if b.w.triggers {
			b.Last = true
		}
		return true
	case *ast.FuncCallExpr:
		if tv.Name != nil {
			if r, ok := b.ref(tv.Name, false); ok && r.where == placed && !r.temp {
				present := notKind(r.written, "function", kindFunction)
				b.ofKind(spaceObject, r.parent, r.n, []string{kindFunction}, note(r.rng, "", present))
				b.ofKind(spaceRelation, r.parent, r.n, []string{kindFunction}, note(r.rng, "", present))
				if b.see(spaceObject, r.parent, r.n).fact == walk.Absent {
					b.needs(spaceObject, r.parent, r.n, []string{kindFunction}, note(r.rng, msgMissing("function", r.written), present))
				}
			}
		}
		b.Last = true
		return true
	}
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
				b.needs(spaceTemp, walk.Path(), r.n, nil, note(r.rng, msgMissing("", r.written), nil))
			}
			break
		}
		if check {
			b.needs(spaceRelation, r.parent, r.n, []string{kindTable, kindView, kindSynonym}, note(r.rng, msgMissing("", r.written), notKind(r.written, "table, view, or synonym", kindTable, kindView, kindSynonym)))
		}
		// A synonym may lead to another database.
		s := b.see(spaceRelation, r.parent, r.n)
		reaches = s.fact != walk.Present || s.kind != kindTable && s.kind != kindView
	default:
		reaches = !r.temp
	}
	if b.w.triggers || reaches || cascades && !r.temp {
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

// selectStmt follows SELECT, which creates a table with INTO and rejects
// SELECT * without a FROM clause (263).
func (b *builder) selectStmt(v *ast.SelectStmt) bool {
	if starRng, bad := selectStarWithoutFrom(v); bad {
		b.fail(starRng, review.UnsupportedForm, "`SELECT *` requires a `FROM` clause.")
		return true
	}
	if b.setWidths(v) {
		return true
	}
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
	b.lacksObject(r, kindTable)
	b.createdTable(r, func(path string) {
		// SELECT ... INTO adds columns the walk does not read, and no
		// index or constraint.
		b.Forget(walk.Scope{Space: spaceColumn, Parent: path})
	})
	return true
}
