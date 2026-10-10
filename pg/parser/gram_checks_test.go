package parser

import "testing"

// TestConstraintAttributesMatchPG pins the constraint attributes each kind
// takes. gram.y's processCASbits refuses DEFERRABLE on a check or a NOT
// NULL, NOT VALID on a key, an exclusion constraint, a constraint trigger,
// or ALTER CONSTRAINT, and NO INHERIT on those and on a foreign key;
// ConstraintAttributeSpec refuses conflicting attributes; and a column's
// constraint list takes no NOT VALID, and NO INHERIT only on a check.
func TestConstraintAttributesMatchPG(t *testing.T) {
	assertSyntaxMatchesPG(t, []syntaxCase{
		{"CREATE TABLE n (id int, PRIMARY KEY (id) NOT VALID)", false},  // PRIMARY KEY constraints cannot be marked NOT VALID
		{"CREATE TABLE n (id int, PRIMARY KEY (id) NO INHERIT)", false}, // PRIMARY KEY constraints cannot be marked NO INHERIT
		{"CREATE TABLE n (id int, PRIMARY KEY (id) DEFERRABLE INITIALLY DEFERRED)", true},
		{"ALTER TABLE t ADD CONSTRAINT k PRIMARY KEY USING INDEX i NOT VALID", false}, // PRIMARY KEY constraints cannot be marked NOT VALID
		{"CREATE TABLE n (id int, UNIQUE (id) NOT VALID)", false},                     // UNIQUE constraints cannot be marked NOT VALID
		{"ALTER TABLE t ADD CONSTRAINT k UNIQUE USING INDEX i NO INHERIT", false},     // UNIQUE constraints cannot be marked NO INHERIT
		{"CREATE TABLE n (id int, UNIQUE (id) DEFERRABLE)", true},
		{"CREATE TABLE n (id int, EXCLUDE (id WITH =) NOT VALID)", false},  // EXCLUDE constraints cannot be marked NOT VALID
		{"CREATE TABLE n (id int, EXCLUDE (id WITH =) NO INHERIT)", false}, // EXCLUDE constraints cannot be marked NO INHERIT
		{"CREATE TABLE n (id int, EXCLUDE (id WITH =) INITIALLY DEFERRED)", true},
		{"CREATE TABLE n (id int, CHECK (id > 0) DEFERRABLE)", false},         // CHECK constraints cannot be marked DEFERRABLE
		{"CREATE TABLE n (id int, CHECK (id > 0) INITIALLY DEFERRED)", false}, // CHECK constraints cannot be marked DEFERRABLE
		{"CREATE TABLE n (id int, CHECK (id > 0) NOT VALID NO INHERIT)", true},
		{"CREATE TABLE n (id int, CHECK (id > 0) NOT DEFERRABLE INITIALLY IMMEDIATE)", true},
		{"CREATE TABLE n (id int, FOREIGN KEY (id) REFERENCES t NO INHERIT)", false}, // FOREIGN KEY constraints cannot be marked NO INHERIT
		{"CREATE TABLE n (id int, FOREIGN KEY (id) REFERENCES t NOT VALID DEFERRABLE)", true},
		{"ALTER TABLE t ADD CONSTRAINT x CHECK (x > 0) DEFERRABLE", false},  // CHECK constraints cannot be marked DEFERRABLE
		{"ALTER TABLE t ADD CONSTRAINT x PRIMARY KEY (x) NOT VALID", false}, // PRIMARY KEY constraints cannot be marked NOT VALID
		{"ALTER TABLE t ADD CONSTRAINT x FOREIGN KEY (x) REFERENCES p NOT VALID", true},
		{"CREATE TABLE n OF ty (PRIMARY KEY (id) NOT VALID)", false},                   // PRIMARY KEY constraints cannot be marked NOT VALID
		{"CREATE TABLE n PARTITION OF t (CHECK (id > 0) DEFERRABLE) DEFAULT", false},   // CHECK constraints cannot be marked DEFERRABLE
		{"CREATE FOREIGN TABLE n (id int, CHECK (id > 0) DEFERRABLE) SERVER s", false}, // CHECK constraints cannot be marked DEFERRABLE

		// Conflicting attributes, refused at the one that conflicts.
		{"CREATE TABLE n (id int, PRIMARY KEY (id) NOT DEFERRABLE INITIALLY DEFERRED)", false},      // constraint declared INITIALLY DEFERRED must be DEFERRABLE
		{"CREATE TABLE n (id int, PRIMARY KEY (id) DEFERRABLE NOT DEFERRABLE)", false},              // conflicting constraint properties
		{"CREATE TABLE n (id int, PRIMARY KEY (id) INITIALLY IMMEDIATE INITIALLY DEFERRED)", false}, // conflicting constraint properties
		{"CREATE TABLE n (id int, PRIMARY KEY (id) DEFERRABLE DEFERRABLE)", true},
		{"CREATE TABLE n (id int, PRIMARY KEY (id) INITIALLY)", false},

		// ALTER CONSTRAINT and constraint triggers change only deferrability.
		{"ALTER TABLE t ALTER CONSTRAINT c DEFERRABLE INITIALLY DEFERRED", true},
		{"ALTER TABLE t ALTER CONSTRAINT c NOT VALID", false},  // FOREIGN KEY constraints cannot be marked NOT VALID
		{"ALTER TABLE t ALTER CONSTRAINT c NO INHERIT", false}, // FOREIGN KEY constraints cannot be marked NO INHERIT
		{"CREATE CONSTRAINT TRIGGER tr AFTER INSERT ON t DEFERRABLE FOR EACH ROW EXECUTE FUNCTION f()", true},
		{"CREATE CONSTRAINT TRIGGER tr AFTER INSERT ON t NOT VALID FOR EACH ROW EXECUTE FUNCTION f()", false},  // TRIGGER constraints cannot be marked NOT VALID
		{"CREATE CONSTRAINT TRIGGER tr AFTER INSERT ON t NO INHERIT FOR EACH ROW EXECUTE FUNCTION f()", false}, // TRIGGER constraints cannot be marked NO INHERIT

		// A domain's check takes NOT VALID and NO INHERIT, its NOT NULL
		// only NO INHERIT.
		{"ALTER DOMAIN d ADD CONSTRAINT c CHECK (VALUE > 0) NOT VALID NO INHERIT", true},
		{"ALTER DOMAIN d ADD CONSTRAINT c CHECK (VALUE > 0) DEFERRABLE", false}, // CHECK constraints cannot be marked DEFERRABLE
		{"ALTER DOMAIN d ADD CONSTRAINT c NOT NULL NO INHERIT", true},
		{"ALTER DOMAIN d ADD CONSTRAINT c NOT NULL NOT VALID", false},  // NOT NULL constraints cannot be marked NOT VALID
		{"ALTER DOMAIN d ADD CONSTRAINT c NOT NULL DEFERRABLE", false}, // NOT NULL constraints cannot be marked DEFERRABLE

		// A column's constraint list.
		{"CREATE TABLE n (id int REFERENCES t DEFERRABLE INITIALLY DEFERRED)", true},
		{"CREATE TABLE n (id int REFERENCES t NOT VALID)", false},
		{"CREATE TABLE n (id int REFERENCES t DEFERRABLE NOT VALID)", false},
		{"CREATE TABLE n (id int REFERENCES t NO INHERIT)", false},
		{"CREATE TABLE n (id int CHECK (id > 0) NO INHERIT)", true},
		{"CREATE TABLE n (id int CHECK (id > 0) NOT VALID)", false},
		{"CREATE TABLE n (id int NOT VALID)", false},
		{"ALTER TABLE t ADD COLUMN id int REFERENCES t NOT VALID", false},
		{"CREATE TABLE n OF ty (id WITH OPTIONS NOT VALID)", false},
		{"CREATE DOMAIN d AS int CHECK (VALUE > 0) NOT VALID", false},
		{"CREATE DOMAIN d AS int NOT NULL NO INHERIT", false},
	})
}

