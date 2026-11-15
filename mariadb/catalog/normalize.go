package catalog

import (
	"fmt"
	nodes "github.com/bytebase/omni/mariadb/ast"
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

// Version is the MariaDB release whose stored form the catalog reproduces. The
// zero value is the pinned release, MariaDB118.
type Version int

const (
	// MariaDB118 is the 11.8 LTS release, mariadb:11.8.8.
	MariaDB118 Version = iota
	// MariaDB107 is mariadb:10.7.8, the oldest release Bytebase supports.
	MariaDB107
	// MariaDB1011 is the 10.11 LTS release, mariadb:10.11.19.
	MariaDB1011
	// MariaDB114 is the 11.4 LTS release, mariadb:11.4.13.
	MariaDB114
)

// release orders the versions by the release each stands for.
func (v Version) release() int {
	switch v {
	case MariaDB107:
		return 1007
	case MariaDB1011:
		return 1011
	case MariaDB114:
		return 1104
	default:
		return 1108
	}
}

// ParseVersion maps a server version, as SELECT VERSION() reports it
// ("10.11.6-MariaDB-1:10.11.6+maria~ubu2204"), to the Version whose stored
// form it has. Each Version is measured on its own release; a release between
// two of them takes the older one's form, and one past 11.8 takes 11.8's. ok
// is false for a version that does not parse or predates 10.7.
func ParseVersion(s string) (v Version, ok bool) {
	var major, minor int
	if _, err := fmt.Sscanf(strings.TrimSpace(s), "%d.%d", &major, &minor); err != nil {
		return MariaDB118, false
	}
	switch release := major*100 + minor; {
	case release < 1007:
		return MariaDB107, false
	case release < 1011:
		return MariaDB107, true
	case release < 1104:
		return MariaDB1011, true
	case release < 1108:
		return MariaDB114, true
	default:
		return MariaDB118, true
	}
}

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

// NormalizerFor returns the Normalizer of a release at its server defaults:
// explicit_defaults_for_timestamp is OFF on 10.7 and ON from 10.10.
func NormalizerFor(v Version) *Normalizer {
	return &Normalizer{Version: v, ExplicitDefaultsForTimestamp: v.release() >= 1010}
}

// Version returns the release whose stored form this catalog reproduces.
func (c *Catalog) Version() Version { return c.session.Version }

// SetVersion makes the catalog reproduce release v from here on: what v
// stores for the objects it creates, and v's server defaults, the collation a
// character set takes and explicit_defaults_for_timestamp among them. It does
// not rewrite what the catalog already holds.
func (c *Catalog) SetVersion(v Version) {
	c.session.Version = v
	c.session.ExplicitDefaultsForTimestamp = NormalizerFor(v).ExplicitDefaultsForTimestamp
	c.defaultCollation = c.charsetCollation(c.defaultCharset)
	c.collationConn = c.charsetCollation(c.charsetClient)
	c.session.CollationConnection = c.collationConn
	for _, db := range c.databases {
		db.version = v
	}
}

// ignoreDescendingKeys makes every key part of the table ref names ascending on
// a release before 10.8, which takes DESC and stores the part ascending.
// Measured on mariadb:10.7.8: KEY ix (b DESC) is stored KEY `ix` (`b`). A
// statement calls it, deferred, for the one table it changes; the table's
// pointer, resolved first, follows an ALTER that renames it.
func (c *Catalog) ignoreDescendingKeys(ref *nodes.TableRef) func() {
	if c.session.Version.release() >= 1008 || ref == nil {
		return func() {}
	}
	tbl := c.tableOf(ref)
	return func() {
		if tbl == nil {
			tbl = c.tableOf(ref)
		}
		if tbl == nil {
			return
		}
		for _, idx := range tbl.Indexes {
			for _, ic := range idx.Columns {
				ic.Descending = false
			}
		}
	}
}

func (c *Catalog) tableOf(ref *nodes.TableRef) *Table {
	db, err := c.resolveDatabase(ref.Schema)
	if err != nil {
		return nil
	}
	return db.Tables[toLower(ref.Name)]
}

// charsetCollation is the collation a character set takes where none is named,
// on the catalog's release. MariaDB 11.4 makes utf8mb4 default to
// utf8mb4_uca1400_ai_ci, and 11.8 the other Unicode character sets as well;
// before 11.4 each takes its _general_ci. Measured on mariadb:10.7.8,
// 10.11.19, 11.4.13 and 11.8.8 (information_schema.COLLATIONS and
// character_set_collations).
func (c *Catalog) charsetCollation(charset string) string {
	return releaseCollation(c.session.Version, charset)
}

// versionCollation is charsetCollation for the release of db's catalog.
func versionCollation(db *Database, charset string) string {
	v := MariaDB118
	if db != nil {
		v = db.version
	}
	return releaseCollation(v, charset)
}

func releaseCollation(v Version, charset string) string {
	cs := toLower(charset)
	coll := defaultCollationForCharset[cs]
	if !strings.Contains(coll, "_uca1400_") {
		return coll
	}
	switch release := v.release(); {
	case release >= 1108:
		return coll
	case release >= 1104 && cs == "utf8mb4":
		return coll
	default:
		return cs + "_general_ci"
	}
}

// CanonicalColumn is a column as SHOW CREATE TABLE writes it, with the check
// it carries, since MariaDB writes a column's check on the column and changes
// it only with the column's definition, and with its character set and
// collation written out: a table's default changes the default of columns
// added later and leaves the existing ones as they were.
func (n *Normalizer) CanonicalColumn(table *Table, c *Column) string {
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
