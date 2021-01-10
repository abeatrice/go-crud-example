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

	tpl, err := template.New("index.html").ParseFiles("index.html")
	logErr(err)
	logErr(tpl.ExecuteTemplate(w, "index.html", users))
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

	tpl, err := template.New("show.html").ParseFiles("show.html")
	logErr(err)
	logErr(tpl.ExecuteTemplate(w, "show.html", u))
}

func main() {
	db, err = sql.Open("mysql", "root:password@tcp(mysql)/local")
	fatalErr(err)

	r := mux.NewRouter()
	r.HandleFunc("/users", index).Methods("GET").Name("users.index")
	r.HandleFunc("/users/{id:[0-9]+}", show).Methods("GET").Name("users.show")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
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
