package review

import (
	"context"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// scan follows the change statement by statement for the rules that ask
// what a statement does to an object of the target's schema:
// DisallowDropConstraint and RequirePrimaryKey. It simulates nothing. It
// checks each statement against the synced schema directly, and keeps
// track of what earlier statements may have changed: a name whose meaning
// a statement may have changed is not looked up again, and a statement
// whose effect the scan cannot bound (procedural code, a rollback, a
// cascading drop of anything but a relation) ends it. Every way the scan
// can be unsure leads to no finding.
//
// A finding is a claim about its own statement, given that the statements
// before it ran: the statement it anchors never fails for a reason the
// SQL and the synced schema show. Whether the change as a whole can run
// is WalkThrough's question, not this scan's. Where the scan notices on
// the way that the server refuses a statement (a missing object, a name of
// the wrong kind, a dependent object without CASCADE) it reports nothing
// for it and ends, since it cannot say what follows; it does not look for
// every way a statement can fail. A finding after such a statement still
// describes the change as written. What user code does when it runs, a
// trigger fired by DML or a function a query calls, is outside what a
// static check can bound, and the scan does not predict it.
type scan struct {
	on    map[review.Rule]bool
	r     *reporter
	index *schemaIndex // nil without a schema

	// path is the search path as last set, "$user" unexpanded, or nil
	// when unknown. user is the current user and session the session
	// user, "" when unknown.
	path    []string
	user    string
	session string

	// touched lists the relation names a statement created, dropped,
	// renamed, or moved, by (schema, name), with "" for a schema the
	// statement left to the search path: the synced schema no longer says
	// what they name. touchedName lists the same names alone.
	touched     map[[2]string]bool
	touchedName map[string]bool
	// lastTouch numbers, by name, the last statement that touched it, as
	// touches counts them.
	lastTouch map[string]int
	touches   int
	// freed lists the (schema, name) a drop the scan resolved removed for
	// certain, until a statement uses the name again.
	freed map[[2]string]bool
	// unsettled lists the (schema, table) whose columns, constraints, or
	// inheritance a statement changed in a way the scan does not follow,
	// with "" for a schema the scan could not tell.
	unsettled map[[2]string]bool
	// constraints lists the (schema, table, constraint) names a statement
	// added, dropped, or renamed, with "" for a schema the scan could not
	// tell, and renamedKeys the names an index rename may have given to or
	// taken from a constraint, on whichever table.
	constraints map[[3]string]bool
	renamedKeys map[string]bool
	// dropped lists the (schema, table, constraint) a DROP CONSTRAINT of a
	// resolved table removed.
	dropped map[[3]string]bool
	// columns records, by (schema, table, column), with "" for a schema
	// the scan could not tell, whether a column a statement added, dropped,
	// or renamed exists now, and renamedColumns the names a column got by a
	// rename, whose identity the scan does not follow.
	columns        map[[3]string]bool
	renamedColumns map[[3]string]bool
	// origins maps a column a rename named, by (schema, table, column), to
	// the synced column it is, "" for one the change added.
	origins map[[3]string]string
	// schemas lists the schemas a statement created or dropped, and
	// schemaGone the ones a DROP SCHEMA removed and nothing made again.
	schemas    map[string]bool
	schemaGone map[string]bool
	// newReferences lists the keys the foreign keys the change created
	// reference, newReads the relations its views read, and newReturns the
	// relations whose row type its functions return, none of which the
	// synced schema records.
	newReferences []newReference
	newReads      []newDependency
	newReturns    []newDependency
	// pathVersion counts search path changes, so two unqualified names
	// are known to resolve alike only under the same path.
	pathVersion int
	// droppedFunctions lists the synced functions, by (schema, name), a
	// DROP FUNCTION certainly removed, and newFunctions counts the functions
	// the change created, by name as written.
	droppedFunctions map[tableRef]bool
	newFunctions     map[tableRef]int
	// madeRoutines lists the names the change gave a routine by creating,
	// renaming, or moving one, movedRoutines the names it took one from by
	// renaming or moving it, and replacedFunctions the synced functions,
	// by (schema, name), CREATE OR REPLACE redefined.
	madeRoutines  map[string]bool
	movedRoutines map[string]bool
	// newSignatures lists the argument signatures of the functions the
	// change created and has not dropped, by name as written.
	newSignatures     map[tableRef][]string
	replacedFunctions map[tableRef]bool
	// renamedFrom maps a name a rename or a move gave a relation the scan
	// resolved to the relation, while no statement used the name since.
	renamedFrom map[string]renamedRelation
	// created lists the tables and views the change certainly created, by
	// (schema, name), while no statement used the name since.
	created map[[2]string]*madeRelation
	// replacedViews lists the synced views CREATE OR REPLACE VIEW
	// redefined, whose synced reads no longer hold, and reowned the
	// sequences, by name, whose owner ALTER SEQUENCE changed.
	replacedViews map[tableRef]bool
	reowned       map[[2]string]bool
	// newKeys lists the (schema, table), with "" for a schema the scan
	// could not tell, a statement may have given a primary key, a unique
	// constraint, or a unique index.
	newKeys map[[2]string]bool
	// generated lists the schemas, "" for one the scan could not tell, a
	// statement made relations in under names the server generates (an
	// unnamed index, a key's index, a serial or identity column's
	// sequence), which the scan does not derive.
	generated map[string]bool

	// triggers are the target's enabled event triggers.
	triggers []*metadata.EventTriggerMetadata

	// pending are the tables RequirePrimaryKey reports at the end unless a
	// later statement settles them.
	pending []*pendingTable
	stopped bool
}

// scanTarget scans the change on one target. RequirePrimaryKey is about
// the end of the change, so it reports only when the scan reaches it.
func scanTarget(ctx context.Context, stmts []statement, on map[review.Rule]bool, target review.Target, r *reporter) error {
	s := &scan{
		on:             on,
		r:              r,
		user:           target.SessionUser,
		session:        target.SessionUser,
		touched:        make(map[[2]string]bool),
		touchedName:    make(map[string]bool),
		lastTouch:      make(map[string]int),
		unsettled:      make(map[[2]string]bool),
		freed:          make(map[[2]string]bool),
		constraints:    make(map[[3]string]bool),
		renamedKeys:    make(map[string]bool),
		dropped:        make(map[[3]string]bool),
		columns:        make(map[[3]string]bool),
		renamedColumns: make(map[[3]string]bool),
		origins:        make(map[[3]string]string),
		schemas:        make(map[string]bool),
		schemaGone:     make(map[string]bool),

		droppedFunctions: make(map[tableRef]bool),
		newFunctions:     make(map[tableRef]int),
		madeRoutines:     make(map[string]bool),
		movedRoutines:    make(map[string]bool),
		newSignatures:    make(map[tableRef][]string),
		renamedFrom:      make(map[string]renamedRelation),
		created:          make(map[[2]string]*madeRelation),
		replacedViews:    make(map[tableRef]bool),
		reowned:          make(map[[2]string]bool),
		newKeys:          make(map[[2]string]bool),
		generated:        make(map[string]bool),

		replacedFunctions: make(map[tableRef]bool),
	}
	if target.Schema != nil {
		s.index = newSchemaIndex(target.Schema)
		// A schema synced without its search path leaves it unknown.
		if path := parseSearchPath(target.Schema.GetSearchPath()); len(path) > 0 {
			s.path = path
		}
	}
	if target.Schema != nil {
		s.triggers = enabledTriggers(target.Schema)
	}
	for i := range stmts {
		if err := ctx.Err(); err != nil {
			return err
		}
		if s.triggerFires(stmts[i].node) {
			return nil
		}
		s.statement(&stmts[i])
		if s.stopped {
			return nil
		}
	}
	for _, p := range s.pending {
		if !p.settled {
			r.report(review.RequirePrimaryKey, p.statement, p.rng, p.message)
		}
	}
	return nil
}

// enabledTriggers lists the target's enabled event triggers that fire on
// statements: all but login triggers, which fire when a session connects.
func enabledTriggers(schema *metadata.DatabaseSchemaMetadata) []*metadata.EventTriggerMetadata {
	var triggers []*metadata.EventTriggerMetadata
	for _, t := range schema.GetEventTriggers() {
		if t.GetEnabled() && !strings.EqualFold(t.GetEvent(), "LOGIN") {
			triggers = append(triggers, t)
		}
	}
	return triggers
}

// triggerFires reports whether an enabled event trigger of the target may
// fire on the statement: one with no tag filter, or one whose tags name
// the statement's command, or any tag-filtered trigger for a DDL statement
// whose command tag the scan does not name. Its function may run more DDL
// or reject the statement, so from it on, what the change does is no
// longer what it says. Every DDL statement raises ddl_command_start and
// ddl_command_end, one that may drop objects sql_drop, and one that may
// rewrite a table table_rewrite.
func (s *scan) triggerFires(n ast.Node) bool {
	if len(s.triggers) == 0 || !isDDL(n) {
		return false
	}
	var tags []string
	if tag := commandTag(n); tag != "" {
		tags = []string{tag}
	}
	for _, t := range s.triggers {
		switch strings.ToUpper(t.GetEvent()) {
		case "SQL_DROP":
			if !mayDrop(n) || s.dropsNothing(n) {
				continue
			}
		case "TABLE_REWRITE":
			if !mayRewrite(n) {
				continue
			}
		}
		if len(t.GetTags()) == 0 || len(tags) == 0 {
			return true
		}
		for _, tag := range t.GetTags() {
			if slices.Contains(tags, strings.ToUpper(tag)) {
				return true
			}
		}
	}
	return false
}

// dropsNothing reports whether a statement is a DROP ... IF EXISTS of
// relations the target certainly lacks, which drops nothing.
func (s *scan) dropsNothing(n ast.Node) bool {
	v, ok := n.(*ast.DropStmt)
	if !ok || !v.Missing_ok || s.index == nil || v.Objects == nil || !isRelationKind(ast.ObjectType(v.RemoveType)) {
		return false
	}
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 || len(parts) > 2 {
			return false
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) == 2 {
			rv.Schemaname = parts[0]
		}
		if !s.missing(rv) {
			return false
		}
	}
	return true
}

// isDDL reports whether a statement may raise an event trigger's event:
// any statement but a query, DML, and the session, transaction, and
// maintenance commands event triggers do not fire on.
//
// pg: doc/src/sgml/event-trigger.sgml — Event Trigger Firing Matrix
func isDDL(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.SelectStmt:
		return selectsInto(v)
	case *ast.ExplainStmt:
		// EXPLAIN ANALYZE runs what it explains.
		if explainAnalyzes(v) {
			return isDDL(v.Query)
		}
		return false
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt,
		*ast.VariableSetStmt, *ast.VariableShowStmt, *ast.TransactionStmt,
		*ast.CopyStmt, *ast.LockStmt, *ast.PrepareStmt, *ast.ExecuteStmt, *ast.DeallocateStmt,
		*ast.NotifyStmt, *ast.ListenStmt, *ast.UnlistenStmt, *ast.CheckPointStmt, *ast.VacuumStmt,
		*ast.ConstraintsSetStmt, *ast.DiscardStmt, *ast.LoadStmt, *ast.DoStmt, *ast.CallStmt, *ast.TruncateStmt:
		return false
	}
	return true
}

// mayDrop reports whether a DDL statement may drop an object, which raises
// sql_drop: a DROP, an ALTER TABLE with a subcommand that may, or a
// statement the scan does not know not to.
func mayDrop(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.CreateStmt, *ast.CreateForeignTableStmt, *ast.CreateTableAsStmt, *ast.SelectStmt, *ast.ViewStmt,
		*ast.CreateSeqStmt, *ast.CompositeTypeStmt, *ast.IndexStmt, *ast.CreateSchemaStmt,
		*ast.RenameStmt, *ast.AlterObjectSchemaStmt, *ast.AlterSeqStmt, *ast.CommentStmt, *ast.GrantStmt:
		return false
	case *ast.AlterTableStmt:
		if v.Cmds == nil {
			return false
		}
		for _, item := range v.Cmds.Items {
			cmd, ok := item.(*ast.AlterTableCmd)
			if !ok || !keepsObjects(cmd) {
				return true
			}
		}
		return false
	}
	return true
}

// keepsObjects reports whether an ALTER TABLE subcommand drops nothing an
// sql_drop trigger sees. Replacing a default or rewriting the table
// deletes catalog rows internally, which the trigger does not see.
func keepsObjects(cmd *ast.AlterTableCmd) bool {
	switch ast.AlterTableType(cmd.Subtype) {
	case ast.AT_AddColumn, ast.AT_AddConstraint, ast.AT_AddIndexConstraint, ast.AT_SetNotNull,
		ast.AT_SetStatistics, ast.AT_ChangeOwner, ast.AT_ValidateConstraint,
		ast.AT_EnableRowSecurity, ast.AT_DisableRowSecurity:
		return true
	case ast.AT_ColumnDefault:
		// SET DEFAULT replaces; DROP DEFAULT drops.
		return cmd.Def != nil
	}
	return togglesTriggers(cmd)
}

// mayRewrite reports whether a DDL statement may rewrite a table, which
// raises table_rewrite: an ALTER TABLE with a subcommand that may, ALTER
// TYPE, or a statement the scan does not name.
func mayRewrite(n ast.Node) bool {
	v, ok := n.(*ast.AlterTableStmt)
	if !ok {
		return commandTag(n) == ""
	}
	if ast.ObjectType(v.ObjType) == ast.OBJECT_TYPE || v.Cmds == nil {
		return ast.ObjectType(v.ObjType) == ast.OBJECT_TYPE
	}
	for _, item := range v.Cmds.Items {
		cmd, ok := item.(*ast.AlterTableCmd)
		if !ok || !keepsRows(cmd) {
			return true
		}
	}
	return false
}

// keepsRows reports whether an ALTER TABLE subcommand never rewrites the
// table: one that changes only the catalog or scans the rows, or ADD
// COLUMN of a plain column.
//
// pg: src/backend/commands/tablecmds.c — ATExecAddColumn, ATRewriteTables
func keepsRows(cmd *ast.AlterTableCmd) bool {
	switch ast.AlterTableType(cmd.Subtype) {
	case ast.AT_AddColumn:
		cd, ok := cmd.Def.(*ast.ColumnDef)
		return ok && plainColumn(cd)
	case ast.AT_AddConstraint, ast.AT_AddIndexConstraint, ast.AT_SetNotNull, ast.AT_DropNotNull, ast.AT_ColumnDefault,
		ast.AT_SetStatistics, ast.AT_ChangeOwner, ast.AT_ValidateConstraint, ast.AT_EnableRowSecurity, ast.AT_DisableRowSecurity,
		ast.AT_DropConstraint, ast.AT_DropColumn:
		return true
	}
	return togglesTriggers(cmd)
}

