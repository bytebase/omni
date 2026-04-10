SELECT id,
    f1.value:type::STRING AS contact_type,
    f1.value:content::STRING AS contact_details
  FROM persons p,
    LATERAL FLATTEN(INPUT => p.c, PATH => 'contact') f,
    LATERAL FLATTEN(INPUT => f.value:business) f1;
