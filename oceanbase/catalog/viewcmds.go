package catalog

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf16"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/deparse"
)

// createView applies CREATE VIEW. source is the text the statement was parsed
// from, or "" when there is none (see deparseViewSelect). checkBody refuses a
// view body the engine rejects (a `?` parameter marker); it is true on the
// CREATE VIEW path (Exec) and false on a DefineView sync, whose body the engine
// already accepted.
func (c *Catalog) createView(stmt *nodes.CreateViewStmt, source string, checkBody, materialized bool) error {
	if err := rejectRemoteRef(stmt.Name); err != nil {
		return err
	}
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	// The query is analyzed once: here when the body is checked, below otherwise.
	var analyzedQuery *Query
	analyzed := false
	if checkBody {
		q, err := c.resolveQueryFirst(stmt.Select, true, db)
		if err != nil {
			return err
		}
		analyzedQuery, analyzed = q, true
	}
	key := toLower(stmt.Name.Name)
	// Tables and views share the same namespace in MySQL.
	if _, exists := db.Tables[key]; exists {
		return errDupTable(stmt.Name.Name)
	}
	if existing, exists := db.Views[key]; exists {
		if !stmt.OrReplace {
			return errDupTable(stmt.Name.Name)
		}
		if existing.Materialized != materialized {
			kind := "VIEW"
			if materialized {
				kind = "MATERIALIZED VIEW"
			}
			return &Error{Code: 1347, SQLState: "HY000", Message: fmt.Sprintf("'%s.%s' is not %s", db.Name, stmt.Name.Name, kind)}
		}
	}

	// A `?` parameter marker is not allowed in a view body; OceanBase rejects
	// CREATE VIEW with one. The grammar permits it (it is a LIMIT operand), so
	// the rejection lives here, not in the parser, and must fail Exec rather
	// than be swallowed as the analyzer error below is.
	if checkBody && stmt.Select != nil && containsParamMarker(stmt.Select) {
		return &Error{
			Code:     1351,
			SQLState: "HY000",
			Message:  "View's SELECT contains a '?' (parameter marker)",
		}
	}

	// opt_table_id pins the view's id, which the pinned server parses and refuses
	// to run. This catalog cannot model that id, so the statement fails rather
	// than install a view the server would not have made. It is checked after the
	// body, so the errors the body itself raises come first, as they do on the
	// server.
	if stmt.HasTableID {
		return errNotSupported("view table id")
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// Analyze the view body for semantic IR (used by lineage, SDL diff).
	// Done before deparseViewSelect because deparse mutates the AST.
	if stmt.Select != nil && !analyzed {
		q, err := c.AnalyzeSelectStmt(stmt.Select)
		if err == nil {
			analyzedQuery = q
		}
		// Swallow error: view analysis may fail for complex views not yet
		// supported by the analyzer. The view still gets created with
		// Definition text but without AnalyzedQuery.
	}

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	aliased := aliasedTargets(stmt.Select)
	derivedCols, finalizeDef := c.prepareViewSelect(stmt.Select, stmt.Columns, stmt.SelectText, source, db)

	// A column list is folded into the body's aliases (prepareViewSelect), so the view
	// states none of its own.
	const hasExplicit = false
	viewCols := derivedCols
	if name, dup := duplicateViewColumn(stmt.Select, stmt.Columns, aliased); dup {
		return errDupColumn(name)
	}
	columnMetadata, err := c.inferViewColumnMetadata(stmt.Select, viewCols, db)
	if err != nil {
		return err
	}
	definition := finalizeDef()

	// OceanBase takes ALGORITHM and SQL SECURITY and keeps neither, and reports every view
	// updatable (TestOracle_OceanBaseViewAttributes).
	db.Views[key] = &View{
		Name:            stmt.Name.Name,
		Database:        db,
		Definition:      definition,
		Definer:         definer,
		CheckOption:     stmt.CheckOption,
		Columns:         viewCols,
		ColumnMetadata:  columnMetadata,
		ExplicitColumns: hasExplicit,
		AnalyzedQuery:   analyzedQuery,
		IsUpdatable:     true,
		Materialized:    materialized,
		namesUncertain:  viewColumnNamesUncertain(stmt.Columns, viewCols, analyzedQuery, stmt.Select),
	}
	return nil
}

// createMaterializedView applies CREATE MATERIALIZED VIEW. The catalog holds no
// materialized view of its own, so it installs the view the query is and marks
// it materialized; a snapshot's materialized view installs the same way
// (metadata_load.go). The refresh options, partitioning, table options and
// primary key change nothing the catalog tracks.
func (c *Catalog) createMaterializedView(stmt *nodes.CreateMaterializedViewStmt, source string) error {
	if err := c.createView(materializedViewAsView(stmt), source, true, true); err != nil {
		return err
	}
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	installed := db.Views[key]
	if stmt.PrimaryKey != nil {
		for _, name := range stmt.PrimaryKey.Columns {
			found := false
			for _, out := range installed.Columns {
				if strings.EqualFold(out, name) {
					found = true
					break
				}
			}
			if !found {
				delete(db.Views, key)
				return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn), Message: fmt.Sprintf("Unknown column '%s' in 'field list'", name)}
			}
		}
		cols := make([]*IndexColumn, 0, len(stmt.PrimaryKey.Columns))
		for _, name := range stmt.PrimaryKey.Columns {
			cols = append(cols, &IndexColumn{Name: name})
		}
		idxType := strings.ToUpper(stmt.PrimaryKey.IndexType)
		if idxType == "" {
			idxType = "BTREE"
		}
		installed.Indexes = append(installed.Indexes, &Index{
			Name:      "PRIMARY",
			Columns:   cols,
			Unique:    true,
			Primary:   true,
			Visible:   true,
			IndexType: idxType,
			Comment:   stmt.PrimaryKey.Comment,
		})
	}
	for _, opt := range stmt.Options {
		if strings.EqualFold(opt.Name, "COMMENT") {
			installed.Comment = opt.Value
		}
	}
	return nil
}

// materializedViewLogTarget checks the table a materialized view log names: a
// log is written for an existing base table, not a view. The catalog models no
// log of its own.
func (c *Catalog) materializedViewLogTarget(ref *nodes.TableRef) error {
	if err := rejectRemoteRef(ref); err != nil {
		return err
	}
	db, err := c.resolveDatabase(ref.Schema)
	if err != nil {
		return err
	}
	key := toLower(ref.Name)
	if db.Views[key] != nil {
		return &Error{Code: 1347, SQLState: "HY000", Message: fmt.Sprintf("'%s.%s' is not TABLE", db.Name, ref.Name)}
	}
	if db.Tables[key] == nil {
		return errUnknownTable(db.Name, ref.Name)
	}
	return nil
}

// createMaterializedViewLog applies CREATE MATERIALIZED VIEW LOG. The log
// itself is not modeled, but its base table must exist and a second log of the
// table is refused unless the statement says OR REPLACE.
func (c *Catalog) createMaterializedViewLog(stmt *nodes.CreateMaterializedViewLogStmt) error {
	if err := c.materializedViewLogTarget(stmt.Table); err != nil {
		return err
	}
	db, err := c.resolveDatabase(stmt.Table.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Table.Name)
	var referenceColumns []string
	if stmt.With != nil {
		referenceColumns = stmt.With.ReferenceColumns
	}
	for _, name := range referenceColumns {
		if db.Tables[key] != nil && db.Tables[key].GetColumn(name) == nil {
			return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn), Message: fmt.Sprintf("Unknown column '%s' in 'field list'", name)}
		}
	}
	if db.MaterializedViewLogsKnown[key] && db.MaterializedViewLogs[key] && !stmt.OrReplace {
		return errDupTable(stmt.Table.Name)
	}
	db.MaterializedViewLogs[key] = true
	db.MaterializedViewLogsKnown[key] = true
	return nil
}

// alterMaterializedViewLog applies ALTER MATERIALIZED VIEW LOG, which needs a
// log of the base table.
func (c *Catalog) alterMaterializedViewLog(stmt *nodes.AlterMaterializedViewLogStmt) error {
	if err := c.materializedViewLogTarget(stmt.Table); err != nil {
		return err
	}
	db, err := c.resolveDatabase(stmt.Table.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Table.Name)
	if db.MaterializedViewLogsKnown[key] && !db.MaterializedViewLogs[key] {
		return errUnknownTable(db.Name, stmt.Table.Name)
	}
	// A successful ALTER proves the log exists.
	db.MaterializedViewLogs[key] = true
	db.MaterializedViewLogsKnown[key] = true
	return nil
}

// dropMaterializedViewLog applies DROP MATERIALIZED VIEW LOG, which needs a
// log of the base table.
func (c *Catalog) dropMaterializedViewLog(stmt *nodes.DropMaterializedViewLogStmt) error {
	if err := c.materializedViewLogTarget(stmt.Table); err != nil {
		return err
	}
	db, err := c.resolveDatabase(stmt.Table.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Table.Name)
	if db.MaterializedViewLogsKnown[key] && !db.MaterializedViewLogs[key] {
		return errUnknownTable(db.Name, stmt.Table.Name)
	}
	delete(db.MaterializedViewLogs, key)
	db.MaterializedViewLogsKnown[key] = true
	return nil
}

