package mocks

// Node represents a single entity in the neo4j graph
type Node struct {
	id     int64
	labels []string
	props  map[string]interface{}
}

// Id returns the internally generated neo4j ID for this node
func (m Node) Id() int64 {
	return m.id
}

// Labels denote the label associated with a ndoe
func (m Node) Labels() []string {
	return m.labels
}

// Props methods details the data returned by the query
func (m Node) Props() map[string]interface{} {
	return m.props
}

// NewNode is a factory method to generate mock nodes
func NewNode(id int64, labels []string, props map[string]interface{}) Node {
	return Node{
		id:     id,
		labels: labels,
		props:  props,
	}
}
