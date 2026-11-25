CALL myprocedure();
CALL schema.myprocedure();
CALL myprocedure(1 + 2, "a", CAST(NULL AS string));
CALL myprocedure(MODEL my.model);
CALL myprocedure(CONNECTION DEFAULT);
CALL myprocedure(CONNECTION my.connection);
CALL myprocedure(TABLE my.table, (
  SELECT
    *
  FROM
    my.another_table
), mytvf(1, 2));
CALL myprocedure(@test_param_bool);
CALL myprocedure(@@sysvar);
