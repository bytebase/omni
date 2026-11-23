package review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// alterTable follows ALTER TABLE. MySQL judges DROP, RENAME, CHANGE, and
// MODIFY against the table as it was and ADD after them, and resolves
// index columns against the final ones (verified): the drops and renames
// are lowered first, then the rest, then a RENAME TO of the table itself,
// each subcommand a part of its own (walk.Passes).
func (b *builder) alterTable(v *ast.AlterTableStmt) bool {
	r, st := b.ref(v.Table)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	if alterIsPartitionMaintenanceOnly(v.Commands) {
		// ANALYZE, CHECK, OPTIMIZE, and REPAIR PARTITION report a status
		// row rather than a SQL error when the table is missing or a view,
		// and change no names the walk tracks (measured on mysql:8.0.46).
		return true
	}
	f, temp, kind := b.resolve(r)
	switch {
	case temp:
		// A temporary table's columns and keys are not followed.
		return false
	case f == walk.Unknown:
		for _, c := range v.Commands {
			if !alterNoFacts(c) {
				return false
			}
		}
		return true
	}
	b.tableNeeds(r)
	if f == walk.Absent || kind == relView {
		return true
	}
	// The table keeps foreign keys as its engine does; a statement that
	// changes the engine claims nothing of them.
	rel := relKey(r.db, r.name)
	engine, sets := alterEngine(v)
	b.fks = !sets && !b.w.ignoresFKs[rel]
	if sets {
		ignores := !keepsFKs(engine)
		b.Session(func() {
			if ignores {
				b.w.ignoresFKs[rel] = true
			} else {
				delete(b.w.ignoresFKs, rel)
			}
			b.w.engines[rel] = engine
		})
	}
	delta := b.buildAlterDelta(r, v, engine, sets)
	if b.alterTextFails(v, r, delta) {
		return true
	}
	var rename *ast.AlterTableCmd
	for _, c := range v.Commands {
		if c.Type == ast.ATRenameTable {
			rename = c
		}
	}
	b.renames = rename != nil
	for _, early := range []bool{true, false} {
		for _, c := range v.Commands {
			if c.Type == ast.ATRenameTable {
				continue
			}
			if alterEarly(c.Type) != early {
				continue
			}
			b.Next()
			b.early = early
			if !b.alterCmd(r, c, delta) {
				return false
			}
		}
	}
	b.applyAlterMeta(r, v, delta)
	if rename != nil {
		b.Next()
		return b.alterRename(r, rename)
	}
	return true
}

// alterDelta summarizes the table-wide changes of one ALTER TABLE statement so
// subcommands can check final column visibility, index coverage, partitioning,
// and foreign key dependencies without false positives across sibling actions.
type alterDelta struct {
	addedCols           map[string]string // colKey -> baseType
	droppedCols         map[string]bool
	renamedFrom         map[string]string // oldColKey -> newColKey
	renamedTo           map[string]string // newColKey -> oldColKey
	modifiedCols        map[string]string // colKey -> new baseType
	droppedIndexes      map[string]bool
	addedIndexes        []*indexMeta
	droppedFKs          map[string]bool
	removesPartitioning bool
	repartition         *ast.PartitionClause
	finalPartitioned    bool
	finalEngine         string
	engineKnown         bool
}

func (b *builder) buildAlterDelta(r tableRef, v *ast.AlterTableStmt, setEngine string, setsEngine bool) *alterDelta {
	rel := relKey(r.db, r.name)
	d := &alterDelta{
		addedCols:      make(map[string]string),
		droppedCols:    make(map[string]bool),
		renamedFrom:    make(map[string]string),
		renamedTo:      make(map[string]string),
		modifiedCols:   make(map[string]string),
		droppedIndexes: make(map[string]bool),
		droppedFKs:     make(map[string]bool),
	}
	if setsEngine {
		d.finalEngine = setEngine
		d.engineKnown = true
	} else if eng, ok := b.w.engines[rel]; ok {
		d.finalEngine = eng
		d.engineKnown = true
	}
	if pi := b.w.partitionState[rel]; pi != nil {
		d.finalPartitioned = pi.partitioned
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATAddColumn:
			cols := c.Columns
			if c.Column != nil {
				cols = append([]*ast.ColumnDef{c.Column}, cols...)
			}
			for _, col := range cols {
				if k, ok := b.key(spaceColumn, col.Name); ok {
					d.addedCols[k] = astBaseType(col.TypeName)
				}
				d.recordInlineColumnKeys(col)
			}
			for _, con := range c.Constraints {
				d.recordAddedConstraint(con)
			}
		case ast.ATDropColumn:
			if k, ok := b.key(spaceColumn, c.Name); ok {
				d.droppedCols[k] = true
			}
		case ast.ATModifyColumn:
			if c.Column != nil {
				if k, ok := b.key(spaceColumn, c.Column.Name); ok {
					d.modifiedCols[k] = astBaseType(c.Column.TypeName)
				}
				d.recordInlineColumnKeys(c.Column)
			}
		case ast.ATChangeColumn:
			if c.Column != nil {
				fk, ok1 := b.key(spaceColumn, c.Name)
				tk, ok2 := b.key(spaceColumn, c.Column.Name)
				if ok1 && ok2 {
					if fk != tk {
						d.renamedFrom[fk] = tk
						d.renamedTo[tk] = fk
					}
					d.modifiedCols[tk] = astBaseType(c.Column.TypeName)
				}
				d.recordInlineColumnKeys(c.Column)
			}
		case ast.ATRenameColumn:
			fk, ok1 := b.key(spaceColumn, c.Name)
			tk, ok2 := b.key(spaceColumn, c.NewName)
			if ok1 && ok2 && fk != tk {
				d.renamedFrom[fk] = tk
				d.renamedTo[tk] = fk
			}
		case ast.ATDropIndex:
			if k, ok := b.key(spaceIndex, c.Name); ok {
				d.droppedIndexes[k] = true
			}
		case ast.ATDropConstraint:
			switch dropKindOf(c) {
			case dropPrimaryKey:
				d.droppedIndexes[primaryIndex] = true
			case dropForeignKey:
				if k, ok := b.key(spaceForeignKey, c.Name); ok {
					d.droppedFKs[k] = true
				}
			case dropConstraint:
				if k, ok := b.key(spaceForeignKey, c.Name); ok {
					d.droppedFKs[k] = true
					d.droppedIndexes[k] = true
				}
			}
		case ast.ATAddConstraint, ast.ATAddIndex:
			d.recordAddedConstraint(c.Constraint)
		case ast.ATRemovePartitioning:
			d.removesPartitioning = true
			d.finalPartitioned = false
		case ast.ATPartitionBy:
			if c.PartitionBy != nil {
				d.repartition = c.PartitionBy
				d.finalPartitioned = true
			}
		}
	}
	return d
}

func (d *alterDelta) recordInlineColumnKeys(col *ast.ColumnDef) {
	if col == nil {
		return
	}
	for _, cc := range col.Constraints {
		if cc != nil && (cc.Type == ast.ColConstrPrimaryKey || cc.Type == ast.ColConstrUnique) {
			d.addedIndexes = append(d.addedIndexes, &indexMeta{
				unique:  true,
				primary: cc.Type == ast.ColConstrPrimaryKey,
				cols:    []string{strings.ToLower(col.Name)},
			})
		}
	}
}

