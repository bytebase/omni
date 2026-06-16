package parser

import "testing"

func TestRoutineBodyRequiresStatementSemicolons(t *testing.T) {
	valid := []string{
		`CREATE PROCEDURE p() BEGIN END`,
		`CREATE PROCEDURE p() BEGIN SELECT 1; END`,
		`CREATE PROCEDURE p() BEGIN IF 1 THEN SELECT 1; END IF; END`,
		`CREATE PROCEDURE p() BEGIN BEGIN SELECT 1; END; END`,
	}
	for _, sql := range valid {
		t.Run("valid/"+sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	invalid := []string{
		`CREATE PROCEDURE p() BEGIN SELECT 1 END`,
		`CREATE PROCEDURE p() BEGIN SELECT 1 SELECT 2; END`,
		`CREATE PROCEDURE p() BEGIN IF 1 THEN SELECT 1 END IF; END`,
		`CREATE PROCEDURE p() BEGIN IF 1 THEN SELECT 1; END IF END`,
		`CREATE PROCEDURE p() BEGIN BEGIN SELECT 1; END END`,
	}
	for _, sql := range invalid {
		t.Run("invalid/"+sql, func(t *testing.T) {
			ParseExpectError(t, sql)
		})
	}
}

func TestTriggerAndEventBodyRequiredKeywords(t *testing.T) {
	valid := []string{
		`CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.v = 1`,
		`CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR DO SELECT 1`,
		`CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE DO SELECT 1`,
		`CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR COMMENT 'tick' DO SELECT 1`,
	}
	for _, sql := range valid {
		t.Run("valid/"+sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	invalid := []string{
		`CREATE TRIGGER tr ON t FOR EACH ROW SET @x = 1`,
		`CREATE TRIGGER tr BEFORE ON t FOR EACH ROW SET @x = 1`,
		`CREATE TRIGGER tr BEFORE INSERT ON t SET @x = 1`,
		`CREATE TRIGGER tr BEFORE INSERT ON t FOR ROW SET @x = 1`,
		`CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH SET @x = 1`,
		`CREATE EVENT ev ON EVERY 1 HOUR DO SELECT 1`,
		`CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR SELECT 1`,
		`CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT DO SELECT 1`,
		`CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR COMMENT DO SELECT 1`,
	}
	for _, sql := range invalid {
		t.Run("invalid/"+sql, func(t *testing.T) {
			ParseExpectError(t, sql)
		})
	}
}

func TestRoutineBodyANSIQuotesIdentifiers(t *testing.T) {
	valid := []string{
		`CREATE PROCEDURE "p_ansi"() BEGIN SELECT 1; END`,
		`CREATE FUNCTION "f_ansi"(a INT) RETURNS INT DETERMINISTIC RETURN a + 1`,
		`CREATE TRIGGER "tr_ansi" BEFORE INSERT ON "t_ansi" FOR EACH ROW SET NEW.v = 1`,
		`CREATE DEFINER="root"@"%:%" EVENT "ev_ansi" ON SCHEDULE EVERY 15 MINUTE STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO DELETE FROM t WHERE expired_at < NOW()`,
	}
	for _, sql := range valid {
		t.Run("valid/"+sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	invalid := []string{
		`CREATE PROCEDURE 'p_string'() BEGIN SELECT 1; END`,
		`CREATE EVENT 'ev_string' ON SCHEDULE EVERY 1 HOUR DO SELECT 1`,
	}
	for _, sql := range invalid {
		t.Run("invalid/"+sql, func(t *testing.T) {
			ParseExpectError(t, sql)
		})
	}
}
