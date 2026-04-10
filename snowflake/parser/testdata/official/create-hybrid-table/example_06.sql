SELECT customer_id, full_name, email, extended_customer_info
  FROM mytable
  WHERE extended_customer_info['state'] = 'CA';
