package auth

import (
	"net/http"
	"time"
)

func SignOut(w http.ResponseWriter, r *http.Request) {
	// Create a new cookie with an expired expiration time
	expiredCookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour), // Expire the cookie in the past
	}
	// Set the expired cookie to effectively remove the token
	http.SetCookie(w, expiredCookie)
	w.WriteHeader(http.StatusOK)
	// You might want to redirect the user to a different page after logout
	http.Redirect(w, r, "/signIn", http.StatusSeeOther)
}
