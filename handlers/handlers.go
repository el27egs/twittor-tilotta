package handlers

import (
	"github.com/el27egs/twittor-tilotta/middlew"
	"github.com/el27egs/twittor-tilotta/routers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	router := mux.NewRouter()

	router.HandleFunc("/registro", middlew.CheckDBConnection(routers.CreateUser)).Methods("POST")
	router.HandleFunc("/login", middlew.CheckDBConnection(routers.Login)).Methods("POST")
	router.HandleFunc("/verperfil", middlew.CheckDBConnection(middlew.ValidateJWT(routers.GetUser))).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	handler := cors.AllowAll().Handler(router)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
