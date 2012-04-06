package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func listen() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/query", query)

	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatalf("Couldn't listen for http requests on %s. %s", *addr, err)
	}
}

func writeJson(w http.ResponseWriter, o interface{}) bool {
	w.Header().Add("Content-Type", "application/json")

	enc := json.NewEncoder(w)

	if err := enc.Encode(o); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{'err': '%s'}", err)

		return false
	}

	return true
}

func query(w http.ResponseWriter, req *http.Request) {
	result := make([]Element, 0)

	queries := req.URL.Query()["q"]

	obj := ""
	pkg := ""

	if len(queries) > 0 {
		obj = strings.ToLower(queries[0])
		parts := strings.Split(obj, ".")

		if len(parts) > 1 {
			pkg = parts[0]
			obj = parts[1]
		}
	}

	fmt.Printf("Package: '%s' - Object: '%s'\n", pkg, obj)

	for _, v := range elements {
		pkgMatch := true

		if len(pkg) > 0 {
			pkgMatch = strings.ToLower(v.ImportPath) == pkg
		}

		if pkgMatch && strings.Contains(strings.ToLower(v.Name), obj) {
			result = append(result, v)
		}
	}

	writeJson(w, result)
}
