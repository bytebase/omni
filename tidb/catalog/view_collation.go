package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Collation derivation for view bodies.
//
// TiDB derives a charset, a collation, a coercibility and a repertoire for
// every expression while it builds a plan, and CREATE VIEW builds the plan of
// its body, so a body whose arguments do not aggregate is refused with 1267,
// 1270 or 1271. The rules here are pingcap/tidb v8.5.5
// pkg/expression/collation.go: deriveCollation picks the arguments of a
// function that take part, inferCollation aggregates them, safeConvert checks
// each conversion, and CheckAndDeriveCollationFromExprs turns a failure into
// the error.
//
// They are not MySQL's rules. Two IMPLICIT collations of one charset, neither
// of them a _bin collation, aggregate to the charset's _bin collation with
// coercibility NONE instead of failing, and only a function that returns no
// string (a comparison, LIKE, STRCMP, ...) refuses NONE. A _bin collation beats
// another collation of its charset at the same coercibility.
//
// Only the expressions inferExprCollation lists are modelled. Anything else
// has no answer (known is false) and never fails, and neither does an
// aggregation one of whose arguments has no answer, so the catalog refuses
// only what TiDB refuses.

// coercibility is TiDB's expression.Coercibility; lower wins.
type coercibility int

const (
	coercibilityExplicit  coercibility = iota // COLLATE clause
	coercibilityNone                          // two collations neither of which wins
	coercibilityImplicit                      // a column, CONVERT ... USING
	coercibilitySysconst                      // USER(), VERSION(), ...
	coercibilityCoercible                     // a string literal
	coercibilityNumeric                       // a number or temporal value
	coercibilityIgnorable                     // NULL
)

var coercibilityNames = [...]string{"EXPLICIT", "NONE", "IMPLICIT", "SYSCONST", "COERCIBLE", "NUMERIC", "IGNORABLE"}

// repertoire is TiDB's expression.Repertoire.
type repertoire uint8

const (
	repertoireASCII    repertoire = 0x01
	repertoireExtended repertoire = 0x02
	repertoireUnicode             = repertoireASCII | repertoireExtended
)

// connectionCharset and connectionCollation are TiDB v8.5.5's defaults for
// character_set_connection and collation_connection, which a string literal
// without an introducer and a string built from numbers take.
const (
	connectionCharset   = "utf8mb4"
	connectionCollation = "utf8mb4_bin"
)

// exprCollation is what TiDB derives for one expression.
type exprCollation struct {
	known     bool
	charset   string
	collation string
	coer      coercibility
	repe      repertoire
	// constant marks an expression TiDB folds to a constant. safeConvert
	// converts a constant by value: value holds it for a string literal and is
	// nil when the catalog does not know it.
	constant bool
	isNull   bool
	value    *string
}

func (e exprCollation) numeric() bool { return e.coer == coercibilityNumeric }

var (
	numericCollation = exprCollation{known: true, charset: "binary", collation: "binary", coer: coercibilityNumeric, repe: repertoireASCII}
	nullCollation    = exprCollation{known: true, charset: "binary", collation: "binary", coer: coercibilityIgnorable, repe: repertoireASCII, constant: true, isNull: true}
)

// viewScope resolves the column references of one SELECT.
type viewScope struct {
	relations map[string]viewRelation
	// opaque is set when the FROM clause has a relation the catalog cannot
	// list the columns of (a derived table, a CTE, a missing table), so an
	// unqualified name may be one of its columns.
	opaque bool
}

// viewRelation lists the derived collation of each column of a FROM item.
type viewRelation map[string]exprCollation

// newViewScope lists the FROM items of sel. ctes holds the CTE names an
// enclosing query makes visible, which a WITH on a set operation or a
// parenthesized query attaches to the outer node.
func (c *Catalog) newViewScope(sel *nodes.SelectStmt, db *Database, ctes map[string]bool) *viewScope {
	s := &viewScope{relations: map[string]viewRelation{}}
	ctes = withCTENames(sel, ctes)
	for _, from := range sel.From {
		c.collectViewRelations(from, db, ctes, s)
	}
	return s
}

// withCTENames adds the names of sel's CTEs to the visible ones, copying
// rather than changing ctes.
func withCTENames(sel *nodes.SelectStmt, ctes map[string]bool) map[string]bool {
	if len(sel.CTEs) == 0 {
		return ctes
	}
	out := make(map[string]bool, len(ctes)+len(sel.CTEs))
	for name := range ctes {
		out[name] = true
	}
	for _, cte := range sel.CTEs {
		out[toLower(cte.Name)] = true
	}
	return out
}

