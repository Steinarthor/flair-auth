package main

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
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

func (r *Response) generateToken() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	JwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	jwtSignedKeyBytes := []byte(JwtSigningKey)

	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "flair-auth",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtSignedKeyBytes)

	r.Token = ss
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

	err, jwtToken := GenerateJWT(loginRequest.Username)
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)
	err = dec.Decode(&loginRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	err = loginRequest.Login(loginRequest.Username, loginRequest.Password, a.Db)

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

	err, jwtToken := GenerateJWT(signupRequest.Username)
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)
	err = dec.Decode(&signupRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	err = signupRequest.Signup(signupRequest.Username, signupRequest.Password, a.Db)

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
