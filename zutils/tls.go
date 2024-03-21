package zutils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

func MakeTLSConfig(tlsCert, tlsKey string, clientCa string) (*tls.Config, error) {
	tlsConfig := &tls.Config{}
	tlsCertValue := atomic.Value{}
	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		return nil, err
	}
	tlsCertValue.Store(cert)
	if clientCa != "" {
		if caBytes, err := os.ReadFile(clientCa); err != nil {
			return nil, err
		} else {
			ca := x509.NewCertPool()
			if ok := ca.AppendCertsFromPEM(caBytes); !ok {
				return nil, fmt.Errorf("failed to parse %v ", clientCa)
			}
			tlsConfig.ClientCAs = ca
		}
	}

	if _, err := Watch(tlsCert, time.Second*10, func() {
		c, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			fmt.Errorf("service failed to load x509 key pair: %v", err)
			return
		}
		tlsCertValue.Store(c)
	}); err != nil {
		return nil, err
	}
	tlsConfig.GetCertificate = func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
		c := tlsCertValue.Load()
		if c == nil {
			return nil, fmt.Errorf("certificate not loaded")
		}
		res := c.(tls.Certificate)
		return &res, nil
	}

	return tlsConfig, nil

}
