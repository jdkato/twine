package tokenize

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jdkato/twine/internal"
)

// A Token represents an individual token of text such as a word or punctuation
// symbol.
type Token struct {
	Text string // The token's actual content.
}

type TokenTester func(string) bool

type Tokenizer interface {
	Tokenize(string) []*Token
}

// iterTokenizer splits a sentence into words.
type iterTokenizer struct {
	specialRE      *regexp.Regexp
	sanitizer      *strings.Replacer
	contractions   []string
	splitCases     []string
	suffixes       []string
	prefixes       []string
	emoticons      map[string]int
	isUnsplittable TokenTester
	noSuffix       bool
}

type TokenizerOptFunc func(*iterTokenizer)

// UsingIsUnsplittableFN gives a function that tests whether a token is splittable or not.
func UsingIsUnsplittable(x TokenTester) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.isUnsplittable = x
	}
}

// Use the provided special regex for unsplittable tokens.
func UsingSpecialRE(x *regexp.Regexp) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.specialRE = x
	}
}

// Use the provided sanitizer.
func UsingSanitizer(x *strings.Replacer) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.sanitizer = x
	}
}

// Use the provided suffixes.
func UsingSuffixes(x []string) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.suffixes = x
	}
}

// Use the provided prefixes.
func UsingPrefixes(x []string) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.prefixes = x
	}
}

// Use the provided map of emoticons.
func UsingEmoticons(x map[string]int) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.emoticons = x
	}
}

// Use the provided contractions.
func UsingContractions(x []string) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.contractions = x
	}
}

// Use the provided splitCases.
func UsingSplitCases(x []string) TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.splitCases = x
	}
}

func WithoutSuffix() TokenizerOptFunc {
	return func(tokenizer *iterTokenizer) {
		tokenizer.noSuffix = true
	}
}

// Constructor for default iterTokenizer
func NewIterTokenizer(opts ...TokenizerOptFunc) *iterTokenizer {
	tok := new(iterTokenizer)

	// Set default parameters
	tok.contractions = contractions
	tok.emoticons = emoticons
	tok.isUnsplittable = func(_ string) bool { return false }
	tok.prefixes = prefixes
	tok.sanitizer = sanitizer
	tok.specialRE = internalRE
	tok.suffixes = suffixes
	tok.noSuffix = false

	// Apply options if provided
	for _, applyOpt := range opts {
		applyOpt(tok)
	}

	tok.splitCases = append(tok.splitCases, tok.contractions...)

	return tok
}

func addToken(s string, toks []string) []string {
	if strings.TrimSpace(s) != "" {
		toks = append(toks, s)
	}
	return toks
}

func (t *iterTokenizer) isSpecial(token string) bool {
	_, found := t.emoticons[token]
	return found || t.specialRE.MatchString(token) || t.isUnsplittable(token)
}

func (t *iterTokenizer) doSplit(token string) []string {
	tokens := []string{}
	suffs := []string{}

	last := 0
	for token != "" && utf8.RuneCountInString(token) != last {
		if t.isSpecial(token) {
			// We've found a special case (e.g., an emoticon) -- so, we add it as a token without
			// any further processing.
			tokens = addToken(token, tokens)
			break
		}
		last = utf8.RuneCountInString(token)
		lower := strings.ToLower(token)
		if internal.HasAnyPrefix(token, t.prefixes) {
			// Remove prefixes -- e.g., $100 -> [$, 100].
			tokens = addToken(string(token[0]), tokens)
			token = token[1:]
		} else if idx := internal.HasAnyIndex(lower, t.splitCases); idx > -1 {
			// Handle "they'll", "I'll", "Don't", "won't", amount($).
			//
			// they'll -> [they, 'll].
			// don't -> [do, n't].
			// amount($) -> [amount, (, $, )].
			tokens = addToken(token[:idx], tokens)
			token = token[idx:]
		} else if internal.HasAnySuffix(token, t.suffixes) {
			// Remove suffixes -- e.g., Well) -> [Well, )].
			suffs = append([]string{string(token[len(token)-1])}, suffs...)
			token = token[:len(token)-1]
		} else {
			tokens = addToken(token, tokens)
		}
	}

	return append(tokens, suffs...)
}