func (d *alterDelta) recordAddedConstraint(con *ast.Constraint) {
	if con == nil || con.Type == ast.ConstrCheck || con.Type == ast.ConstrForeignKey {
		return
	}
	im := &indexMeta{
		unique:   con.Type == ast.ConstrUnique || con.Type == ast.ConstrPrimaryKey,
		primary:  con.Type == ast.ConstrPrimaryKey,
		fulltext: con.Type == ast.ConstrFulltextIndex,
		spatial:  con.Type == ast.ConstrSpatialIndex || isSpatialIndexType(con.IndexType),
	}
	for _, p := range constraintKeyParts(con) {
		if p.functional {
			break
		}
		if p.length > 0 {
			im.hasPrefix = true
		}
		im.cols = append(im.cols, strings.ToLower(p.name))
	}
	d.addedIndexes = append(d.addedIndexes, im)
}

func (b *builder) alterFinalCol(r tableRef, d *alterDelta, colName string) (walk.Fact, string, bool) {
	k, ok := b.key(spaceColumn, colName)
	if !ok {
		return walk.Unknown, "", false
	}
	if bt, ok := d.addedCols[k]; ok {
		return walk.Present, bt, false
	}
	if oldK, ok := d.renamedTo[k]; ok {
		if bt, has := d.modifiedCols[k]; has && bt != "" {
			return walk.Present, bt, false
		}
		return walk.Present, b.w.colTypes[colKey(r.db, r.name, oldK)], false
	}
	if d.droppedCols[k] {
		return walk.Absent, "", true
	}
	if _, moved := d.renamedFrom[k]; moved {
		return walk.Absent, "", true
	}
	ck := colKey(r.db, r.name, k)
	f, _ := b.l.Lookup(ck)
	if bt, has := d.modifiedCols[k]; has && bt != "" {
		return f, bt, false
	}
	return f, b.w.colTypes[ck], false
}

func (b *builder) applyAlterMeta(r tableRef, v *ast.AlterTableStmt, d *alterDelta) {
	rel := relKey(r.db, r.name)
	b.Session(func() {
		for k := range d.droppedCols {
			if _, readded := d.addedCols[k]; !readded {
				delete(b.w.colTypes, colKey(r.db, r.name, k))
			}
		}
		// Copy types for renamed columns before overwriting.
		renamedTypes := make(map[string]string, len(d.renamedTo))
		for newK, oldK := range d.renamedTo {
			if bt, ok := b.w.colTypes[colKey(r.db, r.name, oldK)]; ok {
				renamedTypes[newK] = bt
			}
		}
		for oldK := range d.renamedFrom {
			delete(b.w.colTypes, colKey(r.db, r.name, oldK))
		}
		for newK, bt := range renamedTypes {
			b.w.colTypes[colKey(r.db, r.name, newK)] = bt
		}
		for k, bt := range d.addedCols {
			if bt != "" {
				b.w.colTypes[colKey(r.db, r.name, k)] = bt
			}
		}
		for k, bt := range d.modifiedCols {
			if bt != "" {
				b.w.colTypes[colKey(r.db, r.name, k)] = bt
			}
		}
		if len(d.droppedCols) > 0 {
			b.w.tableIndexesComplete[rel] = false
		}
		for k := range d.droppedIndexes {
			if b.w.tableIndexes[rel] != nil {
				delete(b.w.tableIndexes[rel], k)
			}
		}
		for k := range d.droppedFKs {
			if b.w.tableFKs[rel] != nil {
				delete(b.w.tableFKs[rel], k)
			}
			for pRel, refs := range b.w.referencedByFK {
				delete(refs, fkRefKey{child: rel, fk: k})
				if len(refs) == 0 {
					delete(b.w.referencedByFK, pRel)
				}
			}
		}
		for _, c := range v.Commands {
			if c == nil {
				continue
			}
			switch c.Type {
			case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
				cols := c.Columns
				if c.Column != nil {
					cols = append([]*ast.ColumnDef{c.Column}, cols...)
				}
				for _, col := range cols {
					for _, cc := range col.Constraints {
						switch cc.Type {
						case ast.ColConstrPrimaryKey:
							if b.w.tableIndexes[rel] != nil {
								b.w.tableIndexes[rel][primaryIndex] = &indexMeta{primary: true, unique: true, cols: []string{strings.ToLower(col.Name)}}
							}
						case ast.ColConstrUnique:
							b.w.tableIndexesComplete[rel] = false
						}
					}
				}
				for _, con := range c.Constraints {
					b.applyAddedConstraintMeta(rel, con)
				}
			case ast.ATAddConstraint, ast.ATAddIndex:
				b.applyAddedConstraintMeta(rel, c.Constraint)
			case ast.ATRenameIndex:
				if fk, ok1 := b.key(spaceIndex, c.Name); ok1 {
					if tk, ok2 := b.key(spaceIndex, c.NewName); ok2 && fk != tk {
						if idxs := b.w.tableIndexes[rel]; idxs != nil {
							if im, ok := idxs[fk]; ok {
								delete(idxs, fk)
								idxs[tk] = im
							}
						}
					}
				}
			case ast.ATPartitionBy:
				if c.PartitionBy != nil {
					b.w.partitionState[rel] = partitionInfoFromClause(c.PartitionBy)
				}
			case ast.ATRemovePartitioning:
				b.w.partitionState[rel] = &tablePartitionInfo{partitioned: false}
			case ast.ATAddPartition:
				if pi := b.w.partitionState[rel]; pi != nil && pi.partitioned {
					if len(c.PartitionDefs) == 0 {
						pi.names = nil
						if c.Number > 0 && pi.count > 0 {
							pi.count += c.Number
						} else {
							pi.count = 0
						}
					} else {
						for _, pd := range c.PartitionDefs {
							if pd != nil && pd.Name != "" && pi.names != nil {
								pi.names[strings.ToLower(pd.Name)] = true
							}
						}
						if pi.count > 0 {
							pi.count += len(c.PartitionDefs)
						}
					}
				}
			case ast.ATDropPartition:
				if pi := b.w.partitionState[rel]; pi != nil && pi.partitioned {
					for _, name := range c.PartitionNames {
						if pi.names != nil {
							delete(pi.names, strings.ToLower(name))
						}
					}
					if pi.count > 0 {
						pi.count = max(0, pi.count-len(c.PartitionNames))
					}
				}
			case ast.ATCoalescePartition, ast.ATReorganizePartition:
				if pi := b.w.partitionState[rel]; pi != nil && pi.partitioned {
					pi.names = nil
					pi.count = 0
				}
			}
		}
	})
}

