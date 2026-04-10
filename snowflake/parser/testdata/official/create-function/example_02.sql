create function my_decrement_udf(i numeric(9, 0))
    returns numeric
    language java
    imports = ('@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar')
    handler = 'my_decrement_udf_package.my_decrement_udf_class.my_decrement_udf_method'
    ;
