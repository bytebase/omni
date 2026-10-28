package review

import (
	"context"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// checkWalkThrough reports the first statement of the change that
// certainly fails on the target, given that the statements before it
// succeeded. The walk follows the change on review/walk's ledger, seeded
// with the synced schema as oracleSnapshot trusts it; the scope is
// docs/specs/2026-10-26-walk-through-oracle.md.
func checkWalkThrough(ctx context.Context, sql string, stmts []statement, t review.Target, r *reporter) error {
	v, err := walkThrough(ctx, sql, stmts, t)
	if err != nil || v == nil {
		return err
	}
	note, ok := v.Check.Note.(*walkNote)
	if !ok {
		return nil
	}
	r.report(review.WalkThrough, stmts[v.Statement].index, note.rng, escapeLines(note.message(v)))
	return nil
}

// walkThrough walks the change on one target and returns the first
// violation, or nil. Without the database's name the walk cannot place
// an unqualified name, and does not walk.
func walkThrough(ctx context.Context, sql string, stmts []statement, t review.Target) (*walk.Violation, error) {
	if t.Schema == nil || t.Schema.GetName() == "" {
		return nil, nil
	}
	w := newWalker(sql, stmts, t)
	var snap walk.Snapshot = w.snap
	if t.OtherSpecs {
		// Another spec of the plan may have changed anything first, and
		// added triggers the sync does not show.
		snap = nil
		w.triggers = true
	}
	return walk.Run(ctx, len(stmts), snap, w.lower)
}

// walkNote is a check's finding: where it anchors and what it says when
// the object is certainly absent or certainly present (or of another
// kind).
type walkNote struct {
	rng     review.Range
	absent  string
	present string
}

func (n *walkNote) message(v *walk.Violation) string {
	if v.Got == walk.Absent {
		return n.absent
	}
	return n.present
}

// walker is the session the change runs in, as far as the walk follows
// it.
type walker struct {
	sql   string
	stmts []statement
	snap  *oracleSnapshot
	// current is the schema unqualified names resolve in: the database,
	// which Bytebase makes the session's CURRENT_SCHEMA, until an ALTER
	// SESSION SET CURRENT_SCHEMA moves it; "" after one the walk cannot
	// read.
	//
	// bytebase: backend/plugin/db/oracle/oracle.go — connector.Connect
	current string
	// triggers is set when a DML statement may fire a trigger: the sync
	// lists one, the change created one, or another spec of the plan may
	// have. A trigger can run DDL through an autonomous transaction, and a
	// cascaded delete fires the triggers of other tables (verified).
	triggers bool
}

func newWalker(sql string, stmts []statement, t review.Target) *walker {
	snap := newOracleSnapshot(t.Schema)
	return &walker{
		sql:      sql,
		stmts:    stmts,
		snap:     snap,
		current:  snap.schema,
		triggers: len(snap.triggers) > 0,
	}
}

// lower implements walk.Lower.
func (w *walker) lower(i int, l *walk.Ledger) *walk.Step {
	if namesOutsideASCII(w.sql, w.stmts[i].loc) {
		return nil
	}
	b := &builder{w: w, l: l, s: &w.stmts[i]}
	if !b.statement(w.stmts[i].node) {
		return nil
	}
	if len(b.settled) > 0 {
		b.Next()
		for _, s := range b.settled {
			b.Set(s.key, s.fact, "")
		}
	}
	return b.Step()
}

// Token types the walk reads names and literals by.
var (
	tokenIdent      = parser.Tokenize("x")[0].Type
	tokenQuoted     = parser.Tokenize(`"x"`)[0].Type
	tokenString     = parser.Tokenize("'x'")[0].Type
	tokenNString    = parser.Tokenize("N'x'")[0].Type
	tokenConstraint = parser.Tokenize("CONSTRAINT")[0].Type
	tokenPublic     = parser.Tokenize("PUBLIC")[0].Type
	tokenCreate     = parser.Tokenize("CREATE")[0].Type
	tokenOpenParen  = parser.Tokenize("(")[0].Type
	tokenCloseParen = parser.Tokenize(")")[0].Type
	tokenComma      = parser.Tokenize(",")[0].Type
)

// namesOutsideASCII reports whether a statement's text holds a name with
// a byte outside ASCII, which the walk does not follow: the parser folds
// an unquoted one with Go's rules, which need not be the database's. A
// token that is no string literal and holds such a byte counts.
func namesOutsideASCII(sql string, loc ast.Loc) bool {
	if loc.Start < 0 || loc.End > len(sql) || loc.Start >= loc.End || !outsideASCII(sql[loc.Start:loc.End]) {
		return false
	}
	for _, tok := range tokens(sql, loc) {
		if tok.Type != tokenString && tok.Type != tokenNString && outsideASCII(tok.Str) {
			return true
		}
	}
	return false
}

func outsideASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return true
		}
	}
	return false
}

