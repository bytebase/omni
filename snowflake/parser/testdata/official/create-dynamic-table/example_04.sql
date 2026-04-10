CREATE DYNAMIC TABLE my_cloned_dynamic_table CLONE my_dynamic_table AT (TIMESTAMP => TO_TIMESTAMP_TZ('04/05/2013 01:02:03', 'mm/dd/yyyy hh24:mi:ss'));
