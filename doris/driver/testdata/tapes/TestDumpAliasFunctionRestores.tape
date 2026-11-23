-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestDumpAliasFunctionRestores default#1
> exec
| CREATE ALIAS FUNCTION alias_00(INT, INT) WITH PARAMETER(x, y) AS abs(x - y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_01(INT, INT) WITH PARAMETER(x, y) AS abs((x - y));
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_02(INT, INT) WITH PARAMETER(x, y) AS abs(((x - y)));
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_03(INT, INT) WITH PARAMETER(x, y) AS abs(x + y - 2);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_04(INT, INT) WITH PARAMETER(x, y) AS abs(x - (y - 2));
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_05(INT, INT) WITH PARAMETER(x, y) AS abs(x * 2 + y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_06(INT, INT) WITH PARAMETER(x, y) AS abs(x / y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_07(INT, INT) WITH PARAMETER(x, y) AS abs(x % y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_08(INT, INT) WITH PARAMETER(x, y) AS abs(x DIV y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_09(INT, INT) WITH PARAMETER(x, y) AS abs(x & 1 | y ^ 3);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_10(INT, INT) WITH PARAMETER(x, y) AS abs(~(x + y));
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_11(INT, INT) WITH PARAMETER(x, y) AS abs(-x + y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_12(INT, INT) WITH PARAMETER(x, y) AS abs(-(x - y));
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_13(INT, INT) WITH PARAMETER(x, y) AS abs(- -x + y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_14(INT, INT) WITH PARAMETER(x, y) AS abs(-x * y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_15(INT, INT) WITH PARAMETER(x, y) AS abs(-1 * x + y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_16(INT, INT) WITH PARAMETER(x, y) AS abs((-1) * x + y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_17(INT, INT) WITH PARAMETER(x, y) AS abs(x - -y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_18(INT, INT) WITH PARAMETER(x, y) AS abs(x * -1 + y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_19(INT, INT) WITH PARAMETER(x, y) AS abs(x + - -1 + y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_20(INT, INT) WITH PARAMETER(x, y) AS abs(CAST(x AS BIGINT) + y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_21(INT, INT) WITH PARAMETER(x, y) AS abs(abs(x - 1) - y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_22(INT, INT) WITH PARAMETER(x, y) AS abs(x + 1.50 + y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_23(INT, INT) WITH PARAMETER(x, y) AS if(x + 1 > y, x - y, -y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_24(INT, INT) WITH PARAMETER(x, y) AS if(x & 1 = 1, y, 0);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_25(INT, INT) WITH PARAMETER(x, y) AS if((x > 1) AND (y < 5 OR x = y), 1, 0);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_26(INT, INT) WITH PARAMETER(x, y) AS if(NOT x > y, 1, 0);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_27(INT, INT) WITH PARAMETER(x, y) AS if(x BETWEEN y - 1 AND y + 1, 1, 0);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_28(INT, INT) WITH PARAMETER(x, y) AS if(x IN (y - 1, y + 1), 1, 0);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_29(INT, INT) WITH PARAMETER(x, y) AS if(x IS NULL, -y, y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_30(INT, INT) WITH PARAMETER(x, y) AS coalesce(CASE x WHEN y - 1 THEN 'a' END, 'b');
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_31(INT, INT) WITH PARAMETER(x, y) AS concat('a''b', x, y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_32(INT, INT) WITH PARAMETER(x, y) AS concat('back\\slash', x, y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_33(INT, INT) WITH PARAMETER(x, y) AS concat('trailing\\', x, y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_34(INT, INT) WITH PARAMETER(x, y) AS concat('both\\''quote', x, y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_35(INT, INT) WITH PARAMETER(x, y) AS concat("dq'", x, y);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION alias_36(INT, INT) WITH PARAMETER(x, y) AS concat('tab	and -- dashes', x, y);
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_alias_function
< result
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_alias_function
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_alias_function") ORDER BY Name
< result
< column "Name" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_alias_function
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
> query
| SHOW FULL FUNCTIONS FROM `omni_dump_alias_function`
< result
< column "Signature" "CHAR"
< column "Return Type" "CHAR"
< column "Function Type" "CHAR"
< column "Intermediate Type" "CHAR"
< column "Properties" "CHAR"
< row
: bytes
| alias_00(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` - `y`))"}
< row
: bytes
| alias_01(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` - `y`)))"}
< row
: bytes
| alias_02(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` - `y`)))"}
< row
: bytes
| alias_03(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` + `y`) - 2))"}
< row
: bytes
| alias_04(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` - ((`y` - 2))))"}
< row
: bytes
| alias_05(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` * 2) + `y`))"}
< row
: bytes
| alias_06(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` / `y`))"}
< row
: bytes
| alias_07(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` % `y`))"}
< row
: bytes
| alias_08(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` DIV `y`))"}
< row
: bytes
| alias_09(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((((`x` \u0026 1) | `y`) ^ 3))"}
< row
: bytes
| alias_10(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(~ ((`x` + `y`)))"}
< row
: bytes
| alias_11(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((-1 * `x`) + `y`))"}
< row
: bytes
| alias_12(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((-1 * ((`x` - `y`))))"}
< row
: bytes
| alias_13(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((-1 * (-1 * `x`)) + `y`))"}
< row
: bytes
| alias_14(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((-1 * (`x` * `y`)))"}
< row
: bytes
| alias_15(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((-1 * (1 * `x`)) + `y`))"}
< row
: bytes
| alias_16(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((((-1) * `x`) + `y`))"}
< row
: bytes
| alias_17(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` - (-1 * `y`)))"}
< row
: bytes
| alias_18(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` * -1) + `y`))"}
< row
: bytes
| alias_19(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` + 1) + `y`))"}
< row
: bytes
| alias_20(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((CAST(`x` AS bigint) + `y`))"}
< row
: bytes
| alias_21(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((abs((`x` - 1)) - `y`))"}
< row
: bytes
| alias_22(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` + 1.50) + `y`))"}
< row
: bytes
| alias_23(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(((`x` + 1) \u003e `y`), (`x` - `y`), (-1 * `y`))"}
< row
: bytes
| alias_24(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(((`x` \u0026 1) \u003d 1), `y`, 0)"}
< row
: bytes
| alias_25(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(((`x` \u003e 1) AND ((`y` \u003c 5) OR (`x` \u003d `y`))), 1, 0)"}
< row
: bytes
| alias_26(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if((NOT (`x` \u003e `y`)), 1, 0)"}
< row
: bytes
| alias_27(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(`x` BETWEEN (`y` - 1) AND (`y` + 1), 1, 0)"}
< row
: bytes
| alias_28(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(`x` IN ((`y` - 1), (`y` + 1)), 1, 0)"}
< row
: bytes
| alias_29(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(`x` IS NULL, (-1 * `y`), `y`)"}
< row
: bytes
| alias_30(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"coalesce(CASE `x` WHEN (`y` - 1) THEN `a` END, `b`)"}
< row
: bytes
| alias_31(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`a``b`, `x`, `y`)"}
< row
: bytes
| alias_32(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`back\\slash`, `x`, `y`)"}
< row
: bytes
| alias_33(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`trailing\\`, `x`, `y`)"}
< row
: bytes
| alias_34(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`both\\``quote`, `x`, `y`)"}
< row
: bytes
| alias_35(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`dq```, `x`, `y`)"}
< row
: bytes
| alias_36(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`tab\tand -- dashes`, `x`, `y`)"}
> query
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_alias_function
< result
< column "SPECIFIC_NAME" "CHAR"
> query
| SHOW CREATE FUNCTION alias_00(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_00(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_00(int, int) WITH PARAMETER(x, y) AS abs((`x` - `y`));
> query
| SHOW CREATE FUNCTION alias_01(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_01(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_01(int, int) WITH PARAMETER(x, y) AS abs(((`x` - `y`)));
> query
| SHOW CREATE FUNCTION alias_02(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_02(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_02(int, int) WITH PARAMETER(x, y) AS abs(((`x` - `y`)));
> query
| SHOW CREATE FUNCTION alias_03(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_03(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_03(int, int) WITH PARAMETER(x, y) AS abs(((`x` + `y`) - 2));
> query
| SHOW CREATE FUNCTION alias_04(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_04(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_04(int, int) WITH PARAMETER(x, y) AS abs((`x` - ((`y` - 2))));
> query
| SHOW CREATE FUNCTION alias_05(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_05(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_05(int, int) WITH PARAMETER(x, y) AS abs(((`x` * 2) + `y`));
> query
| SHOW CREATE FUNCTION alias_06(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_06(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_06(int, int) WITH PARAMETER(x, y) AS abs((`x` / `y`));
> query
| SHOW CREATE FUNCTION alias_07(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_07(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_07(int, int) WITH PARAMETER(x, y) AS abs((`x` % `y`));
> query
| SHOW CREATE FUNCTION alias_08(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_08(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_08(int, int) WITH PARAMETER(x, y) AS abs((`x` DIV `y`));
> query
| SHOW CREATE FUNCTION alias_09(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_09(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_09(int, int) WITH PARAMETER(x, y) AS abs((((`x` & 1) | `y`) ^ 3));
> query
| SHOW CREATE FUNCTION alias_10(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_10(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_10(int, int) WITH PARAMETER(x, y) AS abs(~ ((`x` + `y`)));
> query
| SHOW CREATE FUNCTION alias_11(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_11(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_11(int, int) WITH PARAMETER(x, y) AS abs(((-1 * `x`) + `y`));
> query
| SHOW CREATE FUNCTION alias_12(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_12(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_12(int, int) WITH PARAMETER(x, y) AS abs((-1 * ((`x` - `y`))));
> query
| SHOW CREATE FUNCTION alias_13(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_13(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_13(int, int) WITH PARAMETER(x, y) AS abs(((-1 * (-1 * `x`)) + `y`));
> query
| SHOW CREATE FUNCTION alias_14(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_14(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_14(int, int) WITH PARAMETER(x, y) AS abs((-1 * (`x` * `y`)));
> query
| SHOW CREATE FUNCTION alias_15(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_15(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_15(int, int) WITH PARAMETER(x, y) AS abs(((-1 * (1 * `x`)) + `y`));
> query
| SHOW CREATE FUNCTION alias_16(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_16(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_16(int, int) WITH PARAMETER(x, y) AS abs((((-1) * `x`) + `y`));
> query
| SHOW CREATE FUNCTION alias_17(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_17(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_17(int, int) WITH PARAMETER(x, y) AS abs((`x` - (-1 * `y`)));
> query
| SHOW CREATE FUNCTION alias_18(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_18(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_18(int, int) WITH PARAMETER(x, y) AS abs(((`x` * -1) + `y`));
> query
| SHOW CREATE FUNCTION alias_19(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_19(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_19(int, int) WITH PARAMETER(x, y) AS abs(((`x` + 1) + `y`));
> query
| SHOW CREATE FUNCTION alias_20(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_20(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_20(int, int) WITH PARAMETER(x, y) AS abs((CAST(`x` AS bigint) + `y`));
> query
| SHOW CREATE FUNCTION alias_21(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_21(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_21(int, int) WITH PARAMETER(x, y) AS abs((abs((`x` - 1)) - `y`));
> query
| SHOW CREATE FUNCTION alias_22(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_22(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_22(int, int) WITH PARAMETER(x, y) AS abs(((`x` + 1.50) + `y`));
> query
| SHOW CREATE FUNCTION alias_23(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_23(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_23(int, int) WITH PARAMETER(x, y) AS if(((`x` + 1) > `y`), (`x` - `y`), (-1 * `y`));
> query
| SHOW CREATE FUNCTION alias_24(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_24(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_24(int, int) WITH PARAMETER(x, y) AS if(((`x` & 1) = 1), `y`, 0);
> query
| SHOW CREATE FUNCTION alias_25(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_25(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_25(int, int) WITH PARAMETER(x, y) AS if(((`x` > 1) AND ((`y` < 5) OR (`x` = `y`))), 1, 0);
> query
| SHOW CREATE FUNCTION alias_26(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_26(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_26(int, int) WITH PARAMETER(x, y) AS if((NOT (`x` > `y`)), 1, 0);
> query
| SHOW CREATE FUNCTION alias_27(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_27(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_27(int, int) WITH PARAMETER(x, y) AS if(`x` BETWEEN (`y` - 1) AND (`y` + 1), 1, 0);
> query
| SHOW CREATE FUNCTION alias_28(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_28(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_28(int, int) WITH PARAMETER(x, y) AS if(`x` IN ((`y` - 1), (`y` + 1)), 1, 0);
> query
| SHOW CREATE FUNCTION alias_29(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_29(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_29(int, int) WITH PARAMETER(x, y) AS if(`x` IS NULL, (-1 * `y`), `y`);
> query
| SHOW CREATE FUNCTION alias_30(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_30(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_30(int, int) WITH PARAMETER(x, y) AS coalesce(CASE `x` WHEN (`y` - 1) THEN 'a' END, 'b');
> query
| SHOW CREATE FUNCTION alias_31(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_31(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_31(int, int) WITH PARAMETER(x, y) AS concat('a''b', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_32(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_32(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_32(int, int) WITH PARAMETER(x, y) AS concat('back\slash', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_33(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_33(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_33(int, int) WITH PARAMETER(x, y) AS concat('trailing\', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_34(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_34(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_34(int, int) WITH PARAMETER(x, y) AS concat('both\''quote', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_35(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_35(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_35(int, int) WITH PARAMETER(x, y) AS concat('dq''', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_36(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_36(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_36(int, int) WITH PARAMETER(x, y) AS concat('tab	and -- dashes', `x`, `y`);
> exec
| DROP DATABASE `omni_dump_alias_function`
< ok 0 0
> exec
| CREATE DATABASE `omni_dump_alias_function`
< ok 0 0
> exec
| USE `omni_dump_alias_function`
< ok 0 0
> exec
| --
| -- Function structure for alias_00(int, int)
| --
| CREATE ALIAS FUNCTION alias_00(int, int) WITH PARAMETER(x, y) AS abs(`x` - `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_01(int, int)
| --
| CREATE ALIAS FUNCTION alias_01(int, int) WITH PARAMETER(x, y) AS abs((`x` - `y`));
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_02(int, int)
| --
| CREATE ALIAS FUNCTION alias_02(int, int) WITH PARAMETER(x, y) AS abs((`x` - `y`));
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_03(int, int)
| --
| CREATE ALIAS FUNCTION alias_03(int, int) WITH PARAMETER(x, y) AS abs(`x` + `y` - 2);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_04(int, int)
| --
| CREATE ALIAS FUNCTION alias_04(int, int) WITH PARAMETER(x, y) AS abs(`x` - (`y` - 2));
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_05(int, int)
| --
| CREATE ALIAS FUNCTION alias_05(int, int) WITH PARAMETER(x, y) AS abs(`x` * 2 + `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_06(int, int)
| --
| CREATE ALIAS FUNCTION alias_06(int, int) WITH PARAMETER(x, y) AS abs(`x` / `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_07(int, int)
| --
| CREATE ALIAS FUNCTION alias_07(int, int) WITH PARAMETER(x, y) AS abs(`x` % `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_08(int, int)
| --
| CREATE ALIAS FUNCTION alias_08(int, int) WITH PARAMETER(x, y) AS abs(`x` DIV `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_09(int, int)
| --
| CREATE ALIAS FUNCTION alias_09(int, int) WITH PARAMETER(x, y) AS abs(`x` & 1 | `y` ^ 3);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_10(int, int)
| --
| CREATE ALIAS FUNCTION alias_10(int, int) WITH PARAMETER(x, y) AS abs(~ (`x` + `y`));
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_11(int, int)
| --
| CREATE ALIAS FUNCTION alias_11(int, int) WITH PARAMETER(x, y) AS abs(-`x` + `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_12(int, int)
| --
| CREATE ALIAS FUNCTION alias_12(int, int) WITH PARAMETER(x, y) AS abs(-(`x` - `y`));
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_13(int, int)
| --
| CREATE ALIAS FUNCTION alias_13(int, int) WITH PARAMETER(x, y) AS abs(- -`x` + `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_14(int, int)
| --
| CREATE ALIAS FUNCTION alias_14(int, int) WITH PARAMETER(x, y) AS abs(-`x` * `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_15(int, int)
| --
| CREATE ALIAS FUNCTION alias_15(int, int) WITH PARAMETER(x, y) AS abs(-1 * `x` + `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_16(int, int)
| --
| CREATE ALIAS FUNCTION alias_16(int, int) WITH PARAMETER(x, y) AS abs((-1) * `x` + `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_17(int, int)
| --
| CREATE ALIAS FUNCTION alias_17(int, int) WITH PARAMETER(x, y) AS abs(`x` - -`y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_18(int, int)
| --
| CREATE ALIAS FUNCTION alias_18(int, int) WITH PARAMETER(x, y) AS abs(`x` * -1 + `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_19(int, int)
| --
| CREATE ALIAS FUNCTION alias_19(int, int) WITH PARAMETER(x, y) AS abs(`x` + 1 + `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_20(int, int)
| --
| CREATE ALIAS FUNCTION alias_20(int, int) WITH PARAMETER(x, y) AS abs(CAST(`x` AS bigint) + `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_21(int, int)
| --
| CREATE ALIAS FUNCTION alias_21(int, int) WITH PARAMETER(x, y) AS abs(abs(`x` - 1) - `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_22(int, int)
| --
| CREATE ALIAS FUNCTION alias_22(int, int) WITH PARAMETER(x, y) AS abs(`x` + 1.50 + `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_23(int, int)
| --
| CREATE ALIAS FUNCTION alias_23(int, int) WITH PARAMETER(x, y) AS if((`x` + 1 > `y`), `x` - `y`, -`y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_24(int, int)
| --
| CREATE ALIAS FUNCTION alias_24(int, int) WITH PARAMETER(x, y) AS if((`x` & 1 = 1), `y`, 0);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_25(int, int)
| --
| CREATE ALIAS FUNCTION alias_25(int, int) WITH PARAMETER(x, y) AS if(((`x` > 1) AND ((`y` < 5) OR (`x` = `y`))), 1, 0);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_26(int, int)
| --
| CREATE ALIAS FUNCTION alias_26(int, int) WITH PARAMETER(x, y) AS if((NOT (`x` > `y`)), 1, 0);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_27(int, int)
| --
| CREATE ALIAS FUNCTION alias_27(int, int) WITH PARAMETER(x, y) AS if(`x` BETWEEN `y` - 1 AND `y` + 1, 1, 0);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_28(int, int)
| --
| CREATE ALIAS FUNCTION alias_28(int, int) WITH PARAMETER(x, y) AS if(`x` IN (`y` - 1, `y` + 1), 1, 0);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_29(int, int)
| --
| CREATE ALIAS FUNCTION alias_29(int, int) WITH PARAMETER(x, y) AS if(`x` IS NULL, -`y`, `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_30(int, int)
| --
| CREATE ALIAS FUNCTION alias_30(int, int) WITH PARAMETER(x, y) AS coalesce(CASE `x` WHEN `y` - 1 THEN 'a' END, 'b');
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_31(int, int)
| --
| CREATE ALIAS FUNCTION alias_31(int, int) WITH PARAMETER(x, y) AS concat('a''b', `x`, `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_32(int, int)
| --
| CREATE ALIAS FUNCTION alias_32(int, int) WITH PARAMETER(x, y) AS concat('back\\slash', `x`, `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_33(int, int)
| --
| CREATE ALIAS FUNCTION alias_33(int, int) WITH PARAMETER(x, y) AS concat('trailing\\', `x`, `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_34(int, int)
| --
| CREATE ALIAS FUNCTION alias_34(int, int) WITH PARAMETER(x, y) AS concat('both\\''quote', `x`, `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_35(int, int)
| --
| CREATE ALIAS FUNCTION alias_35(int, int) WITH PARAMETER(x, y) AS concat('dq''', `x`, `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for alias_36(int, int)
| --
| CREATE ALIAS FUNCTION alias_36(int, int) WITH PARAMETER(x, y) AS concat('tab	and -- dashes', `x`, `y`);
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_alias_function
< result
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_alias_function
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_alias_function") ORDER BY Name
< result
< column "Name" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_alias_function
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
> query
| SHOW FULL FUNCTIONS FROM `omni_dump_alias_function`
< result
< column "Signature" "CHAR"
< column "Return Type" "CHAR"
< column "Function Type" "CHAR"
< column "Intermediate Type" "CHAR"
< column "Properties" "CHAR"
< row
: bytes
| alias_00(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` - `y`))"}
< row
: bytes
| alias_01(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` - `y`)))"}
< row
: bytes
| alias_02(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` - `y`)))"}
< row
: bytes
| alias_03(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` + `y`) - 2))"}
< row
: bytes
| alias_04(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` - ((`y` - 2))))"}
< row
: bytes
| alias_05(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` * 2) + `y`))"}
< row
: bytes
| alias_06(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` / `y`))"}
< row
: bytes
| alias_07(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` % `y`))"}
< row
: bytes
| alias_08(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` DIV `y`))"}
< row
: bytes
| alias_09(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((((`x` \u0026 1) | `y`) ^ 3))"}
< row
: bytes
| alias_10(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(~ ((`x` + `y`)))"}
< row
: bytes
| alias_11(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((-1 * `x`) + `y`))"}
< row
: bytes
| alias_12(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((-1 * ((`x` - `y`))))"}
< row
: bytes
| alias_13(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((-1 * (-1 * `x`)) + `y`))"}
< row
: bytes
| alias_14(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((-1 * (`x` * `y`)))"}
< row
: bytes
| alias_15(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((-1 * (1 * `x`)) + `y`))"}
< row
: bytes
| alias_16(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((((-1) * `x`) + `y`))"}
< row
: bytes
| alias_17(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((`x` - (-1 * `y`)))"}
< row
: bytes
| alias_18(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` * -1) + `y`))"}
< row
: bytes
| alias_19(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` + 1) + `y`))"}
< row
: bytes
| alias_20(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((CAST(`x` AS bigint) + `y`))"}
< row
: bytes
| alias_21(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs((abs((`x` - 1)) - `y`))"}
< row
: bytes
| alias_22(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"abs(((`x` + 1.50) + `y`))"}
< row
: bytes
| alias_23(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(((`x` + 1) \u003e `y`), (`x` - `y`), (-1 * `y`))"}
< row
: bytes
| alias_24(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(((`x` \u0026 1) \u003d 1), `y`, 0)"}
< row
: bytes
| alias_25(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(((`x` \u003e 1) AND ((`y` \u003c 5) OR (`x` \u003d `y`))), 1, 0)"}
< row
: bytes
| alias_26(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if((NOT (`x` \u003e `y`)), 1, 0)"}
< row
: bytes
| alias_27(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(`x` BETWEEN (`y` - 1) AND (`y` + 1), 1, 0)"}
< row
: bytes
| alias_28(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(`x` IN ((`y` - 1), (`y` + 1)), 1, 0)"}
< row
: bytes
| alias_29(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"if(`x` IS NULL, (-1 * `y`), `y`)"}
< row
: bytes
| alias_30(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"coalesce(CASE `x` WHEN (`y` - 1) THEN `a` END, `b`)"}
< row
: bytes
| alias_31(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`a``b`, `x`, `y`)"}
< row
: bytes
| alias_32(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`back\\slash`, `x`, `y`)"}
< row
: bytes
| alias_33(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`trailing\\`, `x`, `y`)"}
< row
: bytes
| alias_34(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`both\\``quote`, `x`, `y`)"}
< row
: bytes
| alias_35(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`dq```, `x`, `y`)"}
< row
: bytes
| alias_36(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"concat(`tab\tand -- dashes`, `x`, `y`)"}
> query
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_alias_function
< result
< column "SPECIFIC_NAME" "CHAR"
> query
| SHOW CREATE FUNCTION alias_00(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_00(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_00(int, int) WITH PARAMETER(x, y) AS abs((`x` - `y`));
> query
| SHOW CREATE FUNCTION alias_01(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_01(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_01(int, int) WITH PARAMETER(x, y) AS abs(((`x` - `y`)));
> query
| SHOW CREATE FUNCTION alias_02(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_02(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_02(int, int) WITH PARAMETER(x, y) AS abs(((`x` - `y`)));
> query
| SHOW CREATE FUNCTION alias_03(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_03(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_03(int, int) WITH PARAMETER(x, y) AS abs(((`x` + `y`) - 2));
> query
| SHOW CREATE FUNCTION alias_04(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_04(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_04(int, int) WITH PARAMETER(x, y) AS abs((`x` - ((`y` - 2))));
> query
| SHOW CREATE FUNCTION alias_05(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_05(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_05(int, int) WITH PARAMETER(x, y) AS abs(((`x` * 2) + `y`));
> query
| SHOW CREATE FUNCTION alias_06(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_06(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_06(int, int) WITH PARAMETER(x, y) AS abs((`x` / `y`));
> query
| SHOW CREATE FUNCTION alias_07(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_07(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_07(int, int) WITH PARAMETER(x, y) AS abs((`x` % `y`));
> query
| SHOW CREATE FUNCTION alias_08(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_08(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_08(int, int) WITH PARAMETER(x, y) AS abs((`x` DIV `y`));
> query
| SHOW CREATE FUNCTION alias_09(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_09(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_09(int, int) WITH PARAMETER(x, y) AS abs((((`x` & 1) | `y`) ^ 3));
> query
| SHOW CREATE FUNCTION alias_10(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_10(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_10(int, int) WITH PARAMETER(x, y) AS abs(~ ((`x` + `y`)));
> query
| SHOW CREATE FUNCTION alias_11(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_11(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_11(int, int) WITH PARAMETER(x, y) AS abs(((-1 * `x`) + `y`));
> query
| SHOW CREATE FUNCTION alias_12(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_12(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_12(int, int) WITH PARAMETER(x, y) AS abs((-1 * ((`x` - `y`))));
> query
| SHOW CREATE FUNCTION alias_13(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_13(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_13(int, int) WITH PARAMETER(x, y) AS abs(((-1 * (-1 * `x`)) + `y`));
> query
| SHOW CREATE FUNCTION alias_14(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_14(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_14(int, int) WITH PARAMETER(x, y) AS abs((-1 * (`x` * `y`)));
> query
| SHOW CREATE FUNCTION alias_15(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_15(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_15(int, int) WITH PARAMETER(x, y) AS abs(((-1 * (1 * `x`)) + `y`));
> query
| SHOW CREATE FUNCTION alias_16(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_16(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_16(int, int) WITH PARAMETER(x, y) AS abs((((-1) * `x`) + `y`));
> query
| SHOW CREATE FUNCTION alias_17(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_17(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_17(int, int) WITH PARAMETER(x, y) AS abs((`x` - (-1 * `y`)));
> query
| SHOW CREATE FUNCTION alias_18(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_18(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_18(int, int) WITH PARAMETER(x, y) AS abs(((`x` * -1) + `y`));
> query
| SHOW CREATE FUNCTION alias_19(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_19(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_19(int, int) WITH PARAMETER(x, y) AS abs(((`x` + 1) + `y`));
> query
| SHOW CREATE FUNCTION alias_20(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_20(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_20(int, int) WITH PARAMETER(x, y) AS abs((CAST(`x` AS bigint) + `y`));
> query
| SHOW CREATE FUNCTION alias_21(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_21(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_21(int, int) WITH PARAMETER(x, y) AS abs((abs((`x` - 1)) - `y`));
> query
| SHOW CREATE FUNCTION alias_22(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_22(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_22(int, int) WITH PARAMETER(x, y) AS abs(((`x` + 1.50) + `y`));
> query
| SHOW CREATE FUNCTION alias_23(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_23(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_23(int, int) WITH PARAMETER(x, y) AS if(((`x` + 1) > `y`), (`x` - `y`), (-1 * `y`));
> query
| SHOW CREATE FUNCTION alias_24(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_24(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_24(int, int) WITH PARAMETER(x, y) AS if(((`x` & 1) = 1), `y`, 0);
> query
| SHOW CREATE FUNCTION alias_25(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_25(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_25(int, int) WITH PARAMETER(x, y) AS if(((`x` > 1) AND ((`y` < 5) OR (`x` = `y`))), 1, 0);
> query
| SHOW CREATE FUNCTION alias_26(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_26(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_26(int, int) WITH PARAMETER(x, y) AS if((NOT (`x` > `y`)), 1, 0);
> query
| SHOW CREATE FUNCTION alias_27(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_27(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_27(int, int) WITH PARAMETER(x, y) AS if(`x` BETWEEN (`y` - 1) AND (`y` + 1), 1, 0);
> query
| SHOW CREATE FUNCTION alias_28(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_28(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_28(int, int) WITH PARAMETER(x, y) AS if(`x` IN ((`y` - 1), (`y` + 1)), 1, 0);
> query
| SHOW CREATE FUNCTION alias_29(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_29(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_29(int, int) WITH PARAMETER(x, y) AS if(`x` IS NULL, (-1 * `y`), `y`);
> query
| SHOW CREATE FUNCTION alias_30(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_30(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_30(int, int) WITH PARAMETER(x, y) AS coalesce(CASE `x` WHEN (`y` - 1) THEN 'a' END, 'b');
> query
| SHOW CREATE FUNCTION alias_31(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_31(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_31(int, int) WITH PARAMETER(x, y) AS concat('a''b', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_32(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_32(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_32(int, int) WITH PARAMETER(x, y) AS concat('back\slash', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_33(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_33(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_33(int, int) WITH PARAMETER(x, y) AS concat('trailing\', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_34(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_34(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_34(int, int) WITH PARAMETER(x, y) AS concat('both\''quote', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_35(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_35(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_35(int, int) WITH PARAMETER(x, y) AS concat('dq''', `x`, `y`);
> query
| SHOW CREATE FUNCTION alias_36(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| alias_36(int, int)
: bytes
| CREATE ALIAS FUNCTION alias_36(int, int) WITH PARAMETER(x, y) AS concat('tab	and -- dashes', `x`, `y`);