// tokenName returns the name a token spells, as Oracle folds it, and
// whether it is a name: an identifier, quoted or not, or a keyword Oracle
// takes for one.
func tokenName(tok parser.Token) (string, bool) {
	switch {
	case tok.Type == tokenIdent || tok.Type == tokenQuoted:
		return tok.Str, tok.Str != ""
	case tok.Type == tokenString || tok.Type == tokenNString:
		return "", false
	}
	return tok.Str, plainIdent(tok.Str)
}

// builder lowers one statement into a walk.Step.
type builder struct {
	walk.Builder
	w *walker
	l *walk.Ledger
	s *statement
	// settled holds what holds after the statement whether or not it
	// skips an optional part: a part of its own, after the others.
	settled []settledKey
}

type settledKey struct {
	key  walk.Key
	fact walk.Fact
}

// settle records a fact that holds after the statement whatever its
// optional parts do, of an unknown kind when Present.
func (b *builder) settle(k walk.Key, f walk.Fact) {
	b.settled = append(b.settled, settledKey{k, f})
}

func (b *builder) anchor() review.Range { return rangeOf(b.s.loc) }

// need adds a check, with the messages for a key certainly absent and
// certainly present (or of another kind). A name Oracle generates is never
// checked.
func (b *builder) need(k walk.Key, want walk.Fact, kinds []string, rng review.Range, absent, present string) {
	if generated(k.Space, k.Name) {
		return
	}
	b.Need(k, want, kinds, &walkNote{rng: rng, absent: absent, present: present})
}

// fail records that the statement certainly fails, from its text.
func (b *builder) fail(rng review.Range, message string) {
	b.Fail(nil, walk.Absent, &walkNote{rng: rng, absent: message, present: message})
}

// objRef is a reference to a schema object as the walk resolves it: the
// schema and name as Oracle stores them, and how the SQL wrote it.
type objRef struct {
	schema, name string
	written      string
	rng          review.Range
}

// ref resolves a reference to an object of a schema. It returns false for
// one the walk does not follow, which ends the walk: a database link, or
// an unqualified name after an ALTER SESSION the walk cannot read.
func (b *builder) ref(o *ast.ObjectName) (objRef, bool) {
	if o == nil || o.Name == "" || o.DBLink != "" {
		return objRef{}, false
	}
	r := objRef{schema: o.Schema, name: o.Name, written: objectName(o), rng: rangeOf(o.Loc)}
	if r.schema == "" {
		r.schema = b.w.current
	}
	if r.rng == (review.Range{}) {
		r.rng = b.anchor()
	}
	return r, r.schema != ""
}

// refIn resolves a name a statement gives with no schema, such as the new
// name of ALTER INDEX ... RENAME TO, in the schema of another reference.
func refIn(r objRef, name string) objRef {
	return objRef{schema: r.schema, name: name, written: ident(name), rng: r.rng}
}

func (r objRef) key(space string) walk.Key {
	return walk.Key{Space: space, Parent: walk.Path(r.schema), Name: r.name}
}

func (r objRef) path() string { return tablePath(r.schema, r.name) }

// childKey keys a column, the primary key, or a constraint of the table a
// reference names.
func (r objRef) childKey(space, name string) walk.Key {
	return walk.Key{Space: space, Parent: r.path(), Name: name}
}

