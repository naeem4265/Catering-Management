package users

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/naeem4265/Catering-Management/data"
)

func GetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var users []data.Credential

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var temp data.Credential
		if err := rows.Scan(&temp.Username); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		users = append(users, temp)
	}
	userlist, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(userlist)
	w.WriteHeader(http.StatusOK)
}
