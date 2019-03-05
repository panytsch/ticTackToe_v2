package main

import (
	"flag"
	"log"
	"net/http"
	"src/github.com/gorilla/mux"
	"ticTackToe_v2/src/api/v1/routes"
	"time"
)

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "./build", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()
	//get index file or other
	r.Handle(dir, http.FileServer(http.Dir(dir)))

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static"))))
	for route, handler := range routes.RoutesV1 {
		go r.HandleFunc(route, handler)
	}

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
