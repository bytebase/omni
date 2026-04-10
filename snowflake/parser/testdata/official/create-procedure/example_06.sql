CREATE OR ALTER PROCEDURE python_add1(A NUMBER)
  RETURNS NUMBER
  LANGUAGE PYTHON
  HANDLER='main'
  RUNTIME_VERSION=3.10
  EXTERNAL_ACCESS_INTEGRATIONS=(example_integration)
  secrets=('secret_variable_name'=secret_name)
  PACKAGES = ('snowflake-snowpark-python')
  EXECUTE AS CALLER
  AS
$$
def main(session, a):
    return a+1
$$;