// alterMaterializedView applies ALTER MATERIALIZED VIEW. Its actions refresh,
// re-compile or toggle rewrite, which change no query the catalog holds, but
// the target must be a materialized view the catalog has.
func (c *Catalog) alterMaterializedView(stmt *nodes.AlterMaterializedViewStmt) error {
	if err := rejectRemoteRef(stmt.Name); err != nil {
		return err
	}
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	if db.Tables[key] != nil {
		return &Error{Code: 1347, SQLState: "HY000", Message: fmt.Sprintf("'%s.%s' is not MATERIALIZED VIEW", db.Name, stmt.Name.Name)}
	}
	installed, exists := db.Views[key]
	if !exists {
		return errUnknownTable(db.Name, stmt.Name.Name)
	}
	if !installed.Materialized {
		return &Error{Code: 1347, SQLState: "HY000", Message: fmt.Sprintf("'%s.%s' is not MATERIALIZED VIEW", db.Name, stmt.Name.Name)}
	}
	return nil
}

// materializedViewAsView turns CREATE MATERIALIZED VIEW into the CREATE VIEW
// the catalog models.
func materializedViewAsView(stmt *nodes.CreateMaterializedViewStmt) *nodes.CreateViewStmt {
	return &nodes.CreateViewStmt{
		Loc:         stmt.Loc,
		Name:        stmt.Name,
		Columns:     stmt.Columns,
		Select:      stmt.Select,
		SelectText:  stmt.SelectText,
		CheckOption: stmt.CheckOption,
	}
}

// viewColumnNamesUncertain reports whether a view's stored column names may not be
// the ones the engine gives. A column list the statement wrote is the engine's too;
// otherwise the names come from the body, and an analysis that did not finish
// leaves them unknown.
func viewColumnNamesUncertain(columns, viewCols []string, q *Query, sel *nodes.SelectStmt) bool {
	if len(columns) > 0 {
		return false
	}
	if len(viewCols) == 0 || selectHasUnexpandedStar(sel) {
		return true
	}
	if q != nil {
		return q.namesUncertain
	}
	return selectDerivesNames(sel)
}

func selectHasUnexpandedStar(sel *nodes.SelectStmt) bool {
	switch {
	case sel == nil:
		return false
	case sel.SetOp != nodes.SetOpNone:
		return selectHasUnexpandedStar(sel.Left)
	case sel.ParenSource != nil:
		return selectHasUnexpandedStar(sel.ParenSource)
	}
	for _, item := range sel.TargetList {
		expr := item
		if rt, ok := item.(*nodes.ResTarget); ok {
			expr = rt.Val
		}
		switch e := expr.(type) {
		case *nodes.StarExpr:
			return true
		case *nodes.ColumnRef:
			if e.Star {
				return true
			}
		}
	}
	return false
}

// containsParamMarker reports whether the subtree holds a `?` parameter
// marker, at any depth including nested subqueries. It also inspects the rows
// of INSERT ... VALUES and of a VALUES statement, which nodes.Inspect's walker
// skips, so a marker in a stored routine or trigger body cannot slip past.
func containsParamMarker(n nodes.Node) bool {
	found := false
	nodes.Inspect(n, func(x nodes.Node) bool {
		if found {
			return false
		}
		switch e := x.(type) {
		case *nodes.ParamMarkerExpr:
			found = true
		case *nodes.InsertStmt:
			found = rowsHaveParamMarker(e.Values)
		case *nodes.ValuesStmt:
			found = rowsHaveParamMarker(e.Rows)
		}
		return !found
	})
	return found
}

// rowsHaveParamMarker reports whether any expression in the rows of an
// INSERT ... VALUES or a VALUES statement holds a `?` marker. Those rows are
// not reached by nodes.Inspect, so containsParamMarker checks them here.
func rowsHaveParamMarker(rows [][]nodes.ExprNode) bool {
	for _, row := range rows {
		for _, v := range row {
			if v != nil && containsParamMarker(v) {
				return true
			}
		}
	}
	return false
}

