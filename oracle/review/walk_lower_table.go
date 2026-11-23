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
	if hasModifyVarrayShrink(b.w.sql, v.Actions) {
		// Oracle 23 accepts ALTER TABLE ... MODIFY VARRAY ... (SHRINK SPACE [CASCADE])
		// even when the table does not exist or names a non-table object (verified).
		return false
	}
	var cmds []*ast.AlterTableCmd
	unfollowed := false
	if v.Actions != nil {
		for _, item := range v.Actions.Items {
			c, ok := item.(*ast.AlterTableCmd)
			if !ok || !alterFollowed(c) {
				unfollowed = true
				continue
			}
			cmds = append(cmds, c)
		}
	}
	if b.refTooLong(v.Name, r) {
		return true
	}
	b.Order = walk.Lookups
	k := r.key(spaceTable)
	if !unfollowed && alterFollowsSynonym(cmds) {
		// ALTER TABLE ... MOVE and ALTER TABLE ... MODIFY PARTITION BY ...
		// follow private and public synonyms in Oracle 23 (verified).
		if f, _ := b.look(k); f == walk.Absent {
			if objFact, objKind := b.look(r.key(spaceObject)); objFact == walk.Present {
				if objKind == kindSynonym {
					if sr, _, sok := b.resolveSynonymTable(v.Name, r); sok {
						for _, c := range cmds {
							if c.Action == ast.AT_MOVE && b.moveHasPartition(c.Loc) {
								if b.checkPartitionTextContradictions(c.Loc, sr, true, rangeOf(c.Loc)) {
									return true
								}
								if b.checkPartitionMeta(sr, b.syncedPartMeta(sr), c.Loc, rangeOf(c.Loc)) {
									return true
								}
							} else if c.Action == ast.AT_MODIFY_TO_PARTITIONED {
								if b.checkModifyPartitionBy(sr, c.Loc, rangeOf(c.Loc)) {
									return true
								}
							}
						}
					}
					b.Last = true
					return true
				}
				if objKind != "" && objKind != kindTable && objKind != kindMView {
					b.heldByOther(k, r, kindTable)
					return true
				}
			}
			if !generated(spaceTable, r.name) {
				keys := []walk.Key{k, r.key(spaceSynonym)}
				if v.Name.Schema == "" {
					keys = append(keys, publicKey(r.name))
				}
				allAbsent := true
				for _, sk := range keys {
					if sf, _ := b.look(sk); sf != walk.Absent {
						allAbsent = false
						break
					}
				}
				if allAbsent {
					b.Fail(keys, walk.Absent, noteOf(r.rng, review.ObjectNotFound, missingTable(r.written)))
					return true
				}
			}
			b.Last = true
			return true
		}
	} else {
		if b.heldByOther(k, r, kindTable) {
			return true
		}
		f, _ := b.look(k)
		b.need(k, walk.Present, nil, r.rng, missingTable(r.written), "")
		if f == walk.Absent {
			return true
		}
	}
	if b.checkDuplicateAnnotations(v.Loc, r.rng) {
		return true
	}
	if unfollowed {
		// Table existence has been verified, but an unfollowed clause may add
		// hidden columns or constraints, so stop walking after this statement.
		// Still check supplemental log group column lists (ORA-00904/00957) and
		// REF/SCOPE clauses (ORA-22893) on known tables before ending the walk.
		if v.Actions != nil {
			for _, item := range v.Actions.Items {
				c, ok := item.(*ast.AlterTableCmd)
				if !ok {
					continue
				}
				if c.Action == ast.AT_ALTER_PROPERTY && strings.HasPrefix(c.Subtype, "ADD SUPPLEMENTAL LOG") {
					if b.checkSupplementalLogGroup(r, c.Loc, rangeOf(c.Loc)) {
						return true
					}
				}
			}
		}
		b.Last = true
		return true
	}
	_, kind := b.look(k)
	declared, fails := b.alterTextFails(cmds, r)
	if fails {
		return true
	}
	addedCols := make(map[string]bool)
	addedPK := false
	for _, c := range cmds {
		switch c.Action {
		case ast.AT_ADD_COLUMN:
			for _, col := range alterColumns(c) {
				addedCols[col.Name] = true
				if pk, _ := columnKeys(col); pk {
					addedPK = true
				}
			}
		case ast.AT_MODIFY_COLUMN:
			for _, col := range alterColumns(c) {
				if pk, _ := columnKeys(col); pk {
					addedPK = true
				}
			}
		case ast.AT_ADD_CONSTRAINT:
			if c.Constraint != nil && c.Constraint.Type == ast.CONSTRAINT_PRIMARY {
				addedPK = true
			}
		}
	}
	b.Last = b.domain(r)
	for _, c := range cmds {
		b.Next()
		if c.Action == ast.AT_RENAME && kind == kindMView {
			return false
		}
		if !b.alterCmd(r, kind, c, declared, addedCols, addedPK) {
			return false
		}
	}
	return true
}