func (b *builder) applyAddedConstraintMeta(rel walk.Key, con *ast.Constraint) {
	if con == nil || con.Type == ast.ConstrCheck {
		return
	}
	if con.Type == ast.ConstrForeignKey {
		b.w.tableIndexesComplete[rel] = false
		if b.fks {
			if k, ok := b.key(spaceForeignKey, con.Name); ok {
				if b.w.tableFKs[rel] == nil {
					b.w.tableFKs[rel] = make(map[string][]string)
				}
				var cols []string
				for _, p := range constraintKeyParts(con) {
					if !p.functional && p.name != "" {
						cols = append(cols, strings.ToLower(p.name))
					}
				}
				b.w.tableFKs[rel][k] = cols
			}
		}
		return
	}
	im := &indexMeta{
		unique:   con.Type == ast.ConstrUnique || con.Type == ast.ConstrPrimaryKey,
		primary:  con.Type == ast.ConstrPrimaryKey,
		fulltext: con.Type == ast.ConstrFulltextIndex,
		spatial:  con.Type == ast.ConstrSpatialIndex || isSpatialIndexType(con.IndexType),
	}
	for _, p := range constraintKeyParts(con) {
		if p.functional {
			break
		}
		if p.length > 0 {
			im.hasPrefix = true
		}
		im.cols = append(im.cols, strings.ToLower(p.name))
	}
	if con.Type == ast.ConstrPrimaryKey {
		if b.w.tableIndexes[rel] != nil {
			b.w.tableIndexes[rel][primaryIndex] = im
		}
		return
	}
	if k, ok := b.key(spaceIndex, con.Name); ok {
		if b.w.tableIndexes[rel] != nil {
			b.w.tableIndexes[rel][k] = im
		}
	} else {
		b.w.tableIndexesComplete[rel] = false
	}
}

// alterEngine returns the last ENGINE an ALTER TABLE sets, if it sets one.
func alterEngine(v *ast.AlterTableStmt) (string, bool) {
	var eng string
	var sets bool
	for _, c := range v.Commands {
		if c.Type == ast.ATTableOption && c.Option != nil &&
			(strings.EqualFold(c.Option.Name, "ENGINE") || strings.EqualFold(c.Option.Name, "STORAGE ENGINE")) {
			eng = engineOf([]*ast.TableOption{c.Option})
			sets = true
		}
	}
	return eng, sets
}

// alterEarly reports whether a subcommand is judged against the table as
// it was: a drop or a rename.
func alterEarly(t ast.AlterTableCmdType) bool {
	switch t {
	case ast.ATDropColumn, ast.ATDropIndex, ast.ATDropConstraint, ast.ATChangeColumn, ast.ATRenameColumn, ast.ATRenameIndex:
		return true
	}
	return false
}

// alterIsPartitionMaintenanceOnly reports whether an ALTER TABLE consists
// solely of partition maintenance actions (ANALYZE, CHECK, OPTIMIZE, REPAIR
// PARTITION), which MySQL answers with a status row instead of a SQL error
// even when the table is missing or a view.
func alterIsPartitionMaintenanceOnly(cmds []*ast.AlterTableCmd) bool {
	if len(cmds) == 0 {
		return false
	}
	for _, c := range cmds {
		if c == nil {
			return false
		}
		switch c.Type {
		case ast.ATAnalyzePartition, ast.ATCheckPartition, ast.ATOptimizePartition, ast.ATRepairPartition:
		default:
			return false
		}
	}
	return true
}

// alterNoFacts reports whether a subcommand changes no name the walk
// keeps: table options, partitions, tablespaces, and CHECK constraints,
// whose names the walk does not follow.
func alterNoFacts(c *ast.AlterTableCmd) bool {
	switch c.Type {
	case ast.ATTableOption, ast.ATConvertCharset, ast.ATAlgorithm, ast.ATLock, ast.ATForce, ast.ATOrderBy,
		ast.ATEnableKeys, ast.ATDisableKeys, ast.ATWithValidation, ast.ATWithoutValidation, ast.ATAlterCheckEnforced,
		ast.ATSecondaryLoad, ast.ATSecondaryUnload, ast.ATDiscardTablespace, ast.ATImportTablespace,
		ast.ATAddPartition, ast.ATDropPartition, ast.ATCoalescePartition, ast.ATReorganizePartition,
		ast.ATExchangePartition, ast.ATTruncatePartition, ast.ATAnalyzePartition, ast.ATCheckPartition,
		ast.ATOptimizePartition, ast.ATRebuildPartition, ast.ATRepairPartition, ast.ATDiscardPartitionTablespace,
		ast.ATImportPartitionTablespace, ast.ATRemovePartitioning, ast.ATPartitionBy:
		return true
	case ast.ATAddConstraint:
		return c.Constraint != nil && c.Constraint.Type == ast.ConstrCheck
	}
	return false
}