// alterView applies ALTER VIEW; source is as for createView.
func (c *Catalog) alterView(stmt *nodes.AlterViewStmt, source string) error {
	if err := rejectRemoteRef(stmt.Name); err != nil {
		return err
	}
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	analyzedQuery, err := c.resolveQueryFirst(stmt.Select, true, db)
	if err != nil {
		return err
	}
	// opt_table_id pins the view's id, which the pinned server parses and refuses
	// to run; this catalog cannot model it.
	if stmt.HasTableID {
		return errNotSupported("view table id")
	}
	key := toLower(stmt.Name.Name)
	// ALTER VIEW requires an ordinary view to exist.
	existing, exists := db.Views[key]
	if !exists {
		return errUnknownTable(db.Name, stmt.Name.Name)
	}
	if existing.Materialized {
		return &Error{Code: 1347, SQLState: "HY000", Message: fmt.Sprintf("'%s.%s' is not VIEW", db.Name, stmt.Name.Name)}
	}

	// A `?` parameter marker is not allowed in a view body; the engine rejects
	// ALTER VIEW with one just as it rejects CREATE VIEW. ALTER VIEW comes only
	// from Exec, so its body is always checked.
	if stmt.Select != nil && containsParamMarker(stmt.Select) {
		return &Error{
			Code:     1351,
			SQLState: "HY000",
			Message:  "View's SELECT contains a '?' (parameter marker)",
		}
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// analyzedQuery, the semantic IR (used by lineage, SDL diff), was read above,
	// before deparseViewSelect mutates the AST.

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	aliased := aliasedTargets(stmt.Select)
	derivedCols, finalizeDef := c.prepareViewSelect(stmt.Select, stmt.Columns, stmt.SelectText, source, db)

	// A column list is folded into the body's aliases (prepareViewSelect), so the view
	// states none of its own.
	const hasExplicit = false
	viewCols := derivedCols
	if name, dup := duplicateViewColumn(stmt.Select, stmt.Columns, aliased); dup {
		return errDupColumn(name)
	}
	columnMetadata, err := c.inferViewColumnMetadata(stmt.Select, viewCols, db)
	if err != nil {
		return err
	}
	definition := finalizeDef()

	// OceanBase takes ALGORITHM and SQL SECURITY and keeps neither, and reports every view
	// updatable (TestOracle_OceanBaseViewAttributes).
	db.Views[key] = &View{
		Name:            stmt.Name.Name,
		Database:        db,
		Definition:      definition,
		Definer:         definer,
		CheckOption:     stmt.CheckOption,
		Columns:         viewCols,
		ColumnMetadata:  columnMetadata,
		ExplicitColumns: hasExplicit,
		AnalyzedQuery:   analyzedQuery,
		namesUncertain:  viewColumnNamesUncertain(stmt.Columns, viewCols, analyzedQuery, stmt.Select),
		IsUpdatable:     true,
	}
	return nil
}

type viewRelationInfo struct {
	key      string
	schema   string
	hasAlias bool
	table    *Table
	optional bool
}

type viewRelationsScope struct {
	rels       []viewRelationInfo
	hasUnknown bool
	maskedCols map[string]bool
}

func (s *viewRelationsScope) withMasked(params []string) *viewRelationsScope {
	if len(params) == 0 {
		return s
	}
	var cp viewRelationsScope
	if s != nil {
		cp = *s
	}
	masked := make(map[string]bool, len(cp.maskedCols)+len(params))
	for k, v := range cp.maskedCols {
		masked[k] = v
	}
	for _, p := range params {
		if p != "" {
			masked[toLower(p)] = true
		}
	}
	cp.maskedCols = masked
	return &cp
}

type viewCollationDerivation int

const (
	viewDerivationExplicit viewCollationDerivation = iota
	viewDerivationImplicit
	viewDerivationCoercible
)

type viewCollationInfo struct {
	charset    string
	collation  string
	derivation viewCollationDerivation
	valid      bool
}

func (c *Catalog) inferViewColumnMetadata(sel *nodes.SelectStmt, names []string, db *Database) ([]ViewColumn, error) {
	cols := make([]ViewColumn, len(names))
	for i, name := range names {
		cols[i].Name = name
	}
	if sel == nil {
		return cols, nil
	}
	relations := &viewRelationsScope{}
	for _, from := range sel.From {
		c.collectViewRelations(from, db, false, relations)
	}
	for i, target := range sel.TargetList {
		if i >= len(cols) {
			break
		}
		rt, ok := target.(*nodes.ResTarget)
		if !ok {
			continue
		}
		cols[i].Nullable = inferViewExprNullable(rt.Val, relations)
		if ci, err := inferViewExprCollation(rt.Val, relations); err != nil {
			return nil, err
		} else if ci.valid {
			cols[i].Charset = ci.charset
			cols[i].Collation = ci.collation
		}
	}
	if err := validateViewSelectCollations(sel, relations); err != nil {
		return nil, err
	}
	return cols, nil
}

func (c *Catalog) collectViewRelations(te nodes.TableExpr, db *Database, optional bool, out *viewRelationsScope) {
	switch n := te.(type) {
	case *nodes.TableRef:
		if n.Catalog != "" || n.Dblink != "" {
			out.hasUnknown = true
			return
		}
		targetDB := db
		if n.Schema != "" {
			if c == nil {
				return
			}
			targetDB = c.GetDatabase(n.Schema)
		}
		if targetDB == nil {
			return
		}
		tbl := targetDB.GetTable(n.Name)
		if tbl == nil {
			return
		}
		key := n.Name
		if n.Alias != "" {
			key = n.Alias
		}
		out.rels = append(out.rels, viewRelationInfo{
			key:      toLower(key),
			schema:   targetDB.Name,
			hasAlias: n.Alias != "",
			table:    tbl,
			optional: optional,
		})
	case *nodes.ExternalTableRef, *nodes.TableCollectionExpr, *nodes.UnnestExpr, *nodes.RbIterateExpr:
		out.hasUnknown = true
	case *nodes.JoinClause:
		// OceanBase reports a column of an outer join's optional side as it is in its
		// table, where MySQL widens it to nullable.
		c.collectViewRelations(n.Left, db, optional, out)
		c.collectViewRelations(n.Right, db, optional, out)
	}
}

func inferViewExprNullable(expr nodes.ExprNode, relations *viewRelationsScope) bool {
	switch e := expr.(type) {
	case *nodes.ColumnRef:
		return inferViewColumnRefNullable(e, relations)
	case *nodes.CaseExpr:
		// A CASE is nullable when a result is, or when it has no ELSE. IF(c, a, b) is
		// one by now (deparse's rewrite), and reads as it did as a function.
		if e.Default == nil || inferViewExprNullable(e.Default, relations) {
			return true
		}
		for _, w := range e.Whens {
			if inferViewExprNullable(w.Result, relations) {
				return true
			}
		}
		return false
	case *nodes.FuncCallExpr:
		switch strings.ToLower(e.Name) {
		// MOD follows its operands on oceanbase-ce 4.3.5.6 (see
		// isAnalyzedExprNullable).
		case "nullif",
			"date_add", "date_sub", "adddate", "subdate",
			"timestampadd", "timestampdiff", "datediff", "timediff", "addtime", "subtime":
			return true
		}
		if (strings.EqualFold(e.Name, "coalesce") || strings.EqualFold(e.Name, "ifnull")) && len(e.Args) > 0 {
			for _, arg := range e.Args {
				if !inferViewExprNullable(arg, relations) {
					return false
				}
			}
			return true
		}
		if strings.EqualFold(e.Name, "if") && len(e.Args) == 3 {
			return inferViewExprNullable(e.Args[1], relations) || inferViewExprNullable(e.Args[2], relations)
		}
		// COUNT reads nullable in a view of OceanBase's, where MySQL's reads NOT NULL.
		if isViewStringMetadataNullableFunc(e.Name) || isViewNullableAggregateFunc(e.Name) || strings.EqualFold(e.Name, "count") {
			return true
		}
		for _, arg := range e.Args {
			if inferViewExprNullable(arg, relations) {
				return true
			}
		}
		return false
	case *nodes.ParenExpr:
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.OdbcEscapeExpr:
		// The escape is a wrapper around its expression, so the result's
		// nullability is the payload's.
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.QuantifiedExpr:
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.ResTarget:
		return inferViewExprNullable(e.Val, relations)
	case *nodes.BinaryExpr:
		if _, ok := e.Left.(*nodes.IntervalExpr); ok {
			return true
		}
		if _, ok := e.Right.(*nodes.IntervalExpr); ok {
			return true
		}
		if e.Op == nodes.BinOpNullSafeEq {
			return false
		}
		return inferViewExprNullable(e.Left, relations) || inferViewExprNullable(e.Right, relations)
	case *nodes.BetweenExpr:
		return inferViewExprNullable(e.Expr, relations) ||
			inferViewExprNullable(e.Low, relations) ||
			inferViewExprNullable(e.High, relations)
	case *nodes.LikeExpr:
		return inferViewExprNullable(e.Expr, relations) ||
			inferViewExprNullable(e.Pattern, relations) ||
			(e.Escape != nil && inferViewExprNullable(e.Escape, relations))
	case *nodes.IsExpr:
		return false
	case *nodes.SubqueryExpr:
		return true
	case *nodes.InExpr:
		if e.Select != nil || inferViewExprNullable(e.Expr, relations) {
			return true
		}
		for _, item := range e.List {
			if inferViewExprNullable(item, relations) {
				return true
			}
		}
		return false
	case *nodes.UnaryExpr:
		return inferViewExprNullable(e.Operand, relations)
	case *nodes.WeightStringExpr:
		// Propagate like the plain WEIGHT_STRING(expr) function-call path.
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.ExtractExpr:
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.IntervalExpr:
		return inferViewExprNullable(e.Value, relations)
	case *nodes.CollateExpr:
		// Nullable whatever its operand, as the engine reports it (see
		// isAnalyzedExprNullable).
		return true
	case *nodes.CastExpr:
		// Covers plain casts and the WEIGHT_STRING(... AS BINARY(n)) desugar,
		// whose WeightStringExpr.Expr is a CastExpr (oracle 8.0.32:
		// IS_NULLABLE=YES for both over a nullable column).
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.PositionExpr:
		// mysql:8.0, whose view metadata this catalog follows: IS_NULLABLE
		// is YES for position('x' in s) over a nullable s and for
		// position('x' in NULL), NO over a NOT NULL column.
		return inferViewExprNullable(e.Substr, relations) ||
			inferViewExprNullable(e.Str, relations)
	case *nodes.NullLit:
		return true
	default:
		return false
	}
}

func (s *viewRelationsScope) matchQualified(schema, table string) (viewRelationInfo, bool) {
	if s == nil {
		return viewRelationInfo{}, false
	}
	lowerTable := toLower(table)
	var matched viewRelationInfo
	found := false
	for _, rel := range s.rels {
		if rel.key != lowerTable {
			continue
		}
		if schema != "" && (rel.hasAlias || !strings.EqualFold(rel.schema, schema)) {
			continue
		}
		if found {
			return viewRelationInfo{}, false
		}
		found = true
		matched = rel
	}
	return matched, found
}

func inferViewColumnRefNullable(cr *nodes.ColumnRef, relations *viewRelationsScope) bool {
	if cr == nil || cr.Star || cr.Catalog != "" {
		return true
	}
	if cr.Table != "" {
		if rel, ok := relations.matchQualified(cr.Schema, cr.Table); ok {
			if col := rel.table.GetColumn(cr.Column); col != nil {
				return rel.optional || col.Nullable
			}
		}
		return true
	}
	if relations == nil || relations.hasUnknown {
		return true
	}
	found := false
	nullable := true
	for _, rel := range relations.rels {
		col := rel.table.GetColumn(cr.Column)
		if col == nil {
			continue
		}
		if found {
			return true
		}
		found = true
		nullable = rel.optional || col.Nullable
	}
	return nullable
}

func validateViewSelectCollations(sel *nodes.SelectStmt, relations *viewRelationsScope) error {
	if sel == nil {
		return nil
	}
	for _, expr := range []nodes.ExprNode{sel.Where, sel.Having} {
		if expr == nil {
			continue
		}
		if _, err := inferViewExprCollation(expr, relations); err != nil {
			return err
		}
	}
	for _, item := range sel.GroupBy {
		if _, err := inferViewExprCollation(item.Expr, relations); err != nil {
			return err
		}
	}
	for _, item := range sel.OrderBy {
		if item == nil {
			continue
		}
		if _, err := inferViewExprCollation(item.Expr, relations); err != nil {
			return err
		}
	}
	return nil
}

func inferViewExprCollation(expr nodes.ExprNode, relations *viewRelationsScope) (viewCollationInfo, error) {
	switch e := expr.(type) {
	case *nodes.ColumnRef:
		return inferViewColumnRefCollation(e, relations), nil
	case *nodes.StringLit:
		charset := normalizeCharsetName(strings.TrimPrefix(toLower(e.Charset), "_"))
		if charset == "" {
			charset = "utf8mb4"
		}
		return newViewCollation(charset, "", viewDerivationCoercible), nil
	case *nodes.CollateExpr:
		child, err := inferViewExprCollation(e.Expr, relations)
		if err != nil {
			return viewCollationInfo{}, err
		}
		collation := toLower(e.Collation)
		charset := normalizeCharsetName(charsetForCollation(collation))
		if charset == "" && child.valid {
			charset = child.charset
		}
		return newViewCollation(charset, collation, viewDerivationExplicit), nil
	case *nodes.ConvertExpr:
		if e.Charset != "" {
			return newViewCollation(normalizeCharsetName(e.Charset), "", viewDerivationImplicit), nil
		}
		return inferViewExprCollation(e.Expr, relations)
	case *nodes.FuncCallExpr:
		return inferViewFuncCollation(e, relations)
	case *nodes.CaseExpr:
		// The conditions compare what they compare, which can fail: CASE x WHEN y
		// compares x with y. They give the CASE no character set.
		operand, err := inferViewExprCollation(e.Operand, relations)
		if err != nil {
			return viewCollationInfo{}, err
		}
		for _, w := range e.Whens {
			cond, err := inferViewExprCollation(w.Cond, relations)
			if err != nil {
				return viewCollationInfo{}, err
			}
			if e.Operand != nil {
				if _, err := aggregateViewCollationsForComparison(operand, cond, "case"); err != nil {
					return viewCollationInfo{}, err
				}
			}
		}
		// The character set of a CASE is that of its results together.
		results := make([]nodes.ExprNode, 0, len(e.Whens)+1)
		for _, w := range e.Whens {
			results = append(results, w.Result)
		}
		if e.Default != nil {
			results = append(results, e.Default)
		}
		var out viewCollationInfo
		for _, result := range results {
			ci, err := inferViewExprCollation(result, relations)
			if err != nil {
				return viewCollationInfo{}, err
			}
			if out, err = aggregateViewCollations(out, ci, "case"); err != nil {
				return viewCollationInfo{}, err
			}
		}
		return out, nil
	case *nodes.ParenExpr:
		return inferViewExprCollation(e.Expr, relations)
	case *nodes.OdbcEscapeExpr:
		// Temporal escapes (`{d '...'}`, `{t '...'}`, `{ts '...'}`) convert a valid
		// string literal into a temporal value with no string charset or collation;
		// transparent escapes (`{fn ...}`, etc.) and non-temporal fallbacks carry
		// the inner expression's metadata and collation checks.
		if isOdbcTemporalEscapeExpr(e) {
			return viewCollationInfo{}, nil
		}
		return inferViewExprCollation(e.Expr, relations)
	case *nodes.QuantifiedExpr:
		return inferViewExprCollation(e.Expr, relations)
	case *nodes.LambdaExpr:
		if _, err := inferViewExprCollation(e.Body, relations.withMasked(e.Params)); err != nil {
			return viewCollationInfo{}, err
		}
		return viewCollationInfo{}, nil
	case *nodes.ResTarget:
		return inferViewExprCollation(e.Val, relations)
	case *nodes.GroupingSet:
		// A grouping set's items are expressions like any other: each is checked
		// for a comparison of collations it cannot make, though the set itself
		// gives the result no character set.
		for _, item := range e.Items {
			if _, err := inferViewExprCollation(item, relations); err != nil {
				return viewCollationInfo{}, err
			}
		}
		return viewCollationInfo{}, nil
	case *nodes.RowExpr:
		for _, item := range e.Items {
			if _, err := inferViewExprCollation(item, relations); err != nil {
				return viewCollationInfo{}, err
			}
		}
		return viewCollationInfo{}, nil
	case *nodes.BinaryExpr:
		return inferViewBinaryCollation(e, relations)
	case *nodes.UnaryExpr:
		return inferViewExprCollation(e.Operand, relations)
	default:
		return viewCollationInfo{}, nil
	}
}

func isOdbcTemporalEscapeExpr(e *nodes.OdbcEscapeExpr) bool {
	if e == nil {
		return false
	}
	kind, rt := odbcTemporalEscapeKind(e.Name)
	if rt == nil {
		return false
	}
	inner := e.Expr
	for {
		pe, ok := inner.(*nodes.ParenExpr)
		if !ok {
			break
		}
		inner = pe.Expr
	}
	lit, ok := inner.(*nodes.StringLit)
	return ok && odbcTemporalValid(kind, lit.Value)
}

func inferViewFuncCollation(fn *nodes.FuncCallExpr, relations *viewRelationsScope) (viewCollationInfo, error) {
	name := strings.ToLower(fn.Name)
	switch name {
	case "concat", "concat_ws":
		var out viewCollationInfo
		for _, arg := range fn.Args {
			ci, err := inferViewExprCollation(arg, relations)
			if err != nil {
				return viewCollationInfo{}, err
			}
			agg, err := aggregateViewCollations(out, ci, name)
			if err != nil {
				return viewCollationInfo{}, err
			}
			out = agg
		}
		return out, nil
	case "repeat":
		if len(fn.Args) == 0 {
			return viewCollationInfo{}, nil
		}
		return inferViewExprCollation(fn.Args[0], relations)
	case "lpad", "rpad":
		// OceanBase gives the result the character set of the string and of the pad,
		// where MySQL pins it to the string's.
		var out viewCollationInfo
		for i, arg := range fn.Args {
			if i == 1 {
				continue // the length
			}
			ci, err := inferViewExprCollation(arg, relations)
			if err != nil {
				return viewCollationInfo{}, err
			}
			agg, err := aggregateViewCollations(out, ci, name)
			if err != nil {
				return viewCollationInfo{}, err
			}
			out = agg
		}
		return out, nil
	default:
		var out viewCollationInfo
		for _, arg := range fn.Args {
			ci, err := inferViewExprCollation(arg, relations)
			if err != nil {
				return viewCollationInfo{}, err
			}
			agg, err := aggregateViewCollations(out, ci, name)
			if err != nil {
				return viewCollationInfo{}, err
			}
			out = agg
		}
		return out, nil
	}
}

func inferViewBinaryCollation(expr *nodes.BinaryExpr, relations *viewRelationsScope) (viewCollationInfo, error) {
	left, err := inferViewExprCollation(expr.Left, relations)
	if err != nil {
		return viewCollationInfo{}, err
	}
	right, err := inferViewExprCollation(expr.Right, relations)
	if err != nil {
		return viewCollationInfo{}, err
	}
	op := binaryOpToString(expr.Op)
	if isViewCollationComparisonOp(expr.Op) {
		if _, err := aggregateViewCollationsForComparison(left, right, op); err != nil {
			return viewCollationInfo{}, err
		}
		return viewCollationInfo{}, nil
	}
	return aggregateViewCollations(left, right, op)
}

func inferViewColumnRefCollation(cr *nodes.ColumnRef, relations *viewRelationsScope) viewCollationInfo {
	if cr == nil || cr.Star || cr.Catalog != "" {
		return viewCollationInfo{}
	}
	if cr.Table != "" {
		if rel, ok := relations.matchQualified(cr.Schema, cr.Table); ok {
			return viewColumnCollation(rel.table.GetColumn(cr.Column))
		}
		return viewCollationInfo{}
	}
	if relations == nil || relations.hasUnknown || relations.maskedCols[toLower(cr.Column)] {
		return viewCollationInfo{}
	}
	var out viewCollationInfo
	found := false
	for _, rel := range relations.rels {
		ci := viewColumnCollation(rel.table.GetColumn(cr.Column))
		if !ci.valid {
			continue
		}
		if found {
			return viewCollationInfo{}
		}
		found = true
		out = ci
	}
	return out
}

func viewColumnCollation(col *Column) viewCollationInfo {
	if col == nil || col.Charset == "" || !isStringType(col.DataType) {
		return viewCollationInfo{}
	}
	return newViewCollation(col.Charset, col.Collation, viewDerivationImplicit)
}

func newViewCollation(charset, collation string, derivation viewCollationDerivation) viewCollationInfo {
	charset = normalizeCharsetName(charset)
	collation = toLower(collation)
	if charset == "" && collation != "" {
		charset = normalizeCharsetName(charsetForCollation(collation))
	}
	if collation == "" && charset != "" {
		collation = defaultCollationForCharset[toLower(charset)]
	}
	if charset == "" || collation == "" {
		return viewCollationInfo{}
	}
	return viewCollationInfo{
		charset:    charset,
		collation:  collation,
		derivation: derivation,
		valid:      true,
	}
}

func aggregateViewCollations(left, right viewCollationInfo, op string) (viewCollationInfo, error) {
	if !left.valid {
		return right, nil
	}
	if !right.valid {
		return left, nil
	}
	if strings.EqualFold(left.collation, right.collation) {
		return left, nil
	}
	if left.derivation == viewDerivationExplicit && right.derivation == viewDerivationExplicit {
		return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
	}
	if left.derivation != right.derivation {
		if left.derivation < right.derivation {
			return left, nil
		}
		return right, nil
	}
	if isCharsetSuperset(left.charset, right.charset) {
		return left, nil
	}
	if isCharsetSuperset(right.charset, left.charset) {
		return right, nil
	}
	// MySQL 8.0 permits several same-charset string function mixes that the
	// comparison path rejects. Preserve the left side for view metadata.
	if strings.EqualFold(left.charset, right.charset) {
		return left, nil
	}
	return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
}

func aggregateViewCollationsForComparison(left, right viewCollationInfo, op string) (viewCollationInfo, error) {
	if !left.valid || !right.valid || strings.EqualFold(left.collation, right.collation) {
		if left.valid {
			return left, nil
		}
		return right, nil
	}
	if left.derivation == viewDerivationExplicit && right.derivation == viewDerivationExplicit {
		return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
	}
	if left.derivation == right.derivation && strings.EqualFold(left.charset, right.charset) {
		return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
	}
	return aggregateViewCollations(left, right, op)
}

func isViewCollationComparisonOp(op nodes.BinaryOp) bool {
	switch op {
	case nodes.BinOpEq, nodes.BinOpNe, nodes.BinOpLt, nodes.BinOpGt, nodes.BinOpLe, nodes.BinOpGe, nodes.BinOpNullSafeEq:
		return true
	default:
		return false
	}
}

func isCharsetSuperset(charset, other string) bool {
	charset = normalizeCharsetName(charset)
	other = normalizeCharsetName(other)
	if charset == other {
		return true
	}
	switch charset {
	case "utf8mb4":
		return other == "latin1" || other == "ascii" || other == "utf8mb3" || other == "utf8"
	case "utf8mb3":
		return other == "ascii"
	case "latin1":
		return other == "ascii"
	default:
		return false
	}
}

func errIllegalMixCollations(left, right viewCollationInfo, op string) error {
	return &Error{
		Code:     1267,
		SQLState: "HY000",
		Message: fmt.Sprintf("Illegal mix of collations (%s,%s) and (%s,%s) for operation '%s'",
			left.collation, viewDerivationName(left.derivation),
			right.collation, viewDerivationName(right.derivation), op),
	}
}

func viewDerivationName(derivation viewCollationDerivation) string {
	switch derivation {
	case viewDerivationExplicit:
		return "EXPLICIT"
	case viewDerivationImplicit:
		return "IMPLICIT"
	case viewDerivationCoercible:
		return "COERCIBLE"
	default:
		return "NONE"
	}
}

func isViewStringMetadataNullableFunc(name string) bool {
	switch strings.ToLower(name) {
	case "concat", "concat_ws", "ifnull", "coalesce":
		return true
	default:
		return false
	}
}

func isViewNullableAggregateFunc(name string) bool {
	switch strings.ToLower(name) {
	case "sum", "avg", "min", "max", "group_concat":
		return true
	default:
		return false
	}
}

func (c *Catalog) dropView(stmt *nodes.DropViewStmt) error {
	// Every name is checked before any view is dropped, IF EXISTS or not, so a
	// wrong kind late in the list leaves the catalog as it was.
	kindWord := "VIEW"
	if stmt.Materialized {
		kindWord = "MATERIALIZED VIEW"
	}
	seen := map[[2]string]bool{}
	for _, ref := range stmt.Views {
		if err := rejectRemoteRef(ref); err != nil {
			return err
		}
		db, err := c.resolveDatabase(ref.Schema)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		key := toLower(ref.Name)
		if seen[[2]string{toLower(db.Name), key}] {
			continue
		}
		seen[[2]string{toLower(db.Name), key}] = true
		if db.Tables[key] != nil {
			return &Error{Code: 1347, SQLState: "HY000", Message: fmt.Sprintf("'%s.%s' is not %s", db.Name, ref.Name, kindWord)}
		}
		installed, exists := db.Views[key]
		if !exists {
			if stmt.IfExists {
				continue
			}
			return errUnknownTable(db.Name, ref.Name)
		}
		if installed.Materialized != stmt.Materialized {
			return &Error{Code: 1347, SQLState: "HY000", Message: fmt.Sprintf("'%s.%s' is not %s", db.Name, ref.Name, kindWord)}
		}
	}
	for _, ref := range stmt.Views {
		db, err := c.resolveDatabase(ref.Schema)
		if err != nil {
			continue
		}
		key := toLower(ref.Name)
		if _, exists := db.Views[key]; exists {
			delete(db.Views, key)
		}
	}
	return nil
}

// deparseViewSelect resolves, rewrites, and deparses the SELECT AST for a view.
// If the AST is nil (parser didn't produce one), falls back to the raw SelectText.
// Returns the deparsed definition and the derived column names from the resolved
// SELECT target list (used when no explicit column list is specified).
//
// source is the text the statement was parsed from, which the AST's Loc
// offsets index into. An unaliased select item is named after its text in
// source, as the engine names it; with no source ("" from DefineView), the
// name is rebuilt from the AST.
//
// columns is the view's column list. OceanBase folds it into the body, as the alias of
// each select item in every SELECT the body joins: CREATE VIEW v (p, q) AS SELECT a, b
// FROM t is stored as select a AS p, b AS q (TestOracle_ViewIdempotence, explicit-columns
// and explicit-columns-union).
func (c *Catalog) deparseViewSelect(sel *nodes.SelectStmt, columns []string, rawText, source string, db *Database) (string, []string) {
	derivedCols, finalize := c.prepareViewSelect(sel, columns, rawText, source, db)
	return finalize(), derivedCols
}

// prepareViewSelect resolves and rewrites sel for a view, returning the derived
// column names and a finalize callback that applies OceanBase stored-view
// literal normalization (charset introducer stripping and _binary hex folding)
// and deparses to SQL. Splitting finalize lets callers run
// inferViewColumnMetadata while StringLit.Charset is still intact.
func (c *Catalog) prepareViewSelect(sel *nodes.SelectStmt, columns []string, rawText, source string, db *Database) ([]string, func() string) {
	if sel == nil {
		return columns, func() string { return rawText }
	}

	// Build a TableLookup that resolves table names from this database.
	lookup := tableLookupForDB(db)

	// Determine the database charset for CAST resolution.
	charset := db.Charset
	if charset == "" {
		charset = c.defaultCharset
	}

	// Resolve: qualify columns, expand *, normalize JOINs.
	resolver := &deparse.Resolver{
		Lookup:         lookup,
		DefaultCharset: charset,
		Source:         source,
	}
	resolver.Resolve(sel)

	_, bindings := aliasViewColumns(sel, columns)

	// Extract column names from the resolved target list.
	derivedCols := extractViewColumns(sel)

	// Rewrite: NOT folding, boolean context wrapping.
	deparse.RewriteSelectStmt(sel)
	bindings.substitute()

	// Deparse: AST → canonical SQL text.
	return derivedCols, func() string {
		normalizeOceanBaseViewLiterals(sel)
		return deparse.DeparseSelect(sel)
	}
}

// normalizeOceanBaseViewLiterals applies OceanBase 4.3.5.6's stored-view
// literal normalization: non-binary charset introducers are stripped from
// string literals (_utf8's' -> 's') and _binary string literals are stored as
// hex literals (_binary's' -> X'73').
func normalizeOceanBaseViewLiterals(sel *nodes.SelectStmt) {
	if sel == nil {
		return
	}
	rewriteViewExprs(sel)
}

func rewriteViewExprs(sel *nodes.SelectStmt) {
	if sel == nil {
		return
	}
	for _, cte := range sel.CTEs {
		if cte != nil {
			rewriteViewExprs(cte.Select)
		}
	}
	if sel.ParenSource != nil {
		rewriteViewExprs(sel.ParenSource)
	}
	if sel.Left != nil {
		rewriteViewExprs(sel.Left)
	}
	if sel.Right != nil {
		rewriteViewExprs(sel.Right)
	}
	for i, t := range sel.TargetList {
		if rt, ok := t.(*nodes.ResTarget); ok {
			rt.Val = normalizeViewLiteralExpr(rt.Val)
		} else {
			sel.TargetList[i] = normalizeViewLiteralExpr(t)
		}
	}
	for _, te := range sel.From {
		normalizeViewTableExpr(te)
	}
	if sel.Where != nil {
		sel.Where = normalizeViewLiteralExpr(sel.Where)
	}
	for _, item := range sel.GroupBy {
		item.Expr = normalizeViewLiteralExpr(item.Expr)
	}
	if sel.Having != nil {
		sel.Having = normalizeViewLiteralExpr(sel.Having)
	}
	for _, wd := range sel.WindowClause {
		normalizeViewWindowDef(wd)
	}
	for _, item := range sel.OrderBy {
		item.Expr = normalizeViewLiteralExpr(item.Expr)
	}
	if sel.Limit != nil {
		sel.Limit.Count = normalizeViewLiteralExpr(sel.Limit.Count)
		sel.Limit.Offset = normalizeViewLiteralExpr(sel.Limit.Offset)
	}
	if sel.ValuesSource != nil {
		for _, row := range sel.ValuesSource.Rows {
			for i, item := range row {
				row[i] = normalizeViewLiteralExpr(item)
			}
		}
	}
}

func normalizeViewWindowDef(wd *nodes.WindowDef) {
	if wd == nil {
		return
	}
	for i, expr := range wd.PartitionBy {
		wd.PartitionBy[i] = normalizeViewLiteralExpr(expr)
	}
	for _, item := range wd.OrderBy {
		item.Expr = normalizeViewLiteralExpr(item.Expr)
	}
	if wd.Frame != nil {
		if wd.Frame.Start != nil {
			wd.Frame.Start.Offset = normalizeViewLiteralExpr(wd.Frame.Start.Offset)
		}
		if wd.Frame.End != nil {
			wd.Frame.End.Offset = normalizeViewLiteralExpr(wd.Frame.End.Offset)
		}
	}
}

func normalizeViewTableExpr(te nodes.TableExpr) {
	switch t := te.(type) {
	case *nodes.JoinClause:
		normalizeViewTableExpr(t.Left)
		normalizeViewTableExpr(t.Right)
		if on, ok := t.Condition.(*nodes.OnCondition); ok {
			on.Expr = normalizeViewLiteralExpr(on.Expr)
		}
	case *nodes.SubqueryExpr:
		rewriteViewExprs(t.Select)
	case *nodes.JsonTableExpr:
		t.Expr = normalizeViewLiteralExpr(t.Expr)
	case *nodes.UnnestExpr:
		for i, arg := range t.Args {
			t.Args[i] = normalizeViewLiteralExpr(arg)
		}
	case *nodes.RbIterateExpr:
		t.Arg = normalizeViewLiteralExpr(t.Arg)
	case *nodes.TableCollectionExpr:
		t.Expr = normalizeViewLiteralExpr(t.Expr)
	}
}

func normalizeViewLiteralExpr(expr nodes.ExprNode) nodes.ExprNode {
	if expr == nil {
		return nil
	}
	switch n := expr.(type) {
	case *nodes.StringLit:
		if n.Charset != "" {
			cs := strings.ToLower(strings.TrimPrefix(n.Charset, "_"))
			switch cs {
			case "binary":
				return &nodes.HexLit{
					Loc:   n.Loc,
					Value: fmt.Sprintf("%X", []byte(n.Value)),
				}
			case "utf16le":
				// OceanBase 4.3.5.6 reinterprets _utf16le string bytes as UTF-16LE
				// (left-padding odd-length byte runs with a leading 0x00 byte) when
				// storing the literal in SHOW CREATE VIEW.
				raw := []byte(n.Value)
				if len(raw)%2 == 1 {
					padded := make([]byte, len(raw)+1)
					copy(padded[1:], raw)
					raw = padded
				}
				u16 := make([]uint16, len(raw)/2)
				for i := range u16 {
					u16[i] = binary.LittleEndian.Uint16(raw[2*i:])
				}
				n.Value = string(utf16.Decode(u16))
			}
			n.Charset = ""
		}
		return n
	case *nodes.ParenExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
	case *nodes.UnaryExpr:
		n.Operand = normalizeViewLiteralExpr(n.Operand)
	case *nodes.BinaryExpr:
		n.Left = normalizeViewLiteralExpr(n.Left)
		n.Right = normalizeViewLiteralExpr(n.Right)
	case *nodes.FuncCallExpr:
		for i, a := range n.Args {
			n.Args[i] = normalizeViewLiteralExpr(a)
		}
		for _, item := range n.OrderBy {
			item.Expr = normalizeViewLiteralExpr(item.Expr)
		}
		n.Separator = normalizeViewLiteralExpr(n.Separator)
		normalizeViewWindowDef(n.Over)
	case *nodes.CaseExpr:
		n.Operand = normalizeViewLiteralExpr(n.Operand)
		for _, w := range n.Whens {
			w.Cond = normalizeViewLiteralExpr(w.Cond)
			w.Result = normalizeViewLiteralExpr(w.Result)
		}
		n.Default = normalizeViewLiteralExpr(n.Default)
	case *nodes.CastExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
	case *nodes.ConvertExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
	case *nodes.CollateExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
	case *nodes.InExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
		for i, item := range n.List {
			n.List[i] = normalizeViewLiteralExpr(item)
		}
		rewriteViewExprs(n.Select)
	case *nodes.BetweenExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
		n.Low = normalizeViewLiteralExpr(n.Low)
		n.High = normalizeViewLiteralExpr(n.High)
	case *nodes.LikeExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
		n.Pattern = normalizeViewLiteralExpr(n.Pattern)
		n.Escape = normalizeViewLiteralExpr(n.Escape)
	case *nodes.IsExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
	case *nodes.PositionExpr:
		n.Substr = normalizeViewLiteralExpr(n.Substr)
		n.Str = normalizeViewLiteralExpr(n.Str)
	case *nodes.ExtractExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
	case *nodes.WeightStringExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
	case *nodes.IntervalExpr:
		n.Value = normalizeViewLiteralExpr(n.Value)
	case *nodes.RowExpr:
		for i, item := range n.Items {
			n.Items[i] = normalizeViewLiteralExpr(item)
		}
	case *nodes.MemberOfExpr:
		n.Value = normalizeViewLiteralExpr(n.Value)
		n.Array = normalizeViewLiteralExpr(n.Array)
	case *nodes.MatchExpr:
		n.Against = normalizeViewLiteralExpr(n.Against)
	case *nodes.QuantifiedExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
	case *nodes.OdbcEscapeExpr:
		n.Expr = normalizeViewLiteralExpr(n.Expr)
	case *nodes.LambdaExpr:
		n.Body = normalizeViewLiteralExpr(n.Body)
	case *nodes.GroupingSet:
		for i, item := range n.Items {
			n.Items[i] = normalizeViewLiteralExpr(item)
		}
	case *nodes.ResTarget:
		n.Val = normalizeViewLiteralExpr(n.Val)
	case *nodes.SubqueryExpr:
		rewriteViewExprs(n.Select)
	case *nodes.ExistsExpr:
		rewriteViewExprs(n.Select)
	}
	return expr
}

// aliasedTargets returns the items of a view's select list that are written with an alias,
// read before the resolver gives every other item the name it derives.
func aliasedTargets(sel *nodes.SelectStmt) map[*nodes.ResTarget]bool {
	aliased := map[*nodes.ResTarget]bool{}
	if sel == nil {
		return aliased
	}
	for _, target := range nodes.LeftmostQueryLeaf(sel).TargetList {
		if rt, ok := target.(*nodes.ResTarget); ok && rt.Name != "" {
			aliased[rt] = true
		}
	}
	return aliased
}

// duplicateViewColumn returns a name two of a view's columns have, whatever its case, which
// OceanBase refuses as MySQL does (ER_DUP_FIELDNAME, 1060). The names are the view's column
// list or, without one, the names written in the resolved select list: an alias (aliased)
// or a column's own name, a column of * among them. A name the engine makes up for an
// expression is not held against the view: SELECT 1, 1 is taken.
// TestOracle_OceanBaseRefusals holds each form to the engine.
func duplicateViewColumn(sel *nodes.SelectStmt, explicit []string, aliased map[*nodes.ResTarget]bool) (string, bool) {
	names := explicit
	if len(names) == 0 && sel != nil {
		for _, target := range nodes.LeftmostQueryLeaf(sel).TargetList {
			switch t := target.(type) {
			case *nodes.ResTarget:
				if _, column := t.Val.(*nodes.ColumnRef); column || aliased[t] {
					names = append(names, t.Name)
				}
			case *nodes.ColumnRef:
				names = append(names, t.Column)
			}
		}
	}
	seen := make(map[string]bool, len(names))
	for _, name := range names {
		if seen[toLower(name)] {
			return name, true
		}
		seen[toLower(name)] = true
	}
	return "", false
}

// aliasViewColumns gives each SELECT a query joins the names of a view's column list, as
// the aliases of its select items. A list of another length is left for the caller.
//
// What names an item by the alias it had has to follow it: the resolver has written
// ORDER BY 1 as that alias by now, and the old name is no column of the body. On a set
// operation or parentheses, an ORDER BY names the query's output columns, so it takes the
// new names. In a SELECT, an ORDER BY item that is the alias alone takes the new name too,
// and becomes the item's position (deparse's orderByPosition). Anywhere else in a SELECT
// the new name could be a column of a table, which a bare name binds to first, so the
// alias stands for the item's expression, as OceanBase writes it: ORDER BY x + 1, for a
// written a AS x, is stored as order by (t.a + 1). That substitution waits for the rewrite
// pass (aliasBindings.substitute), so no expression is rewritten twice. A query nested in
// any of them is a scope of its own, whose aliases are not these.
func aliasViewColumns(sel *nodes.SelectStmt, columns []string) (map[string]string, aliasBindings) {
	if sel == nil || len(columns) == 0 {
		return nil, nil
	}
	var renamed map[string]string
	var bindings aliasBindings
	switch {
	case sel.ParenSource != nil:
		renamed, bindings = aliasViewColumns(sel.ParenSource, columns)
	case sel.SetOp != nodes.SetOpNone:
		renamed, bindings = aliasViewColumns(sel.Left, columns)
		_, right := aliasViewColumns(sel.Right, columns)
		bindings = append(bindings, right...)
	case len(sel.TargetList) == len(columns):
		renamed = map[string]string{}
		leaf := aliasBinding{sel: sel, items: map[string]aliasItem{}}
		for i, target := range sel.TargetList {
			rt, ok := target.(*nodes.ResTarget)
			if !ok {
				sel.TargetList[i] = &nodes.ResTarget{Name: columns[i], Val: target}
				continue
			}
			if old := toLower(rt.Name); old != "" && old != toLower(columns[i]) {
				if _, twice := renamed[old]; !twice {
					renamed[old] = columns[i]
					leaf.items[old] = aliasItem{target: rt, position: i + 1}
				}
			}
			rt.Name = columns[i]
		}
		for _, item := range sel.OrderBy {
			if name, ok := aliasReference(item.Expr); ok {
				if to, ok := renamed[name]; ok {
					item.Expr.(*nodes.ColumnRef).Column = to
				}
			}
		}
		if len(leaf.items) > 0 {
			bindings = aliasBindings{leaf}
		}
		return renamed, bindings
	default:
		return nil, nil
	}
	for _, item := range sel.OrderBy {
		followAliases(item.Expr, renamed)
	}
	return renamed, bindings
}

// aliasBindings are the SELECTs whose select items a view's column list renamed.
type aliasBindings []aliasBinding

type aliasBinding struct {
	sel   *nodes.SelectStmt
	items map[string]aliasItem // old alias, lower-cased, to the item it named
}

type aliasItem struct {
	target   *nodes.ResTarget
	position int
}

// substitute puts each renamed item's expression where its old alias stands in the
// SELECT's GROUP BY, HAVING and ORDER BY. A GROUP BY item that is the alias alone, of an
// item that is a constant, becomes the item's position: a constant there would read as
// one. A query nested in them is a scope of its own and is left alone.
func (bs aliasBindings) substitute() {
	for _, b := range bs {
		for _, item := range b.sel.GroupBy {
			if name, ok := aliasReference(item.Expr); ok {
				if it, ok := b.items[name]; ok && isConstantLiteral(it.target.Val) {
					item.Expr = &nodes.IntLit{Value: int64(it.position)}
					continue
				}
			}
			item.Expr = b.substituteIn(item.Expr)
		}
		b.sel.Having = b.substituteIn(b.sel.Having)
		for _, item := range b.sel.OrderBy {
			item.Expr = b.substituteIn(item.Expr)
		}
	}
}

// substituteIn returns expr with every reference to a renamed alias replaced by its
// item's expression. It reaches every expression a node holds, in a field, a list or a
// clause of its own (a CASE's WHEN, a window's ORDER BY), and stops at a nested query.
func (b aliasBinding) substituteIn(expr nodes.ExprNode) nodes.ExprNode {
	if expr == nil {
		return nil
	}
	if name, ok := aliasReference(expr); ok {
		if it, ok := b.items[name]; ok {
			return it.target.Val
		}
		return expr
	}
	b.substituteFields(reflect.ValueOf(expr))
	return expr
}

var exprNodeType = reflect.TypeOf((*nodes.ExprNode)(nil)).Elem()

// substituteFields substitutes in the expression fields of the node v points to. Only
// fields of expression type are set: a Loc, a name or a flag is never written.
func (b aliasBinding) substituteFields(v reflect.Value) {
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return
	}
	if _, nested := v.Interface().(*nodes.SelectStmt); nested {
		return
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}
		switch t := f.Type(); {
		case t == exprNodeType:
			b.substituteValue(f)
		case t.Kind() == reflect.Slice && t.Elem() == exprNodeType:
			for j := 0; j < f.Len(); j++ {
				b.substituteValue(f.Index(j))
			}
		case t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct:
			b.substituteFields(f)
		case t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Pointer && t.Elem().Elem().Kind() == reflect.Struct:
			for j := 0; j < f.Len(); j++ {
				b.substituteFields(f.Index(j))
			}
		}
	}
}

