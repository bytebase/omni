package parser

import (
	"testing"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

func TestRaise(t *testing.T) {
	t.Run("bare RAISE (re-raise)", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE; END`)
		if len(block.Body) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(block.Body))
		}
		r, ok := block.Body[0].(*ast.PLRaise)
		if !ok {
			t.Fatalf("expected PLRaise, got %T", block.Body[0])
		}
		if r.Level != ast.RaiseLevelNone {
			t.Errorf("expected level None, got %d", r.Level)
		}
		if r.Message != "" {
			t.Errorf("expected empty message, got %q", r.Message)
		}
		if r.CondName != "" {
			t.Errorf("expected empty condname, got %q", r.CondName)
		}
		if r.SQLState != "" {
			t.Errorf("expected empty sqlstate, got %q", r.SQLState)
		}
	})

	t.Run("RAISE with level and message", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE NOTICE 'hello %', name; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelNotice {
			t.Errorf("expected NOTICE, got %d", r.Level)
		}
		if r.Message != "hello %" {
			t.Errorf("expected 'hello %%', got %q", r.Message)
		}
		if len(r.Params) != 1 || r.Params[0] != "name" {
			t.Errorf("expected params [name], got %v", r.Params)
		}
	})

	t.Run("RAISE DEBUG", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE DEBUG 'debug info'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelDebug {
			t.Errorf("expected DEBUG, got %d", r.Level)
		}
		if r.Message != "debug info" {
			t.Errorf("expected 'debug info', got %q", r.Message)
		}
	})

	t.Run("RAISE LOG", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE LOG 'log entry'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelLog {
			t.Errorf("expected LOG, got %d", r.Level)
		}
		if r.Message != "log entry" {
			t.Errorf("expected 'log entry', got %q", r.Message)
		}
	})

	t.Run("RAISE INFO", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE INFO 'info message'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelInfo {
			t.Errorf("expected INFO, got %d", r.Level)
		}
		if r.Message != "info message" {
			t.Errorf("expected 'info message', got %q", r.Message)
		}
	})

	t.Run("RAISE WARNING", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE WARNING 'warning %', detail; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelWarning {
			t.Errorf("expected WARNING, got %d", r.Level)
		}
		if r.Message != "warning %" {
			t.Errorf("expected 'warning %%', got %q", r.Message)
		}
		if len(r.Params) != 1 || r.Params[0] != "detail" {
			t.Errorf("expected params [detail], got %v", r.Params)
		}
	})

	t.Run("RAISE EXCEPTION (default level)", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE EXCEPTION 'error occurred'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelError {
			t.Errorf("expected EXCEPTION, got %d", r.Level)
		}
		if r.Message != "error occurred" {
			t.Errorf("expected 'error occurred', got %q", r.Message)
		}
	})

	t.Run("RAISE without level (defaults to EXCEPTION)", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE 'error occurred'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelNone {
			t.Errorf("expected level None (defaults to EXCEPTION at runtime), got %d", r.Level)
		}
		if r.Message != "error occurred" {
			t.Errorf("expected 'error occurred', got %q", r.Message)
		}
	})

	t.Run("RAISE with multiple params", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE NOTICE '% and % and %', a, b, c; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelNotice {
			t.Errorf("expected NOTICE, got %d", r.Level)
		}
		if r.Message != "% and % and %" {
			t.Errorf("expected '%% and %% and %%', got %q", r.Message)
		}
		if len(r.Params) != 3 {
			t.Fatalf("expected 3 params, got %d: %v", len(r.Params), r.Params)
		}
		if r.Params[0] != "a" || r.Params[1] != "b" || r.Params[2] != "c" {
			t.Errorf("expected params [a, b, c], got %v", r.Params)
		}
	})

	t.Run("RAISE with USING", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE EXCEPTION 'fail' USING ERRCODE = 'P0001'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelError {
			t.Errorf("expected EXCEPTION, got %d", r.Level)
		}
		if r.Message != "fail" {
			t.Errorf("expected 'fail', got %q", r.Message)
		}
		if len(r.Options) != 1 {
			t.Fatalf("expected 1 option, got %d", len(r.Options))
		}
		opt := r.Options[0].(*ast.PLRaiseOption)
		if opt.OptType != "ERRCODE" {
			t.Errorf("expected ERRCODE, got %q", opt.OptType)
		}
		if opt.Expr != "'P0001'" {
			t.Errorf("expected 'P0001', got %q", opt.Expr)
		}
	})

	t.Run("RAISE with multiple USING options", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE EXCEPTION 'fail' USING ERRCODE = 'P0001', DETAIL = 'more info', HINT = 'try this'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if len(r.Options) != 3 {
			t.Fatalf("expected 3 options, got %d", len(r.Options))
		}
		opt0 := r.Options[0].(*ast.PLRaiseOption)
		opt1 := r.Options[1].(*ast.PLRaiseOption)
		opt2 := r.Options[2].(*ast.PLRaiseOption)
		if opt0.OptType != "ERRCODE" {
			t.Errorf("opt0: expected ERRCODE, got %q", opt0.OptType)
		}
		if opt0.Expr != "'P0001'" {
			t.Errorf("opt0: expected 'P0001', got %q", opt0.Expr)
		}
		if opt1.OptType != "DETAIL" {
			t.Errorf("opt1: expected DETAIL, got %q", opt1.OptType)
		}
		if opt1.Expr != "'more info'" {
			t.Errorf("opt1: expected 'more info', got %q", opt1.Expr)
		}
		if opt2.OptType != "HINT" {
			t.Errorf("opt2: expected HINT, got %q", opt2.OptType)
		}
		if opt2.Expr != "'try this'" {
			t.Errorf("opt2: expected 'try this', got %q", opt2.Expr)
		}
	})

	t.Run("RAISE with condition name", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE division_by_zero; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelNone {
			t.Errorf("expected level None, got %d", r.Level)
		}
		if r.CondName != "division_by_zero" {
			t.Errorf("expected 'division_by_zero', got %q", r.CondName)
		}
	})

	t.Run("RAISE with SQLSTATE", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE SQLSTATE '22012'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelNone {
			t.Errorf("expected level None, got %d", r.Level)
		}
		if r.SQLState != "22012" {
			t.Errorf("expected '22012', got %q", r.SQLState)
		}
	})

	t.Run("RAISE USING only", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE USING MESSAGE = 'dynamic error', ERRCODE = 'P0001'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if r.Level != ast.RaiseLevelNone {
			t.Errorf("expected level None, got %d", r.Level)
		}
		if r.Message != "" {
			t.Errorf("expected empty message, got %q", r.Message)
		}
		if len(r.Options) != 2 {
			t.Fatalf("expected 2 options, got %d", len(r.Options))
		}
		opt0 := r.Options[0].(*ast.PLRaiseOption)
		opt1 := r.Options[1].(*ast.PLRaiseOption)
		if opt0.OptType != "MESSAGE" {
			t.Errorf("opt0: expected MESSAGE, got %q", opt0.OptType)
		}
		if opt0.Expr != "'dynamic error'" {
			t.Errorf("opt0: expected 'dynamic error', got %q", opt0.Expr)
		}
		if opt1.OptType != "ERRCODE" {
			t.Errorf("opt1: expected ERRCODE, got %q", opt1.OptType)
		}
		if opt1.Expr != "'P0001'" {
			t.Errorf("opt1: expected 'P0001', got %q", opt1.Expr)
		}
	})

	t.Run("all USING options", func(t *testing.T) {
		block := parseOK(t, `BEGIN RAISE EXCEPTION 'test' USING
			MESSAGE = 'msg',
			DETAIL = 'det',
			HINT = 'hnt',
			ERRCODE = 'P0001',
			COLUMN = 'col',
			CONSTRAINT = 'cns',
			DATATYPE = 'dt',
			TABLE = 'tbl',
			SCHEMA = 'sch'; END`)
		r := block.Body[0].(*ast.PLRaise)
		if len(r.Options) != 9 {
			t.Fatalf("expected 9 options, got %d", len(r.Options))
		}
		expectedTypes := []string{"MESSAGE", "DETAIL", "HINT", "ERRCODE", "COLUMN", "CONSTRAINT", "DATATYPE", "TABLE", "SCHEMA"}
		for i, expected := range expectedTypes {
			opt := r.Options[i].(*ast.PLRaiseOption)
			if opt.OptType != expected {
				t.Errorf("option %d: expected %q, got %q", i, expected, opt.OptType)
			}
		}
	})
}