// schemaKey keys a constraint name, an index, or a trigger of the schema a
// reference is in.
func (r objRef) schemaKey(space, name string) walk.Key {
	return walk.Key{Space: space, Parent: walk.Path(r.schema), Name: name}
}

// identifierLimit is the longest identifier Oracle takes, in bytes: 128
// from 12.2, and 30 before it.
const identifierLimit = 128

// tooLong fails a statement that writes an identifier longer than Oracle
// allows, which fails with ORA-00972 before anything is looked up
// (verified). Only an ASCII name is judged.
func (b *builder) tooLong(name string, rng review.Range) bool {
	if len(name) <= identifierLimit || outsideASCII(name) {
		return false
	}
	b.fail(rng, "The identifier "+ident(name[:32])+"... is longer than 128 bytes")
	return true
}

// refTooLong checks the parts of a reference.
func (b *builder) refTooLong(o *ast.ObjectName, r objRef) bool {
	return b.tooLong(o.Name, r.rng) || b.tooLong(o.Schema, r.rng)
}

// look reads a key as the statements before this one left it.
func (b *builder) look(k walk.Key) (walk.Fact, string) { return b.l.Lookup(k) }

// arriveObject records that an object of a kind holds a reference's name:
// it is Present in the object space and, for a table, a materialized
// view, a view, or a sequence, in that space, and Absent in the others,
// since the namespace is shared.
func (b *builder) arriveObject(r objRef, kind string) {
	b.Set(r.key(spaceObject), walk.Present, kind)
	for _, s := range []struct {
		space string
		kinds []string
	}{
		{spaceTable, []string{kindTable, kindMView}},
		{spaceView, []string{kindView}},
		{spaceSequence, []string{kindSequence}},
	} {
		if contains(s.kinds, kind) {
			b.Set(r.key(s.space), walk.Present, kind)
		} else {
			b.Set(r.key(s.space), walk.Absent, "")
		}
	}
}

// leaveObject records that no object holds a reference's name.
func (b *builder) leaveObject(r objRef) {
	for _, space := range []string{spaceObject, spaceTable, spaceView, spaceSequence} {
		b.Set(r.key(space), walk.Absent, "")
	}
}

// forgetObject records that the statement may have done anything to a
// reference's name.
func (b *builder) forgetObject(r objRef) {
	for _, space := range []string{spaceObject, spaceTable, spaceView, spaceSequence} {
		b.Set(r.key(space), walk.Unknown, "")
	}
}

// dropOwned records that what was on a table or view may be gone: its
// indexes, constraint names, and triggers, of its path's kind.
func (b *builder) dropOwned(r objRef) {
	b.dropIndexes(r)
	for _, space := range []string{spaceConstraint, spaceTrigger} {
		b.MayRemove(walk.Scope{Space: space, Deep: true, Kinds: []string{r.path()}})
	}
}

// dropIndexes records that a table's indexes may be gone, in any schema:
// an index, and a trigger, may be in another schema than its table.
func (b *builder) dropIndexes(r objRef) {
	b.MayRemove(walk.Scope{Space: spaceIndex, Deep: true, Kinds: []string{r.path()}})
}

// domain reports whether the synced table a reference names has a domain
// index, whose indextype's code runs on the table's DDL and DML.
func (b *builder) domain(r objRef) bool {
	return r.schema == b.w.snap.schema && b.w.snap.domainTables[r.name]
}

// domainIndex reports whether a synced index is a domain index.
func (b *builder) domainIndex(r objRef) bool {
	return r.schema == b.w.snap.schema && b.w.snap.domainIndexes[r.name]
}

// forgetStorage records that a statement on a table of the target's
// schema may have dropped a nested table's storage table, of a name the
// sync lists but does not tie to its parent.
func (b *builder) forgetStorage(r objRef) {
	b.forgetTables(r, b.w.snap.storage)
}

// forgetLogs records that DROP TABLE of a table of the target's schema
// may have dropped its materialized view log's table, of a name the sync
// lists but does not tie to its master.
func (b *builder) forgetLogs(r objRef) {
	b.forgetTables(r, b.w.snap.logs)
}

