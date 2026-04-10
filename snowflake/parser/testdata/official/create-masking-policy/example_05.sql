create masking policy de_email_visibility as
 (email varchar, visibility string) returns varchar ->
   case
     when current_role() = 'ADMIN' and visibility = 'Public' then de_email(email)
     else email
   end;
