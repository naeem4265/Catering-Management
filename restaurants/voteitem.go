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
func GetWinner(w http.ResponseWriter, r *http.Request, db *sql.DB) *data.Item {
	query := "SELECT m.Id, m.Name, m.RestId, r.Name, r.Location, m.Price, m.Vote " +
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
		return nil
	}
	defer rows.Close()

	// An item slice to hold data from returned rows.
	var item data.Item

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		if err := rows.Scan(&item.Id, &item.Name, &item.ResId, &item.ResName, &item.Location, &item.Price, &item.Vote); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}
	}

	// Check if any items were found
	if item.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	itemlist, err := json.Marshal(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	w.Write(itemlist)
	w.WriteHeader(http.StatusOK)
	return &item
}

// daily winner list
func GetResult(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := "SELECT WinnerRestaurantId, WinnerMenuId, Date FROM daily_winner ORDER BY Date DESC"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("query fail")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer rows.Close()

	// An user slice to hold data from returned rows.
	var dailylist []data.Daily_Winner

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var temp data.Daily_Winner
		if err := rows.Scan(&temp.WinnerRestaurantId, &temp.WinnerManuId, &temp.Date); err != nil {
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

// Confirm today's menu
func ConfirmMenu(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Today's winner
	winner := GetWinner(w, r, db)
	if winner == nil {
		fmt.Println("\nGot error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var id = -1
	err := db.QueryRow("SELECT WinnerMenuId FROM daily_winner WHERE Date = CURDATE()").Scan(&id)
	if err == nil {
		fmt.Println("\nAlready confirmed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// insert into daily_winner table
	_, err = db.Exec("INSERT INTO daily_winner (Date, WinnerMenuId, WinnerRestaurantId) VALUES(CURDATE(), ?, ?)", winner.Id, winner.ResId)
	if err != nil {
		fmt.Println("\nInsertion error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("UPDATE menu SET Vote = 0")
	if err != nil {
		fmt.Println("\nVote reset error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("\nConfirmed")
	w.WriteHeader(http.StatusOK)
}