// forgetTables records that a statement on a table of the target's schema
// may have dropped synced tables of that schema, with what was on them.
func (b *builder) forgetTables(r objRef, names []string) {
	if r.schema != b.w.snap.schema {
		return
	}
	for _, name := range names {
		s := objRef{schema: r.schema, name: name}
		b.forgetObject(s)
		b.Forget(walk.Scope{Parent: s.path(), Deep: true})
		b.dropOwned(s)
	}
}

// forgetTypeDependents records that DROP TYPE ... FORCE may have taken the
// columns of the type, or of a type with an attribute of it, in a table of
// any schema, which become unused, and their constraints and indexes
// (verified).
func (b *builder) forgetTypeDependents() {
	for _, space := range []string{spaceColumn, spacePrimary, spaceTableConstraint, spaceConstraint, spaceIndex} {
		b.MayRemove(walk.Scope{Space: space, Deep: true})
	}
}

// cascadeConstraints records that CASCADE CONSTRAINTS may have dropped the
// foreign keys of any table.
func (b *builder) cascadeConstraints() {
	b.MayRemove(walk.Scope{Space: spaceConstraint, Deep: true})
	b.MayRemove(walk.Scope{Space: spaceTableConstraint, Deep: true})
}

func contains(xs []string, x string) bool {
	for _, y := range xs {
		if y == x {
			return true
		}
	}
	return false
}

// article writes an object kind with its article, for a message.
func article(kind string) string {
	return "a " + strings.ToLower(kind)
}

// Messages, worded as Oracle's errors are.
func missingTable(w string) string { return "Table or view " + w + " does not exist" }
func takenName(w string) string    { return "Name " + w + " is already used by an existing object" }
func takenConstraint(n string) string {
	return "Name " + ident(n) + " is already used by an existing constraint"
}
func missingObject(w string) string   { return "Object " + w + " does not exist" }
func missingIndex(w string) string    { return "Index " + w + " does not exist" }
func missingTrigger(w string) string  { return "Trigger " + w + " does not exist" }
func missingSequence(w string) string { return "Sequence " + w + " does not exist" }

// callsUserCode reports whether a statement may run code of a schema: a
// function call that is qualified, quoted, or whose name is not an Oracle
// SQL builtin, or a name that may be a call without parentheses. In SQL a
// builtin takes precedence over a schema function of the same name, while
// a qualified call runs the schema's (verified with a schema UPPER and
// NVL); a quoted name of a function of a syntax of its own is a plain
// call, which a schema function of the name takes ("EXTRACT"(1),
// verified); and a function, called from a query or DML, can run DDL in
// an autonomous transaction (verified).
func (b *builder) callsUserCode(n ast.Node) bool {
	calls := false
	ast.Inspect(n, func(m ast.Node) bool {
		if calls {
			return false
		}
		if fc, ok := m.(*ast.FuncCallExpr); ok {
			name := fc.FuncName
			calls = name == nil || name.Schema != "" || name.DBLink != "" || !sqlBuiltins[name.Name] || b.quoted(name.Loc)
		}
		return !calls
	})
	return calls || b.callsBareFunction(n)
}

// quoted reports whether the text at a location begins with a quoted
// identifier.
func (b *builder) quoted(loc ast.Loc) bool {
	return loc.Start >= 0 && loc.Start < len(b.w.sql) && b.w.sql[loc.Start] == '"'
}

