package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/naeem4265/Catering-Management/auth"
	"github.com/naeem4265/Catering-Management/data"
	"github.com/naeem4265/Catering-Management/restaurants"
	"github.com/naeem4265/Catering-Management/users"
)

func main() {

	fmt.Println("Programm started")
	// create a database object which can be used to connect with the database.
	db, err := sql.Open("mysql", "root:1234@tcp(db:3306)/catering_management")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Now its  time to connect with oru database, database object has a method Ping.
	// Ping returns error, if unable connect to database.
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// Create Table if not exists.
	if err := data.CreateTables(db); err != nil {
		log.Fatal("Database Creating error:", err)
	}
	fmt.Print("Database Connected\n")

	router := chi.NewRouter()

	// Use a closure to capture the 'db' object and pass it to the handler.
	router.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		auth.SignIn(w, r, db)
	})
	router.Get("/signout", auth.SignOut)

	// Get all users
	router.Route("/users", func(r chi.Router) {
		// Create user account
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			users.CreateUser(w, r, db)
		})
		// Get user list
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			users.GetUsers(w, r, db)
		})
	})

	// Add restaurant
	router.Route("/restaurant", func(r chi.Router) {
		r.Use(authentication)
		// Add restaurant
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			restaurants.AddRestaurant(w, r, db)
		})
	})

	// Get all Menus
	router.Route("/menu", func(r chi.Router) {
		r.Use(authentication)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			restaurants.GetMenu(w, r, db)
		})
		// Add Menu for the restaurant
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			restaurants.AddMenu(w, r, db)
		})
		// Vote for menu id
		r.Put("/vote/{id}", func(w http.ResponseWriter, r *http.Request) {
			restaurants.VoteItem(w, r, db)
		})
		// Getting results for the previous day.
		r.Get("/winner", func(w http.ResponseWriter, r *http.Request) {
			restaurants.GetWinner(w, r, db)
		})

		// Getting results for daily winner.
		r.Get("/result", func(w http.ResponseWriter, r *http.Request) {
			restaurants.GetResult(w, r, db)
		})
		// Confirm today's menu and reset vote
		r.Post("/confirm", func(w http.ResponseWriter, r *http.Request) {
			restaurants.ConfirmMenu(w, r, db)
		})
	})

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Check user authentications
func authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for the "token" cookie
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value

		claims := &auth.Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JWTKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				// Token signature is invalid, return unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other error while parsing claims, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			// Token is not valid, return unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// If token is valid, continue to the next handler
		next.ServeHTTP(w, r)
	})
}