func (b *builder) alterCmd(r tableRef, c *ast.AlterTableCmd, d *alterDelta) bool {
	rng := rangeOf(c.Loc)
	rel := relKey(r.db, r.name)
	switch c.Type {
	case ast.ATAddColumn:
		cols := c.Columns
		if c.Column != nil {
			cols = append([]*ast.ColumnDef{c.Column}, cols...)
		}
		for _, col := range cols {
			if !b.addColumn(r, col, c.After, rangeOf(col.Loc)) {
				return false
			}
		}
		// The parenthesized form reads a table element list, so an index or a
		// constraint written beside the columns is part of the same action.
		for _, con := range c.Constraints {
			b.addConstraint(r, con, d)
		}
	case ast.ATDropColumn:
		k, ok := b.key(spaceColumn, c.Name)
		if !ok {
			return false
		}
		ck := colKey(r.db, r.name, k)
		b.need(ck, walk.Present, nil, rng, notFound("column", ident(c.Name), r.written), "") // 1091
		if f, _ := b.l.Lookup(ck); f == walk.Present {
			if !d.removesPartitioning && d.repartition == nil {
				if pi := b.w.partitionState[rel]; pi != nil && pi.partitioned && pi.cols[k] {
					b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(c.Name))+" in "+review.Code(r.written)+" is used by the table's partitioning function.") // 3855
					return true
				}
			}
			if b.fks {
				for fkName, fkCols := range b.w.tableFKs[rel] {
					if d.droppedFKs[fkName] {
						continue
					}
					for _, fc := range fkCols {
						if fc == k {
							b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(c.Name))+" in "+review.Code(r.written)+" is needed by a foreign key constraint.") // 1828
							return true
						}
					}
				}
			}
		}
		b.leave(ck)
		// An index loses the column, and goes with its last one.
		b.MayRemove(walk.Scope{Space: spaceIndex, Parent: tablePath(r.db, r.name)})
	case ast.ATModifyColumn:
		if c.Column == nil || !b.columnNeeds(r, c.Column.Name, rng) {
			return false
		}
		b.after(r, c.After, rng)
		b.alterColumnKeys(r, c.Column, rng)
	case ast.ATChangeColumn:
		if c.Column == nil || !b.renameColumn(r, c.Name, c.Column.Name, rng, d, false) {
			return false
		}
		b.after(r, c.After, rng)
		b.alterColumnKeys(r, c.Column, rng)
	case ast.ATRenameColumn:
		return b.renameColumn(r, c.Name, c.NewName, rng, d, true)
	case ast.ATAlterColumnDefault:
		if !b.columnNeeds(r, c.Name, rng) {
			return false
		}
		if k, ok := b.key(spaceColumn, c.Name); ok && c.DefaultExpr != nil {
			ck := colKey(r.db, r.name, k)
			if f, _ := b.l.Lookup(ck); f == walk.Present {
				b.columnDefaultFails(c.Name, b.w.colTypes[ck], c.DefaultExpr, rng)
			}
		}
	case ast.ATAlterColumnVisible, ast.ATAlterColumnInvisible:
		return b.columnNeeds(r, c.Name, rng)
	case ast.ATAddConstraint, ast.ATAddIndex:
		if c.Constraint == nil {
			return false
		}
		b.addConstraint(r, c.Constraint, d)
	case ast.ATDropIndex:
		if k, ok := b.key(spaceIndex, c.Name); ok {
			if b.dropIndexNeededByFK(r, k, d.droppedIndexes, d.droppedFKs, d.addedIndexes) {
				b.fail(rng, nil, review.IndexBacksConstraint, "Index "+review.Code(ident(c.Name))+" in "+review.Code(r.written)+" is needed by a foreign key constraint.")
				return true
			}
		}
		return b.dropIndexName(r, c.Name, rng)
	case ast.ATDropConstraint:
		switch dropKindOf(c) {
		case dropPrimaryKey:
			if b.dropIndexNeededByFK(r, primaryIndex, d.droppedIndexes, d.droppedFKs, d.addedIndexes) {
				b.fail(rng, nil, review.IndexBacksConstraint, "Primary key of "+review.Code(r.written)+" is needed by a foreign key constraint.")
				return true
			}
			return b.dropIndexName(r, "PRIMARY", rng)
		case dropForeignKey:
			k, ok := b.key(spaceForeignKey, c.Name)
			if !ok {
				return false
			}
			fk, own := fkKeys(r, k)
			if !b.fks {
				b.Set(fk, walk.Unknown, "")
				b.Set(own, walk.Unknown, "")
				break
			}
			// The name must be the table's: another table's foreign key of
			// the name fails as a missing one does (1091, verified).
			msg := notFound("foreign key", ident(c.Name), r.written)
			b.need(fk, walk.Present, nil, rng, msg, "")  // 1091
			b.need(own, walk.Present, nil, rng, msg, "") // 1091
			b.Set(fk, walk.Absent, "")
			b.Set(own, walk.Absent, "")
			b.onForeignKeyDropped(r, k)
		case dropCheck:
			k, ok := b.key(spaceCheckConstraint, c.Name)
			if !ok {
				return false
			}
			ck := chkKey(r.db, r.name, k)
			b.need(ck, walk.Present, nil, rng, notFound("check constraint", ident(c.Name), r.written), "") // 3821
			b.Set(ck, walk.Absent, "")
		case dropConstraint:
			// The foreign key, unique key, or check of the name (not a plain index).
			if k, ok := b.key(spaceForeignKey, c.Name); ok {
				fk, own := fkKeys(r, k)
				ck := chkKey(r.db, r.name, k)
				ik := idxKey(r.db, r.name, k)
				cf, _ := b.l.Lookup(ck)
				ff, _ := b.l.Lookup(own)
				if_, _ := b.l.Lookup(ik)
				idxNonUnique := false
				if im := b.w.tableIndexes[rel][k]; im != nil && !im.unique && !im.primary {
					idxNonUnique = true
				}
				if cf == walk.Absent && ff == walk.Absent && (if_ == walk.Absent || idxNonUnique) {
					b.fail(rng, nil, review.ObjectNotFound, notFound("constraint", ident(c.Name), r.written)) // 3940
					return true
				}
				b.Set(fk, walk.Unknown, "")
				b.Set(own, walk.Unknown, "")
				b.Set(ck, walk.Unknown, "")
				b.Set(ik, walk.Unknown, "")
			}
		default:
			return false
		}
	case ast.ATAlterCheckEnforced:
		k, ok := b.key(spaceCheckConstraint, c.Name)
		if !ok {
			return false
		}
		ck := chkKey(r.db, r.name, k)
		if c.ConstraintKeyword == ast.ConstraintKeywordCheck {
			b.need(ck, walk.Present, nil, rng, notFound("check constraint", ident(c.Name), r.written), "") // 3821
			return true
		}
		// ALTER CONSTRAINT name [NOT] ENFORCED: only CHECK constraints are supported (3941),
		// and a missing constraint fails with 3940.
		_, own := fkKeys(r, k)
		ik := idxKey(r.db, r.name, k)
		cf, _ := b.l.Lookup(ck)
		ff, _ := b.l.Lookup(own)
		if_, _ := b.l.Lookup(ik)
		im := b.w.tableIndexes[rel][k]
		isKnownUnique := im != nil && (im.unique || im.primary)
		isKnownNonUnique := im != nil && !im.unique && !im.primary
		if cf == walk.Absent {
			if ff == walk.Present || (if_ == walk.Present && isKnownUnique) {
				b.fail(rng, nil, review.UnsupportedForm, "Constraint "+review.Code(ident(c.Name))+" in "+review.Code(r.written)+" is not a check constraint.") // 3941
				return true
			}
			if ff == walk.Absent && (if_ == walk.Absent || isKnownNonUnique) {
				b.fail(rng, nil, review.ObjectNotFound, notFound("constraint", ident(c.Name), r.written)) // 3940
				return true
			}
		}
	case ast.ATRenameIndex:
		from, ok1 := b.key(spaceIndex, c.Name)
		to, ok2 := b.key(spaceIndex, c.NewName)
		if !ok1 || !ok2 || from == primaryIndex || to == primaryIndex {
			return false
		}
		fk, tk := idxKey(r.db, r.name, from), idxKey(r.db, r.name, to)
		b.need(fk, walk.Present, nil, rng, notFound("index", ident(c.Name), r.written), "") // 1091
		if from != to {
			b.need(tk, walk.Absent, nil, rng, "", alreadyExists("index", ident(c.NewName), r.written)) // 1061
			b.leave(fk)
			b.arrive(tk)
		}
		b.consumeFKGenerated(fk)
	case ast.ATAlterIndexVisible, ast.ATAlterIndexInvisible:
		k, ok := b.key(spaceIndex, c.Name)
		if !ok {
			return false
		}
		ik := idxKey(r.db, r.name, k)
		b.need(ik, walk.Present, nil, rng, notFound("index", ident(c.Name), r.written), "")
		b.consumeFKGenerated(ik)
	case ast.ATExchangePartition:
		if c.ExchangeTable != nil {
			exch, st := b.ref(c.ExchangeTable)
			switch st {
			case refFailed:
				return true
			case refOK:
				if exch.db == r.db && exch.name == r.name {
					b.fail(rng, nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" cannot exchange a partition with itself.") // 1066
					return true
				}
				if ef, temp, _ := b.resolve(exch); !temp && ef != walk.Unknown {
					b.tableNeeds(exch)
				}
			}
		}
		if pi := b.w.partitionState[rel]; pi != nil && !pi.partitioned {
			b.fail(rng, nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not partitioned.") // 1505
			return true
		}
	case ast.ATAddPartition, ast.ATDropPartition, ast.ATCoalescePartition, ast.ATReorganizePartition,
		ast.ATTruncatePartition, ast.ATRebuildPartition, ast.ATDiscardPartitionTablespace,
		ast.ATImportPartitionTablespace, ast.ATRemovePartitioning:
		pi := b.w.partitionState[rel]
		if pi != nil && !pi.partitioned {
			b.fail(rng, nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not partitioned.") // 1505
			return true
		}
		if pi != nil && pi.partitioned {
			rangeOrList := pi.hasType && (pi.partType == ast.PartitionRange || pi.partType == ast.PartitionList)
			switch c.Type {
			case ast.ATAddPartition:
				if rangeOrList && len(c.PartitionDefs) == 0 {
					b.fail(rng, nil, review.UnsupportedForm, "Each partition must be defined for "+review.Code("RANGE")+" or "+review.Code("LIST")+" partitioning.") // 1492
					return true
				}
				seenPart := make(map[string]bool, len(c.PartitionDefs))
				for _, pd := range c.PartitionDefs {
					if pd == nil {
						continue
					}
					if rangeOrList && pd.Values == nil {
						b.fail(rng, nil, review.UnsupportedForm, "Partition "+review.Code(ident(pd.Name))+" requires a "+review.Code("VALUES")+" definition.") // 1479
						return true
					}
					pk := strings.ToLower(pd.Name)
					if pk != "" {
						if seenPart[pk] {
							b.fail(rng, nil, review.DuplicateName, statesTwice("adds", "partition", ident(pd.Name))) // 1517
							return true
						}
						seenPart[pk] = true
						if pi.names != nil && pi.names[pk] {
							b.fail(rng, nil, review.NameTaken, alreadyExists("partition", ident(pd.Name), r.written)) // 1517
							return true
						}
					}
				}
			case ast.ATReorganizePartition:
				if rangeOrList && len(c.PartitionNames) == 0 && len(c.PartitionDefs) == 0 {
					b.fail(rng, nil, review.UnsupportedForm, review.Code("REORGANIZE PARTITION")+" without parameters is only valid for "+review.Code("HASH")+" or "+review.Code("KEY")+" partitioning.") // 1511
					return true
				}
				for _, pd := range c.PartitionDefs {
					if pd != nil && rangeOrList && pd.Values == nil {
						b.fail(rng, nil, review.UnsupportedForm, "Partition "+review.Code(ident(pd.Name))+" requires a "+review.Code("VALUES")+" definition.") // 1479
						return true
					}
				}
			case ast.ATCoalescePartition:
				if rangeOrList {
					b.fail(rng, nil, review.UnsupportedForm, review.Code("COALESCE PARTITION")+" is only valid for "+review.Code("HASH")+" or "+review.Code("KEY")+" partitioning.") // 1509
					return true
				}
			case ast.ATDropPartition:
				if pi.count > 0 && len(c.PartitionNames) >= pi.count {
					b.fail(rng, nil, review.UnsupportedForm, "Cannot remove all partitions from "+review.Code(r.written)+".") // 1508
					return true
				}
			}
		}
	case ast.ATPartitionBy:
		if c.PartitionBy == nil {
			return true
		}
		if b.fks && b.tableHasSurvivingFK(rel, d.droppedFKs) {
			b.fail(rng, nil, review.UnsupportedForm, "Partitioned table "+review.Code(r.written)+" cannot have foreign keys.") // 1506
			return true
		}
		var uqIdxs map[string]*indexMeta
		idxComplete := b.w.tableIndexesComplete[rel] && len(d.renamedFrom) == 0
		if idxComplete {
			uqIdxs = make(map[string]*indexMeta)
			for name, im := range b.w.tableIndexes[rel] {
				if !d.droppedIndexes[name] {
					uqIdxs[name] = im
				}
			}
			for i, im := range d.addedIndexes {
				uqIdxs["add:"+strconv.Itoa(i)] = im
			}
		}
		if b.partitionClauseColumnsAndUniqueFails(r, c.PartitionBy, rng, func(colName string) walk.Fact {
			st, _, _ := b.alterFinalCol(r, d, colName)
			return st
		}, uqIdxs, idxComplete) {
			return true
		}
	default:
		return alterNoFacts(c)
	}
	return true
}

