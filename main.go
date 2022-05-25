package main

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"net"
	"net/http"
)

//go:embed tmp/ca.crt
var ca []byte

//go:embed tmp/tls_server.crt
var cert []byte

//go:embed tmp/tls_server.priv
var priv []byte

func main() {
	cert, err := tls.X509KeyPair(cert, priv)
	if err != nil {
		panic(err)
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(ca); !ok {
		panic(ok)
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		CipherSuites: []uint16{tls.TLS_AES_256_GCM_SHA384},
	}

	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	tl := tls.NewListener(l, &config)

	http.Serve(tl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello world!"))
	}))
	fmt.Println("end")
}
