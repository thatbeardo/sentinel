package mocks

import "github.com/neo4j/neo4j-go-driver/neo4j"

type MockSession struct {
	SetResponse func() (neo4j.Result, error)
}

func (m MockSession) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	return m.SetResponse()
}

func (m MockSession) Close() error {
	return nil
}
