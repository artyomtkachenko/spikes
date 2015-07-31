package main

import (
  "fmt"
  "gopkg.in/xmlpath.v2"
  "strings"
  "io/ioutil"
  "net/http"
)

func parseWithXpath(data []byte){
  // This funcion works fine, but i can not find a way to use unions in xpath query
  href       := xmlpath.MustCompile("//table//td[1]//@href | //table//td[6]")
  // statusPath := xmlpath.MustCompile("//table//td[6]")
  root, err  := xmlpath.ParseHTML(strings.NewReader(string(data)))

  if err != nil {
    panic(err)
  }

  iter := href.Iter(root) // Search happens here
  for iter.Next() {
    node := iter.Node().String()
    fmt.Println(node)
  }
}

func main() {
  response, err := http.Get("http://localhost:4568/balancer-manager")
  defer response.Body.Close()

  if err != nil { panic(err) }
  body, _ := ioutil.ReadAll(response.Body)
  // fmt.Println(body)
  parseWithXpath(body)
}
