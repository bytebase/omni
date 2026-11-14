package deparse_test

import "testing"

// DROPs of objects created on other paths, DROP SERVICE ... FORCE, and a
// service's stage and RESTORE VOLUME round-trip.
func TestDeparse_DocDropsAndServiceStages(t *testing.T) {
	for _, sql := range []string{
		"DROP SECRET IF EXISTS db.s.sec",
		"DROP PASSWORD POLICY p",
		"DROP AUTHENTICATION POLICY IF EXISTS db.s.p",
		"DROP GIT REPOSITORY r",
		"DROP EXTERNAL VOLUME IF EXISTS v",
		"DROP SEMANTIC VIEW db.s.v",
		"DROP ACCOUNT IF EXISTS a GRACE_PERIOD_IN_DAYS = 14",
		"DROP DYNAMIC TABLE IF EXISTS dt",
		"DROP SERVICE IF EXISTS s FORCE",
		"CREATE SERVICE s IN COMPUTE POOL p FROM @db.s.st/dir SPECIFICATION_FILE = 'a.yaml' MIN_INSTANCES = 2",
		"CREATE SERVICE s IN COMPUTE POOL p FROM @st SPECIFICATION_TEMPLATE_FILE = 't.yaml' USING (a => 1) COMMENT = 'c'",
		"ALTER SERVICE s FROM @st SPECIFICATION_FILE = 'a.yaml'",
		"ALTER SERVICE IF EXISTS s RESTORE VOLUME v INSTANCES 0, 2 FROM SNAPSHOT snap",
		"EXECUTE JOB SERVICE IN COMPUTE POOL p FROM @st SPECIFICATION_FILE = 'a.yaml' NAME = j",
		"EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 's' FROM @db.s.input_stage/images/ MODEL = db.s.m ASYNC = FALSE",
	} {
		assertRoundTrip(t, sql)
	}
}
