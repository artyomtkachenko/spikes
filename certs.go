package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"strings"
	"time"
)

type position struct {
	expected int
	actual   int
}

type certMeta struct {
	name     string
	position int
	step     int
}

func knownRoots(cert string) bool {
	roots := []string{
		"GeoTrust Global CA",
	}
	for _, c := range roots {
		if cert == c {
			return true
		}
	}
	return false
}

func getActualPosition(chain []*x509.Certificate, issuer string) int {
	for i, item := range chain {
		if item.Subject.CommonName == issuer {
			return i
		}
	}
	return -1
}

func verifyChain2(chain []*x509.Certificate, issuer string, result *[]certMeta, step int) {
	next := "none"
	if step == 0 {
		*result = append(*result, certMeta{chain[0].Subject.CommonName, len(*result), step})
	} else {
		for i, item := range chain {
			if item.Issuer.CommonName == issuer {
				*result = append(*result, certMeta{issuer, i, step})
				next = item.Subject.CommonName
			}
		}
		verifyChain2(chain, next, result, step-1)
	}
}

// func verifyChain3(chain []*x509.Certificate, issuer string, result *[]certMeta, step int) {
// 	if len(chain) == 1 || step == 0 {
// 		*result = append(*result, certMeta{chain[0].Subject.CommonName, len(*result)})
// 	} else {
// 		if chain[step-1].Issuer.CommonName == issuer {
// 			fmt.Printf("Step: %d, Size: %d\n", step, len(chain))
// 			*result = append(*result, certMeta{issuer, step + len(*result)}) //WRONG
// 			chain[step-1].Issuer.CommonName = "none"                         // Found it, we do not want to check it again
// 			verifyChain2(chain, chain[step-1].Subject.CommonName, result, step-1)
// 		} else {
// 			fmt.Printf("STEP: %d, SIZE: %d\n", step, len(chain))
// 			verifyChain2(chain, issuer, result, step-1)
// 		}
// 	}
// }

func verifyChain(chain []*x509.Certificate) bool {

	var nativeChain = make(map[string]string)
	var positions = make(map[string]*position)

	rootFound := false
	for i, item := range chain {
		issuer := item.Issuer.CommonName
		subject := item.Subject.CommonName
		positions[subject] = &position{0, i}
		nativeChain[subject] = issuer
		if strings.Contains(issuer, "Root") || strings.Contains(issuer, "GeoTrust Global CA") || issuer == "" {
			rootFound = true
			positions[issuer] = &position{i + 1, i + 1}
			positions[subject].expected = len(chain) - 1
		}
	}

	for subject, issuer := range nativeChain {
		if issuer != "" {
			issuerExpectedPosition := positions[nativeChain[subject]].expected
			positions[subject].expected = issuerExpectedPosition - 1
		}
	}

	// for k, v := range positions {
	// 	fmt.Printf("%s %+v\n", k, *v)
	// }
	// move it upper
	if !rootFound {
		fmt.Println("NO ROOT cerificate detected in the chain")
		return false
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
	var result []certMeta
	chain := conn.ConnectionState().PeerCertificates
	verifyChain2(chain, "GeoTrust Global CA", &result, len(chain))
	/* verifyChain2(chain, "NBN Co Root CA", &result, len(chain)) */
	/* verifyChain2(chain, "DigiCert SHA2 High Assurance Server CA", &result, len(chain)) */
	fmt.Printf("%+v\n", result)
	for _, cert := range chain {
		fmt.Printf("Certificate: %s, Issued by: %s, Expires at:  %s, Days left: %d\n",
			cert.Subject.CommonName, cert.Issuer.CommonName, cert.NotAfter, cert.NotAfter.Sub(time.Now())/time.Hour/24)
	}

	conn.Close()
}
