package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
)

type position struct {
	expected int
	actual   int
}

type certMeta struct {
	name     string
	position int
}

func reverse(numbers []certMeta) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

func verifyPositions(nativeChain []certMeta) {
	reverse(nativeChain)
	for i, item := range nativeChain {
		if i == 0 || i == 0 {
			if item.position != 0 {
				fmt.Printf("Wrong position found for %s, expected position in the chain: %d, actual position in the chain: %d\n", item.name, i, item.position)
			}
		} else {
			if item.position != (i - 1) {
				fmt.Printf("Wrong position found for %s, expected position in the chain: %d, actual position in the chain: %d\n", item.name, i-1, item.position)
			}
		}
	}

}

func findRoot(chain []*x509.Certificate) (error, string) {
	knownRoots := map[string]string{
		"GeoTrust Global CA":                     "google",
		"NBN Co Root CA":                         "nbn",
		"DigiCert SHA2 High Assurance Server CA": "facebook",
	}
	for _, item := range chain {
		_, inIssuer := knownRoots[item.Issuer.CommonName]
		_, inSubject := knownRoots[item.Subject.CommonName]
		if inIssuer {
			return nil, item.Issuer.CommonName
		}
		if inSubject {
			return nil, item.Subject.CommonName
		}
	}
	return errors.New("NO ROOT certificate found in the chain"), ""
}

func buildNativeChain(chain []*x509.Certificate, issuer string, result *[]certMeta, step int) {
	next := "none"
	if step == 0 {
		*result = append(*result, certMeta{chain[0].Subject.CommonName, 0})
	} else {
		for i, item := range chain {
			if item.Issuer.CommonName == issuer {
				*result = append(*result, certMeta{issuer, i})
				next = item.Subject.CommonName
			}
		}
		buildNativeChain(chain, next, result, step-1)
	}
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
	var nativeChain []certMeta
	chain := conn.ConnectionState().PeerCertificates
	err, root := findRoot(chain)
	if err != nil {
		fmt.Println(err)
	}
	buildNativeChain(chain, root, &nativeChain, len(chain))

	verifyPositions(nativeChain)
	fmt.Printf("%+v\n", nativeChain)
	// for _, cert := range chain {
	// 	fmt.Printf("Certificate: %s, Issued by: %s, Expires at:  %s, Days left: %d\n",
	// 		cert.Subject.CommonName, cert.Issuer.CommonName, cert.NotAfter, cert.NotAfter.Sub(time.Now())/time.Hour/24)
	// }

	conn.Close()
}
