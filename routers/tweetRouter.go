package routers

import (
	"encoding/json"
	"github.com/el27egs/twittor-tilotta/db"
	"github.com/el27egs/twittor-tilotta/models"
	"net/http"
	"strconv"
	"time"
)

func CreateTweet(w http.ResponseWriter, r *http.Request) {

	var tweet models.Tweet

	err := json.NewDecoder(r.Body).Decode(&tweet)

	record := models.UserTweet{
		UserID:  IDUser,
		Message: tweet.Message,
		Date:    time.Now(),
	}
	var status bool
	_, status, err = db.CreateTweet(record)

	if err != nil {
		http.Error(w, "Error al crear el tweet "+err.Error(), 400)
		return
	}
	if !status {
		http.Error(w, "No se pudo guardar el tweet", 400)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetTweetsWithPager(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")
	if len(ID) == 0 {
		http.Error(w, "El ID es un parametro requerido", http.StatusBadRequest)
		return
	}

	if len(r.URL.Query().Get("page")) == 0 {
		http.Error(w, "La pagina es un parametro requerido", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		http.Error(w, "Pagina debe ser un valor numerico mayor a 0", http.StatusBadRequest)
		return
	}
	tweets, found := db.GetTweetsWithPager(ID, int64(page))

	if !found {
		http.Error(w, "Error al recuperar los tweets", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tweets)

}

func DeleteTweet(w http.ResponseWriter, r *http.Request) {

	tweetId := r.URL.Query().Get("id")
	if len(tweetId) == 0 {
		http.Error(w, "ID del Tweet es requerido", http.StatusBadRequest)
		return
	}

	err := db.DeleteTweet(tweetId, IDUser)
	if err != nil {
		http.Error(w, "Error al borrar tweet "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
