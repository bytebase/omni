package catalog

import (
	"fmt"
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
)

// compileError is the first error Oracle lists for a PL/SQL unit it cannot
// compile. Oracle creates the unit, INVALID, and reports ORA-24344 ("success
// with compilation error") as a warning rather than a failure, so Exec
// stores the unit and the statement's ExecResult carries no error.
type compileError struct {
	// Kind and Name are the unit's: PROCEDURE, FUNCTION, PACKAGE,
	// PACKAGE BODY, TRIGGER, TYPE or TYPE BODY.
	Kind, Name string
	// Code is PLS-nnnnn, or ORA-nnnnn for a SQL statement in the unit.
	Code    string
	Message string
}

func (e *compileError) Error() string {
	return fmt.Sprintf("%s %s created with compilation errors: %s: %s", e.Kind, e.Name, e.Code, e.Message)
}

// unitState is one PL/SQL unit the catalog holds: its kind, name and
// status, and for a trigger ENABLED or DISABLED.
type unitState struct {
	kind, name, status, enabled string
	// err is the code of the first error of a unit that did not compile.
	err string
}

// unitStates lists the PL/SQL units the catalog holds.
func (c *Catalog) unitStates() []unitState {
	var out []unitState
	add := func(kind, name, enabled string) {
		u := unitState{kind: kind, name: name, status: "VALID", enabled: enabled}
		if c.units != nil && c.units.invalid[unitKey(kind, name)] {
			u.status = "INVALID"
			if ce := c.units.errors[unitKey(kind, name)]; ce != nil {
				u.err = ce.Code
			}
		}
		out = append(out, u)
	}
	if c.units != nil {
		for name, t := range c.units.types {
			if t.spec {
				add("TYPE", name, "")
			}
			if t.body {
				add("TYPE BODY", name, "")
			}
		}
	}
	for _, s := range c.meta.GetSchemas() {
		for _, p := range s.GetProcedures() {
			add("PROCEDURE", p.GetName(), "")
		}
		for _, f := range s.GetFunctions() {
			add("FUNCTION", f.GetName(), "")
		}
		for _, p := range s.GetPackages() {
			spec, body := splitPackageUnits(p.GetDefinition())
			if spec != "" {
				add("PACKAGE", p.GetName(), "")
			}
			if body != "" {
				add("PACKAGE BODY", p.GetName(), "")
			}
		}
		triggers := func(list []*metadatapb.TriggerMetadata) {
			for _, tr := range list {
				enabled := "ENABLED"
				if strings.Contains(tr.GetBody(), "\" DISABLE;") {
					enabled = "DISABLED"
				}
				add("TRIGGER", tr.GetName(), enabled)
			}
		}
		for _, t := range s.GetTables() {
			triggers(t.GetTriggers())
		}
		for _, v := range s.GetViews() {
			triggers(v.GetTriggers())
		}
		for _, v := range s.GetMaterializedViews() {
			triggers(v.GetTriggers())
		}
	}
	return out
}

// unitBook is what the catalog keeps of its PL/SQL units beyond their text,
// none of which the metadata holds: the units Oracle would list INVALID,
// the objects each unit reads, and the object and collection types.
type unitBook struct {
	// invalid holds unitKey(kind, name) of each invalid unit.
	invalid map[string]bool
	// errors holds, by unitKey, the first error of a unit that did not
	// compile.
	errors map[string]*compileError
	// reads holds, by unitKey, the names of the schema objects the unit
	// reads.
	reads map[string]map[string]bool
	// types holds the object and collection types by name.
	types map[string]*objectType
	// typesUnknown marks a catalog loaded from metadata, which holds no
	// types: a type name it does not hold may be one the schema has.
	typesUnknown bool
}

// objectType is a CREATE TYPE the catalog holds: its spec, its body, or
// both, as Oracle creates a body without a spec, INVALID.
type objectType struct {
	spec    bool
	final   bool
	body    bool
	methods []string
	attrs   []string
}

func newUnitBook() *unitBook {
	return &unitBook{invalid: map[string]bool{}, errors: map[string]*compileError{}, reads: map[string]map[string]bool{}, types: map[string]*objectType{}}
}

func (b *unitBook) clone() *unitBook {
	if b == nil {
		return newUnitBook()
	}
	out := newUnitBook()
	out.typesUnknown = b.typesUnknown
	for k := range b.invalid {
		out.invalid[k] = true
	}
	for k, ce := range b.errors {
		out.errors[k] = ce
	}
	for k, names := range b.reads {
		m := make(map[string]bool, len(names))
		for n := range names {
			m[n] = true
		}
		out.reads[k] = m
	}
	for k, t := range b.types {
		cp := *t
		cp.methods = append([]string(nil), t.methods...)
		cp.attrs = append([]string(nil), t.attrs...)
		out.types[k] = &cp
	}
	return out
}