func (t *iterTokenizer) doSplitNoSuffix(token string) []string {
	var tokens []string

	last := 0
	for token != "" && utf8.RuneCountInString(token) != last {
		if t.isSpecial(token) {
			// We've found a special case (e.g., an emoticon) -- so, we add it as a token without
			// any further processing.
			tokens = addToken(token, tokens)
			break
		}
		last = utf8.RuneCountInString(token)
		lower := strings.ToLower(token)
		if internal.HasAnyPrefix(token, t.prefixes) {
			// Remove prefixes -- e.g., $100 -> [$, 100].
			token = token[1:]
		} else if idx := internal.HasAnyIndex(lower, t.splitCases); idx > -1 {
			// Handle "they'll", "I'll", "Don't", "won't", amount($).
			//
			// they'll -> [they, 'll].
			// don't -> [do, n't].
			// amount($) -> [amount, (, $, )].
			tokens = addToken(token[:idx], tokens)
			token = token[idx:]
		} else if internal.HasAnySuffix(token, t.suffixes) {
			// Remove suffixes -- e.g., Well) -> [Well, )].
			token = token[:len(token)-1]
		} else {
			tokens = addToken(token, tokens)
		}
	}

	return tokens
}

// tokenize splits a sentence into a slice of words.
func (t *iterTokenizer) Tokenize(text string) []string {
	var tokens []string

	clean, white := t.sanitizer.Replace(text), false
	length := len(clean)

	start, index := 0, 0
	cache := map[string][]string{}
	for index <= length {
		uc, size := utf8.DecodeRuneInString(clean[index:])
		if size == 0 {
			break
		} else if index == 0 {
			white = unicode.IsSpace(uc)
		}
		if unicode.IsSpace(uc) != white {
			if start < index {
				span := clean[start:index]
				if toks, found := cache[span]; found {
					tokens = append(tokens, toks...)
				} else {
					toks := []string{}
					if t.noSuffix {
						toks = t.doSplitNoSuffix(span)
					} else {
						toks = t.doSplit(span)
					}
					cache[span] = toks
					tokens = append(tokens, toks...)
				}
			}
			if uc == ' ' {
				start = index + 1
			} else {
				start = index
			}
			white = !white
		}
		index += size
	}

	if start < index {
		tokens = append(tokens, t.doSplit(clean[start:index])...)
	}

	return tokens
}

var internalRE = regexp.MustCompile(`^(?:[A-Za-z]\.){2,}$|^[A-Z][a-z]{1,2}\.$`)
var sanitizer = strings.NewReplacer(
	"\u201c", `"`,
	"\u201d", `"`,
	"\u2018", "'",
	"\u2019", "'",
	"&rsquo;", "'")
var contractions = []string{"'ll", "'s", "'re", "'m", "n't"}
var suffixes = []string{",", ")", `"`, "]", "!", ";", ".", "?", ":", "'"}
var prefixes = []string{"$", "(", `"`, "["}
var emoticons = map[string]int{
	"(-8":         1,
	"(-;":         1,
	"(-_-)":       1,
	"(._.)":       1,
	"(:":          1,
	"(=":          1,
	"(o:":         1,
	"(¬_¬)":       1,
	"(ಠ_ಠ)":       1,
	"(╯°□°）╯︵┻━┻": 1,
	"-__-":        1,
	"8-)":         1,
	"8-D":         1,
	"8D":          1,
	":(":          1,
	":((":         1,
	":(((":        1,
	":()":         1,
	":)))":        1,
	":-)":         1,
	":-))":        1,
	":-)))":       1,
	":-*":         1,
	":-/":         1,
	":-X":         1,
	":-]":         1,
	":-o":         1,
	":-p":         1,
	":-x":         1,
	":-|":         1,
	":-}":         1,
	":0":          1,
	":3":          1,
	":P":          1,
	":]":          1,
	":`(":         1,
	":`)":         1,
	":`-(":        1,
	":o":          1,
	":o)":         1,
	"=(":          1,
	"=)":          1,
	"=D":          1,
	"=|":          1,
	"@_@":         1,
	"O.o":         1,
	"O_o":         1,
	"V_V":         1,
	"XDD":         1,
	"[-:":         1,
	"^___^":       1,
	"o_0":         1,
	"o_O":         1,
	"o_o":         1,
	"v_v":         1,
	"xD":          1,
	"xDD":         1,
	"¯\\(ツ)/¯":    1,
}
