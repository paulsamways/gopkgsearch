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

func isMatch(a, b string) bool {
	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}

func query(w http.ResponseWriter, req *http.Request) {
	result := make([]interface{}, 0)

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

	for _, e := range elements {
		if len(pkg) > 0 && strings.ToLower(e.Package) != pkg {
			continue
		}

		if isMatch(e.Name, obj) {
			result = append(result, e)
		}
	}

	writeJson(w, result)
}
