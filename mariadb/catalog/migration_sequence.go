package catalog

import (
	"fmt"
	"math/big"
	"strings"
)

// Sequences are created ahead of tables, since CREATE TABLE refuses a default
// that reads a sequence it cannot find (1146), and dropped after them.
const (
	prioritySequence     = priorityTable - 5
	prioritySequenceDrop = priorityTable + 1
)

func generateSequenceDDL(diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, e := range diff.Sequences {
		switch e.Action {
		case DiffAdd:
			ops = append(ops, MigrationOp{
				Type: OpCreateSequence, Database: e.Database, ObjectName: e.Name,
				SQL:      createSequenceSQL(e.To),
				Phase:    PhaseMain,
				Priority: prioritySequence,
				sortName: e.Database + "." + toLower(e.Name),
			})
		case DiffDrop:
			ops = append(ops, MigrationOp{
				Type: OpDropSequence, Database: e.Database, ObjectName: e.Name,
				SQL:      "DROP SEQUENCE " + sequenceIdent(e.From),
				Phase:    PhasePre,
				Priority: prioritySequenceDrop,
				sortName: e.Database + "." + toLower(e.Name),
			})
		case DiffModify:
			ops = append(ops, MigrationOp{
				Type: OpAlterSequence, Database: e.Database, ObjectName: e.Name,
				SQL:      alterSequenceSQL(e.From, e.To),
				Warning:  sequenceRangeWarning(e.From, e.To),
				Phase:    PhaseMain,
				Priority: prioritySequence,
				sortName: e.Database + "." + toLower(e.Name),
			})
		}
	}
	return ops
}

func sequenceIdent(seq *Sequence) string {
	if seq.Database != nil {
		return quoteIdent(seq.Database.Name) + "." + quoteIdent(seq.Name)
	}
	return quoteIdent(seq.Name)
}

// createSequenceSQL is SHOW CREATE SEQUENCE for the sequence, qualified with
// its database.
func createSequenceSQL(seq *Sequence) string {
	return strings.Replace(showCreateSequence(seq), "CREATE SEQUENCE "+quoteIdent(seq.Name), "CREATE SEQUENCE "+sequenceIdent(seq), 1)
}

// alterSequenceSQL changes a sequence in place, keeping its next value: ALTER
// SEQUENCE for the sequence options, and ALTER TABLE, which a sequence takes
// as the table it is, for ENGINE and COMMENT, which ALTER SEQUENCE has no
// syntax for.
func alterSequenceSQL(from, to *Sequence) string {
	var stmts []string
	var opts []string
	if from.DataType != to.DataType {
		opts = append(opts, "AS "+to.DataType)
	}
	for _, o := range []struct{ name, from, to string }{
		{"MINVALUE", from.MinValue, to.MinValue},
		{"MAXVALUE", from.MaxValue, to.MaxValue},
		{"START WITH", from.Start, to.Start},
		{"INCREMENT BY", from.Increment, to.Increment},
	} {
		if o.from != o.to {
			opts = append(opts, o.name+" "+o.to)
		}
	}
	if from.Cache != to.Cache {
		if to.Cache == "0" {
			opts = append(opts, "NOCACHE")
		} else {
			opts = append(opts, "CACHE "+to.Cache)
		}
	}
	if from.Cycle != to.Cycle {
		if to.Cycle {
			opts = append(opts, "CYCLE")
		} else {
			opts = append(opts, "NOCYCLE")
		}
	}
	if len(opts) > 0 {
		stmts = append(stmts, "ALTER SEQUENCE "+sequenceIdent(to)+" "+strings.Join(opts, " "))
	}
	var tableOpts []string
	if !strings.EqualFold(from.Engine, to.Engine) {
		tableOpts = append(tableOpts, "ENGINE="+to.Engine)
	}
	if from.Comment != to.Comment {
		tableOpts = append(tableOpts, "COMMENT="+quoteStringLiteral(to.Comment))
	}
	if len(tableOpts) > 0 {
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE %s %s", sequenceIdent(to), strings.Join(tableOpts, " ")))
	}
	return strings.Join(stmts, ";\n")
}

// sequenceRangeWarning warns where the new limits leave out the sequence's
// start. ALTER SEQUENCE refuses limits that leave out the next value (4085),
// and the next value is data: a plan that restarted the sequence could hand out
// values it has already issued, so it leaves the restart to the user.
func sequenceRangeWarning(from, to *Sequence) string {
	if !outsideRange(from.Start, to.MinValue, to.MaxValue) {
		return ""
	}
	return fmt.Sprintf("the new limits of sequence %s leave out its start %s: if its next value lies outside them, ALTER SEQUENCE fails (4085) until it is restarted within them", sequenceIdent(to), from.Start)
}

func outsideRange(value, minValue, maxValue string) bool {
	v, ok1 := new(big.Int).SetString(value, 10)
	lo, ok2 := new(big.Int).SetString(minValue, 10)
	hi, ok3 := new(big.Int).SetString(maxValue, 10)
	return ok1 && ok2 && ok3 && (v.Cmp(lo) < 0 || v.Cmp(hi) > 0)
}