// alterTextFails fails an ALTER TABLE its text alone refutes: a column
// added or dropped twice (1060, 1091), an index added or dropped twice
// (1061, 1091), two primary keys (1068), an index named PRIMARY (1280), a
// name too long (1059), or invalid column/option/partition/check definitions.
func (b *builder) alterTextFails(v *ast.AlterTableStmt, r tableRef, d *alterDelta) bool {
	var tableOpts []*ast.TableOption
	hasColOrExchange := false
	hasStandaloneValidation := false
	var validationRng review.Range
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATTableOption:
			if c.Option != nil {
				tableOpts = append(tableOpts, c.Option)
			}
		case ast.ATAddColumn, ast.ATDropColumn, ast.ATModifyColumn, ast.ATChangeColumn, ast.ATExchangePartition:
			hasColOrExchange = true
		case ast.ATWithValidation, ast.ATWithoutValidation:
			hasStandaloneValidation = true
			validationRng = rangeOf(c.Loc)
		}
	}
	if hasStandaloneValidation && !hasColOrExchange {
		b.fail(validationRng, nil, review.UnsupportedForm, review.Code("WITH/WITHOUT VALIDATION")+" requires a column modification or "+review.Code("EXCHANGE PARTITION")+".") // 1221
		return true
	}
	if b.tableOptionsFails(tableOpts, d.finalEngine, d.engineKnown, b.anchor()) {
		return true
	}
	seen := make(map[[2]string]bool)
	twice := func(kind, space, name string, rng review.Range, message string) bool {
		k, ok := b.key(space, name)
		if !ok {
			return false
		}
		if seen[[2]string{kind, k}] {
			b.fail(rng, nil, review.DuplicateName, message)
			return true
		}
		seen[[2]string{kind, k}] = true
		return false
	}
	primaries := 0
	for _, c := range v.Commands {
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ATAddColumn:
			cols := c.Columns
			if c.Column != nil {
				cols = append([]*ast.ColumnDef{c.Column}, cols...)
			}
			for _, col := range cols {
				if b.columnDefTextFails(col) {
					return true
				}
				if b.tooLong(col.Name, rangeOf(col.Loc)) ||
					twice("add column", spaceColumn, col.Name, rangeOf(col.Loc), statesTwice("adds", "column", ident(col.Name))) {
					return true
				}
				if colHasAutoIncrement(col) && !colHasInlineKey(col) {
					keyed := false
					for _, con := range c.Constraints {
						if con != nil && con.Type != ast.ConstrCheck && keyPartsCoverColumn(constraintKeyParts(con), col.Name) {
							keyed = true
							break
						}
					}
					if !keyed {
						for _, other := range v.Commands {
							if other == nil {
								continue
							}
							if (other.Type == ast.ATAddIndex || other.Type == ast.ATAddConstraint) &&
								other.Constraint != nil && other.Constraint.Type != ast.ConstrCheck &&
								keyPartsCoverColumn(constraintKeyParts(other.Constraint), col.Name) {
								keyed = true
								break
							}
							if other.Type == ast.ATAddColumn {
								for _, con := range other.Constraints {
									if con != nil && con.Type != ast.ConstrCheck && keyPartsCoverColumn(constraintKeyParts(con), col.Name) {
										keyed = true
										break
									}
								}
								if keyed {
									break
								}
							}
						}
					}
					if !keyed {
						b.fail(rangeOf(col.Loc), nil, review.UnsupportedForm, "Column "+review.Code(ident(col.Name))+" is "+review.Code("AUTO_INCREMENT")+" and must be defined as a key.") // 1075
						return true
					}
				}
				for _, cc := range col.Constraints {
					if cc.Type == ast.ColConstrPrimaryKey {
						primaries++
					}
				}
			}
			for _, con := range c.Constraints {
				if con != nil && con.Type == ast.ConstrCheck && con.Expr != nil && !checkExprIsBoolean(con.Expr) {
					b.fail(rangeOf(con.Loc), nil, review.UnsupportedForm, "Check constraint in "+review.Code(r.written)+" must be a boolean expression.") // 3812
					return true
				}
			}
		case ast.ATDropColumn:
			if twice("drop column", spaceColumn, c.Name, rng, statesTwice("drops", "column", ident(c.Name))) {
				return true
			}
		case ast.ATDropIndex:
			if twice("drop index", spaceIndex, c.Name, rng, statesTwice("drops", "index", ident(c.Name))) {
				return true
			}
		case ast.ATDropConstraint:
			switch dropKindOf(c) {
			case dropPrimaryKey:
				if twice("drop index", spaceIndex, "PRIMARY", rng, "This statement drops the primary key twice.") {
					return true
				}
			case dropForeignKey:
				if b.fks && twice("drop foreign key", spaceForeignKey, c.Name, rng, statesTwice("drops", "foreign key", ident(c.Name))) {
					return true
				}
			}
		case ast.ATChangeColumn, ast.ATRenameColumn, ast.ATModifyColumn:
			// A column a CHANGE or RENAME names is one more the statement
			// adds (1060, verified).
			to := c.NewName
			if c.Column != nil {
				if b.columnDefTextFails(c.Column) {
					return true
				}
				to = c.Column.Name
				for _, cc := range c.Column.Constraints {
					if cc.Type == ast.ColConstrPrimaryKey {
						primaries++
					}
				}
			}
			from, ok1 := b.key(spaceColumn, c.Name)
			dest, ok2 := b.key(spaceColumn, to)
			if c.Type != ast.ATModifyColumn && ok1 && ok2 && from != dest &&
				twice("add column", spaceColumn, to, rng, statesTwice("adds", "column", ident(to))) {
				return true
			}
		case ast.ATAddConstraint, ast.ATAddIndex:
			con := c.Constraint
			if con == nil {
				continue
			}
			if name := b.constraintName(con); b.tooLong(name, rng) {
				return true
			}
			switch con.Type {
			case ast.ConstrPrimaryKey:
				primaries++
			case ast.ConstrForeignKey:
				// 1826, or 1061 for the index each would add (verified).
				if b.fks && twice("add foreign key", spaceForeignKey, con.Name, rng, statesTwice("adds", "foreign key", ident(con.Name))) {
					return true
				}
			case ast.ConstrCheck:
				if con.Expr != nil && !checkExprIsBoolean(con.Expr) {
					b.fail(rng, nil, review.UnsupportedForm, "Check constraint in "+review.Code(r.written)+" must be a boolean expression.") // 3812
					return true
				}
			default:
				if k, ok := b.key(spaceIndex, con.Name); ok && k == primaryIndex {
					b.fail(rng, nil, review.ReservedName, reservedPrimary)
					return true
				}
				if twice("add index", spaceIndex, con.Name, rng, statesTwice("adds", "index", ident(con.Name))) {
					return true
				}
			}
		case ast.ATAddPartition, ast.ATReorganizePartition:
			if b.partitionDefsBoundFails(c.PartitionDefs, rng) {
				return true
			}
		case ast.ATPartitionBy:
			if b.partitionClausePreFKFails(c.PartitionBy, rng) {
				return true
			}
		}
	}
	if primaries > 1 {
		b.fail(b.anchor(), nil, review.NameTaken, hasPrimaryKey(r.written))
		return true
	}
	return false
}

