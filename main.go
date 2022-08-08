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
)

// returnJsonResponse: return the JSON-formatted response via HTTP
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

	// Shutdown := func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Server Shutdown."))
	// 	go func() {
	// 		if err := http.Shutdown(context.Background());
	// 		err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}()
	// }



	// http.HandleFunc("/", StartingPage)
	// http.HandleFunc("/movies", Moviepage)
	// http.HandleFunc("/persons", Personpage)
	// http.HandleFunc("/actors", Actorpage)
	// http.HandleFunc("/reviewers", Reviewerpage)
	// http.HandleFunc("/rankings", RankingPages)
	// http.HandleFunc("/rankings/reviewranking", ReviewRanking)
	// http.HandleFunc("/rankings/actorranking", ActorRanking)
	// http.HandleFunc("/shutdown", Shutdown)
	// log.Fatal(http.ListenAndServe(":1000", nil))

	m := http.NewServeMux()
	s := http.Server{Addr: ":1000", Handler: m}
	m.HandleFunc("/", StartingPage)
	m.HandleFunc("/movies", Moviepage)
	m.HandleFunc("/persons", Personpage)
	m.HandleFunc("/actors", Actorpage)
	m.HandleFunc("/reviewers", Reviewerpage)
	m.HandleFunc("/rankings", RankingPages)
	m.HandleFunc("/rankings/reviewranking", ReviewRanking)
	m.HandleFunc("/rankings/actorranking", ActorRanking)
	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Goodbye!"))
        go func() {
            if err := s.Shutdown(context.Background()); err != nil {
                log.Fatal(err)
            }
        }()
    })
    if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
    log.Printf("Finished")




	
}