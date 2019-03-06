package auth

import (
	"encoding/json"
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request) {
	data := RegisterRequest{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
func auth(w http.ResponseWriter, r *http.Request) {}

/* handle func. get pr create token*/
func Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		auth(w, r)
	case http.MethodPost:
		register(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}
