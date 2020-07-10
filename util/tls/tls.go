package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/xmlking/toolkit/util/ioutil"
)

// NewTLSConfig returns a TLS config that includes a certificate
// Use for Server TLS config or when using a client certificate
// If caPath is empty, system CAs will be used
func NewTLSConfig(certPath, keyPath, caPath, serverName string, password string) (tlsConfig *tls.Config, err error) {
	var certPEMBlock, keyPEMBlock []byte
	certPEMBlock, err = ioutil.ReadFile(certPath)
	if err != nil {
		return
	}
	keyPEMBlock, err = ioutil.ReadFile(keyPath)
	if err != nil {
		return
	}

	// unwrap keyPEMBlock, if protected with password
	keyDERBlock, _ := pem.Decode(keyPEMBlock)
	log.Debug().Msgf("Is Encrypted Private Key: %v", x509.IsEncryptedPEMBlock(keyDERBlock))
	if x509.IsEncryptedPEMBlock(keyDERBlock) {
		var decryptedKeyBytes []byte
		decryptedKeyBytes, err = x509.DecryptPEMBlock(keyDERBlock, []byte(password))
		if err != nil {
			return
		}
		keyDERBlock = &pem.Block{
			Type:  keyDERBlock.Type,
			Bytes: decryptedKeyBytes,
		}
		keyPEMBlock = pem.EncodeToMemory(keyDERBlock)
	}

	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}

	roots, err := loadRoots(caPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,
		// ClientCAs:    roots,
		NextProtos: []string{"h2"},
		MinVersion: tls.VersionTLS12,
	}, nil
}

func loadRoots(caPath string) (*x509.CertPool, error) {
	if caPath == "" {
		return nil, nil
	}

	roots := x509.NewCertPool()
	pem, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", caPath, err)
	}
	ok := roots.AppendCertsFromPEM(pem)
	if !ok {
		return nil, fmt.Errorf("could not read root certs: %s", err)
	}
	return roots, nil
}
