package review

import (
	"context"
	"errors"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// viewQueryCodes are the errors oracle/catalog finds in a view's query that
// hold whatever else the schema holds, once the tables the query reads are
// known in full: each is pinned by oracle/catalog's query suite
// (testdata/queries/simulate.tsv) on Oracle 23. Left out are the errors an
// object the sync does not list can cause or cure: a missing name (00942,
// and 00904 unless namesNoRoutine), a user-defined aggregate or a
// function's type (viewGroupCodes unless the query calls no user
// function), a collection (22905), and periods, flashback and database
// links.
var viewQueryCodes = map[string]bool{
	"00913": true, // too many values
	"00918": true, // column ambiguously defined
	"00957": true, // duplicate column name
	"00964": true, // schema name not in FROM list
	"00998": true, // must name this expression with a column alias
	"01027": true, // bind variables not allowed for data definition operations
	"01468": true, // a predicate may reference only one outer-joined table
	"01785": true, // ORDER BY item must be the number of a SELECT-list expression
	"01789": true, // query block has incorrect number of result columns
	"14501": true, // object is not partitioned
	"18307": true, // analytic view does not exist
	"18370": true, // HIERARCHIES clause may only be specified for cubes
	"18385": true, // recursive WITH clause in a derived analytic view
	"22848": true, // cannot use a LOB type as a comparison key
	"22849": true, // a LOB in an aggregate
	"22902": true, // CURSOR expression not allowed
	"25154": true, // column part of USING clause cannot have qualifier
	"30483": true, // window functions are not allowed here
	"30732": true, // table contains no user-visible columns
	"32038": true, // number of WITH clause column names does not match
	"32039": true, // recursive WITH clause must have column alias list
	"54058": true, // '*' without FROM
	"56901": true, // non-constant expression in PIVOT values
	"62503": true, // illegal variable definition
}

// viewGroupCodes are the group-function and datatype errors. The catalog
// reads a call it does not know as a scalar function of no known type,
// which a user-defined aggregate is not, so they hold only for a query
// that calls no user function.
// A name that is no column fails first with 00904, which the walk does not
// report, so a function called without parentheses cannot reach them.
var viewGroupCodes = map[string]bool{
	"00932": true, // inconsistent datatypes
	"00934": true, // group function is not allowed here
	"00937": true, // not a single-group group function
	"00978": true, // nested group function without GROUP BY
	"00979": true, // not a GROUP BY expression
	"01790": true, // expression must have same datatype as corresponding expression
	"56902": true, // non-aggregate expression inside PIVOT
}

// viewQueryFails fails a CREATE VIEW whose query oracle/catalog refuses with
// one of viewQueryCodes. The catalog is asked only when every relation the
// query reads is DUAL or a table of the schema whose columns the sync lists
// in full and the change has left alone, and it holds those tables alone.
// A table named with a partition is a heap table, and an analytic view a
// query names is a table, view or materialized view of the schema, which
// leaves no analytic view of that name.
func (b *builder) viewQueryFails(v *ast.CreateViewStmt, r objRef) bool {
	w := b.w
	if v.Force || v.IfNotExists || w.otherSpecs || w.current != w.snap.schema || b.isObjectView(v) {
		return false
	}
	if f, kind := b.look(r.key(spaceObject)); f == walk.Present && !(v.OrReplace && kind == kindView) {
		return false
	}
	// A name a WITH clause gives is a CTE where the clause is in scope. It
	// may name a table outside it too, so a CTE named as a table brings
	// that table along, and one named DUAL is not followed.
	ctes := map[string]bool{}
	known := true
	ast.Inspect(v.Query, func(n ast.Node) bool {
		if cte, ok := n.(*ast.CTE); ok {
			ctes[cte.Name] = true
			if cte.Name == "DUAL" {
				known = false
			}
		}
		return true
	})
	names := map[string]bool{}
	ast.Inspect(v.Query, func(n ast.Node) bool {
		if !known {
			return false
		}
		switch t := n.(type) {
		case *ast.CTE:
			if t.SubAV != nil && !(t.SubAV.Using != nil && t.SubAV.Using.Schema == "" && t.SubAV.Using.Name == t.Name) {
				known = b.namesNoAnalyticView(t.SubAV)
			}
		case *ast.InlineAnalyticView:
			known = b.namesNoAnalyticView(t.Clause)
		case *ast.TableRef:
			switch {
			case t.Name == nil || t.Dblink != "" || t.Name.DBLink != "":
				known = false
			case t.Name.Schema == "" && ctes[t.Name.Name] && w.snap.tables[t.Name.Name] == nil, b.isSysDual(t.Name):
			default:
				tr, ok := b.ref(t.Name)
				st := b.syncedColMeta(tr)
				if !ok || tr.schema != w.snap.schema || st == nil || !st.complete || t.PartitionExt != nil && (b.syncedPartMeta(tr) == nil || !st.heap) {
					known = false
				}
				names[tr.name] = true
			}
		}
		return known
	})
	if !known {
		return false
	}
	schema := &metadata.SchemaMetadata{}
	for _, s := range w.meta.GetSchemas() {
		if s.GetName() != "" && s.GetName() != w.snap.schema {
			continue
		}
		schema.Name = s.GetName()
		for _, t := range s.GetTables() {
			if names[t.GetName()] {
				schema.Tables = append(schema.Tables, t)
			}
		}
	}
	err := viewCatalogRefusal(&metadata.DatabaseSchemaMetadata{Name: w.meta.GetName(), Schemas: []*metadata.SchemaMetadata{schema}}, v, w.sql)
	var refusal *catalog.Error
	if !errors.As(err, &refusal) || !viewQueryCodes[refusal.Code] && !(viewGroupCodes[refusal.Code] && b.callsNoUserFunction(v.Query)) && !(refusal.Code == "00904" && b.namesNoRoutine(refusal.Message)) {
		return false
	}
	b.fail(r.rng, review.UnsupportedForm, "Oracle refuses the view's query: ORA-"+refusal.Code+": "+refusal.Message+".")
	return true
}

// namesNoRoutine reports whether a 00904 message names, with no
// qualifier, what is no function, procedure or package of the schema and
// no private or public synonym: Oracle reads such a name as a call of a
// function with no arguments, and the sync lists those in full, so a name
// none holds is no identifier at all.
func (b *builder) namesNoRoutine(message string) bool {
	name, ok := strings.CutSuffix(message, `": invalid identifier`)
	if !ok || !strings.HasPrefix(name, `"`) || strings.Contains(name[1:], `"`) {
		return false
	}
	r, ok := b.ref(&ast.ObjectName{Name: name[1:]})
	if !ok || r.schema != b.w.snap.schema {
		return false
	}
	for _, k := range []walk.Key{r.key(spaceRoutine), r.key(spaceSynonym), publicKey(r.name)} {
		if f, _ := b.look(k); f != walk.Absent {
			return false
		}
	}
	return true
}

// namesNoAnalyticView reports whether a sub_av_clause's USING names, with
// no schema, a table, view or materialized view of the schema: Oracle
// resolves the name to that object, which is no analytic view (18307 on a
// table, view or materialized view, 02201 on a sequence, verified).
func (b *builder) namesNoAnalyticView(c *ast.SubAVClause) bool {
	if c == nil || c.Using == nil || c.Using.Schema != "" {
		return false
	}
	return b.namesRelation(c.Using.Name)
}

// namesRelation reports whether name, with no schema, is a table, view or
// materialized view of the synced schema.
func (b *builder) namesRelation(name string) bool {
	r, ok := b.ref(&ast.ObjectName{Name: name})
	if !ok || r.schema != b.w.snap.schema {
		return false
	}
	f, kind := b.look(r.key(spaceObject))
	return f == walk.Present && (kind == kindTable || kind == kindView || kind == kindMView)
}

// callsNoUserFunction reports whether every function q calls is one the
// catalog knows as built in, or one whose name, with no schema, is a
// table, view or materialized view of the schema, which Oracle resolves
// to that object and so to no function. q declares no PL/SQL function of
// its own in a WITH clause.
func (b *builder) callsNoUserFunction(q ast.Node) bool {
	ok := true
	ast.Inspect(q, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.WithClause:
			if x.PLSQLDeclarations != nil && len(x.PLSQLDeclarations.Items) > 0 {
				ok = false
			}
		case *ast.FuncCallExpr:
			if x.FuncName == nil || x.FuncName.Schema != "" || x.FuncName.DBLink != "" || !catalog.BuiltinFunction(x.FuncName.Name) && !b.namesRelation(x.FuncName.Name) {
				ok = false
			}
		}
		return ok
	})
	return ok
}

// viewCatalogRefusal checks the CREATE VIEW in a catalog that holds meta.
// A panic in the catalog is no refusal: the walk then accepts the view.
func viewCatalogRefusal(meta *metadata.DatabaseSchemaMetadata, v *ast.CreateViewStmt, sql string) (err error) {
	defer func() {
		if recover() != nil {
			err = nil
		}
	}()
	c := catalog.New()
	c.SetServerConfig(catalog.StandardOracle23ServerConfig())
	if _, err := c.LoadMetadata(context.Background(), meta); err != nil {
		return nil
	}
	return c.CheckView(v, sql)
}
