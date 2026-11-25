package parser

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Actions ALTER TABLE, ALTER VIEW, ALTER DYNAMIC TABLE and ALTER MATERIALIZED
// VIEW share: data metric functions, tags, policies and contacts.
//
//	https://docs.snowflake.com/en/sql-reference/sql/alter-table
//	https://docs.snowflake.com/en/sql-reference/sql/alter-view
//	https://docs.snowflake.com/en/sql-reference/sql/alter-dynamic-table
// ---------------------------------------------------------------------------

// dataMetricClauseCap bounds the comma-separated lists of a data metric
// function action so a malformed input can never spin.
const dataMetricClauseCap = 1024

// startsDataMetricFunction reports whether cur begins DATA METRIC FUNCTION
// (after ADD, DROP or MODIFY). DATA is a non-reserved keyword that can also be
// a column name, so METRIC must follow.
func (p *Parser) startsDataMetricFunction() bool {
	return p.curIs(kwDATA) && p.nextIs(kwMETRIC)
}

// startsDataMetricSchedule reports whether cur is DATA_METRIC_SCHEDULE (after
// SET or UNSET).
func (p *Parser) startsDataMetricSchedule() bool {
	return p.curIsWord("DATA_METRIC_SCHEDULE")
}

// parseDataMetricSchedule parses the part after SET or UNSET:
//
//	SET DATA_METRIC_SCHEDULE = '<schedule>'
//	UNSET DATA_METRIC_SCHEDULE
//
// cur is DATA_METRIC_SCHEDULE; start is the SET/UNSET offset.
func (p *Parser) parseDataMetricSchedule(set bool, start int) (*ast.DataMetricAction, error) {
	p.advance() // consume DATA_METRIC_SCHEDULE
	act := &ast.DataMetricAction{Kind: ast.DataMetricUnsetSchedule, Loc: ast.Loc{Start: start}}
	if set {
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		act.Kind = ast.DataMetricSetSchedule
		act.Schedule = tok.Str
	}
	act.Loc.End = p.prev.Loc.End
	return act, nil
}

// parseDataMetricFunctionAction parses a data metric function action after its
// ADD, DROP or MODIFY keyword; cur is DATA and start is the verb's offset.
//
//	ADD DATA METRIC FUNCTION <metric> ON ( <args> ) [ <options> ] [ , <metric> ON ( <args> ) [ <options> ] ... ]
//	DROP DATA METRIC FUNCTION <metric> ON ( <args> ) [ , <metric> ON ( <args> ) ... ]
//	MODIFY DATA METRIC FUNCTION <metric> ON ( <args> ) <change> [ , <metric> ON ( <args> ) <change> ... ]
//	MODIFY DATA METRIC FUNCTION <metric> SET <property> = <value> [ ... ]
func (p *Parser) parseDataMetricFunctionAction(kind ast.DataMetricActionKind, start int) (*ast.DataMetricAction, error) {
	p.advance() // consume DATA
	p.advance() // consume METRIC
	if !p.curIsWord("FUNCTION") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume FUNCTION
	act := &ast.DataMetricAction{Kind: kind, Loc: ast.Loc{Start: start}}
	for i := 0; ; i++ {
		if i >= dataMetricClauseCap {
			return nil, p.syntaxErrorAtCur()
		}
		clause, err := p.parseDataMetricFunctionClause(kind)
		if err != nil {
			return nil, err
		}
		act.Functions = append(act.Functions, clause)
		if !p.curIs(',') || !p.commaStartsDataMetricClause() {
			break
		}
		p.advance() // consume ','
	}
	act.Loc.End = p.prev.Loc.End
	return act, nil
}

// commaStartsDataMetricClause reports whether the ',' at cur is followed by
// another `<metric> ON` (or `<metric> SET`) clause rather than by something
// that belongs to the enclosing statement.
func (p *Parser) commaStartsDataMetricClause() bool {
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance() // consume ','
	if _, err := p.parseObjectName(); err != nil {
		return false
	}
	return p.curIs(kwON) || p.curIs(kwSET)
}

