package catalog

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Sequence is a TiDB sequence. Sequences share the table namespace: CREATE
// TABLE, CREATE VIEW and CREATE SEQUENCE answer 1050 for a name either of the
// other two holds (read on v8.5.5).
type Sequence struct {
	Name      string
	Database  *Database
	Start     int64
	MinValue  int64
	MaxValue  int64
	Increment int64
	// Cache is the cache size; 0 is NOCACHE.
	Cache    int64
	Cycle    bool
	Comment  string
	SkipDump bool
}

// The defaults TiDB's handleSequenceOptions (pkg/ddl/sequence.go) fills in.
const (
	sequenceDefaultCache       = 1000
	sequencePositiveMaxDefault = math.MaxInt64 - 1
	sequenceNegativeMinDefault = math.MinInt64 + 1
)

func errSequenceInvalidData(db, name string) error {
	return &Error{Code: ErrSequenceInvalidData, SQLState: sqlState(ErrSequenceInvalidData),
		Message: fmt.Sprintf("Sequence '%s.%s' values are conflicting", db, name)}
}

func errUnknownSequence(db, name string) error {
	return &Error{Code: ErrUnknownSequence, SQLState: sqlState(ErrUnknownSequence),
		Message: fmt.Sprintf("Unknown SEQUENCE: '%s.%s'", db, name)}
}

// nameTaken reports whether a table, view or sequence of db has the name. A
// local temporary table lives in its own namespace, TempTables, and does not
// collide with sequences (read on v8.5.5).
func (db *Database) nameTaken(key string) bool {
	return db.nameTakenFor(key, false)
}

// createSequence installs a sequence as TiDB's handleSequenceOptions does: the
// options apply in order, and MINVALUE, MAXVALUE and START that none of them
// set take the defaults of the increment's sign, unless all three are set.
func (c *Catalog) createSequence(stmt *nodes.CreateSequenceStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	if db.nameTaken(key) {
		if stmt.IfNotExists {
			return nil
		}
		return errDupTableQualified(db.Name, stmt.Name.Name)
	}
	seq := &Sequence{Name: stmt.Name.Name, Database: db, Increment: 1, Cache: sequenceDefaultCache}
	var minSet, maxSet, startSet bool
	for _, opt := range stmt.Options {
		if refusedSequenceCache(opt) {
			return errSequenceInvalidData(db.Name, stmt.Name.Name)
		}
		switch opt.Type {
		case nodes.SequenceOptionStartWith:
			startSet = true
		case nodes.SequenceOptionMinValue:
			minSet = true
		case nodes.SequenceOptionMaxValue:
			maxSet = true
		default:
		}
		applySequenceOption(seq, opt)
	}
	if !minSet || !maxSet || !startSet {
		if seq.Increment >= 0 {
			if !minSet {
				seq.MinValue = 1
			}
			if !startSet {
				seq.Start = max(seq.MinValue, 1)
			}
			if !maxSet {
				seq.MaxValue = sequencePositiveMaxDefault
			}
		} else {
			if !maxSet {
				seq.MaxValue = -1
			}
			if !startSet {
				seq.Start = min(seq.MaxValue, -1)
			}
			if !minSet {
				seq.MinValue = sequenceNegativeMinDefault
			}
		}
	}
	for _, opt := range stmt.TableOptions {
		if strings.EqualFold(opt.Name, "COMMENT") {
			seq.Comment = opt.Value
		}
	}
	if !validSequence(seq) {
		return errSequenceInvalidData(db.Name, stmt.Name.Name)
	}
	db.Sequences[key] = seq
	return nil
}

// alterSequence applies the options to the sequence and keeps the result only
// if it is valid. TiDB fills in no defaults here, answers 1347 when a table or
// view holds the name, and IF EXISTS does not spare a missing sequence its
// 1146 (read on v8.5.5).
func (c *Catalog) alterSequence(stmt *nodes.AlterSequenceStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	seq := db.Sequences[key]
	if seq == nil {
		if db.Tables[key] != nil || db.Views[key] != nil {
			return errWrongObjectKind(db.Name, stmt.Name.Name, "SEQUENCE")
		}
		return errNoSuchTable(db.Name, stmt.Name.Name)
	}
	next := *seq
	for _, opt := range stmt.Options {
		if refusedSequenceCache(opt) {
			return errSequenceInvalidData(db.Name, stmt.Name.Name)
		}
		applySequenceOption(&next, opt)
	}
	if !validSequence(&next) {
		return errSequenceInvalidData(db.Name, stmt.Name.Name)
	}
	*seq = next
	return nil
}

