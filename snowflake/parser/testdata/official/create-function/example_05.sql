CREATE OR REPLACE FUNCTION dream(i int)
  RETURNS VARIANT
  LANGUAGE PYTHON
  RUNTIME_VERSION = '3.10'
  HANDLER = 'sleepy.snore'
  IMPORTS = ('@my_stage/sleepy.py')
