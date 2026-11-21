package catalog_test

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	omnidriver "github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/catalog"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

// A4 of the ENTERPRISE smoke axes: a SEEDED, deterministic
// random-migration fuzzer over a 30-table Zabbix slice. Each seed derives K ∈ [3,8]
// independent mutations from a menu (add/drop/modify column, add/drop index, add/drop FK
// with valid targets only, add/drop/modify view, add/drop trigger, add/drop routine,
// partition/departition), applies them TEXTUALLY to the canonical dump to build the
// target SDL, and runs the oracle protocol (diff → apply → converge → idempotence).
//
// Determinism: the mutation stream is a pure function of the seed (math/rand with a
// fixed source; the schema model is iterated in sorted order), so a failing seed replays
// exactly. Every failure prints the seed and a ddmin-minimized mutation list, each
// minimization trial re-running the full oracle on a fresh scratch database.
//
// 20 seeds run. Scratch databases are entsdl_-prefixed and dropped eagerly per trial.

// entFuzzSliceTables is the ~30-table fuzz base: the core slice plus the httptest family
// (6 more triggered tables, FK-chained through httptest/httpstep into hosts and items)
// and token (a 2-FK fan-in on users).
var entFuzzSliceTables = append(append([]string{}, entSliceCoreTables...),
	"httptest", "httpstep", "httptestitem", "httpstepitem", "httptest_field", "httpstep_field",
	"token",
)

// entFuzzAux seeds the droppable/modifiable object pool the menu needs: two views, one
// function, one procedure, and a partitioned table (the departition target).
//
//go:embed testdata/sdl_migration/ent_fuzz_aux.sql
var entFuzzAux string

// entFzProtectedTables are never touched by column-level mutations: the changelog
// trigger bodies name changelog's columns outright (not via NEW/OLD), and ent_fz_part
// belongs to the partition mutations.
var entFzProtectedTables = map[string]bool{"changelog": true, "ent_fz_part": true}

// entFzProtectedColumns are referenced by the seeded views (or their modify-view
// extension columns) and must survive column drops.
var entFzProtectedColumns = map[string]bool{
	"users.username": true, "users.name": true,
	"hosts.host": true, "hosts.status": true,
}

// ----------------------------------------------------------------------------
// Canonical-dump schema model (guides valid mutation choices; the mutations themselves
// are textual).
// ----------------------------------------------------------------------------

type entFzIndex struct {
	name   string
	cols   []string
	unique bool
}

type entFzFK struct {
	name string
	cols []string
}

type entFzTable struct {
	name        string
	cols        []string          // ordered
	colLine     map[string]string // name -> trimmed definition line
	pk          map[string]bool
	indexes     []entFzIndex
	fks         []entFzFK
	partitioned bool
}

type entFzModel struct {
	tables      []*entFzTable // sorted by name
	byName      map[string]*entFzTable
	referenced  map[string]bool            // tables referenced by any FK
	triggerCols map[string]map[string]bool // table -> NEW./OLD.-referenced columns
	triggers    []string                   // sorted trigger names
	views       []string                   // seeded, droppable views
	functions   []string                   // seeded, droppable functions
	procedures  []string                   // seeded, droppable procedures
}

var (
	entFzReColLine = regexp.MustCompile("^`(\\w+)` (.+?),?$")
	entFzRePK      = regexp.MustCompile(`^PRIMARY KEY \((.+)\),?$`)
	entFzReIndex   = regexp.MustCompile("^(UNIQUE )?KEY `(\\w+)` \\((.+)\\),?$")
	entFzReFK      = regexp.MustCompile("^CONSTRAINT `(\\w+)` FOREIGN KEY \\((.+?)\\) REFERENCES `(\\w+)`")
	entFzReTrigger = regexp.MustCompile("(?i)^CREATE TRIGGER `(\\w+)`.* ON `(\\w+)`")
	entFzReNewOld  = regexp.MustCompile(`(?i)\b(?:new|old)\.(\w+)`)
	entFzReCol     = regexp.MustCompile("`(\\w+)`")
	entFzReVarLen  = regexp.MustCompile(`varchar\((\d+)\)`)
)

// entFuzzModel parses the canonical dump into the mutation-planning model.
func entFuzzModel(t *testing.T, source string) *entFzModel {
	t.Helper()
	stmts := splitSQLSegments(source)

	m := &entFzModel{
		byName:      map[string]*entFzTable{},
		referenced:  map[string]bool{},
		triggerCols: map[string]map[string]bool{},
	}
	for _, s := range stmts {
		text := strings.TrimSpace(s.Text)
		upper := strings.ToUpper(text)
		switch {
		case strings.HasPrefix(upper, "CREATE TABLE `"):
			tbl := entFzParseTable(text)
			if tbl != nil {
				m.tables = append(m.tables, tbl)
				m.byName[tbl.name] = tbl
			}
		case strings.HasPrefix(upper, "CREATE TRIGGER `") || (strings.HasPrefix(upper, "CREATE ") && strings.Contains(upper, " TRIGGER `")):
			if mm := entFzReTrigger.FindStringSubmatch(text); mm != nil {
				m.triggers = append(m.triggers, mm[1])
				onTable := mm[2]
				if m.triggerCols[onTable] == nil {
					m.triggerCols[onTable] = map[string]bool{}
				}
				for _, ref := range entFzReNewOld.FindAllStringSubmatch(text, -1) {
					m.triggerCols[onTable][ref[1]] = true
				}
			}
		default:
		}
	}
	// FK-referenced tables (for partition candidacy).
	for _, s := range stmts {
		for _, mm := range entFzReFKRefs.FindAllStringSubmatch(s.Text, -1) {
			m.referenced[mm[1]] = true
		}
	}
	// Seeded droppable objects (fixed names — the fuzz base owns them).
	m.views = []string{"ent_fz_v1", "ent_fz_v2"}
	m.functions = []string{"ent_fz_fn1"}
	m.procedures = []string{"ent_fz_pr1"}
	return m
}

var entFzReFKRefs = regexp.MustCompile("FOREIGN KEY \\(.+?\\) REFERENCES `(\\w+)`")

// entFzParseTable parses one canonical CREATE TABLE statement.
func entFzParseTable(stmt string) *entFzTable {
	header := regexp.MustCompile("^CREATE TABLE `(\\w+)`").FindStringSubmatch(stmt)
	if header == nil {
		return nil
	}
	tbl := &entFzTable{
		name:        header[1],
		colLine:     map[string]string{},
		pk:          map[string]bool{},
		partitioned: strings.Contains(strings.ToUpper(stmt), "PARTITION BY"),
	}
	for _, raw := range strings.Split(stmt, "\n")[1:] {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, ")") {
			break
		}
		if mm := entFzRePK.FindStringSubmatch(line); mm != nil {
			for _, c := range entFzReCol.FindAllStringSubmatch(mm[1], -1) {
				tbl.pk[c[1]] = true
			}
			continue
		}
		if mm := entFzReIndex.FindStringSubmatch(line); mm != nil {
			idx := entFzIndex{name: mm[2], unique: mm[1] != ""}
			for _, c := range entFzReCol.FindAllStringSubmatch(mm[3], -1) {
				idx.cols = append(idx.cols, c[1])
			}
			tbl.indexes = append(tbl.indexes, idx)
			continue
		}
		if mm := entFzReFK.FindStringSubmatch(line); mm != nil {
			fk := entFzFK{name: mm[1]}
			for _, c := range entFzReCol.FindAllStringSubmatch(mm[2], -1) {
				fk.cols = append(fk.cols, c[1])
			}
			tbl.fks = append(tbl.fks, fk)
			continue
		}
		if strings.HasPrefix(line, "CONSTRAINT `") || strings.HasPrefix(line, "FULLTEXT KEY") {
			continue // CHECK constraints / fulltext — not mutation targets
		}
		if mm := entFzReColLine.FindStringSubmatch(line); mm != nil {
			tbl.cols = append(tbl.cols, mm[1])
			tbl.colLine[mm[1]] = strings.TrimSuffix(line, ",")
		}
	}
	return tbl
}

// entFzIndexedCols returns every column participating in any index or FK of tbl.
func entFzIndexedCols(tbl *entFzTable) map[string]bool {
	used := map[string]bool{}
	for c := range tbl.pk {
		used[c] = true
	}
	for _, idx := range tbl.indexes {
		for _, c := range idx.cols {
			used[c] = true
		}
	}
	for _, fk := range tbl.fks {
		for _, c := range fk.cols {
			used[c] = true
		}
	}
	return used
}

// ----------------------------------------------------------------------------
// Mutation generation.
// ----------------------------------------------------------------------------

// entFzMutation is one textual schema mutation with a replayable description.
type entFzMutation struct {
	desc  string
	apply func(t *testing.T, source string) string
}

// entFzState tracks per-round conflicts so the K mutations stay independent.
type entFzState struct {
	usedCols     map[string]bool // "table.col" touched by drop/modify
	usedIndexes  map[string]bool // "table.index" dropped
	usedFKs      map[string]bool
	usedViews    map[string]bool
	usedTriggers map[string]bool
	usedRoutines map[string]bool
	lockedTables map[string]bool // partition/departition/fk-structure locks
}

func entFzNewState() *entFzState {
	return &entFzState{
		usedCols:     map[string]bool{},
		usedIndexes:  map[string]bool{},
		usedFKs:      map[string]bool{},
		usedViews:    map[string]bool{},
		usedTriggers: map[string]bool{},
		usedRoutines: map[string]bool{},
		lockedTables: map[string]bool{},
	}
}

var entFzMenu = []string{
	"add_column", "drop_column", "modify_column",
	"add_index", "drop_index",
	"add_fk", "drop_fk",
	"add_view", "drop_view", "modify_view",
	"add_trigger", "drop_trigger",
	"add_routine", "drop_routine",
	"partition", "departition",
}

// entFuzzMutations derives k independent mutations from the model, deterministically in
// (model order, rng stream).
func entFuzzMutations(t *testing.T, model *entFzModel, rng *rand.Rand, k int) []entFzMutation {
	t.Helper()
	state := entFzNewState()
	var muts []entFzMutation
	for attempt := 0; attempt < 400 && len(muts) < k; attempt++ {
		kind := entFzMenu[rng.Intn(len(entFzMenu))]
		seq := len(muts)
		if m, ok := entFzGenerate(kind, model, rng, state, seq); ok {
			muts = append(muts, m)
		}
	}
	require.Len(t, muts, k, "fuzz generator starved (only %d of %d mutations)", len(muts), k)
	return muts
}

