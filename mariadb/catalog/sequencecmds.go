package catalog

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// Sequence is a MariaDB sequence. It shares the table namespace of its
// database: a table, view, or sequence of the same name answers 1050.
type Sequence struct {
	Name     string
	Database *Database
	// DataType is the AS type as SHOW CREATE SEQUENCE writes it ("bigint",
	// "int unsigned"). Sequences without AS are bigint.
	DataType  string
	Start     string
	MinValue  string
	MaxValue  string
	Increment string
	// Cache is the cache size; 0 is NOCACHE.
	Cache   string
	Cycle   bool
	Engine  string
	Comment string
	// Temporary marks a CREATE TEMPORARY SEQUENCE, which lasts for the
	// session and is no part of the schema.
	Temporary bool
}

// sequenceTypes are the integer types a sequence may be declared AS, with
// their signed ranges.
var sequenceTypes = map[string][2]string{
	"tinyint":   {"-128", "127"},
	"smallint":  {"-32768", "32767"},
	"mediumint": {"-8388608", "8388607"},
	"int":       {"-2147483648", "2147483647"},
	"bigint":    {"-9223372036854775808", "9223372036854775807"},
}

const defaultSequenceCache = "1000"

func (c *Catalog) createSequence(stmt *nodes.CreateSequenceStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	// As with tables, the catalog keeps one object per name, so a temporary
	// sequence that would shadow a permanent object is reported a duplicate.
	if db.Tables[key] != nil || db.Views[key] != nil || db.Sequences[key] != nil {
		switch {
		case stmt.IfNotExists:
			return nil
		case stmt.OrReplace && db.Views[key] == nil && (!stmt.Temporary || temporaryOnly(db, key)):
		default:
			return errDupTable(stmt.Name.Name)
		}
	}
	seq := &Sequence{Name: stmt.Name.Name, Database: db, DataType: "bigint", Cache: defaultSequenceCache, Engine: "InnoDB", Temporary: stmt.Temporary}
	if stmt.DataType != nil {
		seq.DataType = sequenceDataType(stmt.DataType)
		if seq.DataType == "" {
			return errSequenceOutOfRange(db.Name, stmt.Name.Name)
		}
	}
	opts := sequenceOptions{
		start: stmt.Start, increment: stmt.Increment, min: stmt.MinValue, max: stmt.MaxValue,
		noMin: stmt.NoMinValue, noMax: stmt.NoMaxValue, cache: stmt.Cache, noCache: stmt.NoCache, cycle: stmt.Cycle,
	}
	if err := seq.apply(opts, true); err != nil {
		return err
	}
	for _, opt := range stmt.Options {
		switch toLower(opt.Name) {
		case "engine":
			seq.Engine = canonicalEngineName(opt.Value)
		case "comment":
			seq.Comment = opt.Value
		}
	}
	c.dropReplacedTable(db, db.Tables[key])
	db.Sequences[key] = seq
	return nil
}

// temporaryOnly reports whether every object named key in db is temporary.
func temporaryOnly(db *Database, key string) bool {
	if t := db.Tables[key]; t != nil && !t.Temporary {
		return false
	}
	if s := db.Sequences[key]; s != nil && !s.Temporary {
		return false
	}
	return true
}

func (c *Catalog) alterSequence(stmt *nodes.AlterSequenceStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	seq := db.Sequences[toLower(stmt.Name.Name)]
	if seq == nil {
		if stmt.IfExists {
			return nil
		}
		if db.Tables[toLower(stmt.Name.Name)] != nil {
			return errNotSequence(4089, db.Name, stmt.Name.Name)
		}
		return errUnknownSequence(stmt.Name.Name)
	}
	next := *seq
	if stmt.DataType != nil {
		next.DataType = sequenceDataType(stmt.DataType)
		if next.DataType == "" {
			return errSequenceOutOfRange(db.Name, seq.Name)
		}
	}
	opts := sequenceOptions{
		start: stmt.Start, increment: stmt.Increment, min: stmt.MinValue, max: stmt.MaxValue,
		noMin: stmt.NoMinValue, noMax: stmt.NoMaxValue, cache: stmt.Cache, noCache: stmt.NoCache, cycle: stmt.Cycle,
	}
	if err := next.apply(opts, false); err != nil {
		return err
	}
	// RESTART WITH sets the next value, which has to lie within the limits,
	// though SHOW CREATE SEQUENCE does not show it.
	if stmt.RestartWith != nil {
		restart, ok := sequenceNumber(stmt.RestartWith)
		v, _ := new(big.Int).SetString(restart, 10)
		minV, _ := new(big.Int).SetString(next.MinValue, 10)
		maxV, _ := new(big.Int).SetString(next.MaxValue, 10)
		if !ok || v == nil || v.Cmp(minV) < 0 || v.Cmp(maxV) > 0 {
			return errSequenceOutOfRange(db.Name, seq.Name)
		}
	}
	*seq = next
	return nil
}

