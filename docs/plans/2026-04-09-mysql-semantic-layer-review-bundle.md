# MySQL Semantic Layer — Phase 0 Review Bundle

> **For:** user, cc, codex
> **Artifacts under review:**
> 1. `docs/plans/2026-04-09-mysql-semantic-layer.md` — full starmap (sections 1–9)
> 2. `mysql/catalog/query.go` — IR skeleton (types only, compiles clean, no behavior)
>
> **Reviewers do not need prior conversation context.** This bundle provides
> everything needed to weigh in on the open questions.

---

## TL;DR

- omni's MySQL side has L0 (parser) and L1 (DDL state machine) but is missing
  L2 (semantic analyzer) and L3 (deparse-from-IR). PG side has all of them.
- This blocks three downstream initiatives: **bytebase MySQL query span
  migration**, **MySQL SDL diff**, and miscellaneous future advisor / masking
  / information_schema work.
- Plan: build the missing L2/L3/L4 over **6 phases**, ending at **Phase 5
  (SDL diff engine)**. Phase 0 (now) is IR design only.
- **Critical constraint:** unlike PG, we cannot translate MySQL's analyzer —
  MySQL's `sql_resolver.cc` is entangled with the optimizer. We hand-write
  Go using real MySQL 8.0 containers as the behavioral oracle.
- **Out of scope:** any bytebase adapter code. omni-only effort.

---

## What you are reviewing

### Artifact 1: `docs/plans/2026-04-09-mysql-semantic-layer.md`

The full plan. Read top to bottom; it's ~700 lines but well-sectioned:

