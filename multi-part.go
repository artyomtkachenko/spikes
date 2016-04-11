package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gorilla/mux"

	"log"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	routes := Routes{
		Route{
			Name:        "DoBuild",
			Method:      "POST",
			Pattern:     "/release/v1/build/rpm",
			HandlerFunc: DoBuild,
		},
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func DoBuild(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(0)

	file1, header, _ := req.FormFile("config")
	file2 := req.FormValue("meta")
	body1, _ := ioutil.ReadAll(file1)
	fmt.Printf("%+v\n", header)
	fmt.Printf("%+v\n", string(body1))
	fmt.Printf("%+v\n", string(file2))
}

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe("localhost:9999", router))
}
