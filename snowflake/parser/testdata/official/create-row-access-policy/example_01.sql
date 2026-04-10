create or replace row access policy rap_it as (empl_id varchar) returns boolean ->
  case
      when 'it_admin' = current_role() then true
      else false
  end
;
