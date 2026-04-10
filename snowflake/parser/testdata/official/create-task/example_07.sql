CREATE OR REPLACE PROCEDURE my_unload_sp()
  returns string not null
  language javascript
  AS
    $$
      var my_sql_command = ""
      var my_sql_command = my_sql_command.concat("copy into @mystage","/",Date.now(),"/"," from mytable overwrite=true;");
      var statement1 = snowflake.createStatement( {sqlText: my_sql_command} );
      var result_set1 = statement1.execute();
    return my_sql_command;
    $$;

CREATE TASK my_copy_task
  WAREHOUSE = mywh
  SCHEDULE = '60 MINUTES'
  AS
    CALL my_unload_sp();