func unitKey(kind, name string) string { return kind + "\x00" + name }

// Unit kinds, as USER_OBJECTS names them.
const (
	unitProcedure   = "PROCEDURE"
	unitFunction    = "FUNCTION"
	unitPackage     = "PACKAGE"
	unitPackageBody = "PACKAGE BODY"
	unitTrigger     = "TRIGGER"
	unitType        = "TYPE"
	unitTypeBody    = "TYPE BODY"
)

// hasUnit reports whether the schema holds a unit of kind under name.
func (e *oracleMetadataExtractor) hasUnit(kind, name string) bool {
	switch kind {
	case unitProcedure:
		return e.procedures[name] != nil
	case unitFunction:
		return e.functions[name] != nil
	case unitPackage, unitPackageBody:
		pkg := e.packages[name]
		if pkg == nil {
			return false
		}
		spec, body := splitPackageUnits(pkg.Definition)
		if kind == unitPackage {
			return spec != ""
		}
		return body != ""
	case unitType:
		return e.hasType(name)
	case unitTypeBody:
		return e.units.types[name] != nil && e.units.types[name].body
	case unitTrigger:
		return e.triggers[name] != nil
	}
	return false
}

// createsUnit checks the name a CREATE of a unit of kind takes, and
// reports whether the statement goes on to create it. A spec's name is in
// the schema's one namespace with tables and views; a body's is its spec's,
// and Oracle creates a body without a spec, INVALID.
func (e *oracleMetadataExtractor) createsUnit(kind string, name *ast.ObjectName, orReplace, ifNotExists bool) bool {
	if name == nil || name.Name == "" {
		return false
	}
	if e.isForeignSchema(name.Schema) {
		// Creating in another schema takes CREATE ANY privileges.
		e.setError(oraError("01031", "insufficient privileges"))
		return false
	}
	if !e.useSchema(name.Schema) {
		return false
	}
	n := name.Name
	switch kind {
	case unitPackageBody, unitTypeBody:
		if e.hasUnit(kind, n) && !orReplace {
			if ifNotExists {
				return false
			}
			e.setError(oraError("00955", "name %q is already used by an existing object", n))
			return false
		}
		return true
	}
	if e.hasUnit(kind, n) {
		if ifNotExists {
			return false
		}
		if !orReplace {
			e.setError(oraError("00955", "name %q is already used by an existing object", n))
			return false
		}
		return true
	}
	bodyOnly := kind == unitPackage && e.packages[n] != nil || kind == unitType && e.units.types[n] != nil
	if e.objectExistsInNamespace(n) && !bodyOnly {
		if ifNotExists {
			return false
		}
		e.setError(oraError("00955", "name %q is already used by an existing object", n))
		return false
	}
	return true
}

// extractCreateType stores an object or collection type, or its body,
// whose name createsUnit has checked.
func (e *oracleMetadataExtractor) extractCreateType(n *ast.CreateTypeStmt) {
	name := n.Name.Name
	old := e.units.types[name]
	if n.IsBody {
		if old == nil {
			old = &objectType{}
			e.units.types[name] = old
		}
		old.body = true
		return
	}
	t := &objectType{spec: true, final: true}
	if old != nil {
		t.body = old.body
	}
	for _, m := range n.Modifiers {
		switch strings.ToUpper(m) {
		case "NOT FINAL":
			t.final = false
		case "FINAL":
			t.final = true
		}
	}
	for _, item := range listItems(n.Attributes) {
		if col, ok := item.(*ast.ColumnDef); ok {
			t.attrs = append(t.attrs, col.Name)
		}
	}
	for _, item := range listItems(n.Methods) {
		if m, ok := item.(*ast.TypeBodyMember); ok {
			t.methods = append(t.methods, subprogramName(m.Subprog))
		}
	}
	if old != nil && old.spec && t.body && (!slices.Equal(old.attrs, t.attrs) || !slices.Equal(old.methods, t.methods) || old.final != t.final) {
		// A changed spec leaves its body to be compiled again; Oracle
		// takes the same spec again as no change.
		e.units.invalid[unitKey(unitTypeBody, name)] = true
	}
	e.units.types[name] = t
}

// hasType reports an object or collection type's spec under name.
func (e *oracleMetadataExtractor) hasType(name string) bool {
	t := e.units.types[name]
	return t != nil && t.spec
}