// entFzPickTable picks a random mutable (non-protected, non-locked) table.
func entFzPickTable(model *entFzModel, rng *rand.Rand, state *entFzState) *entFzTable {
	var candidates []*entFzTable
	for _, tbl := range model.tables {
		if entFzProtectedTables[tbl.name] || state.lockedTables[tbl.name] {
			continue
		}
		candidates = append(candidates, tbl)
	}
	if len(candidates) == 0 {
		return nil
	}
	return candidates[rng.Intn(len(candidates))]
}

// entFzGenerate builds one mutation of the given kind, or reports it infeasible.
func entFzGenerate(kind string, model *entFzModel, rng *rand.Rand, state *entFzState, seq int) (entFzMutation, bool) {
	switch kind {
	case "add_column":
		tbl := entFzPickTable(model, rng, state)
		if tbl == nil {
			return entFzMutation{}, false
		}
		defs := []string{"varchar(40) DEFAULT NULL", "int NOT NULL DEFAULT '0'", "decimal(8,2) DEFAULT NULL", "bigint unsigned DEFAULT NULL"}
		def := defs[rng.Intn(len(defs))]
		col := fmt.Sprintf("ent_fz_c%d", seq)
		table := tbl.name
		return entFzMutation{
			desc: fmt.Sprintf("add_column %s.%s (%s)", table, col, def),
			apply: func(t *testing.T, source string) string {
				return addColumnToTable(t, source, table, "`"+col+"` "+def)
			},
		}, true

	case "drop_column":
		type cand struct{ table, col string }
		var cands []cand
		for _, tbl := range model.tables {
			if entFzProtectedTables[tbl.name] || state.lockedTables[tbl.name] {
				continue
			}
			indexed := entFzIndexedCols(tbl)
			trigRef := model.triggerCols[tbl.name]
			for _, c := range tbl.cols {
				if indexed[c] || trigRef[c] || entFzProtectedColumns[tbl.name+"."+c] {
					continue
				}
				if state.usedCols[tbl.name+"."+c] {
					continue
				}
				if strings.Contains(tbl.colLine[c], "GENERATED") {
					continue
				}
				cands = append(cands, cand{tbl.name, c})
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.table+"."+pick.col] = true
		return entFzMutation{
			desc: fmt.Sprintf("drop_column %s.%s", pick.table, pick.col),
			apply: func(t *testing.T, source string) string {
				return entDropLineInTable(t, source, pick.table, "`"+pick.col+"`")
			},
		}, true

	case "modify_column":
		type cand struct {
			table, col string
			from, to   int
		}
		var cands []cand
		for _, tbl := range model.tables {
			if entFzProtectedTables[tbl.name] || state.lockedTables[tbl.name] {
				continue
			}
			for _, c := range tbl.cols {
				if state.usedCols[tbl.name+"."+c] {
					continue
				}
				mm := entFzReVarLen.FindStringSubmatch(tbl.colLine[c])
				if mm == nil {
					continue
				}
				n := entFzAtoi(mm[1])
				if n < 8 || n > 150 {
					continue
				}
				cands = append(cands, cand{tbl.name, c, n, n + 37})
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.table+"."+pick.col] = true
		return entFzMutation{
			desc: fmt.Sprintf("modify_column widen %s.%s varchar(%d)->varchar(%d)", pick.table, pick.col, pick.from, pick.to),
			apply: func(t *testing.T, source string) string {
				return entReplaceInTable(t, source, pick.table,
					fmt.Sprintf("`%s` varchar(%d)", pick.col, pick.from),
					fmt.Sprintf("`%s` varchar(%d)", pick.col, pick.to))
			},
		}, true

	case "add_index":
		type cand struct{ table, col string }
		var cands []cand
		for _, tbl := range model.tables {
			if entFzProtectedTables[tbl.name] || state.lockedTables[tbl.name] {
				continue
			}
			for _, c := range tbl.cols {
				line := tbl.colLine[c]
				rest := strings.TrimPrefix(line, "`"+c+"` ")
				ok := strings.HasPrefix(rest, "bigint") || strings.HasPrefix(rest, "int") || strings.HasPrefix(rest, "integer")
				if mm := entFzReVarLen.FindStringSubmatch(rest); mm != nil && strings.HasPrefix(rest, "varchar") {
					ok = entFzAtoi(mm[1]) <= 150
				}
				if ok {
					cands = append(cands, cand{tbl.name, c})
				}
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		idx := fmt.Sprintf("ent_fz_i%d", seq)
		return entFzMutation{
			desc: fmt.Sprintf("add_index %s on %s(%s)", idx, pick.table, pick.col),
			apply: func(t *testing.T, source string) string {
				return addIndexToTable(t, source, pick.table, "KEY `"+idx+"` (`"+pick.col+"`)")
			},
		}, true

	case "drop_index":
		type cand struct{ table, index string }
		var cands []cand
		for _, tbl := range model.tables {
			if state.lockedTables[tbl.name] {
				continue
			}
			fkCols := map[string]bool{}
			for _, fk := range tbl.fks {
				for _, c := range fk.cols {
					fkCols[c] = true
				}
			}
			for _, idx := range tbl.indexes {
				if state.usedIndexes[tbl.name+"."+idx.name] {
					continue
				}
				overlap := false
				for _, c := range idx.cols {
					if fkCols[c] {
						overlap = true
						break
					}
				}
				if !overlap {
					cands = append(cands, cand{tbl.name, idx.name})
				}
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedIndexes[pick.table+"."+pick.index] = true
		return entFzMutation{
			desc: fmt.Sprintf("drop_index %s.%s", pick.table, pick.index),
			apply: func(t *testing.T, source string) string {
				return entDropLineInTable(t, source, pick.table, "`"+pick.index+"`")
			},
		}, true

	case "add_fk":
		src := entFzPickTable(model, rng, state)
		if src == nil || src.partitioned {
			return entFzMutation{}, false
		}
		var refs []*entFzTable
		for _, tbl := range model.tables {
			if tbl.partitioned || state.lockedTables[tbl.name] || len(tbl.pk) != 1 {
				continue
			}
			pkCol := ""
			for c := range tbl.pk {
				pkCol = c
			}
			if strings.HasPrefix(strings.TrimPrefix(tbl.colLine[pkCol], "`"+pkCol+"` "), "bigint") &&
				strings.Contains(tbl.colLine[pkCol], "unsigned") {
				refs = append(refs, tbl)
			}
		}
		if len(refs) == 0 {
			return entFzMutation{}, false
		}
		ref := refs[rng.Intn(len(refs))]
		refPK := ""
		for c := range ref.pk {
			refPK = c
		}
		srcName, refName := src.name, ref.name
		state.lockedTables[srcName] = true
		state.lockedTables[refName] = true
		col := fmt.Sprintf("ent_fz_r%d", seq)
		fk := fmt.Sprintf("ent_fz_fk%d", seq)
		return entFzMutation{
			desc: fmt.Sprintf("add_fk %s.%s -> %s.%s (%s)", srcName, col, refName, refPK, fk),
			apply: func(t *testing.T, source string) string {
				s := addColumnToTable(t, source, srcName, "`"+col+"` bigint unsigned DEFAULT NULL")
				s = addIndexToTable(t, s, srcName, "KEY `"+fk+"` (`"+col+"`)")
				return addIndexToTable(t, s, srcName,
					"CONSTRAINT `"+fk+"` FOREIGN KEY (`"+col+"`) REFERENCES `"+refName+"` (`"+refPK+"`) ON DELETE SET NULL")
			},
		}, true

	case "drop_fk":
		type cand struct{ table, fk string }
		var cands []cand
		for _, tbl := range model.tables {
			if state.lockedTables[tbl.name] {
				continue
			}
			for _, fk := range tbl.fks {
				if !state.usedFKs[tbl.name+"."+fk.name] {
					cands = append(cands, cand{tbl.name, fk.name})
				}
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedFKs[pick.table+"."+pick.fk] = true
		return entFzMutation{
			desc: fmt.Sprintf("drop_fk %s.%s", pick.table, pick.fk),
			apply: func(t *testing.T, source string) string {
				return entDropLineInTable(t, source, pick.table, "CONSTRAINT `"+pick.fk+"`")
			},
		}, true

	case "add_view":
		tbl := entFzPickTable(model, rng, state)
		if tbl == nil || len(tbl.pk) == 0 {
			return entFzMutation{}, false
		}
		pkCol := ""
		for _, c := range tbl.cols {
			if tbl.pk[c] {
				pkCol = c
				break
			}
		}
		if pkCol == "" {
			return entFzMutation{}, false
		}
		name := fmt.Sprintf("ent_fz_v%d", 10+seq)
		table := tbl.name
		return entFzMutation{
			desc: fmt.Sprintf("add_view %s over %s(%s)", name, table, pkCol),
			apply: func(_ *testing.T, source string) string {
				return source + fmt.Sprintf("\nCREATE VIEW %s AS SELECT %s FROM %s;\n", name, pkCol, table)
			},
		}, true

	case "drop_view":
		var cands []string
		for _, v := range model.views {
			if !state.usedViews[v] {
				cands = append(cands, v)
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedViews[pick] = true
		return entFzMutation{
			desc: "drop_view " + pick,
			apply: func(t *testing.T, source string) string {
				return dropObjectBlock(t, source, "VIEW", pick)
			},
		}, true

	case "modify_view":
		extension := map[string]string{
			"ent_fz_v1": "`users`.`name` AS `name`",
			"ent_fz_v2": "`hosts`.`status` AS `status`",
		}
		var cands []string
		for _, v := range model.views {
			if !state.usedViews[v] && extension[v] != "" {
				cands = append(cands, v)
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedViews[pick] = true
		ext := extension[pick]
		return entFzMutation{
			desc: "modify_view " + pick,
			apply: func(t *testing.T, source string) string {
				return addColumnToView(t, source, pick, ext)
			},
		}, true

	case "add_trigger":
		tbl := entFzPickTable(model, rng, state)
		if tbl == nil || len(tbl.pk) == 0 {
			return entFzMutation{}, false
		}
		pkCol := ""
		for _, c := range tbl.cols {
			if tbl.pk[c] {
				pkCol = c
				break
			}
		}
		if pkCol == "" {
			return entFzMutation{}, false
		}
		name := fmt.Sprintf("ent_fz_t%d", seq)
		table := tbl.name
		return entFzMutation{
			desc: fmt.Sprintf("add_trigger %s on %s", name, table),
			apply: func(_ *testing.T, source string) string {
				return source + fmt.Sprintf("\nCREATE TRIGGER %s BEFORE UPDATE ON %s FOR EACH ROW SET NEW.%s = NEW.%s;\n",
					name, table, pkCol, pkCol)
			},
		}, true

	case "drop_trigger":
		var cands []string
		for _, trg := range model.triggers {
			if !state.usedTriggers[trg] {
				cands = append(cands, trg)
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedTriggers[pick] = true
		return entFzMutation{
			desc: "drop_trigger " + pick,
			apply: func(t *testing.T, source string) string {
				return dropTrigger(t, source, pick)
			},
		}, true

	case "add_routine":
		if rng.Intn(2) == 0 {
			name := fmt.Sprintf("ent_fz_fn%d", 10+seq)
			ret := seq + 100
			return entFzMutation{
				desc: "add_routine function " + name,
				apply: func(_ *testing.T, source string) string {
					return source + fmt.Sprintf("\nCREATE FUNCTION %s() RETURNS INT DETERMINISTIC RETURN %d;\n", name, ret)
				},
			}, true
		}
		name := fmt.Sprintf("ent_fz_pr%d", 10+seq)
		sel := seq + 200
		return entFzMutation{
			desc: "add_routine procedure " + name,
			apply: func(_ *testing.T, source string) string {
				return source + fmt.Sprintf("\nCREATE PROCEDURE %s() BEGIN SELECT %d; END;\n", name, sel)
			},
		}, true

	case "drop_routine":
		type cand struct{ kind, name string }
		var cands []cand
		for _, fn := range model.functions {
			if !state.usedRoutines[fn] {
				cands = append(cands, cand{"FUNCTION", fn})
			}
		}
		for _, pr := range model.procedures {
			if !state.usedRoutines[pr] {
				cands = append(cands, cand{"PROCEDURE", pr})
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedRoutines[pick.name] = true
		return entFzMutation{
			desc: "drop_routine " + strings.ToLower(pick.kind) + " " + pick.name,
			apply: func(t *testing.T, source string) string {
				return dropObjectBlock(t, source, pick.kind, pick.name)
			},
		}, true

	case "partition":
		var cands []*entFzTable
		for _, tbl := range model.tables {
			// changelog is column-protected but IS the partition candidate (the only
			// FK-free, unreferenced, single-PK table in the slice); ent_fz_part is
			// already partitioned.
			if tbl.partitioned || state.lockedTables[tbl.name] {
				continue
			}
			if len(tbl.fks) > 0 || model.referenced[tbl.name] || len(tbl.pk) != 1 {
				continue
			}
			// The partition key must be part of every unique key: single-col-PK tables
			// with no other UNIQUE index qualify.
			hasUnique := false
			for _, idx := range tbl.indexes {
				if idx.unique {
					hasUnique = true
					break
				}
			}
			if hasUnique {
				continue
			}
			cands = append(cands, tbl)
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		pkCol := ""
		for c := range pick.pk {
			pkCol = c
		}
		state.lockedTables[pick.name] = true
		table := pick.name
		return entFzMutation{
			desc: fmt.Sprintf("partition %s by hash(%s)", table, pkCol),
			apply: func(t *testing.T, source string) string {
				return entAppendTableClause(t, source, table, "PARTITION BY HASH (`"+pkCol+"`)\nPARTITIONS 4")
			},
		}, true

	case "departition":
		var cands []*entFzTable
		for _, tbl := range model.tables {
			if tbl.partitioned && !state.lockedTables[tbl.name] {
				cands = append(cands, tbl)
			}
		}
		if len(cands) == 0 {
			return entFzMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.lockedTables[pick.name] = true
		table := pick.name
		return entFzMutation{
			desc: "departition " + table,
			apply: func(t *testing.T, source string) string {
				return entSetPartitionClause(t, source, table, "")
			},
		}, true

	default:
		return entFzMutation{}, false
	}
}

// entFzAtoi is a no-error atoi for regex-captured digits.
func entFzAtoi(s string) int {
	n := 0
	for _, r := range s {
		n = n*10 + int(r-'0')
	}
	return n
}

// ----------------------------------------------------------------------------
// Trial runner + ddmin minimization.
// ----------------------------------------------------------------------------

// entFuzzScratchDB creates an eagerly-droppable scratch database (fuzz + minimization
// churn through many databases inside a single test, so cleanup cannot wait for test
// end).
func entFuzzScratchDB(ctx context.Context, t *testing.T, srv *liveServer) (string, func()) {
	t.Helper()
	dbName := newLiveDatabase(ctx, t, srv, "entsdl_fz")
	drop := func() {
		testcontainer.Reset(t, testcontainer.Default, func(ctx context.Context, db *sql.DB) error {
			_, err := db.ExecContext(ctx, "DROP DATABASE IF EXISTS `"+dbName+"`")
			return err
		})
	}
	return dbName, drop
}

// entFuzzTrial runs the whole oracle protocol for one mutation set on a fresh scratch
// database. Infrastructure failures (server down, base DDL broken, mutation anchors
// missing) fail the test hard; ORACLE violations (diff error, apply error, residual
// after apply, non-idempotent dump) are returned as errors so the caller can minimize.
func entFuzzTrial(ctx context.Context, t *testing.T, srv *liveServer, baseDDL string, muts []entFzMutation) error {
	t.Helper()
	dbName, drop := entFuzzScratchDB(ctx, t, srv)
	defer drop()

	driver, err := createLiveMySQLDriver(ctx, srv, dbName)
	require.NoError(t, err)
	_, err = driver.Execute(ctx, entNormalizeDelimiters(baseDDL), omnidriver.ExecuteOptions{})
	driver.Close()
	require.NoError(t, err, "[%s] fuzz base failed to load", srv.name)

	source := dumpSDL(ctx, t, srv, dbName)
	target := source
	for _, m := range muts {
		target = m.apply(t, target)
	}

	plan, err := mysqlDiffSDLMigration(source, target, srv.version)
	if err != nil {
		return fmt.Errorf("diff failed: %w", err)
	}
	if strings.TrimSpace(plan) == "" {
		return errors.New("empty plan for a non-empty mutation set")
	}
	if err := applyDDL(ctx, t, srv, dbName, plan); err != nil {
		return fmt.Errorf("plan failed to apply; plan was:\n%s: %w", plan, err)
	}
	after := dumpSDL(ctx, t, srv, dbName)
	converge, err := mysqlDiffSDLMigration(after, target, srv.version)
	if err != nil {
		return fmt.Errorf("converge diff failed: %w", err)
	}
	if converge != "" {
		return fmt.Errorf("did not converge; residual:\n%s\nplan was:\n%s", converge, plan)
	}
	self, err := mysqlDiffSDLMigration(after, after, srv.version)
	if err != nil {
		return fmt.Errorf("idempotence diff failed: %w", err)
	}
	if self != "" {
		return fmt.Errorf("post-apply dump not idempotent:\n%s", self)
	}
	return nil
}

// entFuzzMinimize greedily shrinks a failing mutation set (ddmin-lite: drop one mutation
// at a time, keep the removal whenever the rest still fails). Returns the minimized set
// and its error.
func entFuzzMinimize(ctx context.Context, t *testing.T, srv *liveServer, baseDDL string, muts []entFzMutation, firstErr error) ([]entFzMutation, error) {
	t.Helper()
	minimized := muts
	lastErr := firstErr
	for i := 0; i < len(minimized) && len(minimized) > 1; {
		trial := make([]entFzMutation, 0, len(minimized)-1)
		trial = append(trial, minimized[:i]...)
		trial = append(trial, minimized[i+1:]...)
		if err := entFuzzTrial(ctx, t, srv, baseDDL, trial); err != nil {
			minimized = trial
			lastErr = err
			continue // same index now names the next mutation
		}
		i++
	}
	return minimized, lastErr
}

func entFzDescs(muts []entFzMutation) []string {
	out := make([]string, len(muts))
	for i, m := range muts {
		out[i] = m.desc
	}
	return out
}

// ----------------------------------------------------------------------------
// The fuzz test: 20 seeds (deterministic; the seed is in the subtest name).
// ----------------------------------------------------------------------------

func TestSDLEnterpriseFuzz(t *testing.T) {
	ctx := context.Background()

	for _, srv := range tapeOracleServers(t) {
		srv := srv
		const seedCount = 20
		t.Run(srv.name, func(t *testing.T) {
			baseDDL := entZabbixSlice(t, entFuzzSliceTables...) + entFuzzAux

			// The model is built once per server from a reference load (the canonical dump
			// is deterministic, so it matches every trial's source).
			refDB, refDrop := entFuzzScratchDB(ctx, t, srv)
			driver, err := createLiveMySQLDriver(ctx, srv, refDB)
			require.NoError(t, err)
			_, err = driver.Execute(ctx, entNormalizeDelimiters(baseDDL), omnidriver.ExecuteOptions{})
			driver.Close()
			require.NoError(t, err, "[%s] fuzz base failed to load", srv.name)
			refSource := dumpSDL(ctx, t, srv, refDB)
			refDrop()
			model := entFuzzModel(t, refSource)
			require.NotEmpty(t, model.tables, "fuzz model parsed no tables")
			require.NotEmpty(t, model.triggers, "fuzz model parsed no triggers")

			for seed := int64(1); seed <= int64(seedCount); seed++ {
				seed := seed
				t.Run(fmt.Sprintf("seed_%d", seed), func(t *testing.T) {
					rng := rand.New(rand.NewSource(seed))
					k := 3 + rng.Intn(6) // K ∈ [3,8]
					muts := entFuzzMutations(t, model, rng, k)
					t.Logf("[FUZZ %s seed=%d] K=%d mutations:\n  %s", srv.name, seed, k, strings.Join(entFzDescs(muts), "\n  "))

					err := entFuzzTrial(ctx, t, srv, baseDDL, muts)
					if err == nil {
						return
					}
					minimized, minErr := entFuzzMinimize(ctx, t, srv, baseDDL, muts, err)
					t.Errorf("[FUZZ %s seed=%d] FAILED\nminimized mutations (%d of %d):\n  %s\nerror:\n%v",
						srv.name, seed, len(minimized), len(muts), strings.Join(entFzDescs(minimized), "\n  "), minErr)
				})
			}
		})
	}
}

// A6 of the ENTERPRISE smoke axes: a SEEDED, deterministic
// MULTI-ROUND STATEFUL fuzzer. The single-round fuzzer (A4)
// mutates a pristine dump once; real projects mutate the LAST dump repeatedly, so drift
// and text-accretion bugs that only appear across dump->mutate->apply->dump cycles are
// invisible to it (real example this campaign: omni used to accrete one paren layer per
// dump cycle around subquery wrappers until the fold fixed it — a per-round differ never
// sees the growth).
//
// Protocol per seed, on a 30-table Zabbix slice base:
//
//	dump_0 = canonical dump of the freshly loaded base
//	for round r in 1..5:
//	    target_r = mutate(dump_{r-1}, K ∈ [2,5] mutations from the menu)
//	    oracle   = diff(dump_{r-1} -> target_r) → apply → converge → idempotence
//	    dump_r   = fresh canonical dump, and ADDITIONALLY:
//	      (a) dump_r reloads through the STRICT LoadSDL path (no LoadSQL fallback mask),
//	      (b) len(dump_r) stays under len(dump_0) + the cumulative added-object budget
//	          (linear ledger — catches unbounded text accretion), and every statement NOT
//	          touched by this round's mutations is BYTE-IDENTICAL to its dump_{r-1} form
//	          (catches per-cycle accretion and collateral churn on untouched objects),
//	      (c) object counts (tables/views/routines/triggers/columns/indexes/FKs/checks)
//	          match the mutation ledger exactly.
//
// The menu is the A4 menu PLUS: add/drop CHECK constraint, add/drop FULLTEXT index,
// add/drop SPATIAL column+index with SRID, add/modify/drop STORED and VIRTUAL generated
// columns, toggle column/index INVISIBLE, and ENUM member append.
//
// Determinism: the mutation stream is a pure function of the seed (the seed is in the
// subtest name); the schema model is parsed from the previous round's dump in statement
// order, so candidate lists are stable. Oracle violations are ddmin-minimized (one
// mutation at a time, each trial reloading dump_{r-1} into a fresh scratch database) and
// reported with seed + round + minimized ledger. Scratch databases are entsdl_-prefixed
// and dropped eagerly.
//
// 12 seeds run.
//
// Shared helpers (liveServers, createLiveMySQLDriver, newLiveDatabase, dumpSDL, applyDDL,
// syncMetaForDB, objectCounts, addColumnToTable, addIndexToTable, addColumnToView,
// dropObjectBlock, dropTrigger, entZabbixSlice, entFuzzSliceTables, entFuzzAux,
// entNormalizeDelimiters, entDropLineInTable, entReplaceLineInTable, entAppendTableClause,
// entSetPartitionClause, entFzProtectedTables, entFzProtectedColumns, entFuzzModel,
// entFzIndexedCols, entFzState, entFzMutation, withDatabaseContext) come from the sibling
// files in this package.

const (
	entSfRounds = 5
	// entSfRoundSlack is the per-round allowance on top of the per-mutation text budgets
	// (absorbs benign re-rendering, e.g. a DEFAULT NULL a fresh dump makes explicit).
	// Deliberately small so systematic accretion across 30+ objects trips the guard.
	entSfRoundSlack = 400
)

// ----------------------------------------------------------------------------
// Stateful-fuzz aux schema: seeds the droppable/modifiable object pool for the NEW menu
// kinds. entSfAux80 seeds the SRID spatial, INVISIBLE, and CHECK kinds.
// ----------------------------------------------------------------------------

//go:embed testdata/sdl_migration/ent_sf_aux_common.sql
var entSfAuxCommon string

//go:embed testdata/sdl_migration/ent_sf_aux_80.sql
var entSfAux80 string

func entSfBaseDDL(t *testing.T) string {
	t.Helper()
	return entZabbixSlice(t, entFuzzSliceTables...) + entFuzzAux + entSfAuxCommon + entSfAux80
}

// ----------------------------------------------------------------------------
// Object-count ledger.
// ----------------------------------------------------------------------------

// entSfCounts is the schema-shape summary compared against the mutation ledger after
// every round. Comparable (all ints) so `!=` works.
type entSfCounts struct {
	tables, views, functions, procedures, triggers int
	columns, indexes, fks, checks                  int
}

func (c entSfCounts) plus(d entSfCounts) entSfCounts {
	return entSfCounts{
		tables: c.tables + d.tables, views: c.views + d.views,
		functions: c.functions + d.functions, procedures: c.procedures + d.procedures,
		triggers: c.triggers + d.triggers, columns: c.columns + d.columns,
		indexes: c.indexes + d.indexes, fks: c.fks + d.fks, checks: c.checks + d.checks,
	}
}

func entSfCountObjects(meta *metadatapb.DatabaseSchemaMetadata) entSfCounts {
	c := entSfCounts{}
	c.tables, c.views, c.functions, c.procedures, c.triggers = objectCounts(meta)
	proto := meta
	if proto == nil {
		return c
	}
	for _, sm := range proto.GetSchemas() {
		for _, tbl := range sm.GetTables() {
			c.columns += len(tbl.GetColumns())
			c.indexes += len(tbl.GetIndexes())
			c.fks += len(tbl.GetForeignKeys())
			c.checks += len(tbl.GetCheckConstraints())
		}
	}
	return c
}

// entSfMutation wraps the replayable A4 mutation with the ledger metadata the stateful
// protocol needs: the expected object-count delta, the dump-text growth budget, and the
// statement keys the mutation may legitimately rewrite (everything else must stay
// byte-identical across the round).
type entSfMutation struct {
	entFzMutation
	delta   entSfCounts
	budget  int
	touched []string
}

// entSfLedger is the cross-round state: expected counts, and the once-only guards for
// mutations that cannot repeat on the same object across rounds.
type entSfLedger struct {
	expect       entSfCounts
	viewModified map[string]bool
}

// ----------------------------------------------------------------------------
// Extended schema model: everything the A4 model parses, plus the constructs the new
// menu mutates (FULLTEXT/SPATIAL keys, CHECK constraints, generated columns, INVISIBLE
// columns/indexes, ENUM columns) and the referential protections they imply.
// ----------------------------------------------------------------------------

type entSfIdxRec struct {
	table, name, line string
	cols              []string
	invisible         bool
	fkOverlap         bool
}

type entSfNamedRec struct{ table, name string }

type entSfColRec struct{ table, col string }

type entSfGenRec struct{ table, name, line string }

type entSfSchema struct {
	base       *entFzModel
	views      []string
	functions  []string
	procedures []string
	idx        []entSfIdxRec // non-PK B-tree indexes, visible AND invisible
	fulltext   []entSfNamedRec
	spatial    []entSfIdxRec // cols[0] is the single spatial column
	checks     []entSfNamedRec
	genCols    []entSfGenRec
	enumCols   []entSfColRec
	invisCols  []entSfColRec
	// visibleCols counts non-INVISIBLE columns per table (a toggle must leave >= 1).
	visibleCols map[string]int
	// protected are "table.col" keys drop/modify must not touch beyond the base model's
	// own exclusions: fulltext/spatial/invisible-index parts, generated-expression and
	// CHECK-expression references.
	protected map[string]bool
	// noPartition marks tables that cannot gain a partition clause (FULLTEXT/SPATIAL
	// indexes are unsupported on partitioned tables).
	noPartition map[string]bool
}

var (
	entSfReIdxLine  = regexp.MustCompile("^(UNIQUE )?KEY `(\\w+)` \\((.+?)\\)( /\\*!80000 INVISIBLE \\*/)?,?$")
	entSfReFtLine   = regexp.MustCompile("^FULLTEXT KEY `(\\w+)` \\((.+?)\\),?$")
	entSfReSpLine   = regexp.MustCompile("^SPATIAL KEY `(\\w+)` \\(`(\\w+)`\\),?$")
	entSfReChkLine  = regexp.MustCompile("^CONSTRAINT `(\\w+)` CHECK \\((.+)\\),?$")
	entSfReColLine  = regexp.MustCompile("^`(\\w+)` (.+?),?$")
	entSfReViewHdr  = regexp.MustCompile("(?i)^CREATE .*?VIEW `(\\w+)`")
	entSfReFnHdr    = regexp.MustCompile("(?i)^CREATE .*?FUNCTION `(\\w+)`")
	entSfRePrHdr    = regexp.MustCompile("(?i)^CREATE .*?PROCEDURE `(\\w+)`")
	entSfReGenConst = regexp.MustCompile(`\+ (\d+)\)`)
	entSfReEnumSpec = regexp.MustCompile(`enum\(([^)]*)\)`)
)

const (
	entSfInvisibleColToken = " /*!80023 INVISIBLE */"
	entSfInvisibleIdxToken = " /*!80000 INVISIBLE */"
)

// entSfParse builds the extended model from a canonical dump.
func entSfParse(t *testing.T, source string) *entSfSchema {
	t.Helper()
	m := &entSfSchema{
		base:        entFuzzModel(t, source),
		visibleCols: map[string]int{},
		protected:   map[string]bool{},
		noPartition: map[string]bool{},
	}
	stmts := splitSQLSegments(source)
	for _, s := range stmts {
		text := strings.TrimSpace(s.Text)
		upper := strings.ToUpper(text)
		switch {
		case strings.HasPrefix(upper, "CREATE TABLE `"):
			m.parseTableExtras(text)
		case entSfReViewHdr.MatchString(text):
			m.views = append(m.views, entSfReViewHdr.FindStringSubmatch(text)[1])
		case entSfReFnHdr.MatchString(text):
			m.functions = append(m.functions, entSfReFnHdr.FindStringSubmatch(text)[1])
		case entSfRePrHdr.MatchString(text):
			m.procedures = append(m.procedures, entSfRePrHdr.FindStringSubmatch(text)[1])
		default:
		}
	}
	return m
}

func (m *entSfSchema) parseTableExtras(stmt string) {
	header := entReTable.FindStringSubmatch(stmt)
	if header == nil {
		return
	}
	table := header[1]
	fkCols := map[string]bool{}
	if tbl := m.base.byName[table]; tbl != nil {
		for _, fk := range tbl.fks {
			for _, c := range fk.cols {
				fkCols[c] = true
			}
		}
	}
	for _, raw := range strings.Split(stmt, "\n")[1:] {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, ")") {
			break
		}
		if mm := entSfReFtLine.FindStringSubmatch(line); mm != nil {
			m.fulltext = append(m.fulltext, entSfNamedRec{table: table, name: mm[1]})
			m.noPartition[table] = true
			for _, c := range entFzReCol.FindAllStringSubmatch(mm[2], -1) {
				m.protected[table+"."+c[1]] = true
			}
			continue
		}
		if mm := entSfReSpLine.FindStringSubmatch(line); mm != nil {
			m.spatial = append(m.spatial, entSfIdxRec{table: table, name: mm[1], cols: []string{mm[2]}})
			m.noPartition[table] = true
			m.protected[table+"."+mm[2]] = true
			continue
		}
		if mm := entSfReChkLine.FindStringSubmatch(line); mm != nil {
			m.checks = append(m.checks, entSfNamedRec{table: table, name: mm[1]})
			for _, c := range entFzReCol.FindAllStringSubmatch(mm[2], -1) {
				m.protected[table+"."+c[1]] = true
			}
			continue
		}
		if mm := entSfReIdxLine.FindStringSubmatch(line); mm != nil {
			rec := entSfIdxRec{table: table, name: mm[2], line: strings.TrimSuffix(line, ","), invisible: mm[4] != ""}
			for _, c := range entFzReCol.FindAllStringSubmatch(mm[3], -1) {
				rec.cols = append(rec.cols, c[1])
				if rec.invisible {
					m.protected[table+"."+c[1]] = true
				}
				if fkCols[c[1]] {
					rec.fkOverlap = true
				}
			}
			m.idx = append(m.idx, rec)
			continue
		}
		if mm := entSfReColLine.FindStringSubmatch(line); mm != nil {
			col, def := mm[1], strings.TrimSuffix(line, ",")
			if genIdx := strings.Index(def, "GENERATED ALWAYS AS"); genIdx >= 0 {
				m.genCols = append(m.genCols, entSfGenRec{table: table, name: col, line: def})
				for _, c := range entFzReCol.FindAllStringSubmatch(def[genIdx:], -1) {
					m.protected[table+"."+c[1]] = true
				}
			}
			if strings.HasPrefix(strings.TrimPrefix(def, "`"+col+"` "), "enum(") {
				m.enumCols = append(m.enumCols, entSfColRec{table: table, col: col})
			}
			if strings.HasSuffix(def, strings.TrimSpace(entSfInvisibleColToken)) {
				m.invisCols = append(m.invisCols, entSfColRec{table: table, col: col})
			} else {
				m.visibleCols[table]++
			}
		}
	}
}

// colDef returns the trimmed definition line of table.col from the base model.
func (m *entSfSchema) colDef(table, col string) string {
	if tbl := m.base.byName[table]; tbl != nil {
		return tbl.colLine[col]
	}
	return ""
}

// ----------------------------------------------------------------------------
// Mutation generation.
// ----------------------------------------------------------------------------

var entSfMenu = []string{
	"add_column", "drop_column", "modify_column",
	"add_index", "drop_index",
	"add_fk", "drop_fk",
	"add_view", "drop_view", "modify_view",
	"add_trigger", "drop_trigger",
	"add_routine", "drop_routine",
	"partition", "departition",
	"add_fulltext", "drop_fulltext",
	"add_gen_stored", "add_gen_virtual", "modify_gen", "drop_gen",
	"enum_append",
	"add_check", "drop_check",
	"add_spatial", "drop_spatial",
	"toggle_col_invisible", "toggle_idx_invisible",
}

// entSfMutations derives k independent mutations for round r, deterministically in
// (model order, rng stream).
func entSfMutations(t *testing.T, m *entSfSchema, rng *rand.Rand, ledger *entSfLedger, r, k int) []entSfMutation {
	t.Helper()
	state := entFzNewState()
	var muts []entSfMutation
	for attempt := 0; attempt < 500 && len(muts) < k; attempt++ {
		kind := entSfMenu[rng.Intn(len(entSfMenu))]
		if mut, ok := entSfGenerate(t, kind, m, rng, state, ledger, r, len(muts)); ok {
			muts = append(muts, mut)
		}
	}
	require.Len(t, muts, k, "stateful-fuzz generator starved (only %d of %d mutations, round %d)", len(muts), k, r)
	return muts
}

// entSfMutableTables returns the non-protected, non-locked tables.
func entSfMutableTables(m *entSfSchema, state *entFzState) []*entFzTable {
	var out []*entFzTable
	for _, tbl := range m.base.tables {
		if entFzProtectedTables[tbl.name] || state.lockedTables[tbl.name] {
			continue
		}
		out = append(out, tbl)
	}
	return out
}

// entSfColCandidates lists table.col pairs whose (dump-form) definition rest satisfies
// accept, excluding protected/used/locked/generated columns.
// entSfIsFKMemberCol reports whether table.col participates in any child-side
// foreign key of the base model. Used to keep STORED generated columns off FK
// member bases (the server refuses that shape).
func entSfIsFKMemberCol(m *entSfSchema, table, col string) bool {
	tbl := m.base.byName[table]
	if tbl == nil {
		return false
	}
	for _, fk := range tbl.fks {
		for _, c := range fk.cols {
			if c == col {
				return true
			}
		}
	}
	return false
}

func entSfColCandidates(m *entSfSchema, state *entFzState, accept func(rest string) bool) []entSfColRec {
	var out []entSfColRec
	for _, tbl := range entSfMutableTables(m, state) {
		for _, c := range tbl.cols {
			key := tbl.name + "." + c
			if state.usedCols[key] || m.protected[key] || entFzProtectedColumns[key] {
				continue
			}
			def := tbl.colLine[c]
			if strings.Contains(def, "GENERATED ALWAYS AS") {
				continue
			}
			rest := strings.TrimPrefix(def, "`"+c+"` ")
			if accept(rest) {
				out = append(out, entSfColRec{table: tbl.name, col: c})
			}
		}
	}
	return out
}

func entSfIsIntRest(rest string) bool {
	return strings.HasPrefix(rest, "bigint") || strings.HasPrefix(rest, "int") || strings.HasPrefix(rest, "integer")
}

// entSfGenerate builds one mutation of the given kind, or reports it infeasible. seq
// namespaces generated object names within the round; r namespaces across rounds.
//
//nolint:gocyclo
func entSfGenerate(t *testing.T, kind string, m *entSfSchema, rng *rand.Rand, state *entFzState, ledger *entSfLedger, r, seq int) (entSfMutation, bool) {
	t.Helper()
	name := func(tag string) string { return fmt.Sprintf("ent_sf_r%d_%s%d", r, tag, seq) }
	tableKey := func(tbl string) []string { return []string{"table:" + tbl} }

	switch kind {
	case "add_column":
		tables := entSfMutableTables(m, state)
		if len(tables) == 0 {
			return entSfMutation{}, false
		}
		tbl := tables[rng.Intn(len(tables))].name
		defs := []string{"varchar(40) DEFAULT NULL", "int NOT NULL DEFAULT '0'", "decimal(8,2) DEFAULT NULL", "bigint unsigned DEFAULT NULL"}
		def := defs[rng.Intn(len(defs))]
		col := name("c")
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("add_column %s.%s (%s)", tbl, col, def),
				apply: func(t *testing.T, source string) string {
					return addColumnToTable(t, source, tbl, "`"+col+"` "+def)
				},
			},
			delta: entSfCounts{columns: 1}, budget: 150, touched: tableKey(tbl),
		}, true

	case "drop_column":
		var cands []entSfColRec
		for _, tbl := range entSfMutableTables(m, state) {
			indexed := entFzIndexedCols(tbl)
			trigRef := m.base.triggerCols[tbl.name]
			for _, c := range tbl.cols {
				key := tbl.name + "." + c
				if indexed[c] || trigRef[c] || entFzProtectedColumns[key] || m.protected[key] || state.usedCols[key] {
					continue
				}
				if strings.Contains(tbl.colLine[c], "GENERATED") {
					continue
				}
				cands = append(cands, entSfColRec{table: tbl.name, col: c})
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.table+"."+pick.col] = true
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("drop_column %s.%s", pick.table, pick.col),
				apply: func(t *testing.T, source string) string {
					return entDropLineInTable(t, source, pick.table, "`"+pick.col+"`")
				},
			},
			delta: entSfCounts{columns: -1}, touched: tableKey(pick.table),
		}, true

	case "modify_column":
		type cand struct {
			rec      entSfColRec
			from, to int
		}
		var cands []cand
		for _, tbl := range entSfMutableTables(m, state) {
			for _, c := range tbl.cols {
				key := tbl.name + "." + c
				if state.usedCols[key] || m.protected[key] {
					continue
				}
				mm := entFzReVarLen.FindStringSubmatch(tbl.colLine[c])
				if mm == nil || strings.Contains(tbl.colLine[c], "GENERATED") {
					continue
				}
				n := entFzAtoi(mm[1])
				if n < 8 || n > 150 {
					continue
				}
				cands = append(cands, cand{rec: entSfColRec{table: tbl.name, col: c}, from: n, to: n + 23})
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.rec.table+"."+pick.rec.col] = true
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("modify_column widen %s.%s varchar(%d)->varchar(%d)", pick.rec.table, pick.rec.col, pick.from, pick.to),
				apply: func(t *testing.T, source string) string {
					return entReplaceInTable(t, source, pick.rec.table,
						fmt.Sprintf("`%s` varchar(%d)", pick.rec.col, pick.from),
						fmt.Sprintf("`%s` varchar(%d)", pick.rec.col, pick.to))
				},
			},
			budget: 80, touched: tableKey(pick.rec.table),
		}, true

	case "add_index":
		cands := entSfColCandidates(m, state, func(rest string) bool {
			if entSfIsIntRest(rest) {
				return true
			}
			mm := entFzReVarLen.FindStringSubmatch(rest)
			return mm != nil && strings.HasPrefix(rest, "varchar") && entFzAtoi(mm[1]) <= 150
		})
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		idx := name("i")
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("add_index %s on %s(%s)", idx, pick.table, pick.col),
				apply: func(t *testing.T, source string) string {
					return addIndexToTable(t, source, pick.table, "KEY `"+idx+"` (`"+pick.col+"`)")
				},
			},
			delta: entSfCounts{indexes: 1}, budget: 150, touched: tableKey(pick.table),
		}, true

	case "drop_index":
		var cands []entSfIdxRec
		for _, rec := range m.idx {
			if rec.invisible || rec.fkOverlap || state.lockedTables[rec.table] || state.usedIndexes[rec.table+"."+rec.name] {
				continue
			}
			cands = append(cands, rec)
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedIndexes[pick.table+"."+pick.name] = true
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("drop_index %s.%s", pick.table, pick.name),
				apply: func(t *testing.T, source string) string {
					return entDropLineInTable(t, source, pick.table, "`"+pick.name+"`")
				},
			},
			delta: entSfCounts{indexes: -1}, touched: tableKey(pick.table),
		}, true

	case "add_fk":
		var srcs []*entFzTable
		for _, tbl := range entSfMutableTables(m, state) {
			if !tbl.partitioned {
				srcs = append(srcs, tbl)
			}
		}
		if len(srcs) == 0 {
			return entSfMutation{}, false
		}
		src := srcs[rng.Intn(len(srcs))]
		var refs []*entFzTable
		for _, tbl := range m.base.tables {
			if tbl.partitioned || state.lockedTables[tbl.name] || len(tbl.pk) != 1 {
				continue
			}
			pkCol := ""
			for c := range tbl.pk {
				pkCol = c
			}
			if strings.HasPrefix(strings.TrimPrefix(tbl.colLine[pkCol], "`"+pkCol+"` "), "bigint") &&
				strings.Contains(tbl.colLine[pkCol], "unsigned") {
				refs = append(refs, tbl)
			}
		}
		if len(refs) == 0 {
			return entSfMutation{}, false
		}
		ref := refs[rng.Intn(len(refs))]
		refPK := ""
		for c := range ref.pk {
			refPK = c
		}
		srcName, refName := src.name, ref.name
		state.lockedTables[srcName] = true
		state.lockedTables[refName] = true
		col, fk := name("r"), name("fk")
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("add_fk %s.%s -> %s.%s (%s)", srcName, col, refName, refPK, fk),
				apply: func(t *testing.T, source string) string {
					s := addColumnToTable(t, source, srcName, "`"+col+"` bigint unsigned DEFAULT NULL")
					s = addIndexToTable(t, s, srcName, "KEY `"+fk+"` (`"+col+"`)")
					return addIndexToTable(t, s, srcName,
						"CONSTRAINT `"+fk+"` FOREIGN KEY (`"+col+"`) REFERENCES `"+refName+"` (`"+refPK+"`) ON DELETE SET NULL")
				},
			},
			delta: entSfCounts{columns: 1, indexes: 1, fks: 1}, budget: 450, touched: tableKey(srcName),
		}, true

	case "drop_fk":
		var cands []entSfNamedRec
		for _, tbl := range m.base.tables {
			if state.lockedTables[tbl.name] {
				continue
			}
			for _, fk := range tbl.fks {
				if !state.usedFKs[tbl.name+"."+fk.name] {
					cands = append(cands, entSfNamedRec{table: tbl.name, name: fk.name})
				}
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedFKs[pick.table+"."+pick.name] = true
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("drop_fk %s.%s", pick.table, pick.name),
				apply: func(t *testing.T, source string) string {
					return entDropLineInTable(t, source, pick.table, "CONSTRAINT `"+pick.name+"`")
				},
			},
			delta: entSfCounts{fks: -1}, touched: tableKey(pick.table),
		}, true

	case "add_view":
		tables := entSfMutableTables(m, state)
		if len(tables) == 0 {
			return entSfMutation{}, false
		}
		tbl := tables[rng.Intn(len(tables))]
		pkCol := ""
		for _, c := range tbl.cols {
			if tbl.pk[c] {
				pkCol = c
				break
			}
		}
		if pkCol == "" {
			return entSfMutation{}, false
		}
		vw, table := name("v"), tbl.name
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("add_view %s over %s(%s)", vw, table, pkCol),
				apply: func(_ *testing.T, source string) string {
					return source + fmt.Sprintf("\nCREATE VIEW %s AS SELECT %s FROM %s;\n", vw, pkCol, table)
				},
			},
			delta: entSfCounts{views: 1}, budget: 600, touched: nil,
		}, true

	case "drop_view":
		var cands []string
		for _, v := range m.views {
			if !state.usedViews[v] {
				cands = append(cands, v)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedViews[pick] = true
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: "drop_view " + pick,
				apply: func(t *testing.T, source string) string {
					return dropObjectBlock(t, source, "VIEW", pick)
				},
			},
			delta: entSfCounts{views: -1}, touched: []string{"view:" + pick},
		}, true

	case "modify_view":
		// The extension is only appendable once per view over the whole seed (a second
		// append would duplicate the column), so the guard lives in the LEDGER, not the
		// per-round state.
		extension := map[string]string{
			"ent_fz_v1": "`users`.`name` AS `name`",
			"ent_fz_v2": "`hosts`.`status` AS `status`",
		}
		var cands []string
		for _, v := range m.views {
			if !state.usedViews[v] && !ledger.viewModified[v] && extension[v] != "" {
				cands = append(cands, v)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedViews[pick] = true
		ledger.viewModified[pick] = true
		ext := extension[pick]
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: "modify_view " + pick,
				apply: func(t *testing.T, source string) string {
					return addColumnToView(t, source, pick, ext)
				},
			},
			budget: 80, touched: []string{"view:" + pick},
		}, true

	case "add_trigger":
		tables := entSfMutableTables(m, state)
		if len(tables) == 0 {
			return entSfMutation{}, false
		}
		tbl := tables[rng.Intn(len(tables))]
		pkCol := ""
		for _, c := range tbl.cols {
			if tbl.pk[c] {
				pkCol = c
				break
			}
		}
		if pkCol == "" {
			return entSfMutation{}, false
		}
		trg, table := name("t"), tbl.name
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("add_trigger %s on %s", trg, table),
				apply: func(_ *testing.T, source string) string {
					return source + fmt.Sprintf("\nCREATE TRIGGER %s BEFORE UPDATE ON %s FOR EACH ROW SET NEW.%s = NEW.%s;\n",
						trg, table, pkCol, pkCol)
				},
			},
			delta: entSfCounts{triggers: 1}, budget: 500, touched: nil,
		}, true

	case "drop_trigger":
		var cands []string
		for _, trg := range m.base.triggers {
			if !state.usedTriggers[trg] {
				cands = append(cands, trg)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedTriggers[pick] = true
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: "drop_trigger " + pick,
				apply: func(t *testing.T, source string) string {
					return dropTrigger(t, source, pick)
				},
			},
			delta: entSfCounts{triggers: -1}, touched: []string{"trigger:" + pick},
		}, true

	case "add_routine":
		if rng.Intn(2) == 0 {
			fn := name("fn")
			ret := 100*r + seq
			return entSfMutation{
				entFzMutation: entFzMutation{
					desc: "add_routine function " + fn,
					apply: func(_ *testing.T, source string) string {
						return source + fmt.Sprintf("\nCREATE FUNCTION %s() RETURNS INT DETERMINISTIC RETURN %d;\n", fn, ret)
					},
				},
				delta: entSfCounts{functions: 1}, budget: 500, touched: nil,
			}, true
		}
		pr := name("pr")
		sel := 200*r + seq
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: "add_routine procedure " + pr,
				apply: func(_ *testing.T, source string) string {
					return source + fmt.Sprintf("\nCREATE PROCEDURE %s() BEGIN SELECT %d; END;\n", pr, sel)
				},
			},
			delta: entSfCounts{procedures: 1}, budget: 500, touched: nil,
		}, true

	case "drop_routine":
		type cand struct{ kind, name string }
		var cands []cand
		for _, fn := range m.functions {
			if !state.usedRoutines[fn] {
				cands = append(cands, cand{"FUNCTION", fn})
			}
		}
		for _, pr := range m.procedures {
			if !state.usedRoutines[pr] {
				cands = append(cands, cand{"PROCEDURE", pr})
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedRoutines[pick.name] = true
		delta := entSfCounts{functions: -1}
		key := "function:" + pick.name
		if pick.kind == "PROCEDURE" {
			delta = entSfCounts{procedures: -1}
			key = "procedure:" + pick.name
		}
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: "drop_routine " + strings.ToLower(pick.kind) + " " + pick.name,
				apply: func(t *testing.T, source string) string {
					return dropObjectBlock(t, source, pick.kind, pick.name)
				},
			},
			delta: delta, touched: []string{key},
		}, true

	case "partition":
		var cands []*entFzTable
		for _, tbl := range m.base.tables {
			if tbl.partitioned || state.lockedTables[tbl.name] || m.noPartition[tbl.name] {
				continue
			}
			if len(tbl.fks) > 0 || m.base.referenced[tbl.name] || len(tbl.pk) != 1 {
				continue
			}
			hasUnique := false
			for _, idx := range tbl.indexes {
				if idx.unique {
					hasUnique = true
					break
				}
			}
			if hasUnique {
				continue
			}
			cands = append(cands, tbl)
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		pkCol := ""
		for c := range pick.pk {
			pkCol = c
		}
		state.lockedTables[pick.name] = true
		table := pick.name
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("partition %s by hash(%s)", table, pkCol),
				apply: func(t *testing.T, source string) string {
					return entAppendTableClause(t, source, table, "PARTITION BY HASH (`"+pkCol+"`)\nPARTITIONS 4")
				},
			},
			budget: 300, touched: tableKey(table),
		}, true

	case "departition":
		var cands []*entFzTable
		for _, tbl := range m.base.tables {
			if tbl.partitioned && !state.lockedTables[tbl.name] {
				cands = append(cands, tbl)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.lockedTables[pick.name] = true
		table := pick.name
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: "departition " + table,
				apply: func(t *testing.T, source string) string {
					return entSetPartitionClause(t, source, table, "")
				},
			},
			touched: tableKey(table),
		}, true

	case "add_fulltext":
		// One FULLTEXT creation per table per round (InnoDB builds them one at a time),
		// and never on a partitioned table — the whole table is locked for the round.
		var cands []entSfColRec
		for _, rec := range entSfColCandidates(m, state, func(rest string) bool {
			return strings.HasPrefix(rest, "varchar") || strings.HasPrefix(rest, "text") ||
				strings.HasPrefix(rest, "mediumtext") || strings.HasPrefix(rest, "longtext")
		}) {
			if tbl := m.base.byName[rec.table]; tbl != nil && tbl.partitioned {
				continue
			}
			cands = append(cands, rec)
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.lockedTables[pick.table] = true
		idx := name("ft")
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("add_fulltext %s on %s(%s)", idx, pick.table, pick.col),
				apply: func(t *testing.T, source string) string {
					return addIndexToTable(t, source, pick.table, "FULLTEXT KEY `"+idx+"` (`"+pick.col+"`)")
				},
			},
			delta: entSfCounts{indexes: 1}, budget: 150, touched: tableKey(pick.table),
		}, true

	case "drop_fulltext":
		var cands []entSfNamedRec
		for _, rec := range m.fulltext {
			if !state.lockedTables[rec.table] && !state.usedIndexes[rec.table+"."+rec.name] {
				cands = append(cands, rec)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedIndexes[pick.table+"."+pick.name] = true
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("drop_fulltext %s.%s", pick.table, pick.name),
				apply: func(t *testing.T, source string) string {
					return entDropLineInTable(t, source, pick.table, "FULLTEXT KEY `"+pick.name+"`")
				},
			},
			delta: entSfCounts{indexes: -1}, touched: tableKey(pick.table),
		}, true

	case "add_spatial":
		// A NOT NULL POINT column with an SRID plus its SPATIAL index, in one mutation.
		// SPATIAL indexes are unsupported on partitioned tables; the table is locked so
		// no same-round mutation interleaves with the two-clause add.
		var cands []*entFzTable
		for _, tbl := range entSfMutableTables(m, state) {
			if !tbl.partitioned {
				cands = append(cands, tbl)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.lockedTables[pick.name] = true
		table, col, idx := pick.name, name("sp"), name("spi")
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("add_spatial %s.%s (SRID 4326) + %s", table, col, idx),
				apply: func(t *testing.T, source string) string {
					s := addColumnToTable(t, source, table, "`"+col+"` point NOT NULL /*!80003 SRID 4326 */")
					return addIndexToTable(t, s, table, "SPATIAL KEY `"+idx+"` (`"+col+"`)")
				},
			},
			delta: entSfCounts{columns: 1, indexes: 1}, budget: 400, touched: tableKey(table),
		}, true

	case "drop_spatial":
		// Index and column leave together (a dangling SPATIAL KEY over a dropped column
		// would be an invalid target).
		var cands []entSfIdxRec
		for _, rec := range m.spatial {
			if !state.lockedTables[rec.table] {
				cands = append(cands, rec)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.lockedTables[pick.table] = true
		col := pick.cols[0]
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("drop_spatial %s.%s (+column %s)", pick.table, pick.name, col),
				apply: func(t *testing.T, source string) string {
					s := entDropLineInTable(t, source, pick.table, "SPATIAL KEY `"+pick.name+"`")
					return entDropLineInTable(t, s, pick.table, "`"+col+"` point")
				},
			},
			delta: entSfCounts{columns: -1, indexes: -1}, touched: tableKey(pick.table),
		}, true

	case "add_check":
		cands := entSfColCandidates(m, state, entSfIsIntRest)
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.table+"."+pick.col] = true
		chk := name("ck")
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("add_check %s on %s(%s)", chk, pick.table, pick.col),
				apply: func(t *testing.T, source string) string {
					return addIndexToTable(t, source, pick.table, "CONSTRAINT `"+chk+"` CHECK ((`"+pick.col+"` >= 0))")
				},
			},
			delta: entSfCounts{checks: 1}, budget: 250, touched: tableKey(pick.table),
		}, true

	case "drop_check":
		var cands []entSfNamedRec
		for _, rec := range m.checks {
			if !state.lockedTables[rec.table] && !state.usedIndexes[rec.table+"."+rec.name] {
				cands = append(cands, rec)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedIndexes[pick.table+"."+pick.name] = true
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("drop_check %s.%s", pick.table, pick.name),
				apply: func(t *testing.T, source string) string {
					return entDropLineInTable(t, source, pick.table, "CONSTRAINT `"+pick.name+"`")
				},
			},
			delta: entSfCounts{checks: -1}, touched: tableKey(pick.table),
		}, true

	case "add_gen_stored", "add_gen_virtual":
		mode, tag := "STORED", "gs"
		if kind == "add_gen_virtual" {
			mode, tag = "VIRTUAL", "gv"
		}
		cands := entSfColCandidates(m, state, entSfIsIntRest)
		if mode == "STORED" {
			// The server refuses a STORED generated column whose base column is a
			// child-side FK member (errno 1215). VIRTUAL stays eligible.
			kept := cands[:0]
			for _, c := range cands {
				if !entSfIsFKMemberCol(m, c.table, c.col) {
					kept = append(kept, c)
				}
			}
			cands = kept
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.table+"."+pick.col] = true
		col := name(tag)
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("%s %s.%s as (%s + 7)", kind, pick.table, col, pick.col),
				apply: func(t *testing.T, source string) string {
					return addColumnToTable(t, source, pick.table,
						"`"+col+"` bigint GENERATED ALWAYS AS ((`"+pick.col+"` + 7)) "+mode)
				},
			},
			delta: entSfCounts{columns: 1}, budget: 300, touched: tableKey(pick.table),
		}, true

	case "modify_gen":
		// Bump the additive constant of a fuzz-shaped generated expression (`col` + N).
		var cands []entSfGenRec
		for _, rec := range m.genCols {
			if state.lockedTables[rec.table] || state.usedCols[rec.table+"."+rec.name] {
				continue
			}
			if entSfReGenConst.MatchString(rec.line) {
				cands = append(cands, rec)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.table+"."+pick.name] = true
		old := entSfReGenConst.FindStringSubmatch(pick.line)
		newLine := entSfReGenConst.ReplaceAllString(pick.line, fmt.Sprintf("+ %d)", entFzAtoi(old[1])+1))
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("modify_gen %s.%s (+%s -> +%d)", pick.table, pick.name, old[1], entFzAtoi(old[1])+1),
				apply: func(t *testing.T, source string) string {
					return entReplaceLineInTable(t, source, pick.table, "`"+pick.name+"`", newLine)
				},
			},
			budget: 80, touched: tableKey(pick.table),
		}, true

	case "drop_gen":
		var cands []entSfGenRec
		for _, rec := range m.genCols {
			if state.lockedTables[rec.table] || state.usedCols[rec.table+"."+rec.name] || entFzProtectedTables[rec.table] {
				continue
			}
			cands = append(cands, rec)
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.table+"."+pick.name] = true
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("drop_gen %s.%s", pick.table, pick.name),
				apply: func(t *testing.T, source string) string {
					return entDropLineInTable(t, source, pick.table, "`"+pick.name+"`")
				},
			},
			delta: entSfCounts{columns: -1}, touched: tableKey(pick.table),
		}, true

	case "toggle_col_invisible":
		// Visible -> invisible needs >= 2 visible columns left; invisible -> visible is
		// always legal. Both directions render/strip the /*!80023 INVISIBLE */ token the
		// dumper and omni loader agree on.
		type cand struct {
			rec entSfColRec
			on  bool
		}
		var cands []cand
		for _, rec := range m.invisCols {
			key := rec.table + "." + rec.col
			if !state.usedCols[key] && !state.lockedTables[rec.table] && !entFzProtectedTables[rec.table] {
				cands = append(cands, cand{rec: rec, on: false})
			}
		}
		for _, rec := range entSfColCandidates(m, state, func(rest string) bool {
			return entSfIsIntRest(rest) || strings.HasPrefix(rest, "varchar")
		}) {
			if m.visibleCols[rec.table] < 2 {
				continue
			}
			if strings.Contains(m.colDef(rec.table, rec.col), "INVISIBLE") {
				continue
			}
			cands = append(cands, cand{rec: rec, on: true})
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.rec.table+"."+pick.rec.col] = true
		line := m.colDef(pick.rec.table, pick.rec.col)
		newLine := line + entSfInvisibleColToken
		dir := "->invisible"
		if !pick.on {
			newLine = strings.TrimSuffix(line, entSfInvisibleColToken)
			dir = "->visible"
		}
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("toggle_col_invisible %s.%s %s", pick.rec.table, pick.rec.col, dir),
				apply: func(t *testing.T, source string) string {
					return entReplaceLineInTable(t, source, pick.rec.table, "`"+pick.rec.col+"`", newLine)
				},
			},
			budget: 60, touched: tableKey(pick.rec.table),
		}, true

	case "toggle_idx_invisible":
		var cands []entSfIdxRec
		for _, rec := range m.idx {
			if rec.fkOverlap || state.lockedTables[rec.table] || state.usedIndexes[rec.table+"."+rec.name] {
				continue
			}
			cands = append(cands, rec)
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedIndexes[pick.table+"."+pick.name] = true
		newLine := pick.line + entSfInvisibleIdxToken
		dir := "->invisible"
		if pick.invisible {
			newLine = strings.TrimSuffix(pick.line, entSfInvisibleIdxToken)
			dir = "->visible"
		}
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("toggle_idx_invisible %s.%s %s", pick.table, pick.name, dir),
				apply: func(t *testing.T, source string) string {
					return entReplaceLineInTable(t, source, pick.table, "`"+pick.name+"`", newLine)
				},
			},
			budget: 60, touched: tableKey(pick.table),
		}, true

	case "enum_append":
		var cands []entSfColRec
		for _, rec := range m.enumCols {
			key := rec.table + "." + rec.col
			if !state.usedCols[key] && !state.lockedTables[rec.table] && !entFzProtectedTables[rec.table] && !m.protected[key] {
				cands = append(cands, rec)
			}
		}
		if len(cands) == 0 {
			return entSfMutation{}, false
		}
		pick := cands[rng.Intn(len(cands))]
		state.usedCols[pick.table+"."+pick.col] = true
		line := m.colDef(pick.table, pick.col)
		member := fmt.Sprintf("m%d_%d", r, seq)
		newLine := entSfReEnumSpec.ReplaceAllString(line, "enum($1,'"+member+"')")
		if newLine == line {
			return entSfMutation{}, false
		}
		return entSfMutation{
			entFzMutation: entFzMutation{
				desc: fmt.Sprintf("enum_append %s.%s += '%s'", pick.table, pick.col, member),
				apply: func(t *testing.T, source string) string {
					return entReplaceLineInTable(t, source, pick.table, "`"+pick.col+"`", newLine)
				},
			},
			budget: 60, touched: tableKey(pick.table),
		}, true

	default:
		return entSfMutation{}, false
	}
}

