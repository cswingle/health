package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// Health holds health data
type Health struct {
	ID       *int64     `db:"id,omitempty" json:"id,omitempty"`
	Username *string    `db:"username,omitempty" json:"username,omitempty"`
	Ts       *time.Time `db:"ts,omitempty" json:"ts,omitempty"`
	Variable *string    `db:"variable,omitempty" json:"variable,omitempty"`
	Value    *float64   `db:"value,omitempty" json:"value,omitempty"`
}

// RefVariables holds ref_variables data
type RefVariables struct {
	ID       *int64  `db:"id,omitempty" json:"id,omitempty"`
	Variable *string `db:"variable,omitempty" json:"variable,omitempty"`
	Units    *string `db:"units,omitempty" json:"units,omitempty"`
	Sequence *int64  `db:"sequence,omitempty" json:"sequence,omitempty"`
}

// User holds a users account information
type User struct {
	Username    string `db:"username"`
	AccessLevel string `db:"access_level"`
}

// store will hold all session data
var store *pgstore.PGStore

var authKeyOne []byte

var db *sqlx.DB

var dbhost string
var dbport string
var dbuser string
var dbname string
var dbpassword string

// login authenticates the user
func login(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user := User{}
	query := fmt.Sprintf(
		`SELECT username, access_level FROM users WHERE username='%s'
		AND password = crypt('%s', password)`, username, password)
	db.Get(&user, query)

	if user.Username != "" {
		session.Values["user"] = user

		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

// logout revokes authentication for a user
func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = User{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logRequest(r)
}

// GetLoggedIn returns User if logged in
func GetLoggedIn(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	var err error

	session, err := store.Get(r, "auth")
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode("access denied"); err != nil {
			panic(err)
		}
	} else {
		// Convert our session data into an instance of User
		user := User{}
		user, _ = session.Values["user"].(User)

		if user.Username != "" && user.AccessLevel == "admin" {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(user); err != nil {
				panic(err)
			}
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode("access denied"); err != nil {
				panic(err)
			}
		}
	}

	logRequest(r)
}

// GetHealth returns data from the health table if logged in
func GetHealth(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	params := mux.Vars(r)

	health := []Health{}

	var err error

	session, err := store.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert our session data into an instance of User
	user := User{}
	user, _ = session.Values["user"].(User)

	if user.Username != "" && user.AccessLevel == "admin" {
		if _, ok := params["id"]; ok {
			err = db.Select(&health, "SELECT id, username, ts, variable, value "+
				"FROM public.health "+
				"WHERE id = $1 ", params["id"])
		} else if _, ok = params["ts"]; ok {
			err = db.Select(&health, "SELECT id, username, ts, variable, value "+
				"FROM public.health "+
				"WHERE ts = $1 ", params["ts"])
		} else if _, ok = params["variable"]; ok {
			err = db.Select(&health, "SELECT id, username, ts, variable, value "+
				"FROM public.health "+
				"WHERE variable = $1 ", params["variable"])
		} else {
			err = db.Select(&health, "SELECT id, username, ts, variable, value "+
				"FROM public.health ")
		}
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(health); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode("access denied"); err != nil {
			panic(err)
		}
	}

	logRequest(r)
}

// GetRefVariables returns data from the ref_variables table
func GetRefVariables(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	params := mux.Vars(r)

	refVariables := []RefVariables{}

	var err error

	if _, ok := params["id"]; ok {
		err = db.Select(&refVariables, "SELECT id, variable, units, sequence "+
			"FROM public.ref_variables "+
			"WHERE id = $1 ", params["id"])
	} else if _, ok = params["variable"]; ok {
		err = db.Select(&refVariables, "SELECT id, variable, units, sequence "+
			"FROM public.ref_variables "+
			"WHERE variable = $1 ", params["variable"])
	} else {
		err = db.Select(&refVariables, "SELECT id, variable, units, sequence "+
			"FROM public.ref_variables ")
	}
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(refVariables); err != nil {
		panic(err)
	}

	logRequest(r)
}

// PostHealth inserts data into the health table
func PostHealth(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	params := mux.Vars(r)
	returnMessages := make(map[string][]string)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	session, err := store.Get(r, "auth")
	if err != nil {
		returnMessages["message"] = append(returnMessages["message"], err.Error())
		returnMessages["status"] = append(returnMessages["status"], "error")
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert our session data into an instance of User
	user := User{}
	user, _ = session.Values["user"].(User)

	if user.Username != "" && user.AccessLevel == "admin" {
		var health []Health

		if err := json.Unmarshal(body, &health); err != nil {
			returnMessages["message"] = append(returnMessages["message"], err.Error())
			returnMessages["status"] = append(returnMessages["status"], "error")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			log.Println(returnMessages)
			if err := json.NewEncoder(w).Encode(returnMessages); err != nil {
				panic(err)
			}
			return
		}

		if len(health) == 0 {
			returnMessages["message"] = append(returnMessages["message"], "Warning: no items in health")
			returnMessages["status"] = append(returnMessages["status"], "warning")
		} else {
			if params["keys"] == "u" {
				for _, item := range health {
					_, err := db.Exec("INSERT INTO public.health ("+
						"username, "+
						"ts, "+
						"variable, "+
						"value"+
						") VALUES ("+
						" $1,"+
						" $2,"+
						" $3,"+
						" $4"+
						") "+"ON CONFLICT ("+
						"ts, "+
						"variable "+
						") DO "+
						"UPDATE SET "+
						" username = $1,"+
						" value = $4",
						user.Username,
						item.Ts,
						item.Variable,
						item.Value)
					if err != nil {
						returnMessages["message"] = append(returnMessages["message"], err.Error())
						returnMessages["status"] = append(returnMessages["status"], "error")
						fmt.Println(err)
					} else {
						returnMessages["message"] = append(returnMessages["message"], "inserted item into health")
						returnMessages["status"] = append(returnMessages["status"], "info")
					}
				}
			}
			if params["keys"] == "p" {
				for _, item := range health {
					_, err := db.Exec("INSERT INTO public.health ("+
						"id, "+
						"username, "+
						"ts, "+
						"variable, "+
						"value"+
						") VALUES ("+
						" $1,"+
						" $2,"+
						" $3,"+
						" $4,"+
						" $5"+
						") "+
						"ON CONFLICT ("+
						"id "+
						") DO "+
						"UPDATE SET "+
						" username = $2,"+
						" ts = $3,"+
						" variable = $4,"+
						" value = $5",
						item.ID,
						user.Username,
						item.Ts,
						item.Variable,
						item.Value)
					if err != nil {
						returnMessages["message"] = append(returnMessages["message"], err.Error())
						returnMessages["status"] = append(returnMessages["status"], "error")
						fmt.Println(err)
					} else {
						returnMessages["message"] = append(returnMessages["message"], "inserted item into health")
						returnMessages["status"] = append(returnMessages["status"], "info")
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		log.Println(returnMessages)
		if err := json.NewEncoder(w).Encode(returnMessages); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode("access denied"); err != nil {
			panic(err)
		}
	}

	logRequest(r)
}

// logRequest logs the request
func logRequest(r *http.Request) {
	forward := r.Header.Get("X-Forwarded-For")
	userAgent := r.UserAgent()

	log.Println(r.Host, r.Method, r.URL, r.Proto, forward, userAgent)
}

// GetHelp returns a json help statement
func GetHelp(w http.ResponseWriter, r *http.Request) {
	options := map[string]string{
		"GET /health/v1":                                     "This help message",
		"GET /health/v1/health":                              "Get data from health table",
		"GET /health/v1/health/id/{ id }":                    "Get health data by id",
		"GET /health/v1/health/ts/{ ts }":                    "Get health data by ts",
		"GET /health/v1/health/variable/{ variable }":        "Get health data by variable",
		"GET /health/v1/logged_in":                           "Get login status (username)",
		"GET /health/v1/ref_variables":                       "Get data from ref_variables table",
		"GET /health/v1/ref_variables/id/{ id }":             "Get ref_variables data by id",
		"GET /health/v1/ref_variables/variable/{ variable }": "Get ref_variables data by variable",
		"POST /health/v1/login":                              "Log in",
		"POST /health/v1/logout":                             "Log out",
		"POST /health/v1/health/keys/u":                      "Post data to health table using unique key for insert/update",
		"POST /health/v1/health/keys/p":                      "Post data to health table using primary key for insert/update",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(options); err != nil {
		panic(err)
	}

	logRequest(r)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbhost = os.Getenv("DB_HOST")
	dbport = os.Getenv("DB_PORT")
	dbname = os.Getenv("DB_NAME")
	dbuser = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASSWORD")

	authKeyOne = securecookie.GenerateRandomKey(64)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s "+
		"sslmode=require password=%s",
		dbhost, dbport, dbuser, dbname, dbpassword)

	// Register User struct so it can be encoded/decoded to/from the Session data
	gob.Register(User{})

	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error

	dsn := fmt.Sprintf("postgres://%s:%s/%s?sslmode=require",
		dbhost, dbport, dbname)
	store, err = pgstore.NewPGStore(dsn, []byte(authKeyOne))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer store.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/health/v1", GetHelp).Methods("GET")

	router.HandleFunc("/health/v1/login", login)
	router.HandleFunc("/health/v1/logout", logout)

	router.HandleFunc("/health/v1/health", func(w http.ResponseWriter, r *http.Request) {
		GetHealth(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/health/v1/health/id/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetHealth(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/health/v1/health/ts/{ts}", func(w http.ResponseWriter, r *http.Request) {
		GetHealth(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/health/v1/health/variable/{variable}", func(w http.ResponseWriter, r *http.Request) {
		GetHealth(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/health/v1/logged_in", func(w http.ResponseWriter, r *http.Request) {
		GetLoggedIn(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/health/v1/ref_variables", func(w http.ResponseWriter, r *http.Request) {
		GetRefVariables(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/health/v1/ref_variables/id/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetRefVariables(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/health/v1/ref_variables/variable/{variable}", func(w http.ResponseWriter, r *http.Request) {
		GetRefVariables(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/health/v1/health/keys/{keys}", func(w http.ResponseWriter, r *http.Request) {
		PostHealth(w, r, db)
	}).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:8080/health/v1/login", "https://example.com", "https://example.com/health/v1/login"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	})

	handler := c.Handler(router)
	listenPort := fmt.Sprintf(":%s", os.Getenv("LISTEN_PORT"))

	log.Printf("listening on %v", listenPort)

	log.Fatal(http.ListenAndServe(listenPort, handler))
}
