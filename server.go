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
	Artists  []Artists
	Relation []Relation
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
	ID             int                 `json:"id"`
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
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var artists []Artists
	json.Unmarshal([]byte(responseData), &artists)

	if err != nil {
		handleRequest400(err, w)
		return
	}

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

	if response.StatusCode != http.StatusOK {
		fmt.Print(err.Error())
		return
	}

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
	// x := map[string]interface{}{					Error 400 test
	// 	"foo": make(chan int),
	// }
	// i, err := json.Marshal(x)
	// fmt.Println(i)
	// fmt.Printf("Marshal error: %s\n", err)

	t.Execute(w, relation.Index)
	if err != nil {
		handleRequest400(err, w)
		return
	}
}

func main() {

	http.HandleFunc("/groupietracker", homePageHandler)
	http.HandleFunc("/groupietracker/", relationHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("Starting server on port", "8080")

}

func handleRequest400(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Error 400 Bad Request"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}
