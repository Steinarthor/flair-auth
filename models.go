package main

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
)

// Login model
type Login struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

// Signup model
type Signup struct {
	Name string `json:"name"`
	Password string `json:"password"`
	Email string `json:"email"`
}

func userExists(email string, db *sql.DB) bool {
	var userQuery string
	checkIfUserExists, err := db.Prepare("SELECT EXISTS(SELECT 1 FROM authentication WHERE email=?)")

	defer checkIfUserExists.Close()

	row := checkIfUserExists.QueryRow(email)
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

// Login accepts email and password and does an authentication check on the users credentials.
func (l *Login) Login(email, password string, db *sql.DB) error {
	var pwd string
	// Look for the user
	userPwd, err := db.Prepare("SELECT password FROM authentication WHERE email=?")
	if err != nil {
		log.Fatal(err)
	}

	// Close prepared statements.
	defer userPwd.Close()

	if userExists(email, db) {

		row := userPwd.QueryRow(email)
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
func (s *Signup) Signup(email, name, password string, db *sql.DB) error {
	// Add a new user to db
	addNewAuth, err := db.Prepare("INSERT INTO authentication(email, name, password) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer addNewAuth.Close()

	if userExists(email, db) {
		return errors.New("user already exits")
	} else {
		_, err := addNewAuth.Exec(email, name, HashAndSalt(password))
		if err != nil {
			log.Fatal(err.Error())
		}
		return nil
	}
}
