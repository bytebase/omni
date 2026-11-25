package catalog

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"unicode/utf8"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// Sequence is a sequence of a database. OceanBase keeps a MySQL-mode
// sequence's numbers in 64 bits, but stores a bound that ALTER SEQUENCE
// clears (NOMAXVALUE, NOMINVALUE) as one past them, 2^63 or -2^63-1, so the
// numbers are kept whole here.
//
// Ref: ObSequenceOption and ObSequenceMaxMinInitializer in
// src/share/sequence/ob_sequence_option.h (oceanbase v4.3.5_CE_BP6).
type Sequence struct {
	Name      string
	Start     *big.Int
	Min       *big.Int
	Max       *big.Int
	Increment *big.Int
	Cache     *big.Int
	Cycle     bool
	Order     bool
}

// The bounds of a MySQL-mode sequence, and the values that stand for none.
var (
	seqMaxValue   = big.NewRat(math.MaxInt64, 1)
	seqMinValue   = big.NewRat(math.MinInt64, 1)
	seqNoMaxValue = new(big.Rat).Add(seqMaxValue, big.NewRat(1, 1))
	seqNoMinValue = new(big.Rat).Sub(seqMinValue, big.NewRat(1, 1))
)

// sequenceOption is ObSequenceOption: the numbers as ObNumber holds them,
// which may be fractions until the option builder checks them.
type sequenceOption struct {
	increment, start, max, min, cache *big.Rat
	cycle, order                      bool
}

func newSequenceOption() *sequenceOption {
	return &sequenceOption{
		increment: big.NewRat(1, 1),
		start:     new(big.Rat),
		max:       new(big.Rat).Set(seqNoMaxValue),
		min:       new(big.Rat).Set(seqNoMinValue),
		cache:     big.NewRat(20, 1),
	}
}

func (o *sequenceOption) hasMax() bool { return o.max.Cmp(seqNoMaxValue) != 0 }
func (o *sequenceOption) hasMin() bool { return o.min.Cmp(seqNoMinValue) != 0 }

// The options a statement set, ObSequenceArg's bit set.
const (
	seqIncrementBy = 1 << iota
	seqStartWith
	seqMaxValueBit
	seqMinValueBit
	seqCacheBit
	seqNoMaxValueBit
	seqNoMinValueBit
	seqNoCacheBit
	seqCycleBit
	seqNoCycleBit
	seqOrderBit
	seqNoOrderBit
	seqRestartBit
)

// resolveSequenceOptions reads a statement's options as
// ObSequenceResolver::resolve_sequence_options does: each in written order,
// refusing one given twice, then the pairs that contradict each other, then
// RESTART outside ALTER SEQUENCE.
//
// Ref: src/sql/resolver/ddl/ob_sequence_resolver.h (oceanbase v4.3.5_CE_BP6).
func resolveSequenceOptions(opts []*nodes.SequenceOption, alter bool) (*sequenceOption, int, error) {
	o, set := newSequenceOption(), 0
	for _, opt := range opts {
		var bit, dup int
		var value **big.Rat
		switch opt.Name {
		case "INCREMENT BY":
			bit, dup, value = seqIncrementBy, seqIncrementBy, &o.increment
		case "START WITH":
			bit, dup, value = seqStartWith, seqStartWith, &o.start
		case "MAXVALUE":
			bit, dup, value = seqMaxValueBit, seqMaxValueBit, &o.max
		case "MINVALUE":
			bit, dup, value = seqMinValueBit, seqMinValueBit, &o.min
		case "CACHE":
			bit, dup, value = seqCacheBit, seqCacheBit, &o.cache
		case "NOMAXVALUE":
			bit, dup = seqNoMaxValueBit, seqNoMaxValueBit
			o.max = new(big.Rat).Set(seqNoMaxValue)
		case "NOMINVALUE":
			bit, dup = seqNoMinValueBit, seqNoMinValueBit
			o.min = new(big.Rat).Set(seqNoMinValue)
		case "NOCACHE":
			bit, dup = seqNoCacheBit, seqNoCacheBit
			o.cache = big.NewRat(1, 1)
		case "CYCLE":
			bit, dup, o.cycle = seqCycleBit, seqCycleBit, true
		case "NOCYCLE":
			bit, dup, o.cycle = seqNoCycleBit, seqNoCycleBit, false
		case "ORDER":
			bit, dup, o.order = seqOrderBit, seqOrderBit, true
		case "NOORDER":
			bit, dup, o.order = seqNoOrderBit, seqNoOrderBit, false
		case "RESTART":
			bit, dup = seqRestartBit, seqRestartBit
		default:
			return nil, 0, fmt.Errorf("sequence option %q", opt.Name)
		}
		if set&dup != 0 {
			return nil, 0, errSequenceDuplicate(opt.Name)
		}
		if value != nil {
			// A number with an exponent is a DOUBLE to OceanBase's lexer, and a
			// sequence takes none.
			if strings.ContainsAny(opt.Value, "eE") {
				return nil, 0, seqError(5114, "42000", "Invalid numeric")
			}
			n, ok := new(big.Rat).SetString(opt.Value)
			if !ok {
				return nil, 0, fmt.Errorf("sequence option %s %q is no number", opt.Name, opt.Value)
			}
			// get_normalized_number: a number past the bounds is the bound.
			switch {
			case n.Cmp(seqMaxValue) > 0:
				n.Set(seqMaxValue)
			case n.Cmp(seqMinValue) < 0:
				n.Set(seqMinValue)
			}
			*value = n
		}
		set |= bit
	}
	for _, c := range []struct {
		a, b int
		err  error
	}{
		{seqMaxValueBit, seqNoMaxValueBit, seqError(4306, "HY000", "conflicting MAXVALUE/NOMAXVALUE specifications")},
		{seqMinValueBit, seqNoMinValueBit, seqError(4307, "HY000", "conflicting MINVALUE/NOMINVALUE specifications")},
		{seqCacheBit, seqNoCacheBit, seqError(4309, "HY000", "conflicting CACHE/NOCACHE specifications")},
		{seqOrderBit, seqNoOrderBit, seqError(4310, "HY000", "conflicting ORDER/NOORDER specifications")},
		{seqCycleBit, seqNoCycleBit, seqError(4308, "HY000", "conflicting CYCLE/NOCYCLE specifications")},
	} {
		if set&c.a != 0 && set&c.b != 0 {
			return nil, 0, c.err
		}
	}
	if !alter && set&seqRestartBit != 0 {
		return nil, 0, seqError(1149, "42000", "You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use")
	}
	return o, set, nil
}

