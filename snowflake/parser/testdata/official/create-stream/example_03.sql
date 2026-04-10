CREATE STREAM mystream ON TABLE mytable AT (TIMESTAMP => TO_TIMESTAMP_TZ('02/02/2019 01:02:03', 'mm/dd/yyyy hh24:mi:ss'));