// callsBareFunction reports whether a name of a statement may call a
// function without parentheses. Oracle takes a name for a column of a
// table in scope first, then for a function: SELECT f FROM dual, VALUES
// (f), SET a = f, WHERE f = 1, pkg.f, and shop.f call one, and WHERE f = 1
// on a table with a column F does not (verified). The walk judges a name
// only in a statement of one query block, where the tables in scope are
// those it names: a name there is a column the walk knows a table of it
// has, a column a table or alias qualifies, a sequence's NEXTVAL or
// CURRVAL, or a builtin Oracle takes first (verified: CURRENT_DATE,
// CURRENT_TIMESTAMP, LOCALTIMESTAMP, DBTIMEZONE, SESSIONTIMEZONE, and
// ORA_ROWSCN win over schema functions of their names). Any other name, or
// any name in a statement of more query blocks, may be a call.
func (b *builder) callsBareFunction(n ast.Node) bool {
	targets := make(map[*ast.ColumnRef]bool)
	var refs []*ast.ColumnRef
	blocks := 0
	ast.Inspect(n, func(m ast.Node) bool {
		switch v := m.(type) {
		case *ast.SelectStmt:
			blocks++
		case *ast.InsertStmt:
			markRefs(targets, v.Columns)
		case *ast.InsertIntoClause:
			markRefs(targets, v.Columns)
		case *ast.SetClause:
			targets[v.Column] = true
			markRefs(targets, v.Columns)
		case *ast.MergeClause:
			markRefs(targets, v.InsertCols)
		case *ast.IndexColumn:
			// A name an index lists alone is a column: one that is not
			// fails with 00904 and calls nothing, while a name in an
			// expression may call a function (verified).
			if c, ok := v.Expr.(*ast.ColumnRef); ok {
				targets[c] = true
			}
		case *ast.ColumnRef:
			refs = append(refs, v)
		}
		return true
	})
	var uses []*ast.ColumnRef
	for _, c := range refs {
		if !targets[c] {
			uses = append(uses, c)
		}
	}
	if len(uses) == 0 {
		return false
	}
	scope, ok := b.scopeOf(n, blocks)
	if !ok {
		return true
	}
	for _, c := range uses {
		if !b.resolves(c, scope) {
			return true
		}
	}
	return false
}

func markRefs(set map[*ast.ColumnRef]bool, l *ast.List) {
	if l == nil {
		return
	}
	for _, item := range l.Items {
		if c, ok := item.(*ast.ColumnRef); ok {
			set[c] = true
		}
	}
}

// scoped is a table in scope: its reference and its alias, "" for none.
type scoped struct {
	ref   objRef
	alias string
}

// scopeOf returns the tables in scope of a statement of one query block,
// false for a statement of more, or of a FROM item that is no table.
func (b *builder) scopeOf(n ast.Node, blocks int) ([]scoped, bool) {
	var scope []scoped
	add := func(o *ast.ObjectName, alias *ast.Alias) bool {
		r, ok := b.ref(o)
		if !ok {
			return false
		}
		s := scoped{ref: r}
		if alias != nil {
			s.alias = alias.Name
		}
		scope = append(scope, s)
		return true
	}
	var from func(e ast.TableExpr) bool
	from = func(e ast.TableExpr) bool {
		switch v := e.(type) {
		case *ast.TableRef:
			return v.Dblink == "" && add(v.Name, v.Alias)
		case *ast.JoinClause:
			return from(v.Left) && from(v.Right)
		}
		return false
	}
	fromList := func(l *ast.List) bool {
		if l == nil {
			return true
		}
		for _, item := range l.Items {
			e, ok := item.(ast.TableExpr)
			if !ok || !from(e) {
				return false
			}
		}
		return true
	}
	switch v := n.(type) {
	case *ast.SelectStmt:
		return scope, blocks == 1 && fromList(v.FromClause)
	case *ast.InsertStmt:
		switch {
		case v.MultiTable != nil:
			return nil, false
		case v.Select != nil:
			// The rows come from the query: VALUES sees no table.
			return scope, blocks == 1 && fromList(v.Select.FromClause)
		}
		return scope, blocks == 0
	case *ast.UpdateStmt:
		return scope, blocks == 0 && from(v.Target) && fromList(v.FromClause)
	case *ast.DeleteStmt:
		return scope, blocks == 0 && from(v.Target)
	case *ast.MergeStmt:
		return scope, blocks == 0 && add(v.Target, v.TargetAlias) && from(v.Source)
	case *ast.CreateIndexStmt:
		return scope, blocks == 0 && add(v.Table, nil)
	}
	return nil, false
}