// togglesTriggers reports whether an ALTER TABLE subcommand enables or
// disables triggers, which changes only the catalog.
func togglesTriggers(cmd *ast.AlterTableCmd) bool {
	switch ast.AlterTableType(cmd.Subtype) {
	case ast.AT_EnableTrig, ast.AT_EnableAlwaysTrig, ast.AT_EnableReplicaTrig, ast.AT_DisableTrig,
		ast.AT_EnableTrigAll, ast.AT_DisableTrigAll, ast.AT_EnableTrigUser, ast.AT_DisableTrigUser:
		return true
	}
	return false
}

// plainColumn reports whether adding a column leaves the rows as they
// are: a built-in type, which no domain constraint checks, and no
// default, identity, or generation expression, which the server may
// compute for every row.
func plainColumn(cd *ast.ColumnDef) bool {
	if cd.RawDefault != nil || cd.CookedDefault != nil || cd.Identity != 0 || cd.Generated != 0 || cd.TypeName == nil {
		return false
	}
	for _, c := range constraintsOf(cd.Constraints) {
		switch c.Contype {
		case ast.CONSTR_DEFAULT, ast.CONSTR_IDENTITY, ast.CONSTR_GENERATED:
			return false
		}
	}
	parts := nameParts(cd.TypeName.Names)
	if len(parts) == 2 && parts[0] == "pg_catalog" {
		parts = parts[1:]
	}
	return len(parts) == 1 && builtinTypes[parts[0]]
}

// builtinTypes are the names of common built-in scalar types.
var builtinTypes = map[string]bool{
	"bool": true, "boolean": true, "int2": true, "smallint": true, "int4": true, "int": true, "integer": true,
	"int8": true, "bigint": true, "float4": true, "real": true, "float8": true, "numeric": true, "decimal": true,
	"text": true, "varchar": true, "bpchar": true, "char": true, "name": true, "bytea": true, "uuid": true,
	"date": true, "time": true, "timetz": true, "timestamp": true, "timestamptz": true, "interval": true,
	"json": true, "jsonb": true, "xml": true, "inet": true, "cidr": true, "macaddr": true, "money": true,
	"bit": true, "varbit": true, "oid": true, "tsvector": true, "tsquery": true,
}

// commandTag is the command tag of a statement, as an event trigger's tag
// filter names it, or "" for one the scan does not name.
//
// pg: src/include/tcop/cmdtaglist.h
func commandTag(n ast.Node) string {
	objectTag := func(t ast.ObjectType) string {
		switch t {
		case ast.OBJECT_TABLE:
			return "TABLE"
		case ast.OBJECT_VIEW:
			return "VIEW"
		case ast.OBJECT_MATVIEW:
			return "MATERIALIZED VIEW"
		case ast.OBJECT_FOREIGN_TABLE:
			return "FOREIGN TABLE"
		case ast.OBJECT_SEQUENCE:
			return "SEQUENCE"
		case ast.OBJECT_INDEX:
			return "INDEX"
		case ast.OBJECT_TYPE:
			return "TYPE"
		case ast.OBJECT_SCHEMA:
			return "SCHEMA"
		case ast.OBJECT_FUNCTION:
			return "FUNCTION"
		case ast.OBJECT_ROUTINE:
			return "ROUTINE"
		case ast.OBJECT_PROCEDURE:
			return "PROCEDURE"
		}
		return ""
	}
	prefixed := func(verb string, t ast.ObjectType) string {
		if o := objectTag(t); o != "" {
			return verb + " " + o
		}
		return ""
	}
	switch v := n.(type) {
	case *ast.CreateStmt:
		return "CREATE TABLE"
	case *ast.CreateForeignTableStmt:
		return "CREATE FOREIGN TABLE"
	case *ast.CreateTableAsStmt:
		if v.Objtype == ast.OBJECT_MATVIEW {
			return "CREATE MATERIALIZED VIEW"
		}
		return "CREATE TABLE AS"
	case *ast.SelectStmt:
		if selectsInto(v) {
			return "SELECT INTO"
		}
	case *ast.ViewStmt:
		return "CREATE VIEW"
	case *ast.CreateSeqStmt:
		return "CREATE SEQUENCE"
	case *ast.CompositeTypeStmt:
		return "CREATE TYPE"
	case *ast.IndexStmt:
		return "CREATE INDEX"
	case *ast.CreateFunctionStmt:
		if isProcedure(v) {
			return "CREATE PROCEDURE"
		}
		return "CREATE FUNCTION"
	case *ast.AlterSeqStmt:
		return "ALTER SEQUENCE"
	case *ast.CreateSchemaStmt:
		return "CREATE SCHEMA"
	case *ast.AlterTableStmt:
		return prefixed("ALTER", ast.ObjectType(v.ObjType))
	case *ast.RenameStmt:
		if v.RenameType == ast.OBJECT_COLUMN || v.RenameType == ast.OBJECT_TABCONSTRAINT {
			if tag := prefixed("ALTER", v.RelationType); tag != "" {
				return tag
			}
			return "ALTER TABLE"
		}
		return prefixed("ALTER", v.RenameType)
	case *ast.AlterObjectSchemaStmt:
		return prefixed("ALTER", v.ObjectType)
	case *ast.DropStmt:
		return prefixed("DROP", ast.ObjectType(v.RemoveType))
	case *ast.CommentStmt:
		return "COMMENT"
	case *ast.GrantStmt:
		if v.IsGrant {
			return "GRANT"
		}
		return "REVOKE"
	case *ast.CreateEnumStmt:
		return "CREATE TYPE"
	case *ast.AlterEnumStmt:
		return "ALTER TYPE"
	case *ast.CreateDomainStmt:
		return "CREATE DOMAIN"
	case *ast.CreateTrigStmt:
		return "CREATE TRIGGER"
	case *ast.CreateExtensionStmt:
		return "CREATE EXTENSION"
	}
	return ""
}

// createdIn returns the schema a CREATE of a routine, a type, or a table
// by query writes its object's name in, or "".
func createdIn(n ast.Node) string {
	var parts []string
	switch v := n.(type) {
	case *ast.CreateFunctionStmt:
		parts = nameParts(v.Funcname)
	case *ast.CreateEnumStmt:
		parts = nameParts(v.TypeName)
	case *ast.CreateDomainStmt:
		parts = nameParts(v.Domainname)
	case *ast.CreateRangeStmt:
		parts = nameParts(v.TypeName)
	case *ast.DefineStmt:
		if v.Kind == ast.OBJECT_TYPE {
			parts = nameParts(v.Defnames)
		}
	case *ast.CompositeTypeStmt:
		if v.Typevar != nil {
			return v.Typevar.Schemaname
		}
	case *ast.CreateForeignTableStmt:
		if v.Base.Relation != nil {
			return v.Base.Relation.Schemaname
		}
	case *ast.CreateTableAsStmt:
		if v.Into != nil && v.Into.Rel != nil {
			return v.Into.Rel.Schemaname
		}
	case *ast.SelectStmt:
		if into := intoOf(v); into != nil && into.Rel != nil {
			return into.Rel.Schemaname
		}
	}
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}

// createAs follows CREATE TABLE AS, CREATE MATERIALIZED VIEW, and SELECT
// INTO, which make a relation from a query. One whose name is taken, or
// whose query reads a relation the target lacks, is refused, and under IF
// NOT EXISTS of a taken name does nothing.
func (s *scan) createAs(rel *ast.RangeVar, ifNotExists bool, query ast.Node, matview bool) {
	if rel == nil {
		return
	}
	taken := s.existsWhereCreated(rel) || s.madeHere(rel)
	switch {
	case s.schemaMissing(rel) || s.typeExists(rel) || s.readsMissing(query) || taken && !ifNotExists:
		s.stop()
		return
	case taken:
		return
	}
	s.touch(rel)
	if matview {
		s.reads(rel, query)
	}
}

// isProcedure reports whether CREATE FUNCTION is CREATE PROCEDURE, which
// the parser marks with an option.
func isProcedure(v *ast.CreateFunctionStmt) bool {
	if v.Options == nil {
		return false
	}
	return slices.ContainsFunc(v.Options.Items, func(n ast.Node) bool {
		d, ok := n.(*ast.DefElem)
		return ok && d.Defname == "isProcedure"
	})
}

// stop ends the scan: nothing after this statement is known, and no table
// is known to end the change without a key.
func (s *scan) stop() {
	s.stopped = true
}

func (s *scan) statement(st *statement) {
	// A CREATE in a schema that certainly does not exist is refused.
	if schema := createdIn(st.node); schema != "" && s.schemaMissing(&ast.RangeVar{Schemaname: schema}) {
		s.stop()
		return
	}
	switch v := st.node.(type) {
	case *ast.AlterTableStmt:
		s.alterTable(st, v)
	case *ast.CreateStmt:
		s.create(st, v)
	case *ast.CreateForeignTableStmt:
		rel := v.Base.Relation
		if rel == nil {
			return
		}
		taken := s.existsWhereCreated(rel) || s.madeHere(rel)
		switch {
		case duplicateNames(&v.Base) || checksRefused(&v.Base) || s.sourceRefused(&v.Base):
			s.stop()
			return
		case s.typeExists(rel) || taken && !v.Base.IfNotExists:
			s.stop()
			return
		case taken:
			// IF NOT EXISTS of a taken name does nothing.
			return
		}
		s.touch(rel)
		if !v.Base.IfNotExists || s.createsNew(rel) {
			s.made(rel, kindForeignTable)
		}
	case *ast.CreateTableAsStmt:
		if v.Into != nil {
			s.createAs(v.Into.Rel, v.IfNotExists, v.Query, v.Objtype == ast.OBJECT_MATVIEW)
		}
	case *ast.ViewStmt:
		// CREATE VIEW of a name its schema holds, or OR REPLACE of a
		// relation that is not a view, is refused.
		if s.schemaMissing(v.View) || !v.Replace && (s.existsWhereCreated(v.View) || s.madeHere(v.View)) || s.typeExists(v.View) || s.readsMissing(v.Query) {
			s.stop()
			return
		}
		if v.Replace {
			schema, kind, ok := s.lookup(v.View)
			if ok && kind != kindView && kind != kindAmbiguous {
				s.stop()
				return
			}
			if c, made := s.madeAt(v.View); made && c.kind != kindView {
				s.stop()
				return
			}
			if ok && kind == kindView {
				s.replacedViews[tableRef{schema, v.View.Relname}] = true
			}
			s.retire(s.newReads, tableRef{v.View.Schemaname, v.View.Relname})
		}
		s.touch(v.View)
		if !v.Replace {
			s.made(v.View, kindView)
		}
		s.reads(v.View, v.Query)
	case *ast.CreateFunctionStmt:
		// A SQL-standard body is resolved when the routine is created, and
		// a signature the schema has is taken without OR REPLACE.
		if v.SqlBody != nil && s.readsMissing(v.SqlBody) || !v.IsOrReplace && s.signatureTaken(v) {
			s.stop()
			return
		}
		s.returns(v)
	case *ast.AlterSeqStmt:
		// ALTER SEQUENCE of a relation the target lacks is refused, and
		// under IF EXISTS does nothing; of another kind, it is refused.
		if v.Sequence != nil && s.index != nil {
			if s.missing(v.Sequence) {
				if !v.MissingOk {
					s.stop()
				}
				return
			}
			_, kind, ok := s.lookup(v.Sequence)
			c, made := s.madeAt(v.Sequence)
			if ok && kind != kindSequence && kind != kindAmbiguous || made && c.kind != kindSequence {
				s.stop()
				return
			}
		}
		if v.Sequence != nil && v.Options != nil && slices.ContainsFunc(v.Options.Items, func(n ast.Node) bool {
			d, ok := n.(*ast.DefElem)
			return ok && d.Defname == "owned_by"
		}) {
			s.reowned[[2]string{s.schemaOf(v.Sequence), v.Sequence.Relname}] = true
		}
	case *ast.CreateSeqStmt:
		exists := s.existsWhereCreated(v.Sequence) || s.madeHere(v.Sequence)
		switch {
		case s.schemaMissing(v.Sequence) || exists && !v.IfNotExists:
			s.stop()
			return
		case exists:
			// IF NOT EXISTS of an existing name does nothing.
			return
		}
		s.touch(v.Sequence)
		if !v.IfNotExists || s.createsNew(v.Sequence) {
			s.made(v.Sequence, kindSequence)
		}
	case *ast.CompositeTypeStmt:
		if v.Typevar != nil && s.typeTaken(v.Typevar) {
			s.stop()
			return
		}
		s.touch(v.Typevar)
		if v.Typevar != nil {
			s.made(v.Typevar, kindCompositeType)
		}
	case *ast.CreateEnumStmt:
		// A type's name is taken for a relation's row type too.
		s.touchType(v.TypeName, "")
	case *ast.CreateDomainStmt:
		s.touchType(v.Domainname, "")
	case *ast.CreateRangeStmt:
		// A range type comes with a multirange type, named after it unless
		// the statement names it.
		s.touchType(v.TypeName, "_multirange")
	case *ast.DefineStmt:
		if v.Kind == ast.OBJECT_TYPE {
			s.touchType(v.Defnames, "")
		}
	case *ast.IndexStmt:
		if s.indexRefused(v) {
			s.stop()
			return
		}
		if v.IfNotExists && v.Idxname != "" && v.Relation != nil {
			if schema, _, ok := s.lookup(v.Relation); ok && s.nameTaken(schema, v.Idxname) {
				// IF NOT EXISTS of a taken name does nothing.
				return
			}
		}
		if (v.Unique || v.Primary) && v.Relation != nil {
			s.newKey(v.Relation)
		}
		if v.Idxname == "" {
			s.generate(v.Relation)
		}
		if v.Relation == nil || s.mayHavePartitions(v.Relation) {
			// A partition's index takes a generated name, in the
			// partition's schema.
			s.generate(nil)
		}
		if v.Idxname != "" && v.Relation != nil {
			s.touchName(v.Relation.Schemaname, v.Idxname)
		}
	case *ast.SelectStmt:
		switch {
		case v.IntoClause != nil:
			q := *v
			q.IntoClause = nil
			s.createAs(v.IntoClause.Rel, false, &q, false)
		case selectsInto(v):
			s.stop()
		default:
			s.data(v)
		}
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
		s.data(v)
	case *ast.ExplainStmt:
		if !explainAnalyzes(v) {
			return
		}
		switch q := v.Query.(type) {
		case *ast.CreateTableAsStmt:
			if q.Into != nil {
				s.createAs(q.Into.Rel, q.IfNotExists, q.Query, q.Objtype == ast.OBJECT_MATVIEW)
			}
		case *ast.SelectStmt:
			if q.IntoClause != nil || selectsInto(q) {
				s.stop()
				return
			}
			s.data(q)
		default:
			s.data(q)
		}
	case *ast.RenameStmt:
		s.rename(v)
	case *ast.AlterObjectSchemaStmt:
		s.setSchema(v)
	case *ast.DropStmt:
		s.drop(v)
	case *ast.CreateSchemaStmt:
		s.createSchema(st, v)
	case *ast.DoStmt, *ast.CallStmt:
		// Procedural code may run any DDL.
		s.stop()
	case *ast.CreateExtensionStmt:
		// An extension's script may create any object. IF NOT EXISTS of an
		// extension the target has does nothing.
		if !v.IfNotExists || s.index == nil || !s.index.extensions[v.Extname] {
			s.stop()
		}
	case *ast.DropOwnedStmt, *ast.ImportForeignSchemaStmt, *ast.AlterExtensionStmt, *ast.DiscardStmt, *ast.LoadStmt:
		// What these drop, create, or reset is not in the statement.
		s.stop()
	case *ast.CreateEventTrigStmt, *ast.AlterEventTrigStmt:
		// An event trigger runs code on every later DDL statement.
		s.stop()
	case *ast.TransactionStmt:
		switch v.Kind {
		case ast.TRANS_STMT_ROLLBACK, ast.TRANS_STMT_ROLLBACK_TO, ast.TRANS_STMT_PREPARE,
			ast.TRANS_STMT_COMMIT_PREPARED, ast.TRANS_STMT_ROLLBACK_PREPARED:
			s.stop()
		}
	case *ast.VariableSetStmt:
		s.setting(v)
	}
}

