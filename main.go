// Test
// Julius Shol
package main

import (
	"log"
	"mime"
	"net/http"
)

// a example middleware
func exampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			our middleware logic goes here.....
		*/
		next.ServeHTTP(w, r)
	})
}

// Make a middleware the
// a)checks for the existence of a content-Type header
// b) if the header exists, check that it has the mime type application/json
// if either of those checks fail. we want our middleware to write an error message
// Then stop the request from reaching our application handlers.
func enforceJsonHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}
			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// /////////////////////////////////////////////////////////////////////////
// Class Codes----------------------
// Middle ware is used to protect the handlers
func middlewareA(next http.Handler) http.Handler {
	// return a http.Handler
	// This is executed on the way down to the handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware A ")

		next.ServeHTTP(w, r) // used to pass the request to the router

		log.Println("Executing middleware A again (after the handler is finished)")
	})
}

func middlewareB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware B ")
		if r.URL.Path == "/cherry" {
			return
		}
		next.ServeHTTP(w, r)
		log.Println("Executing middleware B again ")
	})
}

// Create a hanlder function
func ourHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing a handler...")
	w.Write([]byte("Handler Called :)"))
}

func main() {
	mux := http.NewServeMux() //

	// Online example
	finalHandler := http.HandlerFunc(final)
	mux.Handle("/", enforceJsonHandler(finalHandler))

	// Classs demonstration
	mux.Handle("/check", middlewareA(middlewareB(http.HandlerFunc(ourHandler)))) // chaing handlers
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}

/*
Another example:
*/

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Print("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareTwo")
		next.ServeHTTP(w, r)
		log.Print("Executing middlewareTwo again")
	})
}

// handler function
func finalH(w http.ResponseWriter, r *http.Request) {
	log.Printf("Executing finalHandler")
	w.Write([]byte("OK"))
}
