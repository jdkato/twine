package strcase

// A CaseConverter converts a string to a specific case.
type CaseConverter interface {
	Convert(string) string
}

// CaseOptFunc is a function that modifies a CaseConverter.
type CaseOptFunc func(opts *CaseOpts)

// An IndicatorFunc is a SentenceConverter callback that decides whether or not
// the string word should be capitalized.
type IndicatorFunc func(word string, idx int) bool

// CaseOpts is a struct that holds the options for a CaseConverter.
type CaseOpts struct {
	vocab     []string
	indicator IndicatorFunc
}

// UsingVocab sets the vocab for the CaseConverter.
func UsingVocab(vocab []string) CaseOptFunc {
	return func(opts *CaseOpts) {
		opts.vocab = vocab
	}
}

// UsingIndicator sets the indicator for the CaseConverter.
func UsingIndicator(indicator IndicatorFunc) CaseOptFunc {
	return func(opts *CaseOpts) {
		opts.indicator = indicator
	}
}
