package analysis_test

import (
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestCandidateTypeStringAndUniqueness(t *testing.T) {
	cases := []struct {
		typ  analysis.CandidateType
		want string
	}{
		{analysis.CandidateTypeNone, "none"},
		{analysis.CandidateKeyword, "keyword"},
		{analysis.CandidateCatalog, "catalog"},
		{analysis.CandidateDatabase, "database"},
		{analysis.CandidateSchema, "schema"},
		{analysis.CandidateTable, "table"},
		{analysis.CandidateView, "view"},
		{analysis.CandidateMaterializedView, "materialized_view"},
		{analysis.CandidateColumn, "column"},
		{analysis.CandidateFunction, "function"},
		{analysis.CandidateProcedure, "procedure"},
		{analysis.CandidatePackage, "package"},
		{analysis.CandidateSequence, "sequence"},
		{analysis.CandidateIndex, "index"},
		{analysis.CandidateTrigger, "trigger"},
		{analysis.CandidateEvent, "event"},
		{analysis.CandidateType_, "type"},
		{analysis.CandidateVariable, "variable"},
		{analysis.CandidateCharset, "charset"},
		{analysis.CandidateEngine, "engine"},
		{analysis.CandidatePolicy, "policy"},
		{analysis.CandidateExtension, "extension"},
		{analysis.CandidateLogin, "login"},
		{analysis.CandidateUser, "user"},
		{analysis.CandidateRole, "role"},
		{analysis.CandidateSynonym, "synonym"},
		{analysis.CandidateDatabaseLink, "database_link"},
		{analysis.CandidateConstraint, "constraint"},
		{analysis.CandidateCollection, "collection"},
		{analysis.CandidateMethod, "method"},
		{analysis.CandidateCursorMethod, "cursor_method"},
		{analysis.CandidateAggStage, "agg_stage"},
		{analysis.CandidateQueryOperator, "query_operator"},
		{analysis.CandidateBSONHelper, "bson_helper"},
		{analysis.CandidateShowTarget, "show_target"},
		{analysis.CandidateDbMethod, "db_method"},
		{analysis.CandidateRsMethod, "rs_method"},
		{analysis.CandidateShMethod, "sh_method"},
	}

	if analysis.CandidateTypeNone != 0 {
		t.Fatalf("CandidateTypeNone = %d, want 0", int(analysis.CandidateTypeNone))
	}

	seenTypes := make(map[analysis.CandidateType]bool, len(cases))
	seenNames := make(map[string]bool, len(cases))
	for _, tc := range cases {
		if seenTypes[tc.typ] {
			t.Errorf("duplicate CandidateType numeric value %d (%s)", int(tc.typ), tc.want)
		}
		seenTypes[tc.typ] = true

		got := tc.typ.String()
		if got != tc.want {
			t.Errorf("CandidateType(%d).String() = %q, want %q", int(tc.typ), got, tc.want)
		}
		if seenNames[got] {
			t.Errorf("duplicate CandidateType name %q", got)
		}
		seenNames[got] = true
	}

	for _, invalid := range []analysis.CandidateType{-1, 99, 999} {
		if got := invalid.String(); got != "unknown" {
			t.Errorf("CandidateType(%d).String() = %q, want unknown", int(invalid), got)
		}
	}

	type aliased = analysis.CandidateType
	var ct analysis.CandidateType = aliased(analysis.CandidateTable)
	if ct != analysis.CandidateTable {
		t.Fatalf("aliased CandidateType = %v, want %v", ct, analysis.CandidateTable)
	}
}
