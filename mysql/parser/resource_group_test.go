package parser

import "testing"

// MySQL 8.0's own resource group statements, and the forms it refuses.
//
// MySQL and TiDB share the three statement names CREATE, ALTER and DROP
// RESOURCE GROUP and nothing else: MySQL's options are TYPE, VCPU,
// THREAD_PRIORITY and ENABLE/DISABLE, TiDB's are RU_PER_SEC, PRIORITY,
// BURSTABLE, QUERY_LIMIT and BACKGROUND, and each engine answers 1064 for the
// other's. The TiDB-side change that moved its grammar found three MySQL
// forms this parser was taking that mysql:8.0 does not, each read from the
// server:
//
//   - TYPE is required and its value is SYSTEM or USER: `CREATE RESOURCE
//     GROUP x` and `TYPE = foo` are both 1064.
//   - VCPU, once written, needs at least one CPU number.
//   - ENABLE and DISABLE are one choice, so `ENABLE DISABLE` is 1064.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-resource-group.html
var resourceGroupAccepted = []string{
	`create resource group x type=user`,
	`create resource group x type user`,
	`create resource group x type=system`,
	`CREATE RESOURCE GROUP rg1 TYPE = USER`,
	"create resource group `default` type=user",
	`create resource group x type=user vcpu=0-3`,
	`create resource group x type=user vcpu 0-3`,
	`create resource group x type=user vcpu=0-3,5`,
	`create resource group x type=user vcpu=0`,
	`create resource group x type=system vcpu=0-3 thread_priority=0 enable`,
	`create resource group x type=user vcpu=0-3 thread_priority=-10 disable`,

	`alter resource group x`,
	`alter resource group x vcpu=0-3`,
	`alter resource group x enable`,
	`alter resource group x disable`,
	`alter resource group x force`,
	`alter resource group x enable force`,
	`alter resource group x thread_priority=5 enable`,

	`drop resource group x`,
	`drop resource group x force`,

	`set resource group x`,
	`set resource group x for 1`,
	`set resource group x for 1,2`,
	"set resource group `default`",
}

var resourceGroupRejected = []string{
	// TYPE is required and takes SYSTEM or USER.
	`create resource group x`,
	`create resource group x vcpu=0-3`,
	`create resource group x type=foo`,
	`create resource group x type=1`,
	`create resource group x type=`,
	`create resource group x vcpu=0-3 type=user`,
	`create resource group x type=user type=user`,
	// VCPU needs a CPU number once it is written.
	`create resource group x type=user vcpu=`,
	`create resource group x type=user vcpu`,
	// ENABLE or DISABLE, not both.
	`create resource group x type=user enable disable`,
	`create resource group x type=user disable enable`,
	`alter resource group x disable enable`,
	// FORCE belongs to ALTER and DROP, and comes last on ALTER.
	`create resource group x type=user force`,
	`alter resource group x force enable`,
	// ALTER takes no TYPE.
	`alter resource group x type=user`,
	// TiDB's options, which MySQL has none of.
	`create resource group x ru_per_sec=1000`,
	`create resource group x type=user ru_per_sec=1000`,
	`create resource group x burstable`,
	`create resource group x priority=low`,
	`create resource group x type=user query_limit=(exec_elapsed '1s')`,
	`create resource group x type=user background=(task_types 'br')`,
	`create resource group if not exists x type=user`,
	`drop resource group if exists x`,
	`set resource group default`,
	`calibrate resource`,
	`calibrate resource workload tpcc`,
	`query watch add sql digest 'a'`,
	`query watch remove 1`,
	`create user 'u'@'h' resource group rg`,
	`alter user 'u'@'h' resource group rg`,
}

func TestResourceGroupAccepted(t *testing.T) {
	for _, sql := range resourceGroupAccepted {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

func TestResourceGroupRejected(t *testing.T) {
	for _, sql := range resourceGroupRejected {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but mysql:8.0 rejects it", sql)
			}
		})
	}
}

// TestResourceGroupOracle holds the two tables above to a live mysql:8.0, so
// the verdicts stay a measurement rather than a reading of the manual. The
// cases create groups, so each one drops what it may have made before the
// next runs.
func TestResourceGroupOracle(t *testing.T) {
	o := startParserOracle(t)
	// A slice, not a map: the tape answers requests in the order it recorded them.
	for _, group := range []struct {
		want  bool
		cases []string
	}{
		{true, resourceGroupAccepted},
		{false, resourceGroupRejected},
	} {
		want := group.want
		for _, sql := range group.cases {
			t.Run(sql, func(t *testing.T) {
				for _, name := range []string{"x", "rg1", "`default`"} {
					_, _ = o.db.ExecContext(o.ctx, "DROP RESOURCE GROUP "+name+" FORCE")
				}
				if accepts := o.checkSyntax(t, sql); accepts != want {
					t.Errorf("mysql:8.0 accepts %q = %v, this file records %v", sql, accepts, want)
				}
			})
		}
	}
}