func (c *Catalog) dropSequence(stmt *nodes.DropSequenceStmt) error {
	for _, ref := range stmt.Sequences {
		db, err := c.resolveDatabase(ref.Schema)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		key := toLower(ref.Name)
		if seq := db.Sequences[key]; seq == nil || (stmt.Temporary && !seq.Temporary) {
			if stmt.IfExists {
				continue
			}
			if db.Tables[key] != nil {
				return errNotSequence(4090, db.Name, ref.Name)
			}
			return errUnknownSequence(db.Name + "." + ref.Name)
		}
		delete(db.Sequences, key)
	}
	return nil
}

// sequenceOptions is what a CREATE or ALTER SEQUENCE sets.
type sequenceOptions struct {
	start, increment, min, max, cache nodes.ExprNode
	noMin, noMax, noCache             bool
	cycle                             *bool
}

// apply sets the options on seq and checks the result, as
// sequence_definition::check_and_adjust does (sql/sql_sequence.cc): an option
// CREATE leaves unset takes the value the type and the increment's sign give
// it, while ALTER keeps the sequence's own.
func (seq *Sequence) apply(o sequenceOptions, create bool) error {
	outOfRange := errSequenceOutOfRange(seq.Database.Name, seq.Name)
	value := func(e nodes.ExprNode, into *string) bool {
		if e == nil {
			return true
		}
		n, ok := sequenceNumber(e)
		if !ok {
			return false
		}
		*into = n
		return true
	}
	if create {
		seq.Increment = "1"
	}
	if !value(o.increment, &seq.Increment) {
		return outOfRange
	}
	inc, _ := new(big.Int).SetString(seq.Increment, 10)
	if inc.Sign() == 0 {
		// INCREMENT BY 0 takes auto_increment_increment, 1 by default.
		seq.Increment, inc = "1", big.NewInt(1)
	}
	typeMin, typeMax := sequenceTypeRange(seq.DataType)
	setMin, setMax := o.min != nil && !o.noMin, o.max != nil && !o.noMax
	if create || o.noMin {
		if inc.Sign() < 0 {
			seq.MinValue = new(big.Int).Add(typeMin, big.NewInt(1)).String()
		} else {
			seq.MinValue = "1"
		}
	}
	if create || o.noMax {
		if inc.Sign() > 0 {
			seq.MaxValue = new(big.Int).Sub(typeMax, big.NewInt(1)).String()
		} else {
			seq.MaxValue = "-1"
		}
	}
	if setMin && !value(o.min, &seq.MinValue) {
		return outOfRange
	}
	if setMax && !value(o.max, &seq.MaxValue) {
		return outOfRange
	}
	// A limit beyond the type's is held to the type's own, less one. Measured
	// on mariadb:11.8.8: AS TINYINT MAXVALUE 200 stores maxvalue 126, and
	// MINVALUE -200 minvalue -127.
	lowest, highest := new(big.Int).Add(typeMin, big.NewInt(1)), new(big.Int).Sub(typeMax, big.NewInt(1))
	if v, ok := new(big.Int).SetString(seq.MinValue, 10); ok && v.Cmp(lowest) < 0 {
		seq.MinValue = lowest.String()
	}
	if v, ok := new(big.Int).SetString(seq.MaxValue, 10); ok && v.Cmp(highest) > 0 {
		seq.MaxValue = highest.String()
	}
	if create {
		if inc.Sign() < 0 {
			seq.Start = seq.MaxValue
		} else {
			seq.Start = seq.MinValue
		}
	}
	if !value(o.start, &seq.Start) {
		return outOfRange
	}
	switch {
	case o.noCache:
		seq.Cache = "0"
	case !value(o.cache, &seq.Cache):
		return outOfRange
	}
	if o.cycle != nil {
		seq.Cycle = *o.cycle
	}
	start, _ := new(big.Int).SetString(seq.Start, 10)
	minV, _ := new(big.Int).SetString(seq.MinValue, 10)
	maxV, _ := new(big.Int).SetString(seq.MaxValue, 10)
	cache, _ := new(big.Int).SetString(seq.Cache, 10)
	if start == nil || minV == nil || maxV == nil || cache == nil ||
		maxV.Cmp(start) < 0 || maxV.Cmp(minV) <= 0 || start.Cmp(minV) < 0 || cache.Sign() < 0 {
		return outOfRange
	}
	return nil
}