func (b aliasBinding) substituteValue(f reflect.Value) {
	if f.IsNil() {
		return
	}
	f.Set(reflect.ValueOf(b.substituteIn(f.Interface().(nodes.ExprNode))))
}

// aliasReference returns the lower-cased name a bare column reference holds: one written
// with no table, which the resolver leaves only on an alias.
func aliasReference(expr nodes.ExprNode) (string, bool) {
	cr, ok := expr.(*nodes.ColumnRef)
	if !ok || cr.Table != "" || cr.Schema != "" || cr.Star {
		return "", false
	}
	return toLower(cr.Column), true
}

func isConstantLiteral(expr nodes.ExprNode) bool {
	switch expr.(type) {
	case *nodes.IntLit, *nodes.FloatLit, *nodes.StringLit, *nodes.NullLit, *nodes.BoolLit, *nodes.BitLit:
		return true
	}
	return false
}

// followAliases rewrites the references to renamed aliases in expr to the new names. A
// query nested in expr is a scope of its own and is left alone.
func followAliases(expr nodes.ExprNode, renamed map[string]string) {
	if expr == nil || len(renamed) == 0 {
		return
	}
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if _, nested := n.(*nodes.SelectStmt); nested {
			return false
		}
		if name, ok := aliasReference(asExpr(n)); ok {
			if to, ok := renamed[name]; ok {
				n.(*nodes.ColumnRef).Column = to
			}
		}
		return true
	})
}

