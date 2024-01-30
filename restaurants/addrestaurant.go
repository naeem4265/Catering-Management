package restaurants

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/naeem4265/Catering-Management/data"
)

func AddRestaurant(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var restaurant data.Restaurant
	if err := json.NewDecoder(r.Body).Decode(&restaurant); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT INTO restaurant (id, name, location) VALUES(?, ?, ?)", restaurant.Id, restaurant.Name, restaurant.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
