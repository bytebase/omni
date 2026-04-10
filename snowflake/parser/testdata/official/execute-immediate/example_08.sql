EXECUTE IMMEDIATE $$
DECLARE
  radius_of_circle FLOAT;
  area_of_circle FLOAT;
BEGIN
  radius_of_circle := 3;
  area_of_circle := PI() * radius_of_circle * radius_of_circle;
  RETURN area_of_circle;
END;
$$;
