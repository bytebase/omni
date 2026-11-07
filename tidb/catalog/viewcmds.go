package catalog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/deparse"
)

// createView applies CREATE VIEW. source is the text the statement was parsed
// from, or "" when there is none (see deparseViewSelect). checkBody refuses a
// body TiDB cannot build, as CREATE VIEW does; DefineView installs a view
// TiDB already has, whose body a later ALTER TABLE may have broken, and
// leaves it unset.
func (c *Catalog) createView(stmt *nodes.CreateViewStmt, source string, checkBody bool) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	// Tables and views share the same namespace in MySQL.
	if _, exists := db.Tables[key]; exists {
		return errDupTable(stmt.Name.Name)
	}
	if _, exists := db.Views[key]; exists {
		if !stmt.OrReplace {
			return errDupTable(stmt.Name.Name)
		}
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// Analyze the view body for semantic IR (used by lineage, SDL diff).
	// Done before deparseViewSelect because deparse mutates the AST.
	var analyzedQuery *Query
	if stmt.Select != nil {
		q, err := c.AnalyzeSelectStmt(stmt.Select)
		if err == nil {
			analyzedQuery = q
		}
		// Swallow error: view analysis may fail for complex views not yet
		// supported by the analyzer. The view still gets created with
		// Definition text but without AnalyzedQuery.
	}

	// Derive each column's collation, and refuse a body whose collations do
	// not aggregate. Also before deparseViewSelect, which expands stars and
	// qualifies tables in place.
	if checkBody {
		// A `?` parameter marker is not allowed in a view body; TiDB rejects
		// CREATE VIEW with one (ER_VIEW_SELECT_VARIABLE). The grammar permits
		// it (it is a LIMIT operand), so the rejection lives here, not in the
		// parser, and must fail Exec rather than be swallowed as the analyzer
		// error is.
		if stmt.Select != nil && containsParamMarker(stmt.Select) {
			return &Error{
				Code:     1351,
				SQLState: "HY000",
				Message:  "View's SELECT contains a '?' (parameter marker)",
			}
		}
		if err := c.checkViewSelectCollations(stmt.Select, c.viewBodyDatabase(db), nil); err != nil {
			return err
		}
	}
	collations := c.viewTargetCollations(stmt.Select, c.viewBodyDatabase(db))

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	definition, derivedCols := c.deparseViewSelect(stmt.Select, stmt.SelectText, source, db)

	// Use explicit column list if provided, otherwise derive from SELECT target list.
	hasExplicit := len(stmt.Columns) > 0
	viewCols := stmt.Columns
	if !hasExplicit {
		viewCols = derivedCols
	}
	if err := checkViewColumnNames(viewCols); err != nil {
		return err
	}

	db.Views[key] = &View{
		Name:            stmt.Name.Name,
		Database:        db,
		Definition:      definition,
		Algorithm:       viewAlgorithm(stmt.Algorithm),
		Definer:         definer,
		SqlSecurity:     viewSQLSecurity(stmt.SqlSecurity),
		CheckOption:     stmt.CheckOption,
		Columns:         viewCols,
		ExplicitColumns: hasExplicit,
		AnalyzedQuery:   analyzedQuery,
		ColumnMetadata:  viewColumnMetadata(viewCols, collations),
	}
	return nil
}

// containsParamMarker reports whether the subtree holds a `?` parameter
// marker, at any depth including nested subqueries.
func containsParamMarker(n nodes.Node) bool {
	found := false
	nodes.Inspect(n, func(x nodes.Node) bool {
		if _, ok := x.(*nodes.ParamMarkerExpr); ok {
			found = true
			return false
		}
		return !found
	})
	return found
}

func (c *Catalog) dropView(stmt *nodes.DropViewStmt) error {
	for _, ref := range stmt.Views {
		db, err := c.resolveDatabase(ref.Schema)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		key := toLower(ref.Name)
		if _, exists := db.Views[key]; !exists {
			if stmt.IfExists {
				continue
			}
			return errUnknownTable(db.Name, ref.Name)
		}
		delete(db.Views, key)
	}
	return nil
}

