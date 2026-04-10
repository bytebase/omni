# MySQL Analyzer Phase 1a — Single-Table SELECT Scenarios

> **Scope:** single-table SELECT with no JOINs, no subqueries, no CTEs, no set ops.
> **Oracle:** real MySQL 8.0 container via `information_schema`.
> **Exit criterion:** all scenarios produce golden IR that matches expected output; oracle harness shows zero divergence on lineage.
>
> **Driver skill:** `mysql-analyzer-1a-driver` (TBD — first 5 scenarios are manual to establish pattern)
> **Worker skill:** `mysql-analyzer-1a-worker` (TBD)

## Shared catalog state (used by all scenarios unless overridden)

```sql
CREATE DATABASE testdb;
USE testdb;

CREATE TABLE employees (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(200),
    department_id INT NOT NULL,
    salary DECIMAL(10,2) NOT NULL,
    hire_date DATE NOT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE departments (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    budget DECIMAL(12,2),
    location VARCHAR(100)
);
```

---

## Batch 1 — Foundation (scenarios 1.1–1.5)

**Goal:** establish oracle harness, fixture format, golden IR format, and the most basic `AnalyzeSelectStmt` path.

### 1.1 Bare SELECT with literal

```sql
SELECT 1;
```

**Asserts:**
- Query.CommandType = CmdSelect
- TargetList has 1 entry: ConstExprQ{Value:"1"}, ResName="1", ResNo=1
- RangeTable is empty
- JoinTree.FromList is empty, Quals is nil

### 1.2 SELECT column FROM single table

```sql
SELECT name FROM employees;
```

**Asserts:**
- TargetList has 1 entry: VarExprQ{RangeIdx:0, AttNum:2}, ResName="name"
- RangeTable has 1 RTERelation: DBName="testdb", TableName="employees", ERef="employees"
- RTE ColNames = ["id","name","email","department_id","salary","hire_date","is_active","notes","created_at"]
- JoinTree.FromList = [RangeTableRefQ{RTIndex:0}]

### 1.3 SELECT multiple columns with alias

```sql
SELECT id, name AS employee_name, salary FROM employees;
```

**Asserts:**
- TargetList has 3 entries:
  - VarExprQ{RangeIdx:0, AttNum:1}, ResName="id"
  - VarExprQ{RangeIdx:0, AttNum:2}, ResName="employee_name"
  - VarExprQ{RangeIdx:0, AttNum:5}, ResName="salary"
- ResOrigDB="testdb", ResOrigTable="employees", ResOrigCol matches each column

### 1.4 SELECT * expansion

```sql
SELECT * FROM employees;
```

**Asserts:**
- TargetList has 9 entries (one per column in employees table, in column order)
- Each is a VarExprQ with sequential AttNum 1–9
- ResName matches ColNames of the RTE

### 1.5 SELECT t.* with table alias

```sql
SELECT e.* FROM employees AS e;
```

**Asserts:**
- RangeTable[0]: Alias="e", ERef="e", TableName="employees"
- TargetList identical to 1.4 (same expansion, but resolved through alias)

---

## Batch 2 — WHERE clause (scenarios 2.1–2.5)

**Goal:** expression analysis in WHERE, all major expression node kinds.

### 2.1 Simple WHERE equality

```sql
SELECT name FROM employees WHERE id = 1;
```

**Asserts:**
- JoinTree.Quals = OpExprQ{Op:"=", Left:VarExprQ{AttNum:1}, Right:ConstExprQ{Value:"1"}}

### 2.2 WHERE with AND/OR/NOT

```sql
SELECT name FROM employees WHERE is_active = 1 AND (department_id = 1 OR department_id = 2);
```

**Asserts:**
- JoinTree.Quals = BoolExprQ{Op:BoolAnd, Args:[
    OpExprQ{...is_active=1...},
    BoolExprQ{Op:BoolOr, Args:[...dept=1..., ...dept=2...]}
  ]}

### 2.3 WHERE with IN list

```sql
SELECT name FROM employees WHERE department_id IN (1, 2, 3);
```

**Asserts:**
- JoinTree.Quals = InListExprQ{
    Arg: VarExprQ{AttNum:4},
    List: [ConstExprQ{Value:"1"}, ConstExprQ{Value:"2"}, ConstExprQ{Value:"3"}],
    Negated: false
  }

### 2.4 WHERE with BETWEEN

```sql
SELECT name FROM employees WHERE salary BETWEEN 50000 AND 100000;
```

