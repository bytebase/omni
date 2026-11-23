package review

import (
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// sequenceColumnNames are the eight fixed columns of a MariaDB sequence's
// underlying table.
var sequenceColumnNames = []string{
	"next_not_cached_value",
	"minimum_value",
	"maximum_value",
	"start_value",
	"increment",
	"cache_size",
	"cycle_option",
	"cycle_count",
}

func isSequenceColumn(name string) bool {
	for _, c := range sequenceColumnNames {
		if strings.EqualFold(name, c) {
			return true
		}
	}
	return false
}

func notPartitioned(table string) string {
	return "Table " + review.Code(table) + " is not partitioned."
}

func notSystemVersioned(table string) string {
	return "Table " + review.Code(table) + " is not system-versioned."
}

func sequenceInvalidStructure(name string) string {
	return "Sequence " + review.Code(name) + " doesn't support this " + review.Code("ALTER TABLE") + " operation."
}

func isPartitionMgmtCmd(t ast.AlterTableCmdType) bool {
	switch t {
	case ast.ATAddPartition, ast.ATDropPartition, ast.ATTruncatePartition,
		ast.ATCoalescePartition, ast.ATReorganizePartition, ast.ATExchangePartition,
		ast.ATConvertPartitionToTable, ast.ATConvertTableToPartition,
		ast.ATRemovePartitioning, ast.ATRebuildPartition, ast.ATOptimizePartition,
		ast.ATAnalyzePartition, ast.ATCheckPartition, ast.ATRepairPartition,
		ast.ATDiscardPartitionTablespace, ast.ATImportPartitionTablespace:
		return true
	}
	return false
}

func alterHasPartitionBy(cmds []*ast.AlterTableCmd) bool {
	for _, c := range cmds {
		if c != nil && c.Type == ast.ATPartitionBy {
			return true
		}
	}
	return false
}

func alterAddsSystemVersioning(cmds []*ast.AlterTableCmd) bool {
	for _, c := range cmds {
		if c != nil && c.Type == ast.ATAddSystemVersioning {
			return true
		}
	}
	return false
}

func addsColumnVersioning(cmd *ast.AlterTableCmd) bool {
	switch cmd.Type {
	case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
	default:
		return false
	}
	if cmd.Column != nil && cmd.Column.SystemVersioning != ast.ColVersioningNone {
		return true
	}
	for _, c := range cmd.Columns {
		if c != nil && c.SystemVersioning != ast.ColVersioningNone {
			return true
		}
	}
	return false
}

func addsRowBoundColumn(cmd *ast.AlterTableCmd) bool {
	switch cmd.Type {
	case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
	default:
		return false
	}
	if cmd.Column != nil && cmd.Column.Generated != nil && cmd.Column.Generated.RowBound != ast.RowBoundNone {
		return true
	}
	for _, c := range cmd.Columns {
		if c != nil && c.Generated != nil && c.Generated.RowBound != ast.RowBoundNone {
			return true
		}
	}
	return false
}

// alterSequence follows ALTER SEQUENCE [IF EXISTS] name ...: the target must
// exist as a sequence (4091 missing, 4089 table or view) unless IF EXISTS is
// given, and an out-of-range signed-64-bit INCREMENT always fails (4085).
func (b *builder) alterSequence(v *ast.AlterSequenceStmt) bool {
	r, st := b.ref(v.Name)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	if b.tooLong(v.Name.Name, r.rng) {
		return true
	}
	tk := tempKey(r.db, r.name)
	tf, tkind := b.l.Lookup(tk)
	if tf != walk.Absent {
		if tf == walk.Present && tkind == relSequence {
			b.checkSequenceIncrement(v.Increment, r)
			return true
		}
		return false
	}
	k := relKey(r.db, r.name)
	f, kind := b.l.Lookup(k)
	if !v.IfExists {
		b.needKind(k, []string{relSequence}, r.rng, notFound("sequence", r.written, ""), r.written, "sequence")
	}
	if f == walk.Present && kind == relSequence {
		b.checkSequenceIncrement(v.Increment, r)
		b.checkSequenceValues(v, r, k)
	}
	return true
}

// sequenceValues are a sequence's start, minimum, and maximum values.
type sequenceValues struct {
	start, min, max *big.Int
}

// syncedSequenceValues reads the values the sync lists of a sequence,
// false when it lists none.
func syncedSequenceValues(q *metadata.SequenceMetadata) (sequenceValues, bool) {
	start, ok1 := new(big.Int).SetString(q.GetStart(), 10)
	lo, ok2 := new(big.Int).SetString(q.GetMinValue(), 10)
	hi, ok3 := new(big.Int).SetString(q.GetMaxValue(), 10)
	if !ok1 || !ok2 || !ok3 {
		return sequenceValues{}, false
	}
	return sequenceValues{start: start, min: lo, max: hi}, true
}

// checkSequenceValues follows ALTER SEQUENCE's START, MINVALUE, and
// MAXVALUE on a sequence whose values the walk knows: the values the
// statement leaves must keep min <= start <= max and min < max, or it
// fails (4085, sequence_definition::check_and_adjust; verified on
// mariadb:11.8.8). A statement that changes the value type or asks for a
// default bound leaves the values to the server.
func (b *builder) checkSequenceValues(v *ast.AlterSequenceStmt, r tableRef, k walk.Key) {
	values, known := b.w.sequences[k]
	forget := func() { b.Session(func() { delete(b.w.sequences, k) }) }
	if !known || v.DataType != nil || v.NoMinValue || v.NoMaxValue {
		forget()
		return
	}
	next := values
	for _, set := range []struct {
		e  ast.ExprNode
		to **big.Int
	}{{v.Start, &next.start}, {v.MinValue, &next.min}, {v.MaxValue, &next.max}} {
		if set.e == nil {
			continue
		}
		n, ok := sequenceBigInt(set.e)
		if !ok {
			forget()
			return
		}
		*set.to = n
	}
	if next.max.Cmp(next.start) < 0 || next.max.Cmp(next.min) <= 0 || next.start.Cmp(next.min) < 0 {
		b.fail(r.rng, nil, review.UnsupportedForm, "Sequence "+review.Code(r.written)+" has invalid values.") // 4085
		return
	}
	b.Session(func() { b.w.sequences[k] = next })
}

// checkSequenceIncrement rejects an INCREMENT outside MariaDB's signed 64-bit
// range (<= LONGLONG_MIN or >= LONGLONG_MAX), which MariaDB refuses with 4085
// regardless of the sequence's value data type.
func (b *builder) checkSequenceIncrement(inc ast.ExprNode, r tableRef) {
	if inc == nil {
		return
	}
	val, ok := sequenceBigInt(inc)
	if !ok {
		return
	}
	minAllowed := big.NewInt(math.MinInt64 + 1)
	maxAllowed := big.NewInt(math.MaxInt64 - 1)
	if val.Cmp(minAllowed) < 0 || val.Cmp(maxAllowed) > 0 {
		b.fail(r.rng, nil, review.UnsupportedForm, "Sequence "+review.Code(r.written)+" has an out-of-range increment.")
	}
}

func sequenceBigInt(e ast.ExprNode) (*big.Int, bool) {
	switch n := e.(type) {
	case *ast.IntLit:
		text := n.Text
		if text == "" {
			text = strconv.FormatInt(n.Value, 10)
		}
		return new(big.Int).SetString(text, 10)
	case *ast.UnaryExpr:
		if n.Op == ast.UnaryMinus {
			v, ok := sequenceBigInt(n.Operand)
			if !ok {
				return nil, false
			}
			return v.Neg(v), true
		}
		if n.Op == ast.UnaryPlus {
			return sequenceBigInt(n.Operand)
		}
	}
	return nil, false
}

// alterSequenceTable follows ALTER TABLE on a known base sequence. On MariaDB
// 11.8.8 a sequence is backed by an 8-column keyless, unpartitioned,
// non-versioned table: table options, RENAME, FORCE, ENABLE/DISABLE KEYS,
// CONVERT TO CHARSET, and DISCARD TABLESPACE succeed, while structural,
// indexing, partitioning, and system-versioning subcommands fail.
func (b *builder) alterSequenceTable(v *ast.AlterTableStmt, r tableRef) bool {
	rel := relKey(r.db, r.name)
	keys := []walk.Key{rel}
	hasAddSysVer := alterAddsSystemVersioning(v.Commands)
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		if c.Type == ast.ATPartitionBy {
			if c.PartitionBy != nil && c.PartitionBy.Type == ast.PartitionSystemTime {
				b.fail(rangeOf(c.Loc), keys, review.UnsupportedForm, notSystemVersioned(r.written))
				return true
			}
			b.fail(rangeOf(c.Loc), keys, review.WrongObjectKind, sequenceInvalidStructure(r.written))
			return true
		}
		if addsColumnVersioning(c) || (!hasAddSysVer && addsRowBoundColumn(c)) {
			b.fail(rangeOf(c.Loc), keys, review.UnsupportedForm, notSystemVersioned(r.written))
			return true
		}
	}
	var rename *ast.AlterTableCmd
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ATAddColumn:
			if c.After != "" && !isSequenceColumn(c.After) {
				b.fail(rng, keys, review.ObjectNotFound, notFound("column", ident(c.After), r.written))
				return true
			}
			cols := c.Columns
			if c.Column != nil {
				cols = append([]*ast.ColumnDef{c.Column}, cols...)
			}
			if c.IfNotExists && len(cols) > 0 && len(c.Constraints) == 0 {
				allSeq := true
				for _, col := range cols {
					if col == nil || !isSequenceColumn(col.Name) {
						allSeq = false
						break
					}
				}
				if allSeq {
					continue
				}
			}
			b.fail(rng, keys, review.WrongObjectKind, sequenceInvalidStructure(r.written))
			return true
		case ast.ATDropColumn:
			if !isSequenceColumn(c.Name) {
				if c.IfExists {
					continue
				}
				b.fail(rng, keys, review.ObjectNotFound, notFound("column", ident(c.Name), r.written))
				return true
			}
			b.fail(rng, keys, review.WrongObjectKind, sequenceInvalidStructure(r.written))
			return true
		case ast.ATModifyColumn:
			colName := c.Name
			if colName == "" && c.Column != nil {
				colName = c.Column.Name
			}
			if !isSequenceColumn(colName) {
				if c.IfExists {
					continue
				}
				b.fail(rng, keys, review.ObjectNotFound, notFound("column", ident(colName), r.written))
				return true
			}
			b.fail(rng, keys, review.WrongObjectKind, sequenceInvalidStructure(r.written))
			return true
		case ast.ATChangeColumn, ast.ATAlterColumnDefault, ast.ATRenameColumn:
			if !isSequenceColumn(c.Name) {
				if c.IfExists {
					continue
				}
				b.fail(rng, keys, review.ObjectNotFound, notFound("column", ident(c.Name), r.written))
				return true
			}
			b.fail(rng, keys, review.WrongObjectKind, sequenceInvalidStructure(r.written))
			return true
		case ast.ATAddConstraint:
			if c.Constraint != nil && c.Constraint.Type == ast.ConstrForeignKey {
				for _, col := range c.Constraint.Columns {
					if !isSequenceColumn(col) {
						b.fail(rng, keys, review.ObjectNotFound, notFound("column", ident(col), r.written))
						return true
					}
				}
			}
			b.fail(rng, keys, review.WrongObjectKind, sequenceInvalidStructure(r.written))
			return true
		case ast.ATAddIndex, ast.ATAddSystemVersioning, ast.ATOrderBy:
			b.fail(rng, keys, review.WrongObjectKind, sequenceInvalidStructure(r.written))
			return true
		case ast.ATDropIndex:
			if c.IfExists {
				continue
			}
			b.fail(rng, keys, review.ObjectNotFound, notFound("index", ident(c.Name), r.written))
			return true
		case ast.ATDropConstraint:
			if c.IfExists {
				continue
			}
			switch dropKindOf(b.w.sql, c) {
			case dropPrimaryKey:
				b.fail(rng, keys, review.ObjectNotFound, noPrimaryKey(r.written))
			case dropForeignKey:
				b.fail(rng, keys, review.ObjectNotFound, notFound("foreign key", ident(c.Name), r.written))
			default:
				b.fail(rng, keys, review.ObjectNotFound, notFound("constraint", ident(c.Name), r.written))
			}
			return true
		case ast.ATAlterIndexVisible, ast.ATAlterIndexInvisible, ast.ATRenameIndex:
			if c.IfExists {
				continue
			}
			b.fail(rng, keys, review.ObjectNotFound, notFound("index", ident(c.Name), r.written))
			return true
		case ast.ATDropSystemVersioning:
			b.fail(rng, keys, review.UnsupportedForm, notSystemVersioned(r.written))
			return true
		case ast.ATDropPeriod:
			if c.IfExists {
				continue
			}
			if strings.EqualFold(c.PeriodName, "SYSTEM_TIME") {
				b.fail(rng, keys, review.UnsupportedForm, notSystemVersioned(r.written))
				return true
			}
			b.fail(rng, keys, review.ObjectNotFound, notFound("period", ident(c.PeriodName), r.written))
			return true
		case ast.ATAddPeriod:
			if !strings.EqualFold(c.PeriodName, "SYSTEM_TIME") {
				b.fail(rng, keys, review.WrongObjectKind, sequenceInvalidStructure(r.written))
				return true
			}
		case ast.ATRenameTable:
			rename = c
		case ast.ATTableOption, ast.ATConvertCharset, ast.ATAlgorithm, ast.ATLock, ast.ATForce,
			ast.ATEnableKeys, ast.ATDisableKeys, ast.ATDiscardTablespace, ast.ATImportTablespace:
		default:
			if isPartitionMgmtCmd(c.Type) {
				b.fail(rng, keys, review.UnsupportedForm, notPartitioned(r.written))
				return true
			}
			return false
		}
	}
	if rename != nil {
		b.Next()
		return b.alterRename(r, rename, relSequence)
	}
	return true
}