// newReference is a foreign key the change created: the key it references,
// by the table as written, "" for a schema left to the search path, and
// the columns, nil for the primary key; and the table that owns it, as
// written, with the search path it was written under, and its name, ""
// when generated.
type newReference struct {
	schema, table string
	columns       []string

	owner     tableRef
	ownerPath int
	name      string
	retired   bool
}

// references records the keys new foreign keys of the owner reference.
func (s *scan) references(owner *ast.RangeVar, constraints ...*ast.Constraint) {
	for _, c := range constraints {
		if c.Contype == ast.CONSTR_FOREIGN && c.Pktable != nil {
			s.newReferences = append(s.newReferences, newReference{
				schema: s.schemaOf(c.Pktable), table: c.Pktable.Relname, columns: nameParts(c.PkAttrs),
				owner: tableRef{owner.Schemaname, owner.Relname}, ownerPath: s.pathVersion, name: c.Conname,
			})
		}
	}
}

// retireReferences records that a statement dropped the new foreign keys of
// a table, or the one of that name when name is set. Only a table named
// as the foreign key's owner was, under the same search path when neither
// is qualified, is known to be that owner.
func (s *scan) retireReferences(owner *ast.RangeVar, name string) {
	for i := range s.newReferences {
		r := &s.newReferences[i]
		same := r.owner.table == owner.Relname && r.owner.schema == owner.Schemaname &&
			(owner.Schemaname != "" || r.ownerPath == s.pathVersion)
		if same && (name == "" || r.name == name) {
			r.retired = true
		}
	}
}

// newlyReferenced reports whether a foreign key the change created may
// reference the table's key with these columns, nil when they cannot be
// read; primary tells whether the key is the primary key. A foreign key
// of a table in except, as written, does not count.
func (s *scan) newlyReferenced(t tableRef, columns []string, primary bool, except ...tableRef) bool {
	for _, r := range s.newReferences {
		if r.retired || r.table != t.table || r.schema != "" && r.schema != t.schema {
			continue
		}
		// A foreign key of a table the statement drops too goes with it.
		if slices.ContainsFunc(except, func(o tableRef) bool { return r.owner == o && (o.schema != "" || r.ownerPath == s.pathVersion) }) {
			continue
		}
		if columns == nil || r.columns == nil && primary || r.columns != nil && sameColumns(r.columns, columns) {
			return true
		}
	}
	return false
}

// newlyReferencedColumn reports whether a foreign key the change created
// may reference a key of the table that has the column; keyColumns are the
// table's primary key columns, nil when it has none or they cannot be read.
func (s *scan) newlyReferencedColumn(t tableRef, column string, keyColumns []string, hasKey bool) bool {
	for _, r := range s.newReferences {
		if r.retired || r.table != t.table || r.schema != "" && r.schema != t.schema {
			continue
		}
		if r.columns != nil && slices.Contains(r.columns, column) || r.columns == nil && hasKey && (keyColumns == nil || slices.Contains(keyColumns, column)) {
			return true
		}
	}
	return false
}

// newDependency is an object the change created that depends on a
// relation: the relation, in the schema it resolves to, the one written,
// or "" when neither is known, and the object as written, with the search
// path it was written under.
type newDependency struct {
	relation tableRef
	object   tableRef
	path     int
	retired  bool
	// signature is a new function's argument types, as written.
	signature string
	// home is the synced schema a new view was made in, "" when unknown,
	// and made numbers the touch that made it.
	home string
	made int
}

// reads records the relations a new view's query names. An unqualified
// name of a CTE visible where it is used is the CTE, not a relation: a
// CTE is visible to the query of its WITH clause and to the CTEs after it,
// and, under WITH RECURSIVE, to every CTE of the clause.
func (s *scan) reads(view *ast.RangeVar, query ast.Node) {
	first := len(s.newReads)
	s.dependsIn(&s.newReads, tableRef{view.Schemaname, view.Relname}, query, nil)
	home, _ := s.creationSchema(view)
	for i := first; i < len(s.newReads); i++ {
		s.newReads[i].home = home
		s.newReads[i].made = s.lastTouch[view.Relname]
	}
}

// dropNewReaders records the views the change created that a cascading
// drop of a relation certainly takes: those whose query named it, resolved
// to it. A view whose name no statement used since it was made is freed in
// the schema it was made in, and the views reading it go too.
func (s *scan) dropNewReaders(t tableRef) {
	s.dropReturners(t)
	for i := range s.newReads {
		d := s.newReads[i]
		if d.retired || d.relation != t {
			continue
		}
		for j := range s.newReads {
			if o := &s.newReads[j]; o.object == d.object && o.path == d.path && o.made == d.made {
				o.retired = true
			}
		}
		if d.home != "" && s.lastTouch[d.object.table] == d.made {
			view := tableRef{d.home, d.object.table}
			s.free(view.schema, view.table)
			s.dropNewReaders(view)
		}
	}
}

// dependsIn records, as dependencies of object, the relations a query or
// SQL-standard routine body names, CTE names aside.
func (s *scan) dependsIn(deps *[]newDependency, object tableRef, n ast.Node, visible map[string]bool) {
	relationsIn(n, visible, func(rv *ast.RangeVar) {
		*deps = append(*deps, newDependency{
			relation: tableRef{s.schemaOf(rv), rv.Relname},
			object:   object,
			path:     s.pathVersion,
		})
	})
}

// readsMissing reports whether a query names a relation the target
// certainly lacks, or one a query cannot read (an index or a composite
// type), CTE names aside.
func (s *scan) readsMissing(n ast.Node) bool {
	if s.index == nil {
		return false
	}
	missing := false
	relationsIn(n, nil, func(rv *ast.RangeVar) {
		if s.missing(rv) {
			missing = true
		} else if _, kind, ok := s.lookup(rv); ok && (kind == kindIndex || kind == kindCompositeType) {
			missing = true
		}
	})
	return missing
}

// relationsIn calls fn for each relation a query or SQL-standard routine
// body names. An unqualified name of a CTE visible where it is used is the
// CTE, not a relation.
func relationsIn(n ast.Node, visible map[string]bool, fn func(*ast.RangeVar)) {
	ast.Inspect(n, func(n ast.Node) bool {
		if w, rest := splitWith(n); w != nil {
			scope := make(map[string]bool, len(visible))
			for name := range visible {
				scope[name] = true
			}
			var ctes []*ast.CommonTableExpr
			if w.Ctes != nil {
				for _, item := range w.Ctes.Items {
					if c, ok := item.(*ast.CommonTableExpr); ok {
						ctes = append(ctes, c)
					}
				}
			}
			if w.Recursive {
				for _, c := range ctes {
					scope[c.Ctename] = true
				}
			}
			for _, c := range ctes {
				relationsIn(c.Ctequery, scope, fn)
				scope[c.Ctename] = true
			}
			relationsIn(rest, scope, fn)
			return false
		}
		if rv, ok := n.(*ast.RangeVar); ok && !(rv.Schemaname == "" && visible[rv.Relname]) {
			fn(rv)
		}
		return true
	})
}

// splitWith returns a statement's WITH clause and the statement without
// it, or nil when it has none.
func splitWith(n ast.Node) (*ast.WithClause, ast.Node) {
	switch v := n.(type) {
	case *ast.SelectStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	case *ast.InsertStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	case *ast.UpdateStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	case *ast.DeleteStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	case *ast.MergeStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	}
	return nil, nil
}

// returns records the relations whose row type a new function returns or
// takes, as far as its result and parameter types name one.
func (s *scan) returns(v *ast.CreateFunctionStmt) {
	fn := nameParts(v.Funcname)
	if len(fn) == 0 {
		return
	}
	object := tableRef{table: fn[len(fn)-1]}
	if len(fn) >= 2 {
		object.schema = fn[len(fn)-2]
	}
	if v.IsOrReplace {
		s.replaces(v, object)
	}
	s.madeRoutines[object.table] = true
	// The result type and every parameter's type may be a row type.
	types := []*ast.TypeName{v.ReturnType}
	var args []ast.Node
	if v.Parameters != nil {
		for _, item := range v.Parameters.Items {
			if p, ok := item.(*ast.FunctionParameter); ok {
				types = append(types, p.ArgType)
				// The signature DROP FUNCTION names is the input arguments.
				if p.Mode != ast.FUNC_PARAM_OUT && p.Mode != ast.FUNC_PARAM_TABLE && p.ArgType != nil {
					args = append(args, p.ArgType)
				}
			}
		}
	}
	signature := argSignature(&ast.List{Items: args})
	// OR REPLACE of a signature the change created redefines that function.
	if !slices.Contains(s.newSignatures[object], signature) {
		s.newFunctions[object]++
		s.newSignatures[object] = append(s.newSignatures[object], signature)
	}
	first := len(s.newReturns)
	defer func() {
		for i := first; i < len(s.newReturns); i++ {
			s.newReturns[i].signature = signature
		}
	}()
	// A SQL-standard body is parsed at creation, and the server records
	// the relations it names; a string body records nothing.
	if v.SqlBody != nil {
		s.dependsIn(&s.newReturns, object, v.SqlBody, nil)
	}
	for _, t := range types {
		if t == nil {
			continue
		}
		parts := nameParts(t.Names)
		if len(parts) == 0 || len(parts) > 2 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) == 2 {
			rv.Schemaname = parts[0]
		}
		s.newReturns = append(s.newReturns, newDependency{relation: tableRef{s.schemaOf(rv), rv.Relname}, object: object, path: s.pathVersion})
	}
}

// routineSignature returns the routine a CREATE FUNCTION names, as
// written, and the argument signature DROP FUNCTION matches it by.
func routineSignature(v *ast.CreateFunctionStmt) (tableRef, string, bool) {
	fn := nameParts(v.Funcname)
	if len(fn) == 0 {
		return tableRef{}, "", false
	}
	object := tableRef{table: fn[len(fn)-1]}
	if len(fn) >= 2 {
		object.schema = fn[len(fn)-2]
	}
	var args []ast.Node
	if v.Parameters != nil {
		for _, item := range v.Parameters.Items {
			if p, ok := item.(*ast.FunctionParameter); ok && p.Mode != ast.FUNC_PARAM_OUT && p.Mode != ast.FUNC_PARAM_TABLE && p.ArgType != nil {
				args = append(args, p.ArgType)
			}
		}
	}
	return object, argSignature(&ast.List{Items: args}), true
}

// signatureTaken reports whether CREATE FUNCTION names a signature that
// certainly exists: one the change created under the name as written, or,
// without arguments, a synced routine's of that name where the statement
// creates it, which no statement dropped or renamed.
func (s *scan) signatureTaken(v *ast.CreateFunctionStmt) bool {
	object, signature, ok := routineSignature(v)
	if !ok {
		return false
	}
	if slices.Contains(s.newSignatures[object], signature) {
		return true
	}
	if s.index == nil || s.movedRoutines[object.table] {
		return false
	}
	schema, ok := s.creationSchema(&ast.RangeVar{Schemaname: object.schema, Relname: object.table})
	fn := tableRef{schema, object.table}
	if !ok || s.droppedFunctions[fn] {
		return false
	}
	return signature == "" && s.index.noArgs[fn] || slices.Contains(s.index.signatures[fn], signature)
}

// replaces records a CREATE OR REPLACE FUNCTION that certainly redefines
// a synced function: the only one of its name where the statement creates
// it, which, like the statement, takes no arguments, and no routine of
// the change has had the name. The synced function's dependencies are
// then the statement's.
func (s *scan) replaces(v *ast.CreateFunctionStmt, object tableRef) {
	if s.index == nil || s.madeRoutines[object.table] || s.movedRoutines[object.table] {
		return
	}
	if v.Parameters != nil {
		for _, item := range v.Parameters.Items {
			if p, ok := item.(*ast.FunctionParameter); !ok || p.Mode != ast.FUNC_PARAM_OUT && p.Mode != ast.FUNC_PARAM_TABLE {
				return
			}
		}
	}
	schema, ok := s.creationSchema(&ast.RangeVar{Schemaname: object.schema, Relname: object.table})
	fn := tableRef{schema, object.table}
	if ok && s.index.functions[fn] == 1 && s.index.noArgs[fn] && !s.droppedFunctions[fn] {
		s.replacedFunctions[fn] = true
	}
}

// dependsOnNew reports whether an object the change created and has not
// dropped, other than the ones named in except, may depend on the
// relation.
func (s *scan) dependsOnNew(deps []newDependency, t tableRef, except []tableRef) bool {
	for _, d := range deps {
		if d.retired || d.relation.table != t.table || d.relation.schema != "" && d.relation.schema != t.schema {
			continue
		}
		if slices.ContainsFunc(except, func(o tableRef) bool { return s.sameObject(d, o) }) {
			continue
		}
		return true
	}
	return false
}

// sameObject reports whether a name, as written, is the object a
// dependency belongs to: written alike, under the same search path when
// unqualified.
func (s *scan) sameObject(d newDependency, name tableRef) bool {
	return d.object == name && (name.schema != "" || d.path == s.pathVersion)
}

// retarget renames the object of the dependencies a name, as written,
// refers to.
func (s *scan) retarget(deps []newDependency, old tableRef, name string) {
	for i := range deps {
		if s.sameObject(deps[i], old) {
			deps[i].object.table = name
		}
	}
}

// retireSignature retires the dependencies of a new function written with
// that argument signature.
func (s *scan) retireSignature(fn tableRef, sig string) {
	for i := range s.newReturns {
		if s.sameObject(s.newReturns[i], fn) && s.newReturns[i].signature == sig {
			s.newReturns[i].retired = true
		}
	}
}

