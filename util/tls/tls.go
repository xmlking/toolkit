package tls

import (
	"crypto/tls"
	"crypto/x509"
	"io/fs"

	"github.com/cockroachdb/errors"
)

// NewServerTLSConfig returns a TLS Config for a server connection.
// when verifyPeer=true, strictly verify client certs only issued by the designated CA (caPath)
func NewServerTLSConfig(xfs fs.FS, certPath, keyPath, caPath string, verifyPeer bool) (serverTLSConfig *tls.Config, err error) {
	var cert tls.Certificate
	if cert, err = loadCert(xfs, certPath, keyPath); err != nil {
		return
	}

	var clientCAs = x509.NewCertPool()
	var caPem []byte
	if caPem, err = fs.ReadFile(xfs, caPath); err != nil {
		return
	}
	if ok := clientCAs.AppendCertsFromPEM(caPem); !ok {
		return nil, errors.Newf("error loading caPath: %s", caPath)
	}

	serverTLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    clientCAs,
		NextProtos:   []string{"h2"},
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.VerifyClientCertIfGiven,
	}

	// TODO: Should we set InsecureSkipVerify=true and use VerifyPeerCertificate VerifyConnection?
	if verifyPeer {
		serverTLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return
}

// NewClientTLSConfig returns a TLS Config for a client connection.
// `caPath` can be empty, in that case, RootCAs is nil, TLS uses the host's root CA set.
func NewClientTLSConfig(xfs fs.FS, certPath, keyPath, caPath, serverName string) (clientTLSConfig *tls.Config, err error) {
	var cert tls.Certificate
	if cert, err = loadCert(xfs, certPath, keyPath); err != nil {
		return
	}

	var rootCAs *x509.CertPool
	if caPath != "" {
		if rootCAs, err = x509.SystemCertPool(); err != nil {
			return
		}

		var caPem []byte
		if caPem, err = fs.ReadFile(xfs, caPath); err != nil {
			return
		}

		if ok := rootCAs.AppendCertsFromPEM(caPem); !ok {
			return nil, errors.Newf("error loading caPath: %s", caPath)
		}
	}

	clientTLSConfig = &tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCAs, // if RootCAs is nil, TLS uses the host's root CA set.
		NextProtos:   []string{"h2"},
		MinVersion:   tls.VersionTLS12,
	}

	return
}

func loadCert(xfs fs.FS, certPath, keyPath string) (cert tls.Certificate, err error) {
	var certPEMBlock, keyPEMBlock []byte
	certPEMBlock, err = fs.ReadFile(xfs, certPath)
	if err != nil {
		return
	}
	keyPEMBlock, err = fs.ReadFile(xfs, keyPath)
	if err != nil {
		return
	}

	return tls.X509KeyPair(certPEMBlock, keyPEMBlock)
}
