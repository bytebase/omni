package parser

import (
	"context"
	"testing"
)

// OceanBase has no RESOURCE GROUP statements.
//
// The MySQL tree this engine was forked from has CREATE, ALTER and DROP
// RESOURCE GROUP and SET RESOURCE GROUP, and this parser inherited all four.
// OceanBase 4.3.5-LTS in MySQL mode answers 1064 for every one of them. Its
// own resource model is the tenant's: CREATE RESOURCE UNIT and CREATE RESOURCE
// POOL, which are different statements with different grammars and are not
// what this parser was reading. TiDB's resource control statements are not
// OceanBase's either.
//
// Ref: https://en.oceanbase.com/docs — CREATE RESOURCE UNIT / RESOURCE POOL;
// no RESOURCE GROUP statement exists.
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
				t.Fatalf("Parse(%q) was accepted, but OceanBase has no such statement", sql)
			}
		})
	}
}

// TestResourceGroupOceanBaseOracle reads the table above back from a live
// server, so "OceanBase has no resource groups" stays a measurement rather
// than a memory of the manual. It runs on the package's own OceanBase:
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
func TestResourceGroupOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, sql := range resourceGroupRejected {
		t.Run(sql, func(t *testing.T) {
			accepts, err := engineAccepts(ctx, db, sql)
			if err != nil {
				t.Fatalf("the oracle gave no verdict on %q: %v", sql, err)
			}
			if accepts {
				t.Errorf("OceanBase accepts %q, this file records a rejection", sql)
			}
		})
	}
}
