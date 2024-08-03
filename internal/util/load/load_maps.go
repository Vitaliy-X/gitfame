package load

import (
	_ "embed"
	"encoding/json"

	"gitlab.com/slon/shad-go/gitfame/internal"
	"gitlab.com/slon/shad-go/gitfame/internal/exceptions"
)

//go:embed language_extensions.json
var file []byte

func LoadMaps() []internal.FileMapping {
	var mappings []internal.FileMapping
	err := json.Unmarshal(file, &mappings)
	exceptions.Exception(err, "LoadMaps: could not unmarshal JSON data")

	return mappings
}