// checkIndexColumns checks an index's plain column key parts against the
// ledger (1072) and rejects a column named twice in the same index (1060).
// Functional key parts and application-time WITHOUT OVERLAPS periods are skipped.
func (b *builder) checkIndexColumns(r tableRef, cols []*ast.IndexColumn) {
	seen := make(map[string]bool, len(cols))
	for _, ic := range cols {
		if ic == nil || ic.Functional || ic.WithoutOverlaps {
			continue
		}
		ref, ok := ic.Expr.(*ast.ColumnRef)
		if !ok || ref.Star || ref.Table != "" || ref.Schema != "" {
			continue
		}
		k, ok := b.key(spaceColumn, ref.Column)
		if !ok {
			continue
		}
		if seen[k] {
			b.fail(rangeOf(ic.Loc), nil, review.DuplicateName, statesTwice("names", "column", ident(ref.Column)))
			return
		}
		seen[k] = true
		b.need(colKey(r.db, r.name, k), walk.Present, nil, rangeOf(ic.Loc), notFound("column", ident(ref.Column), r.written), "")
	}
}

// checkCreateIndexColumns checks a CREATE TABLE index's plain column key parts
// against the statement's declared column map (1060 duplicate, 1072 missing).
func (b *builder) checkCreateIndexColumns(cols []*ast.IndexColumn, declared map[string]bool, tableName string) bool {
	seen := make(map[string]bool, len(cols))
	for _, ic := range cols {
		if ic == nil || ic.Functional || ic.WithoutOverlaps {
			continue
		}
		ref, ok := ic.Expr.(*ast.ColumnRef)
		if !ok || ref.Star || ref.Table != "" || ref.Schema != "" {
			continue
		}
		k, ok := b.key(spaceColumn, ref.Column)
		if !ok {
			continue
		}
		if seen[k] {
			b.fail(rangeOf(ic.Loc), nil, review.DuplicateName, statesTwice("names", "column", ident(ref.Column)))
			return true
		}
		seen[k] = true
		if !declared[k] {
			b.fail(rangeOf(ic.Loc), nil, review.ObjectNotFound, notFound("column", ident(ref.Column), tableName))
			return true
		}
	}
	return false
}