func (c *Catalog) collectViewRelations(te nodes.TableExpr, db *Database, ctes map[string]bool, s *viewScope) {
	switch n := te.(type) {
	case *nodes.TableRef:
		key := n.Name
		if n.Alias != "" {
			key = n.Alias
		}
		rel := c.viewRelationFor(n, db, ctes)
		if rel == nil {
			s.opaque = true
		}
		s.relations[toLower(key)] = rel
	case *nodes.JoinClause:
		c.collectViewRelations(n.Left, db, ctes, s)
		c.collectViewRelations(n.Right, db, ctes, s)
	default:
		s.opaque = true
	}
}

func (c *Catalog) viewRelationFor(ref *nodes.TableRef, db *Database, ctes map[string]bool) viewRelation {
	if ref.Schema == "" && ctes[toLower(ref.Name)] {
		return nil
	}
	if ref.Schema != "" {
		db = c.GetDatabase(ref.Schema)
	}
	if db == nil {
		return nil
	}
	tbl := db.GetTable(ref.Name)
	if tbl == nil {
		// A view keeps the coercibility of the expression each column
		// selects, which the catalog does not store, so its columns stay
		// unresolved.
		return nil
	}
	rel := viewRelation{}
	for _, col := range tbl.Columns {
		rel[toLower(col.Name)] = columnCollation(col)
	}
	return rel
}

func (s *viewScope) resolve(cr *nodes.ColumnRef) exprCollation {
	if cr == nil || cr.Star {
		return exprCollation{}
	}
	col := toLower(cr.Column)
	if cr.Table != "" {
		rel := s.relations[toLower(cr.Table)]
		if rel == nil {
			return exprCollation{}
		}
		return rel[col]
	}
	if s.opaque {
		return exprCollation{}
	}
	var out exprCollation
	found := false
	for _, rel := range s.relations {
		ec, ok := rel[col]
		if !ok {
			continue
		}
		if found {
			return exprCollation{}
		}
		found, out = true, ec
	}
	return out
}

// columnCollation is deriveCoercibilityForColumn with the column's type.
func columnCollation(col *Column) exprCollation {
	switch col.DataType {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "enum", "set":
		cs, coll := normalizeCharsetName(col.Charset), normalizeCollationName(col.Collation)
		if coll == "" {
			coll = defaultCollationForCharset[cs]
		}
		if cs == "" {
			cs = tidbCollationCharset[coll]
		}
		if cs == "" || coll == "" {
			return exprCollation{}
		}
		repe := repertoireUnicode
		if cs == "ascii" {
			repe = repertoireASCII
		}
		return exprCollation{known: true, charset: cs, collation: coll, coer: coercibilityImplicit, repe: repe}
	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		return exprCollation{known: true, charset: "binary", collation: "binary", coer: coercibilityImplicit, repe: repertoireUnicode}
	case "tinyint", "smallint", "mediumint", "int", "bigint", "decimal", "float", "double",
		"date", "datetime", "timestamp", "time", "year":
		return numericCollation
	}
	return exprCollation{}
}

// checkViewSelectCollations derives the collation of every expression TiDB
// builds for the view body that the catalog models, and returns the error
// TiDB answers when two of them do not aggregate. db is the database an
// unqualified table names (viewBodyDatabase); ctes is as for newViewScope.
func (c *Catalog) checkViewSelectCollations(sel *nodes.SelectStmt, db *Database, ctes map[string]bool) error {
	if sel == nil {
		return nil
	}
	if sel.ParenSource != nil {
		return c.checkViewSelectCollations(sel.ParenSource, db, withCTENames(sel, ctes))
	}
	if sel.Left != nil || sel.Right != nil {
		ctes = withCTENames(sel, ctes)
		if err := c.checkViewSelectCollations(sel.Left, db, ctes); err != nil {
			return err
		}
		return c.checkViewSelectCollations(sel.Right, db, ctes)
	}
	s := c.newViewScope(sel, db, ctes)
	for _, from := range sel.From {
		if err := s.checkJoinConditions(from); err != nil {
			return err
		}
	}
	for _, target := range sel.TargetList {
		expr := target
		if rt, ok := target.(*nodes.ResTarget); ok {
			expr = rt.Val
		}
		if _, err := s.infer(expr); err != nil {
			return err
		}
	}
	if sel.Where != nil {
		if _, err := s.infer(sel.Where); err != nil {
			return err
		}
	}
	return nil
}

