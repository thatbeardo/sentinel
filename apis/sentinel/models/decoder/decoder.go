package decoder

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/decoder/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/neo4j"
)

// DecodeNode marshalls neo4j nodes into the target interface
func DecodeNode(results map[string]interface{}, field string, target interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()
	if results[field] != nil {
		node := results[field].(neo4j.Node)
		err = injection.MapDecoder(node.Props(), &target)
	}
	return
}

// DecodeEdge unmarshalls a neo4j relationship into target
func DecodeEdge(results map[string]interface{}, field string, target interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()
	if results[field] != nil {
		node := results[field].(neo4j.Relationship)
		err = injection.MapDecoder(node.Props(), &target)
	}
	return
}
