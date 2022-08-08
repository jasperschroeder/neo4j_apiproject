package main

// tag::import[]
import (
	"encoding/json"
	"io/ioutil"
)

func GetConfigData(configpath string) (*Neo4jData, error) {
	file, err := ioutil.ReadFile(configpath)
	if err != nil {
		return nil, err
	}
	config := Neo4jData	{}
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err 
	}
	return &config, nil 
}


type Neo4jData struct {
	Uri      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}