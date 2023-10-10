package tokenize_test

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jdkato/twine/internal"
	"github.com/jdkato/twine/nlp/tokenize"
)

var testdata = "../../testdata"
var tokenizer = tokenize.NewIterTokenizer()

func checkTokens(t *testing.T, tokens []string, expected []string, name string) {
	observed := []string{}
	for i := range tokens {
		observed = append(observed, tokens[i])
	}
	if !reflect.DeepEqual(observed, expected) {
		t.Errorf("%v: unexpected tokens", name)
	}
}

func checkCase(t *testing.T, tokens []string, expected []string, name string) {
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("%v: unexpected tokens", name)
	}
}

func getWordData(file string) ([]string, [][]string) {
	in := internal.ReadDataFile(filepath.Join(testdata, "treebank_sents.json"))
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

func getWordBenchData() []string {
	in := internal.ReadDataFile(filepath.Join(testdata, "treebank_sents.json"))
	input := []string{}

	err := json.Unmarshal(in, &input)
	if err != nil {
		panic(err)
	}

	return input
}

func TestTokenizationEmpty(t *testing.T) {
	doc := tokenizer.Tokenize("")

	l := len(doc)
	if l != 0 {
		t.Errorf("TokenizationEmpty() expected = 0, got = %v", l)
	}
}

func TestTokenizationSimple(t *testing.T) {
	doc := tokenizer.Tokenize("Vale is a natural language linter that supports plain text, markup (Markdown, reStructuredText, AsciiDoc, and HTML), and source code comments. Vale doesn't attempt to offer a one-size-fits-all collection of rules—instead, it strives to make customization as easy as possible.")
	expected := []string{
		"Vale", "is", "a", "natural", "language", "linter", "that", "supports",
		"plain", "text", ",", "markup", "(", "Markdown", ",", "reStructuredText",
		",", "AsciiDoc", ",", "and", "HTML", ")", ",", "and", "source",
		"code", "comments", ".", "Vale", "does", "n't", "attempt", "to",
		"offer", "a", "one-size-fits-all", "collection", "of", "rules—instead",
		",", "it", "strives", "to", "make", "customization", "as", "easy", "as",
		"possible", "."}
	checkCase(t, doc, expected, "TokenizationSimple()")
}

func TestTokenizationTreebank(t *testing.T) {
	input, output := getWordData("treebank_words.json")
	for i, s := range input {
		tokens := tokenizer.Tokenize(s)
		if !reflect.DeepEqual(tokens, output[i]) {
			t.Errorf("TokenizationTreebank(): unexpected tokens")
		}
	}
}

func TestTokenizationWeb(t *testing.T) {
	web := `Independent of current body composition, IGF-I levels at 5 yr were significantly
            associated with rate of weight gain between 0-2 yr (beta=0.19; P&lt;0.0005);
            and children who showed postnatal catch-up growth (i.e. those who showed gains in
            weight or length between 0-2 yr by >0.67 SD score) had higher IGF-I levels than other
			children (P=0.02; http://univ.edu.es/study.html) [20-22].`
	expected := []string{"Independent", "of", "current", "body", "composition", ",", "IGF-I",
		"levels", "at", "5", "yr", "were", "significantly", "associated", "with", "rate", "of",
		"weight", "gain", "between", "0-2", "yr", "(", "beta=0.19", ";", "P&lt;0.0005", ")", ";",
		"and", "children", "who", "showed", "postnatal", "catch-up", "growth", "(", "i.e.", "those",
		"who", "showed", "gains", "in", "weight", "or", "length", "between", "0-2", "yr", "by",
		">0.67", "SD", "score", ")", "had", "higher", "IGF-I", "levels", "than", "other", "children",
		"(", "P=0.02", ";", "http://univ.edu.es/study.html", ")", "[", "20-22", "]", "."}
	doc := tokenizer.Tokenize(web)
	checkCase(t, doc, expected, "TokenizationWeb()")
}

func TestTokenizationWebParagraph(t *testing.T) {
	web := `Independent of current body composition, IGF-I levels at 5 yr were significantly
            associated with rate of weight gain between 0-2 yr (beta=0.19; P&lt;0.0005);
            and children who showed postnatal catch-up growth (i.e. those who showed gains in
            weight or length between 0-2 yr by >0.67 SD score) had higher IGF-I levels than other
			children (P=0.02; http://univ.edu.es/study.html) [20-22].

			Independent of current body composition, IGF-I levels at 5 yr were significantly
            associated with rate of weight gain between 0-2 yr (beta=0.19; P&lt;0.0005);
            and children who showed postnatal catch-up growth (i.e. those who showed gains in
            weight or length between 0-2 yr by >0.67 SD score) had higher IGF-I levels than other
			children (P=0.02; http://univ.edu.es/study.html) [20-22].

			Independent of current body composition, IGF-I levels at 5 yr were significantly
            associated with rate of weight gain between 0-2 yr (beta=0.19; P&lt;0.0005);
            and children who showed postnatal catch-up growth (i.e. those who showed gains in
            weight or length between 0-2 yr by >0.67 SD score) had higher IGF-I levels than other
			children (P=0.02; http://univ.edu.es/study.html) [20-22].`

	expected := []string{"Independent", "of", "current", "body", "composition", ",", "IGF-I",
		"levels", "at", "5", "yr", "were", "significantly", "associated", "with", "rate", "of",
		"weight", "gain", "between", "0-2", "yr", "(", "beta=0.19", ";", "P&lt;0.0005", ")", ";",
		"and", "children", "who", "showed", "postnatal", "catch-up", "growth", "(", "i.e.", "those",
		"who", "showed", "gains", "in", "weight", "or", "length", "between", "0-2", "yr", "by",
		">0.67", "SD", "score", ")", "had", "higher", "IGF-I", "levels", "than", "other", "children",
		"(", "P=0.02", ";", "http://univ.edu.es/study.html", ")", "[", "20-22", "]", ".", "Independent", "of", "current", "body", "composition", ",", "IGF-I",
		"levels", "at", "5", "yr", "were", "significantly", "associated", "with", "rate", "of",
		"weight", "gain", "between", "0-2", "yr", "(", "beta=0.19", ";", "P&lt;0.0005", ")", ";",
		"and", "children", "who", "showed", "postnatal", "catch-up", "growth", "(", "i.e.", "those",
		"who", "showed", "gains", "in", "weight", "or", "length", "between", "0-2", "yr", "by",
		">0.67", "SD", "score", ")", "had", "higher", "IGF-I", "levels", "than", "other", "children",
		"(", "P=0.02", ";", "http://univ.edu.es/study.html", ")", "[", "20-22", "]", ".", "Independent", "of", "current", "body", "composition", ",", "IGF-I",
		"levels", "at", "5", "yr", "were", "significantly", "associated", "with", "rate", "of",
		"weight", "gain", "between", "0-2", "yr", "(", "beta=0.19", ";", "P&lt;0.0005", ")", ";",
		"and", "children", "who", "showed", "postnatal", "catch-up", "growth", "(", "i.e.", "those",
		"who", "showed", "gains", "in", "weight", "or", "length", "between", "0-2", "yr", "by",
		">0.67", "SD", "score", ")", "had", "higher", "IGF-I", "levels", "than", "other", "children",
		"(", "P=0.02", ";", "http://univ.edu.es/study.html", ")", "[", "20-22", "]", "."}

	doc := tokenizer.Tokenize(web)
	checkCase(t, doc, expected, "TokenizationWebParagraph()")
}

func TestTokenizationTwitter(t *testing.T) {
	doc := tokenizer.Tokenize("@twitter, what time does it start :-)")
	expected := []string{"@twitter", ",", "what", "time", "does", "it", "start", ":-)"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(1)")

	doc = tokenizer.Tokenize("Mr. James plays basketball in the N.B.A., do you?")
	expected = []string{
		"Mr.", "James", "plays", "basketball", "in", "the", "N.B.A.", ",",
		"do", "you", "?"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(2)")

	doc = tokenizer.Tokenize("ˌˌ kill the last letter")
	expected = []string{"ˌˌ", "kill", "the", "last", "letter"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(3)")

	doc = tokenizer.Tokenize("ˌˌˌ kill the last letter")
	expected = []string{"ˌˌˌ", "kill", "the", "last", "letter"}
	checkCase(t, doc, expected, "TokenizationWebParagraph(4)")

	doc = tokenizer.Tokenize("March. July. March. June. January.")
	expected = []string{
		"March", ".", "July", ".", "March", ".", "June", ".", "January", "."}
	checkCase(t, doc, expected, "TokenizationWebParagraph(5)")
}

func TestTokenizationSplitCases(t *testing.T) {
	tokenizer := tokenize.NewIterTokenizer(tokenize.UsingSplitCases([]string{"("}))
	tokens := tokenizer.Tokenize("amount($)")
	expected := []string{"amount", "(", "$", ")"}
	checkTokens(t, tokens, expected, "TokenizationSplitCases(custom-found)")
}

func TestTokenizationContractions(t *testing.T) {
	tokens := tokenizer.Tokenize("He's happy")
	expected := []string{"He", "'s", "happy"}
	checkTokens(t, tokens, expected, "TokenizationContraction(default-found)")

	tokens = tokenizer.Tokenize("I've been better")
	expected = []string{"I've", "been", "better"}
	checkTokens(t, tokens, expected, "TokenizationContraction(default-missing)")

	tokenizer = tokenize.NewIterTokenizer(tokenize.UsingContractions([]string{"'ve"}))
	tokens = tokenizer.Tokenize("I've been better")
	expected = []string{"I", "'ve", "been", "better"}
	checkTokens(t, tokens, expected, "TokenizationContraction(custom-found)")

	tokens = tokenizer.Tokenize("He's happy")
	expected = []string{"He's", "happy"}
	checkTokens(t, tokens, expected, "TokenizationContraction(custom-missing)")
}

func BenchmarkTokenization(b *testing.B) {
	in := internal.ReadDataFile(filepath.Join(testdata, "sherlock.txt"))
	text := string(in)
	for n := 0; n < b.N; n++ {
		_ = tokenizer.Tokenize(text)
	}
}

func BenchmarkTokenizationSimple(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, s := range getWordBenchData() {
			_ = tokenizer.Tokenize(s)
		}
	}
}
