package tag

import (
	"bytes"
	_ "embed"
	"encoding/gob"
)

var wts map[string]map[string]float64
var tags map[string]string
var classes []string

//go:embed classes.gob
var encodedClasses []byte

//go:embed tags.gob
var encodedTags []byte

//go:embed weights.gob
var encodedWeights []byte

func init() {
	dec := gob.NewDecoder(bytes.NewReader(encodedClasses))
	err := dec.Decode(&classes)
	if err != nil {
		panic(err)
	}

	dec = gob.NewDecoder(bytes.NewReader(encodedTags))
	err = dec.Decode(&tags)
	if err != nil {
		panic(err)
	}

	dec = gob.NewDecoder(bytes.NewReader(encodedWeights))
	err = dec.Decode(&wts)
	if err != nil {
		panic(err)
	}
}
