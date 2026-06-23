package ast

import "testing"

// TestSystemTimeOutfuncs locks the outfuncs serialization of the FOR SYSTEM_TIME
// range bounds, including the per-bound TRANSACTION qualifier so AST dumps keep
// the timestamp-vs-transaction-id distinction.
func TestSystemTimeOutfuncs(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{
			name: "between transaction bounds",
			node: &SystemTime{
				Kind:            SystemTimeBetween,
				From:            &IntLit{Value: 100},
				To:              &IntLit{Value: 200},
				FromTransaction: true,
				ToTransaction:   true,
			},
			want: "{SYSTEM_TIME :loc 0 :kind 1 :from {INT_LIT :val 100 :loc 0} :from_transaction true :to {INT_LIT :val 200 :loc 0} :to_transaction true}",
		},
		{
			name: "mixed bounds (from transaction only)",
			node: &SystemTime{
				Kind:            SystemTimeBetween,
				From:            &IntLit{Value: 100},
				To:              &IntLit{Value: 200},
				FromTransaction: true,
			},
			want: "{SYSTEM_TIME :loc 0 :kind 1 :from {INT_LIT :val 100 :loc 0} :from_transaction true :to {INT_LIT :val 200 :loc 0}}",
		},
		{
			name: "timestamp bounds (no transaction flags)",
			node: &SystemTime{
				Kind: SystemTimeBetween,
				From: &IntLit{Value: 100},
				To:   &IntLit{Value: 200},
			},
			want: "{SYSTEM_TIME :loc 0 :kind 1 :from {INT_LIT :val 100 :loc 0} :to {INT_LIT :val 200 :loc 0}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NodeToString(tt.node); got != tt.want {
				t.Errorf("NodeToString mismatch\n got: %s\nwant: %s", got, tt.want)
			}
		})
	}
}
