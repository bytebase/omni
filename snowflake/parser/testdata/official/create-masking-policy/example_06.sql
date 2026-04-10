create or replace masking policy governance.policies.email_mask
as (val string) returns string ->
case
  when current_role() in ('ANALYST') then val
  when current_role() in ('SUPPORT') then regexp_replace(val,'.+\@','*****@')
  else '********'
end
comment = 'specify in row access policy'
exempt_other_policies = true
;
