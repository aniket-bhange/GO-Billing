package users

import (
	config "billing-gorilla/core"
	database "billing-gorilla/db"
	"billing-gorilla/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("This is user print")
}

func Create(w http.ResponseWriter, r *http.Request) {

	user := model.Users{}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	newUser, err := user.SaveUser(db.Db)
	if err != nil {
		log.Fatalf("Error while creating user: %v \n", err)
		config.RespondError(w, 400, err)
		return
	}
	config.RespondJSON(w, 200, newUser)

}

func Update(w http.ResponseWriter, r *http.Request) {

	user := model.Users{}

	vars := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := database.ConnectDB()

	upatedUser, err := user.UpdateUser(db.Db, uid)

	if err != nil {

		config.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	config.RespondJSON(w, 200, upatedUser)

}
