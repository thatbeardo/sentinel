package neo4j

import (
	"fmt"
	"os"
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
}

// New is a factory method to isntantiate Runner objects
func NewRunner() Runner {
	return runner{}
}

func (n runner) Run(statement string, parameters map[string]interface{}) (data []map[string]interface{}, err error) {
	fmt.Println(statement)
	cleanup, session := Initialize()
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
func Initialize() (func(), neo4j.Session) {
	session, driver, err := ConnectToDB()
	fmt.Println(err)
	return func() {
		session.Close()
		driver.Close()
	}, session
}

// ConnectToDB establishes connection to the neo4j database
func ConnectToDB() (neo4j.Session, neo4j.Driver, error) {
	// define driver, session and result vars
	var (
		driver  neo4j.Driver
		session neo4j.Session
		err     error
	)
	// initialize driver to connect to DB with ID and password
	dbURI := os.Getenv("DB_URI")
	fmt.Println("Now connecting " + dbURI)
	if driver, err = neo4j.NewDriver(dbURI, neo4j.BasicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD"), ""), func(c *neo4j.Config) {
		c.Encrypted = true
	}); err != nil {
		fmt.Println("Error while establishing graph connection")
	}
	// Open a new session with write access
	if session, err = driver.Session(neo4j.AccessModeWrite); err != nil {
		return nil, nil, err
	}
	return session, driver, nil
}