// subprogramName is the name of a *CreateProcedureStmt or
// *CreateFunctionStmt.
func subprogramName(n ast.Node) string {
	switch s := n.(type) {
	case *ast.CreateProcedureStmt:
		return objectName(s.Name)
	case *ast.CreateFunctionStmt:
		return objectName(s.Name)
	}
	return ""
}

// findsUnit checks the unit a DROP or an ALTER names, and reports whether
// it exists. A missing one is ORA-04043, a missing trigger ORA-04080. A
// catalog loaded from metadata holds no types, so it finds any type.
func (e *oracleMetadataExtractor) findsUnit(kind string, name *ast.ObjectName, ifExists bool) bool {
	if name == nil || !e.useSchema(name.Schema) {
		return false
	}
	if e.hasUnit(kind, name.Name) {
		return true
	}
	if (kind == unitType || kind == unitTypeBody) && e.units.typesUnknown {
		return true
	}
	if !ifExists {
		if kind == unitTrigger {
			e.setError(oraError("04080", "trigger %q does not exist", name.Name))
		} else {
			e.setError(oraError("04043", "object %q does not exist", name.Name))
		}
	}
	return false
}

// forgetUnit removes what the book holds of a dropped unit.
func (e *oracleMetadataExtractor) forgetUnit(kind, name string) {
	key := unitKey(kind, name)
	delete(e.units.invalid, key)
	delete(e.units.errors, key)
	delete(e.units.reads, key)
}

// checkTrigger checks what Oracle checks of a DML trigger before it
// compiles the body, and reports whether the statement goes on.
func (e *oracleMetadataExtractor) checkTrigger(n *ast.CreateTriggerStmt) bool {
	name := n.Name.Name
	table := e.resolveSynonymTarget(n.Table.Name, n.Table.Schema)
	if old := e.triggers[name]; old != nil {
		if n.IfNotExists {
			return false
		}
		if !n.OrReplace {
			e.setError(oraError("04081", "trigger %q already exists", name))
			return false
		}
		if on := e.triggerTable(name); on != "" && on != table {
			e.setError(oraError("04095", "trigger %q already exists on another table, cannot replace it", name))
			return false
		}
	}
	view := e.views[table] != nil
	if !view && e.tables[table] == nil && e.materializedViews[table] == nil {
		e.setError(oraError("00942", "table or view %q does not exist", n.Table.Name))
		return false
	}
	switch {
	case n.Timing == ast.TRIGGER_INSTEAD_OF && !view:
		e.setError(oraError("25002", "cannot create INSTEAD OF triggers on tables"))
		return false
	case n.Timing != ast.TRIGGER_INSTEAD_OF && view:
		e.setError(oraError("25001", "cannot create this trigger type on views"))
		return false
	}
	if t := e.tables[table]; t != nil {
		for _, col := range n.UpdateColumns {
			if !tableHasColumn(t, col) {
				e.setError(oraError("00904", "%q: invalid identifier", col))
				return false
			}
		}
	}
	if !n.ForEachRow && !n.Compound {
		if n.When != nil {
			e.setError(oraError("04077", "WHEN clause cannot be used with table level triggers"))
			return false
		}
		if triggerReadsRow(n) {
			e.setError(oraError("04082", "NEW or OLD references not allowed in table level triggers"))
			return false
		}
	}
	if n.ForEachRow && (n.Timing == ast.TRIGGER_AFTER || onlyDelete(n)) && triggerSetsNew(n) {
		e.setError(oraError("04084", "cannot change NEW values for this trigger type"))
		return false
	}
	return true
}

// triggerTable is the table or view a trigger is on.
func (e *oracleMetadataExtractor) triggerTable(name string) string {
	for tname, t := range e.tables {
		for _, tr := range t.Triggers {
			if tr.Name == name {
				return tname
			}
		}
	}
	for vname, v := range e.views {
		for _, tr := range v.Triggers {
			if tr.Name == name {
				return vname
			}
		}
	}
	return ""
}

func tableHasColumn(t *metadatapb.TableMetadata, name string) bool {
	for _, c := range t.Columns {
		if c.Name == name {
			return true
		}
	}
	return false
}

func onlyDelete(n *ast.CreateTriggerStmt) bool {
	items := listItems(n.Events)
	for _, item := range items {
		if ev, ok := item.(*ast.Integer); !ok || ast.TriggerEvent(ev.Ival) != ast.TRIGGER_DELETE {
			return false
		}
	}
	return len(items) > 0
}

