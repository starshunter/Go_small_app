package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Page struct {
	Title string
	Body  []byte
}

type List struct {
	Title string
	Link  Table
	P     string
}

type Row struct {
	RowItem []template.HTML
	title   string
}

type Table struct {
	Content []Row
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

var queryResult List

func userHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("user.html")
	name := r.URL.Query().Get("name")
	if err != nil {
		return
	}
	p := &Page{Title: name}
	t.Execute(w, p)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/search/", http.StatusFound)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("search.html")
	if err != nil {
		return
	}
	t.Execute(w, queryResult)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	db, err := sql.Open("mysql", "root:Arkreact3@/sakila?charset=utf8")
	checkErr(err)
	rows, err := db.Query(`SELECT actor_id, first_name, last_name FROM sakila.actor WHERE first_name="` + name + `" OR last_name="` + name + `"`)
	var results []Row
	for rows.Next() {
		var first_name string
		var last_name string
		var actor_id int
		var result []template.HTML

		err := rows.Scan(&actor_id, &first_name, &last_name)
		checkErr(err)

		var s string
		s = `<li><text>` + strconv.Itoa(actor_id) + `</text>`
		result = append(result, template.HTML(s))

		user_name := first_name + " " + last_name
		s = `<a href="/user.html?name=` + user_name + `">` + user_name + `</a></li>`
		result = append(result, template.HTML(s))
		results = append(results, Row{RowItem: result, title: "aaa"})
	}
	if len(results) == 0 {
		queryResult = List{}
		http.Redirect(w, r, "/search", http.StatusFound)
	} else {
		queryResult = List{Title: "hahaha", P: results[0].title, Link: Table{Content: results}}
		http.Redirect(w, r, "/search/", http.StatusFound)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/user.html", userHandler)
	http.HandleFunc("/result/", resultHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
