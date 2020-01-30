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

type LoginResponse struct {
	Token        string              `json:"token"`
	Message      string              `json:"message"`
	UserResponse *model.UserResponse `json:"userinfo"`
}

type UserBasicInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        uint   `json:"user_id"`
	Client    model.ClientResponse
}

func Login(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := model.Users{}
	//var response model.UserResponse

	err = json.Unmarshal(body, &user)

	if err != nil {
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	authtoken, err, response := SignIn(user.Email, user.Password)

	if err != nil {
		err = errors.New("Incorrect details")
		config.RespondError(w, http.StatusUnprocessableEntity, err)
		return
	}

	responsebody := LoginResponse{
		Token:        authtoken,
		Message:      "Login is successful",
		UserResponse: response,
	}
	config.RespondJSON(w, http.StatusOK, responsebody)

}

func SignIn(email, password string) (string, error, *model.UserResponse) {
	var err error

	user := model.Users{}
	response := model.UserResponse{}
	db := database.ConnectDB()

	// err = db.Db.Debug().Model(model.Users{}).Where("email = ?", email).Take(&user).Error

	result, err := user.FindOneByField(db.Db, "email", email)

	log.Printf("result %v", result)

	if err != nil {
		return "", err, &response
	}

	isVerified := user.VerifyPassword(result.User.Password, password)

	if !isVerified {
		return "", errors.New("Not matched"), &response
	}

	token, err := core.CreateToken(user.ID)

	return token, err, result
}