// checkSequenceOption is ObSequenceOptionBuilder::check_sequence_option.
//
// Ref: src/share/sequence/ob_sequence_option_builder.cpp (oceanbase
// v4.3.5_CE_BP6).
func checkSequenceOption(set int, o *sequenceOption) error {
	if err := checkSequenceNumbers(set, o); err != nil {
		return err
	}
	step := new(big.Rat).Abs(o.increment)
	valueRange := new(big.Rat).Sub(o.max, o.min)
	cacheRange := new(big.Rat).Mul(new(big.Rat).Sub(o.cache, big.NewRat(1, 1)), step)
	switch {
	case o.hasMin() && o.hasMax() && o.min.Cmp(o.max) >= 0:
		return seqError(4320, "HY000", "MINVALUE must be less than MAXVALUE")
	case o.hasMax() && o.start.Cmp(o.max) > 0:
		return seqError(4324, "HY000", "START WITH cannot be more than MAXVALUE")
	case o.hasMin() && o.start.Cmp(o.min) < 0:
		return seqError(4322, "HY000", "START WITH cannot be less than MINVALUE")
	case valueRange.Cmp(step) < 0:
		return seqError(4321, "HY000", "INCREMENT must be less than MAXVALUE minus MINVALUE")
	}
	if err := checkSequenceCycle(o); err != nil {
		return err
	}
	if o.cycle && valueRange.Cmp(cacheRange) <= 0 {
		return seqError(4328, "HY000", "number to CACHE must be less than one cycle")
	}
	return nil
}

// checkSequenceNumbers is the start of check_sequence_option, which reads
// only the numbers the statement sets: each a whole number, a nonzero
// increment, and a cache of more than one.
func checkSequenceNumbers(set int, o *sequenceOption) error {
	for _, c := range []struct {
		bit  int
		v    *big.Rat
		name string
	}{
		{seqMinValueBit, o.min, "MINVALUE"},
		{seqMaxValueBit, o.max, "MAXVALUE"},
		{seqIncrementBy, o.increment, "INCREMENT"},
		{seqStartWith, o.start, "START WITH"},
		{seqCacheBit, o.cache, "CACHE"},
	} {
		if set&c.bit != 0 && !c.v.IsInt() {
			return seqError(4317, "HY000", fmt.Sprintf("sequence parameter %s must be an integer", c.name))
		}
	}
	switch {
	case o.increment.Sign() == 0:
		return seqError(4318, "HY000", "INCREMENT must be a nonzero integer")
	case o.cache.Cmp(big.NewRat(1, 1)) <= 0 && set&seqCacheBit != 0:
		return seqError(4326, "HY000", "the number of values to CACHE must be greater than 1")
	}
	return nil
}

