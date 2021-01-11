package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type user struct {
	ID        int
	UserName  string
	FirstName string
	LastName  string
	Email     string
}

func home(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("./templates/*")
	logErr(err)
	logErr(tpl.ExecuteTemplate(w, "home.html", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`
		SELECT id, username, first_name, last_name, email
		FROM users
	`)
	logErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	logErr(err)

	var users []user
	for rows.Next() {
		var u user
		err = rows.Scan(&u.ID, &u.UserName, &u.FirstName, &u.LastName, &u.Email)
		logErr(err)
		users = append(users, u)
	}

	display(w, "index.html", users)
}

func show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := db.QueryRow(`
		SELECT id, username, first_name, last_name, email
		FROM users
		WHERE id = ?
	`, id)

	var u user
	err = row.Scan(&u.ID, &u.UserName, &u.FirstName, &u.LastName, &u.Email)
	logErr(err)

	display(w, "show.html", u)
}

func create(w http.ResponseWriter, r *http.Request) {
	display(w, "create.html", nil)
}

func store(w http.ResponseWriter, r *http.Request) {
	logErr(r.ParseForm())
	user := user{
		UserName:  r.Form.Get("UserName"),
		FirstName: r.Form.Get("FirstName"),
		LastName:  r.Form.Get("LastName"),
		Email:     r.Form.Get("Email"),
	}

	stmt, err := db.Prepare(`
		INSERT INTO users (username, first_name, last_name, email)
		VALUES (?, ?, ?, ?)
	`)
	logErr(err)
	defer stmt.Close()

	stmt.Exec(user.UserName, user.FirstName, user.LastName, user.Email)

	index(w, r)
}

func main() {
	db, err = sql.Open("mysql", "root:password@tcp(mysql)/local")
	fatalErr(err)

	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET").Name("home")
	r.HandleFunc("/users", index).Methods("GET").Name("users.index")
	r.HandleFunc("/users/{id:[0-9]+}", show).Methods("GET").Name("users.show")
	r.HandleFunc("/users/create", create).Methods("GET").Name("users.create")
	r.HandleFunc("/users", store).Methods("POST").Name("users.store")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func display(w http.ResponseWriter, view string, data interface{}) {
	tpl, err := template.ParseGlob("./templates/*")
	logErr(err)
	logErr(tpl.ExecuteTemplate(w, view, data))
}

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
