CREATE OR REPLACE FUNCTION echo_varchar(x VARCHAR)
  RETURNS VARCHAR
  LANGUAGE JAVA
  CALLED ON NULL INPUT
  HANDLER = 'TestFunc.echoVarchar'
  TARGET_PATH = '@~/testfunc.jar'
  AS
  'class TestFunc {
    public static String echoVarchar(String x) {
      return x;
    }
  }';
