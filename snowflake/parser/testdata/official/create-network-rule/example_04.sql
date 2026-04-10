CREATE NETWORK RULE external_access_rule
  TYPE = HOST_PORT
  MODE = EGRESS
  VALUE_LIST = ('example.com', 'example.com:443');
