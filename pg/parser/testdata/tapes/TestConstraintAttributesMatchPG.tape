-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestConstraintAttributesMatchPG default#1
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, PRIMARY KEY (id) NOT VALID)
< error
| 0A000 PRIMARY KEY constraints cannot be marked NOT VALID
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, PRIMARY KEY (id) NO INHERIT)
< error
| 0A000 PRIMARY KEY constraints cannot be marked NO INHERIT
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, PRIMARY KEY (id) DEFERRABLE INITIALLY DEFERRED)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CONSTRAINT k PRIMARY KEY USING INDEX i NOT VALID
< error
| 0A000 PRIMARY KEY constraints cannot be marked NOT VALID
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, UNIQUE (id) NOT VALID)
< error
| 0A000 UNIQUE constraints cannot be marked NOT VALID
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CONSTRAINT k UNIQUE USING INDEX i NO INHERIT
< error
| 0A000 UNIQUE constraints cannot be marked NO INHERIT
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, UNIQUE (id) DEFERRABLE)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, EXCLUDE (id WITH =) NOT VALID)
< error
| 0A000 EXCLUDE constraints cannot be marked NOT VALID
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, EXCLUDE (id WITH =) NO INHERIT)
< error
| 0A000 EXCLUDE constraints cannot be marked NO INHERIT
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, EXCLUDE (id WITH =) INITIALLY DEFERRED)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, CHECK (id > 0) DEFERRABLE)
< error
| 0A000 CHECK constraints cannot be marked DEFERRABLE
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, CHECK (id > 0) INITIALLY DEFERRED)
< error
| 0A000 CHECK constraints cannot be marked DEFERRABLE
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, CHECK (id > 0) NOT VALID NO INHERIT)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, CHECK (id > 0) NOT DEFERRABLE INITIALLY IMMEDIATE)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, FOREIGN KEY (id) REFERENCES t NO INHERIT)
< error
| 0A000 FOREIGN KEY constraints cannot be marked NO INHERIT
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, FOREIGN KEY (id) REFERENCES t NOT VALID DEFERRABLE)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CONSTRAINT x CHECK (x > 0) DEFERRABLE
< error
| 0A000 CHECK constraints cannot be marked DEFERRABLE
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CONSTRAINT x PRIMARY KEY (x) NOT VALID
< error
| 0A000 PRIMARY KEY constraints cannot be marked NOT VALID
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CONSTRAINT x FOREIGN KEY (x) REFERENCES p NOT VALID
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE n OF ty (PRIMARY KEY (id) NOT VALID)
< error
| 0A000 PRIMARY KEY constraints cannot be marked NOT VALID
> exec
| SELECT 1/0;
| CREATE TABLE n PARTITION OF t (CHECK (id > 0) DEFERRABLE) DEFAULT
< error
| 0A000 CHECK constraints cannot be marked DEFERRABLE
> exec
| SELECT 1/0;
| CREATE FOREIGN TABLE n (id int, CHECK (id > 0) DEFERRABLE) SERVER s
< error
| 0A000 CHECK constraints cannot be marked DEFERRABLE
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, PRIMARY KEY (id) NOT DEFERRABLE INITIALLY DEFERRED)
< error
| 42601 constraint declared INITIALLY DEFERRED must be DEFERRABLE
| position: 69
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, PRIMARY KEY (id) DEFERRABLE NOT DEFERRABLE)
< error
| 42601 conflicting constraint properties
| position: 65
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, PRIMARY KEY (id) INITIALLY IMMEDIATE INITIALLY DEFERRED)
< error
| 42601 conflicting constraint properties
| position: 74
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, PRIMARY KEY (id) DEFERRABLE DEFERRABLE)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE n (id int, PRIMARY KEY (id) INITIALLY)
< error
| 42601 syntax error at or near ")"
| position: 63
> exec
| SELECT 1/0;
| ALTER TABLE t ALTER CONSTRAINT c DEFERRABLE INITIALLY DEFERRED
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ALTER CONSTRAINT c NOT VALID
< error
| 0A000 FOREIGN KEY constraints cannot be marked NOT VALID
> exec
| SELECT 1/0;
| ALTER TABLE t ALTER CONSTRAINT c NO INHERIT
< error
| 0A000 FOREIGN KEY constraints cannot be marked NO INHERIT
> exec
| SELECT 1/0;
| CREATE CONSTRAINT TRIGGER tr AFTER INSERT ON t DEFERRABLE FOR EACH ROW EXECUTE FUNCTION f()
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE CONSTRAINT TRIGGER tr AFTER INSERT ON t NOT VALID FOR EACH ROW EXECUTE FUNCTION f()
< error
| 0A000 TRIGGER constraints cannot be marked NOT VALID
> exec
| SELECT 1/0;
| CREATE CONSTRAINT TRIGGER tr AFTER INSERT ON t NO INHERIT FOR EACH ROW EXECUTE FUNCTION f()
< error
| 0A000 TRIGGER constraints cannot be marked NO INHERIT
> exec
| SELECT 1/0;
| ALTER DOMAIN d ADD CONSTRAINT c CHECK (VALUE > 0) NOT VALID NO INHERIT
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER DOMAIN d ADD CONSTRAINT c CHECK (VALUE > 0) DEFERRABLE
< error
| 0A000 CHECK constraints cannot be marked DEFERRABLE
> exec
| SELECT 1/0;
| ALTER DOMAIN d ADD CONSTRAINT c NOT NULL NO INHERIT
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER DOMAIN d ADD CONSTRAINT c NOT NULL NOT VALID
< error
| 0A000 NOT NULL constraints cannot be marked NOT VALID
> exec
| SELECT 1/0;
| ALTER DOMAIN d ADD CONSTRAINT c NOT NULL DEFERRABLE
< error
| 0A000 NOT NULL constraints cannot be marked DEFERRABLE
> exec
| SELECT 1/0;
| CREATE TABLE n (id int REFERENCES t DEFERRABLE INITIALLY DEFERRED)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE n (id int REFERENCES t NOT VALID)
< error
| 42601 syntax error at or near "VALID"
| position: 53
> exec
| SELECT 1/0;
| CREATE TABLE n (id int REFERENCES t DEFERRABLE NOT VALID)
< error
| 42601 syntax error at or near "VALID"
| position: 64
> exec
| SELECT 1/0;
| CREATE TABLE n (id int REFERENCES t NO INHERIT)
< error
| 42601 syntax error at or near "NO"
| position: 49
> exec
| SELECT 1/0;
| CREATE TABLE n (id int CHECK (id > 0) NO INHERIT)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE n (id int CHECK (id > 0) NOT VALID)
< error
| 42601 syntax error at or near "VALID"
| position: 55
> exec
| SELECT 1/0;
| CREATE TABLE n (id int NOT VALID)
< error
| 42601 syntax error at or near "VALID"
| position: 40
> exec
| SELECT 1/0;
| ALTER TABLE t ADD COLUMN id int REFERENCES t NOT VALID
< error
| 42601 syntax error at or near "VALID"
| position: 62
> exec
| SELECT 1/0;
| CREATE TABLE n OF ty (id WITH OPTIONS NOT VALID)
< error
| 42601 syntax error at or near "VALID"
| position: 55
> exec
| SELECT 1/0;
| CREATE DOMAIN d AS int CHECK (VALUE > 0) NOT VALID
< error
| 42601 syntax error at or near "VALID"
| position: 58
> exec
| SELECT 1/0;
| CREATE DOMAIN d AS int NOT NULL NO INHERIT
< error
| 42601 syntax error at or near "NO"
| position: 45
