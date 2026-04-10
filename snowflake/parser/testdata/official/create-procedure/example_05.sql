CREATE OR ALTER PROCEDURE python_add1(A NUMBER)
  RETURNS NUMBER
  LANGUAGE PYTHON
  HANDLER='main'
  RUNTIME_VERSION=3.10
  EXTERNAL_ACCESS_INTEGRATIONS=(example_integration)
  PACKAGES = ('snowflake-snowpark-python')
  EXECUTE AS OWNER
  AS
$$
def main(session, a):
    return a+1
$$;
