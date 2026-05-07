package parser_test

import (
	"testing"

	"github.com/bytebase/omni/mysql/parser"
)

func TestIsKeyword(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"select", true},
		{"SELECT", true},
		{"execute", true},
		{"Execute", true},
		{"avg", true},
		{"count", false},
		{"cast", false},
		{"absent", false},
		{"get_master_public_key", true},
		{"MASTER_HOST", true},
		{"source_host", true},
		{"master_compression_algorithms", true},
		{"master_compression_algorithm", false},
		{"master_server_id", false},
		{"customer_id", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := parser.IsKeyword(tt.input); got != tt.want {
			t.Errorf("IsKeyword(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsReservedKeyword(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"select", true},
		{"SELECT", true},
		{"master_bind", true},
		{"master_ssl_verify_server_cert", true},
		{"master_host", false},
		{"execute", false},
		{"customer_id", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := parser.IsReservedKeyword(tt.input); got != tt.want {
			t.Errorf("IsReservedKeyword(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