// argSignature writes the argument types of a routine as written, one per
// argument, for comparing a DROP's argument list with a CREATE's.
func argSignature(args *ast.List) string {
	var types []string
	if args != nil {
		for _, item := range args.Items {
			if t, ok := item.(*ast.TypeName); ok {
				// A built-in type is pg_catalog's, however it is written.
				parts := nameParts(t.Names)
				if len(parts) == 2 && parts[0] == "pg_catalog" {
					parts = parts[1:]
				}
				name := strings.Join(parts, ".")
				if t.PctType {
					name += "%TYPE"
				}
				if t.ArrayBounds != nil && len(t.ArrayBounds.Items) > 0 {
					name += "[]"
				}
				types = append(types, name)
			} else {
				types = append(types, "?")
			}
		}
	}
	return strings.Join(types, ",")
}

// retire records that a statement dropped an object the change created,
// named as it was written, under the same search path when unqualified.
func (s *scan) retire(deps []newDependency, object tableRef) {
	for i := range deps {
		if s.sameObject(deps[i], object) {
			deps[i].retired = true
		}
	}
}

// create follows CREATE TABLE. A CREATE TABLE of a name its schema
// already holds is refused.
func (s *scan) create(st *statement, v *ast.CreateStmt) {
	exists := s.existsWhereCreated(v.Relation)
	switch {
	case v.Relation != nil && v.Relation.Catalogname != "" && (s.index == nil || v.Relation.Catalogname != s.index.database):
		// A table in another database is refused; without a schema the
		// database name is not known, and nothing is followed.
		if s.index != nil {
			s.stop()
		}
		return
	case s.schemaMissing(v.Relation) || exists && !v.IfNotExists || s.typeExists(v.Relation):
		// A type of the name is refused even under IF NOT EXISTS.
		s.stop()
		return
	case exists:
		// IF NOT EXISTS of an existing table does nothing.
		return
	}
	if duplicateNames(v) || s.sourceRefused(v) || s.keyNameTaken(v) || checksRefused(v) {
		s.stop()
		return
	}
	// A name the change created and still has is taken.
	if s.madeHere(v.Relation) {
		if !v.IfNotExists {
			s.stop()
		}
		return
	}
	refused, uncertain := s.referencesRefused(v)
	if refused {
		s.stop()
		return
	}
	pending := len(s.pending)
	// A table that may not be created is not reported.
	if !uncertain && !s.nameUncertain(v.Relation) {
		s.createTable(st, v)
	}
	s.touch(v.Relation)
	for _, p := range s.pending[pending:] {
		p.made = s.lastTouch[v.Relation.Relname]
	}
	if !v.IfNotExists || s.createsNew(v.Relation) {
		s.made(v.Relation, kindTable)
		if def, _, complete := defOf(v); complete {
			if c, ok := s.madeAt(v.Relation); ok {
				c.def = &def
			}
		}
	}
	if generatesNames(v) {
		s.generate(v.Relation)
	}
	if v.TableElts == nil {
		return
	}
	for _, item := range v.TableElts.Items {
		var constraints []*ast.Constraint
		switch e := item.(type) {
		case *ast.ColumnDef:
			constraints = constraintsOf(e.Constraints)
		case *ast.Constraint:
			constraints = []*ast.Constraint{e}
		}
		s.references(v.Relation, constraints...)
		for _, c := range constraints {
			// A key's index takes the constraint's name.
			if makesIndex(c) && c.Conname != "" {
				s.touchName(v.Relation.Schemaname, c.Conname)
			}
		}
	}
}

// createSchema follows CREATE SCHEMA and the statements it holds, each in
// the new schema unless it names another. A schema named after its owner
// takes the owner's name.
func (s *scan) createSchema(st *statement, v *ast.CreateSchemaStmt) {
	name := v.Schemaname
	if name == "" && v.Authrole != nil {
		name = v.Authrole.Rolename
	}
	if name == "" {
		if v.SchemaElts != nil && len(v.SchemaElts.Items) > 0 {
			s.stop()
		}
		return
	}
	// IF NOT EXISTS cannot carry schema elements: the server refuses it.
	if v.IfNotExists && v.SchemaElts != nil && len(v.SchemaElts.Items) > 0 {
		s.stop()
		return
	}
	// CREATE SCHEMA of a schema the target has, and the change did not
	// drop, or one the change created and did not drop, is refused, and
	// with IF NOT EXISTS does nothing.
	if s.index != nil && (s.index.schemas[name] != nil && !s.schemas[name] || s.schemas[name] && !s.schemaGone[name]) {
		if !v.IfNotExists {
			s.stop()
		}
		return
	}
	s.schemas[name] = true
	delete(s.schemaGone, name)
	if v.SchemaElts == nil {
		return
	}
	if elementsCollide(name, v.SchemaElts) {
		s.stop()
		return
	}
	in := func(rv *ast.RangeVar) *ast.RangeVar {
		if rv == nil || rv.Schemaname != "" {
			return rv
		}
		c := *rv
		c.Schemaname = name
		return &c
	}
	for _, elt := range v.SchemaElts.Items {
		switch e := elt.(type) {
		case *ast.CreateStmt:
			c := *e
			c.Relation = in(e.Relation)
			s.create(st, &c)
		case *ast.ViewStmt:
			if s.readsMissing(e.Query) {
				s.stop()
				return
			}
			s.touch(in(e.View))
			s.made(in(e.View), kindView)
			s.reads(in(e.View), e.Query)
		case *ast.CreateSeqStmt:
			s.touch(in(e.Sequence))
			s.made(in(e.Sequence), kindSequence)
		case *ast.IndexStmt:
			c := *e
			c.Relation = in(e.Relation)
			if s.indexRefused(&c) {
				s.stop()
				return
			}
			if e.Idxname != "" && e.Relation != nil {
				s.touchName(in(e.Relation).Schemaname, e.Idxname)
			}
		case *ast.CreateTrigStmt, *ast.GrantStmt:
		default:
			s.stop()
			return
		}
	}
}

// elementsCollide reports whether the server may refuse CREATE SCHEMA for
// its elements, which run as one statement: an element in another schema,
// two elements of one name, or an element named like a name the server
// generates (for a table's key, serial column, or unnamed index) while an
// element may generate one.
//
// pg: src/backend/parser/parse_utilcmd.c — transformCreateSchemaStmtElements
func elementsCollide(schema string, elts *ast.List) bool {
	names := make(map[string]bool)
	generates := false
	for _, elt := range elts.Items {
		var rv *ast.RangeVar
		name := ""
		switch e := elt.(type) {
		case *ast.CreateStmt:
			rv = e.Relation
			generates = generates || generatesNames(e)
		case *ast.ViewStmt:
			rv = e.View
		case *ast.CreateSeqStmt:
			rv = e.Sequence
		case *ast.IndexStmt:
			rv, name = e.Relation, e.Idxname
			if name == "" {
				generates = true
			}
		}
		if rv == nil {
			continue
		}
		if rv.Schemaname != "" && rv.Schemaname != schema {
			return true
		}
		if _, ok := elt.(*ast.IndexStmt); !ok {
			name = rv.Relname
		}
		if name == "" {
			continue
		}
		if names[name] {
			return true
		}
		names[name] = true
	}
	if generates {
		for name := range names {
			if generatedName(name) {
				return true
			}
		}
	}
	return false
}

// touchType records the name of a type a statement creates, and the
// name with suffix, when there is one, for a type it creates with it. A
// name a relation's row type or another type certainly has is refused.
func (s *scan) touchType(name *ast.List, suffix string) {
	parts := nameParts(name)
	if len(parts) == 0 || len(parts) > 2 {
		return
	}
	schema := ""
	if len(parts) == 2 {
		schema = parts[0]
	}
	rv := &ast.RangeVar{Schemaname: schema, Relname: parts[len(parts)-1]}
	if s.typeTaken(rv) {
		s.stop()
		return
	}
	s.touchName(schema, parts[len(parts)-1])
	s.made(rv, kindType)
	if suffix != "" {
		s.touchName(schema, parts[len(parts)-1]+suffix)
	}
}

// typeTaken reports whether the schema a CREATE puts a type in certainly
// holds a type of its name: a relation's row type, an enum, or what the
// change created there.
func (s *scan) typeTaken(rv *ast.RangeVar) bool {
	return s.existsWhereCreated(rv) || s.typeExists(rv) || s.madeHere(rv)
}

// indexRefused reports whether the server certainly refuses CREATE
// INDEX: its relation is one the target lacks, or cannot be indexed, or,
// without IF NOT EXISTS, its name is taken in the relation's schema.
func (s *scan) indexRefused(v *ast.IndexStmt) bool {
	if s.index == nil || v.Relation == nil {
		return false
	}
	if s.missing(v.Relation) {
		return true
	}
	schema, kind, ok := s.lookup(v.Relation)
	if !ok {
		return false
	}
	switch kind {
	case kindView, kindSequence, kindIndex, kindCompositeType, kindForeignTable:
		return true
	}
	// A plain column the index names must be the table's.
	if t, known := s.table(v.Relation); known {
		for _, list := range []*ast.List{v.IndexParams, v.IndexIncludingParams} {
			if list == nil {
				continue
			}
			for _, item := range list.Items {
				if e, ok := item.(*ast.IndexElem); ok && e.Name != "" && !s.hasColumn(t, e.Name) {
					return true
				}
			}
		}
	}
	return v.Idxname != "" && !v.IfNotExists && (s.nameTaken(schema, v.Idxname) || s.madeHere(&ast.RangeVar{Schemaname: schema, Relname: v.Idxname}))
}

// typeExists reports whether the schema a CREATE puts a relation in
// certainly holds an enum type of its name, which the relation's row
// type would collide with.
func (s *scan) typeExists(rv *ast.RangeVar) bool {
	schema, ok := s.creationSchema(rv)
	if !ok || s.isTouched(&ast.RangeVar{Schemaname: schema, Relname: rv.Relname}) {
		return false
	}
	return slices.Contains(s.index.schemas[schema].types, rv.Relname)
}

// mayHavePartitions reports whether a relation may be a partitioned table
// with partitions, whose indexes a key or index added to it gives
// generated names: any but a synced table without partitions, or a
// relation that cannot have any.
func (s *scan) mayHavePartitions(rv *ast.RangeVar) bool {
	schema, kind, ok := s.lookup(rv)
	switch {
	case !ok || kind == kindPartition || kind == kindAmbiguous:
		return true
	case kind != kindTable:
		return false
	case s.isUnsettled(tableRef{schema, rv.Relname}):
		return true
	}
	return len(s.index.schemas[schema].tables[rv.Relname].GetPartitions()) > 0
}

// tableDef is what a CREATE TABLE of the change gives its table, when the
// statement says it all: its columns and its keys, primary or unique.
type tableDef struct {
	columns map[string]bool
	keys    [][]string
	primary bool
}

// defOf reads a CREATE TABLE's columns and keys, and its foreign keys,
// each with its column when declared on one. complete is false when LIKE,
// INHERITS, OF, or PARTITION OF brings more.
func defOf(v *ast.CreateStmt) (def tableDef, fks []foreignKeyDecl, complete bool) {
	def.columns = make(map[string]bool)
	complete = (v.InhRelations == nil || len(v.InhRelations.Items) == 0) && v.OfTypename == nil && v.Partbound == nil
	if v.TableElts == nil {
		return def, nil, complete
	}
	add := func(c *ast.Constraint, column string) {
		cols := nameParts(c.Keys)
		if column != "" {
			cols = []string{column}
		}
		switch c.Contype {
		case ast.CONSTR_PRIMARY:
			def.primary = true
			def.keys = append(def.keys, cols)
		case ast.CONSTR_UNIQUE:
			def.keys = append(def.keys, cols)
		case ast.CONSTR_FOREIGN:
			fks = append(fks, foreignKeyDecl{c, column})
		}
	}
	for _, item := range v.TableElts.Items {
		switch e := item.(type) {
		case *ast.ColumnDef:
			def.columns[e.Colname] = true
			for _, c := range constraintsOf(e.Constraints) {
				add(c, e.Colname)
			}
		case *ast.Constraint:
			add(e, "")
		case *ast.TableLikeClause:
			complete = false
		}
	}
	return def, fks, complete
}

// foreignKeyDecl is a foreign key CREATE TABLE declares, with its column
// when declared on one.
type foreignKeyDecl struct {
	c      *ast.Constraint
	column string
}

// refusedBy reports whether a table of that definition has no key a
// foreign key of local referencing columns can reference with these
// columns: no primary key without columns, or none of them, or no key on
// exactly them.
func (d *tableDef) refusedBy(columns []string, local int) bool {
	switch {
	case len(columns) > 0 && len(columns) != local,
		slices.ContainsFunc(columns, func(c string) bool { return !d.columns[c] }),
		len(columns) == 0 && !d.primary,
		len(columns) > 0 && !slices.ContainsFunc(d.keys, func(k []string) bool { return sameKey(k, columns) }):
		return true
	}
	return false
}

// referencesRefused reports whether the server certainly refuses CREATE
// TABLE for a foreign key it declares: on a column the table lacks, or to
// a relation, column, or key as refusesReference finds, or to the table
// itself, or to one the change created and has not altered, without the
// columns or key it names. uncertain reports a foreign key the scan cannot
// check.
func (s *scan) referencesRefused(v *ast.CreateStmt) (refused, uncertain bool) {
	if s.index == nil {
		return false, false
	}
	own, fks, complete := defOf(v)
	for _, fk := range fks {
		if hasDuplicate(nameParts(fk.c.FkAttrs)) {
			return true, false
		}
		if fk.column == "" && complete && slices.ContainsFunc(nameParts(fk.c.FkAttrs), func(c string) bool { return !own.columns[c] }) {
			return true, false
		}
		pk := fk.c.Pktable
		columns := nameParts(fk.c.PkAttrs)
		local := max(len(nameParts(fk.c.FkAttrs)), 1)
		self := pk != nil && pk.Relname == v.Relation.Relname && pk.Schemaname == v.Relation.Schemaname
		ambiguous := false
		if pk != nil && pk.Relname == v.Relation.Relname && !self {
			// Spelled otherwise, it is the table itself when both names
			// resolve to the schema CREATE puts the table in, and another
			// table when they certainly resolve to different ones.
			a, aok := s.homeOf(pk)
			b, bok := s.homeOf(v.Relation)
			switch {
			case aok && bok && a == b:
				self = true
			case aok && bok:
				pk = &ast.RangeVar{Schemaname: a, Relname: pk.Relname}
			default:
				ambiguous = true
			}
		}
		switch {
		case pk == nil || ambiguous:
			// It may be the table itself, or another of its name.
			uncertain = true
		case self:
			if !complete {
				uncertain = true
			} else if own.refusedBy(columns, local) {
				return true, false
			}
		default:
			if def, ok := s.madeDef(pk); ok {
				if def.refusedBy(columns, local) {
					return true, false
				}
				continue
			}
			r, u := s.refusesReference(tableRef{}, fk.c, nil, nil)
			if r {
				return true, false
			}
			uncertain = uncertain || u
		}
	}
	return false, uncertain
}

