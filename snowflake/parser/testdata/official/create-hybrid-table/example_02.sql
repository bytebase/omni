INSERT INTO mytable (customer_id, full_name, email, extended_customer_info)
  SELECT 100, 'Jane Doe', 'jdoe@example.com',
    parse_json('{"address": "1234 Main St", "city": "San Francisco", "state": "CA", "zip":"94110"}');
