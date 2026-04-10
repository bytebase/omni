CREATE OR ALTER SECURE FUNCTION multiply(a NUMBER, b NUMBER)
  RETURNS NUMBER
  COMMENT = 'Multiply two numbers.'
  AS 'a * b';
