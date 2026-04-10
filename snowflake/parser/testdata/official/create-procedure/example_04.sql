CREATE OR REPLACE PROCEDURE my_proc(fromTable STRING, toTable STRING, count INT)
  RETURNS STRING
  LANGUAGE JAVA
  RUNTIME_VERSION = '11'
  PACKAGES = ('com.snowflake:snowpark:latest')
  IMPORTS = ('@mystage/myjar.jar')
  HANDLER = 'MyClass.myMethod';
