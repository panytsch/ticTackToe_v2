package routes

import (
	"net/http"
	"ticTackToe_v2/src/api/v1/ctrl/auth"
)

const ApiV1Prefix = "/api/v1"

var RoutesV1 = map[string]func(w http.ResponseWriter, r *http.Request){
	"/hello":              Hello,
	ApiV1Prefix + "/auth": auth.Handle,
}

var Hello = func(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello world"))
}
