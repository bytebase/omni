-- DROP MASKING POLICY command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_MASKING_POLICY.html

-- Basic DROP MASKING POLICY syntax
DROP MASKING POLICY email_mask;

-- DROP MASKING POLICY with various policy names
DROP MASKING POLICY ssn_mask;
DROP MASKING POLICY credit_card_mask;
DROP MASKING POLICY phone_mask;
DROP MASKING POLICY personal_data_mask;

-- DROP MASKING POLICY with underscored names
DROP MASKING POLICY customer_email_mask;
DROP MASKING POLICY employee_ssn_mask;
DROP MASKING POLICY payment_card_mask;
DROP MASKING POLICY contact_phone_mask;

-- DROP MASKING POLICY with quoted identifiers
DROP MASKING POLICY "EmailMask";
DROP MASKING POLICY "SSN-Mask";
DROP MASKING POLICY "Credit Card Mask";
DROP MASKING POLICY "Phone Number Mask";

-- DROP MASKING POLICY with mixed case names
DROP MASKING POLICY EmailMaskPolicy;
DROP MASKING POLICY SSNMaskPolicy;
DROP MASKING POLICY CreditCardMaskPolicy;
DROP MASKING POLICY PhoneNumberMaskPolicy;

-- DROP MASKING POLICY with descriptive names
DROP MASKING POLICY pii_email_protection;
DROP MASKING POLICY sensitive_ssn_mask;
DROP MASKING POLICY financial_card_mask;
DROP MASKING POLICY contact_info_mask;

-- DROP MASKING POLICY with domain-specific names
DROP MASKING POLICY healthcare_patient_id_mask;
DROP MASKING POLICY finance_account_number_mask;
DROP MASKING POLICY retail_customer_data_mask;
DROP MASKING POLICY hr_employee_info_mask;

-- DROP MASKING POLICY with security-focused names
DROP MASKING POLICY confidential_data_mask;
DROP MASKING POLICY restricted_info_mask;
DROP MASKING POLICY private_details_mask;
DROP MASKING POLICY secure_field_mask;

-- DROP MASKING POLICY with functional names
DROP MASKING POLICY partial_email_mask;
DROP MASKING POLICY full_ssn_mask;
DROP MASKING POLICY partial_card_mask;
DROP MASKING POLICY hash_phone_mask;

-- DROP MASKING POLICY with organizational names
DROP MASKING POLICY sales_team_mask;
DROP MASKING POLICY marketing_data_mask;
DROP MASKING POLICY support_info_mask;
DROP MASKING POLICY analytics_mask;

-- DROP MASKING POLICY with version indicators
DROP MASKING POLICY email_mask_v1;
DROP MASKING POLICY ssn_mask_v2;
DROP MASKING POLICY card_mask_2023;
DROP MASKING POLICY phone_mask_latest;

-- DROP MASKING POLICY with environment indicators
DROP MASKING POLICY prod_email_mask;
DROP MASKING POLICY dev_test_mask;
DROP MASKING POLICY staging_data_mask;
DROP MASKING POLICY qa_validation_mask;

-- Edge cases with special characters in quoted names
DROP MASKING POLICY "mask-with-hyphens";
DROP MASKING POLICY "mask.with.dots";
DROP MASKING POLICY "mask$with$dollars";
DROP MASKING POLICY "mask_with_underscores";

-- DROP MASKING POLICY with longer descriptive names
DROP MASKING POLICY customer_personally_identifiable_information_mask;
DROP MASKING POLICY employee_social_security_number_protection_policy;
DROP MASKING POLICY financial_payment_card_data_security_mask;

-- DROP MASKING POLICY statements without semicolons (should also be valid)
DROP MASKING POLICY test_mask
DROP MASKING POLICY simple_mask

-- DROP MASKING POLICY with common masking use cases
DROP MASKING POLICY address_partial_mask;
DROP MASKING POLICY salary_range_mask;
DROP MASKING POLICY birth_date_year_mask;
DROP MASKING POLICY ip_address_mask;

-- DROP MASKING POLICY for compliance scenarios
DROP MASKING POLICY gdpr_compliance_mask;
DROP MASKING POLICY hipaa_patient_mask;
DROP MASKING POLICY pci_payment_mask;
DROP MASKING POLICY sox_financial_mask;

-- DROP MASKING POLICY for different data types
DROP MASKING POLICY numeric_id_mask;
DROP MASKING POLICY text_content_mask;
DROP MASKING POLICY date_field_mask;
DROP MASKING POLICY json_data_mask;

-- DROP MASKING POLICY for table-specific masks
DROP MASKING POLICY customers_email_mask;
DROP MASKING POLICY employees_ssn_mask;
DROP MASKING POLICY orders_payment_mask;
DROP MASKING POLICY users_profile_mask;