package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"strings"
	"time"
)

func verifyChain(chain []*x509.Certificate) bool {

	var result = make(map[string]bool)
	rootFound := false
	rootPosition := 0
	for i, cert := range chain {
		issuer := cert.Issuer.CommonName
		subject := cert.Subject.CommonName
		if strings.Contains(issuer, "Root") {
			rootPosition = i
			rootFound = true
		}
		// Do not really like this logic. TODO FIXME
		for j, cert2 := range chain {
			if j == i {
				continue //We do not want to verify the same entry from the outer loop
			}
			if subject == cert2.Issuer.CommonName { // We found a record in the chain
				result[subject] = true
			}
		}
	}

	for _, v := range result {
		if !v {
			return false
		}
	}
	if !rootFound {
		fmt.Println("NO ROOT cerificate detected in the chain")
		return false
	}

	// It is a bit Ugly, as it does not check position for NON ROOT certs
	if rootPosition != len(chain)-1 {
		fmt.Println("Wrong ROOT cerificate position detected")
	}
	return true
}

func main() {
	addr := flag.String("addr", "", "Address in form of host:port")
	flag.Parse()
	conn, err := tls.Dial("tcp", *addr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	chain := conn.ConnectionState().PeerCertificates
	if !verifyChain(chain) {
		for _, cert := range chain {
			fmt.Printf("Certificate: %s, Issued by: %s, Expires at:  %s, Days left: %d\n",
				cert.Subject.CommonName, cert.Issuer.CommonName, cert.NotAfter, cert.NotAfter.Sub(time.Now())/time.Hour/24)
		}
	}

	conn.Close()
}