func (s *viewScope) checkJoinConditions(te nodes.TableExpr) error {
	j, ok := te.(*nodes.JoinClause)
	if !ok {
		return nil
	}
	if err := s.checkJoinConditions(j.Left); err != nil {
		return err
	}
	if err := s.checkJoinConditions(j.Right); err != nil {
		return err
	}
	if on, ok := j.Condition.(*nodes.OnCondition); ok && on.Expr != nil {
		if _, err := s.infer(on.Expr); err != nil {
			return err
		}
	}
	return nil
}

// viewTargetCollations derives the collation of each select item of a view
// body that is a single SELECT, or returns nil. A star or an item the catalog
// does not model leaves its entry unknown.
func (c *Catalog) viewTargetCollations(sel *nodes.SelectStmt, db *Database) []exprCollation {
	var ctes map[string]bool
	for sel != nil && sel.ParenSource != nil {
		ctes = withCTENames(sel, ctes)
		sel = sel.ParenSource
	}
	if sel == nil || sel.Left != nil || sel.Right != nil {
		return nil
	}
	s := c.newViewScope(sel, db, ctes)
	out := make([]exprCollation, len(sel.TargetList))
	for i, target := range sel.TargetList {
		// An item with no alias is the bare expression.
		expr := target
		if rt, ok := target.(*nodes.ResTarget); ok {
			expr = rt.Val
		}
		switch e := expr.(type) {
		case *nodes.StarExpr:
			return nil
		case *nodes.ColumnRef:
			if e.Star {
				return nil
			}
		}
		if ec, err := s.infer(expr); err == nil {
			out[i] = ec
		}
	}
	return out
}

// viewColumnMetadata names the derived collations after the view's columns.
// A column that is not a string has neither charset nor collation.
func viewColumnMetadata(names []string, collations []exprCollation) []ViewColumn {
	cols := make([]ViewColumn, len(names))
	for i, name := range names {
		cols[i].Name = name
		if len(collations) != len(names) {
			continue
		}
		ec := collations[i]
		if !ec.known || ec.numeric() || ec.coer == coercibilityIgnorable {
			continue
		}
		cols[i].Charset, cols[i].Collation = ec.charset, ec.collation
	}
	return cols
}

// infer derives the collation of expr and returns the error TiDB answers for
// it or any expression inside it.
func (s *viewScope) infer(expr nodes.ExprNode) (exprCollation, error) {
	switch e := expr.(type) {
	case nil:
		return exprCollation{}, nil
	case *nodes.ColumnRef:
		return s.resolve(e), nil
	case *nodes.StringLit:
		return stringLiteralCollation(e), nil
	case *nodes.NullLit:
		return nullCollation, nil
	case *nodes.IntLit, *nodes.FloatLit, *nodes.BoolLit:
		ec := numericCollation
		ec.constant = true
		return ec, nil
	case *nodes.ParenExpr:
		return s.infer(e.Expr)
	case *nodes.CollateExpr:
		child, err := s.infer(e.Expr)
		if err != nil || !child.known {
			return exprCollation{}, err
		}
		coll := normalizeCollationName(e.Collation)
		cs, ok := tidbCollationCharset[coll]
		if !ok {
			return exprCollation{}, nil
		}
		child.charset, child.collation, child.coer = cs, coll, coercibilityExplicit
		return child, nil
	case *nodes.ConvertExpr:
		child, err := s.infer(e.Expr)
		if err != nil || e.Charset == "" || !child.known {
			return exprCollation{}, err
		}
		cs := normalizeCharsetName(e.Charset)
		coll, ok := defaultCollationForCharset[cs]
		if !ok {
			return exprCollation{}, nil
		}
		repe := repertoireUnicode
		if cs == "ascii" {
			repe = repertoireASCII
		}
		return exprCollation{known: true, charset: cs, collation: coll, coer: coercibilityImplicit, repe: repe, constant: child.constant}, nil
	case *nodes.FuncCallExpr:
		return s.inferFunc(e)
	case *nodes.BinaryExpr:
		return s.inferBinary(e)
	case *nodes.LikeExpr:
		left, err := s.infer(e.Expr)
		if err != nil {
			return exprCollation{}, err
		}
		right, err := s.infer(e.Pattern)
		if err != nil {
			return exprCollation{}, err
		}
		if _, err := s.infer(e.Escape); err != nil {
			return exprCollation{}, err
		}
		name := "like"
		if e.ILike {
			name = "ilike"
		}
		if _, err := deriveCollation(name, false, left, right); err != nil {
			return exprCollation{}, err
		}
		return numericCollation, nil
	case *nodes.UnaryExpr:
		if _, err := s.infer(e.Operand); err != nil {
			return exprCollation{}, err
		}
		if e.Op == nodes.UnaryBinary {
			return exprCollation{}, nil
		}
		return numericCollation, nil
	case *nodes.SubqueryExpr:
		return exprCollation{}, nil
	}
	return exprCollation{}, s.inferChildren(expr)
}