// columnDefFails checks column-definition text invariants: inline REFERENCES
// referencing more than one column (1239) and DEFAULT / GENERATED ALWAYS AS
// expressions referencing the column itself (4029).
func (b *builder) columnDefFails(col *ast.ColumnDef, tableName string) bool {
	if col == nil {
		return false
	}
	for _, cc := range col.Constraints {
		if cc != nil && cc.Type == ast.ColConstrReferences && len(cc.RefColumns) > 1 {
			b.fail(rangeOf(cc.Loc), nil, review.UnsupportedForm, "Foreign key on "+review.Code(tableName)+" has a mismatched column count.")
			return true
		}
	}
	colKeyName, ok := b.key(spaceColumn, col.Name)
	if !ok {
		return false
	}
	if col.DefaultValue != nil && exprReferencesColumn(col.DefaultValue, colKeyName) {
		b.fail(rangeOf(col.Loc), nil, review.UnsupportedForm, "Expression for column "+review.Code(ident(col.Name))+" cannot reference itself.")
		return true
	}
	if col.Generated != nil && col.Generated.Expr != nil && exprReferencesColumn(col.Generated.Expr, colKeyName) {
		b.fail(rangeOf(col.Loc), nil, review.UnsupportedForm, "Expression for column "+review.Code(ident(col.Name))+" cannot reference itself.")
		return true
	}
	return false
}

