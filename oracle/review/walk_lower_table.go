package review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// alterTable follows ALTER TABLE, which follows no synonym: the table must
// exist (00942), as a table or a materialized view's container. The
// statement is atomic, and its clauses see the table as it was: every
// check holds against the state before it (walk.Lookups), and the repeats
// of one statement are text claims (spec §4). A clause the walk does not
// follow ends the walk before the statement.
func (b *builder) alterTable(v *ast.AlterTableStmt) bool {
	r, ok := b.ref(v.Name)
	if !ok {
		return false
	}
	var cmds []*ast.AlterTableCmd
	if v.Actions != nil {
		for _, item := range v.Actions.Items {
			c, ok := item.(*ast.AlterTableCmd)
			if !ok || !alterFollowed(c) {
				return false
			}
			cmds = append(cmds, c)
		}
	}
	if b.refTooLong(v.Name, r) {
		return true
	}
	b.Order = walk.Lookups
	k := r.key(spaceTable)
	f, kind := b.look(k)
	b.need(k, walk.Present, nil, r.rng, missingTable(r.written), "")
	if f == walk.Absent {
		return true
	}
	if b.alterTextFails(cmds, r) {
		return true
	}
	for _, c := range cmds {
		b.Next()
		if c.Action == ast.AT_RENAME && kind == kindMView {
			return false
		}
		if !b.alterCmd(r, c) {
			return false
		}
	}
	return true
}

// alterProperties are the table properties ALTER TABLE may change, by the
// parser's name, which change no name the walk keeps. ADD SUPPLEMENTAL LOG
// adds a log group, which is a constraint, and is not here.
var alterProperties = map[string]bool{
	"LOGGING": true, "NOLOGGING": true, "FILESYSTEM_LIKE_LOGGING": true, "READ ONLY": true, "READ WRITE": true,
	"PARALLEL": true, "NOPARALLEL": true, "COMPRESS": true, "NOCOMPRESS": true, "ROW STORE COMPRESS": true,
	"COLUMN STORE COMPRESS": true, "CACHE": true, "NOCACHE": true, "RESULT_CACHE": true, "PCTFREE": true,
	"PCTUSED": true, "INITRANS": true, "MAXTRANS": true, "STORAGE": true, "ENABLE ROW MOVEMENT": true,
	"DISABLE ROW MOVEMENT": true, "DEALLOCATE UNUSED": true, "ALLOCATE EXTENT": true, "FLASHBACK ARCHIVE": true,
	"NO FLASHBACK ARCHIVE": true, "MONITORING": true, "NOMONITORING": true, "INMEMORY": true, "NO INMEMORY": true,
	"ROWDEPENDENCIES": true, "NOROWDEPENDENCIES": true, "INDEXING ON": true, "INDEXING OFF": true,
	"SEGMENT CREATION IMMEDIATE": true, "SEGMENT CREATION DEFERRED": true, "MEMOPTIMIZE FOR READ": true,
	"NO MEMOPTIMIZE FOR READ": true, "MEMOPTIMIZE FOR WRITE": true, "NO MEMOPTIMIZE FOR WRITE": true,
	"ILM": true, "UPGRADE": true, "DEFAULT COLLATION": true, "CLUSTERING": true, "NO CLUSTERING": true,
}

// alterFollowed reports whether the walk follows a clause of ALTER TABLE.
func alterFollowed(c *ast.AlterTableCmd) bool {
	switch c.Action {
	case ast.AT_ALTER_PROPERTY:
		return alterProperties[c.Subtype]
	case ast.AT_ADD_PERIOD, ast.AT_DROP_PERIOD:
		// A period adds hidden columns and a constraint.
		return false
	}
	return true
}

