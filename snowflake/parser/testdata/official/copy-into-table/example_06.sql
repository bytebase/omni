CREATE OR REPLACE PIPE TEST_PRECLUSTERED_PIPE
AS
    COPY INTO TEST_PRECLUSTERED_TABLE (num) FROM (
            SELECT $1:num::number as num FROM TABLE(
                DATA_SOURCE(
                    TYPE => 'STREAMING')
            ))
      CLUSTER_AT_INGEST_TIME=TRUE;
