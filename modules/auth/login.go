package auth

import (
	"billing-gorilla/core"
	config "billing-gorilla/core"
	database "billing-gorilla/db"
	"billing-gorilla/model"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func Login(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := model.Users{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	authtoken, err := SignIn(user.Email, user.Password)

	if err != nil {
		err = errors.New("Incorrect details")
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	response := LoginResponse{
		Token:   authtoken,
		Message: "Login is successful",
	}
	config.RespondJSON(w, http.StatusOK, response)

}

func SignIn(email, password string) (string, error) {
	var err error

	user := model.Users{}
	db := database.ConnectDB()

	err = db.Db.Debug().Model(model.Users{}).Where("email = ?", email).Take(&user).Error

	if err != nil {
		return "", err
	}

	isVerified := user.VerifyPassword(user.Password, password)

	if !isVerified {
		return "", errors.New("Not matched")
	}
	return core.CreateToken(user.ID)
}