// checkSequenceCycle refuses a cycle with no bound to turn at.
func checkSequenceCycle(o *sequenceOption) error {
	switch {
	case o.cycle && o.increment.Sign() > 0 && !o.hasMax():
		return seqError(4330, "HY000", "ascending sequences that CYCLE must specify MAXVALUE")
	case o.cycle && o.increment.Sign() < 0 && !o.hasMin():
		return seqError(4329, "HY000", "descending sequences that CYCLE must specify MINVALUE")
	}
	return nil
}

// checkSequenceName refuses a name longer than OceanBase takes, which the
// resolver checks before the options.
func checkSequenceName(name string) error {
	if utf8.RuneCountInString(name) > 64 {
		return seqError(1059, "42000", fmt.Sprintf("Identifier name '%s' is too long", name))
	}
	return nil
}

// resolveSequenceStatement reads a sequence statement as its resolver does:
// the name, then for ALTER SEQUENCE that it sets an option, then the options.
//
// Ref: ob_create_sequence_resolver.cpp, ob_alter_sequence_resolver.cpp and
// ob_drop_sequence_resolver.cpp (oceanbase v4.3.5_CE_BP6).
func resolveSequenceStatement(stmt *nodes.ObjectDDLStmt) (*sequenceOption, int, error) {
	if err := checkSequenceName(stmt.Name); err != nil {
		return nil, 0, err
	}
	alter := stmt.Action == "ALTER SEQUENCE"
	if alter && len(stmt.Options) == 0 {
		return nil, 0, seqError(4314, "HY000", "no options specified for ALTER SEQUENCE")
	}
	return resolveSequenceOptions(stmt.Options, alter)
}

// buildCreateSequenceOption is ObSequenceOptionBuilder::
// build_create_sequence_option: a cycle with no bound is refused before the
// bounds are filled in, then the start defaults to the bound the sequence
// counts from, then the options are checked.
func buildCreateSequenceOption(o *sequenceOption, set int) error {
	if err := checkSequenceCycle(o); err != nil {
		return err
	}
	switch o.increment.Sign() {
	case 1:
		if !o.hasMin() {
			o.min = big.NewRat(1, 1)
		}
		if !o.hasMax() {
			o.max = new(big.Rat).Set(seqMaxValue)
		}
	case -1:
		if !o.hasMax() {
			o.max = big.NewRat(-1, 1)
		}
		if !o.hasMin() {
			o.min = new(big.Rat).Set(seqMinValue)
		}
	}
	if set&seqStartWith == 0 {
		switch {
		case o.hasMin() && o.increment.Sign() > 0:
			o.start = new(big.Rat).Set(o.min)
		case o.hasMax() && o.increment.Sign() < 0:
			o.start = new(big.Rat).Set(o.max)
		}
	}
	return checkSequenceOption(set, o)
}

// CheckSequenceStatement reports the refusal OceanBase gives a CREATE, ALTER
// or DROP SEQUENCE whatever its database holds, nil where that depends on the
// database. Its resolver reads the name and the options before anything else,
// and a statement whose options the builder refuses is refused either way: a
// CREATE SEQUENCE without IF NOT EXISTS with 4297 when the name is taken, an
// ALTER SEQUENCE with 4298 when the sequence is missing, and with the
// builder's error otherwise. Of an ALTER SEQUENCE the builder's checks that
// read only the statement's own options are known: START WITH without
// RESTART, and checkSequenceNumbers, which an unset increment (the
// sequence's, never zero) and an unset cache pass.
func CheckSequenceStatement(stmt *nodes.ObjectDDLStmt) error {
	o, set, err := resolveSequenceStatement(stmt)
	switch {
	case err != nil:
		return err
	case stmt.Action == "CREATE SEQUENCE" && !stmt.IfExists:
		return buildCreateSequenceOption(o, set)
	case stmt.Action == "ALTER SEQUENCE":
		if err := checkAlterSequenceStart(set); err != nil {
			return err
		}
		return checkSequenceNumbers(set, o)
	}
	return nil
}

// checkAlterSequenceStart refuses START WITH without RESTART, which
// build_alter_sequence_option checks before the options.
func checkAlterSequenceStart(set int) error {
	if set&seqRestartBit == 0 && set&seqStartWith != 0 {
		return seqError(4311, "HY000", "cannot alter starting sequence number")
	}
	return nil
}

