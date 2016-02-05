package main

import (
	"github.com/gorilla/mux"
	"github.com/praveenmenon/golang_demo/api/v1/controllers/account"
	"github.com/praveenmenon/golang_demo/api/v1/controllers/users"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	// Account Routes
	r.HandleFunc("/sign_up", account.Registration.Create).Methods("POST")
	r.HandleFunc("/list", users.List.List_users).Methods("GET")
	// r.HandleFunc("/log_out/{devise_token:([a-zA-Z0-9]+)?}", account.Session.Destroy).Methods("GET")
	// r.HandleFunc("/forgot_password/{email}", account.Forgot_password.Check_email).Methods("GET")

	http.Handle("/", r)
	// http.Handle("/", handler)

	// HTTP Listening Port
	log.Println("main : Started : Listening on: http://localhost:3000 ...")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))

}
