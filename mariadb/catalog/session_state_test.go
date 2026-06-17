package catalog

import "testing"

func TestCatalogSessionStateExplicitDefaultsForTimestamp(t *testing.T) {
	c := New()
	if !c.session.ExplicitDefaultsForTimestamp {
		t.Fatal("default explicit_defaults_for_timestamp should be true")
	}

	mustExec(t, c, "SET SESSION explicit_defaults_for_timestamp=0")
	if c.session.ExplicitDefaultsForTimestamp {
		t.Fatal("expected SESSION explicit_defaults_for_timestamp=0 to disable setting")
	}

	mustExec(t, c, "SET explicit_defaults_for_timestamp=ON")
	if !c.session.ExplicitDefaultsForTimestamp {
		t.Fatal("expected explicit_defaults_for_timestamp=ON to enable setting")
	}

	mustExec(t, c, "SET @@session.explicit_defaults_for_timestamp=OFF")
	if c.session.ExplicitDefaultsForTimestamp {
		t.Fatal("expected @@session.explicit_defaults_for_timestamp=OFF to disable setting")
	}

	mustExec(t, c, "SET @@session.explicit_defaults_for_timestamp=DEFAULT")
	if !c.session.ExplicitDefaultsForTimestamp {
		t.Fatal("expected DEFAULT to restore explicit_defaults_for_timestamp=true")
	}

	mustExec(t, c, "SET explicit_defaults_for_timestamp=1")
	if !c.session.ExplicitDefaultsForTimestamp {
		t.Fatal("expected explicit_defaults_for_timestamp=1 to enable setting")
	}
}

func TestCatalogSessionStateSQLModeIsRecorded(t *testing.T) {
	c := New()
	mustExec(t, c, "SET SESSION sql_mode='STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'")
	if got, want := c.session.SQLMode, "STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION"; got != want {
		t.Fatalf("SQLMode = %q, want %q", got, want)
	}

	mustExec(t, c, "SET SESSION sql_mode=DEFAULT")
	if got, want := c.session.SQLMode, "DEFAULT"; got != want {
		t.Fatalf("SQLMode = %q, want %q", got, want)
	}
}
