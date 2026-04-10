CREATE OR REPLACE FUNCTION js_factorial(d double)
  RETURNS double
  LANGUAGE JAVASCRIPT
  STRICT
  AS '
  if (D <= 0) {
    return 1;
  } else {
    var result = 1;
    for (var i = 2; i <= D; i++) {
      result = result * i;
    }
    return result;
  }
  ';
