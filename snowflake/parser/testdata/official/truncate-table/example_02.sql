INSERT INTO temp_test_truncate SELECT seq8() FROM table(generator(rowcount=>20)) v;
