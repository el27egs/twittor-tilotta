package middlew

import (
	"github.com/el27egs/twittor-tilotta/routers"
	"net/http"
)

func ValidateJWT(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		_, _, _, err := routers.ValidateJWT(r.Header.Get("Authorization"))

		if err != nil {
			http.Error(w, "Error al validar el JWT token "+err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	}

}