// hasModifyVarrayShrink reports whether any action is MODIFY VARRAY ... (SHRINK SPACE ...),
// which Oracle 23 accepts as a no-op even when the table is missing or another kind.
func hasModifyVarrayShrink(sql string, actions *ast.List) bool {
	if actions == nil {
		return false
	}
	for _, item := range actions.Items {
		c, ok := item.(*ast.AlterTableCmd)
		if !ok || c.Action != ast.AT_MODIFY_VARRAY {
			continue
		}
		toks := tokens(sql, c.Loc)
		inParen := false
		for i, tok := range toks {
			if tok.Type == tokenOpenParen {
				inParen = true
				continue
			}
			if inParen && tok.Type != tokenQuoted && strings.EqualFold(tok.Str, "SHRINK") &&
				i+1 < len(toks) && toks[i+1].Type != tokenQuoted && strings.EqualFold(toks[i+1].Str, "SPACE") {
				return true
			}
		}
	}
	return false
}

// alterFollowsSynonym reports whether every clause of an ALTER TABLE is one
// that follows synonyms in Oracle 23 (MOVE and MODIFY PARTITION BY).
func alterFollowsSynonym(cmds []*ast.AlterTableCmd) bool {
	if len(cmds) == 0 {
		return false
	}
	for _, c := range cmds {
		if c.Action != ast.AT_MOVE && c.Action != ast.AT_MODIFY_TO_PARTITIONED {
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
	"MODIFY COLUMN SUBSTITUTABLE": true,
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
// (00972). It returns the constraint names the statement's ADD and MODIFY
// clauses declare.
func (b *builder) alterTextFails(cmds []*ast.AlterTableCmd, r objRef) (declared []string, fails bool) {
	// columns holds what the statement does to each column it names:
	// "adds", "modifies", or "drops".
	columns := make(map[string]string)
	named := func(name, verb string, rng review.Range) bool {
		if first, ok := columns[name]; ok {
			if first != verb {
				verb = "names"
			}
			b.fail(rng, review.DuplicateName, twice(verb, "column", name))
			return true
		}
		columns[name] = verb
		return false
	}
	var d tableDecl
	dropped := make(map[string]bool)
	for _, c := range cmds {
		switch c.Action {
		case ast.AT_ADD_COLUMN, ast.AT_MODIFY_COLUMN:
			verb := "adds"
			if c.Action == ast.AT_MODIFY_COLUMN {
				verb = "modifies"
			}
			for _, col := range alterColumns(c) {
				if b.tooLong(col.Name, rangeOf(col.Loc)) || named(col.Name, verb, rangeOf(col.Loc)) {
					return nil, true
				}
				b.declareColumn(&d, col)
			}
		case ast.AT_DROP_COLUMN:
			for _, name := range b.droppedColumns(c) {
				if named(name, "drops", rangeOf(c.Loc)) {
					return nil, true
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
				b.fail(rangeOf(c.Loc), review.DuplicateName, twice("drops", "constraint", c.Constraint.Name))
				return nil, true
			}
			dropped[c.Constraint.Name] = true
		case ast.AT_RENAME_COLUMN, ast.AT_RENAME_CONSTRAINT, ast.AT_RENAME:
			if b.tooLong(c.NewName, rangeOf(c.Loc)) {
				return nil, true
			}
		}
	}
	return d.names, b.namesTooLong(nil, &d) || b.declTextFails(&d, r)
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

func (b *builder) alterCmd(r objRef, kind string, c *ast.AlterTableCmd, declared []string, addedCols map[string]bool, addedPK bool) bool {
	rng := rangeOf(c.Loc)
	if rng == (review.Range{}) {
		rng = r.rng
	}
	switch c.Action {
	case ast.AT_ADD_COLUMN:
		if nestedStorage(c.Options) {
			// The nested table's storage table takes a name.
			return false
		}
		pkFact, _ := b.look(r.childKey(spacePrimary, ""))
		selfHasPK := addedPK || pkFact != walk.Absent
		var d tableDecl
		for _, col := range alterColumns(c) {
			crng := rangeOf(col.Loc)
			if crng == (review.Range{}) {
				crng = rng
			}
			if columnMissingDatatype(col) {
				b.fail(crng, review.UnsupportedForm, "Column "+review.Code(ident(col.Name))+" must specify a datatype.")
				return true
			}
			if b.checkColTypeKnownNonType(col.TypeName, crng) {
				return true
			}
			if b.checkColCollationAndRef(col, nil, crng) {
				return true
			}
			if col.Constraints != nil {
				for _, item := range col.Constraints.Items {
					if cc, ok := item.(*ast.ColumnConstraint); ok && cc.Type == ast.CONSTRAINT_FOREIGN {
						if b.checkForeignKeyRef(r, cc.RefTable, cc.RefColumns, selfHasPK, nil, crng) {
							return true
						}
					}
				}
			}
			ck := r.childKey(spaceColumn, col.Name)
			b.need(ck, walk.Absent, nil, crng, "", takenColumn(col.Name, r.written))
			b.Set(ck, walk.Present, "")
			b.declareColumn(&d, col)
		}
		b.declChecks(&d, r)
		b.declEffects(&d, r)
		b.mutateColMeta(r)
	case ast.AT_MODIFY_COLUMN:
		pkFact, _ := b.look(r.childKey(spacePrimary, ""))
		selfHasPK := addedPK || pkFact != walk.Absent
		st := b.syncedColMeta(r)
		var d tableDecl
		for _, col := range alterColumns(c) {
			crng := rangeOf(col.Loc)
			if crng == (review.Range{}) {
				crng = rng
			}
			b.columnNeeds(r, col.Name, crng)
			if b.checkModifyColMeta(r, st, col, crng) {
				return true
			}
			if col.Constraints != nil {
				for _, item := range col.Constraints.Items {
					if cc, ok := item.(*ast.ColumnConstraint); ok && cc.Type == ast.CONSTRAINT_FOREIGN {
						if b.checkForeignKeyRef(r, cc.RefTable, cc.RefColumns, selfHasPK, nil, crng) {
							return true
						}
					}
				}
			}
			b.declareColumn(&d, col)
			if col.Null {
				// NULL drops the column's NOT NULL constraint, of any name.
				b.dropsConstraints(r)
			}
		}
		b.declChecks(&d, r)
		b.declEffects(&d, r)
		b.mutateColMeta(r)
	case ast.AT_DROP_COLUMN, ast.AT_SET_UNUSED:
		names := b.droppedColumns(c)
		if names == nil {
			return false
		}
		pst := b.syncedPartMeta(r)
		for _, name := range names {
			b.columnNeeds(r, name, rng)
			if pst != nil && pst.partKeyCols[name] {
				b.fail(rng, review.UnsupportedForm, "Cannot drop partitioning key column "+review.Code(ident(name))+" from "+review.Code(r.written)+".")
				return true
			}
			b.Set(r.childKey(spaceColumn, name), walk.Absent, "")
		}
		// The column's constraints and indexes go with it, and a nested
		// table column's storage table.
		b.dropsConstraints(r)
		b.forgetStorage(r)
		b.MayRemove(walk.Scope{Space: spacePrimary, Parent: r.path()})
		b.dropIndexes(r)
		b.mutateColMeta(r)
		if b.cascadesConstraints(c) {
			b.cascadeConstraints()
		}
	case ast.AT_RENAME_COLUMN:
		from, to := r.childKey(spaceColumn, c.ColumnName), r.childKey(spaceColumn, c.NewName)
		b.columnNeeds(r, c.ColumnName, rng)
		b.need(to, walk.Absent, nil, rng, "", takenColumn(c.NewName, r.written))
		if c.ColumnName != c.NewName {
			b.Set(from, walk.Absent, "")
			b.Set(to, walk.Present, "")
		}
		b.mutateTableAll(r)
	case ast.AT_ADD_CONSTRAINT:
		if c.Constraint == nil {
			return false
		}
		if c.Constraint.Columns != nil {
			seen := make(map[string]bool)
			for _, item := range c.Constraint.Columns.Items {
				if s, ok := item.(*ast.String); ok && s.Str != "" && !strings.Contains(s.Str, ".") {
					if seen[s.Str] {
						b.fail(rng, review.DuplicateName, twice("lists", "column", s.Str))
						return true
					}
					seen[s.Str] = true
					if !addedCols[s.Str] {
						b.columnNeeds(r, s.Str, rng)
					}
				}
			}
		}
		if c.Constraint.Type == ast.CONSTRAINT_FOREIGN {
			pkFact, _ := b.look(r.childKey(spacePrimary, ""))
			selfHasPK := addedPK || pkFact != walk.Absent
			if b.checkForeignKeyRef(r, c.Constraint.RefTable, c.Constraint.RefColumns, selfHasPK, nil, rng) {
				return true
			}
		}
		var d tableDecl
		b.declareConstraint(&d, c.Constraint)
		b.declChecks(&d, r)
		b.declEffects(&d, r)
		b.mutateColMeta(r)
		b.mutateIndexes(r)
	case ast.AT_DROP_CONSTRAINT:
		b.mutateColMeta(r)
		b.mutateIndexes(r)
		return b.dropConstraint(r, c, rng)
	case ast.AT_MODIFY_CONSTRAINT, ast.AT_ENABLE_DISABLE:
		b.mutateColMeta(r)
		b.mutateIndexes(r)
		return b.constraintState(r, c, rng, declared, addedPK)
	case ast.AT_RENAME_CONSTRAINT:
		from, to := c.ColumnName, c.NewName
		b.constraintNeeds(r, from, rng)
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
	case ast.AT_DROP_UNUSED_COLUMNS:
		// An unused nested table column's storage table goes now.
		b.forgetStorage(r)
	case ast.AT_MODIFY_NESTED_TABLE:
		if c.ColumnName != "" && !strings.Contains(c.ColumnName, ".") && !strings.EqualFold(c.ColumnName, "COLUMN_VALUE") {
			b.columnNeeds(r, c.ColumnName, rng)
			if st := b.syncedColMeta(r); st != nil {
				if sc, ok := st.colMeta[c.ColumnName]; ok && isOracleScalarType(sc.typ) {
					b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(c.ColumnName))+" in "+review.Code(r.written)+" is not a nested table column.")
					return true
				}
			}
		}
	case ast.AT_MODIFY_LOB:
		if col := b.modifyLobColumn(c.Loc); col != "" && !strings.EqualFold(col, "COLUMN_VALUE") {
			b.columnNeeds(r, col, rng)
			if st := b.syncedColMeta(r); st != nil {
				if sc, ok := st.colMeta[col]; ok && isOracleScalarType(sc.typ) {
					b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(col))+" in "+review.Code(r.written)+" is not a LOB column.")
					return true
				}
			}
		}
	case ast.AT_MODIFY_VARRAY:
		if col := b.modifyVarrayColumn(c.Loc); col != "" && !strings.EqualFold(col, "COLUMN_VALUE") {
			b.columnNeeds(r, col, rng)
			if st := b.syncedColMeta(r); st != nil {
				if sc, ok := st.colMeta[col]; ok && isOracleScalarType(sc.typ) {
					b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(col))+" in "+review.Code(r.written)+" is not a VARRAY column.")
					return true
				}
			}
		}
	case ast.AT_ALTER_PROPERTY:
		switch c.Subtype {
		case "UPGRADE":
			if b.checkAlterUpgrade(c.Loc, rng) {
				return true
			}
		case "INMEMORY", "NO INMEMORY":
			if b.checkAlterInmemory(r, c.Loc, rng) {
				return true
			}
		case "MODIFY COLUMN SUBSTITUTABLE":
			if c.ColumnName != "" && !strings.Contains(c.ColumnName, ".") {
				b.columnNeeds(r, c.ColumnName, rng)
				if st := b.syncedColMeta(r); st != nil {
					if sc, ok := st.colMeta[c.ColumnName]; ok && isOracleScalarType(sc.typ) {
						b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(c.ColumnName))+" in "+review.Code(r.written)+" is not a substitutable object column.")
						return true
					}
				}
			}
		}
	case ast.AT_MODIFY_TO_PARTITIONED:
		if kind == kindMView {
			// ORA-14427: cannot modify a materialized view to partitioned.
			b.need(r.key(spaceTable), walk.Present, []string{kindTable}, r.rng, "", wrongKind(r.written, kindMView, kindWord(kindTable)))
			return true
		}
		if kind == kindTable {
			if b.checkModifyPartitionBy(r, c.Loc, rng) {
				return true
			}
			b.mutatePartitions(r)
		}
	case ast.AT_ADD_PARTITION, ast.AT_DROP_PARTITION, ast.AT_TRUNCATE_PARTITION, ast.AT_SPLIT_PARTITION,
		ast.AT_SPLIT_SUBPARTITION, ast.AT_MERGE_PARTITIONS, ast.AT_MERGE_SUBPARTITIONS, ast.AT_EXCHANGE_PARTITION,
		ast.AT_COALESCE_PARTITION, ast.AT_COALESCE_SUBPARTITION, ast.AT_MODIFY_PARTITION, ast.AT_MODIFY_SUBPARTITION,
		ast.AT_MODIFY_DEFAULT_ATTRS, ast.AT_DROP_SUBPARTITION, ast.AT_RENAME_PARTITION, ast.AT_SET_INTERVAL,
		ast.AT_SET_PARTITIONING, ast.AT_SET_SUBPARTITION_TEMPLATE, ast.AT_MODIFY_PARTITIONSET:
		if b.checkPartitionTextContradictions(c.Loc, r, true, rng) {
			return true
		}
		if kind == kindTable {
			if b.checkPartitionMeta(r, b.syncedPartMeta(r), c.Loc, rng) {
				return true
			}
			if b.checkPartitionColClauses(r, b.syncedColMeta(r), c.Loc, rng) {
				return true
			}
			b.mutatePartitions(r)
		}
	case ast.AT_MOVE:
		if b.moveHasPartition(c.Loc) {
			if b.checkPartitionTextContradictions(c.Loc, r, true, rng) {
				return true
			}
			if kind == kindTable {
				if b.checkPartitionMeta(r, b.syncedPartMeta(r), c.Loc, rng) {
					return true
				}
				b.mutatePartitions(r)
			}
		}
	case ast.AT_ENABLE_DISABLE_TRIGGERS, ast.AT_ENABLE_DISABLE_TABLE_LOCK,
		ast.AT_SHRINK_SPACE, ast.AT_MODIFY_OPAQUE_TYPE, ast.AT_IMMUTABLE_TABLE, ast.AT_BLOCKCHAIN_TABLE, ast.AT_DUPLICATED_REFRESH:
		// Storage and table-level flags: no name the walk keeps.
	default:
		return false
	}
	return true
}

