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
	"github.com/naeem4265/Catering-Management/restaurants"
	"github.com/naeem4265/Catering-Management/users"
)

func main() {

	// create a database object which can be used to connect with database.
	db, err := sql.Open("mysql", "root:@tcp(0.0.0.0:3306)/catering_management")
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
	fmt.Print("Database Connected\n")

	router := chi.NewRouter()

	// Use a closure to capture the 'db' object and pass it to the handler.
	router.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		auth.SignIn(w, r, db)
	})
	router.Get("/signout", auth.SignOut)
	// Create user account
	router.Post("/createuser", func(w http.ResponseWriter, r *http.Request) {
		users.CreateUser(w, r, db)
	})
	// Get all users
	router.Route("/users", func(r chi.Router) {
		r.Use(authentication)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			users.GetUsers(w, r, db)
		})
	})

	// Add restaurant
	router.Post("/addrest", func(w http.ResponseWriter, r *http.Request) {
		restaurants.AddRestaurant(w, r, db)
	})

	// Add Menu for the restaurant
	router.Post("/addmenu", func(w http.ResponseWriter, r *http.Request) {
		restaurants.AddMenu(w, r, db)
	})

	// Get all Menus
	router.Route("/menu", func(r chi.Router) {
		r.Use(authentication)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			restaurants.GetMenu(w, r, db)
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