// sequenceDataType renders an AS type as SHOW CREATE SEQUENCE writes it, or
// "" for a type a sequence cannot take.
func sequenceDataType(dt *nodes.DataType) string {
	name := toLower(dt.Name)
	switch name {
	case "integer":
		name = "int"
	case "middleint":
		name = "mediumint"
	}
	if _, ok := sequenceTypes[name]; !ok {
		return ""
	}
	if dt.Unsigned {
		return name + " unsigned"
	}
	return name
}

// sequenceTypeRange is the value range of a sequence's type.
func sequenceTypeRange(dataType string) (*big.Int, *big.Int) {
	name, unsigned := strings.CutSuffix(dataType, " unsigned")
	r, ok := sequenceTypes[name]
	if !ok {
		r = sequenceTypes["bigint"]
	}
	lo, _ := new(big.Int).SetString(r[0], 10)
	hi, _ := new(big.Int).SetString(r[1], 10)
	if unsigned {
		hi.Sub(hi, lo)
		lo.SetInt64(0)
	}
	return lo, hi
}

// sequenceNumber reads an integer option value, a literal optionally negated,
// in the canonical decimal form the server stores: START WITH 001 is 1.
func sequenceNumber(e nodes.ExprNode) (string, bool) {
	var text string
	switch n := e.(type) {
	case *nodes.IntLit:
		text = n.Text
		if text == "" {
			text = strconv.FormatInt(n.Value, 10)
		}
	case *nodes.FloatLit:
		// A DECIMAL_NUM, a fraction or an integer past uint64, reads as
		// ULONGLONG_MAX, and as LONGLONG_MIN after a minus (see
		// parser.sequenceValue), which the limits then clamp. Measured on
		// mariadb:11.8.8: MAXVALUE 1.5 stores 9223372036854775806.
		text = strconv.FormatUint(math.MaxUint64, 10)
	case *nodes.UnaryExpr:
		if n.Op != nodes.UnaryMinus {
			return "", false
		}
		if _, ok := n.Operand.(*nodes.FloatLit); ok {
			text = strconv.FormatInt(math.MinInt64, 10)
			break
		}
		v, ok := sequenceNumber(n.Operand)
		if !ok {
			return "", false
		}
		text = "-" + v
	default:
		text = strings.TrimSpace(nodeToSQL(e))
	}
	v, ok := new(big.Int).SetString(text, 10)
	if !ok {
		return "", false
	}
	return v.String(), true
}

// ShowCreateSequence renders a sequence as MariaDB 11.8.8's SHOW CREATE
// SEQUENCE does. It returns "" when there is no such sequence.
func (c *Catalog) ShowCreateSequence(database, name string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	seq := db.Sequences[toLower(name)]
	if seq == nil {
		return ""
	}
	return showCreateSequence(seq)
}

func showCreateSequence(seq *Sequence) string {
	var b strings.Builder
	b.WriteString("CREATE SEQUENCE ")
	b.WriteString(quoteIdent(seq.Name))
	b.WriteString(" ")
	if seq.DataType != "" && seq.DataType != "bigint" {
		b.WriteString("as " + seq.DataType + " ")
	}
	fmt.Fprintf(&b, "start with %s minvalue %s maxvalue %s increment by %s ", seq.Start, seq.MinValue, seq.MaxValue, seq.Increment)
	if seq.Cache == "0" {
		b.WriteString("nocache ")
	} else {
		b.WriteString("cache " + seq.Cache + " ")
	}
	if seq.Cycle {
		b.WriteString("cycle ")
	} else {
		b.WriteString("nocycle ")
	}
	b.WriteString("ENGINE=" + seq.Engine)
	if seq.Comment != "" {
		b.WriteString(" COMMENT=" + quoteStringLiteral(seq.Comment))
	}
	return b.String()
}

