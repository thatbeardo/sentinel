package injection

import (
	"github.com/mitchellh/mapstructure"
)

// Reset takes all the injected vaiables back to their origianl values
func Reset() {
	MapDecoder = mapDecoderDefault
}

// MapDecoder converts a map into its corresponding struct
var MapDecoder = mapDecoderDefault
var mapDecoderDefault = mapstructure.Decode
