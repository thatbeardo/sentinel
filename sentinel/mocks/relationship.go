package mocks

// Relationship is a mock edge present in the graphs
type Relationship struct {
	id      int64
	startId int64
	endId   int64
	relType string
	props   map[string]interface{}
}

// Id represents the internally generated unique ID
func (m Relationship) Id() int64 {
	return m.id
}

// StartId is the internally generated id by the database
func (m Relationship) StartId() int64 {
	return m.startId
}

// EndId is the internally generated id by the database
func (m Relationship) EndId() int64 {
	return m.endId
}

// Type is the type of the relationship
func (m Relationship) Type() string {
	return m.relType
}

// Props represents data that defines this edge
func (m Relationship) Props() map[string]interface{} {
	return m.props
}

// NewRelationship is a factory fucntion to generate mock relationships
func NewRelationship(id, startID, endID int64, relType string, props map[string]interface{}) Relationship {
	return Relationship{
		id:      id,
		startId: startID,
		endId:   endID,
		relType: relType,
		props:   props,
	}
}
