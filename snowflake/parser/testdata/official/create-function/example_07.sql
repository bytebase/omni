CREATE OR REPLACE FUNCTION echo_varchar(x VARCHAR)
  RETURNS VARCHAR
  LANGUAGE SCALA
  RUNTIME_VERSION = 2.13
  HANDLER='Echo.echoVarchar'
  AS
  $$
  class Echo {
    def echoVarchar(x : String): String = {
      return x
    }
  }
  $$;
