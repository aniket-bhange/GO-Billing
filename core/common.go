package core

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func RespondError(w http.ResponseWriter, code int, err error) {
	RespondJSON(w, code, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}

func FileUpload(r *http.Request, fileName string) (string, error) {
	file, header, err := r.FormFile(fileName)

	if err != nil {
		return "", err
	}
	defer file.Close()

	var extension = filepath.Ext(header.Filename)
	var path = "/uploads/" + time.Now().Format("MM-DD-YYYY") + "." + extension

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	io.Copy(out, file)

	return path, nil

}
