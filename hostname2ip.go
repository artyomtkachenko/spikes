package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func main() {

	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	arr := strings.Split(string(content), "\n")
	for _, line := range arr {
		if strings.Contains(line, "sv") {
			ip, err := net.LookupIP(line)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s %s\n", line, ip[0])
		}
	}

}
