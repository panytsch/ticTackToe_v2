package main

import (
	"flag"
	"log"
	"net/http"
	"src/github.com/gorilla/mux"
	"ticTackToe_v2/src/api/v1/routes"
	"ticTackToe_v2/src/socketTest"
	"time"
)

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "./build", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()
	//get index file or other
	r.Handle("/", http.FileServer(http.Dir(dir)))

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

	bus := socketTest.NewBus()
	go bus.Run()
	go socketTest.RunJoker(bus)

	r.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		// апгрейд соединения
		ws, err := socketTest.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		bus.Register <- ws

	})

	log.Fatal(srv.ListenAndServe())
}