// canonicalEngineName spells a storage engine as SHOW CREATE writes it.
func canonicalEngineName(name string) string {
	switch toLower(name) {
	case "innodb":
		return "InnoDB"
	case "myisam":
		return "MyISAM"
	case "aria":
		return "Aria"
	case "memory", "heap":
		return "MEMORY"
	}
	return name
}

func quoteIdent(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// quoteStringLiteral quotes s as SHOW CREATE writes a comment.
func quoteStringLiteral(s string) string {
	return "'" + escapeComment(s) + "'"
}

// qualifySequenceCalls returns expr with each sequence it reads named by its
// database, the table's where the call names none, in the form MariaDB stores a
// default in: NEXT VALUE FOR s and PREVIOUS VALUE FOR s become nextval and
// lastval, and setval takes its optional arguments. Measured on
// mariadb:11.8.8: DEFAULT (NEXT VALUE FOR s1) is stored nextval(`db`.`s1`),
// and setval(s1, 5) setval(`db`.`s1`,5,1,0).
func qualifySequenceCalls(expr nodes.ExprNode, tbl *Table) nodes.ExprNode {
	db := ""
	if tbl != nil && tbl.Database != nil {
		db = tbl.Database.Name
	}
	return qualifySequenceCallsInDB(expr, db)
}

func qualifySequenceCallsInDB(expr nodes.ExprNode, db string) nodes.ExprNode {
	ref := func(schema, name string) nodes.ExprNode {
		if schema == "" {
			schema = db
		}
		return &nodes.ColumnRef{Table: schema, Column: name}
	}
	var rewrite func(nodes.ExprNode) nodes.ExprNode
	rewrite = func(expr nodes.ExprNode) nodes.ExprNode {
		switch e := expr.(type) {
		case *nodes.ParenExpr:
			return &nodes.ParenExpr{Loc: e.Loc, Expr: rewrite(e.Expr)}
		case *nodes.BinaryExpr:
			out := *e
			out.Left, out.Right = rewrite(e.Left), rewrite(e.Right)
			return &out
		case *nodes.UnaryExpr:
			out := *e
			out.Operand = rewrite(e.Operand)
			return &out
		case *nodes.NextValueForExpr:
			return &nodes.FuncCallExpr{Loc: e.Loc, Name: "NEXTVAL", RawName: "nextval", Args: []nodes.ExprNode{ref(e.Sequence.Schema, e.Sequence.Name)}}
		case *nodes.PreviousValueForExpr:
			return &nodes.FuncCallExpr{Loc: e.Loc, Name: "LASTVAL", RawName: "lastval", Args: []nodes.ExprNode{ref(e.Sequence.Schema, e.Sequence.Name)}}
		case *nodes.FuncCallExpr:
			out := *e
			out.Args = make([]nodes.ExprNode, len(e.Args))
			for i, a := range e.Args {
				out.Args[i] = rewrite(a)
			}
			switch name := strings.ToUpper(e.Name); {
			case e.Schema != "" || len(e.Args) == 0:
			case name == "NEXTVAL" || name == "LASTVAL" || name == "SETVAL":
				if seq, ok := e.Args[0].(*nodes.ColumnRef); ok && seq.Schema == "" {
					out.Args[0] = ref(seq.Table, seq.Column)
					for name == "SETVAL" && len(out.Args) < 4 {
						out.Args = append(out.Args, &nodes.IntLit{Value: []int64{1, 0}[len(out.Args)-2]})
					}
				}
			}
			return &out
		}
		return expr
	}
	return rewrite(expr)
}

// dropReplacedSequence drops the sequence CREATE OR REPLACE TABLE replaces: the
// two share a name space, and the table takes the name. Measured on
// mariadb:11.8.8.
func dropReplacedSequence(db *Database, key string, stmt *nodes.CreateTableStmt) {
	if stmt.OrReplace && !stmt.Temporary {
		delete(db.Sequences, key)
	}
}
