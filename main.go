package main

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"io"
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
		rr, err := http.Post("https://www.baidu.com", "text", r.Body)
		if err != nil {
			fmt.Printf("###err:%v", err)
		}
		io.CopyBuffer(w, rr.Body, make([]byte, 128))
	}))
	fmt.Println("end")
}