// alterTextFails fails an ALTER TABLE its text alone refutes (spec §4,
// verified): a column twice across ADD, MODIFY, and a DROP list (00957), a
// constraint name twice (02264), two primary keys (02260), the same
// constraint dropped twice (02443), or a name longer than Oracle allows
// (00972).
func (b *builder) alterTextFails(cmds []*ast.AlterTableCmd, r objRef) bool {
	columns := make(map[string]bool)
	twice := func(name string, rng review.Range) bool {
		if columns[name] {
			b.fail(rng, "Column "+ident(name)+" is named twice in the statement")
			return true
		}
		columns[name] = true
		return false
	}
	var d tableDecl
	dropped := make(map[string]bool)
	for _, c := range cmds {
		switch c.Action {
		case ast.AT_ADD_COLUMN, ast.AT_MODIFY_COLUMN:
			for _, col := range alterColumns(c) {
				if b.tooLong(col.Name, rangeOf(col.Loc)) || twice(col.Name, rangeOf(col.Loc)) {
					return true
				}
				b.declareColumn(&d, col)
			}
		case ast.AT_DROP_COLUMN:
			for _, name := range b.droppedColumns(c) {
				if twice(name, rangeOf(c.Loc)) {
					return true
				}
			}
		case ast.AT_ADD_CONSTRAINT:
			if c.Constraint != nil {
				b.declareConstraint(&d, c.Constraint)
			}
		case ast.AT_DROP_CONSTRAINT:
			if c.Constraint == nil || c.Constraint.Name == "" || c.Subtype != "" {
				continue
			}
			if dropped[c.Constraint.Name] {
				b.fail(rangeOf(c.Loc), "Cannot drop constraint "+ident(c.Constraint.Name)+": table "+r.written+" has no such constraint")
				return true
			}
			dropped[c.Constraint.Name] = true
		case ast.AT_RENAME_COLUMN, ast.AT_RENAME_CONSTRAINT, ast.AT_RENAME:
			if b.tooLong(c.NewName, rangeOf(c.Loc)) {
				return true
			}
		}
	}
	return b.namesTooLong(nil, &d) || b.declTextFails(&d, r)
}

// alterColumns lists the column definitions of an ADD or MODIFY.
func alterColumns(c *ast.AlterTableCmd) []*ast.ColumnDef {
	cols := columnDefs(c.ColumnDefs)
	if c.ColumnDef != nil {
		cols = append([]*ast.ColumnDef{c.ColumnDef}, cols...)
	}
	return cols
}

// droppedColumns lists the columns of DROP COLUMN c, DROP (c, ...), SET
// UNUSED COLUMN c, or SET UNUSED (c, ...). The parser keeps only the first
// of a list, so the list is read from the tokens.
func (b *builder) droppedColumns(c *ast.AlterTableCmd) []string {
	toks := tokens(b.w.sql, c.Loc)
	for i, tok := range toks {
		if tok.Type != tokenOpenParen {
			continue
		}
		var out []string
		for _, t := range toks[i+1:] {
			switch t.Type {
			case tokenCloseParen:
				return out
			case tokenComma:
			default:
				name, ok := tokenName(t)
				if !ok {
					return nil
				}
				out = append(out, name)
			}
		}
		return nil
	}
	if c.ColumnName == "" {
		return nil
	}
	return []string{c.ColumnName}
}

// cascadesConstraints reports whether a clause ends with CASCADE
// CONSTRAINTS, or CASCADE for a constraint.
func (b *builder) cascadesConstraints(c *ast.AlterTableCmd) bool {
	return cascades(b.w.sql, c.Loc)
}

