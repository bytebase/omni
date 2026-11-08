package parser

import "testing"

// TestQuotedVariableNames pins that omni accepts the quoted user- and
// system-variable name forms TiDB v8.5.5 accepts (the §3.4 rest-corpus variable
// cluster), that the lexer fold still keeps usernames working (`'u'@'host'`
// shares the `@'...'` shape), and that a form TiDB rejects stays rejected.
func TestQuotedVariableNames(t *testing.T) {
	accept := []string{
		// user variables with quoted names, including := and a comparison
		"select @'a':=1",
		"select @`a``b`=4",
		"select @\"aBC\">1",
		"select @'a b'",
		"set @'a' = 1",
		// system variables with a quoted name component
		"select @@session.'C'",
		"select @@local.\"aBc\"",
		"select @@'c'",
		"set @@session.'C' = 1",
		// usernames, which must keep parsing after the lexer fold
		"CREATE USER 'u'@'host'",
		"CREATE USER 'u'@'%'",
		"CREATE USER u@localhost",
		"GRANT SELECT ON *.* TO 'u'@'host'",
		"DROP USER 'u'@'host'",
		"CREATE DEFINER = 'u'@'host' VIEW v AS SELECT 1",
		"CLONE INSTANCE FROM 'u'@'host':3306 IDENTIFIED BY 'p'",
		// An empty quoted host `@''` is a valid empty host, distinct from a bare
		// `@` that needs a following token. The fold makes both Str "@", so the
		// host consumers tell them apart by the source span.
		"CREATE DEFINER = 'u'@'' VIEW v AS SELECT 1",
		"CREATE USER 'u'@''",
		"CREATE USER 'u'@'' IDENTIFIED BY 'p'",
		"CREATE USER 'u'@''IDENTIFIED BY 'p'",
		"GRANT SELECT ON *.* TO 'u'@''",
		"DROP USER 'u'@''",
		"SELECT @''",
	}
	for _, sql := range accept {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse rejected %q: %v", sql, err)
		}
	}

	// TiDB answers 1064 to a dotted chain of quoted system-variable components;
	// a quoted name is the whole name, so omni must reject it too.
	reject := []string{
		"select @@session.\"a\".\"b\"",
		// CLONE INSTANCE needs a non-empty host: TiDB v8.5.5 answers 1064 to an
		// empty one, so omni must too (the following colon has nothing before it).
		"CLONE INSTANCE FROM 'u'@'':3306 IDENTIFIED BY 'p'",
	}
	for _, sql := range reject {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse accepted %q; TiDB v8.5.5 rejects it", sql)
		}
	}
}
