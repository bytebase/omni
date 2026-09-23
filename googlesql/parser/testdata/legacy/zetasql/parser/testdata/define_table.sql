DEFINE TABLE t1(a = 1, b = "a", c = 1.4, d = true);
DEFINE TABLE a.b.`c 2`(x = '''
foo''', y = "2011-10-22", z = @param, zz = @@sysvar);
DEFINE TABLE t1(a = b);
DEFINE TABLE t1();
DEFINE TABLE t1(a = b + 1);
DEFINE TABLE t1(x = y.z);
DEFINE TABLE x.`all`(a = 1);
DEFINE TABLE t1(a = CAST(1 AS INT32));
DEFINE TABLE t1(`PROTO` = "PROTO<foo>", `hash` = "HASH");
