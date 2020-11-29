package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

// Response type send to users.
type Response struct {
	Status  int
	Token   string
	Message string
}

// Initialize accepts database credentials and initializes the app.
func (a *App) Initialize(database string) {
	db, err := sql.Open(database, "./flair.db")
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.Db = db
	a.initializeRoutes()
}

// Run accepts an address and starts the application.
func (a *App) Run(addr string) {
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server is running and listening on port: 8080")
	log.Fatal(srv.ListenAndServe())
}

// Route handlers
func (a *App) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginRequest Login

	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)
	err := dec.Decode(&loginRequest)

	err, jwtToken := GenerateJWT(loginRequest.Email)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	err = loginRequest.Login(loginRequest.Email, loginRequest.Password, a.Db)

	if err != nil {
		enc.Encode(Response{
			Status:  http.StatusConflict,
			Token:   "",
			Message: err.Error(),
		})

	} else {
		enc.Encode(Response{
			Status:  http.StatusCreated,
			Token:   jwtToken,
			Message: "Success.",
		})
	}
}

func (a *App) signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var signupRequest Signup

	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)
	err := dec.Decode(&signupRequest)

	err, jwtToken := GenerateJWT(signupRequest.Email)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	err = signupRequest.Signup(signupRequest.Email, signupRequest.Name, signupRequest.Password, a.Db)

	if err != nil {
		enc.Encode(Response{
			Status:  http.StatusConflict,
			Token:   "",
			Message: err.Error(),
		})

	} else {
		enc.Encode(Response{
			Status:  http.StatusCreated,
			Token:   jwtToken,
			Message: "Success.",
		})
	}
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/login", a.login).Methods("POST")
	a.Router.HandleFunc("/signup", a.signup).Methods("POST")
}
