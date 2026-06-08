-- CREATE EXTERNAL MODEL test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_create_external_model.html

-- Basic CREATE EXTERNAL MODEL with minimal required settings
CREATE EXTERNAL MODEL customer_review_model
FUNCTION customer_review_llm
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0'
);

-- CREATE EXTERNAL MODEL with DEFAULT IAM_ROLE
CREATE EXTERNAL MODEL sentiment_model
FUNCTION sentiment_analysis_llm
IAM_ROLE DEFAULT
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1'
);

-- CREATE EXTERNAL MODEL with PROMPT setting
CREATE EXTERNAL MODEL product_description_model
FUNCTION generate_product_description
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-haiku-20240307-v1:0',
    PROMPT 'Generate a compelling product description for the following item:'
);

-- CREATE EXTERNAL MODEL with SUFFIX setting
CREATE EXTERNAL MODEL marketing_copy_model
FUNCTION create_marketing_copy
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0',
    SUFFIX 'Please ensure the response is under 500 characters and includes a call to action.'
);

-- CREATE EXTERNAL MODEL with REQUEST_TYPE RAW
CREATE EXTERNAL MODEL raw_text_model
FUNCTION process_raw_text
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1',
    REQUEST_TYPE RAW
);

-- CREATE EXTERNAL MODEL with REQUEST_TYPE UNIFIED
CREATE EXTERNAL MODEL unified_text_model
FUNCTION process_unified_text
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0',
    REQUEST_TYPE UNIFIED
);

-- CREATE EXTERNAL MODEL with RESPONSE_TYPE VARCHAR
CREATE EXTERNAL MODEL varchar_response_model
FUNCTION generate_varchar_response
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-haiku-20240307-v1:0',
    RESPONSE_TYPE VARCHAR
);

-- CREATE EXTERNAL MODEL with RESPONSE_TYPE SUPER
CREATE EXTERNAL MODEL super_response_model
FUNCTION generate_super_response
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1',
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL with all optional settings
CREATE EXTERNAL MODEL comprehensive_text_model
FUNCTION comprehensive_text_processor
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0',
    PROMPT 'You are a helpful AI assistant analyzing customer feedback. Please provide insights on the following:',
    SUFFIX 'Format your response as structured JSON with sentiment, key_themes, and recommendations fields.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for customer service automation
CREATE EXTERNAL MODEL customer_service_model
FUNCTION automate_customer_response
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-haiku-20240307-v1:0',
    PROMPT 'You are a professional customer service representative. Respond to the following customer inquiry with empathy and provide helpful solutions:',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE VARCHAR
);

-- CREATE EXTERNAL MODEL for data summarization
CREATE EXTERNAL MODEL data_summary_model
FUNCTION summarize_business_data
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1',
    PROMPT 'Analyze the following business data and provide a concise executive summary:',
    SUFFIX 'Include key metrics, trends, and actionable insights in your summary.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE VARCHAR
);

-- CREATE EXTERNAL MODEL for content generation
CREATE EXTERNAL MODEL content_generator_model
FUNCTION generate_marketing_content
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0',
    PROMPT 'Create engaging marketing content based on the provided product information:',
    SUFFIX 'Ensure the content is optimized for social media engagement and includes relevant hashtags.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE VARCHAR
);

-- CREATE EXTERNAL MODEL for code review assistance
CREATE EXTERNAL MODEL code_review_model
FUNCTION review_code_snippet
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-haiku-20240307-v1:0',
    PROMPT 'Review the following code snippet and provide feedback on best practices, potential issues, and improvements:',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for translation services
CREATE EXTERNAL MODEL translation_model
FUNCTION translate_text
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1',
    PROMPT 'Translate the following text accurately while preserving the original meaning and tone:',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE VARCHAR
);

-- CREATE EXTERNAL MODEL for financial analysis
CREATE EXTERNAL MODEL financial_analysis_model
FUNCTION analyze_financial_data
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0',
    PROMPT 'Analyze the following financial data and provide professional insights:',
    SUFFIX 'Include risk assessment, performance indicators, and strategic recommendations.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for email automation
CREATE EXTERNAL MODEL email_automation_model
FUNCTION generate_automated_email
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-haiku-20240307-v1:0',
    PROMPT 'Generate a professional email based on the provided context and recipient information:',
    SUFFIX 'Ensure appropriate tone, clear call-to-action, and professional formatting.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE VARCHAR
);

