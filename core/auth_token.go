package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//CreateToken to create tojen with userId
func CreateToken(userID uint) (string, error) {
	climas := jwt.MapClaims{}
	climas["authorized"] = true
	climas["user_id"] = userID
	climas["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, climas)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

//ExtractToken to fetch token from header
func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()

	token := keys.Get("token")

	if token != "" {
		return token
	}

	barrerToken := r.Header.Get("Authorization")
	if len(strings.Split(barrerToken, " ")) == 2 {
		return strings.Split(barrerToken, " ")[1]
	}

	return ""

}

//Pretty display the claims licely in the terminal
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}

//ParseToken return token value from token string
func ParseToken(tokenStr string) (*jwt.Token, error) {

	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected Signing method %v", token.Header["alg"])
		}

		return []byte(os.Getenv("API_SECRET")), nil
	})
}

//ValidateToken will valid token
func ValidateToken(r *http.Request) error {
	tokenStr := ExtractToken(r)

	token, err := ParseToken(tokenStr)

	if err != nil {
		return err
	}

	if climas, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(climas)
	}
	return nil

}

//ExtractTokenID used to extract userID from token
func ExtractTokenID(r *http.Request) (uint32, error) {
	tokenStr := ExtractToken(r)

	token, err := ParseToken(tokenStr)

	if err != nil {
		return 0, err
	}

	climas, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", climas["user_id"]), 10, 32)

		if err != nil {
			return 0, err
		}

		return uint32(uid), nil
	}

	return 0, nil

}
