package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mysql/ast"
)

// TestRoutineBodyRealworld covers the realworld corpus tracked by
// docs/plans/2026-04-20-mysql-routine-body-grammar.md.
//
// Each row is either active (current scanner handles it, test must pass now)
// or t.Skip-gated with a TODO pointing at the plan; after commit 3 of that
// plan swaps the body scanner for parseCompoundStmtOrStmt, the Skips are
// removed and every row passes on all branches.
//
// Sources:
//   - existing in-repo samples (split_realworld_test.go)
//   - MySQL 8.0 Sakila sample database (simplified where schema-dependent)
//   - synthetic cases mirroring every row of the plan's behavior matrix
//   - generalized shape of the totem-dev-style `if(var) then` failures.
func TestRoutineBodyRealworld(t *testing.T) {
	type tc struct {
		name string
		sql  string
		skip string // non-empty ⇒ t.Skip reason (TODO: clear after commit 3)
	}

	cases := []tc{
		// --- grammar matrix: IF / compound vs function, stmt/expr context ---
		{
			name: "if(x) then lowercase no space",
			sql: `DELIMITER ;;
CREATE PROCEDURE p(IN x INT)
BEGIN
    if(x > 0) then
        SET @r = 'positive';
    end if;
END ;;
DELIMITER ;`,
			skip: "TODO(plan 2026-04-20 commit 3): scanner under-counts IF(, grammar swap fixes",
		},
		{
			name: "IF(x) THEN uppercase no space",
			sql: `CREATE PROCEDURE p(IN x INT)
BEGIN
    IF(x > 0) THEN
        SET @r = 1;
    END IF;
END`,
			skip: "TODO(plan 2026-04-20 commit 3): scanner under-counts IF(, grammar swap fixes",
		},
		{
			name: "IF (x) THEN with space",
			sql: `CREATE PROCEDURE p(IN x INT)
BEGIN
    IF (x > 0) THEN
        SET @r = 1;
    END IF;
END`,
			skip: "TODO(plan 2026-04-20 commit 3): scanner under-counts IF(, grammar swap fixes",
		},
		{
			name: "IF cond THEN no parens",
			sql: `CREATE PROCEDURE p(IN x INT)
BEGIN
    IF x > 0 THEN
        SET @r = 1;
    END IF;
END`,
		},
		{
			name: "IF EXISTS (subquery) THEN flow control",
			sql: `CREATE PROCEDURE p()
BEGIN
    IF EXISTS (SELECT 1 FROM t WHERE id = 1) THEN
        DELETE FROM t WHERE id = 1;
    END IF;
END`,
		},
		{
			name: "SET @r = IF(a,b,c) function call in expr",
			sql: `CREATE PROCEDURE p(IN x INT)
BEGIN
    SET @r = IF(x > 0, 'pos', 'neg');
END`,
		},
		{
			name: "SET @r = IF (a,b,c) function call with space",
			sql: `CREATE PROCEDURE p(IN x INT)
BEGIN
    SET @r = IF (x > 0, 'pos', 'neg');
END`,
		},

		// --- WHILE / compound vs ambiguity ---
		{
			name: "WHILE(cond) DO no space",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE(i < 10) DO
        SET i = i + 1;
    END WHILE;
END`,
		},
		{
			name: "WHILE cond DO no parens",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE i < 10 DO
        SET i = i + 1;
    END WHILE;
END`,
		},

		// --- CASE stmt vs expr ---
		{
			name: "CASE stmt simple form",
			sql: `CREATE PROCEDURE p(IN x INT)
BEGIN
    CASE x
        WHEN 1 THEN SET @r = 'one';
        WHEN 2 THEN SET @r = 'two';
        ELSE SET @r = 'other';
    END CASE;
END`,
		},
		{
			name: "CASE stmt searched form",
			sql: `CREATE PROCEDURE p(IN x INT)
BEGIN
    CASE
        WHEN x < 0 THEN SET @r = 'neg';
        WHEN x = 0 THEN SET @r = 'zero';
        ELSE SET @r = 'pos';
    END CASE;
END`,
		},
		{
			name: "CASE expr inside SET",
			sql: `CREATE FUNCTION f(x INT) RETURNS VARCHAR(10)
BEGIN
    DECLARE r VARCHAR(10);
    SET r = CASE x WHEN 1 THEN 'one' ELSE 'other' END;
    RETURN r;
END`,
		},

		// --- LOOP, REPEAT, labels ---
		{
			name: "labeled LOOP with LEAVE and ITERATE",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE i INT DEFAULT 0;
    outer_loop: LOOP
        SET i = i + 1;
        IF i = 3 THEN
            ITERATE outer_loop;
        END IF;
        IF i >= 5 THEN
            LEAVE outer_loop;
        END IF;
    END LOOP outer_loop;
END`,
		},
		{
			name: "REPEAT UNTIL END REPEAT",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE i INT DEFAULT 0;
    REPEAT
        SET i = i + 1;
    UNTIL i >= 10
    END REPEAT;
END`,
		},
		{
			name: "nested compound: IF inside WHILE inside IF",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE i INT DEFAULT 0;
    IF @cond = 1 THEN
        WHILE i < 10 DO
            IF i % 2 = 0 THEN
                CASE i
                    WHEN 0 THEN SET @a = 'zero';
                    ELSE SET @a = 'even';
                END CASE;
            END IF;
            SET i = i + 1;
        END WHILE;
    END IF;
END`,
		},
		{
			name: "nested IF with inner if(x)",
			sql: `CREATE PROCEDURE p(IN x INT, IN y INT)
BEGIN
    IF x > 0 THEN
        if(y > 0) then
            SET @r = 1;
        end if;
    END IF;
END`,
			skip: "TODO(plan 2026-04-20 commit 3): inner if(y) fails scanner depth; grammar swap fixes",
		},

		// --- Cursors + handlers ---
		{
			name: "cursor + continue handler + loop",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE v INT;
    DECLARE cur CURSOR FOR SELECT id FROM t;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO v;
        IF done THEN
            LEAVE read_loop;
        END IF;
        SET @last = v;
    END LOOP read_loop;
    CLOSE cur;
END`,
		},
		{
			name: "named condition + exit handler with BEGIN...END body",
			sql: `CREATE PROCEDURE p()
BEGIN
    DECLARE dup_key CONDITION FOR SQLSTATE '23000';
    DECLARE EXIT HANDLER FOR dup_key
    BEGIN
        SET @err = 'duplicate';
    END;
    INSERT INTO t(id) VALUES (1);
END`,
		},
		{
			name: "SIGNAL and RESIGNAL in handler",
			sql: `CREATE PROCEDURE p(IN x INT)
BEGIN
    IF x < 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'negative';
    END IF;
    BEGIN
        DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            RESIGNAL;
        END;
        SET @y = x * 2;
    END;
END`,
		},

		// --- Dynamic SQL ---
		{
			name: "PREPARE / EXECUTE / DEALLOCATE",
			sql: `CREATE PROCEDURE p(IN tbl VARCHAR(64))
BEGIN
    SET @sql = CONCAT('SELECT * FROM ', tbl);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END`,
		},

		// --- DDL inside body ---
		{
			name: "DROP TABLE IF EXISTS + CREATE TABLE IF NOT EXISTS",
			sql: `CREATE PROCEDURE p()
BEGIN
    DROP TABLE IF EXISTS tmp;
    CREATE TABLE IF NOT EXISTS tmp (id INT);
    INSERT INTO tmp VALUES (1);
END`,
		},

		// --- Comments and whitespace ---
		{
			name: "comments inside body",
			sql: `CREATE PROCEDURE p()
BEGIN
    /* block comment */
    -- line comment
    # hash comment
    SET @x = 1;
END`,
		},
		{
			name: "comment between BEGIN and first stmt",
			sql: `CREATE PROCEDURE p(IN x INT)
BEGIN
    /* explanatory note */
    IF x > 0 THEN
        SET @r = 1;
    END IF;
END`,
		},

		// --- DELIMITER and multi-procedure ---
		{
			name: "multiple procedures in one DELIMITER block",
			sql: `DELIMITER ;;
CREATE PROCEDURE p_a()
BEGIN
    IF @x = 1 THEN SET @y = 2; END IF;
END ;;
CREATE PROCEDURE p_b()
BEGIN
    WHILE @i < 3 DO SET @i = @i + 1; END WHILE;
END ;;
DELIMITER ;`,
		},

		// --- DEFINER and characteristics ---
		{
			name: "procedure with DEFINER and characteristics",
			sql: "CREATE DEFINER=`root`@`%` PROCEDURE p_def(IN x INT)\n" +
				"    READS SQL DATA\n" +
				"    DETERMINISTIC\n" +
				"    SQL SECURITY DEFINER\n" +
				"    COMMENT 'test'\n" +
				"BEGIN\n" +
				"    IF x > 0 THEN SET @y = 1; END IF;\n" +
				"END",
		},

		// --- Functions ---
		{
			name: "function with single RETURN body (no BEGIN)",
			sql:  `CREATE FUNCTION f_simple(a INT) RETURNS INT RETURN a + 1`,
		},
		{
			name: "function with BEGIN body and RETURN",
			sql: `CREATE FUNCTION f(x INT) RETURNS VARCHAR(20)
BEGIN
    DECLARE r VARCHAR(20);
    IF x > 0 THEN
        SET r = 'positive';
    ELSE
        SET r = 'non-positive';
    END IF;
    RETURN r;
END`,
		},

		// --- Triggers ---
		{
			name: "trigger BEFORE INSERT with IF",
			sql: `CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW
BEGIN
    IF NEW.x IS NULL THEN
        SET NEW.x = 0;
    END IF;
END`,
		},
		{
			name: "trigger with WHILE loop",
			sql: `CREATE TRIGGER trg_w BEFORE UPDATE ON t FOR EACH ROW
BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE i < 3 DO
        SET i = i + 1;
    END WHILE;
END`,
		},

		// --- Events ---
		{
			name: "event with IF EXISTS predicate",
			sql: `CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR DO
BEGIN
    IF EXISTS (SELECT 1 FROM t) THEN
        DELETE FROM t;
    END IF;
END`,
		},
		{
			name: "ALTER EVENT DO body with compound IF (plan's missed-by-#97 site)",
			sql: `ALTER EVENT ev DO
BEGIN
    IF @x > 0 THEN SET @y = 1; END IF;
END`,
			skip: "TODO(plan 2026-04-20 commit 3): ALTER EVENT scanner unfixed by PR #97; grammar swap fixes",
		},

		// --- Sakila-inspired cases (schema-simplified, grammar-preserving) ---
		{
			name: "sakila film_in_stock shape: handler + cursor + OPEN/FETCH/CLOSE",
			sql: `CREATE PROCEDURE film_in_stock (IN p_film_id INT, IN p_store_id INT, OUT p_film_count INT)
READS SQL DATA
BEGIN
    SELECT inventory_id FROM inventory
    WHERE film_id = p_film_id AND store_id = p_store_id;
    SELECT COUNT(*) INTO p_film_count FROM inventory WHERE film_id = p_film_id;
END`,
		},
		{
			name: "sakila rewards_report shape: checks + IF + cursor",
			sql: `CREATE PROCEDURE rewards_report (
    IN min_monthly_purchases TINYINT UNSIGNED,
    IN min_dollar_amount_purchased DECIMAL(10,2) UNSIGNED,
    OUT count_rewardees INT
)
READS SQL DATA
BEGIN
    DECLARE last_month_start DATE;
    DECLARE last_month_end DATE;

    IF min_monthly_purchases = 0 THEN
        SELECT 'Minimum monthly purchases parameter must be > 0';
        SET count_rewardees = 0;
    ELSE
        SET last_month_start = CURRENT_DATE - INTERVAL 1 MONTH;
        SET last_month_start = last_month_start - INTERVAL DAY(last_month_start) - 1 DAY;
        SET last_month_end = last_month_start + INTERVAL 1 MONTH - INTERVAL 1 DAY;
        SELECT COUNT(*) INTO count_rewardees FROM customer;
    END IF;
END`,
		},
		{
			name: "sakila inventory_held_by_customer function shape",
			sql: `CREATE FUNCTION inventory_held_by_customer(p_inventory_id INT) RETURNS INT
READS SQL DATA
BEGIN
    DECLARE v_customer_id INT;
    SELECT customer_id INTO v_customer_id
    FROM rental
    WHERE return_date IS NULL AND inventory_id = p_inventory_id;
    RETURN v_customer_id;
END`,
		},
		{
			name: "sakila get_customer_balance function shape with multi-stmt body",
			sql: `CREATE FUNCTION get_customer_balance(p_customer_id INT, p_effective_date DATETIME) RETURNS DECIMAL(5,2)
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE v_rentfees DECIMAL(5,2);
    DECLARE v_overfees INT;
    DECLARE v_payments DECIMAL(5,2);
    SELECT IFNULL(SUM(amount), 0) INTO v_payments
    FROM payment
    WHERE customer_id = p_customer_id AND payment_date <= p_effective_date;
    RETURN v_payments;
END`,
		},
		{
			name: "sakila customer_create_date trigger shape",
			sql: `CREATE TRIGGER customer_create_date BEFORE INSERT ON customer
FOR EACH ROW SET NEW.create_date = NOW()`,
		},
		{
			name: "sakila ins_film trigger shape with multi-stmt body",
			sql: `CREATE TRIGGER ins_film AFTER INSERT ON film
FOR EACH ROW
BEGIN
    INSERT INTO film_text (film_id, title, description)
    VALUES (new.film_id, new.title, new.description);
END`,
		},

		// --- totem-dev style generalized ---
		{
			name: "totem-dev pattern: dynamic SQL built inside if(var) then",
			sql: `DELIMITER ;;
CREATE PROCEDURE lessonList(IN quiz INT, IN course_id INT)
BEGIN
    SET @statement = CONCAT('SELECT id FROM course_lessons WHERE course_id = ', course_id, ' AND course_lesson_group_id IS NULL');
    if(quiz) then
        SET @statement = CONCAT(@statement, ' AND type <> "QUIZ"');
    end if;
    if(course_id) then
        SET @statement = CONCAT(@statement, ' AND active = 1');
    end if;
    PREPARE stmt FROM @statement;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END ;;
DELIMITER ;`,
			skip: "TODO(plan 2026-04-20 commit 3): scanner under-counts if(, grammar swap fixes",
		},

		// --- Labels (scope: B2 commit 4 enforces matching) ---
		{
			name: "labeled BEGIN with matching end-label",
			sql: `CREATE PROCEDURE p()
myblock: BEGIN
    SET @a = 1;
END myblock`,
		},
		{
			name: "labeled BEGIN with case-insensitive end-label match",
			sql: `CREATE PROCEDURE p()
myBlock: BEGIN
    SET @a = 1;
END MYBLOCK`,
			// passes under current scanner (scanner doesn't compare labels);
			// still passes after commit 4 because match is case-insensitive.
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.skip != "" {
				t.Skip(c.skip)
			}
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v\nsql:\n%s", err, c.sql)
			}
			if list == nil || len(list.Items) == 0 {
				t.Fatal("expected at least one parsed statement")
			}
		})
	}
}

