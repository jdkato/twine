package strcase

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var reNumberList = regexp.MustCompile(`\d+\.`)
var titleCased = cases.Title(language.English)

// A SentenceConverter converts a string to sentence case.
type SentenceConverter struct {
	vocab     []string
	indicator IndicatorFunc
}

// SentenceOptFunc is a function that modifies a SentenceConverter.
type SentenceOptFunc func(*SentenceConverter)

// An IndicatorFunc is a SentenceConverter callback that decides whether or not
// the string word should be capitalized.
type IndicatorFunc func(word string, idx int) bool

// UsingVocab sets the vocab for the SentenceConverter.
func UsingVocab(vocab []string) SentenceOptFunc {
	return func(converter *SentenceConverter) {
		converter.vocab = vocab
	}
}

// UsingIndicator sets the indicator for the SentenceConverter.
func UsingIndicator(indicator IndicatorFunc) SentenceOptFunc {
	return func(converter *SentenceConverter) {
		converter.indicator = indicator
	}
}

// NewSentenceConverter returns a new SentenceConverter.
func NewSentenceConverter(opts ...SentenceOptFunc) *SentenceConverter {
	sent := new(SentenceConverter)

	sent.indicator = wasIndicator
	for _, opt := range opts {
		opt(sent)
	}

	return sent
}

// Sentence returns a copy of the string s in sentence case format.
func (sc *SentenceConverter) Sentence(s string) string {
	var made []string

	ps := `[\p{N}\p{L}*]+[^\s]*`
	if len(sc.vocab) > 0 {
		ps = fmt.Sprintf(`(?:%s)|%s`, strings.Join(sc.vocab, "|"), ps)
	}
	re := regexp.MustCompile(`(?i)` + ps)

	tokens := re.FindAllString(s, -1)
	for i, token := range tokens {
		prev := ""
		if i-1 >= 0 {
			prev = tokens[i-1]
		}

		if entry := sc.inVocab(token); entry != "" {
			made = append(made, entry)
		} else if i == 0 || sc.indicator(prev, i-1) {
			made = append(made, titleCased.String(strings.ToLower(token)))
		} else {
			made = append(made, strings.ToLower(token))
		}
	}

	return strings.Join(made, " ")
}

func (sc *SentenceConverter) inVocab(s string) string {
	for _, token := range sc.vocab {
		if strings.ToLower(token) == strings.ToLower(s) {
			return token
		}
	}
	return ""
}

func wasIndicator(word string, idx int) bool {
	if strings.HasSuffix(word, ":") {
		return true
	} else if idx == 0 && reNumberList.MatchString(word) {
		return true
	}
	return false
}
