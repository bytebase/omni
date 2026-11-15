package catalog

import (
	"fmt"
	"strconv"
	"strings"
)

func (g *generator) addOp(typ MigrationOpType, db *Database, name, sql string) {
	g.plan.Ops = append(g.plan.Ops, MigrationOp{Type: typ, Database: databaseName(db), ObjectName: name, SQL: sql})
}

// dropViews drops the views the diff drops. TiDB does not check the views
// that read a view when it is dropped, so they go in any order.
func (g *generator) dropViews(diff *SchemaDiff) {
	for _, e := range diff.Views {
		if e.Action == DiffDrop {
			g.addOp(OpDropView, e.From.Database, e.From.Name, "DROP VIEW "+qualifiedIdent(e.From.Database, e.From.Name))
		}
	}
}

// views creates the new views and replaces the changed ones, each after the
// views of the plan its body names.
func (g *generator) views(diff *SchemaDiff) {
	var pending []ViewDiffEntry
	for _, e := range diff.Views {
		if e.Action != DiffDrop {
			pending = append(pending, e)
		}
	}
	for len(pending) > 0 {
		// The first view that names no other pending view goes next; in a
		// cycle, which TiDB cannot create either, the first one does.
		next := 0
		for i, e := range pending {
			if !namesAnyOf(e.To, pending, i) {
				next = i
				break
			}
		}
		e := pending[next]
		pending = append(pending[:next], pending[next+1:]...)
		g.addOp(OpCreateView, e.To.Database, e.To.Name, createViewSQL(e.To, e.Action == DiffModify))
	}
}

// namesAnyOf reports whether v's body reads another pending view. The
// catalog writes every relation a body reads with its database, `db`.`v`,
// and no alias or column with one, so the qualified name tells a view apart
// from an alias or column spelled like it.
func namesAnyOf(v *View, pending []ViewDiffEntry, self int) bool {
	for i, e := range pending {
		if i != self && strings.Contains(v.Definition, qualifiedIdent(e.To.Database, e.To.Name)) {
			return true
		}
	}
	return false
}

// createViewSQL is SHOW CREATE VIEW without the definer, which is an
// environment's account rather than schema, with the name qualified and the
// CHECK OPTION, which SHOW CREATE VIEW leaves out.
func createViewSQL(v *View, replace bool) string {
	definer := v.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}
	text := strings.Replace(showCreateView(v), " DEFINER="+formatDefiner(definer), "", 1)
	text = strings.Replace(text, " VIEW "+backtickIDQ(v.Name)+" (", " VIEW "+qualifiedIdent(v.Database, v.Name)+" (", 1)
	if replace {
		text = "CREATE OR REPLACE " + strings.TrimPrefix(text, "CREATE ")
	}
	if v.CheckOption != "" {
		text += " WITH " + strings.ToUpper(v.CheckOption) + " CHECK OPTION"
	}
	return text
}

// sequences drops the sequences the diff drops, then creates the new ones and
// changes the changed ones. ALTER SEQUENCE cannot change a comment, so a
// sequence whose comment changes is dropped and created again, which
// restarts it.
func (g *generator) sequences(diff *SchemaDiff) {
	for _, e := range diff.Sequences {
		if e.Action == DiffDrop || e.Action == DiffModify && e.From.Comment != e.To.Comment {
			g.addOp(OpDropSequence, e.From.Database, e.From.Name, "DROP SEQUENCE "+qualifiedIdent(e.From.Database, e.From.Name))
		}
	}
	for _, e := range diff.Sequences {
		switch {
		case e.Action == DiffAdd || e.Action == DiffModify && e.From.Comment != e.To.Comment:
			text := strings.Replace(showCreateSequence(e.To), " "+quoteIdent(e.To.Name)+" ", " "+qualifiedIdent(e.To.Database, e.To.Name)+" ", 1)
			g.addOp(OpCreateSequence, e.To.Database, e.To.Name, text)
		case e.Action == DiffModify:
			g.addOp(OpAlterSequence, e.To.Database, e.To.Name, "ALTER SEQUENCE "+qualifiedIdent(e.To.Database, e.To.Name)+" "+sequenceChanges(e.From, e.To))
		}
	}
}

// sequenceChanges lists the options that differ, all in one statement, since
// TiDB checks the sequence they leave and not each step.
func sequenceChanges(a, b *Sequence) string {
	var opts []string
	add := func(changed bool, opt string) {
		if changed {
			opts = append(opts, opt)
		}
	}
	add(a.Start != b.Start, "START WITH "+strconv.FormatInt(b.Start, 10))
	add(a.MinValue != b.MinValue, "MINVALUE "+strconv.FormatInt(b.MinValue, 10))
	add(a.MaxValue != b.MaxValue, "MAXVALUE "+strconv.FormatInt(b.MaxValue, 10))
	add(a.Increment != b.Increment, "INCREMENT BY "+strconv.FormatInt(b.Increment, 10))
	if a.Cache != b.Cache {
		if b.Cache == 0 {
			opts = append(opts, "NOCACHE")
		} else {
			opts = append(opts, fmt.Sprintf("CACHE %d", b.Cache))
		}
	}
	if a.Cycle != b.Cycle {
		if b.Cycle {
			opts = append(opts, "CYCLE")
		} else {
			opts = append(opts, "NOCYCLE")
		}
	}
	return strings.Join(opts, " ")
}