func asExpr(n nodes.Node) nodes.ExprNode {
	e, _ := n.(nodes.ExprNode)
	return e
}

// extractViewColumns extracts column names from a resolved SELECT target list.
// This produces the column list that MySQL would derive for a view.
func extractViewColumns(sel *nodes.SelectStmt) []string {
	if sel == nil {
		return nil
	}
	// The result columns live on the leftmost output SELECT. Unwrap any
	// parenthesized-query wrappers and walk set-operation left arms to reach it,
	// so a parenthesized or set-op view body (CREATE VIEW v AS (SELECT a, b) or
	// (SELECT 1 AS a) UNION (SELECT 2 AS a)) derives columns from that SELECT.
	sel = nodes.LeftmostQueryLeaf(sel)
	var cols []string
	for _, target := range sel.TargetList {
		if rt, ok := target.(*nodes.ResTarget); ok {
			if rt.Name != "" {
				cols = append(cols, rt.Name)
			} else if cr, ok := rt.Val.(*nodes.ColumnRef); ok {
				cols = append(cols, cr.Column)
			} else {
				cols = append(cols, nodeToSQL(rt.Val))
			}
		} else if cr, ok := target.(*nodes.ColumnRef); ok {
			cols = append(cols, cr.Column)
		} else {
			cols = append(cols, nodeToSQL(target))
		}
	}
	return cols
}

