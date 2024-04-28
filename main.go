package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// cannot use ":=" here, see https://stackoverflow.com/a/50875851/902327
var listen_port = 8081

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

func index(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("Handler for / called\n")
	fmt.Fprintf(w, "Information forthcoming")
}

func main() {
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
