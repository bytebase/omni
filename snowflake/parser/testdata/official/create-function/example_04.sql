CREATE OR REPLACE FUNCTION py_udf()
  RETURNS VARIANT
  LANGUAGE PYTHON
  RUNTIME_VERSION = '3.10'
  PACKAGES = ('numpy','pandas','xgboost==1.5.0')
  HANDLER = 'udf'
AS $$
import numpy as np
import pandas as pd
import xgboost as xgb
def udf():
    return [np.__version__, pd.__version__, xgb.__version__]
$$;