// homeOf returns the schema a name means for a table CREATE TABLE makes:
// the one written, or the one the search path creates in.
func (s *scan) homeOf(rv *ast.RangeVar) (string, bool) {
	if rv.Schemaname != "" {
		return rv.Schemaname, true
	}
	return s.creationSchema(rv)
}

// checksRefused reports whether a constraint CREATE TABLE declares is
// certainly refused: a CHECK or a key naming a column the table does not
// have, when the statement says all its columns, a key naming a column
// twice, or a second primary key.
func checksRefused(v *ast.CreateStmt) bool {
	own, _, complete := defOf(v)
	if v.TableElts == nil {
		return false
	}
	primaries := 0
	for _, item := range v.TableElts.Items {
		var cs []*ast.Constraint
		switch e := item.(type) {
		case *ast.ColumnDef:
			cs = constraintsOf(e.Constraints)
		case *ast.Constraint:
			cs = []*ast.Constraint{e}
		}
		for _, c := range cs {
			missing := func(column string) bool { return complete && !own.columns[column] }
			switch c.Contype {
			case ast.CONSTR_CHECK:
				if slices.ContainsFunc(columnsIn(c.RawExpr), missing) {
					return true
				}
			case ast.CONSTR_PRIMARY, ast.CONSTR_UNIQUE:
				// A key names columns the table has, each once, and a table
				// has one primary key.
				keys := nameParts(c.Keys)
				if hasDuplicate(keys) || slices.ContainsFunc(keys, missing) {
					return true
				}
				if c.Contype == ast.CONSTR_PRIMARY {
					primaries++
				}
			}
		}
	}
	return primaries > 1
}

// madeDef returns the definition of a table a name certainly means that
// the change created, said all of, and has not altered.
func (s *scan) madeDef(rv *ast.RangeVar) (*tableDef, bool) {
	c, ok := s.madeAt(rv)
	if !ok || c.kind != kindTable || c.def == nil {
		return nil, false
	}
	return c.def, true
}

// madeHere reports whether a CREATE puts a relation under a name the
// change created a table or view under, and no statement used since.
func (s *scan) madeHere(rv *ast.RangeVar) bool {
	_, ok := s.madeAt(rv)
	return ok
}

// madeIn returns the schema of a relation the change created that a name
// certainly means.
func (s *scan) madeIn(rv *ast.RangeVar) (string, bool) {
	if _, ok := s.madeAt(rv); !ok {
		return "", false
	}
	if rv.Schemaname != "" {
		return rv.Schemaname, true
	}
	return s.creationSchema(rv)
}

// madeAt returns what the change created under a name, in the schema
// written or the one the search path creates in, while no statement used
// the name since.
func (s *scan) madeAt(rv *ast.RangeVar) (*madeRelation, bool) {
	schema := rv.Schemaname
	if schema == "" {
		var ok bool
		if schema, ok = s.creationSchema(rv); !ok {
			return nil, false
		}
	}
	c, ok := s.created[[2]string{schema, rv.Relname}]
	if !ok || s.lastTouch[rv.Relname] != c.made {
		return nil, false
	}
	return c, true
}

// sourceRefused reports whether CREATE TABLE names a relation the server
// refuses to take columns from: an INHERITS parent or a LIKE source the
// target certainly lacks, a parent that is not a table, or a LIKE source
// that is a sequence or an index.
func (s *scan) sourceRefused(v *ast.CreateStmt) bool {
	if s.index == nil {
		return false
	}
	var parents, sources []*ast.RangeVar
	if v.InhRelations != nil {
		for _, item := range v.InhRelations.Items {
			if rv, ok := item.(*ast.RangeVar); ok {
				parents = append(parents, rv)
			}
		}
	}
	if v.TableElts != nil {
		for _, item := range v.TableElts.Items {
			if like, ok := item.(*ast.TableLikeClause); ok && like.Relation != nil {
				sources = append(sources, like.Relation)
			}
		}
	}
	for _, rv := range append(parents, sources...) {
		if s.missing(rv) {
			return true
		}
	}
	// OF names a composite type, which the snapshot lists as a relation.
	if v.OfTypename != nil {
		parts := nameParts(v.OfTypename.Names)
		if len(parts) == 1 || len(parts) == 2 {
			rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
			if len(parts) == 2 {
				rv.Schemaname = parts[0]
			}
			if _, kind, ok := s.lookup(rv); ok && kind != kindCompositeType || !ok && s.missing(rv) {
				return true
			}
		}
	}
	for _, rv := range parents {
		if _, kind, ok := s.lookup(rv); ok && (kind == kindView || kind == kindMatView || kind == kindSequence || kind == kindIndex || kind == kindCompositeType) {
			return true
		}
	}
	for _, rv := range sources {
		if _, kind, ok := s.lookup(rv); ok && (kind == kindSequence || kind == kindIndex) {
			return true
		}
	}
	return false
}

// duplicateNames reports whether CREATE TABLE names a column, or a
// constraint, twice, which the server refuses. A column LIKE or INHERITS
// brings may merge with one the statement names.
func duplicateNames(v *ast.CreateStmt) bool {
	if v.TableElts == nil {
		return false
	}
	columns := make(map[string]bool)
	constraints := make(map[string]bool)
	named := func(cs []*ast.Constraint) bool {
		for _, c := range cs {
			if c.Conname == "" {
				continue
			}
			if constraints[c.Conname] {
				return true
			}
			constraints[c.Conname] = true
		}
		return false
	}
	for _, item := range v.TableElts.Items {
		switch e := item.(type) {
		case *ast.ColumnDef:
			if columns[e.Colname] || named(constraintsOf(e.Constraints)) {
				return true
			}
			columns[e.Colname] = true
		case *ast.Constraint:
			if named([]*ast.Constraint{e}) {
				return true
			}
		}
	}
	return false
}

// keyNameTaken reports whether a key CREATE TABLE names takes for its
// index a name its schema certainly holds: the table's own, or a
// relation's the synced schema or the change has.
func (s *scan) keyNameTaken(v *ast.CreateStmt) bool {
	if s.index == nil || v.TableElts == nil || v.Relation == nil {
		return false
	}
	schema := v.Relation.Schemaname
	if schema == "" {
		schema, _ = s.creationSchema(v.Relation)
	}
	for _, item := range v.TableElts.Items {
		var cs []*ast.Constraint
		switch e := item.(type) {
		case *ast.ColumnDef:
			cs = constraintsOf(e.Constraints)
		case *ast.Constraint:
			cs = []*ast.Constraint{e}
		}
		for _, c := range cs {
			if !makesIndex(c) || c.Conname == "" {
				continue
			}
			if c.Conname == v.Relation.Relname || schema != "" && (s.nameTaken(schema, c.Conname) || s.madeHere(&ast.RangeVar{Schemaname: schema, Relname: c.Conname})) {
				return true
			}
		}
	}
	return false
}

// generatesNames reports whether CREATE TABLE makes a relation under a
// name the server generates: the index of an unnamed key, the sequence
// of a serial or identity column, what LIKE copies, or the indexes a
// partition takes from its parent.
func generatesNames(v *ast.CreateStmt) bool {
	if v.Partbound != nil {
		return true
	}
	if v.TableElts == nil {
		return false
	}
	for _, item := range v.TableElts.Items {
		switch e := item.(type) {
		case *ast.ColumnDef:
			if columnGeneratesNames(e) {
				return true
			}
		case *ast.Constraint:
			if makesIndex(e) && e.Conname == "" {
				return true
			}
		case *ast.TableLikeClause:
			return true
		}
	}
	return false
}

// columnGeneratesNames reports whether a column makes a relation under a
// generated name: a serial or identity column's sequence, or the index of
// an unnamed key on it.
func columnGeneratesNames(cd *ast.ColumnDef) bool {
	if cd.TypeName != nil {
		if parts := nameParts(cd.TypeName.Names); len(parts) == 1 {
			switch parts[0] {
			case "smallserial", "serial", "bigserial", "serial2", "serial4", "serial8":
				return true
			}
		}
	}
	for _, c := range constraintsOf(cd.Constraints) {
		if makesIndex(c) && c.Conname == "" || c.Contype == ast.CONSTR_IDENTITY {
			return true
		}
	}
	return false
}

// generate records that a statement may have made relations under
// generated names in the schema of a relation: the one the name resolves
// to or names, or the one CREATE puts it in, or any schema when rv is nil
// or the schema is not known.
func (s *scan) generate(rv *ast.RangeVar) {
	schema := ""
	if rv != nil {
		if schema = s.schemaOf(rv); schema == "" {
			schema, _ = s.creationSchema(rv)
		}
	}
	s.generated[schema] = true
}

// mayBeGenerated reports whether a CREATE's name may be one a statement
// generated for a relation of the schema the CREATE puts it in.
func (s *scan) mayBeGenerated(rv *ast.RangeVar) bool {
	if !generatedName(rv.Relname) || len(s.generated) == 0 {
		return false
	}
	schema := rv.Schemaname
	if schema == "" {
		var ok bool
		if schema, ok = s.creationSchema(rv); !ok {
			return true
		}
	}
	return s.generated[""] || s.generated[schema]
}

// touch records a relation name a statement creates.
func (s *scan) touch(rv *ast.RangeVar) {
	if rv != nil {
		s.touchName(rv.Schemaname, rv.Relname)
	}
}

// touchName records a relation name a statement changed, in schema, or
// in whichever schema the search path finds when schema is "".
func (s *scan) touchName(schema, name string) {
	s.touched[[2]string{schema, name}] = true
	s.touchedName[name] = true
	s.touches++
	s.lastTouch[name] = s.touches
	for key := range s.freed {
		if key[1] == name && (schema == "" || key[0] == schema) {
			delete(s.freed, key)
		}
	}
}

// free records a name a resolved drop removed from its schema.
func (s *scan) free(schema, name string) {
	s.touchName(schema, name)
	s.freed[[2]string{schema, name}] = true
}

// newKey records a table a statement may have given a key or a unique
// index.
func (s *scan) newKey(rv *ast.RangeVar) {
	s.newKeys[[2]string{s.schemaOf(rv), rv.Relname}] = true
}

// unsettle records a table whose columns or constraints changed in a way
// the scan does not follow.
func (s *scan) unsettle(rv *ast.RangeVar) {
	s.unsettled[[2]string{s.schemaOf(rv), rv.Relname}] = true
}

// isUnsettled reports whether a statement may have changed the table's
// columns or constraints.
func (s *scan) isUnsettled(t tableRef) bool {
	return s.unsettled[[2]string{t.schema, t.table}] || s.unsettled[[2]string{"", t.table}]
}

// columnNow reports whether a statement of the change left the table's
// column existing, and known whether one recorded anything for it.
func (s *scan) columnNow(t tableRef, column string) (exists, known bool) {
	if v, ok := s.columns[[3]string{t.schema, t.table, column}]; ok {
		return v, true
	}
	v, ok := s.columns[[3]string{"", t.table, column}]
	return v, ok
}

// originOf returns what a known table's column, by its name now, is in the
// synced table: itself when no statement renamed it, the column it was
// renamed from, or "" for a column the change added. ok is false when the
// scan cannot tell.
func (s *scan) originOf(t tableRef, column string) (origin string, ok bool) {
	if s.columnRenamed(t, column) {
		origin, ok = s.origins[[3]string{t.schema, t.table, column}]
		return origin, ok
	}
	if _, changed := s.columnNow(t, column); changed {
		return "", true
	}
	return column, true
}

// columnRenamed reports whether the table's column has its name from a
// rename.
func (s *scan) columnRenamed(t tableRef, column string) bool {
	return s.renamedColumns[[3]string{t.schema, t.table, column}] || s.renamedColumns[[3]string{"", t.table, column}]
}

// setColumn records that a column exists or is gone after a statement.
func (s *scan) setColumn(schema, table, column string, exists, renamed bool) {
	key := [3]string{schema, table, column}
	s.columns[key] = exists
	if renamed {
		s.renamedColumns[key] = true
	} else {
		delete(s.renamedColumns, key)
	}
}

// isTouched reports whether a statement may have changed what the name
// means: a qualified name when a statement changed it in that schema or
// in an unknown one, an unqualified name when a statement changed it in
// any schema, which may come first on the search path.
func (s *scan) isTouched(rv *ast.RangeVar) bool {
	if rv.Schemaname == "" {
		return s.touchedName[rv.Relname]
	}
	return s.touched[[2]string{rv.Schemaname, rv.Relname}] || s.touched[[2]string{"", rv.Relname}]
}

// lookup resolves a relation name against the synced schema. ok is false
// when the name may no longer mean what the schema says, or names nothing
// the schema holds. An unqualified name needs a known search path, and a
// name beginning with pg_ may be a system catalog, which pg_catalog finds
// first.
func (s *scan) lookup(rv *ast.RangeVar) (schema string, kind relationKind, ok bool) {
	if s.index == nil || rv == nil || s.isTouched(rv) || s.schemas[rv.Schemaname] {
		return "", 0, false
	}
	if rv.Catalogname != "" && rv.Catalogname != s.index.database {
		return "", 0, false
	}
	if rv.Schemaname != "" {
		ns := s.index.schemas[rv.Schemaname]
		if ns == nil {
			return "", 0, false
		}
		kind, ok := ns.relations[rv.Relname]
		return rv.Schemaname, kind, ok
	}
	if strings.HasPrefix(rv.Relname, "pg_") {
		return "", 0, false
	}
	path, ok := s.searchPath()
	if !ok {
		return "", 0, false
	}
	for _, name := range path {
		if ns := s.index.schemas[name]; ns != nil {
			if kind, ok := ns.relations[rv.Relname]; ok {
				return name, kind, true
			}
			continue
		}
		// A system schema the snapshot leaves out may hold the name.
		if name == "information_schema" {
			return "", 0, false
		}
	}
	return "", 0, false
}

// searchPath expands "$user" in the current search path. Without a known
// current user, "$user" names no schema only when every schema of the
// database is listed by name.
func (s *scan) searchPath() ([]string, bool) {
	if s.path == nil {
		return nil, false
	}
	path := make([]string, 0, len(s.path))
	for _, name := range s.path {
		switch {
		case name != "$user":
			path = append(path, name)
		case s.user != "":
			path = append(path, s.user)
		case s.index.unlisted(s.path):
			return nil, false
		}
	}
	return path, true
}

// schemaOf returns the schema a relation name means now: the one the
// synced schema resolves it to, the one the SQL wrote, or "" when neither
// is known.
func (s *scan) schemaOf(rv *ast.RangeVar) string {
	if schema, _, ok := s.lookup(rv); ok {
		return schema
	}
	return rv.Schemaname
}

