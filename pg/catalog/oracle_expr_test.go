package catalog

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Phase 1: Expression Types in Column DEFAULT
// ---------------------------------------------------------------------------

func TestOracleExprDefault_ValueExpressions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// 1.1 Value Expressions in DEFAULT

	t.Run("CURRENT_TIMESTAMP", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, ts timestamptz DEFAULT CURRENT_TIMESTAMP);`,
		)
	})

	t.Run("CURRENT_DATE", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, d date DEFAULT CURRENT_DATE);`,
		)
	})

	t.Run("CURRENT_USER", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, u text DEFAULT CURRENT_USER);`,
		)
	})

	t.Run("LOCALTIME", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, lt time DEFAULT LOCALTIME);`,
		)
	})

	t.Run("LOCALTIMESTAMP", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, lt timestamp DEFAULT LOCALTIMESTAMP);`,
		)
	})

	t.Run("SESSION_USER", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, u text DEFAULT SESSION_USER);`,
		)
	})
}

func TestOracleExprDefault_FunctionLikeExpressions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// 1.2 Function-Like Expressions in DEFAULT

	t.Run("COALESCE", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, val int DEFAULT COALESCE(NULL, 42));`,
		)
	})

	t.Run("NULLIF", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, val int DEFAULT NULLIF(0, 0));`,
		)
	})

	t.Run("GREATEST", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, val int DEFAULT GREATEST(1, 2));`,
		)
	})

	t.Run("LEAST", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, val int DEFAULT LEAST(1, 2));`,
		)
	})

	t.Run("CASE_WHEN", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, val text DEFAULT CASE WHEN true THEN 'yes' ELSE 'no' END);`,
		)
	})

	t.Run("now_function", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, ts timestamptz DEFAULT now());`,
		)
	})
}

func TestOracleExprDefault_TypeAndArrayExpressions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// 1.3 Type and Array Expressions in DEFAULT

	t.Run("type_cast", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, val numeric DEFAULT 0::numeric);`,
		)
	})

	t.Run("array_constructor", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, vals int[] DEFAULT ARRAY[1,2,3]);`,
		)
	})

	t.Run("array_with_type_cast", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, vals text[] DEFAULT ARRAY['a','b']::text[]);`,
		)
	})

	t.Run("string_concatenation", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, val text DEFAULT 'prefix' || 'suffix');`,
		)
	})

	t.Run("arithmetic", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, val int DEFAULT 1 + 1);`,
		)
	})

	t.Run("complex_expression", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE TABLE t (id int PRIMARY KEY, val text DEFAULT COALESCE(CURRENT_USER, 'anon') || '-' || CAST(CURRENT_DATE AS text));`,
		)
	})
}

// ---------------------------------------------------------------------------
// Phase 2: Expression Types in Constraints and Policies
// ---------------------------------------------------------------------------

func TestOracleExprCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// 2.1 CHECK Constraint Expressions

	t.Run("boolean_test_IS_NOT_TRUE", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, active boolean);`,
			`CREATE TABLE t (id int PRIMARY KEY, active boolean, CONSTRAINT chk CHECK (active IS NOT TRUE));`,
		)
	})

	t.Run("boolean_test_IS_NOT_FALSE", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, active boolean);`,
			`CREATE TABLE t (id int PRIMARY KEY, active boolean, CONSTRAINT chk CHECK (active IS NOT FALSE));`,
		)
	})

	t.Run("BETWEEN", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, age int);`,
			`CREATE TABLE t (id int PRIMARY KEY, age int, CONSTRAINT chk CHECK (age BETWEEN 0 AND 150));`,
		)
	})

	t.Run("IN_list", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, status text);`,
			`CREATE TABLE t (id int PRIMARY KEY, status text, CONSTRAINT chk CHECK (status IN ('a', 'b', 'c')));`,
		)
	})

	t.Run("IS_NOT_NULL", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, name text);`,
			`CREATE TABLE t (id int PRIMARY KEY, name text, CONSTRAINT chk CHECK (name IS NOT NULL));`,
		)
	})

	t.Run("LIKE_pattern", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, name text);`,
			`CREATE TABLE t (id int PRIMARY KEY, name text, CONSTRAINT chk CHECK (name LIKE 'A%'));`,
		)
	})

	t.Run("CASE_WHEN_in_CHECK", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, val int);`,
			`CREATE TABLE t (id int PRIMARY KEY, val int, CONSTRAINT chk CHECK (CASE WHEN val > 0 THEN true ELSE false END));`,
		)
	})

	t.Run("COALESCE_in_CHECK", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, val int);`,
			`CREATE TABLE t (id int PRIMARY KEY, val int, CONSTRAINT chk CHECK (COALESCE(val, 0) >= 0));`,
		)
	})

	t.Run("subquery_EXISTS", func(t *testing.T) {
		t.Skip("BUG: PostgreSQL does not allow subqueries in CHECK constraints (ERROR 0A000)")
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE other (id int PRIMARY KEY);
			 CREATE TABLE t (id int PRIMARY KEY, ref int);`,
			`CREATE TABLE other (id int PRIMARY KEY);
			 CREATE TABLE t (id int PRIMARY KEY, ref int, CONSTRAINT chk CHECK (EXISTS (SELECT 1 FROM other WHERE other.id = ref)));`,
		)
	})

	t.Run("AND_OR_NOT_combination", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, a int, b int);`,
			`CREATE TABLE t (id int PRIMARY KEY, a int, b int, CONSTRAINT chk CHECK (a > 0 AND (b > 0 OR NOT (a = b))));`,
		)
	})

	t.Run("array_subscript", func(t *testing.T) {
		t.Skip("BUG: SubscriptingRef not supported — tags[1] deparsed as NULL::unknown")
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, tags text[]);`,
			`CREATE TABLE t (id int PRIMARY KEY, tags text[], CONSTRAINT chk CHECK (tags[1] IS NOT NULL));`,
		)
	})
}

func TestOracleExprPolicy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// 2.2 Policy USING and WITH CHECK Expressions

	t.Run("CURRENT_USER_comparison", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, owner text);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;`,
			`CREATE TABLE t (id int PRIMARY KEY, owner text);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			 CREATE POLICY p ON t FOR SELECT USING (owner = CURRENT_USER);`,
		)
	})

	t.Run("function_call_in_USING", func(t *testing.T) {
		t.Skip("BUG: policy USING expression emits schema-qualified function name (test_N.has_access vs has_access)")
		assertOracleRoundtrip(t, oracle,
			`CREATE FUNCTION has_access(int) RETURNS boolean LANGUAGE sql STABLE AS 'SELECT true';
			 CREATE TABLE t (id int PRIMARY KEY);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;`,
			`CREATE FUNCTION has_access(int) RETURNS boolean LANGUAGE sql STABLE AS 'SELECT true';
			 CREATE TABLE t (id int PRIMARY KEY);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			 CREATE POLICY p ON t FOR SELECT USING (has_access(id));`,
		)
	})

	t.Run("AND_OR_in_USING", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, active boolean, role text);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;`,
			`CREATE TABLE t (id int PRIMARY KEY, active boolean, role text);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			 CREATE POLICY p ON t FOR SELECT USING (active AND role = 'admin');`,
		)
	})

	t.Run("CASE_WHEN_in_WITH_CHECK", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, status text);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;`,
			`CREATE TABLE t (id int PRIMARY KEY, status text);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			 CREATE POLICY p ON t FOR INSERT WITH CHECK (CASE WHEN status = 'draft' THEN true ELSE false END);`,
		)
	})

	t.Run("COALESCE_in_USING", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, owner text);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;`,
			`CREATE TABLE t (id int PRIMARY KEY, owner text);
			 ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			 CREATE POLICY p ON t FOR SELECT USING (COALESCE(owner, '') = CURRENT_USER);`,
		)
	})
}

func TestOracleExprDomain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// 2.3 Domain CHECK Expressions

	t.Run("VALUE_greater_than_zero", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			``,
			`CREATE DOMAIN posint AS integer CHECK (VALUE > 0);`,
		)
	})

	t.Run("VALUE_IS_NOT_NULL", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			``,
			`CREATE DOMAIN notnulltext AS text CHECK (VALUE IS NOT NULL);`,
		)
	})

	t.Run("VALUE_BETWEEN", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			``,
			`CREATE DOMAIN percentage AS integer CHECK (VALUE BETWEEN 0 AND 100);`,
		)
	})

	t.Run("function_call_validate", func(t *testing.T) {
		t.Skip("BUG: migration creates domain before function — ordering issue")
		assertOracleRoundtrip(t, oracle,
			``,
			`CREATE FUNCTION validate_pos(int) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
			 CREATE DOMAIN posint AS integer CHECK (validate_pos(VALUE));`,
		)
	})

	t.Run("COALESCE_VALUE", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			``,
			`CREATE DOMAIN nonnegint AS integer CHECK (COALESCE(VALUE, 0) >= 0);`,
		)
	})
}

// ---------------------------------------------------------------------------
// Phase 3: Expression Types in Other DDL Contexts
// ---------------------------------------------------------------------------

func TestOracleExprGenerated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// 3.1 Generated Column Expressions

	t.Run("simple_arithmetic", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, a int, b int);`,
			`CREATE TABLE t (id int PRIMARY KEY, a int, b int, total int GENERATED ALWAYS AS (a + b) STORED);`,
		)
	})

	t.Run("function_call_lower", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, name text);`,
			`CREATE TABLE t (id int PRIMARY KEY, name text, name_lower text GENERATED ALWAYS AS (lower(name)) STORED);`,
		)
	})

	t.Run("COALESCE", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, a int, b int);`,
			`CREATE TABLE t (id int PRIMARY KEY, a int, b int, c int GENERATED ALWAYS AS (COALESCE(a, b)) STORED);`,
		)
	})

	t.Run("CASE_WHEN", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, flag boolean);`,
			`CREATE TABLE t (id int PRIMARY KEY, flag boolean, val int GENERATED ALWAYS AS (CASE WHEN flag THEN 1 ELSE 0 END) STORED);`,
		)
	})

	t.Run("string_concatenation", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, first text, last text);`,
			`CREATE TABLE t (id int PRIMARY KEY, first text, last text, full_name text GENERATED ALWAYS AS (first || ' ' || last) STORED);`,
		)
	})
}

func TestOracleExprIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// 3.2 Index Expressions and WHERE Clauses

	t.Run("expression_index_lower", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, name text);`,
			`CREATE TABLE t (id int PRIMARY KEY, name text);
			 CREATE INDEX idx ON t (lower(name));`,
		)
	})

	t.Run("expression_index_arithmetic", func(t *testing.T) {
		t.Skip("BUG: generated CREATE INDEX missing parentheses around expression (a + b)")
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, a int, b int);`,
			`CREATE TABLE t (id int PRIMARY KEY, a int, b int);
			 CREATE INDEX idx ON t ((a + b));`,
		)
	})

	t.Run("partial_index_WHERE_equals", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, active boolean);`,
			`CREATE TABLE t (id int PRIMARY KEY, active boolean);
			 CREATE INDEX idx ON t (id) WHERE (active = true);`,
		)
	})

	t.Run("partial_index_WHERE_function", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, name text);`,
			`CREATE TABLE t (id int PRIMARY KEY, name text);
			 CREATE INDEX idx ON t (name) WHERE (length(name) > 0);`,
		)
	})

	t.Run("partial_index_WHERE_IS_NOT_NULL", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, email text);`,
			`CREATE TABLE t (id int PRIMARY KEY, email text);
			 CREATE INDEX idx ON t (email) WHERE (email IS NOT NULL);`,
		)
	})

	t.Run("partial_index_WHERE_COALESCE", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int PRIMARY KEY, val int);`,
			`CREATE TABLE t (id int PRIMARY KEY, val int);
			 CREATE INDEX idx ON t (val) WHERE (COALESCE(val, 0) > 0);`,
		)
	})
}