// addColumn follows ADD COLUMN: the column must be new (1060), and a
// column AFTER names must exist (1054).
func (b *builder) addColumn(r tableRef, col *ast.ColumnDef, after string, rng review.Range) bool {
	k, ok := b.key(spaceColumn, col.Name)
	if !ok {
		return false
	}
	ck := colKey(r.db, r.name, k)
	b.need(ck, walk.Absent, nil, rng, "", alreadyExists("column", ident(col.Name), r.written)) // 1060
	b.after(r, after, rng)
	b.Set(ck, walk.Present, "")
	b.alterColumnKeys(r, col, rng)
	return true
}

// after checks the column a FIRST/AFTER clause names (1054).
func (b *builder) after(r tableRef, after string, rng review.Range) {
	if k, ok := b.key(spaceColumn, after); ok {
		b.need(colKey(r.db, r.name, k), walk.Present, nil, rng, notFound("column", ident(after), r.written), "") // 1054
	}
}

// columnNeeds checks that a column the subcommand names exists (1054).
func (b *builder) columnNeeds(r tableRef, name string, rng review.Range) bool {
	k, ok := b.key(spaceColumn, name)
	if !ok {
		return false
	}
	b.need(colKey(r.db, r.name, k), walk.Present, nil, rng, notFound("column", ident(name), r.written), "") // 1054
	return true
}

// renameColumn follows CHANGE a b and RENAME COLUMN a TO b.
func (b *builder) renameColumn(r tableRef, from, to string, rng review.Range, d *alterDelta, alwaysCheckPart bool) bool {
	fk, ok1 := b.key(spaceColumn, from)
	tk, ok2 := b.key(spaceColumn, to)
	if !ok1 || !ok2 {
		return false
	}
	if b.tooLong(to, rng) {
		return true
	}
	fc, tc := colKey(r.db, r.name, fk), colKey(r.db, r.name, tk)
	b.need(fc, walk.Present, nil, rng, notFound("column", ident(from), r.written), "") // 1054
	if f, _ := b.l.Lookup(fc); f == walk.Present && (alwaysCheckPart || fk != tk) {
		if !d.removesPartitioning && d.repartition == nil {
			if pi := b.w.partitionState[relKey(r.db, r.name)]; pi != nil && pi.partitioned && pi.cols[fk] {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(from))+" in "+review.Code(r.written)+" is used by the table's partitioning function.") // 3855
				return true
			}
		}
	}
	if fk != tk {
		b.need(tc, walk.Absent, nil, rng, "", alreadyExists("column", ident(to), r.written)) // 1060
		b.leave(fc)
		b.arrive(tc)
	}
	return true
}

