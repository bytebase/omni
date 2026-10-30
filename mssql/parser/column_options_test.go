package parser

import (
	"errors"
	"strings"
	"testing"
)

// columnOptionCases pin the order and repetition of column options that SQL
// Server 2022 accepts in a column definition, as SET PARSEONLY ON reports
// it. omni used to take every option other than COLLATE in any order and
// any number of times. What SQL Server takes:
//
//	column_name data_type [ COLLATE collation_name ]
//	    { GENERATED ALWAYS AS ... { START | END } [ HIDDEN ]
//	    | [ SPARSE | FILESTREAM | SPARSE FILESTREAM | FILESTREAM SPARSE ]
//	      [ MASKED WITH ( ... ) | ENCRYPTED WITH ( ... ) ] }   -- ENCRYPTED not after FILESTREAM
//	    { ROWGUIDCOL | IDENTITY [ ( seed , increment ) ] [ NOT FOR REPLICATION ]
//	    | [ CONSTRAINT name ] { NULL | NOT NULL | DEFAULT expression
//	      | PRIMARY KEY ... | UNIQUE ... | CHECK ... | [ FOREIGN KEY ] REFERENCES ... } } ...
//
// The options after the first line's group come in any order, any number
// of times: SQL Server's parser takes NULL NOT NULL and two IDENTITY
// clauses, which fail later. (omni does not take the FOREIGN KEY keywords
// of an inline foreign key, nor an inline INDEX, yet.) near is the token
// omni's syntax error is at.
var columnOptionCases = []struct {
	sql    string
	accept bool
	near   string
}{
	{sql: "CREATE TABLE dbo.t9 (a VARBINARY(MAX) FILESTREAM)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a VARBINARY(MAX) SPARSE FILESTREAM NULL)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a VARBINARY(MAX) FILESTREAM SPARSE NULL)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a VARBINARY(MAX) COLLATE Latin1_General_CI_AS SPARSE FILESTREAM MASKED WITH (FUNCTION = 'default()') NOT NULL)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') NULL)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a DATETIME2 COLLATE Latin1_General_CI_AS GENERATED ALWAYS AS ROW END HIDDEN NOT NULL)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a INT IDENTITY NOT FOR REPLICATION PRIMARY KEY)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a INT IDENTITY(1, 1) NOT FOR REPLICATION IDENTITY NOT FOR REPLICATION)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a INT NOT NULL REFERENCES dbo.t (a) NOT FOR REPLICATION DEFAULT 1)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a INT NULL NOT NULL)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a UNIQUEIDENTIFIER ROWGUIDCOL DEFAULT NEWID() ROWGUIDCOL)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 'x' COLLATE Latin1_General_BIN NOT NULL)", accept: true},
	{sql: "CREATE TABLE dbo.t9 (a INT CHECK NOT FOR REPLICATION (a > 0) CONSTRAINT nn NOT NULL UNIQUE)", accept: true},
	{sql: "ALTER TABLE dbo.t ADD c INT SPARSE NULL CONSTRAINT ck CHECK (c > 0)", accept: true},

	{sql: "CREATE TABLE dbo.t9 (a VARBINARY(MAX) FILESTREAM FILESTREAM)", near: "FILESTREAM"},
	{sql: "CREATE TABLE dbo.t9 (a VARBINARY(MAX) NULL FILESTREAM)", near: "FILESTREAM"},
	{sql: "CREATE TABLE dbo.t9 (a VARBINARY(MAX) MASKED WITH (FUNCTION = 'default()') FILESTREAM)", near: "FILESTREAM"},
	{sql: "CREATE TABLE dbo.t9 (a VARBINARY(MAX) FILESTREAM ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))", near: "ENCRYPTED"},
	{sql: "CREATE TABLE dbo.t9 (a VARBINARY(MAX) SPARSE FILESTREAM ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))", near: "ENCRYPTED"},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE SPARSE)", near: "SPARSE"},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL SPARSE)", near: "SPARSE"},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL SPARSE)", near: "SPARSE"},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') SPARSE)", near: "SPARSE"},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') MASKED WITH (FUNCTION = 'default()'))", near: "MASKED"},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))", near: "ENCRYPTED"},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') MASKED WITH (FUNCTION = 'default()'))", near: "MASKED"},
	{sql: "CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 MASKED WITH (FUNCTION = 'default()'))", near: "MASKED"},
	{sql: "CREATE TABLE dbo.t9 (a DATETIME2 SPARSE GENERATED ALWAYS AS ROW START)", near: "GENERATED"},
	{sql: "CREATE TABLE dbo.t9 (a DATETIME2 GENERATED ALWAYS AS ROW START SPARSE)", near: "SPARSE"},
	{sql: "CREATE TABLE dbo.t9 (a DATETIME2 GENERATED ALWAYS AS ROW START MASKED WITH (FUNCTION = 'default()'))", near: "MASKED"},
	{sql: "CREATE TABLE dbo.t9 (a DATETIME2 NOT NULL GENERATED ALWAYS AS ROW START)", near: "GENERATED"},
	{sql: "CREATE TABLE dbo.t9 (a DATETIME2 HIDDEN)", near: "HIDDEN"},
	{sql: "CREATE TABLE dbo.t9 (a DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL HIDDEN)", near: "HIDDEN"},
	{sql: "CREATE TABLE dbo.t9 (a DATETIME2 GENERATED ALWAYS AS ROW START HIDDEN HIDDEN)", near: "HIDDEN"},
	{sql: "CREATE TABLE dbo.t9 (a INT NOT FOR REPLICATION)", near: "NOT"},
	{sql: "CREATE TABLE dbo.t9 (a INT NOT FOR REPLICATION IDENTITY)", near: "NOT"},
	{sql: "CREATE TABLE dbo.t9 (a INT IDENTITY NULL NOT FOR REPLICATION)", near: "NOT"},
	{sql: "CREATE TABLE dbo.t9 (a INT IDENTITY NOT FOR REPLICATION NOT FOR REPLICATION)", near: "NOT"},
	{sql: "CREATE TABLE dbo.t9 (a INT CHECK (a > 0) NOT FOR REPLICATION)", near: "NOT"},
	{sql: "ALTER TABLE dbo.t ADD c INT NOT NULL SPARSE", near: "SPARSE"},
	{sql: "ALTER TABLE dbo.t ADD c INT NOT FOR REPLICATION", near: "NOT"},
	{sql: "DECLARE @t TABLE (a INT NOT NULL HIDDEN)", near: "HIDDEN"},
	{sql: "CREATE TYPE dbo.tt AS TABLE (a INT NULL SPARSE)", near: "SPARSE"},
	{sql: "CREATE FUNCTION dbo.f() RETURNS @t TABLE (a INT NOT FOR REPLICATION) AS BEGIN RETURN END", near: "NOT"},
}

