CREATE OR REPLACE FUNCTION get_countries_for_user ( id NUMBER )
  RETURNS TABLE (country_code CHAR, country_name VARCHAR)
  AS 'SELECT DISTINCT c.country_code, c.country_name
      FROM user_addresses a, countries c
      WHERE a.user_id = id
      AND c.country_code = a.country_code';
