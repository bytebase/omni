package parser

import (
	"testing"
)

// --------------------------------------------------------------------------
// Section 1.3: Compiler Options
// --------------------------------------------------------------------------

func TestOptionDump(t *testing.T) {
	result := parseOK(t, "#option dump\nBEGIN END")
	if len(result.Options) != 1 {
		t.Fatalf("expected 1 option, got %d", len(result.Options))
	}
	if result.Options[0].Name != "dump" {
		t.Errorf("expected option name 'dump', got %q", result.Options[0].Name)
	}
}

func TestOptionPrintStrictParams(t *testing.T) {
	tests := []struct {
		name  string
		body  string
		value string
	}{
		{"on", "#print_strict_params on\nBEGIN END", "on"},
		{"off", "#print_strict_params off\nBEGIN END", "off"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Options) != 1 {
				t.Fatalf("expected 1 option, got %d", len(result.Options))
			}
			if result.Options[0].Name != "print_strict_params" {
				t.Errorf("expected option name 'print_strict_params', got %q", result.Options[0].Name)
			}
			if result.Options[0].Value != tt.value {
				t.Errorf("expected value %q, got %q", tt.value, result.Options[0].Value)
			}
		})
	}
}

func TestOptionVariableConflict(t *testing.T) {
	tests := []struct {
		name  string
		body  string
		value string
	}{
		{"error", "#variable_conflict error\nBEGIN END", "error"},
		{"use_variable", "#variable_conflict use_variable\nBEGIN END", "use_variable"},
		{"use_column", "#variable_conflict use_column\nBEGIN END", "use_column"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Options) != 1 {
				t.Fatalf("expected 1 option, got %d", len(result.Options))
			}
			if result.Options[0].Name != "variable_conflict" {
				t.Errorf("expected option name 'variable_conflict', got %q", result.Options[0].Name)
			}
			if result.Options[0].Value != tt.value {
				t.Errorf("expected value %q, got %q", tt.value, result.Options[0].Value)
			}
		})
	}
}

func TestOptionMultiple(t *testing.T) {
	body := "#option dump\n#print_strict_params on\n#variable_conflict error\nBEGIN END"
	result := parseOK(t, body)
	if len(result.Options) != 3 {
		t.Fatalf("expected 3 options, got %d", len(result.Options))
	}
	if result.Options[0].Name != "dump" {
		t.Errorf("option[0]: expected 'dump', got %q", result.Options[0].Name)
	}
	if result.Options[1].Name != "print_strict_params" {
		t.Errorf("option[1]: expected 'print_strict_params', got %q", result.Options[1].Name)
	}
	if result.Options[2].Name != "variable_conflict" {
		t.Errorf("option[2]: expected 'variable_conflict', got %q", result.Options[2].Name)
	}
}

func TestOptionUnknown(t *testing.T) {
	parseErr(t, "#option foobar\nBEGIN END", "unrecognized")
}

func TestEmptyBlock(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"empty block", "BEGIN END"},
		{"empty block with semicolon", "BEGIN END;"},
		{"declare with no variables", "DECLARE BEGIN END"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Body) != 0 {
				t.Errorf("expected empty body, got %d statements", len(result.Body))
			}
		})
	}
}

func TestLabeledBlock(t *testing.T) {
	tests := []struct {
		name  string
		body  string
		label string
	}{
		{"trailing label", "<<lbl>> BEGIN END lbl", "lbl"},
		{"labeled nested block", "BEGIN <<inner>> BEGIN END inner; END", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if result.Label != tt.label {
				t.Errorf("expected label %q, got %q", tt.label, result.Label)
			}
		})
	}
}

func TestNestedBlock(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"nested block", "BEGIN BEGIN END; END"},
		{"labeled nested block", "BEGIN <<inner>> BEGIN END inner; END"},
		{"multiple statements in block", "BEGIN BEGIN END; BEGIN END; END"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOK(t, tt.body)
			if len(result.Body) == 0 {
				t.Error("expected at least one statement in body")
			}
		})
	}
}

func TestBlockError(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantContains string
	}{
		{"label mismatch", "<<a>> BEGIN END b", "does not match"},
		{"END IF in block", "BEGIN END IF", "syntax error"},
		{"unexpected token at top level", "BEGIN END; GARBAGE", "syntax error"},
		{"missing END", "BEGIN", "expected END"},
		{"label between declare and begin", "DECLARE <<lbl>> BEGIN END;", "syntax error"},
		{"extra DECLARE within section", "DECLARE DECLARE x int; BEGIN END", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantContains == "" {
				// Should parse OK
				parseOK(t, tt.body)
				return
			}
			parseErr(t, tt.body, tt.wantContains)
		})
	}
}
