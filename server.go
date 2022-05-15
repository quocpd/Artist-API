package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	
)

type Response []struct {
	Artists   []Artists
	Relation  []Relation
}

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Next         int
	Previous     int
}
type Relation struct {
	Index []Relations `json:"index"`
}
type Relations struct {
	ID             int                `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if !(r.URL.Path == "/groupietracker" || r.URL.Path == "/groupietracker/") {
		http.Error(w, "404 Not Found.", http.StatusNotFound)
		return
	}
	t, _ := template.ParseFiles("homepage.html")

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var artists []Artists
	json.Unmarshal([]byte(responseData), &artists)

	for i := range artists {

		if i == 51 {
			artists[51].Next = 1
		} else {
			artists[i].Next = i + 2
		}
	}
	for i := range artists {

		if i == 0 {
			artists[0].Previous = 52
		} else {
			artists[i].Previous = i
		}
	}
	t.Execute(w, artists)
}

func relationHandler(w http.ResponseWriter, r *http.Request) {
	if !(r.URL.Path == "/groupietracker" || r.URL.Path == "/groupietracker/") {
		http.Error(w, "404 Not Found.", http.StatusNotFound)
		return
	}
	t, _ := template.ParseFiles("relation.html")

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	var relation Relation
	json.Unmarshal(responseData, &relation)
	t.Execute(w, relation.Index)

}

func main() {

	http.HandleFunc("/groupietracker", homePageHandler)
	http.HandleFunc("/groupietracker/", relationHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("Starting server on port", "8080")
}
