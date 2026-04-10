CREATE OR REPLACE FUNCTION echo_varchar(x VARCHAR)
  RETURNS VARCHAR
  LANGUAGE SCALA
  RUNTIME_VERSION = 2.12
  IMPORTS = ('@udf_libs/echohandler.jar')
  HANDLER='Echo.echoVarchar';