// deparseViewSelect builds the definition TiDB stores for a view body and
// the view's column names. If the AST is nil (parser didn't produce one),
// falls back to the raw SelectText.
//
// TiDB stores the body as its parser restores it after the planner has
// prepared it (pingcap/tidb pkg/ddl/create_table.go, BuildViewInfo; see
// deparse.RestoreSelect): every * expanded to schema-qualified columns, every
// select field named, every table that is not a CTE qualified by its schema,
// and otherwise the query as written, with keywords in upper case.
//
// source is the text the statement was parsed from, which the AST's Loc
// offsets index into. An unaliased select item is named after its text in
// source, as the engine names it; with no source ("" from DefineView), the
// name is rebuilt from the AST.
// viewBodyDatabase returns the database an unqualified table in the body of
// a view created in db names. It is the session's, not the view's:
// USE a; CREATE VIEW b.v AS SELECT id FROM t stores FROM `a`.`t` (verified on
// v8.5.5). With no database selected the view's own is the only one there is.
func (c *Catalog) viewBodyDatabase(db *Database) *Database {
	if cur := c.GetDatabase(c.CurrentDatabase()); cur != nil {
		return cur
	}
	return db
}

func (c *Catalog) deparseViewSelect(sel *nodes.SelectStmt, rawText, source string, db *Database) (string, []string) {
	if sel == nil {
		return rawText, nil
	}

	from := c.viewBodyDatabase(db)
	lookup := tableLookupForDB(from)
	resolver := &deparse.Resolver{
		Lookup:   lookup,
		Source:   source,
		Restore:  true,
		Database: from.Name,
		// A table written with its schema is that schema's: its columns,
		// and the widths CAST AS CHAR takes from them, are not the
		// session's table of the same name.
		LookupIn: func(schema, name string) *deparse.ResolverTable {
			db := c.GetDatabase(schema)
			if db == nil {
				return nil
			}
			return tableLookupForDB(db)(name)
		},
	}
	resolver.Resolve(sel)

	return deparse.RestoreSelect(sel), viewColumnNames(extractViewColumns(sel, lookup))
}

// viewAlgorithm returns the algorithm TiDB stores for a view: the one
// declared, UNDEFINED when none is. TiDB keeps MERGE even for a view that
// cannot merge, where MySQL downgrades it.
func viewAlgorithm(declared string) string {
	if declared == "" {
		return "UNDEFINED"
	}
	return strings.ToUpper(declared)
}

// viewSQLSecurity returns the SQL SECURITY TiDB stores for a view: the one
// declared, DEFINER when none is.
func viewSQLSecurity(declared string) string {
	if declared == "" {
		return "DEFINER"
	}
	return strings.ToUpper(declared)
}

// checkViewColumnNames refuses a view whose columns repeat a name, compared
// without case, as TiDB does (error 1060, Duplicate column name): SELECT
// t1.a, t2.a cannot be a view.
func checkViewColumnNames(cols []string) error {
	seen := make(map[string]bool, len(cols))
	for _, c := range cols {
		k := strings.ToLower(c)
		if seen[k] {
			return errDupColumn(c)
		}
		seen[k] = true
	}
	return nil
}

// viewColumnNames returns the column names TiDB gives a view from the names
// of its select fields: a name longer than a column name may be becomes
// name_exp_N, N the column's position (planner/core
// adjustOverlongViewColname). The field keeps its long name in the body.
func viewColumnNames(names []string) []string {
	for i, name := range names {
		if len(strings.ToLower(name)) > maxColumnNameLength {
			names[i] = fmt.Sprintf("name_exp_%d", i+1)
		}
	}
	return names
}

// maxColumnNameLength is the longest column name TiDB accepts, in bytes.
const maxColumnNameLength = 64

