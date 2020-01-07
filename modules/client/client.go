package client

import (
	config "billing-gorilla/core"
	database "billing-gorilla/db"
	"billing-gorilla/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clients := model.Client{}

	db := database.ConnectDB()
	allclients, err := clients.FindAll(db.Db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.RespondJSON(w, 200, allclients)

}

func CreateClient(w http.ResponseWriter, r *http.Request) {
	client := model.Client{}

	err := json.NewDecoder(r.Body).Decode(&client)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()
	newClient, err := client.Save(db.Db)

	if err != nil {
		log.Fatalf("Error for creating client: %v \n", err)
		config.RespondError(w, 400, err)
		return
	}

	config.RespondJSON(w, 200, newClient)

}

func Update(w http.ResponseWriter, r *http.Request) {
	client := model.Client{}

	vars := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&client)

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
	updatedClient, err := client.Update(db.Db, uid)

	if err != nil {

		config.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	config.RespondJSON(w, 200, updatedClient)
}
