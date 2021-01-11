package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// User ...
type User struct {
	gorm.Model
	ID        int
	UserName  string
	FirstName string
	LastName  string
	Email     string
}

func home(w http.ResponseWriter, r *http.Request) {
	display(w, "home", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.Scopes(Paginate(r)).Find(&users)
	display(w, "index", users)
}

func show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user User
	db.First(&user, vars["id"])
	display(w, "show", user)
}

func create(w http.ResponseWriter, r *http.Request) {
	display(w, "create", nil)
}

func store(w http.ResponseWriter, r *http.Request) {
	logErr(r.ParseForm())
	user := User{
		UserName:  r.Form.Get("UserName"),
		FirstName: r.Form.Get("FirstName"),
		LastName:  r.Form.Get("LastName"),
		Email:     r.Form.Get("Email"),
	}
	db.Create(&user)
	index(w, r)
}

func main() {
	db, err = gorm.Open(mysql.Open("root:password@tcp(mysql)/local?parseTime=true"), &gorm.Config{})
	fatalErr(err)

	db.AutoMigrate(&User{})

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
	tpl := template.Must(template.ParseFiles(
		"templates/layout.html",
		fmt.Sprintf("templates/%s.html", view),
	))
	logErr(tpl.Execute(w, data))
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

// Paginate ...
func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	query := r.URL.Query()
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(query.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(query.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
