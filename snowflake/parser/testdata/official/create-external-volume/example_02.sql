CREATE EXTERNAL VOLUME exvol
  STORAGE_LOCATIONS =
    (
      (
        NAME = 'my-us-east-1'
        STORAGE_PROVIDER = 'GCS'
        STORAGE_BASE_URL = 'gcs://mybucket1/path1/'
        ENCRYPTION=(TYPE='GCS_SSE_KMS' KMS_KEY_ID = '1234abcd-12ab-34cd-56ef-1234567890ab')
      )
    )
  ALLOW_WRITES = TRUE;