// checkExprTableColumns checks that bare column references in a table-scoped
// expression (CHECK, DEFAULT) exist in the table (1054), skipping columns in
// skip (such as columns added earlier in the same ADD (...) subcommand).
func (b *builder) checkExprTableColumns(r tableRef, expr ast.ExprNode, skip map[string]bool) {
	if expr == nil {
		return
	}
	for _, ref := range exprBareColumnRefs(expr) {
		k, ok := b.key(spaceColumn, ref.Column)
		if !ok || skip[k] {
			continue
		}
		b.need(colKey(r.db, r.name, k), walk.Present, nil, rangeOf(ref.Loc), notFound("column", ident(ref.Column), r.written), "")
	}
}

func exprReferencesColumn(expr ast.ExprNode, targetKey string) bool {
	for _, ref := range exprBareColumnRefs(expr) {
		if strings.EqualFold(ref.Column, targetKey) {
			return true
		}
	}
	return false
}

// exprBareColumnRefs returns unqualified ColumnRefs in expr outside subqueries
// and sequence/special function calls (NEXTVAL, LASTVAL, SETVAL).
func exprBareColumnRefs(expr ast.ExprNode) []*ast.ColumnRef {
	var out []*ast.ColumnRef
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectStmt:
			return false
		case *ast.FuncCallExpr:
			switch strings.ToUpper(x.Name) {
			case "NEXTVAL", "LASTVAL", "SETVAL":
				return false
			}
		case *ast.ColumnRef:
			if !x.Star && x.Schema == "" && x.Table == "" && x.Column != "" {
				out = append(out, x)
			}
			return false
		}
		return true
	})
	return out
}

