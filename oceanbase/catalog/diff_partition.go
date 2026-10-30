package catalog

import "strings"

// MySQL SDL diff — table partitioning (the PARTITION BY clause, compared as one spec).
//
// MySQL partitioning is a single per-table clause, not a set of independent sub-objects, so
// the diff signal is a coarse bool (TableDiffEntry.PartitionChanged): "does the partition
// spec differ?". diffPartitions is wired into compareTable (diff_table.go) and folded into
// the "is this table modified?" decision via tableSubdiffsChanged.
//
// IDEMPOTENCE is the whole point of this node. A partitioned table's SHOW CREATE emits a
// canonical partition clause whose surface form differs from the user's DDL — most visibly,
// every explicitly-defined partition is echoed with a trailing `ENGINE = <engine>` even when
// the user wrote none. The omni loader (buildPartitionInfo in tablecmds.go) already absorbs
// the OTHER surface differences into a single canonical PartitionInfo: it normalizes the
// partition EXPRESSION identically for the user form and the engine readback (so 8.0's
// backtick-quoted, lower-cased `RANGE (year(`dt`))`, 5.7's bare `RANGE (YEAR(dt))`, and the
// user's `RANGE (YEAR(dt))` all load to the same Expr), folds RANGE/LIST COLUMNS and KEY
// columns into Columns, and expands `PARTITIONS N` / `SUBPARTITIONS N` defaulting into the
// UseDefault* flags. (Verified on live 5.7.25 + 8.0.32: for every partition form the ONLY
// field that differs between the user DDL and its SHOW CREATE readback is the per-partition /
// per-subpartition Engine — empty in the user form, the table's storage engine in the
// readback.)
//
// So canonicalization here is exactly: resolve each partition's (and subpartition's) effective
// storage engine — empty means "the table's engine" — and otherwise compare the partition spec
// rendered through the SAME deparser show.go uses for the stored form (showPartitioning). Using
// the deparser as the comparison key keeps a single source of truth: the form the diff treats
// as canonical is byte-for-byte the form generatePartitionDDL re-emits, which is the form the
// engine stores — so apply-correctness and idempotence hold by construction.
//
// implemented by omni:partitions breadth node
func diffPartitions(from, to *Table, n *Normalizer) bool {
	return canonicalPartitionSpec(from, n) != canonicalPartitionSpec(to, n)
}

// canonicalPartitionSpec returns the canonical signature of a table's partitioning, or "" when
// the table is not partitioned. Two tables have the same partitioning iff their signatures are
// equal. The signature is showPartitioning over a copy whose per-partition/subpartition engine
// is resolved to the effective engine (empty → the table's storage engine), folding the only
// user-vs-stored surface difference; every other surface difference is already canonicalized by
// the loader (see diffPartitions doc).
func canonicalPartitionSpec(t *Table, n *Normalizer) string {
	if t == nil || t.Partitioning == nil {
		return ""
	}
	return showPartitioning(partitionSpecWithResolvedEngine(t, n))
}

// keyAlgorithmDefault is the KEY-partitioning ALGORITHM that MySQL omits from SHOW CREATE: the
// modern algorithm (2) is the default and is NOT echoed, while the legacy algorithm (1) IS echoed
// (via a /*!50611 ALGORITHM = 1 */ split comment). The catalog uses 0 to mean "no ALGORITHM
// written"; the loader stores 2 only when the user wrote it explicitly. Folding 2 → 0 makes an
// explicit `KEY ALGORITHM=2` compare equal to the engine's stripped `KEY` form (verified on live
// 5.7.25 + 8.0.32). Algorithm 1 is left intact so its echoed form round-trips. (Concern raised by
// cross-family review and confirmed on both engines.)
const keyAlgorithmDefault = 2