// execSequence runs CREATE, ALTER or DROP SEQUENCE: the statement's name and
// options first, as the resolver reads them, then its database, then whether
// the sequence exists, then what the options make of it.
//
// Ref: src/share/sequence/ob_sequence_ddl_proxy.cpp (oceanbase v4.3.5_CE_BP6).
func (c *Catalog) execSequence(stmt *nodes.ObjectDDLStmt) error {
	o, set, err := resolveSequenceStatement(stmt)
	if err != nil {
		return err
	}
	db, err := c.resolveDatabase(stmt.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name)
	seq, exists := db.Sequences[key]
	switch stmt.Action {
	case "CREATE SEQUENCE":
		if exists {
			if stmt.IfExists {
				return nil
			}
			return seqError(4297, "HY000", "name is already used by an existing object")
		}
		if err := buildCreateSequenceOption(o, set); err != nil {
			return err
		}
		db.Sequences[key] = o.sequence(stmt.Name)
	case "ALTER SEQUENCE":
		if !exists {
			return seqError(4298, "HY000", "sequence does not exist")
		}
		if err := buildAlterSequenceOption(o, set, seq.option()); err != nil {
			return err
		}
		db.Sequences[key] = o.sequence(seq.Name)
	case "DROP SEQUENCE":
		if !exists {
			if stmt.IfExists {
				return nil
			}
			return seqError(4298, "HY000", "sequence does not exist")
		}
		delete(db.Sequences, key)
	}
	return nil
}

// buildAlterSequenceOption is ObSequenceOptionBuilder::
// build_alter_sequence_option: what the statement leaves out is the
// sequence's own, and START WITH moves only with RESTART, which without it
// restarts at the bound the sequence counts from.
func buildAlterSequenceOption(o *sequenceOption, set int, old *sequenceOption) error {
	if set&seqIncrementBy == 0 {
		o.increment = old.increment
	}
	if set&seqStartWith == 0 {
		o.start = old.start
	}
	if set&(seqMaxValueBit|seqNoMaxValueBit) == 0 {
		o.max = old.max
	}
	if set&(seqMinValueBit|seqNoMinValueBit) == 0 {
		o.min = old.min
	}
	if set&(seqCacheBit|seqNoCacheBit) == 0 {
		o.cache = old.cache
	}
	if set&(seqCycleBit|seqNoCycleBit) == 0 {
		o.cycle = old.cycle
	}
	if set&(seqOrderBit|seqNoOrderBit) == 0 {
		o.order = old.order
	}
	if err := checkAlterSequenceStart(set); err != nil {
		return err
	}
	if set&seqRestartBit != 0 && set&seqStartWith == 0 {
		if o.increment.Sign() > 0 {
			o.start = new(big.Rat).Set(o.min)
		} else {
			o.start = new(big.Rat).Set(o.max)
		}
	}
	return checkSequenceOption(set, o)
}

// sequence stores the option, whose numbers the builder has checked are whole.
func (o *sequenceOption) sequence(name string) *Sequence {
	whole := func(r *big.Rat) *big.Int { return new(big.Int).Set(r.Num()) }
	return &Sequence{
		Name:      name,
		Start:     whole(o.start),
		Min:       whole(o.min),
		Max:       whole(o.max),
		Increment: whole(o.increment),
		Cache:     whole(o.cache),
		Cycle:     o.cycle,
		Order:     o.order,
	}
}

func (s *Sequence) option() *sequenceOption {
	rat := func(i *big.Int) *big.Rat { return new(big.Rat).SetInt(i) }
	return &sequenceOption{
		increment: rat(s.Increment),
		start:     rat(s.Start),
		max:       rat(s.Max),
		min:       rat(s.Min),
		cache:     rat(s.Cache),
		cycle:     s.Cycle,
		order:     s.Order,
	}
}

func errSequenceDuplicate(option string) error {
	switch option {
	case "INCREMENT BY":
		return seqError(4312, "HY000", "duplicate INCREMENT BY specifications")
	case "START WITH":
		return seqError(4313, "HY000", "duplicate START WITH specifications")
	case "MAXVALUE", "NOMAXVALUE":
		return seqError(4301, "HY000", "duplicate MAXVALUE/NOMAXVALUE specifications")
	case "MINVALUE", "NOMINVALUE":
		return seqError(4302, "HY000", "duplicate MINVALUE/NOMINVALUE specifications")
	case "CYCLE", "NOCYCLE":
		return seqError(4303, "HY000", "duplicate CYCLE/NOCYCLE specifications")
	case "CACHE", "NOCACHE":
		return seqError(4304, "HY000", "duplicate CACHE/NOCACHE specifications")
	case "ORDER", "NOORDER":
		return seqError(4305, "HY000", "duplicate ORDER/NOORDER specifications")
	}
	return seqError(4382, "HY000", "duplicate RESTART specified")
}

// seqError is one of OceanBase's own sequence errors, with the SQLSTATE and
// text OceanBase 4.3.5.6 answers.
func seqError(code int, state, message string) error {
	return &Error{Code: code, SQLState: state, Message: message}
}
