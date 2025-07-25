package authcontroller

import (
	"encoding/json"
	"eticket/models"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("secret"))

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		http.Error(w, "Gagal decode input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var user models.User
	if err := models.DB.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
		http.Error(w, "Email tidak ditemukan", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		http.Error(w, "Password salah", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "session-id")
	session.Values["user_id"] = user.Id
	session.Save(r, w)

	response, _ := json.Marshal(map[string]string{"message": "Login berhasil"})
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func Register(w http.ResponseWriter, r *http.Request) {

	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		log.Fatal("Gagal mendecode JSON")
	}
	defer r.Body.Close()

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	if err := models.DB.Create(&userInput).Error; err != nil {
		log.Fatal("Gagal menyimpan data")
	}

	response, _ := json.Marshal(map[string]string{"message": "User berhasil terdaftar"})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	session.Options.MaxAge = -1
	session.Save(r, w)

	response, _ := json.Marshal(map[string]string{"message": "Logout berhasil"})
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}
