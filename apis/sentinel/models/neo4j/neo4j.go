package neo4j

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Relationship represents a relationship in the neo4j graph database
type Relationship interface {
	Id() int64
	StartId() int64
	EndId() int64
	Type() string
	Props() map[string]interface{}
}

// Node represents a node in the database
type Node interface {
	Id() int64
	Labels() []string
	Props() map[string]interface{}
}

// Runner encapsulates operations that are needed to be carried out at the DB level
type Runner interface {
	Run(string, map[string]interface{}) ([]map[string]interface{}, error)
}

type neo4jSession struct {
	session neo4j.Session
}

// New is a factory method to isntantiate Runner objects
func New(session neo4j.Session) Runner {
	return neo4jSession{
		session: session,
	}
}

func (n neo4jSession) Run(statement string, parameters map[string]interface{}) (data []map[string]interface{}, err error) {
	result, err := n.session.Run(statement, parameters)
	fmt.Println(statement)
	data = []map[string]interface{}{}

	if err != nil {
		fmt.Println("Could not connect database ", err.Error())
		return
	}
	for result.Next() {
		data = append(data, result.Record().GetByIndex(0).(map[string]interface{}))
	}
	return
}