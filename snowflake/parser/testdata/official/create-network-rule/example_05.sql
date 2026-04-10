CREATE OR REPLACE NETWORK RULE ext_network_access_db.network_rules.azure_sql_private_rule
  MODE = EGRESS
  TYPE = PRIVATE_HOST_PORT
  VALUE_LIST = ('externalaccessdemo.database.windows.net');
