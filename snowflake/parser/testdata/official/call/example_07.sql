DECLARE
  ret1 NUMBER;
BEGIN
  CALL sv_proc1('Manitoba', 127.4) into :ret1;
  RETURN ret1;
END;