// modifyLobColumn extracts the top-level column name from MODIFY LOB ( col ... ).
func (b *builder) modifyLobColumn(loc ast.Loc) string {
	toks := tokens(b.w.sql, loc)
	for i, tok := range toks {
		if tok.Type == tokenOpenParen && i+1 < len(toks) {
			if name, ok := tokenName(toks[i+1]); ok {
				return name
			}
			return ""
		}
	}
	return ""
}

// modifyVarrayColumn extracts the top-level column name from MODIFY VARRAY col (...).
func (b *builder) modifyVarrayColumn(loc ast.Loc) string {
	toks := tokens(b.w.sql, loc)
	for i, tok := range toks {
		if tok.Type != tokenQuoted && strings.EqualFold(tok.Str, "VARRAY") && i+1 < len(toks) {
			if name, ok := tokenName(toks[i+1]); ok {
				return name
			}
			return ""
		}
	}
	return ""
}

// partitionByColumns extracts the partitioning key columns from
// MODIFY PARTITION BY {RANGE|LIST|HASH} (col1, col2, ...).
// PARTITION BY REFERENCE names a constraint, not a column, and returns nil.
func (b *builder) partitionByColumns(loc ast.Loc) []string {
	toks := tokens(b.w.sql, loc)
	isWord := func(i int, w string) bool {
		return i < len(toks) && toks[i].Type != tokenQuoted && strings.EqualFold(toks[i].Str, w)
	}
	for i := 0; i+3 < len(toks); i++ {
		if !isWord(i, "PARTITION") || !isWord(i+1, "BY") {
			continue
		}
		if !isWord(i+2, "RANGE") && !isWord(i+2, "LIST") && !isWord(i+2, "HASH") {
			return nil
		}
		j := i + 3
		if toks[j].Type != tokenOpenParen {
			return nil
		}
		var out []string
		for j++; j < len(toks) && toks[j].Type != tokenCloseParen; j++ {
			if toks[j].Type == tokenComma {
				continue
			}
			if name, ok := tokenName(toks[j]); ok {
				out = append(out, name)
			} else {
				return nil
			}
		}
		return out
	}
	return nil
}

