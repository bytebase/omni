package catalog

import (
	"strconv"
	"strings"
)

// The diff compares what MariaDB stores, and the catalog already holds each
// object in that form: Exec and LoadMetadata both resolve types, defaults,
// charsets and key types the way the server does, and the SHOW CREATE
// renderers write the result as the server writes it (held to the server by
// the container and simulate.tsv tests). So a column, a key, a partitioning
// and a table's options are compared by the text SHOW CREATE TABLE gives them,
// and a schema text and a synced snapshot of the same schema compare equal
// without a canonicalizer of their own.

// Version is the MariaDB release whose stored form the diff and the generated
// migration reproduce.
type Version int

const (
	// MariaDB118 is the 11.8 LTS release, mariadb:11.8.8.
	MariaDB118 Version = iota
)

// Normalizer carries what a comparison depends on beyond the catalogs
// themselves: the release, and explicit_defaults_for_timestamp, a session
// setting rather than a release trait.
type Normalizer struct {
	Version                      Version
	ExplicitDefaultsForTimestamp bool

	// checks caches columnChecks per table for one comparison, which renders
	// every column of a table: the map is the table's whole constraint list.
	checks map[*Table]map[string]*Constraint
}

// NormalizerFor returns the Normalizer of a release at its server defaults.
func NormalizerFor(v Version) *Normalizer {
	return &Normalizer{Version: v, ExplicitDefaultsForTimestamp: true}
}

// Version returns the release whose stored form this catalog's diff and
// migration reproduce.
func (c *Catalog) Version() Version { return c.session.Version }

// SetVersion fixes the release whose stored form this catalog's diff and
// migration reproduce.
func (c *Catalog) SetVersion(v Version) { c.session.Version = v }

// CanonicalColumn is a column as SHOW CREATE TABLE writes it, with the check
// it carries: MariaDB writes a column's check on the column, and it changes
// only with the column's definition.
func (n *Normalizer) CanonicalColumn(table *Table, c *Column) string {
	// The column's character set and collation are written out: a table's
	// default changes the default of columns added later and leaves the
	// existing ones as they were.
	context := Table{}
	if table != nil {
		context = *table
	}
	context.Charset, context.Collation = "", ""
	return showColumnWithTable(c, &context, n.columnChecks(table), false)
}

func (n *Normalizer) columnChecks(t *Table) map[string]*Constraint {
	if t == nil {
		return nil
	}
	if checks, ok := n.checks[t]; ok {
		return checks
	}
	if n.checks == nil {
		n.checks = make(map[*Table]map[string]*Constraint)
	}
	n.checks[t] = columnChecks(t)
	return n.checks[t]
}

// CanonicalIndex is a key as SHOW CREATE TABLE writes it.
func (n *Normalizer) CanonicalIndex(idx *Index) string {
	return showIndex(idx)
}

// CanonicalIndexColumn is one key part as SHOW CREATE TABLE writes it.
func (n *Normalizer) CanonicalIndexColumn(ic *IndexColumn) string {
	return showIndexColumn(ic)
}

// CanonicalTableOptions is the option list SHOW CREATE TABLE ends a table
// with, less AUTO_INCREMENT=N: the next value is data, not schema.
func (n *Normalizer) CanonicalTableOptions(t *Table) string {
	withoutCounter := *t
	withoutCounter.AutoIncrement = 0
	return showTableOptions(&withoutCounter)
}

// CanonicalPartitionExpr and CanonicalPartitionValue are a partitioning
// expression and a bound as the catalog stores them, which the loader has
// already folded the way the server does.
func (n *Normalizer) CanonicalPartitionExpr(expr string) string { return strings.TrimSpace(expr) }

// CanonicalPartitionValue: see CanonicalPartitionExpr.
func (n *Normalizer) CanonicalPartitionValue(value string) string { return strings.TrimSpace(value) }

// CanonicalCheckExpr is a check's expression. The catalog stores it as the
// server writes it.
func (n *Normalizer) CanonicalCheckExpr(expr string) string {
	return strings.TrimSpace(expr)
}

// CanonicalComment is a comment's text.
func (n *Normalizer) CanonicalComment(comment string) string {
	return comment
}

// CanonicalEngine is a table's engine, lower-cased, InnoDB where none is set.
func (n *Normalizer) CanonicalEngine(t *Table) string {
	if strings.TrimSpace(t.Engine) == "" {
		return "innodb"
	}
	return toLower(strings.TrimSpace(t.Engine))
}

// canonicalRoutineTypeText is a routine parameter's or RETURNS type with its
// keywords lower-cased and the display width an integer type takes by default
// left out, so INT and int(11) compare equal; ENUM and SET members keep their
// case.
func canonicalRoutineTypeText(typ string) string {
	dt, err := parseTypeName(typ)
	if err != nil {
		return strings.TrimSpace(typ)
	}
	return formatColumnType(dt)
}

// encodeKeyFields joins name/value pairs into one comparison key that no two
// different lists can share: each value is written with its length.
func encodeKeyFields(pairs ...string) string {
	var b strings.Builder
	for i := 0; i+1 < len(pairs); i += 2 {
		b.WriteString(pairs[i])
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(len(pairs[i+1])))
		b.WriteByte(':')
		b.WriteString(pairs[i+1])
		b.WriteByte('|')
	}
	return b.String()
}