func (b *builder) alterCmd(r objRef, c *ast.AlterTableCmd) bool {
	rng := rangeOf(c.Loc)
	switch c.Action {
	case ast.AT_ADD_COLUMN:
		var d tableDecl
		for _, col := range alterColumns(c) {
			ck := r.childKey(spaceColumn, col.Name)
			b.need(ck, walk.Absent, nil, rangeOf(col.Loc), "", "Column "+ident(col.Name)+" being added already exists in table "+r.written)
			b.Set(ck, walk.Present, "")
			b.declareColumn(&d, col)
		}
		b.declChecks(&d, r)
		b.declEffects(&d, r)
	case ast.AT_MODIFY_COLUMN:
		var d tableDecl
		for _, col := range alterColumns(c) {
			b.columnNeeds(r, col.Name, rangeOf(col.Loc))
			b.declareColumn(&d, col)
			if col.Null {
				// NULL drops the column's NOT NULL constraint, of any name.
				b.dropsConstraints(r)
			}
		}
		b.declChecks(&d, r)
		b.declEffects(&d, r)
	case ast.AT_DROP_COLUMN, ast.AT_SET_UNUSED:
		names := b.droppedColumns(c)
		if names == nil {
			return false
		}
		for _, name := range names {
			b.columnNeeds(r, name, rng)
			b.Set(r.childKey(spaceColumn, name), walk.Absent, "")
		}
		// The column's constraints and indexes go with it.
		b.dropsConstraints(r)
		b.MayRemove(walk.Scope{Space: spacePrimary, Parent: r.path()})
		b.dropIndexes(r)
		if b.cascadesConstraints(c) {
			b.cascadeConstraints()
		}
	case ast.AT_RENAME_COLUMN:
		from, to := r.childKey(spaceColumn, c.ColumnName), r.childKey(spaceColumn, c.NewName)
		b.columnNeeds(r, c.ColumnName, rng)
		b.need(to, walk.Absent, nil, rng, "", "Column "+ident(c.NewName)+" already exists in table "+r.written)
		if c.ColumnName != c.NewName {
			b.Set(from, walk.Absent, "")
			b.Set(to, walk.Present, "")
		}
	case ast.AT_ADD_CONSTRAINT:
		if c.Constraint == nil {
			return false
		}
		var d tableDecl
		b.declareConstraint(&d, c.Constraint)
		b.declChecks(&d, r)
		b.declEffects(&d, r)
	case ast.AT_DROP_CONSTRAINT:
		return b.dropConstraint(r, c, rng)
	case ast.AT_MODIFY_CONSTRAINT, ast.AT_ENABLE_DISABLE:
		return b.constraintState(r, c, rng)
	case ast.AT_RENAME_CONSTRAINT:
		from, to := c.ColumnName, c.NewName
		b.constraintNeeds(r, from, rng, "Cannot rename constraint "+ident(from)+": table "+r.written+" has no such constraint")
		b.need(r.schemaKey(spaceConstraint, to), walk.Absent, nil, rng, "", takenConstraint(to))
		if from != to {
			b.Set(r.childKey(spaceTableConstraint, from), walk.Absent, "")
			b.Set(r.schemaKey(spaceConstraint, from), walk.Absent, "")
			b.Set(r.childKey(spaceTableConstraint, to), walk.Present, "")
			b.Set(r.schemaKey(spaceConstraint, to), walk.Present, r.path())
		}
	case ast.AT_RENAME:
		if strings.Contains(c.NewName, ".") {
			return false
		}
		to := refIn(r, c.NewName)
		b.need(to.key(spaceObject), walk.Absent, nil, rng, "", takenName(to.written))
		if to.name != r.name {
			b.moveTable(r, to)
		}
	case ast.AT_DROP_UNUSED_COLUMNS, ast.AT_ENABLE_DISABLE_TRIGGERS, ast.AT_ENABLE_DISABLE_TABLE_LOCK, ast.AT_MOVE,
		ast.AT_SHRINK_SPACE, ast.AT_ALTER_PROPERTY, ast.AT_MODIFY_LOB, ast.AT_MODIFY_VARRAY, ast.AT_MODIFY_NESTED_TABLE,
		ast.AT_MODIFY_OPAQUE_TYPE, ast.AT_IMMUTABLE_TABLE, ast.AT_BLOCKCHAIN_TABLE, ast.AT_DUPLICATED_REFRESH,
		ast.AT_ADD_PARTITION, ast.AT_DROP_PARTITION, ast.AT_TRUNCATE_PARTITION, ast.AT_SPLIT_PARTITION,
		ast.AT_SPLIT_SUBPARTITION, ast.AT_MERGE_PARTITIONS, ast.AT_MERGE_SUBPARTITIONS, ast.AT_EXCHANGE_PARTITION,
		ast.AT_COALESCE_PARTITION, ast.AT_COALESCE_SUBPARTITION, ast.AT_MODIFY_PARTITION, ast.AT_MODIFY_SUBPARTITION,
		ast.AT_MODIFY_DEFAULT_ATTRS, ast.AT_DROP_SUBPARTITION, ast.AT_RENAME_PARTITION, ast.AT_SET_INTERVAL,
		ast.AT_SET_PARTITIONING, ast.AT_SET_SUBPARTITION_TEMPLATE, ast.AT_MODIFY_TO_PARTITIONED, ast.AT_MODIFY_PARTITIONSET:
		// Storage, partitions, and properties: no name the walk keeps.
	default:
		return false
	}
	return true
}

