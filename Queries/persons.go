package queries
import(
	"fmt"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"encoding/json"
)

type Person struct {
	ID int64 `json:"ID"`
	Name string `json:"Name"`
}

type Actor struct {
	ID 			int64 `json:"ID"`
	Name 		string `json:"Name"`
	Appearances int64 `json:"Appearances"`
}

type Reviewer struct {
	ID 			int64 `json:"ID"`
	Name 		string `json:"Name"`
	Reviews 	int64 `json:"Reviews"`
}

func PersonTransaction(session neo4j.Session) string {
	var persons []Person
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		getPersons := `
			MATCH (p:Person)
			RETURN id(p) AS id, p.name AS name`
		result, err := tx.Run(getPersons, map[string]interface{}{})
		if err != nil {return nil, err}
		
		for result.Next() {
			persons = append(persons, 
				Person{
					ID: result.Record().Values[0].(int64),
					Name: result.Record().Values[1].(string),
				})}
		return nil, result.Err()
	})
	
	if err != nil {panic(err)}
	record, error := json.Marshal(persons)
	stringrecord := string(record)
	if error != nil {fmt.Println(error)}
	return stringrecord
}

func ActorTransaction(session neo4j.Session) string {
	var actors []Actor 
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		getActors := `
			MATCH (p:Person)-[r:ACTED_IN]->(m:Movie)
			RETURN id(p) AS id, p.name AS name, COUNT(r) AS appearances`
		result, err := tx.Run(getActors, map[string]interface{}{})
		if err != nil {return nil, err }

		for result.Next() {
			actors = append(actors,
				Actor{
					ID: result.Record().Values[0].(int64),
					Name: result.Record().Values[1].(string),
					Appearances: result.Record().Values[2].(int64),
				})}
		return nil, result.Err()
	})

	if err != nil {panic(err)}
	record, error := json.Marshal(actors)
	stringrecord := string(record)
	if error != nil {fmt.Println(error)}
	return stringrecord
}

func ReviewerTransaction(session neo4j.Session) string {
	var reviewers []Reviewer
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		getReviewers := `
			MATCH (p:Person)-[r:REVIEWED]->(m:Movie)
			RETURN id(p) AS id, p.name AS name, COUNT(r) AS reviews`
		result, err := tx.Run(getReviewers, map[string]interface{}{})
		if err != nil {return nil, err}

		for result.Next() {
			reviewers = append(reviewers, 
				Reviewer{
					ID: result.Record().Values[0].(int64),
					Name: result.Record().Values[1].(string),
					Reviews: result.Record().Values[2].(int64),
				})}
		return nil, result.Err()
	})
	if err != nil {panic(err)}
	record, error := json.Marshal(reviewers)
	stringrecord := string(record)
	if error != nil {fmt.Println(error)}
	return stringrecord
}