// noParenBuiltins are the builtins the parser takes for names, which Oracle
// takes before a schema function of the name (verified).
var noParenBuiltins = map[string]bool{
	"CURRENT_DATE": true, "CURRENT_TIMESTAMP": true, "LOCALTIMESTAMP": true,
	"DBTIMEZONE": true, "SESSIONTIMEZONE": true, "ORA_ROWSCN": true, "UID": true,
}

// resolves reports whether a name of a statement certainly reads a column
// or a sequence, not a function.
func (b *builder) resolves(c *ast.ColumnRef, scope []scoped) bool {
	if c.Column == "NEXTVAL" || c.Column == "CURRVAL" {
		schema := c.Schema
		if schema == "" {
			schema = b.w.current
		}
		if c.Table != "" && schema != "" {
			f, _ := b.look(walk.Key{Space: spaceSequence, Parent: walk.Path(schema), Name: c.Table})
			if f == walk.Present {
				return true
			}
		}
	}
	switch {
	case c.Schema != "":
		for _, s := range scope {
			if s.alias == "" && s.ref.schema == c.Schema && s.ref.name == c.Table {
				return true
			}
		}
		return false
	case c.Table != "":
		for _, s := range scope {
			if s.alias == c.Table || s.alias == "" && s.ref.name == c.Table {
				return true
			}
		}
		return false
	}
	if noParenBuiltins[c.Column] {
		return true
	}
	for _, s := range scope {
		if f, _ := b.look(s.ref.childKey(spaceColumn, c.Column)); f == walk.Present {
			return true
		}
	}
	return false
}

