package strcase

import (
	"fmt"
	"strings"

	"github.com/errata-ai/regexp2"
	"github.com/jdkato/twine/internal"
)

var reNumberList = regexp2.MustCompileStd(`\d+\.`)
var defaultSentOpts = CaseOpts{
	vocab: []string{},
	indicator: func(word string, idx int) bool {
		if strings.HasSuffix(word, ":") {
			return true
		} else if idx == 0 && reNumberList.MatchStringStd(word) {
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

	ps := `[\p{N}\p{L}*]+[^\s]*`
	if len(sc.vocab) > 0 {
		ps = fmt.Sprintf(`\b(?:%s)\b|%s`, strings.Join(sc.vocab, "|"), ps)
	}
	re := regexp2.MustCompileStd(`(?i)` + ps)

	tokens := re.FindAllString(s, -1)
	// NOTE: We have to do this *after* tokenizing the string in order to
	// respect the case of would-be exceptions.
	s = strings.ToLower(s)

	for i, token := range tokens {
		prev := ""
		if i-1 >= 0 {
			prev = tokens[i-1]
		}

		if entry := sc.inVocab(token); entry != "" {
			made = append(made, entry)
		} else if i == 0 || sc.indicator(prev, i-1) {
			made = append(made, internal.ToTitle(token, true))
		} else {
			made = append(made, strings.ToLower(token))
		}
	}

	return strings.Join(made, " ")
}

func (sc *SentenceConverter) inVocab(s string) string {
	for _, token := range sc.vocab {
		matched, _ := regexp2.MatchString(token, s)
		if strings.ToLower(token) == strings.ToLower(s) {
			return token
		} else if matched {
			return s
		}
	}
	return ""
}