// tableLookupForDB returns a deparse.TableLookup function that resolves table
// and view names from the given database's Tables and Views maps.
func tableLookupForDB(db *Database) deparse.TableLookup {
	return func(tableName string) *deparse.ResolverTable {
		key := toLower(tableName)
		// Try tables first.
		tbl := db.Tables[key]
		if tbl != nil {
			cols := make([]deparse.ResolverColumn, len(tbl.Columns))
			for i, c := range tbl.Columns {
				cols[i] = deparse.ResolverColumn{
					Name:     c.Name,
					Position: c.Position,
				}
			}
			return &deparse.ResolverTable{
				Name:    tbl.Name,
				Columns: cols,
			}
		}
		// Fall back to views.
		v := db.Views[key]
		if v != nil {
			cols := make([]deparse.ResolverColumn, len(v.Columns))
			for i, colName := range v.Columns {
				cols[i] = deparse.ResolverColumn{
					Name:     colName,
					Position: i + 1,
				}
			}
			return &deparse.ResolverTable{
				Name:    v.Name,
				Columns: cols,
			}
		}
		return nil
	}
}

// ShowCreateView produces MySQL 8.0-compatible SHOW CREATE VIEW output.
// Returns "" if the database or view does not exist.
func (c *Catalog) ShowCreateView(database, name string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	v := db.Views[toLower(name)]
	if v == nil {
		return ""
	}
	return showCreateView(v)
}

