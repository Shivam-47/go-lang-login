package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	//"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Person struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "<USER>"
	dbPass := "<PASS>"
	dbHost := "<DB-HOST>"
	dbName := "<DB-NAME>"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("SQL connection success")
	}

	return db

}

func home(w http.ResponseWriter, req *http.Request) {

	fmt.Println("Welcome to Home page")

	http.ServeFile(w, req, "static/index.html")
}

func register(w http.ResponseWriter, req *http.Request) {

	fmt.Println("Welcome to Registration page")

	http.ServeFile(w, req, "static/register.html")
}

func addRecord(w http.ResponseWriter, req *http.Request) {
	db := dbConn()

	if req.Method == "POST" {
		name := req.FormValue("name")
		email := req.FormValue("email")
		password := req.FormValue("password")

		sql, err := db.Prepare("INSERT INTO login(name,email,password) VALUES(?,?,?)")

		if err != nil {
			panic(err.Error())
		}

		sql.Exec(name, email, password)
		log.Println("Name: " + name + " added successfully")
	}

	defer db.Close()
	http.Redirect(w, req, "/", 301)
}

func dashboard(w http.ResponseWriter, req *http.Request) {

	db := dbConn()
	var Name string
	var Email string
	var person Person
	var count int

	if req.Method == "POST" {
		email := req.FormValue("email")
		password := req.FormValue("password")

		sql := db.QueryRow("SELECT COUNT(*) FROM login WHERE email=? AND password=?", email, password).Scan(&count)

		if sql != nil {
			panic(sql.Error())
		}

		if count == 0 {
			http.ServeFile(w, req, "static/error.html")
		}

		err := db.QueryRow("SELECT id,name,email,password FROM login WHERE email=? AND password =?", email, password).Scan(&person.Id, &person.Name, &person.Email, &person.Password)

		if err != nil {
			fmt.Fprintf(w, "OOPS! Invalid login credentials")
			panic(err.Error())

		}

		Name = person.Name
		Email = person.Email

		fmt.Println("Name: " + Name + " Email: " + Email)
	}

	//fmt.Fprintf(w, "Welcome "+Name+"\n")
	//fmt.Fprintf(w, "Email ID: "+Email+"\n")

	template, err1 := template.ParseFiles("static/dashboard.html")

	if err1 != nil {
		log.Println("ERROR")
	}
	template.Execute(w, person)

	defer db.Close()

}

func main() {

	router := mux.NewRouter()

	//home
	router.HandleFunc("/", home).Methods("GET")

	//home - POST
	router.HandleFunc("/", dashboard).Methods("POST")

	//register -GET
	router.HandleFunc("/register", register).Methods("GET")

	//register -POST
	router.HandleFunc("/register", addRecord).Methods("POST")

	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(":"+string(port), router))

}
