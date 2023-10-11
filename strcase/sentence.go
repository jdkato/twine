package strcase

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jdkato/twine/internal"
)

var reNumberList = regexp.MustCompile(`\d+\.`)
var defaultSentOpts = CaseOpts{
	vocab: []string{},
	indicator: func(word string, idx int) bool {
		if strings.HasSuffix(word, ":") {
			return true
		} else if idx == 0 && reNumberList.MatchString(word) {
			return true
		}
		return false
	},
}

// A SentenceConverter converts a string to sentence case.
type SentenceConverter struct {
	CaseOpts
}

// NewSentenceConverter returns a new SentenceConverter.
func NewSentenceConverter(opts ...CaseOptFunc) *SentenceConverter {
	sent := new(SentenceConverter)

	base := defaultSentOpts
	for _, opt := range opts {
		opt(&base)
	}

	sent.vocab = base.vocab
	sent.indicator = base.indicator

	return sent
}

// Convert returns a copy of the string s in sentence case format.
func (sc *SentenceConverter) Convert(s string) string {
	var made []string

	s = strings.ToLower(s)

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
			made = append(made, internal.ToTitle(token))
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