// commaStartsExpectation reports whether the ',' at cur is followed by another
// `<name> ( <expr> )` expectation rather than by another metric clause.
func (p *Parser) commaStartsExpectation() bool {
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance() // consume ','
	if _, err := p.parseIdent(); err != nil {
		return false
	}
	return p.curIs('(')
}

// parseDataMetricFunctionClause parses one `<metric> ...` clause of a data
// metric function action.
func (p *Parser) parseDataMetricFunctionClause(kind ast.DataMetricActionKind) (*ast.DataMetricFunctionClause, error) {
	clause := &ast.DataMetricFunctionClause{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	clause.Name = name

	// MODIFY DATA METRIC FUNCTION <metric> SET
	//   { ANOMALY_DETECTION = ... | SENSITIVITY = ... | DATA_QUALITY_NOTIFICATION = ... } [ ... ]
	if kind == ast.DataMetricModify && p.curIs(kwSET) {
		p.advance() // consume SET
		var opts []*ast.CopyOption
		for p.curIsWord("ANOMALY_DETECTION") || p.curIsWord("SENSITIVITY") || p.curIsWord("DATA_QUALITY_NOTIFICATION") {
			opt, err := p.parseAssignment()
			if err != nil {
				return nil, err
			}
			opts = append(opts, opt)
		}
		if len(opts) == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		clause.Modify = ast.DataMetricModifySet
		clause.Options = opts
		clause.Loc.End = p.prev.Loc.End
		return clause, nil
	}

	if err := p.parseDataMetricArgs(clause); err != nil {
		return nil, err
	}

	switch kind {
	case ast.DataMetricAdd:
		if err := p.parseDataMetricAddOptions(clause); err != nil {
			return nil, err
		}
	case ast.DataMetricModify:
		if err := p.parseDataMetricModifyChange(clause); err != nil {
			return nil, err
		}
	}
	clause.Loc.End = p.prev.Loc.End
	return clause, nil
}

// parseDataMetricArgs parses ON ( [ <col> [ , ... ] ] [ , TABLE <table> ( <col> [ , ... ] ) ] ).
// The argument list may be empty: a table-level metric such as ROW_COUNT takes
// ON ( ).
func (p *Parser) parseDataMetricArgs(clause *ast.DataMetricFunctionClause) error {
	if _, err := p.expect(kwON); err != nil {
		return err
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	clause.HasOn = true
	for i := 0; i > 0 || !p.curIs(')'); i++ {
		if i >= dataMetricClauseCap {
			return p.syntaxErrorAtCur()
		}
		if p.curIs(kwTABLE) {
			p.advance() // consume TABLE
			ref, err := p.parseObjectName()
			if err != nil {
				return err
			}
			cols, err := p.parseIdentListInParens()
			if err != nil {
				return err
			}
			clause.RefTable = ref
			clause.RefColumns = cols
			break // the referenced table closes the list
		}
		col, err := p.parseIdent()
		if err != nil {
			return err
		}
		clause.Columns = append(clause.Columns, col)
		if !p.curIs(',') {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(')')
	return err
}

// parseDataMetricAddOptions parses the options of an ADD clause, in the order
// the page lists them:
//
//	[ WITHIN GROUP ( <col> [ , ... ] ) [ GROUP LIMIT <n> ] ]
//	[ FILTER ( <expr> ) ]
//	[ EXPECTATION <name> ( <expr> ) [ , <name> ( <expr> ) ... ] ]
//	[ EXECUTE AS ROLE <role> ]
//	[ ANOMALY_DETECTION = { TRUE | FALSE } ]
//	[ SENSITIVITY = { 'LOW' | 'MEDIUM' | 'HIGH' } ]
func (p *Parser) parseDataMetricAddOptions(clause *ast.DataMetricFunctionClause) error {
	if p.curIs(kwWITHIN) {
		p.advance() // consume WITHIN
		if _, err := p.expect(kwGROUP); err != nil {
			return err
		}
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return err
		}
		clause.WithinGroup = cols
		if p.curIs(kwGROUP) && p.nextIs(kwLIMIT) {
			p.advance() // consume GROUP
			p.advance() // consume LIMIT
			tok, err := p.expect(tokInt)
			if err != nil {
				return err
			}
			n := tok.Ival
			clause.GroupLimit = &n
		}
	}
	if p.curIsWord("FILTER") && p.nextIs('(') {
		if err := p.parseDataMetricFilter(clause, false); err != nil {
			return err
		}
	}
	if p.curIsWord("EXPECTATION") {
		p.advance() // consume EXPECTATION
		exps, err := p.parseDataMetricExpectations(true)
		if err != nil {
			return err
		}
		clause.Expectations = exps
	}
	if p.curIsWord("EXECUTE") && p.nextIs(kwAS) {
		p.advance() // consume EXECUTE
		p.advance() // consume AS
		if _, err := p.expect(kwROLE); err != nil {
			return err
		}
		role, err := p.parseIdent()
		if err != nil {
			return err
		}
		clause.ExecuteAs = role
	}
	for _, word := range []string{"ANOMALY_DETECTION", "SENSITIVITY"} {
		if p.curIsWord(word) {
			opt, err := p.parseAssignment()
			if err != nil {
				return err
			}
			clause.Options = append(clause.Options, opt)
		}
	}
	return nil
}

// parseDataMetricModifyChange parses what a MODIFY clause changes:
//
//	{ SUSPEND | RESUME }
//	{ ADD | MODIFY } EXPECTATION <name> ( <expr> ) [ , ... ]
//	DROP EXPECTATION <name> [ , ... ]
//	FILTER ( [ <expr> ] )
func (p *Parser) parseDataMetricModifyChange(clause *ast.DataMetricFunctionClause) error {
	switch {
	case p.curIs(kwSUSPEND):
		p.advance()
		clause.Modify = ast.DataMetricModifySuspend
	case p.curIs(kwRESUME):
		p.advance()
		clause.Modify = ast.DataMetricModifyResume
	case p.curIs(kwADD) || p.curIs(kwMODIFY) || p.curIs(kwDROP):
		verb := p.cur.Type
		p.advance() // consume ADD / MODIFY / DROP
		if !p.curIsWord("EXPECTATION") {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume EXPECTATION
		exps, err := p.parseDataMetricExpectations(verb != kwDROP)
		if err != nil {
			return err
		}
		clause.Expectations = exps
		switch verb {
		case kwADD:
			clause.Modify = ast.DataMetricModifyAddExpectation
		case kwMODIFY:
			clause.Modify = ast.DataMetricModifyModifyExpectation
		default:
			clause.Modify = ast.DataMetricModifyDropExpectation
		}
	case p.curIsWord("FILTER"):
		if err := p.parseDataMetricFilter(clause, true); err != nil {
			return err
		}
		clause.Modify = ast.DataMetricModifyFilter
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseDataMetricFilter parses FILTER ( <expr> ); with allowEmpty, MODIFY's
// FILTER ( ) that removes the filter. cur is FILTER.
func (p *Parser) parseDataMetricFilter(clause *ast.DataMetricFunctionClause, allowEmpty bool) error {
	p.advance() // consume FILTER
	if _, err := p.expect('('); err != nil {
		return err
	}
	clause.HasFilter = true
	if !(allowEmpty && p.curIs(')')) {
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		clause.Filter = expr
	}
	_, err := p.expect(')')
	return err
}

// parseDataMetricExpectations parses `<name> ( <expr> ) [ , ... ]`, or with
// withExpr false the bare names of DROP EXPECTATION. A comma continues the list
// only when another expectation follows, so `, <metric> ON (...)` is left to the
// caller.
func (p *Parser) parseDataMetricExpectations(withExpr bool) ([]*ast.DataMetricExpectation, error) {
	var exps []*ast.DataMetricExpectation
	for i := 0; ; i++ {
		if i >= dataMetricClauseCap {
			return nil, p.syntaxErrorAtCur()
		}
		start := p.cur.Loc.Start
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		exp := &ast.DataMetricExpectation{Name: name}
		if withExpr {
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			exp.Expr = expr
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		exp.Loc = ast.Loc{Start: start, End: p.prev.Loc.End}
		exps = append(exps, exp)
		if !p.curIs(',') {
			break
		}
		if withExpr && !p.commaStartsExpectation() {
			break
		}
		p.advance() // consume ','
	}
	return exps, nil
}

// ---------------------------------------------------------------------------
// KEY = value assignments
// ---------------------------------------------------------------------------

// parseAssignment parses one `KEY = <value>` property through parseCopyOption
// and rejects the bare (value-less) form, which no ALTER property takes.
//
// A property whose value the pages document (propertyValues) must have one of
// those values; any other property keeps the open-ended value forms.
func (p *Parser) parseAssignment() (*ast.CopyOption, error) {
	if !p.startsCopyOption() || !p.nextIs('=') {
		return nil, p.syntaxErrorAtCur()
	}
	if err := p.checkPropertyValue(); err != nil {
		return nil, err
	}
	return p.parseCopyOption()
}

// propertyValue is the documented value of an ALTER property: one of words
// (unquoted), one of strs (single-quoted), or, with anyString, any string.
type propertyValue struct {
	words     []string
	strs      []string
	anyString bool
}

var (
	boolValue   = propertyValue{words: []string{"TRUE", "FALSE"}}
	stringValue = propertyValue{anyString: true}
)

// propertyValues holds the values the ALTER TABLE, ALTER ICEBERG TABLE (incl.
// CONVERT TO MANAGED) and data metric function pages give their properties.
var propertyValues = map[string]propertyValue{
	// alter-table
	"ICEBERG_DEFAULT_DDL_COLLATION": stringValue,
	"ENABLE_SCHEMA_EVOLUTION":       boolValue,
	"ERROR_LOGGING":                 boolValue,
	"ROW_TIMESTAMP":                 boolValue,
	// alter-iceberg-table
	"REPLACE_INVALID_CHARACTERS":     boolValue,
	"CATALOG_SYNC":                   stringValue,
	"AUTO_REFRESH":                   boolValue,
	"TARGET_FILE_SIZE":               {words: []string{"AUTO"}, strs: []string{"AUTO", "16MB", "32MB", "64MB", "128MB"}},
	"PATH_LAYOUT":                    {words: []string{"FLAT", "HIERARCHICAL"}},
	"STORAGE_SERIALIZATION_POLICY":   {words: []string{"COMPATIBLE", "OPTIMIZED"}},
	"LOG_EVENT_LEVEL":                {words: []string{"ERROR", "WARN", "INFO", "DEBUG"}},
	"ENABLE_DATA_COMPACTION":         boolValue,
	"ICEBERG_MERGE_ON_READ_BEHAVIOR": {strs: []string{"AUTO", "ENABLED", "DISABLED"}},
	"ENABLE_ICEBERG_MERGE_ON_READ":   boolValue,
	"MATCH_BY_COLUMN_NAME":           {words: []string{"CASE_INSENSITIVE", "CASE_SENSITIVE"}},
	"ALLOW_V3_EXPRESSION_DEFAULTS":   boolValue,
	"TIME_TRAVEL_TIMESTAMP_SOURCE":   {strs: []string{"SNOWFLAKE_METADATA", "ICEBERG_METADATA"}},
	// alter-iceberg-table-convert-to-managed
	"FILE_PATH":                stringValue,
	"BASE_LOCATION":            stringValue,
	"DROP_FROM_SOURCE_CATALOG": boolValue,
	// data metric functions
	"ANOMALY_DETECTION":         boolValue,
	"SENSITIVITY":               {strs: []string{"LOW", "MEDIUM", "HIGH"}},
	"DATA_QUALITY_NOTIFICATION": boolValue,
}

// tableSetProperties and icebergTableSetProperties are the SET properties the
// ALTER TABLE and ALTER ICEBERG TABLE pages list beyond the ones
// parseOneTableProp names itself (DATA_RETENTION_TIME_IN_DAYS, COMMENT, ...):
//
//	https://docs.snowflake.com/en/sql-reference/sql/alter-table
//	https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table
var (
	tableSetProperties = map[string]bool{
		"ICEBERG_DEFAULT_DDL_COLLATION": true,
		"ENABLE_SCHEMA_EVOLUTION":       true,
		"ERROR_LOGGING":                 true,
		"ROW_TIMESTAMP":                 true,
	}
	icebergTableSetProperties = map[string]bool{
		"REPLACE_INVALID_CHARACTERS":     true,
		"CATALOG_SYNC":                   true,
		"AUTO_REFRESH":                   true,
		"TARGET_FILE_SIZE":               true,
		"PATH_LAYOUT":                    true,
		"STORAGE_SERIALIZATION_POLICY":   true,
		"LOG_EVENT_LEVEL":                true,
		"ERROR_LOGGING":                  true,
		"ENABLE_DATA_COMPACTION":         true,
		"ICEBERG_MERGE_ON_READ_BEHAVIOR": true,
		"ENABLE_ICEBERG_MERGE_ON_READ":   true,
		"ENABLE_SCHEMA_EVOLUTION":        true,
		"MATCH_BY_COLUMN_NAME":           true,
		"ALLOW_V3_EXPRESSION_DEFAULTS":   true,
		"TIME_TRAVEL_TIMESTAMP_SOURCE":   true,
	}
)

// checkPropertyValue checks, without consuming anything, that the `KEY =
// <value>` at cur gives a documented property one of its documented values.
func (p *Parser) checkPropertyValue() error {
	spec, ok := propertyValues[strings.ToUpper(p.cur.Str)]
	if !ok {
		return nil
	}
	snap := p.snapshot()
	p.advance() // consume the name
	p.advance() // consume '='
	val := p.cur
	next := p.peekNext().Type
	if p.curIs(tokString) {
		if spec.anyString || slices.Contains(spec.strs, strings.ToUpper(val.Str)) {
			p.restore(snap)
			return nil
		}
	} else if val.Type != tokQuotedIdent && val.Str != "" && next != '.' &&
		slices.Contains(spec.words, strings.ToUpper(val.Str)) {
		p.restore(snap)
		return nil
	}
	err := p.syntaxErrorAtCur()
	p.restore(snap)
	return err
}

// startsAssignment reports whether cur begins a `KEY = <value>` property.
func (p *Parser) startsAssignment() bool {
	return p.startsCopyOption() && p.nextIs('=')
}

// ---------------------------------------------------------------------------
// Tags
// ---------------------------------------------------------------------------

// parseSetTagClause parses the TAG part of SET TAG, unparenthesized as the
// ALTER TABLE / VIEW pages write it, or parenthesized as omni has always read
// it:
//
//	TAG <tag> = '<value>' [ , <tag> = '<value>' ... ]
//	TAG ( <tag> = '<value>' [ , ... ] )
//
// In the unparenthesized form a comma continues the list only when another
// `<tag> =` follows, so a column list (`..., <col> SET TAG ...`) is left to the
// caller. cur is TAG.
func (p *Parser) parseSetTagClause() ([]*ast.TagAssignment, error) {
	if p.nextIs('(') {
		return p.parseTagAssignments()
	}
	p.advance() // consume TAG
	var tags []*ast.TagAssignment
	for {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		val, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &ast.TagAssignment{Name: name, Value: val.Str})
		if !p.curIs(',') || !p.commaStartsTagAssignment() {
			return tags, nil
		}
		p.advance() // consume ','
	}
}

// parseUnsetTagClause parses the TAG part of UNSET TAG, unparenthesized or
// parenthesized (see parseSetTagClause). In the unparenthesized form a comma
// followed by `<col> { SET | UNSET }` belongs to the caller's column list. cur
// is TAG.
func (p *Parser) parseUnsetTagClause() ([]*ast.ObjectName, error) {
	if p.nextIs('(') {
		return p.parseUnsetTagList()
	}
	p.advance() // consume TAG
	var names []*ast.ObjectName
	for {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if !p.curIs(',') || !p.commaStartsUnsetTagName() {
			return names, nil
		}
		p.advance() // consume ','
	}
}

// commaStartsTagAssignment reports whether the ',' at cur is followed by
// `<tag> =`.
func (p *Parser) commaStartsTagAssignment() bool {
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance() // consume ','
	if _, err := p.parseObjectName(); err != nil {
		return false
	}
	return p.curIs('=')
}

// commaStartsUnsetTagName reports whether the ',' at cur is followed by another
// tag name rather than by `<col> { SET | UNSET }`.
func (p *Parser) commaStartsUnsetTagName() bool {
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance() // consume ','
	if p.curIs(kwCOLUMN) {
		// `, COLUMN <col> ...`: the next column action. (parseObjectName would
		// take the reserved COLUMN as a name part.)
		return false
	}
	if _, err := p.parseObjectName(); err != nil {
		return false
	}
	return !p.curIs(kwSET) && !p.curIs(kwUNSET)
}

// ---------------------------------------------------------------------------
// Policies and contacts
// ---------------------------------------------------------------------------

// parseSetAggregationPolicy parses
//
//	AGGREGATION POLICY <name> [ ENTITY KEY ( <col> [ , ... ] ) ] [ FORCE ]
//
// after SET; cur is AGGREGATION (an identifier).
func (p *Parser) parseSetAggregationPolicy() (*ast.ObjectName, []ast.Ident, bool, error) {
	p.advance() // consume AGGREGATION
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, nil, false, err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, nil, false, err
	}
	var keys []ast.Ident
	if p.curIsWord("ENTITY") {
		p.advance() // consume ENTITY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, nil, false, err
		}
		keys, err = p.parseIdentListInParens()
		if err != nil {
			return nil, nil, false, err
		}
	}
	force := false
	if p.curIs(kwFORCE) {
		p.advance() // consume FORCE
		force = true
	}
	return name, keys, force, nil
}

// startsPolicyWord reports whether cur is the given policy-kind word (an
// identifier such as AGGREGATION or PROJECTION) followed by POLICY.
func (p *Parser) startsPolicyWord(word string) bool {
	return p.curIsWord(word) && p.nextIs(kwPOLICY)
}

// expectPolicyWord consumes `<word> POLICY`; the caller checked startsPolicyWord.
func (p *Parser) expectPolicyWord() error {
	p.advance() // consume the kind word
	_, err := p.expect(kwPOLICY)
	return err
}

// parseContactValue parses the value of a CONTACT property after the CONTACT
// word, either form the pages write:
//
//	CONTACT <purpose> = <contact> [ , <purpose> = <contact> ... ]
//	CONTACT ( <purpose> = <contact> [ , ... ] )
//
// and returns its verbatim source. cur is CONTACT.
func (p *Parser) parseContactValue() (string, error) {
	p.advance() // consume CONTACT
	start := p.cur.Loc.Start
	paren := p.curIs('(')
	if paren {
		p.advance() // consume '('
	}
	for {
		if _, err := p.parseIdent(); err != nil {
			return "", err
		}
		if _, err := p.expect('='); err != nil {
			return "", err
		}
		if _, err := p.parseObjectName(); err != nil {
			return "", err
		}
		if !p.curIs(',') {
			break
		}
		if !paren && !p.commaStartsContactEntry() {
			break
		}
		p.advance() // consume ','
	}
	if paren {
		if _, err := p.expect(')'); err != nil {
			return "", err
		}
	}
	return p.srcSlice(start, p.prev.Loc.End), nil
}

// commaStartsContactEntry reports whether the ',' at cur is followed by
// another `<purpose> =` entry.
func (p *Parser) commaStartsContactEntry() bool {
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance() // consume ','
	if _, err := p.parseIdent(); err != nil {
		return false
	}
	return p.curIs('=')
}

// parseUnsetContact parses CONTACT <purpose> after UNSET and returns
// "CONTACT <purpose>". cur is CONTACT.
func (p *Parser) parseUnsetContact() (string, error) {
	p.advance() // consume CONTACT
	purpose, err := p.parseIdent()
	if err != nil {
		return "", err
	}
	return "CONTACT " + purpose.String(), nil
}

// ---------------------------------------------------------------------------
// Row access policies and column policy actions
// ---------------------------------------------------------------------------

// parseAddRowAccessPolicyTail parses ROW ACCESS POLICY <p> ON ( <col> [ , ... ] )
// after ADD.
func (p *Parser) parseAddRowAccessPolicyTail() (*ast.ObjectName, []ast.Ident, error) {
	for _, t := range []int{kwROW, kwACCESS, kwPOLICY} {
		if _, err := p.expect(t); err != nil {
			return nil, nil, err
		}
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, nil, err
	}
	if _, err := p.expect(kwON); err != nil {
		return nil, nil, err
	}
	cols, err := p.parseIdentListInParens()
	if err != nil {
		return nil, nil, err
	}
	return name, cols, nil
}

// parseDropRowAccessPolicyTail parses, after DROP,
//
//	ROW ACCESS POLICY <p> [ , ADD ROW ACCESS POLICY <p2> ON ( <col> [ , ... ] ) ]
//	ALL ROW ACCESS POLICIES
//
// and reports all for the second form. With legacy, DROP ROW ACCESS POLICIES
// (no ALL), which omni has always read for views, is the second form too.
func (p *Parser) parseDropRowAccessPolicyTail(legacy bool) (all bool, name, addName *ast.ObjectName, addCols []ast.Ident, err error) {
	if p.curIs(kwALL) {
		p.advance() // consume ALL
		for _, t := range []int{kwROW, kwACCESS, kwPOLICIES} {
			if _, err := p.expect(t); err != nil {
				return false, nil, nil, nil, err
			}
		}
		return true, nil, nil, nil, nil
	}
	for _, t := range []int{kwROW, kwACCESS} {
		if _, err := p.expect(t); err != nil {
			return false, nil, nil, nil, err
		}
	}
	if legacy && p.curIs(kwPOLICIES) {
		p.advance()
		return true, nil, nil, nil, nil
	}
	if _, err := p.expect(kwPOLICY); err != nil {
		return false, nil, nil, nil, err
	}
	if name, err = p.parseObjectName(); err != nil {
		return false, nil, nil, nil, err
	}
	// The documented drop-and-add combination.
	if p.curIs(',') && p.nextIs(kwADD) {
		p.advance() // consume ','
		p.advance() // consume ADD
		if addName, addCols, err = p.parseAddRowAccessPolicyTail(); err != nil {
			return false, nil, nil, nil, err
		}
	}
	return false, name, addName, addCols, nil
}

// parseColumnPolicyActions parses the comma-separated column action list after
// ALTER / MODIFY that ALTER VIEW and ALTER DYNAMIC TABLE share (see
// parseAlterViewColumnAction).
func (p *Parser) parseColumnPolicyActions() ([]*ast.AlterViewColumnAction, error) {
	var actions []*ast.AlterViewColumnAction
	for i := 0; ; i++ {
		if i >= dataMetricClauseCap {
			return nil, p.syntaxErrorAtCur()
		}
		action, err := p.parseAlterViewColumnAction()
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
		if !p.curIs(',') {
			return actions, nil
		}
		p.advance() // consume ','
	}
}
