CREATE TASK t1
  USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = 'XSMALL'
  SCHEDULE = '15 SECONDS'
  AS
    DECLARE
      radius_of_circle float;
      area_of_circle float;
    BEGIN
      radius_of_circle := 3;
      area_of_circle := pi() * radius_of_circle * radius_of_circle;
      return area_of_circle;
    END;
