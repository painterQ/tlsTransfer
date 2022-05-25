
SHELL := /bin/bash
DOMAIN := 'subjectAltName=DNS:*.xixi201314.cn,DNS:localhost\nnsCertType=server\nextendedKeyUsage=serverAuth\nsubjectKeyIdentifier=hash\nauthorityKeyIdentifier=keyid,issuer'

cert:
	@rm -rf ./tmp
	@mkdir tmp
	@openssl ecparam -name prime256v1 -genkey > ./tmp/ca.priv
	@openssl req -new -key ./tmp/ca.priv -subj "/C=CN/O=root/CN=Root.Painter" | \
	    openssl x509 -req  -days 36500 -signkey ./tmp/ca.priv \
	    -extfile <(echo -e 'nsCertType = sslCA\nkeyUsage = cRLSign, keyCertSign\nbasicConstraints = critical,CA:true\nsubjectKeyIdentifier=hash\nauthorityKeyIdentifier=keyid,issuer') -out ./tmp/ca.crt

	@openssl ecparam -name prime256v1 -genkey > ./tmp/tls_server.priv
	@openssl req -new -key ./tmp/tls_server.priv -subj "/C=CN/O=server/CN=Painter" | \
		openssl x509 -req  -days 36500 -CA ./tmp/ca.crt -CAcreateserial -CAkey ./tmp/ca.priv  \
		-extfile <(echo -e $(DOMAIN)) -out ./tmp/tls_server.crt

	@openssl ecparam -name prime256v1 -genkey > ./tmp/tls_client.priv
	@openssl req -new -key ./tmp/tls_client.priv -subj "/C=CN/O=client/CN=Painter" | \
		openssl x509 -req  -days 36500 -CA ./tmp/ca.crt -CAcreateserial -CAkey ./tmp/ca.priv \
		-extfile <(echo -e 'nsCertType=client,email,objsign\nextendedKeyUsage=clientAuth,emailProtection\nsubjectKeyIdentifier=hash\nauthorityKeyIdentifier=keyid,issuer')  \
		-out ./tmp/tls_client.crt
	@openssl pkcs12 -export -out ./tmp/tls_client.pfx -CAfile ./tmp/ca.crt -inkey ./tmp/tls_client.priv -in ./tmp/tls_client.crt
	@rm ./tmp/*.srl || echo ""