// table resolves the table an ALTER TABLE acts on, with its constraints
// as the synced schema lists them, when the change has not altered it in
// a way the scan does not follow. A partition is not a table here: its
// keys come from its parent.
func (s *scan) table(rv *ast.RangeVar) (tableRef, bool) {
	schema, kind, ok := s.lookup(rv)
	t := tableRef{schema, rv.Relname}
	if !ok || kind != kindTable || s.isUnsettled(t) {
		return tableRef{}, false
	}
	return t, true
}

// setting follows SET and RESET. The scan follows SET search_path to a
// list of plain names and SET ROLE. A transaction-local setting reverts at
// commit, and RESET returns to the session's default, which the target
// does not record, so after those the path or the user is unknown.
func (s *scan) setting(v *ast.VariableSetStmt) {
	if v.Kind == ast.VAR_RESET_ALL {
		s.path = nil
		s.pathVersion++
		return
	}
	switch v.Name {
	case "search_path":
		s.pathVersion++
		if v.Kind != ast.VAR_SET_VALUE || v.IsLocal || !plainNames(v.Args) {
			s.path = nil
			return
		}
		s.path = nil
		for _, item := range v.Args.Items {
			name, _ := constString(item)
			s.path = append(s.path, name)
		}
	case "role":
		s.pathVersion++
		switch {
		case v.IsLocal:
			s.user = ""
		case v.Kind == ast.VAR_SET_DEFAULT || v.Kind == ast.VAR_RESET:
			s.user = s.session
		case v.Kind == ast.VAR_SET_VALUE && v.Args != nil && len(v.Args.Items) == 1:
			name, _ := constString(v.Args.Items[0])
			if name == "none" {
				name = s.session
			}
			s.user = name
		default:
			s.user = ""
		}
	case "session_authorization":
		s.pathVersion++
		s.user, s.session = "", ""
	}
}

// plainNames reports whether every SET search_path value is a name the
// setting reads back unchanged: lower-case letters, digits, underscores,
// and dollar signs. A value such as 'MySchema' or 'a, b' is read as a list
// of identifiers, which the scan does not redo.
func plainNames(args *ast.List) bool {
	if args == nil {
		return false
	}
	for _, item := range args.Items {
		name, ok := constString(item)
		if !ok || name == "" {
			return false
		}
		for i := 0; i < len(name); i++ {
			c := name[i]
			if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '$') {
				return false
			}
		}
	}
	return true
}

// data follows a statement that changes data, not the schema. A call of
// set_config on search_path changes how later names resolve: the form
// pg_dump writes, a lone SELECT of it with constant arguments, sets the
// path, and any other call leaves it unknown.
func (s *scan) data(n ast.Node) {
	calls := searchPathCalls(n)
	if len(calls) == 0 {
		return
	}
	s.pathVersion++
	s.path = nil
	if path, ok := dumpSearchPath(n, calls); ok {
		// An empty path is a known path that finds nothing.
		s.path = append([]string{}, path...)
	}
}

// searchPathCalls lists the set_config calls that may set search_path:
// every call whose setting name is not a constant naming another setting.
func searchPathCalls(n ast.Node) []*ast.FuncCall {
	var calls []*ast.FuncCall
	ast.Inspect(n, func(n ast.Node) bool {
		fc, ok := n.(*ast.FuncCall)
		if !ok {
			return true
		}
		parts := nameParts(fc.Funcname)
		builtin := len(parts) == 1 || len(parts) == 2 && parts[0] == "pg_catalog"
		if !builtin || parts[len(parts)-1] != "set_config" {
			return true
		}
		if name, ok := constString(firstArg(fc)); ok && !strings.EqualFold(name, "search_path") {
			return true
		}
		calls = append(calls, fc)
		return true
	})
	return calls
}

// dumpSearchPath reads SELECT pg_catalog.set_config('search_path', '...',
// false), the only statement in n, and returns the path it sets.
func dumpSearchPath(n ast.Node, calls []*ast.FuncCall) ([]string, bool) {
	sel, ok := n.(*ast.SelectStmt)
	if !ok || len(calls) != 1 || sel.FromClause != nil || sel.WhereClause != nil || sel.WithClause != nil ||
		sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
		return nil, false
	}
	target, ok := sel.TargetList.Items[0].(*ast.ResTarget)
	if !ok || target.Val != calls[0] || calls[0].Args == nil || len(calls[0].Args.Items) != 3 {
		return nil, false
	}
	value, ok := constString(calls[0].Args.Items[1])
	if !ok {
		return nil, false
	}
	local, ok := calls[0].Args.Items[2].(*ast.A_Const)
	if !ok || local.Isnull {
		return nil, false
	}
	if b, ok := local.Val.(*ast.Boolean); !ok || b.Boolval {
		return nil, false
	}
	return parseSearchPath(value), true
}

func firstArg(fc *ast.FuncCall) ast.Node {
	if fc.Args == nil || len(fc.Args.Items) == 0 {
		return nil
	}
	return fc.Args.Items[0]
}

// constString returns the value of a string constant.
func constString(n ast.Node) (string, bool) {
	c, ok := n.(*ast.A_Const)
	if !ok || c.Isnull {
		return "", false
	}
	str, ok := c.Val.(*ast.String)
	if !ok {
		return "", false
	}
	return str.Str, true
}

// selectsInto reports whether a branch of a set operation carries INTO.
func selectsInto(s *ast.SelectStmt) bool {
	if s == nil {
		return false
	}
	return s.IntoClause != nil || selectsInto(s.Larg) || selectsInto(s.Rarg)
}

// isTemp reports whether a relation is created in the temporary schema.
func isTemp(rv *ast.RangeVar) bool {
	return rv != nil && (rv.Relpersistence == 't' || rv.Schemaname == "pg_temp" || strings.HasPrefix(rv.Schemaname, "pg_temp_"))
}

// isRelationKind reports whether an object type lives in the relation
// namespace.
func isRelationKind(t ast.ObjectType) bool {
	switch t {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_SEQUENCE, ast.OBJECT_INDEX:
		return true
	}
	return false
}

// rename follows RENAME. A relation renamed is touched under both names,
// and a pending table it may be keeps the new name too. An index rename
// renames the constraint the index backs, if any.
func (s *scan) rename(v *ast.RenameStmt) {
	// A rename on a relation the target certainly lacks is refused, and
	// under IF EXISTS does nothing.
	if v.Relation != nil && s.index != nil && s.missing(v.Relation) {
		if !v.MissingOk {
			s.stop()
		}
		return
	}
	switch {
	case isRelationKind(v.RenameType) && v.Relation != nil:
		if s.alterKindRefused(v.RenameType, v.Relation) {
			s.stop()
			return
		}
		schema, _, resolved := s.lookup(v.Relation)
		// A name its schema holds is refused.
		if resolved && s.nameTaken(schema, v.Newname) {
			s.stop()
			return
		}
		if home, ok := s.madeIn(v.Relation); ok && (s.nameTaken(home, v.Newname) || s.madeHere(&ast.RangeVar{Schemaname: home, Relname: v.Newname})) {
			s.stop()
			return
		}
		s.follow(v.Relation, "", v.Newname)
		s.followParents(v.Relation, tableRef{s.schemaOf(v.Relation), v.Newname})
		s.touchName(v.Relation.Schemaname, v.Relation.Relname)
		s.touchName(v.Relation.Schemaname, v.Newname)
		if resolved {
			s.renamedFrom[v.Newname] = renamedRelation{tableRef{schema, v.Relation.Relname}, schema, s.lastTouch[v.Newname]}
		}
		for _, p := range s.matching(v.Relation) {
			// A rename naming the table's own known schema is that table's;
			// otherwise it may be another table of that name.
			if v.Relation.Schemaname != "" && v.Relation.Schemaname == p.schema {
				delete(p.names, v.Relation.Relname)
			}
			p.names[v.Newname] = true
		}
		if v.RenameType == ast.OBJECT_INDEX {
			s.renamedKeys[v.Relation.Relname] = true
			s.renamedKeys[v.Newname] = true
		}
		// The change's own objects of that name, as written, keep their
		// dependencies under the new name.
		old := tableRef{v.Relation.Schemaname, v.Relation.Relname}
		s.retarget(s.newReads, old, v.Newname)
		s.retarget(s.newReturns, old, v.Newname)
		for i := range s.newReferences {
			r := &s.newReferences[i]
			if r.owner == old && (old.schema != "" || r.ownerPath == s.pathVersion) {
				r.owner.table = v.Newname
			}
		}
	case v.RenameType == ast.OBJECT_TYPE:
		if parts := nameParts(listOf(v.Object)); len(parts) > 0 {
			schema := ""
			if len(parts) >= 2 {
				schema = parts[len(parts)-2]
			}
			s.touchName(schema, parts[len(parts)-1])
			s.touchName(schema, v.Newname)
		}
	case v.RenameType == ast.OBJECT_TABCONSTRAINT && v.Relation != nil:
		if t, ok := s.table(v.Relation); ok && s.renameConstraintRefused(t, v.Subname, v.Newname) {
			s.stop()
			return
		}
		schema := s.schemaOf(v.Relation)
		s.constraints[[3]string{schema, v.Relation.Relname, v.Subname}] = true
		s.constraints[[3]string{schema, v.Relation.Relname, v.Newname}] = true
	case v.RenameType == ast.OBJECT_COLUMN && v.Relation != nil:
		// A known table must have the column, and not the new name.
		t, known := s.table(v.Relation)
		if known && (!s.hasColumn(t, v.Subname) || s.hasColumn(t, v.Newname)) {
			s.stop()
			return
		}
		if known {
			// The column keeps what it was in the synced table.
			if origin, ok := s.originOf(t, v.Subname); ok {
				s.origins[[3]string{t.schema, t.table, v.Newname}] = origin
			} else {
				delete(s.origins, [3]string{t.schema, t.table, v.Newname})
			}
		}
		schema := s.schemaOf(v.Relation)
		s.setColumn(schema, v.Relation.Relname, v.Subname, false, false)
		s.setColumn(schema, v.Relation.Relname, v.Newname, true, true)
	case v.RenameType == ast.OBJECT_FUNCTION || v.RenameType == ast.OBJECT_PROCEDURE || v.RenameType == ast.OBJECT_ROUTINE:
		if fn, ok := routineOf(v.Object); ok && s.routineMissing(fn, v.RenameType) {
			s.stop()
			return
		}
		s.renameRoutine(v)
	case v.RenameType == ast.OBJECT_SCHEMA:
		s.stop()
	}
}

// renameConstraintRefused reports whether RENAME CONSTRAINT of a known
// table certainly fails: the table has a constraint of the new name, or
// the constraint owns an index the new name of which its schema holds.
func (s *scan) renameConstraintRefused(t tableRef, from, to string) bool {
	table := s.index.schemas[t.schema].tables[t.table]
	if _, ok := constraint(table, to); ok && s.constraintKnown(t, to) && !s.dropped[[3]string{t.schema, t.table, to}] {
		return true
	}
	return s.index.constraintIndexes[tableRef{t.schema, from}] && s.constraintKnown(t, from) && s.nameTaken(t.schema, to)
}

// alterKindRefused reports whether ALTER of an object type certainly names
// a relation of another kind, as the snapshot or the change's own creation
// tells: ALTER VIEW, MATERIALIZED VIEW, SEQUENCE, or FOREIGN TABLE takes
// only its own kind, ALTER TABLE any relation but a composite type, and
// ALTER INDEX ... RENAME any relation at all.
//
// pg: src/backend/commands/tablecmds.c — RangeVarCallbackForAlterRelation
func (s *scan) alterKindRefused(t ast.ObjectType, rv *ast.RangeVar) bool {
	kind := relationKind(0)
	if _, k, ok := s.lookup(rv); ok {
		kind = k
	} else if c, made := s.madeAt(rv); made {
		kind = c.kind
	}
	switch {
	case kind == 0 || kind == kindAmbiguous || t == ast.OBJECT_INDEX:
		return false
	case t == ast.OBJECT_TABLE:
		return kind == kindCompositeType || kind == kindType
	}
	return !dropKindMatches(t, kind)
}

// renamedRelation is a relation of the synced schema a rename or a move
// gave a new name: the relation as synced, the schema it is in now, and
// the touch that gave the name.
type renamedRelation struct {
	relation tableRef
	schema   string
	made     int
}

// renamedAway returns the synced relation a name a rename or a move gave
// may mean, while no statement used the name since.
func (s *scan) renamedAway(rv *ast.RangeVar) (tableRef, bool) {
	r, ok := s.renamedFrom[rv.Relname]
	if !ok || s.lastTouch[rv.Relname] != r.made || rv.Schemaname != "" && rv.Schemaname != r.schema {
		return tableRef{}, false
	}
	return r.relation, true
}

// movesWith returns what SET SCHEMA moves with a synced relation: a
// table's own indexes and the sequences it owns, a materialized view's
// indexes. A table's partitions stay.
func (s *scan) movesWith(schema, name string) []string {
	ns := s.index.schemas[schema]
	t := ns.tables[name]
	var out []string
	for _, d := range ns.dependents[name] {
		switch {
		case t == nil:
			out = append(out, d)
		case ns.relations[d] == kindSequence || slices.ContainsFunc(t.GetIndexes(), func(i *metadata.IndexMetadata) bool { return i.GetName() == d }):
			out = append(out, d)
		}
	}
	return out
}

// nameTaken reports whether a synced schema certainly holds a relation or
// an enum type of that name, which a relation cannot take.
func (s *scan) nameTaken(schema, name string) bool {
	ns := s.index.schemas[schema]
	if ns == nil || s.schemas[schema] || s.isTouched(&ast.RangeVar{Schemaname: schema, Relname: name}) {
		return false
	}
	_, exists := ns.relations[name]
	return exists || slices.Contains(ns.types, name)
}

// follow records that what the change's own objects depend on keeps
// depending on a relation renamed or moved to another schema: each
// dependency on a relation the statement's may be is copied to the new
// schema, when one is given, and name.
func (s *scan) follow(rv *ast.RangeVar, newSchema, name string) {
	schema := s.schemaOf(rv)
	may := func(t tableRef) bool {
		return t.table == rv.Relname && (t.schema == "" || schema == "" || t.schema == schema)
	}
	moved := func(t tableRef) tableRef {
		if newSchema != "" {
			t.schema = newSchema
		}
		t.table = name
		return t
	}
	for _, deps := range []*[]newDependency{&s.newReads, &s.newReturns} {
		for i, n := 0, len(*deps); i < n; i++ {
			if d := (*deps)[i]; !d.retired && may(d.relation) {
				d.relation = moved(d.relation)
				*deps = append(*deps, d)
			}
		}
	}
	for i, n := 0, len(s.newReferences); i < n; i++ {
		if r := s.newReferences[i]; !r.retired && may(tableRef{r.schema, r.table}) {
			t := moved(tableRef{r.schema, r.table})
			r.schema, r.table = t.schema, t.table
			s.newReferences = append(s.newReferences, r)
		}
	}
}

