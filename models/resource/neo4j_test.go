package resource_test

// import (
// 	"errors"
// 	"testing"

// 	"github.com/neo4j/neo4j-go-driver/neo4j"
// 	"github.com/thatbeardo/go-sentinel/mocks"
// 	"github.com/thatbeardo/go-sentinel/models/resource"
// )

// func TestGetResourcesOk(t testing.T) {
// 	response := func() (neo4j.Result, error) {
// 		return nil, errors.New("Database error")
// 	}
// 	mockSession := mocks.MockSession{
// 		SetResponse: response,
// 	}
// 	resource.NewNeo4jRepository(mockSession)
// }
