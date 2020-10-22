package main

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
)

// Login model
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Signup model
type Signup struct {
	Name string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
}

func userExists(username string, db *sql.DB) bool {
	var userQuery string
	checkIfUserExists, err := db.Prepare("SELECT EXISTS(SELECT 1 FROM auth WHERE username=?)")

	defer checkIfUserExists.Close()

	row := checkIfUserExists.QueryRow(username)
	err = row.Scan(&userQuery)
	if err != nil {
		log.Fatal(err)
	}
	userExists, err := strconv.ParseBool(userQuery)
	if err != nil {
		log.Fatal(err)
	}
	return userExists
}

// Login accepts username and password and does an authentication check on the users credentials.
func (l *Login) Login(username, password string, db *sql.DB) error {
	var pwd string
	// Look for the user
	userPwd, err := db.Prepare("SELECT password FROM auth WHERE username=?")
	if err != nil {
		log.Fatal(err)
	}

	// Close prepared statements.
	defer userPwd.Close()

	if userExists(username, db) {

		row := userPwd.QueryRow(username)
		err := row.Scan(&pwd)
		if err != nil {
			log.Fatal(err)
		}

		// Validate user credentials
		switch ComparePassword(pwd, password) {
		case true:
			return nil
		case false:
			return errors.New("invalid credentials")
		}
	}
	return errors.New("invalid credentials")
}

// Signup accepts username, password and an email and registers an creates a new account.
func (s *Signup) Signup(username, password string, db *sql.DB) error {

	// Add a new user to db
	addNewAuth, err := db.Prepare("INSERT INTO auth(username, password) values(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer addNewAuth.Close()

	if userExists(username, db) {
		return errors.New("user already exits")
	} else {
		_, err := addNewAuth.Exec(username, HashAndSalt(password))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}
}
