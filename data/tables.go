package data

import (
	"database/sql"
)

func CreateTables(db *sql.DB) error {
	// Users table to store all the users
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (Username varchar(255) PRIMARY KEY, Password varchar(255) NOT NULL)")
	if err != nil {
		return nil
	}
	// Restaurant table to store all the restaurants name
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS restaurant (Id INT AUTO_INCREMENT PRIMARY KEY, Name varchar(255) NOT NULL, Location varchar(255))")
	if err != nil {
		return nil
	}
	// menu table to store all the menues
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS menu (Id INT AUTO_INCREMENT PRIMARY KEY, RestId INT NOT NULL, Name varchar(255) NOT NULL, Price INT NOT NULL, Vote INT)")
	if err != nil {
		return nil
	}
	// Daily_winner table to store daily record
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS daily_winner (Date date, WinnerMenuId INT NOT NULL, WinnerRestaurantId INT NOT NULL)")
	if err != nil {
		return nil
	}

	// all the databases created successfully.
	return nil
}