// ----------------------------------------------------------------------------
// Round runner, trial (for ddmin), and the untouched-statement stability guard.
// ----------------------------------------------------------------------------

// entSfScratchDB creates an eagerly-droppable scratch database (minimization churns
// through many databases inside a single test).
func entSfScratchDB(ctx context.Context, t *testing.T, srv *liveServer, prefix string) (string, func()) {
	t.Helper()
	dbName := newLiveDatabase(ctx, t, srv, prefix)
	drop := func() {
		c, err := createLiveMySQLDriver(ctx, srv, "")
		if err != nil {
			return
		}
		defer c.Close()
		_, _ = c.Execute(ctx, "DROP DATABASE IF EXISTS `"+dbName+"`", omnidriver.ExecuteOptions{})
	}
	return dbName, drop
}

// entSfRound runs one full stateful round against an already-loaded database: mutate
// source textually, oracle (diff → apply → converge → idempotence), then the round-local
// guards — (a) the fresh dump reloads through the STRICT LoadSDL path, (c) object counts
// match want. Infrastructure failures fail t hard; protocol violations return an error
// so the caller can ddmin. Returns the fresh canonical dump.
func entSfRound(ctx context.Context, t *testing.T, srv *liveServer, dbName, source string, muts []entSfMutation, want entSfCounts) (string, error) {
	t.Helper()
	target := source
	for _, m := range muts {
		target = m.apply(t, target)
	}
	plan, err := mysqlDiffSDLMigration(source, target, srv.version)
	if err != nil {
		return "", fmt.Errorf("diff failed: %w", err)
	}
	if strings.TrimSpace(plan) == "" {
		return "", errors.New("empty plan for a non-empty mutation set")
	}
	if err := applyDDL(ctx, t, srv, dbName, plan); err != nil {
		return "", fmt.Errorf("plan failed to apply; plan was:\n%s: %w", plan, err)
	}
	after := dumpSDL(ctx, t, srv, dbName)
	converge, err := mysqlDiffSDLMigration(after, target, srv.version)
	if err != nil {
		return "", fmt.Errorf("converge diff failed: %w", err)
	}
	if converge != "" {
		return "", fmt.Errorf("did not converge; residual:\n%s\nplan was:\n%s", converge, plan)
	}
	self, err := mysqlDiffSDLMigration(after, after, srv.version)
	if err != nil {
		return "", fmt.Errorf("idempotence diff failed: %w", err)
	}
	if self != "" {
		return "", fmt.Errorf("post-apply dump not idempotent:\n%s", self)
	}
	// (a) LoadSDL must accept its own output — checked directly (DiffSDLMigration would
	// mask an SDL rejection behind the LoadSQL fallback).
	if _, err := catalog.LoadSDL("CREATE DATABASE IF NOT EXISTS `" + catalog.DefaultSDLDatabase + "`;\nUSE `" + catalog.DefaultSDLDatabase + "`;\n" + after); err != nil {
		return "", fmt.Errorf("canonical dump does not reload through LoadSDL: %w", err)
	}
	// (c) object counts must match the mutation ledger.
	got := entSfCountObjects(syncMetaForDB(ctx, t, srv, dbName))
	if got != want {
		return "", fmt.Errorf("object counts diverge from the mutation ledger:\n want %+v\n got  %+v\nplan was:\n%s", want, got, plan)
	}
	return after, nil
}