// sqlBuiltins are Oracle SQL functions, as the parser folds their names.
//
// https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/About-SQL-Functions.html
var sqlBuiltins = func() map[string]bool {
	m := make(map[string]bool)
	for _, f := range strings.Fields(`
ABS ACOS ADD_MONTHS APPROX_COUNT APPROX_COUNT_DISTINCT APPROX_COUNT_DISTINCT_AGG APPROX_COUNT_DISTINCT_DETAIL
APPROX_MEDIAN APPROX_PERCENTILE APPROX_PERCENTILE_AGG APPROX_PERCENTILE_DETAIL APPROX_RANK APPROX_SUM ASCII ASCIISTR
ASIN ATAN ATAN2 AVG BFILENAME BIN_TO_NUM BITAND BIT_AND_AGG BIT_OR_AGG BIT_XOR_AGG BOOLEAN_AND_AGG BOOLEAN_OR_AGG
CARDINALITY CAST CEIL CHARTOROWID CHECKSUM CHR CLUSTER_ID CLUSTER_PROBABILITY CLUSTER_SET COALESCE COLLATION COLLECT
COMPOSE CON_DBID_TO_ID CON_GUID_TO_ID CON_ID_TO_CON_NAME CON_ID_TO_DBID CON_NAME_TO_ID CON_UID_TO_ID CONCAT CONVERT
CORR CORR_K CORR_S COS COSH COUNT COVAR_POP COVAR_SAMP CUME_DIST CURRENT_DATE CURRENT_TIMESTAMP CV DBTIMEZONE DECODE
DECOMPOSE DENSE_RANK DEPTH DEREF DUMP EMPTY_BLOB EMPTY_CLOB EXISTSNODE EXP EXTRACT EXTRACTVALUE FIRST FIRST_VALUE
FLOOR FROM_TZ FROM_VECTOR GREATEST GREATEST_UB GROUP_ID GROUPING GROUPING_ID HEXTORAW INITCAP INSTR INSTR2 INSTR4
INSTRB INSTRC ITERATION_NUMBER JSON JSON_ARRAY JSON_ARRAYAGG JSON_DATAGUIDE JSON_ID JSON_MERGEPATCH JSON_OBJECT
JSON_OBJECTAGG JSON_QUERY JSON_SCALAR JSON_SERIALIZE JSON_TABLE JSON_TRANSFORM JSON_VALUE KURTOSIS_POP KURTOSIS_SAMP
LAG LAST LAST_DAY LAST_VALUE LEAD LEAST LEAST_UB LENGTH LENGTH2 LENGTH4 LENGTHB LENGTHC LISTAGG LN LNNVL
LOCALTIMESTAMP LOG LOWER LPAD LTRIM MAKE_REF MAX MEDIAN MIN MOD MONTHS_BETWEEN NANVL NCHR NEW_TIME NEXT_DAY
NLS_CHARSET_DECL_LEN NLS_CHARSET_ID NLS_CHARSET_NAME NLS_COLLATION_ID NLS_COLLATION_NAME NLS_INITCAP NLS_LOWER
NLS_UPPER NLSSORT NTH_VALUE NTILE NULLIF NUMTODSINTERVAL NUMTOYMINTERVAL NVL NVL2 ORA_DM_PARTITION_NAME
ORA_DST_AFFECTED ORA_DST_CONVERT ORA_DST_ERROR ORA_HASH ORA_INVOKING_USER ORA_INVOKING_USERID PERCENT_RANK
PERCENTILE_CONT PERCENTILE_DISC POWER POWERMULTISET POWERMULTISET_BY_CARDINALITY RANK RATIO_TO_REPORT RAWTOHEX
RAWTONHEX REF REFTOHEX REGEXP_COUNT REGEXP_INSTR REGEXP_LIKE REGEXP_REPLACE REGEXP_SUBSTR REGR_AVGX REGR_AVGY
REGR_COUNT REGR_INTERCEPT REGR_R2 REGR_SLOPE REGR_SXX REGR_SXY REGR_SYY REMAINDER REPLACE ROUND ROUND_TIES_TO_EVEN
ROW_NUMBER ROWIDTOCHAR ROWIDTONCHAR RPAD RTRIM SCN_TO_TIMESTAMP SESSIONTIMEZONE SET SIGN SIN SINH SKEWNESS_POP
SKEWNESS_SAMP SOUNDEX SQRT STANDARD_HASH STATS_BINOMIAL_TEST STATS_CROSSTAB STATS_F_TEST STATS_KS_TEST STATS_MODE
STATS_MW_TEST STATS_ONE_WAY_ANOVA STATS_T_TEST_INDEP STATS_T_TEST_INDEPU STATS_T_TEST_ONE STATS_T_TEST_PAIRED
STATS_WSR_TEST STDDEV STDDEV_POP STDDEV_SAMP SUBSTR SUBSTR2 SUBSTR4 SUBSTRB SUBSTRC SUM SYS_CONNECT_BY_PATH
SYS_CONTEXT SYS_DBURIGEN SYS_EXTRACT_UTC SYS_GUID SYS_OP_ZONE_ID SYS_TYPEID SYS_XMLAGG SYS_XMLGEN SYSDATE
SYSTIMESTAMP TAN TANH TIME_BUCKET TIMESTAMP_TO_SCN TO_APPROX_COUNT_DISTINCT TO_APPROX_PERCENTILE TO_BINARY_DOUBLE
TO_BINARY_FLOAT TO_BLOB TO_BOOLEAN TO_CHAR TO_CLOB TO_DATE TO_DSINTERVAL TO_LOB TO_MULTI_BYTE TO_NCHAR TO_NCLOB
TO_NUMBER TO_SINGLE_BYTE TO_TIMESTAMP TO_TIMESTAMP_TZ TO_UTC_TIMESTAMP_TZ TO_VECTOR TO_YMINTERVAL TRANSLATE TREAT
TRIM TRUNC TZ_OFFSET UID UNISTR UPPER USER USERENV VALIDATE_CONVERSION VALUE VAR_POP VAR_SAMP VARIANCE VECTOR
VECTOR_DIMENSION_COUNT VECTOR_DIMENSION_FORMAT VECTOR_DIMS VECTOR_DISTANCE VECTOR_NORM VECTOR_SERIALIZE VSIZE
WIDTH_BUCKET XMLAGG XMLCAST XMLCDATA XMLCOLATTVAL XMLCOMMENT XMLCONCAT XMLDIFF XMLELEMENT XMLEXISTS XMLFOREST
XMLISVALID XMLPARSE XMLPATCH XMLPI XMLQUERY XMLROOT XMLSEQUENCE XMLSERIALIZE XMLTABLE XMLTRANSFORM`) {
		m[f] = true
	}
	return m
}()
