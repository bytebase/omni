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
