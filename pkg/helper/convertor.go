package helper

import "github.com/mitchellh/mapstructure"

// Map2Struct map转struct
func Map2Struct(mapData interface{}, structPoint interface{}) error {
	return mapstructure.Decode(mapData, structPoint)
}
