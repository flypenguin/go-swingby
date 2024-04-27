package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// cannot use ":=" here, see https://stackoverflow.com/a/50875851/902327
var listen_port = 8081

func redirect_path(w http.ResponseWriter, r *http.Request) {
	info := fmt.Sprintf("REDIRECT 'PATH': %s\n", r.URL.Path)
	fmt.Printf(info)
	fmt.Fprintf(w, info)
}

func redirect_rest(w http.ResponseWriter, r *http.Request) {
	rest := mux.Vars(r)["rest"]
	info := fmt.Sprintf("REDIRECT 'REST': %s\n", rest)
	fmt.Printf(info)
	fmt.Fprintf(w, info)
}

func index(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("Handler for / called\n")
	fmt.Fprintf(w, "Information forthcoming")
}

func handle_test(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("Handler for /test called\n")
	fmt.Fprintf(w, "Test handler\n")
}

func main() {
	fmt.Printf("Listening on port :%d ...\n", listen_port)

	// that router is then used in the http package's "ListenAndServe()" function
	// note we're using the gorilla.mux router
	r := mux.NewRouter()

	// in "basic" examples, we use "http.HandleFunc()", but not here.
	// this is so good.
	// HandlerFunc is a _type_ !!  (see https://www.willem.dev/articles/http-handler-func/)
	// what we do here is a type conversion, and the resulting type ("HandlerFunc")
	// has a member function "ServeHTTP(w, r)", which simply calls "redirect(w, r)"
	// so, basically, a function with a member function calling the actual function.
	// (this here works very much like a python decorator)
	r.PathPrefix("/-").
		Handler(
			http.StripPrefix(
				"/-/",
				http.HandlerFunc(redirect_path)))

	// this works too, but we want to strip the prefix before that.
	//r.PathPrefix("/-").HandlerFunc(redirect)

	r.HandleFunc("/r/{rest:.+}", redirect_rest)

	// catch-all
	r.HandleFunc("/", index)

	// "mux.NewRouter()" is used here – see at the end.
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listen_port), r))
}