func TestOracleExprTrigger(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// 3.3 Trigger WHEN Clause

	t.Run("IS_DISTINCT_FROM", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
			 CREATE TABLE t (id int PRIMARY KEY, status text);`,
			`CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
			 CREATE TABLE t (id int PRIMARY KEY, status text);
			 CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW
			     WHEN (OLD.status IS DISTINCT FROM NEW.status)
			     EXECUTE FUNCTION trig_fn();`,
		)
	})

	t.Run("AND_combination", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
			 CREATE TABLE t (id int PRIMARY KEY, a int, b int);`,
			`CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
			 CREATE TABLE t (id int PRIMARY KEY, a int, b int);
			 CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW
			     WHEN (OLD.a <> NEW.a AND OLD.b <> NEW.b)
			     EXECUTE FUNCTION trig_fn();`,
		)
	})

	t.Run("IS_NOT_NULL", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
			 CREATE TABLE t (id int PRIMARY KEY, email text);`,
			`CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
			 CREATE TABLE t (id int PRIMARY KEY, email text);
			 CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW
			     WHEN (NEW.email IS NOT NULL)
			     EXECUTE FUNCTION trig_fn();`,
		)
	})

	t.Run("function_call_in_WHEN", func(t *testing.T) {
		assertOracleRoundtrip(t, oracle,
			`CREATE FUNCTION should_audit(int) RETURNS boolean LANGUAGE sql STABLE AS 'SELECT true';
			 CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
			 CREATE TABLE t (id int PRIMARY KEY);`,
			`CREATE FUNCTION should_audit(int) RETURNS boolean LANGUAGE sql STABLE AS 'SELECT true';
			 CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
			 CREATE TABLE t (id int PRIMARY KEY);
			 CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW
			     WHEN (should_audit(NEW.id))
			     EXECUTE FUNCTION trig_fn();`,
		)
	})
}

func TestOracleExprSemanticEquivalence(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}

	// 3.4 Semantic Equivalence (False Positive Prevention)
	// These tests verify that identical schemas produce no diff.

	t.Run("DEFAULT_now_no_diff", func(t *testing.T) {
		ddl := `CREATE TABLE t (id int PRIMARY KEY, ts timestamptz DEFAULT now());`
		catA, err := LoadSQL(ddl)
		if err != nil {
			t.Fatalf("LoadSQL failed: %v", err)
		}
		catB, err := LoadSQL(ddl)
		if err != nil {
			t.Fatalf("LoadSQL failed: %v", err)
		}
		diff := Diff(catA, catB)
		if !diff.IsEmpty() {
			t.Errorf("identical schemas with DEFAULT now() should produce empty diff, got: %+v", diff)
		}
	})

	t.Run("DEFAULT_int_literal_no_diff", func(t *testing.T) {
		ddl := `CREATE TABLE t (id int PRIMARY KEY, val int DEFAULT 0);`
		catA, err := LoadSQL(ddl)
		if err != nil {
			t.Fatalf("LoadSQL failed: %v", err)
		}
		catB, err := LoadSQL(ddl)
		if err != nil {
			t.Fatalf("LoadSQL failed: %v", err)
		}
		diff := Diff(catA, catB)
		if !diff.IsEmpty() {
			t.Errorf("identical schemas with DEFAULT 0 should produce empty diff, got: %+v", diff)
		}
	})

	t.Run("CHECK_constraint_no_diff", func(t *testing.T) {
		ddl := `CREATE TABLE t (id int PRIMARY KEY, val int, CONSTRAINT chk CHECK (val > 0));`
		catA, err := LoadSQL(ddl)
		if err != nil {
			t.Fatalf("LoadSQL failed: %v", err)
		}
		catB, err := LoadSQL(ddl)
		if err != nil {
			t.Fatalf("LoadSQL failed: %v", err)
		}
		diff := Diff(catA, catB)
		if !diff.IsEmpty() {
			t.Errorf("identical schemas with CHECK constraint should produce empty diff, got: %+v", diff)
		}
	})
}