**Asserts:**
- JoinTree.Quals = BetweenExprQ{
    Arg: VarExprQ{AttNum:5},
    Lower: ConstExprQ{Value:"50000"},
    Upper: ConstExprQ{Value:"100000"},
    Negated: false
  }

### 2.5 WHERE with IS NULL and IS NOT NULL

```sql
SELECT name FROM employees WHERE email IS NOT NULL AND notes IS NULL;
```

**Asserts:**
- BoolExprQ{Op:BoolAnd, Args:[
    NullTestExprQ{Arg:VarExprQ{AttNum:3}, IsNull:false},
    NullTestExprQ{Arg:VarExprQ{AttNum:8}, IsNull:true}
  ]}

---

## Batch 3 — GROUP BY / HAVING / aggregates (scenarios 3.1–3.5)

**Goal:** aggregate detection, GROUP BY binding, HAVING analysis.

### 3.1 Simple GROUP BY with COUNT

```sql
SELECT department_id, COUNT(*) FROM employees GROUP BY department_id;
```

**Asserts:**
- TargetList[0] = VarExprQ{AttNum:4}, ResName="department_id"
- TargetList[1] = FuncCallExprQ{Name:"count", IsAggregate:true, Args:[]}, ResName="COUNT(*)"
- GroupClause = [SortGroupClauseQ{TargetIdx:1}]
- Query.HasAggs = true

### 3.2 GROUP BY with multiple aggregates

```sql
SELECT department_id, COUNT(*) AS cnt, SUM(salary) AS total_salary, AVG(salary) AS avg_salary
FROM employees
GROUP BY department_id;
```

**Asserts:**
- 4 TargetEntryQ entries
- FuncCallExprQ for count, sum, avg — all with IsAggregate=true
- GroupClause references TargetIdx=1 (the department_id column)

### 3.3 HAVING clause

```sql
SELECT department_id, COUNT(*) AS cnt
FROM employees
GROUP BY department_id
HAVING COUNT(*) > 5;
```

**Asserts:**
- Query.HavingQual = OpExprQ{Op:">",
    Left: FuncCallExprQ{Name:"count", IsAggregate:true},
    Right: ConstExprQ{Value:"5"}
  }

### 3.4 COUNT DISTINCT

```sql
SELECT COUNT(DISTINCT department_id) FROM employees;
```

**Asserts:**
- FuncCallExprQ{Name:"count", IsAggregate:true, Distinct:true, Args:[VarExprQ{AttNum:4}]}

### 3.5 GROUP BY ordinal

```sql
SELECT department_id, COUNT(*) FROM employees GROUP BY 1;
```

**Asserts:**
- GroupClause = [SortGroupClauseQ{TargetIdx:1}]
- Same as 3.1 — ordinal `1` resolves to first TargetEntryQ

---

## Batch 4 — ORDER BY / LIMIT / DISTINCT (scenarios 4.1–4.5)

**Goal:** sorting, limit/offset, distinct, and ResJunk for ORDER BY helpers.

### 4.1 ORDER BY column

```sql
SELECT name, salary FROM employees ORDER BY salary DESC;
```

**Asserts:**
- SortClause = [SortGroupClauseQ{TargetIdx:2, Descending:true, NullsFirst:false}]

### 4.2 ORDER BY expression not in select list (ResJunk)

```sql
SELECT name FROM employees ORDER BY salary;
```

**Asserts:**
- TargetList has 2 entries:
  - ResNo=1: VarExprQ{AttNum:2}, ResName="name", ResJunk=false
  - ResNo=2: VarExprQ{AttNum:5}, ResName="salary", ResJunk=true
- SortClause references TargetIdx=2 (the junk entry)

### 4.3 LIMIT and OFFSET

```sql
SELECT name FROM employees ORDER BY id LIMIT 10 OFFSET 20;
```

**Asserts:**
- Query.LimitCount = ConstExprQ{Value:"10"}
- Query.LimitOffset = ConstExprQ{Value:"20"}

### 4.4 MySQL-style LIMIT offset, count

```sql
SELECT name FROM employees LIMIT 20, 10;
```

**Asserts:**
- Same as 4.3: LimitCount=10, LimitOffset=20 (normalized form)

### 4.5 SELECT DISTINCT

```sql
SELECT DISTINCT department_id FROM employees;
```

**Asserts:**
- Query.Distinct = true
- TargetList has 1 non-junk entry

---

## Batch 5 — Expressions (scenarios 5.1–5.6)

**Goal:** various expression node kinds that appear in SELECT list and WHERE.

