# PG Expression Diversity Oracle Tests

> Goal: Verify that every PG expression type that can appear in DDL contexts is correctly handled by the analyze+DeparseExpr pipeline, validated against real PostgreSQL.
> Verification: assertOracleRoundtrip — execute migration DDL on real PG, compare schemas.
> Reference: PG ruleutils.c expression types, PG DDL expression contexts

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Expression Types in Column DEFAULT

### 1.1 Value Expressions in DEFAULT

- [x] CURRENT_TIMESTAMP — SQLValueFunction
- [x] CURRENT_DATE — SQLValueFunction
- [x] CURRENT_USER — SQLValueFunction
- [x] LOCALTIME — SQLValueFunction
- [x] LOCALTIMESTAMP — SQLValueFunction
- [x] SESSION_USER — SQLValueFunction

### 1.2 Function-Like Expressions in DEFAULT

- [x] COALESCE(a, b) — CoalesceExpr
- [x] NULLIF(a, b) — NullIfExpr
- [x] GREATEST(a, b) — MinMaxExpr
- [x] LEAST(a, b) — MinMaxExpr
- [x] CASE WHEN cond THEN val ELSE val END — CaseExpr
- [x] Regular function call: now() — FuncExpr (baseline)

### 1.3 Type and Array Expressions in DEFAULT

- [x] Type cast: 0::numeric — RelabelType / CoerceViaIO
- [x] Array constructor: ARRAY[1,2,3] — ArrayExpr
- [x] Array with type cast: ARRAY['a','b']::text[] — ArrayCoerceExpr
- [x] String concatenation: 'prefix' || 'suffix' — OpExpr
- [x] Arithmetic: 1 + 1 — OpExpr
- [x] Complex: COALESCE(CURRENT_USER, 'anon') || '-' || CAST(CURRENT_DATE AS text)

---

## Phase 2: Expression Types in Constraints and Policies

### 2.1 CHECK Constraint Expressions

- [ ] Boolean test: col IS NOT TRUE — BooleanTest
- [ ] Boolean test: col IS NOT FALSE — BooleanTest
- [ ] BETWEEN: val BETWEEN 0 AND 100 — A_Expr/OpExpr
- [ ] IN list: status IN ('a','b','c') — ScalarArrayOpExpr
- [ ] IS NULL / IS NOT NULL — NullTest
- [ ] LIKE pattern: name LIKE 'A%' — OpExpr
- [ ] CASE WHEN in CHECK — CaseExpr
- [ ] COALESCE in CHECK: COALESCE(val, 0) >= 0 — CoalesceExpr
- [ ] Subquery: EXISTS (SELECT 1 FROM other WHERE ...) — SubLink
- [ ] AND / OR / NOT combination — BoolExpr
- [ ] Array subscript: tags[1] IS NOT NULL — SubscriptingRef

### 2.2 Policy USING and WITH CHECK Expressions

- [ ] CURRENT_USER comparison: owner = CURRENT_USER — SQLValueFunction + OpExpr
- [ ] Function call in USING: has_access(id) — FuncExpr
- [ ] AND/OR in USING: active AND role = 'admin' — BoolExpr
- [ ] CASE WHEN in WITH CHECK — CaseExpr
- [ ] COALESCE in USING: COALESCE(owner, '') = CURRENT_USER

### 2.3 Domain CHECK Expressions

- [ ] VALUE > 0 — OpExpr with CoerceToDomainValue
- [ ] VALUE IS NOT NULL — NullTest with CoerceToDomainValue
- [ ] VALUE BETWEEN range — OpExpr with CoerceToDomainValue
- [ ] Function call: validate_func(VALUE) — FuncExpr
- [ ] COALESCE(VALUE, 0) >= 0 — CoalesceExpr with CoerceToDomainValue

---

## Phase 3: Expression Types in Other DDL Contexts

### 3.1 Generated Column Expressions

- [ ] Simple arithmetic: col1 + col2 — OpExpr with Var
- [ ] Function call: lower(name) — FuncExpr with Var
- [ ] COALESCE: COALESCE(a, b) — CoalesceExpr with Var
- [ ] CASE WHEN: CASE WHEN flag THEN 1 ELSE 0 END — CaseExpr with Var
- [ ] String concatenation: first || ' ' || last — OpExpr with Var

### 3.2 Index Expressions and WHERE Clauses

- [ ] Expression index: lower(name) — FuncExpr
- [ ] Expression index: (a + b) — OpExpr
- [ ] Partial index WHERE: active = true — OpExpr
- [ ] Partial index WHERE with function: length(name) > 0 — FuncExpr
- [ ] Partial index WHERE with IS NOT NULL — NullTest
- [ ] Partial index WHERE with COALESCE — CoalesceExpr

### 3.3 Trigger WHEN Clause

- [ ] Simple comparison: OLD.status IS DISTINCT FROM NEW.status — DistinctExpr
- [ ] AND combination: OLD.a <> NEW.a AND OLD.b <> NEW.b — BoolExpr
- [ ] IS NOT NULL: NEW.email IS NOT NULL — NullTest
- [ ] Function call: should_audit(NEW.id) — FuncExpr

### 3.4 Semantic Equivalence (False Positive Prevention)

- [ ] DEFAULT now() vs DEFAULT now() — same expression produces no diff
- [ ] DEFAULT 0 on int column — no spurious type cast diff
- [ ] CHECK (x > 0) vs CHECK (x > 0) — same constraint produces no diff
