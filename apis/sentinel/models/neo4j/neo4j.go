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

type runner struct {
	driver neo4j.Driver
}

// NewRunner is a factory method to isntantiate Runner objects
func NewRunner(driver neo4j.Driver) Runner {
	return runner{
		driver: driver,
	}
}

func (n runner) Run(statement string, parameters map[string]interface{}) (data []map[string]interface{}, err error) {
	fmt.Println(statement)
	cleanup, session := Initialize(n.driver)
	defer cleanup()
	
	result, err := session.Run(statement, parameters)
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

// Initialize connects to the database and returns a shut down function
func Initialize(driver neo4j.Driver) (func(), neo4j.Session) {
	session, err := ConnectToDB(driver)
	fmt.Println(err)
	return func() {
		session.Close()
	}, session
}

// ConnectToDB establishes connection to the neo4j database
func ConnectToDB(driver neo4j.Driver) (neo4j.Session, error) {
	var (
		session neo4j.Session
		err     error
	)
	
	// Open a new session with write access
	if session, err = driver.Session(neo4j.AccessModeWrite); err != nil {
		return nil, err
	}
	return session, nil
}
