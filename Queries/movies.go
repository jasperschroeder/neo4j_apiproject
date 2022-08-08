package queries 

import(
	"fmt"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"encoding/json"
)

type Movie struct {
	ID int64 `json:"ID"`
	Title string `json:"Title"`
}

func MovieTransaction(session neo4j.Session) string {
	var movies []Movie 
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		getMovies := `
			MATCH (m:Movie)
			RETURN id(m) AS id, m.title AS title`
		result, err := tx.Run(getMovies, map[string]interface{}{})
		if err != nil {return nil, err}

		for result.Next() {
			movies = append(movies, 
				Movie{
					ID: result.Record().Values[0].(int64),
					Title: result.Record().Values[1].(string),
				})}
		return nil, result.Err()
	})

	if err != nil {fmt.Println(err)}
	record, error := json.Marshal(movies)
	stringrecord := string(record)
	if error != nil {panic(error)}	
	return stringrecord
}