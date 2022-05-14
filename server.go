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
	Locations []Locations
	Dates     []Dates
	Relation  []Relation
}

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
	Next         int
	Previous     int
	// test         []interface{}
}

type Locations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

type Dates struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type Relation struct {
	Index []struct {
		ID             int `json:"id"`
		DatesLocations string
	}
}

type groupieTrackerPage struct {
	Title string
	Test  string
}

func homePageHandler(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("homepage.html")

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Print(err.Error())
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

func groupieTrackerHandler(w http.ResponseWriter, r *http.Request) {
	// p := groupieTrackerPage{Title: "Groupie Tracker"}
	t, _ := template.ParseFiles("artists.html")

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// var response Response
	// json.Unmarshal(body, &Response)
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

func locationsHandler(w http.ResponseWriter, r *http.Request) {
	// p := groupieTrackerPage{Title: "Groupie Tracker"}
	t, _ := template.ParseFiles("locations.html")

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		fmt.Print(err.Error())
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
	var locations Locations
	json.Unmarshal([]byte(responseData), &locations)
	// fmt.Println(locations)
	fmt.Print(locations.Index)

	t.Execute(w, locations)
	// t.Execute(w, locations.Index[1].ID)
}

func datesHandler(w http.ResponseWriter, r *http.Request) {
	// p := groupieTrackerPage{Title: "Groupie Tracker"}
	t, _ := template.ParseFiles("dates.html")

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		fmt.Print(err.Error())
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

	// var response Response
	// json.Unmarshal(body, &Response)
	var dates Dates
	json.Unmarshal([]byte(responseData), &dates)
	// fmt.Println(locations)

	fmt.Print(dates.Index)

	t.Execute(w, dates)
	// t.Execute(w, locations.Index[1].ID)

}

// 	if err := r.ParseForm(); err != nil {
// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		return
// 	}
// 	input := r.FormValue("input")
// 	banner := r.FormValue("banner")
// 	path := "banners/" + banner + ".txt"

// 	// if input != "" && banner == "standard" || banner == "shadow" || banner == "thinkertoy" {
// 	// 	fmt.Fprintf(w, "%s\n", newline(input, path))
// 	// }
// 	if input == "" {
// 		http.Error(w, "400 Bad Request Error.", http.StatusBadRequest)
// 	} else if !(banner == "standard" || banner == "shadow" || banner == "thinkertoy") {
// 		http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
// 	} else {
// 		fmt.Fprintf(w, "%s\n", newline(input, path))

// 	}

// }

func main() {

	// fmt.Println(len(responseObject.Artists))

	// for _, p := range artists {
	// 	fmt.Println("Members:", p.Members[0])
	// }

	// for i := 0; i < len(responseObject.Artists); i++ {
	//     fmt.Println(responseObject.Artists[i].ID.Name)
	// }
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/artists", groupieTrackerHandler)
	http.HandleFunc("/locations", locationsHandler)
	http.HandleFunc("/dates", datesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("Starting server on port", "8080")
}

// func artistAPI() []Artists {

// 	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}
// 	defer response.Body.Close()

// 	responseData, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// var response Response
// 	// json.Unmarshal(body, &Response)
// 	var artists []Artists
// 	json.Unmarshal([]byte(responseData), &artists)

// 	fmt.Println(artists)
// 	return artists

// }

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// )

// // A Response struct to map the Entire Response
// type Response struct {
// 	Artists   string `json:"artists"`
// 	Locations string `json:"locations"`
// 	Dates     string `json:"dates"`
// 	Relation  string `json:"relation"`
// }

// func groupieTracker(w http.ResponseWriter, r *http.Request) {

// 	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}

// 	responseData, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var responseObject Response
// 	json.Unmarshal(responseData, &responseObject)

// 	fmt.Println(responseObject.Artists)
// 	fmt.Println(len(responseObject.Artists))

// 	// for i := 0; i < len(responseObject.Artists); i++ {
// 	// 	fmt.Println(responseObject.Artists[i].Species.Name)
// 	// }
// 	// fmt.Println(responseObject.Artists[52].Species.Name)

// 	if !(r.URL.Path == "/") {
// 		http.Error(w, "404 Not Found.", http.StatusNotFound)
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		return
// 	}
// 	// input := r.FormValue("input")
// 	// banner := r.FormValue("banner")

// 	// if input == "" {
// 	// 	http.Error(w, "400 Bad Request Error.", http.StatusBadRequest)
// 	// } else if !(banner == "standard" || banner == "shadow" || banner == "thinkertoy") {
// 	// 	http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
// 	// } else {
// 	// 	fmt.Println("do nothing")
// 	// }
// }

// func main() {
// 	// fileServer := http.FileServer(http.Dir("./static"))
// 	// http.Handle("/", fileServer)
// 	http.HandleFunc("/", groupieTracker)
// 	fmt.Printf("Starting server at port 8080\n")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatal(err)
// 	}
// }

// // implementing our newline function on the argument which is the string after main.go

// // Also logging errors if too many or too little arguments.

// package main

// import (
// 	"net/http"
// 	"text/template"
// )

// type movie struct {
// 	Title string
// }

// type page struct {
// 	Title     string
// 	TopMovies []movie
// }

// func main() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
// 		w.Header().Add("Content Type", "text/html")

// 		templates := template.New("template")
// 		templates.New("Body").Parse(doc)
// 		templates.New("List").Parse(docList)

// 		topMovies := []movie{{Title: "Movie 1"}, {Title: "Movie 2"}, {Title: "Movie 3"}}

// 		page := page{Title: "My Title", TopMovies: topMovies}
// 		templates.Lookup("Body").Execute(w, page)

// 	})

// 	http.ListenAndServe(":8000", nil)
// }

// const docList = `
// <ul >
//     {{range .}}
//     <li>{{.Title}}</li>
//     {{end}}
// </ul>
// `

// const doc = `
// <!DOCTYPE html>
// <html>
//     <head><title>{{.Title}}</title></head>
//     <body>
//         <h1>Hello Templates</h1>
//         {{template "List" .TopMovies}}
//     </body>
// </html>
// `
