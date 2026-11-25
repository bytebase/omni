package completion

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpecialForms(t *testing.T) {
	forms := specialForms()
	for _, want := range []string{"coalesce", "nullif", "if", "ifnull", "iferror", "greatest", "substring", "trim", "current_date", "session_user"} {
		require.Contains(t, forms, want)
	}
	// Keywords that only look like calls are not functions.
	for _, not := range []string{"not", "all", "distinct", "array", "row", "left", "join"} {
		require.NotContains(t, forms, not)
	}
}

func TestTypeKeywords(t *testing.T) {
	kws := typeKeywords()
	for _, want := range []string{"INT", "INTEGER", "BIGINT", "BOOLEAN", "REAL", "STRING", "DECIMAL", "TIMESTAMPTZ"} {
		require.Contains(t, kws, want)
	}
	for _, not := range []string{"TRIGGER", "UNKNOWN", "SELECT", "ABORT"} {
		require.NotContains(t, kws, not)
	}
}

func TestFormatNames(t *testing.T) {
	in := []Candidate{
		{Text: "SUBSTRING", Type: CandidateKeyword},
		{Text: "substring", Type: CandidateFunction},
		{Text: "order", Type: CandidateFunction},
		{Text: "CalculateTotal", Type: CandidateFunction},
		{Text: "now", Type: CandidateFunction},
	}
	require.Equal(t, []Candidate{
		{Text: "substring()", Type: CandidateFunction},
		{Text: `"order"()`, Type: CandidateFunction},
		{Text: `"CalculateTotal"()`, Type: CandidateFunction},
		{Text: "now()", Type: CandidateFunction},
	}, formatNames("SELECT ", pending(in), false, true))

	// Before an existing "(", and inside double quotes, no parentheses.
	require.Equal(t, "substring", formatNames("SELECT ", pending(in[1:2]), false, false)[0].Text)
	require.Equal(t, "CalculateTotal", formatNames("SELECT ", pending(in[3:4]), true, true)[0].Text)
}

func TestFormatNamesTypes(t *testing.T) {
	in := []Candidate{
		{Text: "INT8", Type: CandidateType_},
		{Text: "INT", Type: CandidateType_},
		{Text: "OrderStatus", Type: CandidateType_},
		{Text: "order-status", Type: CandidateType_},
		{Text: "status", Type: CandidateType_},
	}
	require.Equal(t, []Candidate{
		{Text: "INT8", Type: CandidateType_},
		{Text: "INT", Type: CandidateType_},
		{Text: `"OrderStatus"`, Type: CandidateType_},
		{Text: `"order-status"`, Type: CandidateType_},
		{Text: "status", Type: CandidateType_},
	}, formatNames("SELECT 1::", pending(in), false, false))
	require.Equal(t, "OrderStatus", formatNames("SELECT 1::", pending(in[2:3]), true, false)[0].Text)
}

func pending(cands []Candidate) []pendingName {
	out := make([]pendingName, len(cands))
	for i, c := range cands {
		out[i] = pendingName{Candidate: c}
	}
	return out
}

func TestFormatNamesQualifies(t *testing.T) {
	in := []pendingName{
		{Candidate: Candidate{Text: "status", Type: CandidateType_}, schema: "analytics"},
		{Candidate: Candidate{Text: "Calc", Type: CandidateFunction}, schema: "My Schema"},
		{Candidate: Candidate{Text: "order", Type: CandidateFunction}, schema: "analytics"},
	}
	require.Equal(t, []Candidate{
		{Text: "analytics.status", Type: CandidateType_},
		{Text: `"My Schema"."Calc"()`, Type: CandidateFunction},
		{Text: `analytics."order"()`, Type: CandidateFunction},
	}, formatNames("SELECT ", in, false, true))
}