// TestRoutineBodyRealworld_BodyTextRoundTrip verifies that stmt.BodyText
// (currently populated by the scanner) can itself be fed back through
// Parse() wrapped as a procedure body. After commit 3 this becomes a
// stronger assertion via parseCompoundStmtOrStmt directly, but even against
// the current scanner the round-trip of the captured body text should hold
// for cases the scanner handles.
func TestRoutineBodyRealworld_BodyTextRoundTrip(t *testing.T) {
	cases := []string{
		`CREATE PROCEDURE p()
BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE i < 3 DO
        SET i = i + 1;
    END WHILE;
END`,
		`CREATE FUNCTION f(x INT) RETURNS INT
BEGIN
    IF x > 0 THEN RETURN x; END IF;
    RETURN 0;
END`,
	}
	for _, sql := range cases {
		t.Run(firstLine(sql), func(t *testing.T) {
			list, err := Parse(sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			var bodyText string
			switch s := list.Items[0].(type) {
			case *ast.CreateFunctionStmt:
				bodyText = s.BodyText
			case *ast.CreateTriggerStmt:
				bodyText = s.BodyText
			case *ast.CreateEventStmt:
				bodyText = s.BodyText
			default:
				t.Fatalf("unexpected top-level type %T", list.Items[0])
			}
			if strings.TrimSpace(bodyText) == "" {
				t.Fatal("body text is empty")
			}
			// After commit 3, round-trip bodyText through
			// parseCompoundStmtOrStmt directly. For now just assert non-empty.
		})
	}
}

func firstLine(s string) string {
	before, _, _ := strings.Cut(s, "\n")
	return before
}