// extractViewColumns extracts column names from a resolved SELECT target list.
// This produces the column list that MySQL would derive for a view. lookup
// resolves the table of a TABLE query primary, which has no target list to
// read.
func extractViewColumns(sel *nodes.SelectStmt, lookup deparse.TableLookup) []string {
	if sel == nil {
		return nil
	}
	// The result columns live on the leftmost output SELECT. Unwrap any
	// parenthesized-query wrappers and walk set-operation left arms to reach it,
	// so a parenthesized or set-op view body (CREATE VIEW v AS (SELECT a, b) or
	// (SELECT 1 AS a) UNION (SELECT 2 AS a)) derives columns from that SELECT.
	// A TABLE query primary is SELECT * FROM t and a VALUES primary a row
	// constructor, so neither has a target list to read. The resolver leaves
	// the primary alone — TiDB's SHOW CREATE VIEW keeps the body as
	// TABLE `t` — so its columns are derived here, by the same code the
	// resolver uses for a CTE or derived table over one, which searches the
	// query's own CTE scope before lookup's tables.
	//
	// Container-verified: CREATE VIEW v AS TABLE t declares the columns of t,
	// and WITH x AS (SELECT 1 AS a) TABLE x declares a. A VALUES primary's
	// columns are ValuesColumnName's; TiDB cannot run a VALUES query at all
	// -- every form answers error 1051, Unknown table '', a view of one
	// included -- so the names are MySQL's, whose VALUES statement TiDB's
	// grammar implements.
	if primary := deparse.PrimaryOutputColumns(sel, lookup); primary != nil {
		cols := make([]string, 0, len(primary))
		for _, c := range primary {
			cols = append(cols, c.Name)
		}
		return cols
	}
	sel = nodes.LeftmostQueryLeaf(sel)
	var cols []string
	for _, target := range sel.TargetList {
		switch t := target.(type) {
		case *nodes.ResTarget:
			// The resolver names every field; only an empty string
			// literal keeps an empty name, which TiDB gives its column.
			if cr, ok := t.Val.(*nodes.ColumnRef); ok && t.Name == "" {
				cols = append(cols, cr.Column)
			} else {
				cols = append(cols, t.Name)
			}
		case *nodes.ColumnRef:
			// An unaliased column parses as a bare column reference.
			if !t.Star {
				cols = append(cols, t.Column)
			}
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
				length, integer := castCharLength(c.ColumnType)
				cols[i] = deparse.ResolverColumn{
					Name:       c.Name,
					Position:   c.Position,
					CharLength: length,
					Integer:    integer,
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

// castCharLength returns the length TiDB gives CAST(column AS CHAR) for a
// column of the given type, which is the type's display width, and whether
// the type is an integer; 0 when the width is not known. The widths are the
// ones v8.5.5 stores in a view body, whatever display width the column
// declares: INT 11, INT UNSIGNED 10, TINYINT 4,
// BIGINT 20, VARCHAR(10) 10, DECIMAL(10,2) 12, DATETIME 19, DATE 10,
// TIME(2) 13, TEXT and BLOB 262140, DOUBLE 370.
func castCharLength(columnType string) (int, bool) {
	m := columnTypeRE.FindStringSubmatch(strings.ToLower(strings.TrimSpace(columnType)))
	if m == nil {
		return 0, false
	}
	name, unsigned := m[1], strings.Contains(m[4], "unsigned")
	length, _ := strconv.Atoi(m[2])
	scale, _ := strconv.Atoi(m[3])
	intWidths := map[string][2]int{
		"tinyint": {4, 3}, "smallint": {6, 5}, "mediumint": {9, 8},
		"int": {11, 10}, "integer": {11, 10}, "bigint": {20, 20},
	}
	if w, ok := intWidths[name]; ok {
		// The display width does not bound the value, so TiDB ignores it
		// (expression/builtin_cast.go, adjustRetFtForCastString): INT(1)
		// ZEROFILL casts as CHAR(10), TINYINT(1) as CHAR(4).
		if unsigned {
			return w[1], true
		}
		return w[0], true
	}
	switch name {
	case "char", "varchar", "binary", "varbinary":
		if length > 0 {
			return length, false
		}
		if name == "char" || name == "binary" {
			return 1, false
		}
	case "text", "blob":
		return 262140, false
	case "decimal", "numeric":
		if length == 0 {
			length = 10
		}
		width := length + 1
		if scale > 0 {
			width++
		}
		return width, false
	case "datetime", "timestamp":
		if length > 0 {
			return 20 + length, false
		}
		return 19, false
	case "time":
		if length > 0 {
			return 11 + length, false
		}
		return 10, false
	case "date":
		return 10, false
	case "double":
		return 370, false
	}
	return 0, false
}

// columnTypeRE splits a column type into its name, length, scale, and the
// attributes after them.
var columnTypeRE = regexp.MustCompile(`^([a-z]+)(?:\s*\(\s*(\d+)\s*(?:,\s*(\d+)\s*)?\))?(.*)$`)

// ShowCreateView produces TiDB's SHOW CREATE VIEW output.
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

// showCreateView produces the SHOW CREATE VIEW output for a view, as TiDB
// writes it (pkg/executor/show.go, fetchShowCreateTable4View):
//
//	CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`a`, `b`) AS SELECT ...
//
// The column list is always written, explicit or derived, and the CHECK
// OPTION never is: TiDB stores it but SHOW CREATE VIEW leaves it out.
func showCreateView(v *View) string {
	var b strings.Builder

	algorithm := v.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	fmt.Fprintf(&b, "CREATE ALGORITHM=%s", strings.ToUpper(algorithm))

	definer := v.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}
	fmt.Fprintf(&b, " DEFINER=%s", formatDefiner(definer))

	sqlSecurity := v.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}
	fmt.Fprintf(&b, " SQL SECURITY %s", strings.ToUpper(sqlSecurity))

	b.WriteString(" VIEW " + backtickIDQ(v.Name) + " (")
	for i, c := range v.Columns {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(backtickIDQ(c))
	}
	b.WriteString(") AS ")
	b.WriteString(v.Definition)
	return b.String()
}
