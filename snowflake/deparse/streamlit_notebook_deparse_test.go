package deparse_test

import "testing"

// STREAMLIT, NOTEBOOK and EXECUTE NOTEBOOK, read by their documented shape,
// round-trip.
func TestDeparse_StreamlitNotebook(t *testing.T) {
	for _, sql := range []string{
		// STREAMLIT
		"CREATE OR REPLACE STREAMLIT db.s.app FROM '@db.s.stage/src' MAIN_FILE = 'main.py' QUERY_WAREHOUSE = wh TITLE = 't' IMPORTS = ('@s/a.py', '@s/b.py') EXTERNAL_ACCESS_INTEGRATIONS = (i1, i2) SECRETS = ('k' = db.s.secret)",
		"CREATE STREAMLIT IF NOT EXISTS a ROOT_LOCATION = '@s/app' MAIN_FILE = 'main.py'",
		"CREATE STREAMLIT a FROM @db.s.repo/branches/main/app/ MAIN_FILE = 'main.py'",
		"ALTER STREAMLIT IF EXISTS a SET MAIN_FILE = 'm.py' SECRETS = ('k' = s, 'l' = t)",
		"ALTER STREAMLIT a UNSET SECRETS, TITLE",
		"ALTER STREAMLIT a RENAME TO b",
		"ALTER STREAMLIT a COMMIT",
		"ALTER STREAMLIT a ABORT",
		"ALTER STREAMLIT a PULL",
		"ALTER STREAMLIT a ADD LIVE VERSION FROM LAST",
		"ALTER STREAMLIT a PUSH",
		"ALTER STREAMLIT a PUSH TO 'https://g/r/main' USERNAME = u PASSWORD = 'p' NAME = n EMAIL = 'e@x.com' COMMENT = 'c'",
		"ALTER STREAMLIT a PUSH GIT_CREDENTIALS = db.s.secret NAME = n EMAIL = e",
		"DROP STREAMLIT IF EXISTS a",
		// NOTEBOOK, NOTEBOOK PROJECT
		"CREATE OR REPLACE NOTEBOOK db.s.nb FROM '@db.s.stage' MAIN_FILE = 'nb.ipynb' QUERY_WAREHOUSE = wh IDLE_AUTO_SHUTDOWN_TIME_SECONDS = 3600 COMPUTE_POOL = 'pool'",
		"CREATE NOTEBOOK PROJECT",
		"CREATE NOTEBOOK PROJECT IF NOT EXISTS db.s.p FROM '@db.s.stage' COMMENT = 'c'",
		"ALTER NOTEBOOK IF EXISTS a RENAME TO b",
		"ALTER NOTEBOOK a UNSET QUERY_WAREHOUSE, COMMENT",
		"ALTER NOTEBOOK a SET COMMENT = 'c' SECRETS = ('k' = db.s.secret)",
		"DROP NOTEBOOK a",
		// EXECUTE NOTEBOOK [ PROJECT ]
		"EXECUTE NOTEBOOK db.s.nb()",
		"EXECUTE NOTEBOOK nb('a', 'b c')",
		"EXECUTE NOTEBOOK nb(1, -2.5, TRUE, FALSE, NULL)",
		"ALTER STREAMLIT a UNSET TAG t, u",
		"ALTER NOTEBOOK a UNSET TAG t",
		"EXECUTE NOTEBOOK PROJECT p MAIN_FILE = 'm' COMPUTE_POOL = 'c' QUERY_WAREHOUSE = 'w' RUNTIME = 'r' ARGUMENTS = :nb_args",
		"EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'nb.ipynb' COMPUTE_POOL = 'cp' QUERY_WAREHOUSE = 'wh' RUNTIME = '2.2-CPU-PY3.11' SECRETS = (db.s.k1, db.s.k2)",
	} {
		assertRoundTrip(t, sql)
	}
}