// madeRelation is a table or view the change created, and the touch
// that made its name.
type madeRelation struct {
	kind relationKind
	made int
	// def is a table's definition when its CREATE said it all and no
	// ALTER TABLE changed it since.
	def *tableDef
}

// made records a table or view a statement certainly created under a
// name: the schema written, or the one the search path creates it in.
func (s *scan) made(rv *ast.RangeVar, kind relationKind) {
	schema := rv.Schemaname
	if schema == "" {
		var ok bool
		if schema, ok = s.creationSchema(rv); !ok {
			return
		}
	}
	s.created[[2]string{schema, rv.Relname}] = &madeRelation{kind: kind, made: s.lastTouch[rv.Relname]}
}

// dropsOwn reports whether a qualified DROP certainly drops a table or
// view the change created, under the name it made, which nothing the
// change created depends on, or CASCADE takes with it.
func (s *scan) dropsOwn(rv *ast.RangeVar, kind ast.ObjectType, cascade bool, written []tableRef) bool {
	if rv.Schemaname == "" {
		return false
	}
	t := tableRef{rv.Schemaname, rv.Relname}
	c, ok := s.created[[2]string{t.schema, t.table}]
	if !ok || s.lastTouch[t.table] != c.made || !dropKindMatches(kind, c.kind) {
		return false
	}
	return cascade || !s.dependsOnNew(s.newReads, t, written) && !s.dependsOnNew(s.newReturns, t, nil) && !s.newlyReferenced(t, nil, true, written...)
}

// setSchema follows SET SCHEMA. A pending table it may be can then be in
// any schema.
func (s *scan) setSchema(v *ast.AlterObjectSchemaStmt) {
	switch {
	case isRelationKind(v.ObjectType) && v.Relation != nil:
		if s.index != nil && s.missing(v.Relation) {
			if !v.MissingOk {
				s.stop()
			}
			return
		}
		if s.alterKindRefused(v.ObjectType, v.Relation) {
			s.stop()
			return
		}
		schema, _, resolved := s.lookup(v.Relation)
		// A schema that lacks the name must exist to take it.
		if s.index != nil && (s.index.schemas[v.Newschema] == nil && !s.schemas[v.Newschema] || s.schemaGone[v.Newschema]) ||
			resolved && s.nameTaken(v.Newschema, v.Relation.Relname) {
			s.stop()
			return
		}
		s.follow(v.Relation, v.Newschema, v.Relation.Relname)
		s.followParents(v.Relation, tableRef{v.Newschema, v.Relation.Relname})
		if resolved {
			// The relation leaves its schema, with a table's indexes and
			// owned sequences.
			s.free(schema, v.Relation.Relname)
			for _, d := range s.movesWith(schema, v.Relation.Relname) {
				s.free(schema, d)
				s.touchName(v.Newschema, d)
			}
		} else {
			s.touch(v.Relation)
		}
		s.touchName(v.Newschema, v.Relation.Relname)
		if resolved {
			s.renamedFrom[v.Relation.Relname] = renamedRelation{tableRef{schema, v.Relation.Relname}, v.Newschema, s.lastTouch[v.Relation.Relname]}
		}
		for _, p := range s.matching(v.Relation) {
			// The move is known only when the statement names the table's
			// own schema; otherwise it may be another table of that name.
			if v.Relation.Schemaname != "" && v.Relation.Schemaname == p.schema {
				p.schema = v.Newschema
			} else {
				p.schema = ""
			}
		}
	case v.ObjectType == ast.OBJECT_TYPE:
		if parts := nameParts(listOf(v.Object)); len(parts) > 0 {
			s.touchName("", parts[len(parts)-1])
		}
	case v.ObjectType == ast.OBJECT_FUNCTION || v.ObjectType == ast.OBJECT_PROCEDURE || v.ObjectType == ast.OBJECT_ROUTINE:
		if fn, ok := routineOf(v.Object); ok && s.routineMissing(fn, v.ObjectType) {
			if !v.MissingOk {
				s.stop()
			}
			return
		}
		if owa, ok := v.Object.(*ast.ObjectWithArgs); ok {
			if parts := nameParts(owa.Objname); len(parts) > 0 {
				s.madeRoutines[parts[len(parts)-1]] = true
				s.movedRoutines[parts[len(parts)-1]] = true
			}
		}
	case v.ObjectType == ast.OBJECT_EXTENSION:
		s.stop()
	}
}

// drop follows DROP. A cascading drop of a relation also takes the
// relations that go with it (its partitions, indexes, owned sequences,
// and the views that read it), the foreign keys that reference it, and
// any inheritance child, which the synced schema does not record; a
// cascading drop of anything else may take columns of any table, which
// ends the scan. DROP EXTENSION drops the extension's objects.
func (s *scan) drop(v *ast.DropStmt) {
	kind := ast.ObjectType(v.RemoveType)
	cascade := v.Behavior == int(ast.DROP_CASCADE)
	if kind == ast.OBJECT_EXTENSION || cascade && !isRelationKind(kind) {
		s.stop()
		return
	}
	if kind == ast.OBJECT_FUNCTION || kind == ast.OBJECT_PROCEDURE || kind == ast.OBJECT_ROUTINE {
		s.dropFunctions(v)
		return
	}
	if kind == ast.OBJECT_SCHEMA && v.Objects != nil {
		// Without CASCADE the schema must be empty; its name may then be
		// created again. A schema the target lacks and the change did not
		// create is not there to drop: IF EXISTS drops nothing, and
		// without it the statement is refused.
		for _, obj := range v.Objects.Items {
			parts := nameParts(listOf(obj))
			if len(parts) != 1 {
				continue
			}
			if s.index != nil && (s.index.schemas[parts[0]] == nil && !s.schemas[parts[0]] || s.schemaGone[parts[0]]) {
				if !v.Missing_ok {
					s.stop()
					return
				}
				continue
			}
			if s.schemaHolds(parts[0]) {
				s.stop()
				return
			}
			s.schemas[parts[0]] = true
			s.schemaGone[parts[0]] = true
		}
		return
	}
	// A type a column, a function, or another type may use is not dropped
	// without CASCADE; the snapshot does not record every such use.
	if kind == ast.OBJECT_TYPE || kind == ast.OBJECT_DOMAIN {
		s.stop()
		return
	}
	if !isRelationKind(kind) || v.Objects == nil {
		return
	}
	// A name in another database is refused.
	for _, obj := range v.Objects.Items {
		if parts := nameParts(listOf(obj)); len(parts) == 3 && (s.index == nil || parts[0] != s.index.database) {
			if s.index != nil {
				s.stop()
			}
			return
		}
	}
	if !cascade && s.dropRefused(v) || !v.Missing_ok && s.dropsMissing(v) || s.dropsWrongKind(v) {
		s.stop()
		return
	}
	// The names the statement drops, as written: a view among them does
	// not hold the drop of another.
	var written []tableRef
	for _, obj := range v.Objects.Items {
		if parts := nameParts(listOf(obj)); len(parts) > 0 {
			t := tableRef{table: parts[len(parts)-1]}
			if len(parts) >= 2 {
				t.schema = parts[len(parts)-2]
			}
			written = append(written, t)
		}
	}
	var droppedRefs []tableRef
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) >= 2 {
			rv.Schemaname = parts[len(parts)-2]
		}
		schema, relKind, resolved := s.lookup(rv)
		ref := tableRef{rv.Schemaname, rv.Relname}
		if resolved {
			ref.schema = schema
		}
		switch {
		case resolved:
			s.free(schema, rv.Relname)
		case v.Missing_ok && s.index != nil && s.missing(rv):
			// IF EXISTS of a name the target lacks drops nothing.
			continue
		case s.dropsOwn(rv, kind, cascade, written):
			// The change's own table or view, under the name it made.
			s.free(rv.Schemaname, rv.Relname)
			schema = ""
		default:
			s.touch(rv)
			schema = ""
		}
		droppedRefs = append(droppedRefs, ref)
		if kind == ast.OBJECT_INDEX {
			s.renamedKeys[rv.Relname] = true
		}
		for _, p := range s.matching(rv) {
			p.settle()
		}
		switch kind {
		case ast.OBJECT_TABLE:
			s.retireReferences(rv, "")
		case ast.OBJECT_VIEW, ast.OBJECT_MATVIEW:
			s.retire(s.newReads, tableRef{rv.Schemaname, rv.Relname})
		}
		s.dropOwned(rv, schema, resolved && relKind == kindTable)
		if cascade {
			s.dropDependents(rv, schema)
			if resolved {
				s.dropNewReaders(tableRef{schema, rv.Relname})
			}
		}
	}
	if cascade && len(droppedRefs) > 0 && (kind == ast.OBJECT_TABLE || kind == ast.OBJECT_FOREIGN_TABLE) {
		// An inheritance child of a dropped table goes with it. The synced
		// schema records no inheritance, so a table of it may be a child.
		for _, p := range s.pending {
			if p.existed || slices.ContainsFunc(p.parents, func(parent tableRef) bool {
				return slices.ContainsFunc(droppedRefs, func(d tableRef) bool { return mayBe(parent, d) })
			}) {
				p.settle()
			}
		}
	}
}

// dropFunctions follows DROP FUNCTION: a synced function whose name the
// scan resolves to one function, and a function the change created, no
// longer hold the relation whose row type they return.
func (s *scan) dropFunctions(v *ast.DropStmt) {
	if v.Objects == nil {
		return
	}
	for _, obj := range v.Objects.Items {
		owa, ok := obj.(*ast.ObjectWithArgs)
		if !ok {
			continue
		}
		parts := nameParts(owa.Objname)
		if len(parts) == 0 || len(parts) > 2 {
			continue
		}
		fn := tableRef{table: parts[len(parts)-1]}
		if len(parts) == 2 {
			fn.schema = parts[0]
		}
		// A name alone that several routines have is refused, IF EXISTS
		// or not.
		if owa.ArgsUnspecified && s.routinesOf(fn) > 1 {
			s.stop()
			return
		}
		// A routine the target certainly lacks is not there to drop; IF
		// EXISTS then drops nothing.
		if s.routineMissing(fn, ast.ObjectType(v.RemoveType)) {
			if !v.Missing_ok {
				s.stop()
				return
			}
			continue
		}
		// A function is retired only when the statement leaves no choice
		// of which one it drops: an argument list matching the one the
		// change wrote, or no argument list on a name only one function
		// has. A synced function's argument types are not compared.
		if sigs := s.newSignatures[fn]; len(sigs) > 0 && s.syncedFunctions(fn) == 0 {
			sig, i := argSignature(owa.Objargs), -1
			if owa.ArgsUnspecified {
				if len(sigs) == 1 {
					sig, i = sigs[0], 0
				}
			} else {
				i = slices.Index(sigs, sig)
			}
			if i >= 0 {
				s.retireSignature(fn, sig)
				s.newSignatures[fn] = slices.Delete(slices.Clone(sigs), i, i+1)
				s.newFunctions[fn]--
			}
		}
		// A synced function is certainly dropped only by its name alone,
		// when no routine of the change has the name and none moved away
		// from it.
		if s.index == nil || s.madeRoutines[fn.table] || s.movedRoutines[fn.table] || !owa.ArgsUnspecified {
			continue
		}
		if fn.schema != "" {
			if s.index.functions[fn] == 1 {
				s.droppedFunctions[fn] = true
			}
			continue
		}
		path, ok := s.searchPath()
		if !ok {
			continue
		}
		// The server finds the first schema on the path with a function of
		// that name; one function there is the one dropped.
		for _, name := range path {
			if n := s.index.functions[tableRef{name, fn.table}]; n > 0 {
				if n == 1 {
					s.droppedFunctions[tableRef{name, fn.table}] = true
				}
				break
			}
		}
	}
}

// syncedFunctions counts the synced functions a routine name may mean: of
// that name in the schema written, or in any schema.
func (s *scan) syncedFunctions(fn tableRef) int {
	if s.index == nil {
		return 0
	}
	n := 0
	for name, count := range s.index.functions {
		if name.table == fn.table && (fn.schema == "" || name.schema == fn.schema) {
			n += count
		}
	}
	return n
}

// routinesOf counts the routines a qualified name certainly has: the
// synced ones of the name, none dropped or renamed by the change, and the
// signatures the change created under the name as written. It is 0 for a
// name the scan cannot count.
func (s *scan) routinesOf(fn tableRef) int {
	if s.index == nil || fn.schema == "" || s.movedRoutines[fn.table] || s.droppedFunctions[fn] || s.schemas[fn.schema] {
		return 0
	}
	return s.index.functions[fn] + len(s.newSignatures[fn])
}

// routineMissing reports whether no routine of that name certainly exists
// where the name resolves: the change created, renamed, or moved none of
// that name, and the synced schema written, or every schema of a known
// search path, has none. A name only the system catalogs have cannot be
// dropped either.
func (s *scan) routineMissing(fn tableRef, kind ast.ObjectType) bool {
	if s.index == nil || s.madeRoutines[fn.table] {
		return false
	}
	// DROP FUNCTION refuses a procedure, and DROP PROCEDURE a function.
	want := kindFunction | kindProcedure
	switch kind {
	case ast.OBJECT_FUNCTION:
		want = kindFunction
	case ast.OBJECT_PROCEDURE:
		want = kindProcedure
	}
	if fn.schema != "" {
		return s.index.schemas[fn.schema] != nil && !s.schemas[fn.schema] && s.index.routineKinds[fn]&want == 0
	}
	path, ok := s.searchPath()
	if !ok {
		return false
	}
	for _, name := range path {
		if s.schemas[name] || name == "information_schema" || s.index.routineKinds[tableRef{name, fn.table}]&want != 0 {
			return false
		}
	}
	return true
}

// routineOf returns the routine an ALTER names, as written.
func routineOf(n ast.Node) (tableRef, bool) {
	owa, ok := n.(*ast.ObjectWithArgs)
	if !ok {
		return tableRef{}, false
	}
	parts := nameParts(owa.Objname)
	if len(parts) == 0 || len(parts) > 2 {
		return tableRef{}, false
	}
	fn := tableRef{table: parts[len(parts)-1]}
	if len(parts) == 2 {
		fn.schema = parts[0]
	}
	return fn, true
}

