package summarize

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/jdkato/twine/internal"
)

func BenchmarkEasyWordsLookupMap(b *testing.B) {
	cases := internal.ReadDataFile(filepath.Join(testdata, "syllables.json"))
	tests := make(map[string]int)

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		for word := range tests {
			if _, ok := easyWords[word]; ok {

			}
		}
	}
}
