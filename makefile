
SHELL := /bin/bash
PWD := 123456
cert:
	@rm -rf ./tmp
	@mkdir tmp
	@openssl ecparam -name prime256v1 -genkey > ./tmp/ca.priv
	@openssl req -new -key ./tmp/ca.priv -subj "/C=CN/O=root/CN=painter" | \
	    openssl x509 -req  -days 36500 -signkey ./tmp/ca.priv \
	    -extfile <(echo -e 'nsCertType = sslCA\nkeyUsage = cRLSign, keyCertSign\nbasicConstraints = critical,CA:true') -out ./tmp/ca.cert

	@openssl ecparam -name prime256v1 -genkey > ./tmp/tls_server.priv
	@openssl req -new -key ./tmp/tls_server.priv -subj "/C=CN/O=server/CN=painter" | \
		openssl x509 -req  -days 36500 -CA ./tmp/ca.cert -CAcreateserial -CAkey ./tmp/ca.priv  \
		-extfile <(echo -e 'subjectAltName=DNS:*.com,DNS:*.cn,IP:0.0.0.0') -out ./tmp/tls_server.cert

	@openssl ecparam -name prime256v1 -genkey > ./tmp/tls_client.priv
	@openssl req -new -key ./tmp/tls_client.priv -subj "/C=CN/O=client/CN=painter" | \
		openssl x509 -req  -days 36500 -CA ./tmp/ca.cert -CAcreateserial -CAkey ./tmp/ca.priv  \
		-out ./tmp/tls_client.cert
	@openssl pkcs12 -export -out ./tmp/tls_client.pfx -CAfile ./tmp/ca.cert -inkey ./tmp/tls_client.priv -in ./tmp/tls_client.cert
	@rm ./.srl || echo ""