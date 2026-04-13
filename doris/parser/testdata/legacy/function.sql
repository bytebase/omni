CREATE FUNCTION java_udf_add_one(INT) RETURNS INT PROPERTIES (
    "file"="file:///path/to/java-udf-demo-jar-with-dependencies.jar",
    "symbol"="org.apache.doris.udf.AddOne",
    "always_nullable"="true",
    "type"="JAVA_UDF"
);
CREATE AGGREGATE FUNCTION simple_sum(INT) RETURNS INT PROPERTIES (
    "file"="file:///pathTo/java-udaf.jar",
    "symbol"="org.apache.doris.udf.demo.SimpleDemo",
    "always_nullable"="true",
    "type"="JAVA_UDF"
);
CREATE ALIAS FUNCTION id_masking(INT) WITH PARAMETER(id) AS CONCAT(LEFT(id, 3), '****', RIGHT(id, 4));
CREATE GLOBAL ALIAS FUNCTION id_masking(INT) WITH PARAMETER(id) AS CONCAT(LEFT(id, 3), '****', RIGHT(id, 4));
DROP FUNCTION my_add(INT, INT);
DROP GLOBAL FUNCTION my_add(INT, INT);
SHOW FULL FUNCTIONS IN testDb;
SHOW BUILTIN FUNCTIONS IN testDb LIKE 'year%';
SHOW GLOBAL FULL FUNCTIONS;
SHOW GLOBAL FUNCTIONS
