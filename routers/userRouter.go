package routers

import (
	"encoding/json"
	"fmt"
	"github.com/el27egs/twittor-tilotta/db"
	"github.com/el27egs/twittor-tilotta/models"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Input data error - "+err.Error(), 400)
		return
	}
	if len(u.Email) == 0 {
		http.Error(w, "email is required", 400)
		return
	}
	if len(u.Password) < 6 {
		http.Error(w, "password length must be grater than 6", 400)
		return
	}
	/*
			Aqui se puede ver que cuando se genera el decode a u,
			sl valor del ObjectID es 0000...00 y asi se va a MongoDB cuando _id en el bson no tenia omitempty
			ahora con el omitempty, ese valor no se manda a Mongo y mongo genera uno nuevo automaticamente.

		    Asi como esta contruido el modelo sin usar bson.M, lo que pasa es que en Go cuando no se asigna nada
		    se da un valor por default dependiendo del tipo de datos, por ejemplo si avatar no se entrega en la
		    peticion de entrada, se da el dafult de cadena vacia, si en el modelo dice omitempty, entonces no se
			guardara en absoluto en MongoDB, si el omitempty, se guardara el valor por default que se asigno en Go
			el cual es la cadena vacia, es decir el campo avatar estara en la coleccion de MongoDB, esto mismo
			pasaba con el campo _id, sele daba el valor 000...00 y como no tenia el omitempty, se guardaba ese valor
			en la coleccion en MongoDB, como sigue el codigo sin asiganarle un valor pero se agrego el omitempty,
			ese valor por default se omite y Mongo por default asigna un valor a ese campo especial _id.
			Falta ver que pasa si construyo un objeto con bson.M en lugar de con el Decode(&models.User)
	*/
	fmt.Println("nuevo usuario ", u)
	_, userFound, _ := db.SearchUserByEmail(u.Email)
	if userFound == true {
		http.Error(w, "User already exists", 400)
		return
	}
	_, status, err := db.SaveUser(u)
	if err != nil {
		http.Error(w, "Error on saving user "+err.Error(), 400)
		return
	}
	if status == false {
		http.Error(w, "Status was false on saving user ", 400)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")
	if len(ID) == 0 {
		http.Error(w, "ingresar un id valido a buscar", http.StatusBadRequest)
		return
	}
	user, err := db.SearchUserByID(ID)

	if err != nil {
		http.Error(w, "usuario no encontrado "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, "datos de entrada incorrectos", 400)
		return
	}

	var status bool
	status, err = db.UpdateUser(u, IDUser)

	if err != nil {
		http.Error(w, "ocurrio un error al actualizar los datos "+err.Error(), 400)
		return
	}

	if !status {
		http.Error(w, "datos del usuario no se actualizacon, intente nuevamente mas tarde", 400)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {

	sourceFile, fileHeader, _ := r.FormFile("avatar")
	var ext = strings.Split(fileHeader.Filename, ".")[1]
	var targetFileName = IDUser + "." + ext
	var targetPath = "uploads/avatars/" + targetFileName

	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error al subir la imagen"+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		http.Error(w, "Error al copia la imagen"+err.Error(), http.StatusBadRequest)
		return
	}
	var user models.User
	user.Avatar = targetFileName
	_, err = db.UpdateUser(user, IDUser)

	if err != nil {
		http.Error(w, "Error al gravar la imagen en la BD"+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func UploadBanner(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10000 * time.Millisecond)
	sourceFile, fileHeader, err := r.FormFile("banner")
	if err != nil {
		http.Error(w, "Error al enviar imagen al servidor"+err.Error(), http.StatusBadRequest)
		return
	}

	var ext = strings.Split(fileHeader.Filename, ".")[1]
	var targetFileName = IDUser + "." + ext
	var targetPath = "uploads/banners/" + targetFileName

	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error al subir la imagen"+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		http.Error(w, "Error al copia la imagen"+err.Error(), http.StatusBadRequest)
		return
	}
	var user models.User
	user.Banner = targetFileName
	_, err = db.UpdateUser(user, IDUser)

	if err != nil {
		http.Error(w, "Error al gravar la imagen en la BD"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func GetAvatar(w http.ResponseWriter, r *http.Request) {

	userId := r.URL.Query().Get("id")
	if len(userId) == 0 {
		http.Error(w, "Id del usuario requerido", http.StatusBadRequest)
		return
	}

	user, err := db.SearchUserByID(userId)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusBadRequest)
		return
	}

	file, err := os.Open("uploads/avatars/" + user.Avatar)
	if err != nil {
		http.Error(w, "Imagen no encontrada", http.StatusBadRequest)
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error al copiar la imagen", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetBanner(w http.ResponseWriter, r *http.Request) {

	userId := r.URL.Query().Get("id")
	if len(userId) == 0 {
		http.Error(w, "Id del usuario requerido", http.StatusBadRequest)
		return
	}

	user, err := db.SearchUserByID(userId)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusBadRequest)
		return
	}

	file, err := os.Open("uploads/banners/" + user.Avatar)
	if err != nil {
		http.Error(w, "Imagen no encontrada", http.StatusBadRequest)
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error al copiar la imagen", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
