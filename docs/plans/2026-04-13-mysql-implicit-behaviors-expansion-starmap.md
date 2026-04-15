# MySQL Implicit Behaviors — Coverage Expansion Starmap

> **Status:** Wave 1 dispatching (2026-04-13)
> **Driver:** this document + per-wave task tracking
> **Workers:** general-purpose sub-agents, one per category
> **Final target:** ~200 scenarios (4x expansion from the current 49)

## 0. Why this exists

The initial `SCENARIOS-mysql-implicit-behavior.md` (49 scenarios) was generated
as a "representative sample" drawn from the Round 1+2 source-code scan, with
heavy filtering for what was already spot-checked. It's a good seed, but it
**under-covers by design**: e.g. C2 type normalization has only 2 scenarios
when MySQL's type system has ~20 normalizations in scope, and entire sections
(C17 string function charset, C22 ALTER algorithm, C23 string NULL context)
were dropped during scenario filtering.

This starmap's job is to transform that seed into **systematic coverage**.
The new principle: **proactively discover gaps before bytebase hits them**,
by walking MySQL docs and source code section-by-section rather than relying
on what we happened to grep in the initial sweep.

## 1. Core principle: 3-channel discovery

Every new scenario must be anchored to **at least one** (ideally two) of:

1. **MySQL 8.0 Reference Manual** — section URL with the specific sentence(s)
   establishing the implicit behavior (`dev.mysql.com/doc/refman/8.0/...`)
2. **MySQL 8.0 source code** — `file:line` in `/Users/rebeliceyang/Github/mysql-server/`
   (branch `8.0`, version 8.0.45)
3. **omni implementation gap** — direct inspection of `mysql/catalog/*.go`
   revealing that omni currently handles the case wrong, partially, or not
   at all

Scenarios that only come from "I think this might exist" without one of the
three anchors are rejected. The anchor requirement is what makes the expansion
**systematic** rather than speculative.

## 2. Coverage targets

Current: **49 scenarios** across 17 sections.
Target: **~200 scenarios** across 25+ sections.
Expansion ratio: **~4x overall**, heavier in under-covered categories.

