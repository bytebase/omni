// Package main implements the omni conformance runner: it diffs omni parser
// verdicts against engine ground truth (pre-labeled upstream test corpora,
// optionally adjudicated by a live container) and emits JSONL + a
// deterministic scoreboard. Design:
// CLAUDE_BB/plans/2026-07-07-omni-full-compatibility-conformance-design.md
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

// Verdict is a parse verdict from either side (omni or engine).
type Verdict string

const (
	VerdictAccept Verdict = "accept"
	VerdictReject Verdict = "reject"
	VerdictNone   Verdict = "" // not yet adjudicated / not applicable
)

// Class is the per-row classification produced by classify().
type Class string

const (
	ClassAgreeAccept   Class = "AGREE_ACCEPT"
	ClassAgreeReject   Class = "AGREE_REJECT"
	ClassGap           Class = "GAP"           // engine accepts, omni rejects — the hard-bar metric
	ClassOver          Class = "OVER"          // engine rejects, omni accepts
	ClassIndeterminate Class = "INDETERMINATE" // conflicting/unknown signals; manual queue
	ClassSkip          Class = "SKIP"
)

// Row is one corpus statement's full provenance — the JSONL schema from the
// design doc §3. Field names are the wire format; keep them stable.
type Row struct {
	Engine           string  `json:"engine"`
	Lane             string  `json:"lane"` // "upstream" | "generated"
	SourcePath       string  `json:"source_path"`
	Line             int     `json:"line"`
	TestName         string  `json:"test_name,omitempty"`
	SQL              string  `json:"sql"`
	StmtHash         string  `json:"stmt_hash"`
	Expected         Verdict `json:"expected,omitempty"`       // upstream label
	EngineVerdict    Verdict `json:"engine_verdict,omitempty"` // container adjudication
	OmniVerdict      Verdict `json:"omni_verdict"`
	OmniError        string  `json:"omni_error,omitempty"`
	RawErrorCode     int     `json:"raw_error_code,omitempty"` // engine error code (adjudicated rows)
	RawErrorMessage  string  `json:"raw_error_message,omitempty"`
	ClassifierReason string  `json:"classifier_reason,omitempty"`
	Family           string  `json:"family"`
	DivergenceKey    string  `json:"divergence_key,omitempty"`
	SkipReason       string  `json:"skip_reason,omitempty"`
	Class            Class   `json:"class"`
}

// classify computes the row class. Ground truth precedence: container
// adjudication beats the upstream label; label-vs-container disagreement is
// INDETERMINATE (extraction bug / stale label / context loss), per design §2.
func classify(r *Row) {
	truth := r.Expected
	if r.EngineVerdict != VerdictNone {
		if r.Expected != VerdictNone && r.Expected != r.EngineVerdict {
			r.Class = ClassIndeterminate
			r.ClassifierReason = "label_container_disagree"
			return
		}
		truth = r.EngineVerdict
	}
	switch {
	case truth == VerdictNone:
		r.Class = ClassIndeterminate
		r.ClassifierReason = "no_ground_truth"
	case truth == VerdictAccept && r.OmniVerdict == VerdictAccept:
		r.Class = ClassAgreeAccept
	case truth == VerdictReject && r.OmniVerdict == VerdictReject:
		r.Class = ClassAgreeReject
	case truth == VerdictAccept && r.OmniVerdict == VerdictReject:
		r.Class = ClassGap
		r.DivergenceKey = clusterKey(r.OmniError)
	default:
		r.Class = ClassOver
		// omni accepted: no omni error to key on. Pre-adjudication, key on the
		// leading tokens; adjudication upgrades this to the engine message.
		if r.RawErrorMessage != "" {
			r.DivergenceKey = clusterKey(r.RawErrorMessage)
		} else {
			r.DivergenceKey = leadingTokens(r.SQL, 4)
		}
	}
}

var familyPrefixes = []struct{ prefix, family string }{
	{"CREATE TABLE", "CREATE TABLE"}, {"CREATE INDEX", "CREATE INDEX"},
	{"CREATE UNIQUE INDEX", "CREATE INDEX"}, {"CREATE DATABASE", "CREATE DATABASE"},
	{"CREATE VIEW", "CREATE VIEW"}, {"CREATE USER", "DCL"},
	{"ALTER TABLE", "ALTER TABLE"}, {"ALTER DATABASE", "ALTER DATABASE"},
	{"DROP", "DROP"}, {"RENAME", "RENAME"}, {"TRUNCATE", "TRUNCATE"},
	{"SELECT", "SELECT"}, {"TABLE", "SELECT"}, {"VALUES", "SELECT"},
	{"WITH", "SELECT"}, {"(", "SELECT"},
	{"INSERT", "INSERT"}, {"REPLACE", "INSERT"}, {"UPDATE", "UPDATE"},
	{"DELETE", "DELETE"}, {"LOAD", "LOAD"},
	{"SET", "SET"}, {"SHOW", "SHOW"}, {"EXPLAIN", "EXPLAIN"}, {"DESC", "EXPLAIN"},
	{"ADMIN", "ADMIN"}, {"GRANT", "DCL"}, {"REVOKE", "DCL"},
	{"BEGIN", "TXN"}, {"START", "TXN"}, {"COMMIT", "TXN"}, {"ROLLBACK", "TXN"},
	{"SAVEPOINT", "TXN"}, {"LOCK", "TXN"}, {"UNLOCK", "TXN"},
	{"PREPARE", "PREPARED"}, {"EXECUTE", "PREPARED"}, {"DEALLOCATE", "PREPARED"},
	{"ANALYZE", "STATS"}, {"FLASHBACK", "ADMIN"}, {"RECOVER", "ADMIN"},
	{"USE", "SET"}, {"CALL", "ROUTINE"}, {"DO", "DO"},
	{"BATCH", "DML-BATCH"}, {"IMPORT", "LOAD"},
	{"CREATE", "CREATE OTHER"},
}

var leadingComment = regexp.MustCompile(`^(\s*(/\*([^*]|\*[^/])*\*/|--[^\n]*\n|#[^\n]*\n))*\s*`)

func classifyFamily(sql string) string {
	s := strings.ToUpper(leadingComment.ReplaceAllString(sql, ""))
	s = strings.TrimSpace(s)
	if s == "" {
		return "UNKNOWN"
	}
	for _, p := range familyPrefixes {
		if strings.HasPrefix(s, p.prefix) {
			return p.family
		}
	}
	return "OTHER"
}

var (
	numRe    = regexp.MustCompile(`\d+`)
	quotedRe = regexp.MustCompile("(`[^`]*`|'[^']*'|\"[^\"]*\")")
	spaceRe  = regexp.MustCompile(`\s+`)
)

// clusterKey normalizes an error message so one grammar divergence maps to
// one cluster: strips positions, numbers, and quoted identifiers.
func clusterKey(msg string) string {
	m := quotedRe.ReplaceAllString(msg, "?")
	m = numRe.ReplaceAllString(m, "N")
	m = spaceRe.ReplaceAllString(strings.TrimSpace(m), " ")
	return m
}

func leadingTokens(sql string, n int) string {
	s := strings.ToUpper(leadingComment.ReplaceAllString(sql, ""))
	fields := strings.Fields(s)
	if len(fields) > n {
		fields = fields[:n]
	}
	return strings.Join(fields, " ")
}

func stmtHash(sql string) string {
	h := sha256.Sum256([]byte(sql))
	return hex.EncodeToString(h[:8])
}
