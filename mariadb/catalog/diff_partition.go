package catalog

import "strings"

// diffPartitions reports whether the partitioning of two tables differs.
// Partitioning is one clause of the table, so the result is a single bool
// (TableDiffEntry.PartitionChanged). The two sides compare as the text
// showPartitioning gives the spec once each partition's engine is resolved
// (canonicalPartitionSpec), so a partition that names no engine and one that names
// the table's engine are equal.
func diffPartitions(from, to *Table, n *Normalizer) bool {
	return canonicalPartitionSpec(from, n) != canonicalPartitionSpec(to, n)
}

// canonicalPartitionSpec returns showPartitioning of the table's partitioning
// with engines and the KEY algorithm folded (partitionSpecWithResolvedEngine), or
// "" when the table is not partitioned.
func canonicalPartitionSpec(t *Table, n *Normalizer) string {
	if t == nil || t.Partitioning == nil {
		return ""
	}
	return showPartitioning(partitionSpecWithResolvedEngine(t, n))
}

// keyAlgorithmDefault is the KEY ALGORITHM that compares equal to none written.
// The catalog uses 0 for no ALGORITHM written.
const keyAlgorithmDefault = 2

// partitionSpecWithResolvedEngine returns a copy of the table's PartitionInfo
// with each partition's and subpartition's engine resolved to the table's
// effective engine, the KEY algorithm folded, and the expression and bounds passed
// through the Normalizer.
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

// canonicalPartitionValueExpr passes a partition bound through the Normalizer,
// leaving MAXVALUE and an empty bound as they are.
func canonicalPartitionValueExpr(value string, n *Normalizer) string {
	if value == "" || value == "MAXVALUE" {
		return value
	}
	return n.CanonicalPartitionValue(value)
}

// resolvePartitionEngine returns the table's engine for a partition that declares
// none, or declares the same engine in another case; otherwise the declared
// engine.
func resolvePartitionEngine(declared, tableEngine string) string {
	d := strings.TrimSpace(declared)
	if d == "" || strings.EqualFold(d, tableEngine) {
		return tableEngine
	}
	return d
}

// foldKeyAlgorithm maps KEY ALGORITHM 2 to 0 (none written); other values pass
// through.
func foldKeyAlgorithm(algo int) int {
	if algo == keyAlgorithmDefault {
		return 0
	}
	return algo
}