// constraintNeeds checks that a constraint the clause names is the
// table's: on the table, and, since constraint names are unique in the
// schema, not another table's (02443, 23292, 02430, 02431). Another
// table's constraint of the name leaves this table without one, the
// condition of review.ObjectNotFound.
func (b *builder) constraintNeeds(r objRef, name string, rng review.Range) {
	message := missingIn("Constraint", name, r.written)
	b.need(r.childKey(spaceTableConstraint, name), walk.Present, nil, rng, message, "")
	b.needNote(r.schemaKey(spaceConstraint, name), walk.Present, []string{r.path()}, noteOf(rng, review.ObjectNotFound, message))
}

// nestedStorage reports whether an ADD names a nested table's storage
// table.
func nestedStorage(opts *ast.List) bool {
	if opts == nil {
		return false
	}
	for _, item := range opts.Items {
		if o, ok := item.(*ast.DDLOption); ok && o.Key == "NESTED" {
			return true
		}
	}
	return false
}

// columnNeeds checks that a column the clause names exists (00904).
func (b *builder) columnNeeds(r objRef, name string, rng review.Range) {
	b.need(r.childKey(spaceColumn, name), walk.Present, nil, rng, missingIn("Column", name, r.written), "")
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
		b.need(pk, walk.Present, nil, rng, review.Code(r.written)+" has no primary key.", "")
		b.Set(pk, walk.Absent, "")
		b.dropsConstraints(r)
	case c.Subtype == "UNIQUE":
		if c.Constraint != nil && c.Constraint.Columns != nil {
			seen := make(map[string]bool)
			for _, item := range c.Constraint.Columns.Items {
				if s, ok := item.(*ast.String); ok && s.Str != "" && !strings.Contains(s.Str, ".") {
					if seen[s.Str] {
						b.fail(rng, review.DuplicateName, twice("lists", "column", s.Str))
						return true
					}
					seen[s.Str] = true
					b.columnNeeds(r, s.Str, rng)
				}
			}
		}
		b.dropsConstraints(r)
	case c.Constraint != nil && c.Constraint.Name != "":
		name := c.Constraint.Name
		ck := r.childKey(spaceTableConstraint, name)
		b.constraintNeeds(r, name, rng)
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
// modify, verified). A statement's trailing ENABLE or DISABLE CONSTRAINT
// also sees the constraints its ADD and MODIFY clauses declare, which
// MODIFY CONSTRAINT does not (02431, verified). DISABLE drops the index
// the constraint made, and ENABLE makes it again, under the constraint's
// name or another (spec §4).
func (b *builder) constraintState(r objRef, c *ast.AlterTableCmd, rng review.Range, declared []string, addedPK bool) bool {
	name := c.ColumnName
	if c.Action == ast.AT_MODIFY_CONSTRAINT && c.Constraint != nil {
		name = c.Constraint.Name
	}
	switch {
	case name != "":
		if c.Action == ast.AT_MODIFY_CONSTRAINT || !contains(declared, name) {
			b.constraintNeeds(r, name, rng)
		}
		b.Set(r.schemaKey(spaceIndex, name), walk.Unknown, "")
		if b.createsIndex(c.Loc) {
			b.MayAdd(walk.Scope{Space: spaceIndex, Parent: walk.Path(r.schema)})
		}
	case strings.HasSuffix(c.Subtype, "PRIMARY KEY"), strings.HasSuffix(c.Subtype, "UNIQUE"):
		if strings.HasSuffix(c.Subtype, "PRIMARY KEY") && !addedPK {
			b.need(r.childKey(spacePrimary, ""), walk.Present, nil, rng, review.Code(r.written)+" has no primary key.", "")
		}
		if strings.HasSuffix(c.Subtype, "UNIQUE") {
			cols := b.uniqueClauseColumns(c)
			if dup := firstDupString(cols); dup != "" {
				b.fail(rng, review.DuplicateName, twice("lists", "column", dup))
				return true
			}
			for _, col := range cols {
				b.columnNeeds(r, col, rng)
			}
		}
		// The constraint's name, and so its index's, is not in the text.
		b.MayAdd(walk.Scope{Space: spaceIndex, Parent: walk.Path(r.schema)})
	default:
		return false
	}
	b.dropIndexes(r)
	return true
}

// uniqueClauseColumns returns the column names listed in an ENABLE/DISABLE/MODIFY
// UNIQUE (col1, col2, ...) clause.
func (b *builder) uniqueClauseColumns(c *ast.AlterTableCmd) []string {
	if c.Constraint != nil && c.Constraint.Columns != nil {
		var out []string
		for _, item := range c.Constraint.Columns.Items {
			if s, ok := item.(*ast.String); ok && s.Str != "" && !strings.Contains(s.Str, ".") {
				out = append(out, s.Str)
			}
		}
		return out
	}
	toks := tokens(b.w.sql, c.Loc)
	for i, tok := range toks {
		if tok.Type != tokenQuoted && strings.EqualFold(tok.Str, "UNIQUE") && i+1 < len(toks) && toks[i+1].Type == tokenOpenParen {
			var out []string
			for j := i + 2; j < len(toks) && toks[j].Type != tokenCloseParen; j++ {
				if toks[j].Type == tokenComma {
					continue
				}
				if j+1 < len(toks) && toks[j+1].Str == "." {
					// Skip dotted object-attribute paths (e.g. o.x).
					for j < len(toks) && toks[j].Type != tokenComma && toks[j].Type != tokenCloseParen {
						j++
					}
					if j < len(toks) && toks[j].Type == tokenCloseParen {
						break
					}
					continue
				}
				if name, ok := tokenName(toks[j]); ok {
					out = append(out, name)
				} else {
					return nil
				}
			}
			return out
		}
	}
	return nil
}