// TestColumnOptionOrderOracle: omni accepts the column option orders SQL
// Server accepts and rejects the others at the option out of place.
func TestColumnOptionOrderOracle(t *testing.T) {
	oracle := startParserOracle(t)
	for _, tc := range columnOptionCases {
		t.Run(tc.sql, func(t *testing.T) {
			serverErr, err := oracle.parseError(tc.sql)
			if err != nil {
				t.Fatal(err)
			}
			if (serverErr == nil) != tc.accept {
				t.Fatalf("SQL Server: %v, want accept = %v", serverErr, tc.accept)
			}
			_, omniErr := Parse(tc.sql)
			if (omniErr == nil) != tc.accept {
				t.Fatalf("omni: %v, want accept = %v", omniErr, tc.accept)
			}
			if tc.accept {
				return
			}
			var pe *ParseError
			if !errors.As(omniErr, &pe) {
				t.Fatalf("omni error %v is not a *ParseError", omniErr)
			}
			if !strings.HasPrefix(strings.ToUpper(tc.sql[pe.Position:]), tc.near) {
				t.Errorf("omni error %q at offset %d, want it at %s", pe.Message, pe.Position, tc.near)
			}
		})
	}
}

// columnOptionAtoms are the column options TestColumnOptionPairsOracle puts
// in every order, by twos.
var columnOptionAtoms = []string{
	"COLLATE Latin1_General_CI_AS",
	"SPARSE",
	"MASKED WITH (FUNCTION = 'default()')",
	"ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')",
	"GENERATED ALWAYS AS ROW START",
	"HIDDEN",
	"NOT NULL",
	"NULL",
	"ROWGUIDCOL",
	"IDENTITY(1, 1)",
	"NOT FOR REPLICATION",
	"DEFAULT 1",
	"CONSTRAINT df DEFAULT 1",
	"CONSTRAINT nn NOT NULL",
	"PRIMARY KEY",
	"UNIQUE",
	"CHECK (a > 0)",
	"REFERENCES dbo.t (a)",
}

// TestColumnOptionPairsOracle: for every option of columnOptionAtoms alone
// and every ordered pair, CREATE TABLE and ALTER TABLE ... ADD agree with SQL
// Server. FILESTREAM, FOREIGN KEY, and an inline INDEX are left out:
// FILESTREAM depends on the data type, and omni rejects the other two.
func TestColumnOptionPairsOracle(t *testing.T) {
	oracle := startParserOracle(t)
	var options []string
	for _, a := range columnOptionAtoms {
		options = append(options, a)
		for _, b := range columnOptionAtoms {
			options = append(options, a+" "+b)
		}
	}
	for _, format := range []string{"CREATE TABLE dbo.t9 (a VARCHAR(10) %s)", "ALTER TABLE dbo.t ADD c VARCHAR(10) %s"} {
		for _, opts := range options {
			sql := strings.Replace(format, "%s", opts, 1)
			serverErr, err := oracle.parseError(sql)
			if err != nil {
				t.Fatal(err)
			}
			if _, omniErr := Parse(sql); (omniErr == nil) != (serverErr == nil) {
				t.Errorf("%s\n\tSQL Server: %v\n\tomni: %v", sql, serverErr, omniErr)
			}
		}
	}
}
