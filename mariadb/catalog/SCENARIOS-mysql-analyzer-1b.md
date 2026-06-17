# MySQL Analyzer Phase 1b — JOIN + FROM Subquery Scenarios

> **Scope:** multi-table JOINs, USING/NATURAL coalescing, FROM subqueries.
> **Prerequisite:** Phase 1a complete (30/30 scenarios passing).
> **Exit criterion:** all scenarios produce correct golden IR; oracle harness zero divergence.

## Shared catalog state

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

CREATE TABLE projects (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    department_id INT NOT NULL,
    lead_id INT
);
```

---

## Batch 7 — Basic JOINs (scenarios 7.1–7.5)

### 7.1 INNER JOIN with ON

```sql
SELECT e.name, d.name AS dept_name
FROM employees e
INNER JOIN departments d ON e.department_id = d.id;
```

**Asserts:**
- RangeTable has 3 entries: [0] RTERelation employees (ERef="e"), [1] RTERelation departments (ERef="d"), [2] RTEJoin
- JoinTree.FromList has 1 JoinExprNodeQ
- JoinExprNodeQ: JoinType=JoinInner, Quals=OpExprQ{Op:"=", Left:VarExprQ{RangeIdx:0,AttNum:4}, Right:VarExprQ{RangeIdx:1,AttNum:1}}
- TargetList[0]: VarExprQ{RangeIdx:0, AttNum:2}, TargetList[1]: VarExprQ{RangeIdx:1, AttNum:2}

### 7.2 LEFT JOIN

```sql
SELECT e.name, d.name
FROM employees e
LEFT JOIN departments d ON e.department_id = d.id;
```

**Asserts:**
- JoinExprNodeQ.JoinType = JoinLeft

### 7.3 RIGHT JOIN

```sql
SELECT e.name, d.name
FROM employees e
RIGHT JOIN departments d ON e.department_id = d.id;
```

**Asserts:**
- JoinExprNodeQ.JoinType = JoinRight

### 7.4 CROSS JOIN

```sql
SELECT e.name, d.name FROM employees e CROSS JOIN departments d;
```

**Asserts:**
- JoinExprNodeQ: JoinType=JoinCross, Quals=nil

### 7.5 Implicit cross join (comma syntax)

```sql
SELECT e.name, d.name FROM employees e, departments d WHERE e.department_id = d.id;
```

**Asserts:**
- JoinTree.FromList has **2** RangeTableRefQ (not a JoinExprNodeQ — comma is not explicit JOIN)
- RangeTable has 2 entries (no RTEJoin)
- JoinTree.Quals contains the WHERE condition

---

## Batch 8 — USING / NATURAL JOIN (scenarios 8.1–8.4)

### 8.1 JOIN USING

```sql
SELECT e.name, p.name AS project_name
FROM employees e
JOIN projects p USING (department_id);
```

**Asserts:**
- RTEJoin created with JoinUsing=["department_id"]
- JoinExprNodeQ.UsingClause=["department_id"]
- Column resolution: `department_id` from either table resolves (no ambiguity via USING coalescing)

### 8.2 NATURAL JOIN

```sql
SELECT * FROM employees e NATURAL JOIN departments d;
```

**Asserts:**
- JoinExprNodeQ.Natural=true
- RTEJoin has coalesced columns: shared columns (id, name) appear once, then remaining from each side
- Star expansion produces coalesced order

### 8.3 NATURAL LEFT JOIN

```sql
SELECT * FROM employees e NATURAL LEFT JOIN departments d;
```

**Asserts:**
- JoinExprNodeQ: Natural=true, JoinType=JoinLeft
- Same coalescing behavior as 8.2

### 8.4 JOIN USING with star expansion

```sql
SELECT * FROM employees e JOIN departments d USING (id, name);
```

**Asserts:**
- USING columns (id, name) appear once in star expansion
- Remaining columns from each table follow

---

## Batch 9 — FROM subqueries (scenarios 9.1–9.3)

### 9.1 Simple FROM subquery

```sql
SELECT sub.total FROM (SELECT COUNT(*) AS total FROM employees) AS sub;
```

**Asserts:**
- RangeTable[0]: RTESubquery, ERef="sub"
- RTE.Subquery is a *Query (the analyzed inner SELECT)
- RTE.ColNames=["total"]
- TargetList[0]: VarExprQ{RangeIdx:0, AttNum:1}

### 9.2 FROM subquery with multiple columns

```sql
SELECT x.dept, x.cnt
FROM (SELECT department_id AS dept, COUNT(*) AS cnt FROM employees GROUP BY department_id) AS x;
```

**Asserts:**
- RTESubquery with ColNames=["dept", "cnt"]
- Inner Query has GROUP BY, HasAggs=true

### 9.3 FROM subquery joined with table

```sql
SELECT e.name, sub.avg_sal
FROM employees e
JOIN (SELECT department_id, AVG(salary) AS avg_sal FROM employees GROUP BY department_id) AS sub
  ON e.department_id = sub.department_id;
```

**Asserts:**
- RangeTable: [0] RTERelation employees, [1] RTESubquery sub, [2] RTEJoin
- Join ON condition references both RTEs

---

## Batch 10 — Multi-table edge cases (scenarios 10.1–10.3)

### 10.1 Three-way JOIN

```sql
SELECT e.name, d.name AS dept, p.name AS project
FROM employees e
JOIN departments d ON e.department_id = d.id
JOIN projects p ON p.department_id = d.id;
```

**Asserts:**
- RangeTable has 5 entries: 3 RTERelation + 2 RTEJoin (one per JOIN)
- JoinTree has nested JoinExprNodeQ structure

### 10.2 Star expansion across multiple tables

```sql
SELECT * FROM employees e JOIN departments d ON e.department_id = d.id;
```

**Asserts:**
- Star expansion: 9 columns from employees + 4 from departments = 13 total
- Each VarExprQ points to the correct RTE

### 10.3 Column resolution ambiguity

```sql
SELECT name FROM employees e JOIN departments d ON e.department_id = d.id;
```

**Asserts:**
- Should return error: column 'name' is ambiguous (both tables have 'name')

---

## Progress tracker

| Batch | Scenarios | Status | Notes |
|-------|-----------|--------|-------|
| 7 — Basic JOINs | 7.1–7.5 | pending | |
| 8 — USING/NATURAL | 8.1–8.4 | pending | |
| 9 — FROM subqueries | 9.1–9.3 | pending | |
| 10 — Multi-table edges | 10.1–10.3 | pending | |
