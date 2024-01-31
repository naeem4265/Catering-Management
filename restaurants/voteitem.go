package restaurants

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
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