// formatDefiner ensures the definer string is backtick-quoted per MySQL 8.0 format.
// Input can be: `root`@`%`, root@%, 'root'@'%', etc.
// Output: `root`@`%`
func formatDefiner(definer string) string {
	// If already formatted with backticks, return as-is.
	if strings.HasPrefix(definer, "`") && strings.Contains(definer, "@") {
		return definer
	}
	// Split on @
	parts := strings.SplitN(definer, "@", 2)
	if len(parts) == 1 {
		// No @ — just backtick-quote the whole thing.
		return "`" + strings.Trim(parts[0], "`'") + "`"
	}
	user := strings.Trim(parts[0], "`'")
	host := strings.Trim(parts[1], "`'")
	return fmt.Sprintf("`%s`@`%s`", user, host)
}

// showCreateView produces the SHOW CREATE VIEW output for a view.
// MySQL 8.0 format:
//
//	CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `view_name` AS select_statement
//	WITH CASCADED CHECK OPTION
func showCreateView(v *View) string {
	var b strings.Builder

	b.WriteString("CREATE")

	// ALGORITHM — MySQL 8.0 always shows ALGORITHM, defaults to UNDEFINED.
	algorithm := v.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	b.WriteString(fmt.Sprintf(" ALGORITHM=%s", strings.ToUpper(algorithm)))

	// DEFINER — MySQL 8.0 always shows DEFINER with backtick-quoted user@host.
	if v.Definer != "" {
		b.WriteString(fmt.Sprintf(" DEFINER=%s", formatDefiner(v.Definer)))
	}

	// SQL SECURITY — MySQL 8.0 always shows SQL SECURITY, defaults to DEFINER.
	sqlSecurity := v.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}
	b.WriteString(fmt.Sprintf(" SQL SECURITY %s", strings.ToUpper(sqlSecurity)))

	// VIEW name
	b.WriteString(" VIEW " + backtickIDQ(v.Name))

	// Column list (only if explicitly specified by user in CREATE VIEW).
	if v.ExplicitColumns && len(v.Columns) > 0 {
		cols := make([]string, len(v.Columns))
		for i, c := range v.Columns {
			cols[i] = backtickIDQ(c)
		}
		b.WriteString(fmt.Sprintf(" (%s)", strings.Join(cols, ",")))
	}

	// AS select_statement
	b.WriteString(" AS ")
	b.WriteString(v.Definition)

	// WITH CHECK OPTION
	if v.CheckOption != "" {
		b.WriteString(fmt.Sprintf(" WITH %s CHECK OPTION", strings.ToUpper(v.CheckOption)))
	}

	return b.String()
}

// resolveQueryFirst refuses a view's or CREATE TABLE ... SELECT's query that
// violates OceanBase's view/query semantic rules: pre-analysis checks (syntax,
// VALUES, windows, charsets/collations, JSON_QUERY paths, SQL_CALC_FOUND_ROWS),
// table/column/star/alias resolution (1146, 1054, 1051, 1109, 1096, 1066, 1052,
// 1582), and post-analysis checks (view-only variables 1351, remote sequences
// 5835, vector APPROX 1235). OceanBase resolves the query before it looks at
// the name being created or altered, so these come before 1050, and before
// ALTER VIEW's 1146 or 1347 for a name that is no view: CREATE VIEW a AS
// SELECT zz FROM a, with a a table, answers 1054
// (TestOracle_OceanBaseQueryFirst).
//
// The analyzer's other errors are left alone, as createView has always left them:
// it does not yet read every query OceanBase takes. The analyzed query is returned,
// or nil when analysis did not finish, so the caller need not analyze it again.
func (c *Catalog) resolveQueryFirst(sel *nodes.SelectStmt, isView bool, targetDB *Database) (*Query, error) {
	if sel == nil {
		return nil, nil
	}
	blocks := collectSelectBlocks(sel)
	if err := validateViewQueryPreAnalysis(blocks); err != nil {
		return nil, err
	}
	q, err := c.AnalyzeSelectStmt(sel)
	var refusal *Error
	if errors.As(err, &refusal) {
		switch refusal.Code {
		case ErrNoSuchTable, ErrNoSuchColumn, ErrUnknownTable, ErrUnknownTableInFieldList,
			ErrNoTablesUsed, 1066, 1052, ErrWrongParamCount, 1583, ErrWrongArguments, ErrUnknownCharacterSet:
			return nil, refusal
		}
	}
	if postErr := c.validateViewQueryPostAnalysis(blocks, isView, targetDB); postErr != nil {
		return nil, postErr
	}
	if err != nil {
		return nil, nil
	}
	return q, nil
}

// collectSelectBlocks collects all SelectStmt query blocks in sel in order,
// including set-op arms, parenthesized queries, CTE bodies, derived tables,
// and subqueries in expressions (including inside VALUES ROW(...) items, which
// nodes.Inspect skips).
func collectSelectBlocks(root *nodes.SelectStmt) []*nodes.SelectStmt {
	var blocks []*nodes.SelectStmt
	var walk func(nodes.Node)
	walk = func(node nodes.Node) {
		nodes.Inspect(node, func(n nodes.Node) bool {
			if s, ok := n.(*nodes.SelectStmt); ok {
				blocks = append(blocks, s)
				if s.ValuesSource != nil {
					for _, row := range s.ValuesSource.Rows {
						for _, item := range row {
							if item != nil {
								walk(item)
							}
						}
					}
				}
			}
			return true
		})
	}
	walk(root)
	return blocks
}

