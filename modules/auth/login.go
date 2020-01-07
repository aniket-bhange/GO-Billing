package auth

import (
	"billing-gorilla/core"
	config "billing-gorilla/core"
	database "billing-gorilla/db"
	"billing-gorilla/model"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

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

	token, err := SignIn(user.Email, user.Password)

	if err != nil {
		err = errors.New("Incorrect details")
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	config.RespondJSON(w, http.StatusOK, token)

}

func SignIn(email, password string) (string, error) {
	var err error

	user := model.Users{}

	db := database.ConnectDB()

	err = db.Db.Debug().Model(model.Users{}).Where("email = ?", email).Take(&user).Error

	log.Print(user.Password)

	if err != nil {
		return "", err
	}

	isVerified := user.VerifyPassword(user.Password, password)

	if !isVerified {
		return "", errors.New("Not matched")
	}
	return core.CreateToken(user.ID)
}
