package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/naeem4265/Catering-Management/data"
)

func CreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fmt.Println("Here")
	var user data.Credential
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT INTO users (username, password) VALUES(?, ?)", user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
