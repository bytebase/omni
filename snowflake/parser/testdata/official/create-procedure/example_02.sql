CREATE OR REPLACE PROCEDURE stproc1(FLOAT_PARAM1 FLOAT)
    RETURNS STRING
    LANGUAGE JAVASCRIPT
    STRICT
    EXECUTE AS OWNER
    AS
    $$
    var sql_command = 
     "INSERT INTO stproc_test_table1 (num_col1) VALUES (" + FLOAT_PARAM1 + ")";
    try {
        snowflake.execute (
            {sqlText: sql_command}
            );
        return "Succeeded.";
        }
    catch (err)  {
        return "Failed: " + err;
        }
    $$
    ;
