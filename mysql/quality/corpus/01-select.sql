-- MySQL SELECT corpus — quality gate scenarios
-- Each entry: -- @name / -- @valid / SQL

-- ============================================================
-- Section 2.1: JOIN keyword enforcement
-- ============================================================

-- @name: INNER JOIN accepted
-- @valid: true
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id

-- @name: LEFT JOIN accepted
-- @valid: true
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id

-- @name: CROSS JOIN accepted
-- @valid: true
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 CROSS JOIN t2

-- @name: NATURAL JOIN accepted
-- @valid: true
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 NATURAL JOIN t2

-- @name: INNER without JOIN rejected
-- @valid: false
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 INNER t2

-- @name: LEFT without JOIN rejected
-- @valid: false
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 LEFT t2

-- @name: RIGHT without JOIN rejected
-- @valid: false
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 RIGHT t2

-- @name: CROSS without JOIN rejected
-- @valid: false
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 CROSS t2

-- @name: NATURAL without JOIN rejected
-- @valid: false
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 NATURAL t2

-- @name: LEFT OUTER without JOIN rejected
-- @valid: false
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 LEFT OUTER t2

-- @name: RIGHT OUTER without JOIN rejected
-- @valid: false
-- @source: 2.1 join-keyword-enforcement
SELECT * FROM t1 RIGHT OUTER t2

-- ============================================================
-- Section 2.2: Clause validation
-- ============================================================

-- @name: WHERE without FROM rejected
-- @valid: false
-- @source: 2.2 clause-validation
SELECT * WHERE x = 1

-- @name: GROUP BY without FROM rejected
-- @valid: false
-- @source: 2.2 clause-validation
SELECT * GROUP BY x

-- @name: HAVING without FROM rejected
-- @valid: false
-- @source: 2.2 clause-validation
SELECT * HAVING count(*) > 1

-- @name: duplicate WHERE rejected
-- @valid: false
-- @source: 2.2 clause-validation
SELECT * FROM t WHERE 1=1 WHERE 2=2

-- @name: duplicate GROUP BY rejected
-- @valid: false
-- @source: 2.2 clause-validation
SELECT * FROM t GROUP BY a GROUP BY b

-- @name: duplicate ORDER BY rejected
-- @valid: false
-- @source: 2.2 clause-validation
SELECT * FROM t ORDER BY a ORDER BY b

-- @name: FROM without WHERE accepted
-- @valid: true
-- @source: 2.2 clause-validation
SELECT * FROM t

-- @name: bare SELECT without FROM accepted
-- @valid: true
-- @source: 2.2 clause-validation
SELECT 1

-- @name: normal SELECT with WHERE accepted
-- @valid: true
-- @source: 2.2 clause-validation
SELECT * FROM t WHERE x = 1

-- ============================================================
-- Section 3.1: Parenthesis balance
-- ============================================================

-- @name: unclosed parenthesis in expression rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.1
SELECT (1 + 2 FROM t

-- @name: unclosed nested parenthesis rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.1
SELECT ((1 + 2) FROM t

-- @name: unclosed parenthesis in WHERE rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.1
SELECT * FROM t WHERE (a = 1 AND b = 2

-- @name: unclosed subquery parenthesis rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.1
SELECT * FROM (SELECT 1

-- @name: unclosed function call parenthesis rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.1
SELECT COUNT(

-- @name: balanced parentheses accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.1
SELECT (1 + 2)

-- @name: nested balanced parentheses accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.1
SELECT ((1 + 2))

-- @name: subquery with closing paren accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.1
SELECT * FROM (SELECT 1) AS sub

-- ============================================================
-- Section 3.2: Binary operator right-operand
-- ============================================================

-- @name: missing right operand for + rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.2
SELECT 1 +

-- @name: missing right operand for * rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.2
SELECT 1 *

-- @name: missing right operand for = rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.2
SELECT * FROM t WHERE a =

-- @name: missing right operand for AND rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.2
SELECT * FROM t WHERE a AND

-- @name: missing right operand for OR rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.2
SELECT * FROM t WHERE a OR

-- @name: missing operand for NOT rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.2
SELECT * FROM t WHERE NOT

-- @name: complete binary expression accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.2
SELECT 1 + 2

-- @name: complete compound condition accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.2
SELECT * FROM t WHERE a = 1 AND b = 2

-- ============================================================
-- Section 3.3: IN/BETWEEN/LIKE expression completeness
-- ============================================================

-- @name: missing parenthesized list after IN rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE id IN

-- @name: unclosed IN list rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE id IN (

-- @name: missing AND in BETWEEN rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE id BETWEEN 1

-- @name: missing range in BETWEEN rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE id BETWEEN

-- @name: missing pattern after LIKE rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE name LIKE

-- @name: missing pattern after REGEXP rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE name REGEXP

-- @name: IN with value list accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE id IN (1, 2, 3)

-- @name: BETWEEN with range accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE id BETWEEN 1 AND 10

-- @name: NOT IN accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE id NOT IN (1, 2)

-- @name: NOT BETWEEN accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE id NOT BETWEEN 1 AND 10

-- @name: LIKE with pattern accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.3
SELECT * FROM t WHERE name LIKE '%foo%'

-- ============================================================
-- Section 3.4: CASE expression completeness
-- ============================================================

-- @name: CASE missing THEN and END rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.4
SELECT CASE WHEN 1 FROM t

-- @name: CASE missing END rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.4
SELECT CASE WHEN 1 THEN 'a'

-- @name: CASE missing WHEN clause rejected
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 3.4
SELECT CASE END

-- @name: CASE WHEN THEN END accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.4
SELECT CASE WHEN 1 THEN 'a' END

-- @name: simple CASE with ELSE accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 3.4
SELECT CASE x WHEN 1 THEN 'a' ELSE 'b' END