| Section | Topic | Current | Target | Expansion | Primary doc | Primary source |
|---|---|---:|---:|---:|---|---|
| **C1** | Name auto-generation | 6 | 12 | 2x | [create-table.html#foreign-key-naming](https://dev.mysql.com/doc/refman/8.0/en/create-table-foreign-keys.html) | `sql/sql_table.cc` |
| **C2** | Type normalization | 2 | **20** | **10x** | [data-types.html](https://dev.mysql.com/doc/refman/8.0/en/data-types.html) | `sql/field.cc`, `sql/sql_yacc.yy` |
| **C3** | Nullability promotion | 3 | 8 | 2.5x | [create-table.html](https://dev.mysql.com/doc/refman/8.0/en/create-table.html) | `sql/sql_table.cc` |
| **C4** | Charset/collation inheritance | 2 | **12** | **6x** | [charset-conversion.html](https://dev.mysql.com/doc/refman/8.0/en/charset-conversion.html) | `sql/sql_table.cc`, `sql/item.cc` |
| **C5** | Constraint defaults | 3 | 10 | 3x | [create-table-foreign-keys.html](https://dev.mysql.com/doc/refman/8.0/en/create-table-foreign-keys.html) | `sql/sql_table.cc` |
| **C6** | Partition defaults | 3 | **15** | **5x** | [partitioning.html](https://dev.mysql.com/doc/refman/8.0/en/partitioning.html) | `sql/partition_info.cc`, `sql/sql_partition.cc` |
| **C7** | Index defaults | 2 | 10 | **5x** | [create-index.html](https://dev.mysql.com/doc/refman/8.0/en/create-index.html) | `sql/sql_table.cc` |
| **C8** | Table option defaults | 3 | 10 | 3x | [create-table.html](https://dev.mysql.com/doc/refman/8.0/en/create-table.html) | `sql/handler.cc` |
| **C9** | Generated columns | 2 | 8 | 4x | [create-table-generated-columns.html](https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html) | `sql/sql_table.cc` |
| **C10** | View metadata | 2 | 8 | 4x | [view-syntax.html](https://dev.mysql.com/doc/refman/8.0/en/view-syntax.html) | `sql/sql_view.cc` |
| **C11** | Trigger defaults | 1 | 6 | 6x | [trigger-syntax.html](https://dev.mysql.com/doc/refman/8.0/en/trigger-syntax.html) | `sql/sql_trigger.cc` |
| C12 | sql_mode interactions | 0 | — | (deferred) | out of scope — omni catalog has no session state |
| C13 | (reserved) | 0 | 0 | — | — | — |
| **C14** | Constraint enforcement | 1 | 4 | 4x | [create-table-check-constraints.html](https://dev.mysql.com/doc/refman/8.0/en/create-table-check-constraints.html) | `sql/sql_table.cc` |
| **C15** | Column positioning | 1 | 5 | 5x | [alter-table.html](https://dev.mysql.com/doc/refman/8.0/en/alter-table.html) | `sql/sql_table.cc` |
| **C16** | Date/time function precision | 1 | **10** | **10x** | [date-and-time-functions.html](https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html) | `sql/item_timefunc.cc` |
| **C17** | String function charset propagation | 0 | **8** | new | [charset-unicode.html](https://dev.mysql.com/doc/refman/8.0/en/charset-unicode.html) | `sql/item_strfunc.cc` |
| **C18** | SHOW CREATE elision | 6 | **15** | 2.5x | [show-create-table.html](https://dev.mysql.com/doc/refman/8.0/en/show-create-table.html) | `sql/sql_show.cc` |
| **C19** | Virtual/functional indexes | 0 | **6** | new | [create-index.html#create-index-functional-key-parts](https://dev.mysql.com/doc/refman/8.0/en/create-index.html) | `sql/sql_table.cc` |
| **C20** | Field-type-specific defaults | 0 | **8** | new | per-type doc pages | `sql/field.cc` |
| **C21** | Parser-level defaults | 1 | **8** | 8x | n/a (source only) | `sql/sql_yacc.yy` |
| **C22** | ALTER algorithm / lock defaults | 0 | **6** | new | [alter-table.html](https://dev.mysql.com/doc/refman/8.0/en/alter-table.html) | `sql/sql_alter.cc` |
| **C23** | NULL semantics in string context | 0 | 3 | new | [string-functions.html](https://dev.mysql.com/doc/refman/8.0/en/string-functions.html) | `sql/item_strfunc.cc` |
| **C24** | Invisible PK / skip_gipk | 1 | 4 | 4x | [create-table-gipks.html](https://dev.mysql.com/doc/refman/8.0/en/create-table-gipks.html) | `sql/sql_table.cc` |
| **C25** | Numeric type defaults | 1 | 5 | 5x | [fixed-point-types.html](https://dev.mysql.com/doc/refman/8.0/en/fixed-point-types.html) | `sql/field.cc` |
| **PS** | Path-split behaviors | 8 | 10 | +2 | (source-only) | `sql/sql_table.cc` caller analysis |
| **NEW** | ALTER TABLE sub-commands section | 0 | **15** | new | [alter-table.html](https://dev.mysql.com/doc/refman/8.0/en/alter-table.html) | `sql/sql_alter.cc` |
| | **TOTAL** | **49** | **~200** | **~4x** | | |

### Priority tagging

Every new scenario gets a priority when written:

- **HIGH** — affects SDL diff correctness, deparse round-trip, or blocks a
  bytebase query span case. MUST be oracle-verified in Wave 5.
- **MED** — affects catalog state accuracy for typical schemas
- **LOW** — rare / edge case, worth documenting but not prioritizing

## 3. Worker procedure (per category)

Each worker is dispatched with:
- Category number (`C2`, `C6`, etc.)
- Current scenario count and target count
- The full text of the existing category section (verbatim copy)
- Primary doc URL + primary source file
- 3-channel discovery instructions

Worker output: **a replacement section** for that category, formatted to match
the existing SCENARIOS file, with:
- Each existing scenario preserved unchanged (unless the worker finds it wrong)
- All new scenarios appended with clear "new in expansion" marker
- Each new scenario has at least one of: doc URL + section quote, source
  file:line, or omni code pointer
- Each new scenario has `priority` tag
- Each new scenario has `status: pending` (Wave 5 changes to verified)

### Discovery checklist for each worker

1. **Walk the doc section linearly.** Open the primary doc URL. Read every
   subsection. For each paragraph, ask: "does this establish an implicit/
   automatic behavior omni would need to replicate?" Extract matching
   sentences verbatim with source URL.

2. **Grep source for category-specific patterns.** The worker gets a list of
   grep patterns tailored to the category (e.g., for C2 type normalization:
   `grep -n 'MYSQL_TYPE_\|sql_type\|real_type' sql/field.cc sql/sql_yacc.yy`).
   Each non-trivial hit is a candidate scenario.

3. **Cross-reference existing catalog.** Read the corresponding section(s) of
   `docs/plans/2026-04-13-mysql-implicit-behaviors-catalog.md` (Rounds 1/2/3
   found ~70 behaviors). Any not yet in SCENARIOS is a candidate.

4. **omni spot-check.** Grep `mysql/catalog/*.go` for related logic. If omni
   handles the case, note it. If not, mark as gap.

5. **Deduplicate and format.** Consolidate candidates into scenarios,
   avoiding overlap with existing entries. Target the scenario count but
   don't pad — if the real discovery only yields 80% of target, say so.

6. **Produce output** as a markdown section ready to be concatenated into
   SCENARIOS file.

## 4. Dispatch waves

Dispatching all 20+ workers at once would be chaos. Waves let us calibrate
and iterate.

### Wave 1 (calibration — 4 diverse categories)
Validates the worker procedure on categories with different characteristics:
- **C2 Type normalization** (2→20): doc-heavy, source patterns clear
- **C6 Partition defaults** (3→15): source-heavy, deep nesting
- **C16 Date/time precision** (1→10): mixed doc+source, function-oriented
- **C22 ALTER algorithm** (0→6): new section, minimal prior work

Review Wave 1 before Wave 2. Adjust the worker prompt if issues appear.

### Wave 2 (deep expansions on existing sections)
- **C1** name auto-generation (6→12)
- **C4** charset inheritance (2→12)
- **C7** index defaults (2→10)
- **C18** SHOW CREATE elision (6→15)

### Wave 3 (new sections)
- **C17** string function charset (0→8)
- **C19** virtual/functional indexes (0→6)
- **C20** field-type-specific defaults (0→8)
- **C21** parser-level defaults (1→8)

### Wave 4 (remaining)
- **C3** nullability (3→8)
- **C5** constraint defaults (3→10)
- **C8** table options (3→10)
- **C9** generated columns (2→8)
- **C10** view metadata (2→8)
- **C11** trigger defaults (1→6)
- **C14** enforcement (1→4)
- **C15** column positioning (1→5)
- **C23** NULL string context (0→3)
- **C24** invisible PK (1→4)
- **C25** DECIMAL (1→5)
- **NEW** ALTER TABLE sub-commands (0→15)

### Wave 5 (verification)
- Oracle-run every new HIGH-priority scenario against MySQL 8.0 container
- Update `status: verified` in progress tracker
- Patch any catalog misreadings discovered during verification
- Update `catalog_spotcheck_test.go` with new sub-tests

## 5. Completion criteria (per category)

A category is **done** when:

1. Target scenario count reached, OR the worker explicitly justifies falling
   short (e.g., "C5 had 3 pre-existing + only 4 more legitimate cases in
   MySQL docs, not 10")
2. Every new scenario has at least one of:
   - Doc URL with quoted section
   - MySQL source file:line
   - omni code pointer
3. Every HIGH-priority new scenario has been oracle-verified (Wave 5)
4. The category's doc section TOC has been walked end-to-end (not grep-only)
5. A "what I did NOT find" note at the end of the section, listing
   deliberately-excluded areas (e.g., "skipped NDB-specific partition rules")

## 6. Overall completion criteria (starmap done)

The starmap is complete when:

1. All waves dispatched and reviewed
2. SCENARIOS file reaches ~200 scenarios
3. Oracle-verified rate ≥ 80% (higher than current 71%, reflecting more rigor)
4. The catalog file and SCENARIOS file are cross-referenced — every catalog
   entry either has a corresponding scenario or an explicit "not in SCENARIOS:
   <reason>" note
5. `catalog_spotcheck_test.go` has been grown from 34 sub-tests to 100+
6. omni gaps discovered during expansion are tracked in `docs/plans/` as
   follow-up work items

## 7. Meta-goal

Every new scenario should serve at least one of:

- **A**: Unblock a bytebase query span / SDL use case (HIGH priority)
- **B**: Verify omni deparse round-trip fidelity (HIGH priority)
- **C**: Lock in a rule omni currently implements silently (MED priority)
- **D**: Expose a known-unknown — a case where we're not sure what omni does
  and want explicit assertion (MED priority)

If a candidate scenario doesn't serve any of these, it doesn't get written.

## 8. What this is NOT

- **Not a fuzzing effort.** We're not generating random DDL and seeing what
  breaks. Every scenario comes from a documented MySQL behavior.
- **Not 100% MySQL coverage.** MySQL's DDL surface is enormous. We're targeting
  "everything omni's users could reasonably hit" not "every possible MySQL
  construct". Coverage is guided by bytebase's real needs and SDL diff
  correctness, not academic completeness.
- **Not a replacement for walkthrough tests.** Walkthrough tests
  (`SCENARIOS-mysql-catalog-wt.md`) cover feature-level correctness of CREATE
  TABLE, ALTER TABLE, etc. This starmap complements them with the implicit
  behavior layer that walkthroughs don't systematically cover.

## 9. Progress tracker

Updated as waves complete. Entries: `pending` → `dispatched` → `reviewed` → `verified`.

| Wave | Category | Current | Target | Status | Notes |
|---|---|---:|---:|---|---|
| 1 | C2 Type normalization | 2 | 20 | pending | |
| 1 | C6 Partition defaults | 3 | 15 | pending | |
| 1 | C16 Date/time precision | 1 | 10 | pending | |
| 1 | C22 ALTER algorithm | 0 | 6 | pending | new section |
| 2 | C1 Name auto-gen | 6 | 12 | pending | |
| 2 | C4 Charset inheritance | 2 | 12 | pending | |
| 2 | C7 Index defaults | 2 | 10 | pending | |
| 2 | C18 SHOW CREATE elision | 6 | 15 | pending | critical for deparse |
| 3 | C17 String fn charset | 0 | 8 | pending | new section |
| 3 | C19 Functional indexes | 0 | 6 | pending | new section |
| 3 | C20 Field-type defaults | 0 | 8 | pending | new section |
| 3 | C21 Parser defaults | 1 | 8 | pending | |
| 4 | C3 Nullability | 3 | 8 | pending | |
| 4 | C5 Constraint defaults | 3 | 10 | pending | |
| 4 | C8 Table options | 3 | 10 | pending | |
| 4 | C9 Generated columns | 2 | 8 | pending | |
| 4 | C10 View metadata | 2 | 8 | pending | |
| 4 | C11 Trigger | 1 | 6 | pending | |
| 4 | C14 Enforcement | 1 | 4 | pending | |
| 4 | C15 Column positioning | 1 | 5 | pending | |
| 4 | C23 NULL in string ctx | 0 | 3 | pending | new section |
| 4 | C24 Invisible PK | 1 | 4 | pending | |
| 4 | C25 DECIMAL | 1 | 5 | pending | |
| 4 | NEW ALTER sub-commands | 0 | 15 | pending | **large new section** |
| 5 | Oracle verification of HIGH | n/a | 100%|| pending | after waves 1-4 |

## 10. Open questions

- **Scope of ALTER TABLE section**: worth 15 scenarios as a standalone section,
  or should it be folded into C15 column positioning? → decision after Wave 1
  calibration.
- **sql_mode interaction (C12)**: omni currently doesn't track sql_mode. Should
  we scope C12 in and implement sql_mode tracking, or explicitly mark it out of
  scope? → revisit after Wave 4.
- **Worker prompt length**: the worker needs to see the existing SCENARIOS
  section verbatim + doc URL + source file list + grep patterns. This could
  easily be 2000+ tokens per dispatch. Acceptable.
- **How much oracle verification during Waves 1-4 vs. Wave 5**: workers
  produce scenarios with expected values from docs/source; Wave 5 runs them
  against container. Risk: workers produce wrong expected values that Wave 5
  must catch and fix. Mitigation: require every HIGH-priority scenario to
  include a tentative oracle query that Wave 5 runs as-is.