// entSfTrial reproduces one round from scratch: load baseSDL (a canonical dump) into a
// fresh database and run entSfRound over it. Used by ddmin.
func entSfTrial(ctx context.Context, t *testing.T, srv *liveServer, baseSDL string, muts []entSfMutation) error {
	t.Helper()
	dbName, drop := entSfScratchDB(ctx, t, srv, "entsdl_sf")
	defer drop()

	driver, err := createLiveMySQLDriver(ctx, srv, dbName)
	require.NoError(t, err)
	// FOREIGN_KEY_CHECKS off for the session: baseSDL here is a CANONICAL dump
	// (alphabetical table order), so inline foreign keys may forward-reference.
	_, err = driver.Execute(ctx, "SET FOREIGN_KEY_CHECKS=0;\n"+entNormalizeDelimiters(baseSDL), omnidriver.ExecuteOptions{})
	driver.Close()
	require.NoError(t, err, "[%s] stateful-fuzz trial base failed to load", srv.name)

	want := entSfCountObjects(syncMetaForDB(ctx, t, srv, dbName))
	for _, m := range muts {
		want = want.plus(m.delta)
	}
	source := dumpSDL(ctx, t, srv, dbName)
	_, err = entSfRound(ctx, t, srv, dbName, source, muts, want)
	return err
}