// inferChildren derives the collation of each expression directly inside an
// expression the catalog does not model, for the errors they answer.
func (s *viewScope) inferChildren(expr nodes.ExprNode) error {
	var err error
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if err != nil || n == nil {
			return false
		}
		if n == nodes.Node(expr) {
			return true
		}
		switch n.(type) {
		case *nodes.SubqueryExpr, *nodes.SelectStmt:
			return false
		}
		if child, ok := n.(nodes.ExprNode); ok {
			_, err = s.infer(child)
			return false
		}
		return true
	})
	return err
}

func stringLiteralCollation(lit *nodes.StringLit) exprCollation {
	cs, coll := connectionCharset, connectionCollation
	if lit.Charset != "" {
		cs = normalizeCharsetName(strings.TrimPrefix(lit.Charset, "_"))
		var ok bool
		if coll, ok = defaultCollationForCharset[cs]; !ok {
			return exprCollation{}
		}
	}
	repe := repertoireASCII
	for i := 0; i < len(lit.Value); i++ {
		if lit.Value[i] >= 0x80 {
			repe = repertoireUnicode
			break
		}
	}
	v := lit.Value
	return exprCollation{known: true, charset: cs, collation: coll, coer: coercibilityCoercible, repe: repe, constant: true, value: &v}
}

// inferFunc is deriveCollation for the functions the catalog models.
func (s *viewScope) inferFunc(fn *nodes.FuncCallExpr) (exprCollation, error) {
	args := make([]exprCollation, len(fn.Args))
	for i, arg := range fn.Args {
		ec, err := s.infer(arg)
		if err != nil {
			return exprCollation{}, err
		}
		args[i] = ec
	}
	if _, err := s.infer(fn.Separator); err != nil {
		return exprCollation{}, err
	}
	if fn.Schema != "" || fn.Over != nil || fn.Star {
		return exprCollation{}, nil
	}
	name := strings.ToLower(fn.Name)
	var (
		part         []exprCollation
		returnString = true
	)
	switch name {
	case "concat", "concat_ws", "lower", "lcase", "reverse", "upper", "ucase", "quote":
		part = args
	case "coalesce", "greatest", "least":
		part = args
		returnString = false
	case "left", "right", "repeat", "ltrim", "rtrim", "substr", "substring", "substring_index", "replace", "mid", "translate":
		part = argsAt(args, 0)
	case "trim", "trim_both", "trim_leading", "trim_trailing":
		// TRIM([remstr FROM] str): the parser keeps str last.
		part = argsAt(args, len(args)-1)
	case "insert":
		part = argsAt(args, 0, 3)
	case "lpad", "rpad":
		part = argsAt(args, 0, 2)
	case "elt", "make_set":
		if len(args) > 0 {
			part = args[1:]
		}
	case "if":
		part = argsAt(args, 1, 2)
		returnString = false
	case "ifnull":
		part = argsAt(args, 0, 1)
		returnString = false
	case "strcmp", "find_in_set":
		if _, err := deriveCollation(name, false, args...); err != nil {
			return exprCollation{}, err
		}
		return numericCollation, nil
	case "locate", "instr", "position", "regexp_like", "regexp_instr":
		if _, err := deriveCollation(name, false, argsAt(args, 0, 1)...); err != nil {
			return exprCollation{}, err
		}
		return numericCollation, nil
	default:
		return exprCollation{}, nil
	}
	if part == nil {
		return exprCollation{}, nil
	}
	if !returnString {
		// COALESCE, IF and the rest return a string only when their
		// arguments do; the catalog models only that case.
		for _, a := range part {
			if a.numeric() {
				return exprCollation{}, nil
			}
		}
	}
	return deriveCollation(name, true, part...)
}

// argsAt returns args at the given positions, or nil when one is missing.
func argsAt(args []exprCollation, idx ...int) []exprCollation {
	out := make([]exprCollation, 0, len(idx))
	for _, i := range idx {
		if i < 0 || i >= len(args) {
			return nil
		}
		out = append(out, args[i])
	}
	return out
}

