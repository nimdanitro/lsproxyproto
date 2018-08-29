package main

import (
	"crypto/tls"
	"log"
)

func getTLSConfig(cert, key string) *tls.Config {
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatalf("error in tls.LoadX509KeyPair: %s", err)
		return nil
	}

	config := &tls.Config{
		Certificates:       []tls.Certificate{c},
		InsecureSkipVerify: false,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
			tls.X25519,
		},

		PreferServerCipherSuites: true,
	}

	return config
}