// alterColumnKeys records the keys a column definition of ALTER TABLE
// declares inline: a primary key the table must not have yet (1068).
func (b *builder) alterColumnKeys(r tableRef, col *ast.ColumnDef, rng review.Range) {
	for _, cc := range col.Constraints {
		if cc.Type == ast.ColConstrPrimaryKey {
			b.need(idxKey(r.db, r.name, primaryIndex), walk.Absent, nil, rng, "", hasPrimaryKey(r.written)) // 1068
		}
	}
	b.columnKeys(r, col)
}

// addConstraint follows ADD of a key, index, check, or foreign key.
func (b *builder) addConstraint(r tableRef, c *ast.Constraint, d *alterDelta) {
	rng := rangeOf(c.Loc)
	if rng == (review.Range{}) {
		rng = b.anchor()
	}
	switch c.Type {
	case ast.ConstrPrimaryKey:
		b.need(idxKey(r.db, r.name, primaryIndex), walk.Absent, nil, rng, "", hasPrimaryKey(r.written)) // 1068
	case ast.ConstrCheck:
		b.constraint(r, c)
		return
	case ast.ConstrForeignKey:
		parts := constraintKeyParts(c)
		if b.checkIndexParts(r, parts, false, false, true, d.finalPartitioned, rng, func(colName string) (walk.Fact, string, bool) {
			return b.alterFinalCol(r, d, colName)
		}) {
			return
		}
		if !b.fks {
			break
		}
		if len(c.RefColumns) == 0 || (len(parts) > 0 && len(parts) != len(c.RefColumns)) {
			b.fail(rng, nil, review.UnsupportedForm, "Foreign key in "+review.Code(r.written)+" must name matching referenced columns.") // 1239
			return
		}
		if d.finalPartitioned {
			b.fail(rng, nil, review.UnsupportedForm, "Partitioned table "+review.Code(r.written)+" cannot define foreign keys.") // 1506
			return
		}
		if k, ok := b.key(spaceForeignKey, c.Name); ok {
			b.need(objKey(spaceForeignKey, r.db, k), walk.Absent, nil, rng, "", alreadyExists("foreign key", ident(c.Name), "")) // 1826; 1022 on 5.7
		}
		b.foreignKeyTarget(r, c)
	default:
		if k, ok := b.key(spaceIndex, c.Name); ok {
			ik := idxKey(r.db, r.name, k)
			if !b.isFKGenerated(ik) {
				b.need(ik, walk.Absent, nil, rng, "", alreadyExists("index", ident(c.Name), r.written)) // 1061
			}
		}
	}
	if c.Type != ast.ConstrForeignKey && c.Type != ast.ConstrCheck {
		parts := constraintKeyParts(c)
		fulltext := c.Type == ast.ConstrFulltextIndex
		spatial := c.Type == ast.ConstrSpatialIndex || isSpatialIndexType(c.IndexType)
		if b.checkIndexParts(r, parts, fulltext, spatial, false, d.finalPartitioned, rng, func(colName string) (walk.Fact, string, bool) {
			return b.alterFinalCol(r, d, colName)
		}) {
			return
		}
	}
	b.constraint(r, c)
}

// dropIndexName follows DROP INDEX, of the primary key when the name is
// PRIMARY.
func (b *builder) dropIndexName(r tableRef, name string, rng review.Range) bool {
	k, ok := b.key(spaceIndex, name)
	if !ok {
		return false
	}
	ik := idxKey(r.db, r.name, k)
	if k == primaryIndex {
		b.need(ik, walk.Present, nil, rng, noPrimaryKey(r.written), "") // 1091
	} else {
		b.need(ik, walk.Present, nil, rng, notFound("index", ident(name), r.written), "") // 1091
	}
	b.leave(ik)
	b.consumeFKGenerated(ik)
	return true
}

// alterRename follows ALTER TABLE ... RENAME [TO] t2.
func (b *builder) alterRename(r tableRef, c *ast.AlterTableCmd) bool {
	dest := c.NewTable
	if dest == nil {
		dest = &ast.TableRef{Name: c.NewName, Loc: c.Loc}
	}
	to, st := b.ref(dest)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	if b.tooLong(dest.Name, to.rng) {
		return true
	}
	if to.db == r.db && to.name == r.name {
		return true
	}
	if dest.Schema != "" {
		b.databaseNeeds(to.db, dest.Schema, to.rng) // 1049; 1025 on 5.7
	}
	b.needFree(relKey(to.db, to.name), to.rng, to.written, relTable) // 1050
	b.moveTable(r, to, relTable)
	return true
}

// moveTable records a relation renamed: what it holds goes with it, and
// so, to another database, do its foreign keys and triggers, under names
// the walk does not list.
func (b *builder) moveTable(from, to tableRef, kind string) {
	toPath := tablePath(to.db, to.name)
	b.Move(tablePath(from.db, from.name), toPath)
	b.Set(relKey(from.db, from.name), walk.Absent, "")
	b.Set(relKey(to.db, to.name), walk.Present, kind)
	if from.db != to.db {
		for _, space := range []string{spaceForeignKey, spaceTrigger} {
			b.MayRemove(walk.Scope{Space: space, Parent: walk.Path(from.db)})
			b.MayAdd(walk.Scope{Space: space, Parent: walk.Path(to.db)})
		}
	}
	// The foreign keys named <from>_ibfk_... become <to>_ibfk_...
	// (verified on 8.0 and 8.4).
	old, renamed := b.generatedFKs(from.name), b.generatedFKs(to.name)
	b.MayRemove(walk.Scope{Space: spaceForeignKey, Parent: walk.Path(from.db), Prefix: old})
	b.MayAdd(walk.Scope{Space: spaceForeignKey, Parent: walk.Path(to.db), Prefix: renamed})
	b.MayRemove(walk.Scope{Space: spaceTableForeignKey, Parent: toPath, Prefix: old})
	b.MayAdd(walk.Scope{Space: spaceTableForeignKey, Parent: toPath, Prefix: renamed})
	fk, tk := relKey(from.db, from.name), relKey(to.db, to.name)
	fromPath := tablePath(from.db, from.name)
	b.Session(func() {
		if b.w.triggered[fk] {
			b.w.triggered[tk] = true
		}
		ignores := b.w.ignoresFKs[fk]
		delete(b.w.ignoresFKs, fk)
		if ignores {
			b.w.ignoresFKs[tk] = true
		} else {
			delete(b.w.ignoresFKs, tk)
		}
		type movedGenEntry struct {
			name   string
			owners map[string]bool
		}
		var movedGen []movedGenEntry
		for gk, owners := range b.w.fkGeneratedIndexes {
			if gk.Parent == fromPath {
				movedGen = append(movedGen, movedGenEntry{name: gk.Name, owners: owners})
				delete(b.w.fkGeneratedIndexes, gk)
			}
		}
		for _, entry := range movedGen {
			rewritten := make(map[string]bool, len(entry.owners))
			for fkName := range entry.owners {
				if old != "" && renamed != "" && strings.HasPrefix(fkName, old) {
					fkName = renamed + strings.TrimPrefix(fkName, old)
				}
				rewritten[fkName] = true
			}
			b.w.fkGeneratedIndexes[idxKey(to.db, to.name, entry.name)] = rewritten
		}
		if eng, ok := b.w.engines[fk]; ok {
			delete(b.w.engines, fk)
			b.w.engines[tk] = eng
		}
		if pi, ok := b.w.partitionState[fk]; ok {
			delete(b.w.partitionState, fk)
			b.w.partitionState[tk] = pi
		}
		if idxs, ok := b.w.tableIndexes[fk]; ok {
			delete(b.w.tableIndexes, fk)
			b.w.tableIndexes[tk] = idxs
		}
		if comp, ok := b.w.tableIndexesComplete[fk]; ok {
			delete(b.w.tableIndexesComplete, fk)
			b.w.tableIndexesComplete[tk] = comp
		}
		if fks, ok := b.w.tableFKs[fk]; ok {
			delete(b.w.tableFKs, fk)
			rewrittenFKs := make(map[string][]string, len(fks))
			for fkName, cols := range fks {
				if old != "" && renamed != "" && strings.HasPrefix(fkName, old) {
					fkName = renamed + strings.TrimPrefix(fkName, old)
				}
				rewrittenFKs[fkName] = cols
			}
			b.w.tableFKs[tk] = rewrittenFKs
		}
		if refs, ok := b.w.referencedByFK[fk]; ok {
			delete(b.w.referencedByFK, fk)
			b.w.referencedByFK[tk] = refs
		}
		movedCols := make(map[string]string)
		for ck, bt := range b.w.colTypes {
			if ck.Parent == fromPath {
				movedCols[ck.Name] = bt
				delete(b.w.colTypes, ck)
			}
		}
		for colName, bt := range movedCols {
			b.w.colTypes[colKey(to.db, to.name, colName)] = bt
		}
	})
}

