ALTER SYSTEM ADD BROKER broker_name "host1:port", "host2:port";
ALTER SYSTEM ADD BROKER broker_fqdn1 "broker_fqdn1:port";
ALTER SYSTEM DROP ALL BROKER broker_name;
ALTER SYSTEM DROP BROKER broker_name "10.10.10.1:8000";
SHOW BROKER