// comparisonOpNames are the operators deriveCollation checks as comparisons,
// spelled as TiDB's GetDisplayName writes them in error messages.
var comparisonOpNames = map[nodes.BinaryOp]string{
	nodes.BinOpEq:         "=",
	nodes.BinOpNe:         "!=",
	nodes.BinOpLt:         "<",
	nodes.BinOpGt:         ">",
	nodes.BinOpLe:         "<=",
	nodes.BinOpGe:         ">=",
	nodes.BinOpNullSafeEq: "<=>",
}

func (s *viewScope) inferBinary(e *nodes.BinaryExpr) (exprCollation, error) {
	left, err := s.infer(e.Left)
	if err != nil {
		return exprCollation{}, err
	}
	right, err := s.infer(e.Right)
	if err != nil {
		return exprCollation{}, err
	}
	if name, ok := comparisonOpNames[e.Op]; ok {
		// Only a comparison of two strings has a collation; any other
		// pair compares as numbers or times.
		if !left.numeric() && !right.numeric() {
			if _, err := deriveCollation(name, false, left, right); err != nil {
				return exprCollation{}, err
			}
		}
		return numericCollation, nil
	}
	switch e.Op {
	case nodes.BinOpRegexp:
		if _, err := deriveCollation("regexp", false, left, right); err != nil {
			return exprCollation{}, err
		}
		return numericCollation, nil
	case nodes.BinOpAdd, nodes.BinOpSub, nodes.BinOpMul, nodes.BinOpDiv, nodes.BinOpMod,
		nodes.BinOpAnd, nodes.BinOpOr, nodes.BinOpXor, nodes.BinOpBitAnd, nodes.BinOpBitOr,
		nodes.BinOpBitXor, nodes.BinOpShiftLeft, nodes.BinOpShiftRight, nodes.BinOpDivInt:
		return numericCollation, nil
	}
	return exprCollation{}, nil
}

// deriveCollation is CheckAndDeriveCollationFromExprs: it aggregates args for
// the function name, which returns a string when returnString is set. An
// argument with no answer makes the result have none.
func deriveCollation(name string, returnString bool, args ...exprCollation) (exprCollation, error) {
	for _, a := range args {
		if !a.known {
			return exprCollation{}, nil
		}
	}
	ec, ok := inferCollation(args)
	if !ok {
		return exprCollation{}, errIllegalMixCollations(name, args)
	}
	if !returnString && ec.coer == coercibilityNone {
		return exprCollation{}, errIllegalMixCollations(name, args)
	}
	if returnString && ec.coer == coercibilityNumeric {
		ec.charset, ec.collation = connectionCharset, connectionCollation
		ec.coer, ec.repe = coercibilityCoercible, repertoireASCII
	}
	if !safeConvert(ec, args) {
		return exprCollation{}, errIllegalMixCollations(name, args)
	}
	ec.constant = len(args) > 0
	for _, a := range args {
		ec.constant = ec.constant && a.constant
	}
	return ec, nil
}

