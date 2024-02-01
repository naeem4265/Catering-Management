package restaurants

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/naeem4265/Catering-Management/data"
)

func VoteItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")

	// Increase vote by 1 for id
	result, err := db.Exec("UPDATE menu SET Vote = Vote+1 WHERE Id = ?", id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Winner for today
func GetWinner(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := "SELECT m.Id, m.Name, r.Name, r.Location, m.Price, m.Vote " +
		"FROM menu m " +
		"INNER JOIN restaurant r ON m.RestId = r.Id " +
		"HAVING 2 > ( " +
		"    SELECT COUNT(d.WinnerRestaurantId) " +
		"    FROM daily_winner d " +
		"    WHERE d.WinnerRestaurantId = m.RestId AND d.Date BETWEEN CURDATE() - INTERVAL 2 DAY AND CURDATE() - INTERVAL 1 DAY" +
		") " +
		"ORDER BY m.Vote DESC LIMIT 1"

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer rows.Close()

	// An user slice to hold data from returned rows.
	var items []data.Item

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var temp data.Item
		if err := rows.Scan(&temp.Id, &temp.Name, &temp.ResName, &temp.Location, &temp.Price, &temp.Vote); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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

// List of the winner list
func GetResult(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT m.Id, m.Name, r.Name, r.Location, m.Price, m.Vote, d.date FROM daily_winner d INNER JOIN menu m ON d.winner_id = m.Id INNER JOIN restaurant r ON m.RestId = r.Id ORDER BY d.date DESC")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer rows.Close()

	// An user slice to hold data from returned rows.
	var dailylist []data.Item

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var temp data.Item
		if err := rows.Scan(&temp.Id, &temp.Name, &temp.ResName, &temp.Location, &temp.Price, &temp.Vote, &temp.Date); err != nil {
			fmt.Println("query fail")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dailylist = append(dailylist, temp)
	}
	dailylists, err := json.Marshal(dailylist)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(dailylists)
	w.WriteHeader(http.StatusOK)
}
