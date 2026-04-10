CREATE OR REPLACE PIPE my_pipe_with_transforms
AS COPY INTO my_table (col1, col2, col3)
  FROM (
    SELECT
      $1:field1::STRING AS col1,
      $1:field2::NUMBER AS col2,
      CURRENT_TIMESTAMP() AS col3
    FROM TABLE (DATA_SOURCE(TYPE => 'STREAMING'))
  );
