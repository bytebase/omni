-- DROP LIBRARY test cases for Amazon Redshift
-- Based on AWS documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_LIBRARY.html

-- Basic DROP LIBRARY statement
DROP LIBRARY my_python_library;

-- DROP LIBRARY with different identifier formats
DROP LIBRARY MyLibrary;
DROP LIBRARY my_lib_123;
DROP LIBRARY library_with_underscores;
DROP LIBRARY library123;

-- DROP LIBRARY with quoted identifiers (for case-sensitive or special character names)
DROP LIBRARY "CaseSensitiveLibrary";
DROP LIBRARY "library-with-dashes";
DROP LIBRARY "library with spaces";
DROP LIBRARY "123library";

-- DROP LIBRARY with system-generated library names (typical format)
DROP LIBRARY lib_12345;
DROP LIBRARY python_lib_v1;

-- Real-world examples based on typical Python library names
DROP LIBRARY numpy_custom;
DROP LIBRARY pandas_helper;
DROP LIBRARY sklearn_utils;
DROP LIBRARY custom_ml_lib;

-- DROP LIBRARY with longer descriptive names
DROP LIBRARY data_processing_utilities;
DROP LIBRARY machine_learning_algorithms;
DROP LIBRARY statistical_functions;

-- DROP LIBRARY with mixed case (identifiers are case-insensitive unless quoted)
DROP LIBRARY MyCustomLibrary;
DROP LIBRARY UPPERCASE_LIBRARY;
DROP LIBRARY lowercase_library;

-- Additional test cases for edge cases
DROP LIBRARY a;          -- Single character library name
DROP LIBRARY lib_;       -- Ending with underscore
DROP LIBRARY _lib;       -- Starting with underscore
DROP LIBRARY lib$123;    -- With dollar sign (valid in identifiers)

-- Testing with comments
DROP LIBRARY test_lib; -- This is a comment

/* Multi-line comment before DROP LIBRARY */
DROP LIBRARY commented_lib;

-- Note: Amazon Redshift will discontinue support for creating new Python UDFs
-- after November 1, 2025, though existing Python UDFs will continue to function.

-- Error cases that would fail (included as comments for documentation):
-- DROP LIBRARY;                    -- Missing library name
-- DROP LIBRARY IF EXISTS my_lib;   -- IF EXISTS not supported for DROP LIBRARY
-- DROP LIBRARY my_lib CASCADE;     -- CASCADE not supported for DROP LIBRARY
-- DROP LIBRARY my_lib RESTRICT;    -- RESTRICT not supported for DROP LIBRARY