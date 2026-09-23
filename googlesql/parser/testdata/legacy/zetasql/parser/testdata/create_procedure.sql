CREATE PROCEDURE procedure_name()
BEGIN
END;
CREATE OR REPLACE TEMP PROCEDURE IF NOT EXISTS procedure_name()
OPTIONS
  ()
BEGIN
END;
CREATE PROCEDURE procedure_name()
OPTIONS
  (a = 1, b = "2")
BEGIN
END;
CREATE PROCEDURE procedure_name(param_a string, param_b int32, param_c numeric, param_d TABLE< int32,
int32 >, param_e ANY TYPE, OUT param_f ANY TYPE, param_g ANY TABLE, INOUT param_h ANY TABLE)
BEGIN
END;
CREATE PROCEDURE procedure_name(OUT param_a string)
BEGIN
  DECLARE a int32 ;
  SET a = 1 ;
  SET param_a = "test" ;
END;
CREATE PROCEDURE procedure_name(`OUT` int32)
BEGIN
END;
CREATE PROCEDURE procedure_name(OUT `OUT` int32)
BEGIN
END;
CREATE PROCEDURE procedure_name(param_a int32)
BEGIN
END;
CREATE PROCEDURE `OUT`(param_1 INT32)
BEGIN
END;
CREATE PROCEDURE procedure_name(param_a int32)
BEGIN
  IF param_a > 0 THEN
    RETURN ;
  END IF ;
END;
CREATE PROCEDURE procedure_name()
BEGIN
  SELECT
    1
  ;
END
;
CREATE PROCEDURE procedure_name()
BEGIN
  SELECT
    1
  ;
END
;
CREATE PROCEDURE procedure_name()
BEGIN
  DECLARE ABC INTERVAL ;
END;
CREATE PROCEDURE procedure_name(param_a int32, OUT param_b string, INOUT param_c ANY TYPE)
WITH CONNECTION DEFAULT OPTIONS
  (a = b, c = d)
LANGUAGE PYTHON AS "python code";
CREATE PROCEDURE procedure_name(param_a int32, OUT param_b string, INOUT param_c ANY TYPE)
WITH CONNECTION connection_id OPTIONS
  (a = b, c = d)
LANGUAGE PYTHON AS "python code";
CREATE PROCEDURE procedure_name(param_a int32, OUT param_b string, INOUT param_c ANY TYPE)
WITH CONNECTION connection_id OPTIONS
  ()
LANGUAGE PYTHON AS "python code";
CREATE PROCEDURE procedure_name(param_a int32, OUT param_b string, INOUT param_c ANY TYPE)
OPTIONS
  (a = b, c = d)
LANGUAGE PYTHON AS "python code";
CREATE PROCEDURE procedure_name(param_a int32, OUT param_b string, INOUT param_c ANY TYPE)
OPTIONS
  (a = b, c = d)
LANGUAGE PYTHON;
CREATE PROCEDURE procedure_name()
EXTERNAL SECURITY INVOKER LANGUAGE PYTHON;
