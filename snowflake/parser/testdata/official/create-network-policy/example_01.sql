CREATE NETWORK POLICY allow_vpceid_block_public_policy
  ALLOWED_NETWORK_RULE_LIST = ('allow_vpceid_access')
  BLOCKED_NETWORK_RULE_LIST = ('block_public_access');

DESC NETWORK POLICY rule_based_policy;
