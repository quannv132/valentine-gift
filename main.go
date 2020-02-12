package main

import (
	"app/src"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

//Confession storge data
type Confession struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

//DataSending send data to html template
type DataSending struct {
	Data []string
	Date string
}

var client *db.Client

var ctx = context.Background()

func init() {
	opt := option.WithCredentialsFile("valentine-4a342-firebase-adminsdk-iexzx-1e5a658a38.json")
	// ctx := context.Background()
	config := &firebase.Config{
		DatabaseURL: "https://valentine-4a342.firebaseio.com",
	}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatal(err)
	}

	client, err = app.Database(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	templates := template.Must(template.ParseFiles("templates/index.html"))

	param := mux.Vars(r)

	var confes Confession
	if err := client.NewRef("confession/"+param["id"]+"").Get(ctx, &confes); err != nil {
		log.Fatal(err)
	}

	strParts := strings.Split(confes.Content, "\n")

	data := DataSending{
		Data: strParts,
		Date: confes.Date,
	}

	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func form(w http.ResponseWriter, r *http.Request) {

	templates := template.Must(template.ParseFiles("templates/form.html"))

	if r.Method == "GET" {

		if err := templates.Execute(w, "form.html"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	if r.Method == "POST" {
		r.ParseForm()

		var hosting = "http://localhost:8080/"

		conf := Confession{
			Name:    r.FormValue("crushName"),
			Content: r.FormValue("content"),
			Date:    r.FormValue("crushDate"),
		}
		key := src.RandStringRunes(20)
		hosting += key

		fmt.Println(hosting)

		if err := client.NewRef("confession/"+key+"").Set(ctx, conf); err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(hosting))
	}
}

//Go application entrypoint
func main() {

	router := mux.NewRouter()

	// http.Handle("/static/",
	// 	http.StripPrefix("/static/",
	// 		http.FileServer(http.Dir("static"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/{id}", index)
	router.HandleFunc("/", form)

	fmt.Println("Listening")
	http.ListenAndServe(":8080", router)
}
