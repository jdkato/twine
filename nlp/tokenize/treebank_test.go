package tokenize_test

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jdkato/twine/internal"
	"github.com/jdkato/twine/nlp/tokenize"
)

func getOldWordData(file string) ([]string, [][]string) {
	in := internal.ReadDataFile(filepath.Join(testdata, "treebank_sents_old.json"))
	out := internal.ReadDataFile(filepath.Join(testdata, file))

	input := []string{}
	output := [][]string{}

	err := json.Unmarshal(in, &input)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(out, &output)
	if err != nil {
		panic(err)
	}

	return input, output
}

func TestTreebankWordTokenizer(t *testing.T) {
	input, output := getOldWordData("treebank_words_old.json")
	word := tokenize.NewTreebankWordTokenizer()
	for i, s := range input {
		words := word.Tokenize(s)
		if !reflect.DeepEqual(output[i], words) {
			t.Errorf("%q, want %q", words, output[i])
		}
	}
}

func BenchmarkTreebankWordTokenizer(b *testing.B) {
	word := tokenize.NewTreebankWordTokenizer()
	for n := 0; n < b.N; n++ {
		for _, s := range getWordBenchData() {
			word.Tokenize(s)
		}
	}
}
