# MySQL Analyzer Phase 1c — Subqueries, CTEs, Set Operations

> **Scope:** correlated subqueries, CTE (WITH/WITH RECURSIVE), UNION/INTERSECT/EXCEPT, IN subquery, EXISTS, scalar subquery.
> **Prerequisite:** Phase 1b complete (45 scenarios passing).
> **Exit criterion:** 100-query golden IR corpus passes.

## Shared catalog state (same as Phase 1a/1b)

```sql
CREATE DATABASE testdb; USE testdb;
CREATE TABLE employees (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(100) NOT NULL, email VARCHAR(200), department_id INT NOT NULL, salary DECIMAL(10,2) NOT NULL, hire_date DATE NOT NULL, is_active TINYINT(1) NOT NULL DEFAULT 1, notes TEXT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);
CREATE TABLE departments (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(100) NOT NULL, budget DECIMAL(12,2), location VARCHAR(100));
CREATE TABLE projects (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(100) NOT NULL, department_id INT NOT NULL, lead_id INT);
```

---

## Batch 11 — WHERE subqueries (11.1–11.5)

### 11.1 Scalar subquery in WHERE
```sql
SELECT name FROM employees WHERE salary > (SELECT AVG(salary) FROM employees);
```

### 11.2 IN subquery
```sql
SELECT name FROM employees WHERE department_id IN (SELECT id FROM departments WHERE budget > 100000);
```

### 11.3 EXISTS subquery
```sql
SELECT name FROM employees e WHERE EXISTS (SELECT 1 FROM projects p WHERE p.lead_id = e.id);
```
Asserts: SubLinkExprQ{Kind:SubLinkExists}, correlated VarExprQ has LevelsUp=1.

### 11.4 NOT IN subquery
```sql
SELECT name FROM employees WHERE department_id NOT IN (SELECT department_id FROM projects);
```

### 11.5 Scalar subquery in SELECT list
```sql
SELECT name, (SELECT COUNT(*) FROM projects p WHERE p.lead_id = e.id) AS project_count FROM employees e;
```

---

## Batch 12 — CTEs (12.1–12.5)

### 12.1 Simple CTE
```sql
WITH dept_stats AS (SELECT department_id, COUNT(*) AS cnt FROM employees GROUP BY department_id)
SELECT d.name, ds.cnt FROM departments d JOIN dept_stats ds ON d.id = ds.department_id;
```

### 12.2 Multiple CTEs
```sql
WITH
  active AS (SELECT * FROM employees WHERE is_active = 1),
  dept_totals AS (SELECT department_id, SUM(salary) AS total FROM employees GROUP BY department_id)
SELECT a.name, dt.total FROM active a JOIN dept_totals dt ON a.department_id = dt.department_id;
```

### 12.3 CTE with explicit column names
```sql
WITH emp_summary(dept, cnt) AS (SELECT department_id, COUNT(*) FROM employees GROUP BY department_id)
SELECT * FROM emp_summary;
```
Asserts: RTECTE ColNames=["dept","cnt"] (from explicit column list, not inner query).

### 12.4 CTE referenced multiple times
```sql
WITH cte AS (SELECT department_id, COUNT(*) AS cnt FROM employees GROUP BY department_id)
SELECT a.department_id, a.cnt, b.cnt FROM cte a JOIN cte b ON a.department_id = b.department_id;
```
Asserts: Two RTECTE entries in RangeTable, both referencing CTEIndex=0.

### 12.5 Recursive CTE
```sql
WITH RECURSIVE nums AS (
  SELECT 1 AS n
  UNION ALL
  SELECT n + 1 FROM nums WHERE n < 10
)
SELECT * FROM nums;
```
Asserts: Query.IsRecursive=true, CTEList[0].Recursive=true, inner Query uses SetOp=SetOpUnion.

---

## Batch 13 — Set operations (13.1–13.5)

### 13.1 UNION
```sql
SELECT name, 'employee' AS source FROM employees UNION SELECT name, 'department' FROM departments;
```
Asserts: Query.SetOp=SetOpUnion, AllSetOp=false, LArg and RArg are both *Query.

### 13.2 UNION ALL
```sql
SELECT name FROM employees UNION ALL SELECT name FROM departments;
```
Asserts: AllSetOp=true.

### 13.3 INTERSECT
```sql
SELECT department_id FROM employees INTERSECT SELECT id FROM departments;
```

### 13.4 EXCEPT
```sql
SELECT id FROM employees EXCEPT SELECT lead_id FROM projects;
```

### 13.5 Nested set operations (UNION + ORDER BY + LIMIT)
```sql
(SELECT name FROM employees) UNION ALL (SELECT name FROM departments) ORDER BY name LIMIT 10;
```
Asserts: Top-level Query has SetOp=SetOpUnion, SortClause populated, LimitCount populated.

---

## Progress tracker

| Batch | Scenarios | Status |
|-------|-----------|--------|
| 11 — WHERE subqueries | 11.1–11.5 | pending |
| 12 — CTEs | 12.1–12.5 | pending |
| 13 — Set operations | 13.1–13.5 | pending |
