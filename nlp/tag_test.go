package nlp

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

func makeTagger(text string) (*Document, error) {
	return NewDocument(
		text,
		WithSegmentation(false))
}

func ExampleReadTagged() {
	tagged := "Pierre|NNP Vinken|NNP ,|, 61|CD years|NNS"
	fmt.Println(ReadTagged(tagged, "|"))
	// Output: [[[Pierre Vinken , 61 years] [NNP NNP , CD NNS]]]
}

func TestTagSimple(t *testing.T) {
	doc, err := makeTagger("Pierre Vinken, 61 years old, will join the board as a nonexecutive director Nov. 29.")
	if err != nil {
		panic(err)
	}
	tags := []string{}
	for _, tok := range doc.Tokens() {
		tags = append(tags, tok.Tag)
	}
	if !reflect.DeepEqual([]string{
		"NNP", "NNP", ",", "CD", "NNS", "JJ", ",", "MD", "VB", "DT", "NN",
		"IN", "DT", "JJ", "NN", "NNP", "CD", "."}, tags) {
		t.Errorf("TagSimple() got = %v", tags)
	}
}

func TestTagTreebank(t *testing.T) {
	tokens, expected := []*Token{}, []string{}

	tagger, err := newPerceptronTagger()
	if err != nil {
		panic(err)
	}

	tags := readDataFile(filepath.Join(testdata, "treebank_tags.json"))
	err = json.Unmarshal(tags, &expected)
	if err != nil {
		panic(err)
	}

	treebank := readDataFile(filepath.Join(testdata, "treebank_tokens.json"))
	err = json.Unmarshal(treebank, &tokens)
	if err != nil {
		panic(err)
	}

	correct := 0.0
	for i, tok := range tagger.tag(tokens) {
		if expected[i] == tok.Tag {
			correct++
		}
	}

	v := correct / float64(len(expected))
	if v < 0.957477 {
		t.Errorf("TagTreebank() expected >= 0.957477, got = %v", v)
	}
}

func BenchmarkTag(b *testing.B) {
	tagger, err := newPerceptronTagger()
	if err != nil {
		panic(err)
	}
	tokens := []*Token{}

	treebank := readDataFile(filepath.Join(testdata, "treebank_tokens.json"))

	err = json.Unmarshal(treebank, &tokens)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		_ = tagger.tag(tokens)
	}
}