// renameTable follows RENAME TABLE, whose pairs run in order: a swap
// through a third name succeeds (verified). Each pair is a part, and a
// pair's check on a name another pair changes is left out.
func (b *builder) renameTable(v *ast.RenameTableStmt) bool {
	moved := make(map[walk.Key]string)
	for _, p := range v.Pairs {
		b.Next()
		from, st1 := b.ref(p.Old)
		to, st2 := b.ref(p.New)
		switch {
		case st1 == refFailed || st2 == refFailed:
			return true
		case st1 == refUnknown || st2 == refUnknown:
			return false
		}
		if b.tooLong(p.New.Name, to.rng) {
			return true
		}
		fk := relKey(from.db, from.name)
		f, temp, kind := b.resolve(from)
		movedHere := false
		if k, ok := moved[fk]; ok {
			f, temp, kind = walk.Present, false, k
			movedHere = true
		}
		if _, tt, _ := b.resolve(to); temp || tt || f == walk.Unknown {
			return false
		}
		if from.db == to.db && from.name == to.name {
			if f == walk.Absent {
				b.need(fk, walk.Present, nil, from.rng, notFound("", from.written, ""), "") // 1146
				return true
			}
			if movedHere {
				b.fail(to.rng, nil, review.NameTaken, alreadyExists(kind, to.written, "")) // 1050
				return true
			}
			b.needFree(fk, to.rng, to.written, kind) // 1050
			return true
		}
		b.need(fk, walk.Present, nil, from.rng, notFound("", from.written, ""), "") // 1146; 1017 on 5.7
		if p.New.Schema != "" {
			b.databaseNeeds(to.db, p.New.Schema, to.rng) // 1049; 1025 on 5.7
		}
		if kind == relView && from.db != to.db {
			b.fail(to.rng, []walk.Key{fk}, review.UnsupportedForm, engineName+" doesn't accept "+review.Code("RENAME TABLE")+" of view "+review.Code(from.written)+" to another database.") // 1450 (verified)
			return true
		}
		b.needFree(relKey(to.db, to.name), to.rng, to.written, kind) // 1050
		b.moveTable(from, to, kind)
		delete(moved, fk)
		moved[relKey(to.db, to.name)] = kind
	}
	return true
}

// dropTable follows DROP [TEMPORARY] TABLE [IF EXISTS], which looks every
// table up before it drops any (walk.Lookups; 8.0 and 8.4 are atomic, 5.7
// drops the earlier tables first, and the statement fails either way).
// Without TEMPORARY a temporary table of the name is dropped instead of
// the base table.
func (b *builder) dropTable(v *ast.DropTableStmt) bool {
	b.Order = walk.Lookups
	for _, t := range v.Tables {
		b.Next()
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			return false
		}
		if b.tooLong(t.Name, r.rng) {
			return true
		}
		tk := tempKey(r.db, r.name)
		tf, _ := b.l.Lookup(tk)
		if v.Temporary {
			if !v.IfExists {
				b.need(tk, walk.Present, nil, r.rng, notFound("temporary table", r.written, ""), "") // 1051
			}
			b.Set(tk, walk.Absent, "")
			continue
		}
		path := walk.Scope{Parent: tablePath(r.db, r.name), Deep: true}
		k := relKey(r.db, r.name)
		switch tf {
		case walk.Present:
			b.Set(tk, walk.Absent, "")
			continue
		case walk.Unknown:
			b.Set(tk, walk.Unknown, "")
			b.Set(k, walk.Unknown, "")
			b.Forget(path)
			b.droppedOwned(r)
			continue
		}
		f, kind := b.l.Lookup(k)
		if !v.IfExists {
			// A view is no table to DROP TABLE (1051).
			b.needKind(k, []string{relTable}, r.rng, notFound("table", r.written, ""), r.written, "table")
		}
		switch {
		case f == walk.Present && kind == relView:
		case f == walk.Present && kind == "", f == walk.Unknown && v.IfExists:
			// It may be a view, which DROP TABLE IF EXISTS leaves with a
			// note (verified).
			b.Set(k, walk.Unknown, "")
			b.Forget(path)
			b.droppedOwned(r)
		default:
			b.Set(k, walk.Absent, "")
			b.Clear(path)
			b.droppedOwned(r)
			b.clearTableFKGenerated(r.db, r.name)
		}
	}
	return true
}

// droppedOwned records that a dropped table's foreign keys and triggers
// went with it, under names the walk does not list.
func (b *builder) droppedOwned(r tableRef) {
	b.MayRemove(walk.Scope{Space: spaceForeignKey, Parent: walk.Path(r.db)})
	b.MayRemove(walk.Scope{Space: spaceTrigger, Parent: walk.Path(r.db)})
}

// truncate follows TRUNCATE: each table must exist and be a table (1146,
// a view too).
func (b *builder) truncate(v *ast.TruncateStmt) bool {
	for _, t := range v.Tables {
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			continue
		}
		if f, temp, _ := b.resolve(r); !temp && f != walk.Unknown {
			b.needKind(relKey(r.db, r.name), []string{relTable}, r.rng, notFound("table", r.written, ""), r.written, "table")
		}
	}
	return true
}
