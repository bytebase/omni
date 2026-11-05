package parser

import (
	"slices"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Role names, as TiDB v8.5.5 reads them, every verdict measured with PREPARE on
// pingcap/tidb:v8.5.5 on 2026-11-05: a 1064 is the parser's refusal, and any other answer
// means the text parsed.
//
// A role written on its own is RoleNameString, parser.y:11135: a string or the bare
// identifier token. That token is not the same set as omni's tokIDENT, because this lexer
// holds words as keywords that upstream does not have: `CREATE ROLE loop` parses, where
// `CREATE ROLE transaction` is a 1064 though both are unreserved. tidbPlainIdentifiers lists
// them. Only role@host reaches StringName, which takes any unreserved word.
//
// A privilege written with a keyword is PrivType, parser.y:14512, a closed list that
// tidbPrivilegeNames holds, and a longer one is ExtendedPriv, a run of the same bare
// identifiers. PROXY is neither: GRANT spells it out, and REVOKE PROXY is a 1064.
//
// Not in this table, because omni still disagrees, all by accepting forms TiDB does not have:
// `ALTER USER u DEFAULT ROLE r`, `CREATE USER u DEFAULT ROLE r`, `GRANT a TO u WITH ADMIN
// OPTION`, `GRANT SELECT ON t TO u AS a WITH ROLE r`, `REVOKE IF EXISTS` and `CREATE USER IF
// a`.

// roleNamePositions are the positions the oracle sweep tries every keyword in.
var roleNamePositions = []string{
	"create role %s",
	"create role %s@h",
	"drop role %s",
	"set role %s",
	"set default role %s to u",
	"grant %s to u",
	"revoke %s from u",
	"grant %s on *.* to u",
	"show grants for u using %s",
}

// TestRoleNameForms holds omni to the recorded verdicts without a server, so a run with no
// Docker still reports a widened or narrowed grammar.
func TestRoleNameForms(t *testing.T) {
	for _, c := range roleNameCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestPlainIdentifiersAreKeywords pins that the table names only words the lexer holds as
// keywords and that the rest of the grammar treats as names, so a stale entry fails here.
func TestPlainIdentifiersAreKeywords(t *testing.T) {
	for word := range tidbPlainIdentifiers {
		if word != strings.ToLower(word) {
			t.Errorf("%q is not in lower case, and the lookup lowers what it reads", word)
		}
		tok, ok := keywords[word]
		if !ok {
			t.Errorf("%q is not a keyword here, so it does not need listing", word)
			continue
		}
		if isReserved(tok) {
			t.Errorf("%q is reserved here, and a plain identifier never is", word)
		}
	}
}

// TestRoleNameTrees pins where the names landed, including a run of identifiers read as one
// privilege name, which a verdict cannot tell from a list.
func TestRoleNameTrees(t *testing.T) {
	parse := func(t *testing.T, sql string) nodes.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		return list.Items[0]
	}
	if got := parse(t, "create role loop, get").(*nodes.CreateRoleStmt).Roles; !slices.Equal(got, []string{"loop", "get"}) {
		t.Errorf("CREATE ROLE roles = %v", got)
	}
	if got := parse(t, "drop role loop").(*nodes.DropRoleStmt).Roles; !slices.Equal(got, []string{"loop"}) {
		t.Errorf("DROP ROLE roles = %v", got)
	}
	if got := parse(t, "set role loop").(*nodes.SetRoleStmt).Roles; !slices.Equal(got, []string{"loop"}) {
		t.Errorf("SET ROLE roles = %v", got)
	}
	grant, ok := parse(t, "grant loop, a to u").(*nodes.GrantRoleStmt)
	if !ok {
		t.Fatalf("GRANT with TO is %T", parse(t, "grant loop, a to u"))
	}
	if !slices.Equal(grant.Roles, []string{"LOOP", "A"}) {
		t.Errorf("GRANT roles = %v", grant.Roles)
	}
	priv, ok := parse(t, "grant loop get on *.* to u").(*nodes.GrantStmt)
	if !ok {
		t.Fatalf("GRANT with ON is %T", parse(t, "grant loop get on *.* to u"))
	}
	if !slices.Equal(priv.Privileges, []string{"LOOP GET"}) {
		t.Errorf("a run of identifiers read as %v, want one privilege", priv.Privileges)
	}
	proxy := parse(t, "grant proxy on 'a'@'b' to u").(*nodes.GrantStmt)
	if proxy.ProxyUser == "" {
		t.Errorf("GRANT PROXY has no proxy user: %+v", proxy)
	}
	temporary := parse(t, "grant create temporary tables on *.* to u").(*nodes.GrantStmt)
	if !slices.Equal(temporary.Privileges, []string{"CREATE TEMPORARY TABLES"}) {
		t.Errorf("privileges = %v", temporary.Privileges)
	}
}

var roleNameCases = []tableSampleCase{
	{"create role 'transaction'", true},
	{"create role LOOP", true},
	{"create role Loop, GET", true},
	{"create role TRANSACTION", false},
	{"create role `loop`", true},
	{"create role `select`", true},
	{"create role `transaction`", true},
	{"create role a, loop", true},
	{"create role a, transaction", false},
	{"create role a@h, loop", true},
	{"create role count", true},
	{"create role count@h", true},
	{"create role if", false},
	{"create role if exists a", false},
	{"create role if not a", false},
	{"create role if not exists a", true},
	{"create role if not exists loop", true},
	{"create role if@h", false},
	{"create role loop", true},
	{"create role loop, a", true},
	{"create role loop@'h'", true},
	{"create role loop@h", true},
	{"create role select", false},
	{"create role select@h", false},
	{"create role transaction@h", true},
	{"drop role LOOP", true},
	{"drop role a, loop", true},
	{"drop role if", false},
	{"drop role if a", false},
	{"drop role if exists a", true},
	{"drop role if exists loop", true},
	{"drop role if not exists a", false},
	{"drop role if@h", false},
	{"drop role loop", true},
	{"drop role transaction", false},
	{"drop role transaction@h", true},
	{"grant 'loop' to u", true},
	{"grant Create View on *.* to u", true},
	{"grant LOOP to u", true},
	{"grant Loop Get on *.* to u", true},
	{"grant PROXY on a to u", true},
	{"grant SELECT on *.* to u", true},
	{"grant TRANSACTION to u", false},
	{"grant `admin` on *.* to u", true},
	{"grant `loop` to u", true},
	{"grant `select` on *.* to u", true},
	{"grant `select`(a) on t to u", false},
	{"grant `transaction` to u", true},
	{"grant a loop on *.* to u", true},
	{"grant a loop to u", true},
	{"grant a(b) on t to u", false},
	{"grant a, all to u", false},
	{"grant a, loop to u", true},
	{"grant add on *.* to u", false},
	{"grant add(a) on t to u", false},
	{"grant admin on *.* to u", false},
	{"grant all on *.* to u", true},
	{"grant all on *.* to u with grant option", true},
	{"grant all on t to u", true},
	{"grant all privileges on *.* to u", true},
	{"grant all privileges on t to u", true},
	{"grant all privileges to u", false},
	{"grant all to u", false},
	{"grant all, a to u", false},
	{"grant alter on *.* to u", true},
	{"grant alter routine on *.* to u", true},
	{"grant alter routine to u", false},
	{"grant alter to u", false},
	{"grant backup_admin on *.* to u", true},
	{"grant comment on *.* to u1", false},
	{"grant comment, select on *.* to u1", false},
	{"grant concat to u", true},
	{"grant config on *.* to u", true},
	{"grant config to u", false},
	{"grant count on *.* to u", true},
	{"grant count to u", true},
	{"grant count(a) on t to u", false},
	{"grant create on *.* to u", true},
	{"grant create role on *.* to u", true},
	{"grant create role to u", false},
	{"grant create routine on *.* to u", true},
	{"grant create routine to u", false},
	{"grant create tablespace on *.* to u", true},
	{"grant create tablespace to u", false},
	{"grant create temporary on *.* to u", false},
	{"grant create temporary tables on *.* to u", true},
	{"grant create temporary tables to u", false},
	{"grant create to u", false},
	{"grant create user on *.* to u", true},
	{"grant create user to u", false},
	{"grant create view on *.* to u", true},
	{"grant create view to u", false},
	{"grant delete on *.* to u", true},
	{"grant delete to u", false},
	{"grant drop on *.* to u", true},
	{"grant drop role on *.* to u", true},
	{"grant drop role to u", false},
	{"grant drop to u", false},
	{"grant event on *.* to u", true},
	{"grant event to u", false},
	{"grant execute on *.* to u", true},
	{"grant execute to u", false},
	{"grant file on *.* to u", true},
	{"grant file to u", false},
	{"grant grant on *.* to u", false},
	{"grant grant option on *.* to u", true},
	{"grant grant option to u", false},
	{"grant grant to u", false},
	{"grant grant, select to u", false},
	{"grant index on *.* to u", true},
	{"grant index to u", false},
	{"grant insert on *.* to u", true},
	{"grant insert to u", false},
	{"grant insert(a), update(b) on t to u", true},
	{"grant invoke `lambda` on *.* to u", true},
	{"grant invoke lambda on *.* to u", true},
	{"grant invoke loop on *.* to u", true},
	{"grant load from s3 on *.* to u", true},
	{"grant load on *.* to u", false},
	{"grant lock on *.* to u", false},
	{"grant lock tables on *.* to u", true},
	{"grant lock tables to u", false},
	{"grant loop a on *.* to u", true},
	{"grant loop a to u", true},
	{"grant loop get on *.* to u", true},
	{"grant loop on *.* to u", true},
	{"grant loop to loop", true},
	{"grant loop to transaction", true},
	{"grant loop to u", true},
	{"grant loop transaction on *.* to u", false},
	{"grant loop(a) on t to u", false},
	{"grant loop, a to u", true},
	{"grant loop, select on *.* to u", true},
	{"grant loop, transaction to u", false},
	{"grant loop@'h' to u", true},
	{"grant loop@h to u", true},
	{"grant process on *.* to u", true},
	{"grant process to u", false},
	{"grant proxy on 'a'@'b' to u", true},
	{"grant proxy on *.* to u", false},
	{"grant proxy on a to u", true},
	{"grant proxy to u", false},
	{"grant proxy(a) on t to u", false},
	{"grant proxy, select on *.* to u", false},
	{"grant references on *.* to u", true},
	{"grant references to u", false},
	{"grant reload on *.* to u", true},
	{"grant reload to u", false},
	{"grant replication client on *.* to u", true},
	{"grant replication client to u", false},
	{"grant replication client, replication slave on *.* to u", true},
	{"grant replication on *.* to u", false},
	{"grant replication slave on *.* to u", true},
	{"grant replication slave to u", false},
	{"grant restricted_tables_admin on *.* to u", true},
	{"grant s3 on *.* to u1", false},
	{"grant select (a, b), insert(c) on t to u", true},
	{"grant select into s3 on *.* to u", true},
	{"grant select on *.* to u", true},
	{"grant select to u", false},
	{"grant select(a) on t to u", true},
	{"grant select(a) to u", false},
	{"grant select, grant on *.* to u", false},
	{"grant select, loop on *.* to u", true},
	{"grant select, loop to u", false},
	{"grant select, proxy on *.* to u", false},
	{"grant show databases on *.* to u", true},
	{"grant show databases to u", false},
	{"grant show on *.* to u", false},
	{"grant show view on *.* to u", true},
	{"grant show view to u", false},
	{"grant shutdown on *.* to u", true},
	{"grant shutdown to u", false},
	{"grant super on *.* to u", true},
	{"grant super to u", false},
	{"grant system_variables_admin on *.* to u", true},
	{"grant transaction loop on *.* to u", false},
	{"grant transaction, loop to u", false},
	{"grant transaction@h to u", true},
	{"grant trigger on *.* to u", true},
	{"grant trigger to u", false},
	{"grant update on *.* to u", true},
	{"grant update to u", false},
	{"grant update(a) on t to u", true},
	{"grant usage on *.* to u", true},
	{"grant usage to u", false},
	{"revoke LOOP from u", true},
	{"revoke a, all from u", false},
	{"revoke a, transaction from u", false},
	{"revoke all from u", false},
	{"revoke all on *.* from u", true},
	{"revoke all privileges from u", false},
	{"revoke all privileges on *.* from u", true},
	{"revoke all privileges, grant option from u", true},
	{"revoke all, grant option from u", true},
	{"revoke all, grant option from u, v", true},
	{"revoke loop from u", true},
	{"revoke loop on *.* from u", true},
	{"revoke loop, a from u", true},
	{"revoke loop@h from u", true},
	{"revoke proxy from u", false},
	{"revoke proxy on 'a'@'b' from u", false},
	{"revoke proxy on *.* from u", false},
	{"set default role a, loop to u", true},
	{"set default role all to u", true},
	{"set default role loop to loop", true},
	{"set default role loop to u", true},
	{"set default role none to u", true},
	{"set default role transaction to u", false},
	{"set default role transaction@h to u", true},
	{"set role LOOP", true},
	{"set role a, loop", true},
	{"set role all", true},
	{"set role all except a, loop", true},
	{"set role all except loop", true},
	{"set role all except transaction", false},
	{"set role default", true},
	{"set role loop", true},
	{"set role loop, a", true},
	{"set role none", true},
	{"set role transaction@h", true},
	{"show grants", true},
	{"show grants for u", true},
	{"show grants for u using 'a'", true},
	{"show grants for u using a", true},
	{"show grants for u using a, loop", true},
	{"show grants for u using a@h", true},
	{"show grants for u using loop", true},
	{"show grants for u using select", false},
	{"show grants for u using transaction", false},
	{"show grants for u using transaction@h", true},
}