// TestCreateSchemaIfNotExistsMatchesPG pins that the IF NOT EXISTS forms of
// CREATE SCHEMA take no schema elements, refused at the first one.
func TestCreateSchemaIfNotExistsMatchesPG(t *testing.T) {
	assertSyntaxMatchesPG(t, []syntaxCase{
		{"CREATE SCHEMA IF NOT EXISTS s CREATE TABLE n (id int)", false},                                     // CREATE SCHEMA IF NOT EXISTS cannot include schema elements
		{"CREATE SCHEMA IF NOT EXISTS AUTHORIZATION r CREATE TABLE n (id int)", false},                       // CREATE SCHEMA IF NOT EXISTS cannot include schema elements
		{"CREATE SCHEMA IF NOT EXISTS s AUTHORIZATION r CREATE VIEW v AS SELECT 1 CREATE SEQUENCE q", false}, // CREATE SCHEMA IF NOT EXISTS cannot include schema elements
		{"CREATE SCHEMA IF NOT EXISTS s", true},
		{"CREATE SCHEMA IF NOT EXISTS s AUTHORIZATION r", true},
		{"CREATE SCHEMA IF NOT EXISTS AUTHORIZATION r", true},
		{"CREATE SCHEMA s CREATE TABLE n (id int)", true},
		{"CREATE SCHEMA AUTHORIZATION r CREATE TABLE n (id int)", true},
	})
}
