package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Pokemon represents a Pokemon entity
type Pokemon struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	HP        int    `json:"hp"`
	Attack    int    `json:"attack"`
	Defense   int    `json:"defense"`
	SpAttack  int    `json:"sp_attack"`
	SpDefense int    `json:"sp_defense"`
	Speed     int    `json:"speed"`
}

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Configure DB connection
	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT"),
		DBName:               os.Getenv("MYSQL_DATABASE"),
		AllowNativePasswords: true,
		Params: map[string]string{
			"tls": "false",
		},
		ParseTime: true,
	}

	// Initialize DB connection
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to RDS MySQL database!")

	// Initialize database tables
	initDB()

	// Setup API routes
	router := mux.NewRouter()
	router.HandleFunc("/pokemon", getAllPokemon).Methods("GET")
	router.HandleFunc("/pokemon/{id}", getPokemon).Methods("GET")
	router.HandleFunc("/pokemon", createPokemon).Methods("POST")
	router.HandleFunc("/pokemon/{id}", updatePokemon).Methods("PUT")
	router.HandleFunc("/pokemon/{id}", deletePokemon).Methods("DELETE")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(router)))
}

// Initialize database tables
func initDB() {
	query := `CREATE TABLE IF NOT EXISTS pokemon (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        type VARCHAR(50) NOT NULL,
        hp INT NOT NULL,
        attack INT NOT NULL,
        defense INT NOT NULL,
        sp_attack INT NOT NULL,
        sp_defense INT NOT NULL,
        speed INT NOT NULL
    )`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized!")
}

// CRUD Operations

// GET all Pokemon
func getAllPokemon(w http.ResponseWriter, r *http.Request) {
	var pokemons []Pokemon

	rows, err := db.Query("SELECT id, name, type, hp, attack, defense, sp_attack, sp_defense, speed FROM pokemon")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Pokemon
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.HP, &p.Attack, &p.Defense, &p.SpAttack, &p.SpDefense, &p.Speed); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		pokemons = append(pokemons, p)
	}

	respondWithJSON(w, http.StatusOK, pokemons)
}

// GET a specific Pokemon
func getPokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Pokemon ID")
		return
	}

	var p Pokemon
	query := "SELECT id, name, type, hp, attack, defense, sp_attack, sp_defense, speed FROM pokemon WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Type, &p.HP, &p.Attack, &p.Defense, &p.SpAttack, &p.SpDefense, &p.Speed)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Pokemon not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// CREATE a new Pokemon
func createPokemon(w http.ResponseWriter, r *http.Request) {
	var p Pokemon
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	query := `INSERT INTO pokemon 
              (name, type, hp, attack, defense, sp_attack, sp_defense, speed) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, p.Name, p.Type, p.HP, p.Attack, p.Defense, p.SpAttack, p.SpDefense, p.Speed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	p.ID = int(id)
	respondWithJSON(w, http.StatusCreated, p)
}

// UPDATE an existing Pokemon
func updatePokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Pokemon ID")
		return
	}

	var p Pokemon
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.ID = id

	query := `UPDATE pokemon 
              SET name = ?, type = ?, hp = ?, attack = ?, defense = ?, 
                  sp_attack = ?, sp_defense = ?, speed = ? 
              WHERE id = ?`
	_, err = db.Exec(query, p.Name, p.Type, p.HP, p.Attack, p.Defense, p.SpAttack, p.SpDefense, p.Speed, p.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// DELETE a Pokemon
func deletePokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Pokemon ID")
		return
	}

	query := "DELETE FROM pokemon WHERE id = ?"
	_, err = db.Exec(query, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Helper functions for HTTP responses
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