// renameRoutine follows ALTER FUNCTION, PROCEDURE, or ROUTINE ... RENAME:
// the change's own function the statement certainly names keeps its
// dependencies under the new name.
func (s *scan) renameRoutine(v *ast.RenameStmt) {
	s.madeRoutines[v.Newname] = true
	owa, ok := v.Object.(*ast.ObjectWithArgs)
	if !ok {
		return
	}
	parts := nameParts(owa.Objname)
	if len(parts) == 0 || len(parts) > 2 {
		return
	}
	s.movedRoutines[parts[len(parts)-1]] = true
	fn := tableRef{table: parts[len(parts)-1]}
	if len(parts) == 2 {
		fn.schema = parts[0]
	}
	sigs := s.newSignatures[fn]
	if len(sigs) == 0 || s.syncedFunctions(fn) > 0 || owa.ArgsUnspecified && len(sigs) > 1 {
		return
	}
	sig := argSignature(owa.Objargs)
	if owa.ArgsUnspecified {
		sig = sigs[0]
	}
	i := slices.Index(sigs, sig)
	if i < 0 {
		return
	}
	for j := range s.newReturns {
		d := &s.newReturns[j]
		if !d.retired && s.sameObject(*d, fn) && d.signature == sig {
			d.object.table = v.Newname
		}
	}
	to := tableRef{fn.schema, v.Newname}
	s.newSignatures[fn] = slices.Delete(slices.Clone(sigs), i, i+1)
	s.newFunctions[fn]--
	s.newSignatures[to] = append(s.newSignatures[to], sig)
	s.newFunctions[to]++
}

// dropReturners records the functions a cascading drop of a relation
// takes with its row type: a synced function that depends on it, when it
// is the only function of its name, and a function the change created
// whose result, arguments, or SQL-standard body named the relation,
// resolved to it.
func (s *scan) dropReturners(t tableRef) {
	if s.index != nil {
		for _, fn := range s.index.returnedBy[t] {
			if s.index.functions[fn] == 1 {
				s.droppedFunctions[fn] = true
			}
		}
	}
	for i := range s.newReturns {
		d := s.newReturns[i]
		if d.retired || d.relation != t {
			continue
		}
		for j := range s.newReturns {
			if o := &s.newReturns[j]; o.object == d.object && o.path == d.path && o.signature == d.signature {
				o.retired = true
			}
		}
	}
}

// dropsWrongKind reports whether a DROP names a relation of another kind,
// which the server refuses: DROP TABLE of a view, DROP VIEW of a table, or
// DROP INDEX of the index a constraint owns.
func (s *scan) dropsWrongKind(v *ast.DropStmt) bool {
	kind := ast.ObjectType(v.RemoveType)
	if s.index == nil || v.Objects == nil || !isRelationKind(kind) {
		return false
	}
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 || len(parts) > 2 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) == 2 {
			rv.Schemaname = parts[0]
		}
		schema, relKind, ok := s.lookup(rv)
		if !ok {
			// A relation the change created has the kind it was made as.
			if c, made := s.madeAt(rv); made && !dropKindMatches(kind, c.kind) {
				return true
			}
			continue
		}
		if relKind == kindAmbiguous {
			continue
		}
		if !dropKindMatches(kind, relKind) {
			return true
		}
		if kind == ast.OBJECT_INDEX && s.index.constraintIndexes[tableRef{schema, rv.Relname}] {
			return true
		}
	}
	return false
}

// freeView records that a cascade dropped a synced view. Its name is freed
// only when the change has not touched it: a view it replaced may no
// longer depend on what was dropped, and one it renamed is not under that
// name.
func (s *scan) freeView(view tableRef) bool {
	if s.viewHolds(view) && !s.isTouched(&ast.RangeVar{Schemaname: view.schema, Relname: view.table}) {
		s.free(view.schema, view.table)
		return true
	}
	s.touchName(view.schema, view.table)
	return false
}

// cascadeViews records the views a cascade takes: the certain ones are
// freed where the change left them as synced and touched otherwise, the
// others touched. The views reading a freed view go for certain; those
// reading a view the scan could not free go only maybe.
func (s *scan) cascadeViews(certain, maybe []tableRef) {
	seen := make(map[tableRef]bool)
	for len(certain) > 0 || len(maybe) > 0 {
		if len(certain) > 0 {
			view := certain[0]
			certain = certain[1:]
			if seen[view] {
				continue
			}
			seen[view] = true
			if s.freeView(view) {
				certain = append(certain, s.index.readers[view]...)
				s.dropNewReaders(view)
			} else {
				maybe = append(maybe, s.index.readers[view]...)
			}
			maybe = append(maybe, s.index.mayRead[view]...)
			continue
		}
		view := maybe[0]
		maybe = maybe[1:]
		if seen[view] {
			continue
		}
		seen[view] = true
		s.touchName(view.schema, view.table)
		maybe = append(append(maybe, s.index.readers[view]...), s.index.mayRead[view]...)
	}
}

// viewHolds reports whether a synced view may still read what the snapshot
// says it reads: the change neither dropped it nor redefined it.
func (s *scan) viewHolds(view tableRef) bool {
	return !s.freed[[2]string{view.schema, view.table}] && !s.replacedViews[view]
}

// dropKindMatches reports whether a DROP of an object type can drop a
// relation of that kind.
func dropKindMatches(t ast.ObjectType, kind relationKind) bool {
	switch t {
	case ast.OBJECT_TABLE:
		return kind == kindTable || kind == kindPartition
	case ast.OBJECT_VIEW:
		return kind == kindView
	case ast.OBJECT_MATVIEW:
		return kind == kindMatView
	case ast.OBJECT_FOREIGN_TABLE:
		return kind == kindForeignTable
	case ast.OBJECT_SEQUENCE:
		return kind == kindSequence
	case ast.OBJECT_INDEX:
		return kind == kindIndex
	}
	return true
}

// schemaHolds reports whether a schema may still hold an object, so that
// DROP SCHEMA without CASCADE may be refused and the scan ends: a name the
// change created, renamed, or moved that may be in it and no resolved drop
// removed, a function the change created that may be in it, or, in a
// synced schema, a relation of it no statement touched, an enum type no
// statement named, or a routine no DROP FUNCTION certainly removed.
func (s *scan) schemaHolds(name string) bool {
	if s.index == nil {
		return false
	}
	// What the change created, renamed, or moved may be there: under that
	// schema, or under none when the search path may put it there, unless
	// a drop the scan resolved removed it.
	path, known := s.searchPath()
	onPath := !known || slices.Contains(path, name)
	for key := range s.touched {
		if (key[0] == name || key[0] == "" && onPath) && !s.freed[[2]string{name, key[1]}] {
			return true
		}
	}
	for fn, n := range s.newFunctions {
		if n > 0 && (fn.schema == name || fn.schema == "" && onPath) {
			return true
		}
	}
	ns := s.index.schemas[name]
	if ns == nil || s.schemas[name] {
		return false
	}
	for rel := range ns.relations {
		if !s.isTouched(&ast.RangeVar{Schemaname: name, Relname: rel}) {
			return true
		}
	}
	for _, t := range ns.types {
		if !s.isTouched(&ast.RangeVar{Schemaname: name, Relname: t}) {
			return true
		}
	}
	dropped := 0
	for fn := range s.droppedFunctions {
		if fn.schema == name {
			dropped += s.index.functions[fn]
		}
	}
	return ns.routines > dropped
}

// schemaMissing reports whether a qualified name's schema certainly does
// not exist: the target lacks it and the change did not create it.
func (s *scan) schemaMissing(rv *ast.RangeVar) bool {
	return s.index != nil && rv != nil && rv.Schemaname != "" && !isTemp(rv) &&
		(s.index.schemas[rv.Schemaname] == nil && !s.schemas[rv.Schemaname] || s.schemaGone[rv.Schemaname])
}

// dropsMissing reports whether a DROP without IF EXISTS names a relation
// the target certainly lacks: no statement of the change made the name,
// and the schema the SQL wrote, or every schema of a known search path,
// holds no relation of that name.
func (s *scan) dropsMissing(v *ast.DropStmt) bool {
	if s.index == nil || v.Objects == nil || !isRelationKind(ast.ObjectType(v.RemoveType)) {
		return false
	}
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 || len(parts) > 2 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) == 2 {
			rv.Schemaname = parts[0]
		}
		if s.missing(rv) {
			return true
		}
	}
	return false
}

// missing reports whether the target certainly has no relation of that
// name where the name resolves.
func (s *scan) missing(rv *ast.RangeVar) bool {
	// A resolved drop removed the name, and nothing used it since.
	if rv.Schemaname != "" && s.freed[[2]string{rv.Schemaname, rv.Relname}] {
		return true
	}
	if s.isTouched(rv) {
		return false
	}
	if rv.Schemaname != "" {
		// A schema the change created started empty: what no statement
		// named in it, and no generated name may be, is not there.
		if s.schemas[rv.Schemaname] && !s.schemaGone[rv.Schemaname] {
			return !s.generated[rv.Schemaname] && !s.generated[""]
		}
		ns := s.index.schemas[rv.Schemaname]
		if ns == nil || s.schemas[rv.Schemaname] {
			return false
		}
		_, ok := ns.relations[rv.Relname]
		return !ok
	}
	if strings.HasPrefix(rv.Relname, "pg_") {
		return false
	}
	path, ok := s.searchPath()
	if !ok {
		return false
	}
	for _, name := range path {
		if s.schemas[name] || name == "information_schema" {
			return false
		}
		if ns := s.index.schemas[name]; ns != nil {
			if _, ok := ns.relations[rv.Relname]; ok {
				return false
			}
		}
	}
	return true
}

// dropRefused reports whether the server refuses a DROP without CASCADE of
// relations the scan resolves: a view the statement does not drop reads
// one, a function returns its row type, or, for a table, a foreign key of
// a table the statement does not drop references it, as the synced schema
// or the change itself records.
func (s *scan) dropRefused(v *ast.DropStmt) bool {
	kind := ast.ObjectType(v.RemoveType)
	if s.index == nil || kind != ast.OBJECT_TABLE && kind != ast.OBJECT_VIEW && kind != ast.OBJECT_MATVIEW && kind != ast.OBJECT_FOREIGN_TABLE {
		return false
	}
	var dropped, written []tableRef
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) >= 2 {
			rv.Schemaname = parts[len(parts)-2]
		}
		written = append(written, tableRef{rv.Schemaname, rv.Relname})
		if schema, _, ok := s.lookup(rv); ok {
			dropped = append(dropped, tableRef{schema, rv.Relname})
		} else if t, ok := s.renamedAway(rv); ok {
			// What depends on a renamed relation still does.
			dropped = append(dropped, t)
		} else if home, ok := s.madeIn(rv); ok {
			// So does what depends on a relation the change created.
			dropped = append(dropped, tableRef{home, rv.Relname})
		}
	}
	for _, t := range dropped {
		// A new view the statement drops too does not hold the drop.
		if s.dependsOnNew(s.newReads, t, written) || s.dependsOnNew(s.newReturns, t, nil) {
			return true
		}
		for _, fn := range s.index.returnedBy[t] {
			if !s.droppedFunctions[fn] && !s.replacedFunctions[fn] {
				return true
			}
		}
		// A view whose definition the scan cannot read may read anything.
		if s.index.opaqueViews {
			return true
		}
		for _, view := range append(slices.Clone(s.index.readers[t]), s.index.mayRead[t]...) {
			if s.viewHolds(view) && !slices.Contains(dropped, view) {
				return true
			}
		}
		if kind != ast.OBJECT_TABLE {
			continue
		}
		if s.newlyReferenced(t, nil, true, written...) {
			return true
		}
		for _, fk := range s.index.foreignKeys {
			if fk.referenced == t && !fk.in(s.dropped) && !slices.Contains(dropped, fk.owner) {
				return true
			}
		}
	}
	return false
}

// dropOwned records what any drop of a relation takes with it, CASCADE or
// not: its partitions, indexes, and owned sequences, freed in its schema
// when the scan resolved it and touched in every schema the name may be in
// otherwise, and, for a table the scan resolved, its own foreign keys,
// which then no longer hold a key of another table.
func (s *scan) dropOwned(rv *ast.RangeVar, schema string, table bool) {
	if s.index == nil {
		return
	}
	for name, ns := range s.index.schemas {
		switch {
		case schema != "" && name == schema:
			for _, d := range ns.dependents[rv.Relname] {
				if s.reowned[[2]string{name, d}] || s.reowned[[2]string{"", d}] {
					s.touchName(name, d)
				} else {
					s.free(name, d)
				}
			}
		case schema == "" && (rv.Schemaname == "" || rv.Schemaname == name):
			for _, d := range ns.dependents[rv.Relname] {
				s.touchName(name, d)
			}
		}
	}
	if !table {
		return
	}
	for _, fk := range s.index.foreignKeys {
		if fk.owner == (tableRef{schema, rv.Relname}) {
			s.dropped[[3]string{schema, rv.Relname, fk.name}] = true
		}
	}
}

// dropDependents records what a cascading drop of a relation takes with
// it besides what it owns, as far as the synced schema records it: the
// views that read it, freed when the scan resolved the relation to schema
// and touched in every schema the name may be in otherwise, and the
// foreign keys that reference it.
func (s *scan) dropDependents(rv *ast.RangeVar, schema string) {
	s.dropReferences(rv, tableRef{schema, rv.Relname}, schema != "", nil)
	if s.index == nil {
		return
	}
	if schema != "" {
		t := tableRef{schema, rv.Relname}
		s.cascadeViews(s.index.readers[t], s.index.mayRead[t])
		return
	}
	var maybe []tableRef
	for name := range s.index.schemas {
		if rv.Schemaname == "" || rv.Schemaname == name {
			t := tableRef{name, rv.Relname}
			maybe = append(append(maybe, s.index.readers[t]...), s.index.mayRead[t]...)
		}
	}
	s.cascadeViews(nil, maybe)
}

// dropReferences records that foreign keys a cascading drop takes may be
// gone: those referencing t when the scan resolved it (known), and those
// referencing a table of that name, in the schema the SQL wrote if any,
// otherwise. affected, when set, narrows them to the ones that may depend
// on what was dropped, and says which certainly do; those are recorded as
// dropped.
func (s *scan) dropReferences(rv *ast.RangeVar, t tableRef, known bool, affected func(foreignKeyRef) (maybe, certainly bool)) {
	if s.index == nil {
		return
	}
	for _, fk := range s.index.foreignKeys {
		switch {
		case known && fk.referenced != t:
			continue
		case !known && (fk.referenced.table != rv.Relname || rv.Schemaname != "" && fk.referenced.schema != rv.Schemaname):
			continue
		}
		maybe, certainly := true, false
		if affected != nil {
			maybe, certainly = affected(fk)
		}
		if !maybe {
			continue
		}
		s.constraints[[3]string{fk.owner.schema, fk.owner.table, fk.name}] = true
		if known && certainly {
			s.dropped[[3]string{fk.owner.schema, fk.owner.table, fk.name}] = true
		}
	}
}

// constraintKnown reports whether a table's constraint of that name is
// still the one the synced schema lists.
func (s *scan) constraintKnown(t tableRef, name string) bool {
	return !s.constraints[[3]string{t.schema, t.table, name}] && !s.constraints[[3]string{"", t.table, name}] && !s.renamedKeys[name]
}
