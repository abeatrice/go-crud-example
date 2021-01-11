package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error
var validate *validator.Validate

// User ...
type User struct {
	gorm.Model
	UserName  string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `validate:"required,email"`
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
	logErr(r.ParseForm())
	data := struct {
		Old User
	}{
		User{
			UserName:  r.Form.Get("UserName"),
			FirstName: r.Form.Get("FirstName"),
			LastName:  r.Form.Get("LastName"),
			Email:     r.Form.Get("Email"),
		},
	}
	display(w, "create", data)
}

func store(w http.ResponseWriter, r *http.Request) {
	logErr(r.ParseForm())
	user := User{
		UserName:  r.Form.Get("UserName"),
		FirstName: r.Form.Get("FirstName"),
		LastName:  r.Form.Get("LastName"),
		Email:     r.Form.Get("Email"),
	}
	err = validate.Struct(user)
	if err != nil {
		create(w, r)
		return
	}

	db.Create(&user)
	index(w, r)
}

func main() {
	validate = validator.New()
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
