package main

import (
	"path/filepath"
	"go/build"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const importpath = "github.com/PaulSamways/gopkgsearch"

func webdir() string {
	pkg, err := build.Import(importpath, "", build.FindOnly)
	if err != nil {
		log.Printf("Couldn't determine web directory: %s", err)
		return "./web"
	}
	return filepath.Join(pkg.Dir, "web")
}

func listen() {
	http.Handle("/", http.FileServer(http.Dir(webdir())))
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

// Matches two values, returning a score 0..1. 0 no match, 1 exact match.
func match(a, b string) float32 {
	idx := strings.Index(a, b)

	if idx == -1 {
		return 0
	}

	la, lb := float32(len(a)), float32(len(b))

	if idx == 0 && la == lb {
		return 1
	}

	lenScore := 1.0 - ((la - lb) / la)
	posScore := (la - float32(idx)) / la

	return (lenScore + posScore) / 2.0
}

type score struct {
	element *Element
	score   float32
}
type scores []score

func (s scores) Len() int           { return len(s) }
func (s scores) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s scores) Less(i, j int) bool { return s[i].score > s[j].score }

func query(w http.ResponseWriter, req *http.Request) {
	result := make(scores, 0)

	start := time.Now()

	queries := req.URL.Query()["q"]

	obj := ""
	pkg := ""

	if len(queries) > 0 {
		parts := strings.Split(queries[0], ".")

		if len(parts) > 1 {
			pkg = parts[0]
			obj = strings.ToLower(parts[1])
		} else {
			obj = strings.ToLower(parts[0])
		}
	}

	for _, e := range elements {
		if len(pkg) > 0 && e.Package != pkg {
			continue
		}

		if m := match(e.name, obj); m > 0 {
			result = append(result, score{e, m})
		}
	}

	sort.Sort(result)

	c := len(result)
	if c > 50 {
		c = 50
	}

	r := make([]*Element, c)

	for i, v := range result[:c] {
		r[i] = v.element
	}

	fmt.Printf("Found %d results in %s.\n", c, time.Since(start))

	writeJson(w, r)
}