### 5.1 Arithmetic expression

```sql
SELECT name, salary * 12 AS annual_salary FROM employees;
```

**Asserts:**
- TargetList[1] = OpExprQ{Op:"*", Left:VarExprQ{AttNum:5}, Right:ConstExprQ{Value:"12"}}
- ResName="annual_salary"
- ResOrigDB/Table/Col all empty (computed column, no single-source provenance)

### 5.2 Function call (scalar, non-aggregate)

```sql
SELECT CONCAT(name, ' <', email, '>') AS display FROM employees;
```

**Asserts:**
- FuncCallExprQ{Name:"concat", IsAggregate:false, Args:[
    VarExprQ{AttNum:2}, ConstExprQ{Value:" <"}, VarExprQ{AttNum:3}, ConstExprQ{Value:">"}
  ]}

### 5.3 CASE expression (searched form)

```sql
SELECT name,
       CASE WHEN salary > 100000 THEN 'high'
            WHEN salary > 50000 THEN 'mid'
            ELSE 'low'
       END AS salary_band
FROM employees;
```

**Asserts:**
- CaseExprQ{TestExpr:nil, Args:[
    CaseWhenQ{Cond:OpExprQ{...salary>100000}, Then:ConstExprQ{Value:"high"}},
    CaseWhenQ{Cond:OpExprQ{...salary>50000}, Then:ConstExprQ{Value:"mid"}},
  ], Default:ConstExprQ{Value:"low"}}

### 5.4 CASE expression (simple form)

```sql
SELECT CASE department_id WHEN 1 THEN 'eng' WHEN 2 THEN 'sales' END FROM employees;
```

**Asserts:**
- CaseExprQ{TestExpr:VarExprQ{AttNum:4}, Args:[...], Default:nil}

### 5.5 COALESCE and IFNULL

```sql
SELECT COALESCE(email, 'no-email') AS email, IFNULL(notes, '') AS notes FROM employees;
```

**Asserts:**
- TargetList[0] = CoalesceExprQ{Args:[VarExprQ{AttNum:3}, ConstExprQ{Value:"no-email"}]}
- TargetList[1] = CoalesceExprQ{Args:[VarExprQ{AttNum:8}, ConstExprQ{Value:""}]}
  (IFNULL normalized to COALESCE)

### 5.6 CAST expression

```sql
SELECT CAST(salary AS SIGNED) AS salary_int FROM employees;
```

**Asserts:**
- CastExprQ{Arg:VarExprQ{AttNum:5}, TargetType:&ResolvedType{...}}
  (TargetType populated even in Phase 1 because it's explicit in the source)

---

## Batch 6 — Edge cases (scenarios 6.1–6.4)

**Goal:** tricky single-table scenarios that tend to trip up hand-written analyzers.

### 6.1 Column name ambiguity with alias

```sql
SELECT name AS id FROM employees WHERE id = 1;
```

**Asserts:**
- WHERE references the *column* `id` (AttNum:1), not the *alias* `id`
- MySQL resolves unqualified names in WHERE against base columns, not SELECT aliases

### 6.2 SELECT expression referencing same column twice

```sql
SELECT salary, salary + 1000 AS raised FROM employees;
```

**Asserts:**
- Two separate VarExprQ entries for `salary`, both pointing to AttNum:5
- The VarExprQ are distinct objects (not shared pointers)

### 6.3 Qualified column with database name

```sql
SELECT testdb.employees.name FROM testdb.employees;
```

**Asserts:**
- VarExprQ{RangeIdx:0, AttNum:2} (resolves through db.table.column)
- RTE: DBName="testdb", TableName="employees"

### 6.4 SELECT with no FROM clause (dual)

```sql
SELECT 1 + 2, 'hello', NOW();
```

**Asserts:**
- RangeTable is empty
- JoinTree.FromList is empty
- TargetList: OpExprQ, ConstExprQ, FuncCallExprQ
- No VarExprQ anywhere (no table context)

---

## Progress tracker

| Batch | Scenarios | Status | Notes |
|-------|-----------|--------|-------|
| 1 — Foundation | 1.1–1.5 | pending | Manual (establish pattern) |
| 2 — WHERE | 2.1–2.5 | pending | |
| 3 — GROUP/HAVING | 3.1–3.5 | pending | |
| 4 — ORDER/LIMIT/DISTINCT | 4.1–4.5 | pending | |
| 5 — Expressions | 5.1–5.6 | pending | |
| 6 — Edge cases | 6.1–6.4 | pending | |
