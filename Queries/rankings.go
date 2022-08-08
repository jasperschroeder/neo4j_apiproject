package queries 

import(
	"fmt"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"encoding/json"
)

type ReviewRank struct {
	AvgReview 	float64 `json:"AvgReview"`
	Rank		int64 `json:"Rank"`
	Title 		string `json:"Title"`
}

type ActorRank struct {
	Appearances 	int64 `json:"Appearances"`
	Rank 			int64 `json:"Rank"`
	Name 			string `json:"Name"`
}


func ReviewsRanking(session neo4j.Session, limit int) string {
	var reviews []ReviewRank
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		getReviewSummary := `
			MATCH (r:Person)-[s:REVIEWED]->(m:Movie)
			WITH ROUND(AVG(s.rating), 1) AS avgRating, m AS m 
			ORDER by avgRating DESC
			WITH collect(m) AS movies, collect(avgRating) AS ratings
			UNWIND movies as movie
			WITH movie, ratings, apoc.coll.indexOf(movies, movie) AS rank
			RETURN ratings[rank] AS rating, rank+1 AS rank, movie.title AS title
			LIMIT $limit`
		result, err := tx.Run(getReviewSummary, map[string]interface{}{"limit": limit})
		if err != nil {return nil, err}

		for result.Next() {
			reviews = append(reviews, 
				ReviewRank{
					AvgReview: result.Record().Values[0].(float64),
					Rank: result.Record().Values[1].(int64), 
					Title: result.Record().Values[2].(string),

				})}
		return nil, result.Err()
	})
	if err != nil {fmt.Println(err)}
	record, error := json.Marshal(reviews)
	stringrecord := string(record)
	if error != nil {panic(error)}
	return stringrecord
}


func ActorRanking(session neo4j.Session, limit int) string {
	var appearances []ActorRank 
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		getAppearancesSummary := `
			MATCH (p:Person)-[a:ACTED_IN]->(m:Movie)
			WITH COUNT(a) AS appearances, p AS p
			ORDER BY appearances DESC
			WITH collect(p) AS actors, collect(appearances) AS appearances
			UNWIND actors as actor
			WITH actor, appearances, apoc.coll.indexOf(actors, actor) AS rank
			RETURN appearances[rank] AS appearances, rank+1 AS rank, actor.name AS name
			LIMIT $limit`
		result, err := tx.Run(getAppearancesSummary, map[string]interface{}{"limit": limit})
		if err != nil {return nil, err}

		for result.Next() {
			appearances = append(appearances, 
				ActorRank{
					Appearances: result.Record().Values[0].(int64),
					Rank: result.Record().Values[1].(int64),
					Name: result.Record().Values[2].(string),
				})}
		return nil, result.Err()
		})
		if err != nil {fmt.Println(err)}
		record, error := json.Marshal(appearances)
		stringrecord := string(record)
		if error != nil {panic(error)}
		return stringrecord
}

