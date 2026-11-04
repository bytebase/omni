package parser

import "testing"

// MariaDB has no resource group statements at all.
//
// The MySQL tree this engine was forked from has CREATE, ALTER and DROP
// RESOURCE GROUP and SET RESOURCE GROUP, and this parser inherited all four.
// mariadb:11.8.8 answers 1064 for every one of them: there is no such
// statement in MariaDB, the word RESOURCE reaching no production of its own
// after CREATE, ALTER, DROP or SET. TiDB's own resource control statements are
// not MariaDB's either.
//
// Ref: https://mariadb.com/docs/server/reference/sql-statements — no resource
// group statement is listed.
var resourceGroupRejected = []string{
	// MySQL's resource group grammar.
	`create resource group x`,
	`create resource group x type=user`,
	`CREATE RESOURCE GROUP rg1 TYPE = USER`,
	`create resource group x type=system vcpu=0-3 thread_priority=0 enable`,
	`create resource group x type=user vcpu=0-3`,
	`alter resource group x`,
	`alter resource group x vcpu=0-3`,
	`alter resource group x enable`,
	`alter resource group x force`,
	`drop resource group x`,
	`drop resource group x force`,
	`set resource group x`,
	`set resource group x for 1`,
	"set resource group `default`",

	// TiDB's resource control grammar.
	`create resource group x ru_per_sec=1000`,
	`create resource group if not exists x ru_per_sec=1000`,
	`create resource group x burstable`,
	`create resource group x priority=low`,
	`create resource group x ru_per_sec=1000 query_limit=(exec_elapsed '1s')`,
	`create resource group x ru_per_sec=1000 background=(task_types 'br')`,
	`alter resource group x ru_per_sec=1000`,
	`drop resource group if exists x`,
	`set resource group default`,
	`calibrate resource`,
	`calibrate resource workload tpcc`,
	`calibrate resource start_time now()`,
	`query watch add sql digest 'a'`,
	`query watch add sql text similar to 'select 1'`,
	`query watch remove 1`,
	`create user 'u'@'h' resource group rg`,
	`alter user 'u'@'h' resource group rg`,
	`show create resource group x`,
}

func TestResourceGroupRejected(t *testing.T) {
	for _, sql := range resourceGroupRejected {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but mariadb:11.8.8 has no such statement", sql)
			}
		})
	}
}

// TestResourceGroupOracle reads the table above back from a live
// mariadb:11.8.8, so "MariaDB has no resource groups" stays a measurement
// rather than a memory of the manual.
func TestResourceGroupOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, sql := range resourceGroupRejected {
		t.Run(sql, func(t *testing.T) {
			if o.prepareAccepts(t, sql) {
				t.Errorf("mariadb:11.8.8 accepts %q, this file records a rejection", sql)
			}
		})
	}
}
