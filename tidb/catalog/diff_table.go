package catalog

import (
	"strconv"
	"strings"
)

// diffTables compares the base tables of two catalogs. Temporary tables are
// left out: they last for the session and are no part of a schema.
func diffTables(from, to *Catalog) []TableDiffEntry {
	return diffObjects(tablesOf(from), tablesOf(to),
		func(a, b *Table) bool {
			e := compareTable(a, b)
			return tableChanged(&e)
		},
		func(action DiffAction, key objectKey, f, t *Table) TableDiffEntry {
			e := TableDiffEntry{Action: action, Database: key.db, From: f, To: t}
			if t != nil {
				e.Name = t.Name
			} else {
				e.Name = f.Name
			}
			if action == DiffModify {
				e = compareTable(f, t)
				e.Database = key.db
			}
			return e
		},
		func(e TableDiffEntry) (string, string, DiffAction) { return e.Database, toLower(e.Name), e.Action })
}

func tablesOf(c *Catalog) map[objectKey]*Table {
	m := make(map[objectKey]*Table)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		for name, tbl := range db.Tables {
			if tbl == nil || tbl.Temporary {
				continue
			}
			m[objectKey{db: toLower(db.Name), name: name}] = tbl
		}
	}
	return m
}

// compareTable returns the DiffModify entry of two same-identity tables, with
// every sub-object diff filled in.
func compareTable(from, to *Table) TableDiffEntry {
	return TableDiffEntry{
		Action:           DiffModify,
		Name:             to.Name,
		From:             from,
		To:               to,
		Columns:          diffColumns(from, to),
		Indexes:          diffIndexes(from, to),
		ForeignKeys:      diffForeignKeys(from, to),
		Checks:           diffChecks(from, to),
		PartitionChanged: canonicalPartitioning(from) != canonicalPartitioning(to),
		OptionsChanged:   canonicalTableOptions(from) != canonicalTableOptions(to),
	}
}

func tableChanged(e *TableDiffEntry) bool {
	return len(e.Columns) > 0 || len(e.Indexes) > 0 || len(e.ForeignKeys) > 0 ||
		len(e.Checks) > 0 || e.PartitionChanged || e.OptionsChanged
}

// tableCharsetCollation is the table's charset and collation as SHOW CREATE
// TABLE resolves them.
func tableCharsetCollation(tbl *Table) (charset, collation string) {
	charset = tbl.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	collation = tbl.Collation
	if collation == "" {
		collation = defaultCollationForCharset[toLower(charset)]
	}
	return charset, collation
}

// canonicalTableOptions is the option list SHOW CREATE TABLE ends a table
// with, less the AUTO_INCREMENT and AUTO_RANDOM_BASE counters, which are data
// rather than schema, then TTL, a global temporary table's ON COMMIT, and the
// TiFlash replicas, which SHOW CREATE TABLE leaves out.
func canonicalTableOptions(tbl *Table) string {
	withoutCounters := *tbl
	withoutCounters.AutoIncrement = 0
	withoutCounters.AutoRandomBase = 0
	charset, collation := tableCharsetCollation(tbl)
	var b strings.Builder
	b.WriteString(showTableOptions(&withoutCounters, charset, collation))
	if tbl.TTLColumn != "" {
		b.WriteString(showTTL(tbl))
	}
	if tbl.GlobalTemporary {
		b.WriteString(" GLOBAL TEMPORARY ON COMMIT " + tbl.OnCommit)
	}
	b.WriteString(tiflashReplica(tbl))
	return b.String()
}

// tiflashReplica is the table's SET TIFLASH REPLICA clause, "" for none.
func tiflashReplica(tbl *Table) string {
	if tbl.TiFlashReplica == 0 && len(tbl.TiFlashLocationLabels) == 0 {
		return ""
	}
	s := " SET TIFLASH REPLICA " + strconv.Itoa(tbl.TiFlashReplica)
	if len(tbl.TiFlashLocationLabels) > 0 {
		labels := make([]string, len(tbl.TiFlashLocationLabels))
		for i, l := range tbl.TiFlashLocationLabels {
			labels[i] = "'" + outputFormat(l) + "'"
		}
		s += " LOCATION LABELS " + strings.Join(labels, ", ")
	}
	return s
}

// canonicalPartitioning is the partition clause SHOW CREATE TABLE writes, ""
// for a table without one.
func canonicalPartitioning(tbl *Table) string {
	if tbl.Partitioning == nil {
		return ""
	}
	return showPartitioning(tbl.Partitioning)
}

