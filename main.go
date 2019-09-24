package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/zwergon/pwallet/utils"

	_ "github.com/go-sql-driver/mysql"
)

//RegistrationInfo to store data for web identification
type RegistrationInfo struct {
	Id      int
	Comment string
	Company string
	Detail  string
	Login   string
	Passwd  string
}

var tmpl = template.Must(template.ParseGlob("form/*"))

//IndexHandler request to show all elements
type IndexHandler struct {
	db *sql.DB
}

//Index read all data in registrationinfo table
func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	selDB, err := h.db.Query("SELECT id,company,comment FROM registrationinfo ORDER BY company ASC")
	if err != nil {
		panic(err.Error())
	}
	rInfo := RegistrationInfo{}
	res := []RegistrationInfo{}
	for selDB.Next() {
		var id int
		var company string
		var comment sql.NullString
		err = selDB.Scan(&id, &company, &comment)
		if err != nil {
			panic(err.Error())
		}
		rInfo.Id = id
		if comment.Valid {
			rInfo.Comment = comment.String
		} else {
			rInfo.Comment = ""
		}
		rInfo.Company = company
		res = append(res, rInfo)
	}

	json.NewEncoder(w).Encode(res)

}

type ShowHandler struct {
	db *sql.DB
}

//Show show one registrationinfy selected by id
func (s *ShowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	nID := r.Header.Get("id")
	selDB, err := s.db.Query("SELECT id,company,login,passwd FROM registrationinfo WHERE id=?", nID)
	if err != nil {
		panic(err.Error())
	}
	rInfo := RegistrationInfo{}
	for selDB.Next() {
		var id int
		var company, login, passwd string
		err = selDB.Scan(&id, &company, &login, &passwd)
		if err != nil {
			panic(err.Error())
		}
		rInfo.Id = id
		rInfo.Company = company
		rInfo.Login = login
		rInfo.Passwd = passwd
	}

	json.NewEncoder(w).Encode(rInfo)
}

/*
//New create New
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

//Edit edit a registrationinfo
func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	rInfo := RegistrationInfo{}
	for selDB.Next() {
		var id int
		var comment, company, detail, login, passwd string
		err = selDB.Scan(&id, &comment, &company, &detail, &login, &passwd)
		if err != nil {
			panic(err.Error())
		}
		rInfo.Id = id
		rInfo.Comment = comment
		rInfo.Company = company
		rInfo.Detail = detail
		rInfo.Login = login
		rInfo.Passwd = passwd
	}
	tmpl.ExecuteTemplate(w, "Edit", rInfo)
	defer db.Close()
}

//Insert insert a new element in database
func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city)
		log.Println("INSERT: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

//Update update an element in database
func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city, id)
		log.Println("UPDATE: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

//Delete remove an element in database
func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	rInfo := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(rInfo)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
*/

func main() {

	database := "local"
	if len(os.Args) > 1 {
		database = os.Args[1]
	}

	host := os.Getenv("IP")
	port := os.Getenv("PORT")

	adress := fmt.Sprintf("%s:%s", host, port)
	dbInfo := utils.NewDB(database)

	db := dbInfo.DbConn()
	log.Println("Server started on: http://" + adress)

	idxHandler := IndexHandler{db: db}
	http.Handle("/idx", &idxHandler)

	ShowHandler := ShowHandler{db: db}
	http.Handle("/show", &ShowHandler)
	//http.HandleFunc("/new", New)
	//http.HandleFunc("/edit", Edit)
	//http.HandleFunc("/insert", Insert)
	//http.HandleFunc("/update", Update)
	//http.HandleFunc("/delete", Delete)
	http.ListenAndServe(adress, nil)
	log.Println("Server stopped!")

	defer db.Close()
}
