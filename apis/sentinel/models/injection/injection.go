package injection

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/models/decoder"
	"github.com/mitchellh/mapstructure"
)

// Reset takes all the injected vaiables back to their origianl values
func Reset() {
	MapDecoder = mapDecoderDefault
	NodeDecoder = nodeDecoderDefault
	EdgeDecoder = edgeDecoderDefault
}

// MapDecoder converts a map into its corresponding struct
var MapDecoder = mapDecoderDefault
var mapDecoderDefault = mapstructure.Decode

// NodeDecoder marshalls neo4j nodes into the target interface
var NodeDecoder = nodeDecoderDefault
var nodeDecoderDefault = decoder.DecodeNode

// EdgeDecoder marshalls neo4j nodes into the target interface
var EdgeDecoder = edgeDecoderDefault
var edgeDecoderDefault = decoder.DecodeEdge