// partitionSpecWithResolvedEngine returns a copy of the table's PartitionInfo canonicalized so that
// the user form and the engine's SHOW CREATE readback collapse onto the same showPartitioning
// output. Two surface differences are folded:
//
//   - per-partition / per-subpartition Engine: resolved to the table's effective storage engine
//     when empty. The table default (InnoDB unless ENGINE=... says otherwise) is what SHOW CREATE
//     echoes per partition, so folding empty → that value makes the user form (no per-partition
//     ENGINE) and the stored form (every partition carries ENGINE = <engine>) compare equal —
//     including a MyISAM table, whose readback echoes ENGINE = MyISAM, not InnoDB.
//   - KEY partitioning ALGORITHM = 2: the default algorithm, stripped from SHOW CREATE, is folded
//     to 0 (unspecified) so an explicit `ALGORITHM=2` matches the readback (keyAlgorithmDefault).
//
// HASH/KEY partitions defined only by PARTITIONS N (no explicit definitions) read back as
// `PARTITIONS N` with NO per-partition engine on either side; showPartitioning renders that
// PARTITIONS-N form without touching per-partition engines, so the engine resolution is a no-op
// for them and they round-trip regardless.
//
// Partition BOUND values (`VALUES LESS THAN (...)` / `VALUES IN (...)`) and the partition
// EXPRESSION are routed through n.CanonicalPartitionValue / n.CanonicalPartitionExpr. MySQL
// constant-folds a non-literal bound at storage time (`5+5`→`10`, `TO_DAYS('2020-01-01')`→`737790`)
// and echoes only the folded literal; the loader (partitionValueItemSQL) folds the user form the
// same way at load, and these canonicalizers collapse incidental whitespace so the two compare
// equal. A bound the folder cannot evaluate (an unsupported/non-deterministic function) is kept
// verbatim — a flagged round-trip limitation (see normalize.go entry partition-constant-folding /
// the PR flag); literal bounds (the common case) are unaffected.
func partitionSpecWithResolvedEngine(t *Table, n *Normalizer) *PartitionInfo {
	pi := t.Partitioning
	engine := n.CanonicalEngine(t) // lower-cased effective engine; never returns ""

	out := *pi
	out.Algorithm = foldKeyAlgorithm(pi.Algorithm)
	out.SubAlgo = foldKeyAlgorithm(pi.SubAlgo)
	out.Expr = n.CanonicalPartitionExpr(pi.Expr)
	out.SubExpr = n.CanonicalPartitionExpr(pi.SubExpr)
	out.Partitions = make([]*PartitionDefInfo, len(pi.Partitions))
	for i, pd := range pi.Partitions {
		cp := *pd
		cp.Engine = resolvePartitionEngine(pd.Engine, engine)
		cp.ValueExpr = canonicalPartitionValueExpr(pd.ValueExpr, n)
		if len(pd.SubPartitions) > 0 {
			cp.SubPartitions = make([]*SubPartitionDefInfo, len(pd.SubPartitions))
			for j, sp := range pd.SubPartitions {
				csp := *sp
				csp.Engine = resolvePartitionEngine(sp.Engine, engine)
				cp.SubPartitions[j] = &csp
			}
		}
		out.Partitions[i] = &cp
	}
	return &out
}

// canonicalPartitionValueExpr canonicalizes a partition definition's bound value, preserving the
// MAXVALUE sentinel verbatim (CanonicalPartitionValue's whitespace collapse must never rewrite it,
// and showPartitioning special-cases it). Empty (HASH/KEY partitions carry no bound) stays empty.
func canonicalPartitionValueExpr(value string, n *Normalizer) string {
	if value == "" || value == "MAXVALUE" {
		return value
	}
	return n.CanonicalPartitionValue(value)
}

// resolvePartitionEngine resolves a partition's declared engine against the table's effective
// engine, case-insensitively. An empty declared engine inherits the table engine; a declared
// engine equal to the table engine is normalized to the table-engine spelling (so a user's
// lower-case `engine=innodb` and a readback's `ENGINE = InnoDB` agree). The comparison is
// case-insensitive because MySQL engine names are.
func resolvePartitionEngine(declared, tableEngine string) string {
	d := strings.TrimSpace(declared)
	if d == "" || strings.EqualFold(d, tableEngine) {
		return tableEngine
	}
	return d
}

// foldKeyAlgorithm maps the default KEY-partitioning ALGORITHM (2) to 0 (unspecified), so an
// explicit `ALGORITHM=2` and the engine's stripped readback compare equal. Algorithm 1 (legacy,
// echoed by SHOW CREATE) and 0 (unwritten) pass through unchanged.
func foldKeyAlgorithm(algo int) int {
	if algo == keyAlgorithmDefault {
		return 0
	}
	return algo
}
