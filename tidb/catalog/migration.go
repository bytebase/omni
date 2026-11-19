package catalog

import (
	"strings"
)

// GenerateMigration turns a SchemaDiff into the DDL that brings the `from`
// catalog to the `to` one on TiDB. When the diff is empty the plan is too.
//
// Every statement makes one change. TiDB refuses many changes in one ALTER
// TABLE that MySQL takes: two clauses on one column, an index over a column
// the same statement adds, a second CHECK (8200), and before v6.2 any second
// clause at all. The order is:
//
//  1. foreign keys that go or change, those of dropped tables, and those
//     around a column whose type changes (3780 otherwise);
//  2. dropped views, then on each changed table its dropped checks and keys
//     (TiDB refuses to drop a column a composite key covers, 8200) and the
//     keys around a column an expression key reads (3837 otherwise), then
//     its dropped columns;
//  3. dropped tables and sequences;
//  4. new sequences, which a column default may read, and changed ones;
//  5. new tables without their foreign keys, then on each changed table its
//     options, new columns, changed columns, keys, checks and partitioning;
//  6. new and changed views, each after the views it reads, since TiDB
//     checks a view's body when it is created;
//  7. foreign keys, once every table they reference exists.
//
// Some changes TiDB cannot make in place: a clustered primary key, AUTO_RANDOM
// and AUTO_INCREMENT added or removed, a stored generated column, a column
// another generated column reads, a narrower charset. The plan still states
// them, and the op carries a Warning; TiDB refuses the statement.
func GenerateMigration(from, to *Catalog, diff *SchemaDiff) *MigrationPlan {
	plan := &MigrationPlan{}
	if diff == nil || diff.IsEmpty() {
		return plan
	}
	g := &generator{plan: plan}
	// TiDB refuses every DDL statement against a cached table (8242, measured
	// on v8.5.5), DROP TABLE included, so a table cached in the source is
	// uncached before anything else and one cached in the target is cached
	// again once the plan is done. A table cached on both sides therefore
	// brackets the rest of its changes.
	for i := range diff.Tables {
		if e := &diff.Tables[i]; e.From != nil && e.From.Cached {
			g.add(OpAlterTable, e.From, alterTable(e.From, "NOCACHE"), "")
		}
	}
	g.dropForeignKeys(diff, to)
	g.dropViews(diff)
	for i := range diff.Tables {
		if e := &diff.Tables[i]; e.Action == DiffModify {
			g.dropTableParts(e)
		}
	}
	for i := range diff.Tables {
		if e := &diff.Tables[i]; e.Action == DiffDrop {
			g.add(OpDropTable, e.From, "DROP TABLE "+tableIdent(e.From), "")
		}
	}
	g.sequences(diff)
	for i := range diff.Tables {
		if e := &diff.Tables[i]; e.Action == DiffAdd {
			g.createTable(e.To)
		}
	}
	for i := range diff.Tables {
		if e := &diff.Tables[i]; e.Action == DiffModify {
			g.alterTable(e)
		}
	}
	g.views(diff)
	g.addForeignKeys(diff)
	for i := range diff.Tables {
		if e := &diff.Tables[i]; e.To != nil && e.To.Cached {
			g.add(OpAlterTable, e.To, alterTable(e.To, "CACHE"), "")
		}
	}
	return plan
}

// MigrationOpType classifies one statement of a plan.
type MigrationOpType string

const (
	OpCreateTable    MigrationOpType = "CreateTable"
	OpDropTable      MigrationOpType = "DropTable"
	OpAlterTable     MigrationOpType = "AlterTable" // table options and partitioning
	OpAddColumn      MigrationOpType = "AddColumn"
	OpDropColumn     MigrationOpType = "DropColumn"
	OpModifyColumn   MigrationOpType = "ModifyColumn"
	OpAddIndex       MigrationOpType = "AddIndex"
	OpDropIndex      MigrationOpType = "DropIndex"
	OpAddForeignKey  MigrationOpType = "AddForeignKey"
	OpDropForeignKey MigrationOpType = "DropForeignKey"
	OpAddCheck       MigrationOpType = "AddCheck"
	OpDropCheck      MigrationOpType = "DropCheck"
	OpCreateView     MigrationOpType = "CreateView"
	OpDropView       MigrationOpType = "DropView"
	OpCreateSequence MigrationOpType = "CreateSequence"
	OpAlterSequence  MigrationOpType = "AlterSequence"
	OpDropSequence   MigrationOpType = "DropSequence"
)

// MigrationOp is one statement of a plan. ObjectName is the table, view or
// sequence it changes, and Database the database that holds it. Warning
// says why TiDB refuses the statement, when it does.
type MigrationOp struct {
	Type       MigrationOpType
	Database   string
	ObjectName string
	SQL        string
	Warning    string
}

// MigrationPlan is the ordered statements of a migration.
type MigrationPlan struct {
	Ops []MigrationOp
}

// SQL joins the statements with ";\n", or returns "" for an empty plan.
func (p *MigrationPlan) SQL() string {
	parts := make([]string, len(p.Ops))
	for i, op := range p.Ops {
		parts[i] = op.SQL
	}
	return strings.Join(parts, ";\n")
}

// Warnings returns the operations that carry a warning.
func (p *MigrationPlan) Warnings() []MigrationOp {
	var out []MigrationOp
	for _, op := range p.Ops {
		if op.Warning != "" {
			out = append(out, op)
		}
	}
	return out
}

type generator struct {
	plan *MigrationPlan
	// readdOrder holds the foreign keys of `to` dropped in step 1 only to
	// change a column under them, which step 7 adds back.
	readdOrder []*Constraint
}

// add appends a statement on a table or another object of a database.
func (g *generator) add(typ MigrationOpType, tbl *Table, sql, warning string) {
	g.plan.Ops = append(g.plan.Ops, MigrationOp{Type: typ, Database: databaseName(tbl.Database), ObjectName: tbl.Name, SQL: sql, Warning: warning})
}

func databaseName(db *Database) string {
	if db == nil {
		return ""
	}
	return db.Name
}

// tableIdent is the table's name qualified with its database, both as
// declared: TiDB compares table names case-sensitively unless
// lower_case_table_names says otherwise.
func tableIdent(tbl *Table) string {
	return qualifiedIdent(tbl.Database, tbl.Name)
}

func qualifiedIdent(db *Database, name string) string {
	if db == nil || db.Name == "" {
		return quoteIdent(name)
	}
	return quoteIdent(db.Name) + "." + quoteIdent(name)
}

func alterTable(tbl *Table, clause string) string {
	return "ALTER TABLE " + tableIdent(tbl) + " " + clause
}