// rowNames are the names a trigger's body reads its row by: NEW and OLD,
// or what REFERENCING calls them.
func rowNames(n *ast.CreateTriggerStmt) map[string]bool {
	names := map[string]bool{"NEW": true, "OLD": true}
	if r := n.Referencing; r != nil {
		for _, a := range []string{r.NewAlias, r.OldAlias} {
			if a != "" {
				names[strings.ToUpper(a)] = true
			}
		}
	}
	return names
}

// triggerReadsRow reports whether a trigger's body names :NEW or :OLD.
func triggerReadsRow(n *ast.CreateTriggerStmt) bool {
	names := rowNames(n)
	found := false
	ast.Inspect(n.Body, func(x ast.Node) bool {
		if b, ok := x.(*ast.BindVariable); ok && names[strings.ToUpper(b.Name)] {
			found = true
		}
		return !found
	})
	return found
}

// triggerSetsNew reports whether a trigger's body assigns to a :NEW field.
func triggerSetsNew(n *ast.CreateTriggerStmt) bool {
	newName := "NEW"
	if r := n.Referencing; r != nil && r.NewAlias != "" {
		newName = strings.ToUpper(r.NewAlias)
	}
	found := false
	ast.Inspect(n.Body, func(x ast.Node) bool {
		if a, ok := x.(*ast.PLSQLAssign); ok {
			if b, ok := a.Target.(*ast.BindVariable); ok && strings.ToUpper(b.Name) == newName {
				found = true
			}
		}
		return !found
	})
	return found
}

// compileUnit compiles the unit a statement created as Oracle does when it
// creates it: the unit is INVALID, with its first error, when a name in it
// does not resolve, and VALID otherwise. It records what the unit reads.
func (e *oracleMetadataExtractor) compileUnit(kind, name string, n ast.Node) {
	key := unitKey(kind, name)
	delete(e.units.invalid, key)
	delete(e.units.errors, key)
	c := newPLChecker(e, kind, name)
	c.unit(n)
	e.units.reads[key] = c.reads
	if c.err != nil {
		e.units.invalid[key] = true
		e.units.errors[key] = c.err
	}
}

// invalidateReaders marks INVALID every unit that reads name, and every
// unit that reads one of those, as Oracle does when an object changes.
func (e *oracleMetadataExtractor) invalidateReaders(name string) {
	seen := map[string]bool{}
	queue := []string{name}
	for len(queue) > 0 {
		obj := queue[0]
		queue = queue[1:]
		if seen[obj] {
			continue
		}
		seen[obj] = true
		for key, reads := range e.units.reads {
			if !reads[obj] || e.units.invalid[key] {
				continue
			}
			e.units.invalid[key] = true
			// A spec's readers read the spec, not its body; nothing reads
			// a trigger.
			kind, unit, _ := strings.Cut(key, "\x00")
			if kind != unitPackageBody && kind != unitTypeBody && kind != unitTrigger {
				queue = append(queue, unit)
			}
		}
	}
}

// sameSignature reports whether a stored function's text declares the
// parameters and return type sig renders.
func sameSignature(definition, sig string) bool {
	text := strings.TrimSpace(definition)
	if !strings.HasPrefix(strings.ToUpper(text), "CREATE") {
		text = "CREATE " + text
	}
	f, ok := parseSaved(text).(*ast.CreateFunctionStmt)
	return ok && signatureText(f.Parameters, f.ReturnType) == sig
}

// invalidateColumnReaders marks INVALID the units that read a column an
// ALTER TABLE drops, renames or modifies, or read all of the table's
// columns (%ROWTYPE, SELECT *, INSERT without a column list). A unit that
// reads only other columns stays VALID, as Oracle tracks what it reads.
func (e *oracleMetadataExtractor) invalidateColumnReaders(n *ast.AlterTableStmt) {
	if n.Name == nil {
		return
	}
	table := n.Name.Name
	var columns []string
	for _, item := range listItems(n.Actions) {
		action, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch action.Action {
		case ast.AT_DROP_COLUMN, ast.AT_SET_UNUSED, ast.AT_RENAME_COLUMN:
			if action.Options != nil && len(action.Options.Items) > 0 {
				columns = append(columns, stringList(action.Options)...)
			} else {
				columns = append(columns, action.ColumnName)
			}
		case ast.AT_MODIFY_COLUMN:
			if action.ColumnDef != nil {
				columns = append(columns, action.ColumnDef.Name)
			}
			for _, item := range listItems(action.ColumnDefs) {
				if col, ok := item.(*ast.ColumnDef); ok {
					columns = append(columns, col.Name)
				}
			}
		}
	}
	if len(columns) == 0 {
		return
	}
	e.invalidateReaders(table + ".*")
	for _, col := range columns {
		e.invalidateReaders(table + "." + col)
	}
}