// constraintNeeds checks that a constraint the clause names is the
// table's: on the table, and, since constraint names are unique in the
// schema, not another table's (02443, 23292, 02430, 02431).
func (b *builder) constraintNeeds(r objRef, name string, rng review.Range, message string) {
	b.need(r.childKey(spaceTableConstraint, name), walk.Present, nil, rng, message, "")
	b.need(r.schemaKey(spaceConstraint, name), walk.Present, []string{r.path()}, rng, message, message)
}

// columnNeeds checks that a column the clause names exists (00904).
func (b *builder) columnNeeds(r objRef, name string, rng review.Range) {
	b.need(r.childKey(spaceColumn, name), walk.Present, nil, rng, "Column "+ident(name)+" does not exist in table "+r.written, "")
}

// dropsConstraints records that a clause may drop constraints of the
// table, of names the walk does not know.
func (b *builder) dropsConstraints(r objRef) {
	b.MayRemove(walk.Scope{Space: spaceTableConstraint, Parent: r.path()})
	b.MayRemove(walk.Scope{Space: spaceConstraint, Parent: walk.Path(r.schema), Kinds: []string{r.path()}})
}

// dropConstraint follows DROP CONSTRAINT c, DROP PRIMARY KEY, and DROP
// UNIQUE (...) [CASCADE]: c must be the table's (02443, another table's
// too), and the table must have a primary key to drop (02441) (verified).
// The constraint's index goes with it, and with CASCADE the foreign keys
// that refer to it.
func (b *builder) dropConstraint(r objRef, c *ast.AlterTableCmd, rng review.Range) bool {
	switch {
	case c.Subtype == "PRIMARY KEY":
		pk := r.childKey(spacePrimary, "")
		b.need(pk, walk.Present, nil, rng, "Table "+r.written+" has no primary key to drop", "")
		b.Set(pk, walk.Absent, "")
		b.dropsConstraints(r)
	case c.Subtype == "UNIQUE":
		b.dropsConstraints(r)
	case c.Constraint != nil && c.Constraint.Name != "":
		name := c.Constraint.Name
		ck := r.childKey(spaceTableConstraint, name)
		b.constraintNeeds(r, name, rng, "Cannot drop constraint "+ident(name)+": table "+r.written+" has no such constraint")
		b.Set(ck, walk.Absent, "")
		b.Set(r.schemaKey(spaceConstraint, name), walk.Absent, "")
		// It may be the primary key.
		b.MayRemove(walk.Scope{Space: spacePrimary, Parent: r.path()})
	default:
		return false
	}
	b.dropIndexes(r)
	if b.cascadesConstraints(c) {
		b.cascadeConstraints()
	}
	return true
}

// constraintState follows ENABLE, DISABLE, and MODIFY of a constraint:
// a named one must be the table's (02430 to enable, 02431 to disable or
// modify, verified). DISABLE drops the index the constraint made, and
// ENABLE makes it again, under the constraint's name or another (spec §4).
func (b *builder) constraintState(r objRef, c *ast.AlterTableCmd, rng review.Range) bool {
	name := c.ColumnName
	if c.Action == ast.AT_MODIFY_CONSTRAINT && c.Constraint != nil {
		name = c.Constraint.Name
	}
	switch {
	case name != "":
		verb := "disable"
		if strings.HasPrefix(c.Subtype, "ENABLE") {
			verb = "enable"
		}
		b.constraintNeeds(r, name, rng, "Cannot "+verb+" constraint "+ident(name)+": table "+r.written+" has no such constraint")
		b.Set(r.schemaKey(spaceIndex, name), walk.Unknown, "")
		if b.createsIndex(c.Loc) {
			b.MayAdd(walk.Scope{Space: spaceIndex, Parent: walk.Path(r.schema)})
		}
	case strings.HasSuffix(c.Subtype, "PRIMARY KEY"), strings.HasSuffix(c.Subtype, "UNIQUE"):
		// The constraint's name, and so its index's, is not in the text.
		b.MayAdd(walk.Scope{Space: spaceIndex, Parent: walk.Path(r.schema)})
	default:
		return false
	}
	b.dropIndexes(r)
	return true
}
