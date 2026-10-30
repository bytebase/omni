package parser

func init() {
	keywords["get_master_public_key"] = kwGET_MASTER_PUBLIC_KEY
	keywords["master_auto_position"] = kwMASTER_AUTO_POSITION
	keywords["master_bind"] = kwMASTER_BIND
	keywords["master_compression_algorithms"] = kwMASTER_COMPRESSION_ALGORITHMS
	keywords["master_connect_retry"] = kwMASTER_CONNECT_RETRY
	keywords["master_delay"] = kwMASTER_DELAY
	keywords["master_heartbeat_period"] = kwMASTER_HEARTBEAT_PERIOD
	keywords["master_host"] = kwMASTER_HOST
	keywords["master_log_file"] = kwMASTER_LOG_FILE
	keywords["master_log_pos"] = kwMASTER_LOG_POS
	keywords["master_password"] = kwMASTER_PASSWORD
	keywords["master_port"] = kwMASTER_PORT
	keywords["master_public_key_path"] = kwMASTER_PUBLIC_KEY_PATH
	keywords["master_retry_count"] = kwMASTER_RETRY_COUNT
	keywords["master_ssl"] = kwMASTER_SSL
	keywords["master_ssl_ca"] = kwMASTER_SSL_CA
	keywords["master_ssl_capath"] = kwMASTER_SSL_CAPATH
	keywords["master_ssl_cert"] = kwMASTER_SSL_CERT
	keywords["master_ssl_cipher"] = kwMASTER_SSL_CIPHER
	keywords["master_ssl_crl"] = kwMASTER_SSL_CRL
	keywords["master_ssl_crlpath"] = kwMASTER_SSL_CRLPATH
	keywords["master_ssl_key"] = kwMASTER_SSL_KEY
	keywords["master_ssl_verify_server_cert"] = kwMASTER_SSL_VERIFY_SERVER_CERT
	keywords["master_tls_ciphersuites"] = kwMASTER_TLS_CIPHERSUITES
	keywords["master_tls_version"] = kwMASTER_TLS_VERSION
	keywords["master_user"] = kwMASTER_USER
	keywords["master_zstd_compression_level"] = kwMASTER_ZSTD_COMPRESSION_LEVEL

	keywordCategories[kwMASTER_BIND] = kwCatReserved
	keywordCategories[kwMASTER_SSL_VERIFY_SERVER_CERT] = kwCatReserved
}
