SET time_zone = "Asia/Shanghai";
SET GLOBAL exec_mem_limit = 137438953472;
SET @@exec_mem_limit = 137438953472;
SHOW VARIABLES LIKE 'max_connections';
SHOW VARIABLES LIKE '%connec%';
SHOW VARIABLES WHERE variable_name = 'version'
