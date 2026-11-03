package review

import (
	"regexp"
	"testing"
)

// rulePrefixes is the two-letter prefix of each rule's codes.
var rulePrefixes = map[Rule]string{
	Syntax: "SY", WalkThrough: "WT", OnlineMigration: "OM", PriorBackup: "PB",
	RequireIsNull: "NU", RequireWhere: "WH", DisallowDropObject: "DO", DisallowTruncate: "TR",
	DisallowDropConstraint: "DC", DisallowRename: "RN", RequirePrimaryKey: "PK",
}

func TestCodeTable(t *testing.T) {
	codeShape := regexp.MustCompile(`^[A-Z]{2}[1-9][0-9]{2}$`)
	nameShape := regexp.MustCompile(`^[a-z]+(-[a-z]+)*$`)
	seenCode, seenName := map[FindingCode]bool{}, map[string]bool{}
	rules := map[Rule]bool{}
	var last FindingCode
	for _, info := range Codes() {
		if !codeShape.MatchString(string(info.Code)) {
			t.Errorf("%s: not two letters and three digits", info.Code)
		}
		if prefix := rulePrefixes[info.Rule]; prefix == "" || string(info.Code[:2]) != prefix {
			t.Errorf("%s: prefix does not match rule %s", info.Code, info.Rule)
		}
		if seenCode[info.Code] {
			t.Errorf("%s: listed twice", info.Code)
		}
		if !nameShape.MatchString(info.Name) || seenName[info.Name] {
			t.Errorf("%s: name %q is not unique kebab case", info.Code, info.Name)
		}
		if info.Title == "" {
			t.Errorf("%s: no title", info.Code)
		}
		if last != "" && info.Code[:2] == last[:2] && info.Code <= last {
			t.Errorf("%s: out of order after %s", info.Code, last)
		}
		if got, ok := info.Code.Info(); !ok || got != info {
			t.Errorf("%s: Info() = %+v, %v", info.Code, got, ok)
		}
		seenCode[info.Code], seenName[info.Name], rules[info.Rule] = true, true, true
		last = info.Code
	}
	for rule := range rulePrefixes {
		if !rules[rule] {
			t.Errorf("rule %s has no code", rule)
		}
	}
	if FindingCode("XX999").Rule() != "" {
		t.Error("an unknown code has a rule")
	}
}

func TestDropCode(t *testing.T) {
	tests := map[string]FindingCode{
		"table": DropsData, "temporary table": DropsData, "foreign table": DropsData, "external table": DropsData,
		"column": DropsData, "columns": DropsData, "partition": DropsData, "subpartition": DropsData,
		"schema": DropsData, "database": DropsData, "pluggable database": DropsData, "sequence": DropsData,
		"system versioning": DropsData, "objects owned by": DropsData,
		"view": DropsDefinition, "materialized view": DropsDefinition, "index": DropsDefinition,
		"function": DropsDefinition, "trigger": DropsDefinition, "type": DropsDefinition, "extension": DropsDefinition,
		"tablespace": DropsDefinition, "placement policy": DropsDefinition, "cast": DropsDefinition,
		"user": DropsAccess, "role": DropsAccess, "server role": DropsAccess, "login": DropsAccess,
		"policy": DropsAccess, "row access policy": DropsAccess, "masking policy": DropsAccess,
		"master key": DropsAccess, "symmetric key": DropsAccess, "database scoped credential": DropsAccess,
		"server audit": DropsAccess, "database audit specification": DropsAccess, "user mapping": DropsAccess,
		"lockdown profile": DropsAccess, "certificate": DropsAccess,
	}
	for kind, want := range tests {
		if got := DropCode(kind); got != want {
			t.Errorf("DropCode(%q) = %s, want %s", kind, got, want)
		}
	}
}

func TestCheckCodes(t *testing.T) {
	ok := Finding{Rule: WalkThrough, Codes: []FindingCode{ObjectNotFound}}
	if err := CheckCodes(ok); err != nil {
		t.Errorf("CheckCodes(%+v) = %v", ok, err)
	}
	for _, bad := range []Finding{
		{Rule: WalkThrough},
		{Rule: WalkThrough, Codes: []FindingCode{DropsData}},
		{Rule: WalkThrough, Codes: []FindingCode{"WT999"}},
		{Rule: OnlineMigration, Codes: []FindingCode{UnsupportedChange, UnsupportedChange}},
	} {
		if CheckCodes(bad) == nil {
			t.Errorf("CheckCodes(%+v) = nil, want an error", bad)
		}
	}
}
