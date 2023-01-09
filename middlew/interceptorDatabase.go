package middlew

import (
	"github.com/el27egs/twittor-tilotta/db"
	"net/http"
)

func CheckDatabaseConnection(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db.CheckConnection() == false {
			http.Error(w, "No se pudo conectar a la base de datos", 500)
			return
		}
		next.ServeHTTP(w, r)
	}
}
