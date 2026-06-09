CREATE WORKLOAD GROUP IF NOT EXISTS g1 PROPERTIES (
    "max_cpu_percent"="10%",
    "max_memory_percent"="30%"
);
ALTER WORKLOAD GROUP g1 PROPERTIES (
    "max_cpu_percent"="20%",
    "max_memory_percent"="40%"
);
DROP WORKLOAD GROUP IF EXISTS g1;
SHOW WORKLOAD GROUPS;
SHOW WORKLOAD GROUPS LIKE "normal%"
