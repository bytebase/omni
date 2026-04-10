CREATE OR REPLACE EXTERNAL VOLUME ext_vol_s3_compat
  STORAGE_LOCATIONS = (
    (
      NAME = 'my_s3_compat_storage_location'
      STORAGE_PROVIDER = 'S3COMPAT'
      STORAGE_BASE_URL = 's3compat://mybucket/unload/mys3compatdata'
      CREDENTIALS = (
        AWS_KEY_ID = '1a2b3c...'
        AWS_SECRET_KEY = '4x5y6z...'
      )
      STORAGE_ENDPOINT = 'example.com'
    )
  );
