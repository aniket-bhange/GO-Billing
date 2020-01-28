package product

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

	product := model.Product{}

	db := database.ConnectDB()

	products, err := product.FindAll(db.Db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.RespondJSON(w, 200, products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	product := model.Product{}

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.ConnectDB()

	newProduct, err := product.Create(db.Db)

	if err != nil {
		log.Fatalf("Error while create product: %v \n", err)
		config.RespondError(w, 400, err)
		return
	}

	config.RespondJSON(w, 200, newProduct)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	product := model.Product{}

	vars := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&product)

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

	updatedProduct, err := product.Update(db.Db, uid)

	if err != nil {
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	config.RespondJSON(w, 200, updatedProduct)

}
