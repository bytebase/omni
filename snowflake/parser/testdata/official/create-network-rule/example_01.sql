CREATE NETWORK RULE corporate_network
  TYPE = AWSVPCEID
  VALUE_LIST = ('vpce-123abc3420c1931')
  MODE = INTERNAL_STAGE
  COMMENT = 'corporate privatelink endpoint';