-- CREATE EXTERNAL MODEL for legal document analysis
CREATE EXTERNAL MODEL legal_analysis_model
FUNCTION analyze_legal_document
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1',
    PROMPT 'Analyze the following legal document and extract key information:',
    SUFFIX 'Highlight important clauses, potential risks, and compliance requirements.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for research assistance
CREATE EXTERNAL MODEL research_assistant_model
FUNCTION assist_with_research
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0',
    PROMPT 'Provide comprehensive research assistance on the following topic:',
    SUFFIX 'Include relevant sources, key findings, and actionable recommendations.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for healthcare data analysis
CREATE EXTERNAL MODEL healthcare_analysis_model
FUNCTION analyze_healthcare_data
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-haiku-20240307-v1:0',
    PROMPT 'Analyze the following healthcare data while maintaining patient privacy:',
    SUFFIX 'Provide insights on trends, patterns, and potential interventions.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for educational content creation
CREATE EXTERNAL MODEL educational_content_model
FUNCTION create_educational_content
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1',
    PROMPT 'Create engaging educational content based on the provided learning objectives:',
    SUFFIX 'Ensure content is age-appropriate, engaging, and includes interactive elements.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE VARCHAR
);

-- CREATE EXTERNAL MODEL with Amazon Titan model
CREATE EXTERNAL MODEL titan_text_model
FUNCTION process_with_titan
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'amazon.titan-text-express-v1'
);

-- CREATE EXTERNAL MODEL with Cohere model
CREATE EXTERNAL MODEL cohere_command_model
FUNCTION process_with_cohere
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'cohere.command-text-v14'
);

-- CREATE EXTERNAL MODEL with AI21 Labs model
CREATE EXTERNAL MODEL ai21_jurassic_model
FUNCTION process_with_ai21
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'ai21.j2-ultra-v1'
);

-- CREATE EXTERNAL MODEL for complex prompt engineering
CREATE EXTERNAL MODEL advanced_prompt_model
FUNCTION advanced_text_processing
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0',
    PROMPT 'You are an expert data analyst. Given the following structured data, please perform a comprehensive analysis following these steps: 1) Data quality assessment 2) Statistical summary 3) Pattern identification 4) Anomaly detection 5) Predictive insights. Data to analyze:',
    SUFFIX 'Present your findings in a structured format with clear sections for each analysis step. Include confidence levels for your predictions and recommendations for data-driven decision making.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for multi-language support
CREATE EXTERNAL MODEL multilingual_model
FUNCTION process_multilingual_text
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-haiku-20240307-v1:0',
    PROMPT 'Process the following text which may be in multiple languages. Detect the language(s) and provide appropriate analysis:',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for real-time recommendations
CREATE EXTERNAL MODEL recommendation_model
FUNCTION generate_recommendations
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1',
    PROMPT 'Based on the provided user behavior and preferences data, generate personalized recommendations:',
    SUFFIX 'Rank recommendations by relevance score and include explanations for each suggestion.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for sentiment analysis with detailed output
CREATE EXTERNAL MODEL detailed_sentiment_model
FUNCTION analyze_sentiment_detailed
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-sonnet-20240229-v1:0',
    PROMPT 'Perform detailed sentiment analysis on the following text. Consider emotional tone, intent, urgency, and contextual factors:',
    SUFFIX 'Provide sentiment score (-1 to 1), confidence level, key emotional indicators, and actionable insights.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for competitive analysis
CREATE EXTERNAL MODEL competitive_analysis_model
FUNCTION analyze_competition
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-3-haiku-20240307-v1:0',
    PROMPT 'Analyze the following competitive landscape data and provide strategic insights:',
    SUFFIX 'Include market positioning, competitive advantages, threats, and strategic recommendations.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);

-- CREATE EXTERNAL MODEL for document classification
CREATE EXTERNAL MODEL document_classifier_model
FUNCTION classify_documents
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftBedrockRole'
MODEL_TYPE BEDROCK
SETTINGS (
    MODEL_ID 'anthropic.claude-v2:1',
    PROMPT 'Classify the following document into appropriate categories and extract key metadata:',
    SUFFIX 'Provide classification confidence, document type, key entities, and suggested filing category.',
    REQUEST_TYPE UNIFIED,
    RESPONSE_TYPE SUPER
);