- **§1** — Layered architecture diagram (L0–L4)
- **§2** — Why MySQL's analyzer can't be translated like PG's
- **§3** — Six phases with concrete deliverables and exit criteria
- **§4** — Phase 0 IR design proposal (the thing you're reviewing)
  - §4.1 type catalog
  - **§4.2 decision log** ← 14 numbered decisions, each with rationale + open question
  - §4.3 explicit non-goals
  - §4.4 PG↔MySQL type mapping
  - §4.5 catalog struct field additions
- **§5** — Risks and mitigations
- **§6** — Out of scope
- **§7** — **10 open questions** for this review
- **Appendix B** — minimum reading order for cold reviewers

### Artifact 2: `mysql/catalog/query.go`

The IR skeleton. ~700 lines, types only, no methods beyond interface tags.
`go build ./mysql/catalog/... && go vet ./mysql/catalog/...` is clean.

Read the package comment first — it explains naming, status, and the one
known deviation from the plan doc (ResolvedType vs ColumnType, see § "Known
deviations" below).

---

## Known deviations from the plan doc

### Renamed: `ColumnType` → `ResolvedType`

**Why:** the plan doc § 4.1 names the type `ColumnType`. There is already a
field `Column.ColumnType` (string) on `mysql/catalog/table.go:63` that holds
the full type string like `"varchar(100)"`. Defining a struct named
`ColumnType` in the same package would create constant cognitive friction
for readers.

**Action requested:** approve the rename, or propose a different name.
Candidates considered: `MySQLType`, `ResolvedType`, `ExprType`, `TypeRef`,
`ColumnTypeInfo`. I chose `ResolvedType` because it captures the *role* (the
type as resolved by the analyzer at output time) rather than restating the
domain.

**Affected references:** `query.go` § 10 (the type itself), every
`AnalyzedExpr` field that previously read `*ColumnType`, and the `ColTypes
[]*ColumnType` field on `RangeTableEntry`.

---

## Reading order for cold reviewers (~30 min)

1. **`pg/catalog/query.go`** — the PG IR we are adapting. Skim the type
   definitions; you don't need to read every method. (5 min)
2. **`mysql/catalog/catalog.go`** + **`mysql/catalog/table.go`** — the
   current MySQL L1 surface. Note that `Column.ColumnType` is a string;
   that's the naming-collision context. (3 min)
3. **`mysql/catalog/viewcmds.go::createView`** — the seam where Phase 2 will
   plug the analyzer in. Currently `View.Columns` is derived via a textual
   `deparseViewSelect` helper; Phase 2 replaces it with `Query.TargetList`. (2 min)
4. **`mysql/deparse/resolver.go`** — existing scope handling (CTE virtual
   tables, USING/NATURAL coalescing). Phase 1 will refactor this into
   `mysql/catalog/scope.go`. Read the first ~150 lines. (5 min)
5. **`pg/catalog/query_span_test.go`** — the in-test lineage walker pattern
   that Phase 2 mirrors for MySQL. (3 min)
6. **`docs/plans/2026-04-09-mysql-semantic-layer.md`** §1–§4 (10 min)
7. **`mysql/catalog/query.go`** the new file (5 min)
8. **`docs/plans/2026-04-09-mysql-semantic-layer.md`** §7 — answer the 10
   open questions (your main deliverable as reviewer)

---

## Decision points needing your input

These are the key things for the review to converge on. They are restated
from plan doc § 4.2 and § 7 with the most important context inlined.

### Critical (block Phase 1 start)

1. **D1 — Schema identifier.** RangeTableEntry uses `(DBName, TableName)`
   (string pair) instead of an OID. PG uses OIDs everywhere. Pros of
   strings: no synthetic OID space to maintain, matches MySQL's actual
   identification scheme. Cons: every comparison is a string compare; some
   hot paths might want OIDs later. **Question:** stick with strings, or
   mint synthetic OIDs now for forward compatibility?

2. **D5 — MERGE view substitution.** When the analyzer encounters
   `FROM some_view`, and the view has ALGORITHM=MERGE (or UNDEFINED treated
   as MERGE), it can either:
   - **(a)** substitute the view's `AnalyzedQuery` directly into the
     RangeTableEntry's `Subquery` field at analyze time, making the RTE
     behave structurally like an RTESubquery for downstream consumers, or
   - **(b)** leave the RTE as opaque RTERelation with `IsView=true`, and
     let consumers (lineage walker, deparse) decide whether to expand.
   The plan currently proposes (a). **Question:** confirm (a), or argue for (b).

3. **D6 — Implicit casts.** MySQL inserts implicit casts pervasively
   (string → int, int → datetime, charset coercion, collation coercion).
   PG's analyzer materializes these as `RelabelType` / `CoerceViaIO` nodes.
   Doing the same for MySQL would multiply the IR node count significantly
   and require building MySQL's full type-coercion table early. The plan
   currently says **don't materialize implicit casts** until Phase 5 forces
   us to. **Question:** is anyone confident SDL diff will need them?
   CHECK constraint diffing in particular might catch us out.

4. **D11 — `query_span.go` ownership.** PG's lineage extraction lives in
   bytebase, not omni; omni only ships the IR and a test-only walker. The
   plan proposes the same for MySQL: `mysql/catalog/query_span_test.go`
   exists, but no `query_span.go` ships. **Question:** confirm this. The
   alternative is to ship a reference implementation under
   `mysql/catalog/internal/lineage/` that bytebase can vendor — but that
   creates an API surface we have to maintain.

### Important (resolve before Phase 1 starts but lower urgency)

5. **D12 — Scope refactor.** Phase 1 plans to extract the scope-handling
   code from `mysql/deparse/resolver.go` into a new `mysql/catalog/scope.go`
   so analyzer + deparse + completion all share one implementation. The
   alternative is to copy-paste and reconcile later. **Question:** appetite
   for refactor risk in Phase 1? It's about ~200 lines of code with decent
   test coverage.

6. **§ 4.3 omissions — `WithCheckOption`.** `CREATE VIEW ... WITH CHECK
   OPTION` makes the view enforce its WHERE clause on INSERT/UPDATE
   through the view. Where should this live in the IR — on `View`, on
   `Query`, or on `RangeTableEntry`? **Question:** weigh in, or defer.

7. **Phase 1 expression scope.** Phase 1 leaves `AnalyzedExpr.Type` as nil
   for everything. The Phase 2 lineage walker only needs VarExpr's RangeIdx
   and AttNum, so this should be fine — but **question:** can anyone think
   of a Phase 2 consumer that would need types?

### Optional (can be answered in line review of the file)

8. **§ 7 Q9** — should there be a Phase 2.5 that imports bytebase's golden
   query span corpus into omni as a test fixture, before declaring Phase 2
   done?

9. **§ 7 Q10** — milestone visibility: do we want explicit checkpoints at
   Phase 2 done and Phase 4 done, or just run Phase 0 → Phase 5 as one
   continuous starmap?

10. **D8, D14, etc.** — minor decisions documented in plan § 4.2 that you
    can skim and either accept or comment on.

---

## What I need from you

For each of decisions D1, D5, D6, D11, D12, and open questions §7 Q6/Q7/Q8/Q9/Q10:

1. **Approve / Reject / Need more info**
2. **If reject:** propose an alternative + 1-sentence rationale
3. **If need more info:** ask the specific question

For the IR skeleton (`mysql/catalog/query.go`):

1. **Naming** — does `ResolvedType` work for you?
2. **Coverage** — anything in PG's IR that MySQL also needs but I omitted?
   Specifically check § 4.3 of the plan doc for the explicit non-goals list.
3. **Field organization** — any RangeTableEntry / Query fields you want
   added or removed?
4. **Doc comments** — sufficient for someone with no prior context?

Add your feedback as inline comments on the doc, or as a § 9 review notes
appendix in the main plan doc.

---

## What this bundle is NOT

- It is **not** an implementation review. Phase 1 has not been written yet.
- It is **not** a request to write code. Reviewers should respond with
  decisions, not patches.
- It is **not** a spec for bytebase's adapter. That work is out of scope
  by user direction.

Phase 1 cannot start until § 7 of the main plan doc is fully answered.