// entSfMinimize greedily shrinks a failing round's mutation set (ddmin-lite, mirroring
// entFuzzMinimize but over the stateful trial).
func entSfMinimize(ctx context.Context, t *testing.T, srv *liveServer, baseSDL string, muts []entSfMutation, firstErr error) ([]entSfMutation, error) {
	t.Helper()
	minimized := muts
	lastErr := firstErr
	for i := 0; i < len(minimized) && len(minimized) > 1; {
		candidate := make([]entSfMutation, 0, len(minimized)-1)
		candidate = append(candidate, minimized[:i]...)
		candidate = append(candidate, minimized[i+1:]...)
		if err := entSfTrial(ctx, t, srv, baseSDL, candidate); err != nil {
			minimized = candidate
			lastErr = err
			continue
		}
		i++
	}
	return minimized, lastErr
}

func entSfDescs(muts []entSfMutation) []string {
	out := make([]string, len(muts))
	for i, m := range muts {
		out[i] = m.desc
	}
	return out
}

// entSfStmtKeys splits a canonical dump and keys every statement by kind:name so the
// stability guard can pair statements across rounds.
func entSfStmtKeys(t *testing.T, dump string) map[string]string {
	t.Helper()
	stmts := splitSQLSegments(dump)
	out := map[string]string{}
	for i, s := range stmts {
		text := strings.TrimSpace(s.Text)
		if text == "" || text == ";" {
			continue
		}
		key := ""
		switch {
		case entReTable.MatchString(text):
			key = "table:" + entReTable.FindStringSubmatch(text)[1]
		case entSfReViewHdr.MatchString(text):
			key = "view:" + entSfReViewHdr.FindStringSubmatch(text)[1]
		case entSfReFnHdr.MatchString(text):
			key = "function:" + entSfReFnHdr.FindStringSubmatch(text)[1]
		case entSfRePrHdr.MatchString(text):
			key = "procedure:" + entSfRePrHdr.FindStringSubmatch(text)[1]
		case entFzReTrigger.MatchString(text):
			key = "trigger:" + entFzReTrigger.FindStringSubmatch(text)[1]
		default:
			key = fmt.Sprintf("other:%d", i)
		}
		out[key] = text
	}
	return out
}

