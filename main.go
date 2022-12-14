package main 

import (
	"fmt"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"encoding/json"
	queries "github.com/jasperschroeder/apiproject/queries"
	"net/http"
	"log"
	"strconv"
	"context"
	"time"
)

func returnJsonResponse(res http.ResponseWriter, httpCode int, resMessage []byte) {
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(resMessage)
}

func StartingPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! If you can read this, things worked out. \nWelcome to the starting page!")
}

func RankingPages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, 
`Hello, on the following pages you can see different rankings. 
/rankings/reviewranking shows movies with highest reviews.
/rankings/actorranking shows actors with most appearances.
By passing the limit parameter ?limit= you can specify how many results should be returned (max).`)
}


func main() {
	configData, err := GetConfigData("config.json")
	if err != nil {
		fmt.Println(err)
	}

	auth := neo4j.BasicAuth(configData.Username, configData.Password, "")
	driver, err := neo4j.NewDriver(configData.Uri, auth)
	if err != nil {
		panic(err)
	}
	defer driver.Close()
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close()

	Moviepage := func(w http.ResponseWriter, r *http.Request) {
		moviestring := queries.MovieTransaction(session)	
		json.NewEncoder(w).Encode(moviestring)
	}
	
	Personpage := func(w http.ResponseWriter, r *http.Request) {
		personstring := queries.PersonTransaction(session)
		json.NewEncoder(w).Encode(personstring)
	}

	Actorpage := func(w http.ResponseWriter, r *http.Request) {
		actorstring := queries.ActorTransaction(session)
		json.NewEncoder(w).Encode(actorstring)
	}

	Reviewerpage := func(w http.ResponseWriter, r *http.Request) {
		reviewerstring := queries.ReviewerTransaction(session)
		json.NewEncoder(w).Encode(reviewerstring)
	}

	ReviewRanking := func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query()["limit"]
		var limitParam int
		if limit == nil {
			limitParam = 5 
		} else {
			limitParam, err = strconv.Atoi(limit[0])
			if err != nil {limitParam = 5}
		}
		reviewrankingstring := queries.ReviewsRanking(session, limitParam)
		json.NewEncoder(w).Encode(reviewrankingstring)
	}

	ActorRanking := func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query()["limit"]
		var limitParam int 
		if limit == nil {
			limitParam = 5 
		} else {
			limitParam, err = strconv.Atoi(limit[0])
			if err != nil {limitParam = 5}
		}
		actorrankingstring := queries.ActorRanking(session, limitParam)
		json.NewEncoder(w).Encode(actorrankingstring)
	}

	m := http.NewServeMux()
	s := http.Server{Addr: ":1000", Handler: m}
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "" {
			http.NotFound(w, r)
			return 
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	m.HandleFunc("/welcome", StartingPage)
	m.HandleFunc("/redirect/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/welcome", 303)
	})

	expires := time.Now().AddDate(1, 0, 0)
	ck := http.Cookie{
		Name: "COOKIE",
		Path: "/cookie",
		Expires: expires,
	}
	ck.Value = "Test value."
		m.HandleFunc("/cookie", func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &ck)
			fmt.Fprintf(w, "Cookie is being set.")
	})

	m.HandleFunc("/movies", Moviepage)
	m.HandleFunc("/persons", Personpage)
	m.HandleFunc("/actors", Actorpage)
	m.HandleFunc("/reviewers", Reviewerpage)
	m.HandleFunc("/rankings", RankingPages)
	m.HandleFunc("/rankings/reviewranking", ReviewRanking)
	m.HandleFunc("/rankings/actorranking", ActorRanking)
	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		password := r.URL.Query()["password"]
		if password == nil {
			fmt.Fprintf(w, "Please provide a password as a parameter via the query.")
		} else if  password[0] !=  "ShutdownPassword" {
			fmt.Fprintf(w, "That was not the correct password.")
		} else {
			w.Write([]byte("Goodbye!"))
			go func() {if err := s.Shutdown(context.Background()); err != nil {log.Fatal(err)}}()
		}
    })
    if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {log.Fatal(err)}
    log.Printf("Finished")
}