package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//go:embed html/index.html
var index_html string

// cannot use ":=" here, see https://stackoverflow.com/a/50875851/902327
var listen_port = 8081

var host_override string

func redirect_xxx(w http.ResponseWriter, r *http.Request, redirect_method int) {
	u, _ := url.Parse(mux.Vars(r)["url"])
	u.RawQuery = r.URL.RawQuery
	u.RawFragment = r.URL.RawFragment
	redirect_target := fmt.Sprintf("//%s", u.String())

	fmt.Printf("[DEBUG] redirect_target=%s\n", redirect_target)
	http.Redirect(w, r, redirect_target, redirect_method)
}

func redirect_301(w http.ResponseWriter, r *http.Request) {
	redirect_xxx(w, r, http.StatusMovedPermanently)
}

func index(w http.ResponseWriter, r *http.Request) {
	var use_host string
	fmt.Printf("Handler for / called\n")
	fmt.Printf("[DEBUG] r.Host=%s\n", r.Host)
	if host_override != "" {
		use_host = host_override
	} else {
		use_host = r.Host
	}
	fmt.Fprintf(w, strings.ReplaceAll(index_html, "this.url", use_host))
}

func main() {
	host_override = os.Getenv("SWINGBY_HOST")
	fmt.Printf("Listening on port :%d ...\n", listen_port)

	// that router is then used in the http package's "ListenAndServe()" function
	// note we're using the gorilla.mux router
	r := mux.NewRouter()

	// https://stackoverflow.com/a/49312740/902327
	r.Use(handlers.ProxyHeaders)

	// that's all we need?
	r.HandleFunc("/{marker:[^/]+}/{url:.+}", redirect_301)

	// catch-all
	r.HandleFunc("/", index)

	// "mux.NewRouter()" is used here – see at the end.
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listen_port), r))
}