// entSfAssertUntouchedStable proves every statement NOT named by this round's mutations
// is byte-identical across the round — the sharp form of the accretion guard (the
// historical paren-accretion bug regenerated untouched objects with one extra wrapper
// per cycle while converging and self-diffing empty).
func entSfAssertUntouchedStable(t *testing.T, srv *liveServer, seed int64, r int, prev, after string, muts []entSfMutation) {
	t.Helper()
	touched := map[string]bool{}
	for _, m := range muts {
		for _, k := range m.touched {
			touched[k] = true
		}
	}
	prevKeys := entSfStmtKeys(t, prev)
	afterKeys := entSfStmtKeys(t, after)
	keys := make([]string, 0, len(prevKeys))
	for k := range prevKeys {
		if _, ok := afterKeys[k]; ok && !touched[k] {
			keys = append(keys, k)
		}
	}
	slices.Sort(keys)
	for _, k := range keys {
		require.Equal(t, prevKeys[k], afterKeys[k],
			"[STATEFUZZ %s seed=%d round=%d] untouched object %s changed text across the round (collateral churn / accretion)",
			srv.name, seed, r, k)
	}
}

// ----------------------------------------------------------------------------
// The stateful fuzz test: 12 seeds, 5 rounds each.
// ----------------------------------------------------------------------------