// selectBareColumnNoFrom returns the first bare ColumnRef in the target list
// of a SELECT branch that has no FROM clause, outside subqueries and function
// calls. Such a reference always fails in MariaDB with 1054.
func selectBareColumnNoFrom(sel *ast.SelectStmt) *ast.ColumnRef {
	sel = ast.UnwrapParenSource(sel)
	if sel == nil {
		return nil
	}
	if sel.SetOp != 0 {
		if ref := selectBareColumnNoFrom(sel.Left); ref != nil {
			return ref
		}
		return selectBareColumnNoFrom(sel.Right)
	}
	if sel.ValuesSource != nil || sel.TableSource != nil || len(sel.From) > 0 {
		return nil
	}
	for _, target := range sel.TargetList {
		expr := target
		if rt, ok := target.(*ast.ResTarget); ok {
			expr = rt.Val
		}
		if ref := bareColumnRefOutsideFuncs(expr); ref != nil {
			return ref
		}
	}
	return nil
}

// eventScheduleColumnRef returns the first bare ColumnRef in an event's ON
// SCHEDULE expressions (AT, EVERY, STARTS, ENDS), which have no tables in
// scope and fail in MariaDB with 1054 before event existence is checked.
func eventScheduleColumnRef(sched *ast.EventSchedule) *ast.ColumnRef {
	if sched == nil {
		return nil
	}
	for _, expr := range []ast.ExprNode{sched.At, sched.Every, sched.Starts, sched.Ends} {
		if ref := bareColumnRefOutsideFuncs(expr); ref != nil {
			return ref
		}
	}
	return nil
}

// bareColumnRefOutsideFuncs finds a bare ColumnRef in expr while skipping
// subqueries and function-call arguments (which may hold sequence names or
// temporal unit keywords).
func bareColumnRefOutsideFuncs(expr ast.ExprNode) *ast.ColumnRef {
	if expr == nil {
		return nil
	}
	var found *ast.ColumnRef
	ast.Inspect(expr, func(n ast.Node) bool {
		if found != nil {
			return false
		}
		switch x := n.(type) {
		case *ast.SelectStmt, *ast.FuncCallExpr:
			return false
		case *ast.ColumnRef:
			if !x.Star && x.Column != "" {
				found = x
			}
			return false
		}
		return true
	})
	return found
}
