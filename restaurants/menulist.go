package restaurants

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/naeem4265/Catering-Management/data"
)

func GetMenu(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	rows, err := db.Query("SELECT m.Id, m.Name, r.Name, r.Location, m.Price, m.Vote FROM menu m INNER JOIN restaurant r ON m.RestId = r.Id ORDER BY m.Vote DESC, r.Name ASC")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// An user slice to hold data from returned rows.
	var items []data.Item

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var temp data.Item
		if err := rows.Scan(&temp.Id, &temp.Name, &temp.ResName, &temp.Location, &temp.Price, &temp.Vote); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		items = append(items, temp)
	}
	itemlist, err := json.Marshal(items)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(itemlist)
	w.WriteHeader(http.StatusOK)
}
