CREATE TASK triggeredTask
  WAREHOUSE = my_warehouse
  WHEN system$stream_has_data('my_stream')
  AS
    INSERT INTO my_downstream_table
    SELECT * FROM my_stream;

ALTER TASK triggeredTask RESUME;
