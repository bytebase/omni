-- DROP MODEL Test Cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_MODEL.html

-- Basic DROP MODEL
DROP MODEL customer_churn;

-- DROP MODEL with schema qualification
DROP MODEL demo_ml.customer_churn;

-- DROP MODEL with IF EXISTS
DROP MODEL IF EXISTS customer_churn;

-- DROP MODEL with IF EXISTS and schema qualification
DROP MODEL IF EXISTS demo_ml.customer_churn;

-- DROP MODEL with quoted identifiers
DROP MODEL "my_schema"."my_model";

-- DROP MODEL IF EXISTS with quoted identifiers
DROP MODEL IF EXISTS "my_schema"."my_model";

-- DROP MODEL with mixed case identifiers
DROP MODEL MySchema.MyModel;

-- DROP MODEL IF EXISTS with mixed case
DROP MODEL IF EXISTS MySchema.MyModel;

-- DROP MODEL with underscore in name
DROP MODEL customer_churn_v1;

-- DROP MODEL with numbers in name
DROP MODEL model_v2;

-- DROP MODEL with schema containing underscore
DROP MODEL ml_models.customer_segmentation;

-- DROP MODEL IF EXISTS with complex name
DROP MODEL IF EXISTS analytics_db.ml_models.sentiment_analysis_v3;

-- DROP MODEL with keywords as identifiers (quoted)
DROP MODEL "model"."drop";

-- DROP MODEL with fully qualified name
DROP MODEL public.my_model;

-- Multiple DROP MODEL statements
DROP MODEL model1;
DROP MODEL IF EXISTS model2;
DROP MODEL schema1.model3;

-- DROP MODEL with different casing
drop model test_model;
DROP model if exists test_model2;
Drop Model If Exists schema.model_name;