// dropSequence drops each named sequence it finds and answers 1347 when a
// non-temporary table or view holds the name, or 4139 when nothing does, for
// the first name it fails on, unless IF EXISTS (which suppresses both on
// v8.5.5).
func (c *Catalog) dropSequence(stmt *nodes.DropSequenceStmt) error {
	var first error
	for _, ref := range stmt.Sequences {
		db, err := c.resolveDatabase(ref.Schema)
		if err != nil {
			if first == nil {
				first = err
			}
			continue
		}
		key := toLower(ref.Name)
		if db.Sequences[key] == nil {
			if !stmt.IfExists && first == nil {
				if db.Tables[key] != nil || db.Views[key] != nil {
					first = errWrongObjectKind(db.Name, ref.Name, "SEQUENCE")
				} else {
					first = errUnknownSequence(db.Name, ref.Name)
				}
			}
			continue
		}
		delete(db.Sequences, key)
	}
	return first
}

func applySequenceOption(seq *Sequence, opt *nodes.SequenceOption) {
	switch opt.Type {
	case nodes.SequenceOptionIncrementBy:
		seq.Increment = opt.Value
	case nodes.SequenceOptionStartWith:
		seq.Start = opt.Value
	case nodes.SequenceOptionMinValue:
		seq.MinValue = opt.Value
	case nodes.SequenceOptionMaxValue:
		seq.MaxValue = opt.Value
	case nodes.SequenceOptionCache:
		seq.Cache = opt.Value
	case nodes.SequenceOptionNoCache:
		seq.Cache = 0
	case nodes.SequenceOptionCycle:
		seq.Cycle = true
	case nodes.SequenceOptionNoCycle:
		seq.Cycle = false
	default:
		// NO MINVALUE, NO MAXVALUE and RESTART change no option TiDB keeps.
	}
}

// refusedSequenceCache reports a CACHE option below 1, which TiDB refuses
// (4136) while NOCACHE sets the same cache size of 0.
func refusedSequenceCache(opt *nodes.SequenceOption) bool {
	return opt.Type == nodes.SequenceOptionCache && opt.Value <= 0
}

// validSequence is TiDB's validateSequenceOptions.
func validSequence(seq *Sequence) bool {
	if seq.Increment == 0 {
		return false
	}
	step := seq.Increment
	if step < 0 {
		step = -step
	}
	return seq.MaxValue >= seq.Start &&
		seq.MaxValue > seq.MinValue &&
		seq.Start >= seq.MinValue &&
		seq.MaxValue != math.MaxInt64 &&
		seq.MinValue != math.MinInt64 &&
		seq.Cache < (math.MaxInt64-step)/step
}

// ShowCreateSequence produces TiDB's SHOW CREATE SEQUENCE output, or "" when
// the database or sequence does not exist.
func (c *Catalog) ShowCreateSequence(database, name string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	seq := db.Sequences[toLower(name)]
	if seq == nil {
		return ""
	}
	return showCreateSequence(seq)
}

// showCreateSequence writes a sequence as v8.5.5's SHOW CREATE SEQUENCE does:
// every option, in lower case, then ENGINE=InnoDB and the comment.
func showCreateSequence(seq *Sequence) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE SEQUENCE %s start with %d minvalue %d maxvalue %d increment by %d",
		quoteIdent(seq.Name), seq.Start, seq.MinValue, seq.MaxValue, seq.Increment)
	if seq.Cache == 0 {
		b.WriteString(" nocache")
	} else {
		b.WriteString(" cache " + strconv.FormatInt(seq.Cache, 10))
	}
	if seq.Cycle {
		b.WriteString(" cycle")
	} else {
		b.WriteString(" nocycle")
	}
	b.WriteString(" ENGINE=InnoDB")
	if seq.Comment != "" {
		fmt.Fprintf(&b, " COMMENT='%s'", outputFormat(seq.Comment))
	}
	return b.String()
}
