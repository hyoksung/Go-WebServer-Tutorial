package main

import (
	"GO-WebServer-Tutorial/02-REST-GET-WebServer/countryCapitals"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const port = ":9000"

func home(w http.ResponseWriter, r *http.Request) {
	urlPathElements := strings.Split(r.URL.Path, "/")

	if urlPathElements[1] == "capital" {
		capital := countryCapitals.Capitals[urlPathElements[2]]

		//If not match found, empty string is returned. Use it for validation!
		if capital != "" {
			fmt.Fprint(w, capital)
		} else {
			//Returns 404, Not found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Sorry, Resource Not Found!"))
		}
	} else {
		//Returns 400, Bad Request
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Sorry, Bad API Request!"))
	}
}

func getServerTime(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, time.Now().String())
}

func getRandomNumber(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, rand.Int())
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User belongs to category: %v\n", vars["category"])
	fmt.Fprintf(w, "User ID is: %v\n", vars["id"])
}

func main() {
	//Introducing NewServeMux for multi-routing
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", home)
	rtr.HandleFunc("/servertime", getServerTime)
	rtr.HandleFunc("/random", getRandomNumber)
	rtr.HandleFunc("/users/{category}/{id:[0-9]+}", handleUsers)
	rtr.HandleFunc("/capital/{country}", home)
	err := rtr.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server listening on port " + port)
	log.Fatal(http.ListenAndServe(port, rtr))
}