// diffColumns compares the columns of two versions of a table by the line
// SHOW CREATE TABLE writes for each. Position is not compared: a reorder is no
// change to a column.
func diffColumns(from, to *Table) []ColumnDiffEntry {
	return diffNamed(columnsOf(from), columnsOf(to),
		func(a, b *Column) bool { return canonicalColumn(a) != canonicalColumn(b) },
		func(action DiffAction, f, t *Column) ColumnDiffEntry {
			e := ColumnDiffEntry{Action: action, From: f, To: t}
			if t != nil {
				e.Name = t.Name
			} else {
				e.Name = f.Name
			}
			return e
		},
		func(e ColumnDiffEntry) (string, DiffAction) { return e.Name, e.Action })
}

func columnsOf(t *Table) map[string]*Column {
	m := make(map[string]*Column, len(t.Columns))
	for _, col := range t.Columns {
		m[toLower(col.Name)] = col
	}
	return m
}

// canonicalColumn is a column as SHOW CREATE TABLE writes it, with its
// character set and collation always written: a table's default changes the
// default of columns added later and leaves existing ones as they were.
func canonicalColumn(col *Column) string {
	return showColumn(col, "", "")
}

// diffIndexes compares the keys of two versions of a table by the line SHOW
// CREATE TABLE writes for each, the primary key named PRIMARY. A key that
// backs a foreign key is compared like any other: TiDB keeps it when the
// foreign key is dropped.
func diffIndexes(from, to *Table) []IndexDiffEntry {
	return diffNamed(indexesOf(from), indexesOf(to),
		func(a, b *Index) bool { return showIndex(from, a) != showIndex(to, b) },
		func(action DiffAction, f, t *Index) IndexDiffEntry {
			e := IndexDiffEntry{Action: action, From: f, To: t}
			if t != nil {
				e.Name = t.Name
			} else {
				e.Name = f.Name
			}
			return e
		},
		func(e IndexDiffEntry) (string, DiffAction) { return e.Name, e.Action })
}

func indexesOf(t *Table) map[string]*Index {
	m := make(map[string]*Index, len(t.Indexes))
	for _, idx := range t.Indexes {
		name := idx.Name
		if idx.Primary {
			name = "PRIMARY"
		}
		m[toLower(name)] = idx
	}
	return m
}

// diffForeignKeys compares the foreign keys of two versions of a table by the
// line SHOW CREATE TABLE writes for each.
func diffForeignKeys(from, to *Table) []ForeignKeyDiffEntry {
	return diffNamed(constraintsOf(from, ConForeignKey), constraintsOf(to, ConForeignKey),
		func(a, b *Constraint) bool { return showForeignKey(from, a) != showForeignKey(to, b) },
		func(action DiffAction, f, t *Constraint) ForeignKeyDiffEntry {
			e := ForeignKeyDiffEntry{Action: action, From: f, To: t}
			if t != nil {
				e.Name = t.Name
			} else {
				e.Name = f.Name
			}
			return e
		},
		func(e ForeignKeyDiffEntry) (string, DiffAction) { return e.Name, e.Action })
}

// diffChecks compares the CHECK constraints of two versions of a table by
// expression and ENFORCED state.
func diffChecks(from, to *Table) []CheckDiffEntry {
	return diffNamed(constraintsOf(from, ConCheck), constraintsOf(to, ConCheck),
		func(a, b *Constraint) bool { return canonicalCheck(a) != canonicalCheck(b) },
		func(action DiffAction, f, t *Constraint) CheckDiffEntry {
			e := CheckDiffEntry{Action: action, From: f, To: t}
			if t != nil {
				e.Name = t.Name
			} else {
				e.Name = f.Name
			}
			return e
		},
		func(e CheckDiffEntry) (string, DiffAction) { return e.Name, e.Action })
}

func constraintsOf(t *Table, typ ConstraintType) map[string]*Constraint {
	m := make(map[string]*Constraint)
	for _, con := range t.Constraints {
		if con.Type == typ && con.Name != "" {
			m[toLower(con.Name)] = con
		}
	}
	return m
}

// canonicalCheck is a check as SHOW CREATE TABLE writes it, less its name.
func canonicalCheck(con *Constraint) string {
	enforced := "enforced"
	if con.NotEnforced {
		enforced = "not enforced"
	}
	return encodeKeyFields("expr", stripOuterParens(con.CheckExpr), "state", enforced)
}
