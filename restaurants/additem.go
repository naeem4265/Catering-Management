package restaurants

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/naeem4265/Catering-Management/data"
)

func AddMenu(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var menu data.Menu
	err := json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var Id int32
	err = db.QueryRow("SELECT Id FROM restaurant WHERE Id = ?", menu.RestId).Scan(&Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	menu.Vote = 0
	_, err = db.Exec("INSERT INTO menu (Id, RestId, Name, Price, Vote) VALUES(?, ?, ?, ?, ?)", menu.Id, menu.RestId, menu.Name, menu.Price, menu.Vote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