// inferCollation aggregates args left to right, agg(a, b, c) = agg(agg(a, b), c).
// ok is false where TiDB's inferCollation returns nil.
func inferCollation(args []exprCollation) (exprCollation, bool) {
	if len(args) == 0 {
		return exprCollation{known: true, charset: "utf8mb4", collation: "utf8mb4_bin", coer: coercibilityIgnorable, repe: repertoireUnicode}, true
	}
	repe, coer := args[0].repe, args[0].coer
	cs, coll := args[0].charset, args[0].collation
	unknownCS := false
	for _, arg := range args[1:] {
		// A binary argument goes with any charset; at the same
		// coercibility the binary one wins.
		if coll == "binary" || arg.collation == "binary" {
			if coer > arg.coer || (coer == arg.coer && arg.collation == "binary") {
				coer, cs, coll = arg.coer, arg.charset, arg.collation
			}
			repe |= arg.repe
			continue
		}
		if cs != arg.charset {
			switch {
			case coer < arg.coer:
				if arg.repe == repertoireASCII || arg.coer >= coercibilitySysconst || isUnicodeCharset(cs) {
					repe |= arg.repe
					continue
				}
			case coer == arg.coer:
				if (isUnicodeCharset(cs) && !isUnicodeCharset(arg.charset)) || (cs == "utf8mb4" && arg.charset == "utf8") {
					repe |= arg.repe
					continue
				} else if (isUnicodeCharset(arg.charset) && !isUnicodeCharset(cs)) || (arg.charset == "utf8mb4" && cs == "utf8") {
					coer, cs, coll = arg.coer, arg.charset, arg.collation
					repe |= arg.repe
					continue
				} else if repe == repertoireASCII && arg.repe != repertoireASCII {
					coer, cs, coll = arg.coer, arg.charset, arg.collation
					repe |= arg.repe
					continue
				} else if repe != repertoireASCII && arg.repe == repertoireASCII {
					repe |= arg.repe
					continue
				}
			case coer > arg.coer:
				if repe == repertoireASCII || coer >= coercibilitySysconst || isUnicodeCharset(arg.charset) {
					coer, cs, coll = arg.coer, arg.charset, arg.collation
					repe |= arg.repe
					continue
				}
			}
			// No conversion without loss.
			repe |= arg.repe
			coer, cs, coll = coercibilityNone, "binary", "binary"
			unknownCS = true
			continue
		}
		// One charset: the lower coercibility wins; at the same one, a
		// _bin collation wins, and two others give the charset's _bin
		// collation with coercibility NONE.
		switch {
		case coer == arg.coer:
			switch {
			case coll == arg.collation:
			case coer == coercibilityExplicit:
				return exprCollation{}, false
			case isBinCollation(coll):
			case isBinCollation(arg.collation):
				coer, cs, coll = arg.coer, arg.charset, arg.collation
			default:
				coer, cs, coll = coercibilityNone, arg.charset, binCollationFor(arg.charset)
			}
		case coer > arg.coer:
			coer, cs, coll = arg.coer, arg.charset, arg.collation
		}
		repe |= arg.repe
	}
	if unknownCS && coer != coercibilityExplicit {
		return exprCollation{}, false
	}
	return exprCollation{known: true, charset: cs, collation: coll, coer: coer, repe: repe}, true
}

// safeConvert reports whether every argument converts to ec's charset.
func safeConvert(ec exprCollation, args []exprCollation) bool {
	for _, arg := range args {
		if arg.charset == ec.charset || arg.repe == repertoireASCII || arg.charset == "binary" {
			continue
		}
		if arg.constant {
			// TiDB converts a constant by value. Of the charsets it has,
			// only ascii and gbk refuse some UTF-8 text; latin1 is UTF-8
			// inside TiDB. A constant the catalog has no value for, or a
			// gbk one, is taken to convert.
			if arg.isNull || arg.value == nil || ec.charset != "ascii" {
				continue
			}
			for i := 0; i < len(*arg.value); i++ {
				if (*arg.value)[i] >= 0x80 {
					return false
				}
			}
			continue
		}
		if arg.collation != "binary" && ec.charset != "binary" && !isUnicodeCharset(ec.charset) {
			return false
		}
	}
	return true
}

func isUnicodeCharset(cs string) bool { return cs == "utf8" || cs == "utf8mb4" }

func isBinCollation(coll string) bool {
	switch coll {
	case "ascii_bin", "latin1_bin", "utf8_bin", "utf8mb4_bin", "gbk_bin":
		return true
	}
	return false
}

func binCollationFor(cs string) string {
	switch cs {
	case "utf8":
		return "utf8_bin"
	case "gbk":
		return "gbk_bin"
	}
	return "utf8mb4_bin"
}

// errIllegalMixCollations is illegalMixCollationErr: 1267 for two arguments,
// 1270 for three, 1271 otherwise.
func errIllegalMixCollations(name string, args []exprCollation) error {
	part := func(a exprCollation) string {
		return fmt.Sprintf("(%s,%s)", a.collation, coercibilityNames[a.coer])
	}
	switch len(args) {
	case 2:
		return &Error{Code: ErrCantAggregate2Collations, SQLState: sqlState(ErrCantAggregate2Collations),
			Message: fmt.Sprintf("Illegal mix of collations %s and %s for operation '%s'", part(args[0]), part(args[1]), name)}
	case 3:
		return &Error{Code: ErrCantAggregate3Collations, SQLState: sqlState(ErrCantAggregate3Collations),
			Message: fmt.Sprintf("Illegal mix of collations %s, %s, %s for operation '%s'", part(args[0]), part(args[1]), part(args[2]), name)}
	}
	return &Error{Code: ErrCantAggregateNCollations, SQLState: sqlState(ErrCantAggregateNCollations),
		Message: fmt.Sprintf("Illegal mix of collations for operation '%s'", name)}
}
