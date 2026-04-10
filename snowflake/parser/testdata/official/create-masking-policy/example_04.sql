create masking policy email_visibility as
(email varchar, visibility string) returns varchar ->
  case
    when current_role() = 'ADMIN' then email
    when visibility = 'Public' then email
    else '***MASKED***'
  end;