func entSfRunSeed(ctx context.Context, t *testing.T, srv *liveServer, baseDDL string, seed int64) {
	t.Helper()
	rng := rand.New(rand.NewSource(seed))

	dbName, drop := entSfScratchDB(ctx, t, srv, "entsdl_sf")
	defer drop()
	driver, err := createLiveMySQLDriver(ctx, srv, dbName)
	require.NoError(t, err)
	_, err = driver.Execute(ctx, entNormalizeDelimiters(baseDDL), omnidriver.ExecuteOptions{})
	driver.Close()
	require.NoError(t, err, "[%s] stateful-fuzz base failed to load", srv.name)

	dump0 := dumpSDL(ctx, t, srv, dbName)
	ledger := &entSfLedger{
		expect:       entSfCountObjects(syncMetaForDB(ctx, t, srv, dbName)),
		viewModified: map[string]bool{},
	}
	prev := dump0
	budget := 0
	var history []string

	for r := 1; r <= entSfRounds; r++ {
		m := entSfParse(t, prev)
		k := 2 + rng.Intn(4) // K ∈ [2,5]
		muts := entSfMutations(t, m, rng, ledger, r, k)
		for _, mut := range muts {
			history = append(history, fmt.Sprintf("r%d: %s", r, mut.desc))
		}
		t.Logf("[STATEFUZZ %s seed=%d round=%d] K=%d mutations:\n  %s",
			srv.name, seed, r, k, strings.Join(entSfDescs(muts), "\n  "))

		want := ledger.expect
		for _, mut := range muts {
			want = want.plus(mut.delta)
		}
		after, err := entSfRound(ctx, t, srv, dbName, prev, muts, want)
		if err != nil {
			minimized, minErr := entSfMinimize(ctx, t, srv, prev, muts, err)
			t.Fatalf("[STATEFUZZ %s seed=%d round=%d] FAILED\nminimized mutations (%d of %d):\n  %s\nerror:\n%v\nseed history:\n  %s",
				srv.name, seed, r, len(minimized), len(muts), strings.Join(entSfDescs(minimized), "\n  "),
				minErr, strings.Join(history, "\n  "))
		}
		ledger.expect = want

		// (b) accretion guards: the linear text budget and the untouched-statement
		// stability check.
		for _, mut := range muts {
			budget += mut.budget
		}
		require.Less(t, len(after), len(dump0)+budget+entSfRoundSlack*r,
			"[STATEFUZZ %s seed=%d round=%d] dump grew past the added-object budget: len(dump_%d)=%d, len(dump_0)=%d, budget=%d — text accretion suspected.\nseed history:\n  %s",
			srv.name, seed, r, r, len(after), len(dump0), budget+entSfRoundSlack*r, strings.Join(history, "\n  "))
		entSfAssertUntouchedStable(t, srv, seed, r, prev, after, muts)

		t.Logf("[STATEFUZZ %s seed=%d round=%d] ok: dump %d -> %d bytes (budget %d), counts %+v",
			srv.name, seed, r, len(prev), len(after), len(dump0)+budget+entSfRoundSlack*r, ledger.expect)
		prev = after
	}
}

func TestSDLEnterpriseStatefulFuzz(t *testing.T) {
	ctx := context.Background()

	for _, srv := range tapeOracleServers(t) {
		srv := srv
		const seedCount = 12
		t.Run(srv.name, func(t *testing.T) {
			baseDDL := entSfBaseDDL(t)
			for seed := int64(1); seed <= int64(seedCount); seed++ {
				seed := seed
				t.Run(fmt.Sprintf("seed_%d", seed), func(t *testing.T) {
					entSfRunSeed(ctx, t, srv, baseDDL, seed)
				})
			}
		})
	}
}