// inspectBlockOwnNodes visits AST nodes belonging directly to block without
// descending into nested SelectStmt query blocks or CTE definitions.
// Also walks ValuesSource.Rows, which nodes.Inspect skips.
func inspectBlockOwnNodes(block *nodes.SelectStmt, fn func(nodes.Node) bool) {
	walk := func(root nodes.Node) {
		nodes.Inspect(root, func(n nodes.Node) bool {
			switch e := n.(type) {
			case *nodes.SelectStmt:
				if e != block {
					return false
				}
			case *nodes.CommonTableExpr:
				return false
			}
			return fn(n)
		})
	}
	walk(block)
	if block.ValuesSource != nil {
		for _, row := range block.ValuesSource.Rows {
			for _, item := range row {
				if item != nil {
					walk(item)
				}
			}
		}
	}
}

func validateViewQueryPreAnalysis(blocks []*nodes.SelectStmt) error {
	for _, block := range blocks {
		for _, hint := range block.OptimizerHints {
			if err := validateOptimizerHint(hint); err != nil {
				return err
			}
		}
		if block.CalcFoundRows {
			return errWrongUsage("SQL_CALC_FOUND_ROWS")
		}
		if block.ValuesSource != nil {
			if block.ValuesSource.LeadingComma {
				return errEmptyRowInTVC()
			}
			for _, row := range block.ValuesSource.Rows {
				if len(row) == 0 {
					return errEmptyRowInTVC()
				}
				for _, item := range row {
					if d, ok := item.(*nodes.DefaultExpr); ok && d.Column == "" {
						return errDefaultAsValInTVC()
					}
				}
			}
		}
		// Collect WINDOW clause declarations and check for duplicate names.
		declaredWindows := make(map[string]bool, len(block.WindowClause))
		for _, wd := range block.WindowClause {
			if wd == nil {
				continue
			}
			if wd.Name != "" {
				key := toLower(wd.Name)
				if declaredWindows[key] {
					return &Error{Code: ErrNotSupportedYet, SQLState: "0A000", Message: "duplicate window name not supported"}
				}
				declaredWindows[key] = true
			}
		}
		for _, wd := range block.WindowClause {
			if err := validateWindowDefPreAnalysis(wd, declaredWindows); err != nil {
				return err
			}
		}

		var err error
		inspectBlockOwnNodes(block, func(n nodes.Node) bool {
			if err != nil {
				return false
			}
			switch e := n.(type) {
			case *nodes.FuncCallExpr:
				if e.Schema == "" {
					lower := strings.ToLower(e.Name)
					switch lower {
					case "_st_asmvt", "check":
						err = errOBSyntaxError()
						return false
					case "nth_value":
						// OceanBase 4.3.5.6 validates that NTH_VALUE's second argument
						// is a constant before resolving window names (5929).
						if len(e.Args) == 2 && exprContainsColumnRef(e.Args[1]) {
							err = &Error{
								Code:     ErrWrongArguments,
								SQLState: sqlState(ErrWrongArguments),
								Message:  "Incorrect arguments to nth_value",
							}
							return false
						}
					case "calc_partition_id":
						// OceanBase 4.3.5.6 validates CALC_PARTITION_ID's first
						// argument as a constant before resolving column references.
						if len(e.Args) < 2 || len(e.Args) > 3 || exprContainsColumnRef(e.Args[0]) {
							err = &Error{
								Code:     ErrWrongArguments,
								SQLState: sqlState(ErrWrongArguments),
								Message:  "Invalid argument",
							}
							return false
						}
					}
				}
				if derr := checkDistinctAllowed(e); derr != nil {
					err = derr
					return false
				}
				if e.Over != nil {
					if werr := validateWindowDefPreAnalysis(e.Over, declaredWindows); werr != nil {
						err = werr
						return false
					}
				}
				if e.JSONQueryArm && len(e.Args) >= 2 {
					if lit, ok := e.Args[1].(*nodes.StringLit); ok && !isValidJSONPathLit(lit.Value) {
						err = errOBJSONPathSyntax(lit.Value)
						return false
					}
				}
				if e.Using != "" && !e.UsingQuoted {
					if cerr := checkCharsetName(e.Using); cerr != nil {
						err = cerr
						return false
					}
				}
			case *nodes.CastExpr:
				if e.Ignore {
					err = errOBSyntaxError()
					return false
				}
				if e.TypeName != nil && e.TypeName.Charset != "" {
					if cerr := checkCharsetName(e.TypeName.Charset); cerr != nil {
						err = cerr
						return false
					}
				}
			case *nodes.ConvertExpr:
				// A quoted name is checked after the argument (analyzeExpr).
				if e.Charset != "" && !e.CharsetQuoted {
					if cerr := checkCharsetName(e.Charset); cerr != nil {
						err = cerr
						return false
					}
				}
				if e.TypeName != nil && e.TypeName.Charset != "" {
					if cerr := checkCharsetName(e.TypeName.Charset); cerr != nil {
						err = cerr
						return false
					}
				}
			case *nodes.CollateExpr:
				if cerr := checkCollationName(e.Collation); cerr != nil {
					err = cerr
					return false
				}
			case *nodes.MatchExpr:
				if strings.Contains(strings.ToUpper(e.Modifier), "QUERY EXPANSION") {
					err = &Error{Code: ErrNotSupportedYet, SQLState: "0A000", Message: "search modes other than NATURAL_LANGUAGE_MODE or BOOLEAN_MODE not supported"}
					return false
				}
			}
			return true
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func validateWindowDefPreAnalysis(wd *nodes.WindowDef, declared map[string]bool) error {
	if wd == nil {
		return nil
	}
	if wd.RefName != "" && !declared[toLower(wd.RefName)] {
		return errOBWindowNotDefined(wd.RefName)
	}
	if wd.Frame != nil {
		f := wd.Frame
		if f.Type == nodes.FrameRows {
			if (f.Start != nil && isIntervalExpr(f.Start.Offset)) ||
				(f.End != nil && isIntervalExpr(f.End.Offset)) {
				return errWindowRangeFrameNumericType()
			}
		}
		if f.End == nil && f.Start != nil &&
			(f.Start.Type == nodes.BoundFollowing || f.Start.Type == nodes.BoundUnboundedFollowing) {
			return errOBSyntaxError()
		}
	}
	return nil
}

func isIntervalExpr(e nodes.ExprNode) bool {
	for {
		p, ok := e.(*nodes.ParenExpr)
		if !ok {
			break
		}
		e = p.Expr
	}
	_, ok := e.(*nodes.IntervalExpr)
	return ok
}

func isValidJSONPathLit(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0 && s[0] == '$'
}

func (c *Catalog) validateViewQueryPostAnalysis(blocks []*nodes.SelectStmt, isView bool, targetDB *Database) error {
	for _, block := range blocks {
		if block.Approx && !hasVectorSortExpr(block.OrderBy) {
			return &Error{Code: ErrNotSupportedYet, SQLState: "0A000", Message: "Using vector index without vector_sort_expr is not supported"}
		}
		var err error
		inspectBlockOwnNodes(block, func(n nodes.Node) bool {
			if err != nil {
				return false
			}
			switch e := n.(type) {
			case *nodes.FuncCallExpr:
				if ferr := c.validateViewFuncExists(e, targetDB); ferr != nil {
					err = ferr
					return false
				}
			case *nodes.VariableRef:
				if isView {
					err = errViewSelectVariable()
					return false
				}
			case *nodes.RemoteSequenceExpr:
				err = errOBRemoteDBNotFound()
				return false
			}
			return true
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) validateViewFuncExists(fc *nodes.FuncCallExpr, targetDB *Database) error {
	if fc == nil || fc.Name == "" || c.loadingSDL {
		return nil
	}
	lower := toLower(fc.Name)
	if fc.Schema != "" {
		if isSystemDatabaseName(fc.Schema) {
			return nil
		}
		db := c.GetDatabase(fc.Schema)
		if db == nil || db.Functions[lower] == nil {
			return errNoSuchFunction(fc.Name)
		}
		return nil
	}
	if isBuiltinOceanBaseFunction(lower) {
		return nil
	}
	curDB := c.GetDatabase(c.CurrentDatabase())
	if curDB == nil {
		curDB = targetDB
	}
	if curDB != nil && curDB.Functions[lower] != nil {
		return nil
	}
	return errNoSuchFunction(fc.Name)
}

func hasVectorSortExpr(orderBy []*nodes.OrderByItem) bool {
	for _, item := range orderBy {
		if item == nil {
			continue
		}
		found := false
		nodes.Inspect(item.Expr, func(n nodes.Node) bool {
			if found {
				return false
			}
			if fn, ok := n.(*nodes.FuncCallExpr); ok {
				switch strings.ToLower(fn.Name) {
				case "l2_distance", "cosine_distance", "inner_product", "negative_inner_product",
					"l1_distance", "l2_squared", "approx_l2_distance", "approx_cosine_distance",
					"approx_inner_product", "approx_l1_distance", "approx_l2_squared",
					"vec_l2_distance", "vec_cosine_distance", "vec_inner_product",
					"vec_negative_inner_product", "vec_l1_distance", "vec_l2_squared":
					found = true
				}
			}
			return !found
		})
		if found {
			return true
		}
	}
	return false
}
