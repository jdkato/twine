// Copyright 2013 Matthew Honnibal
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package nlp

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

// TupleSlice is a slice of tuples in the form (words, tags).
type TupleSlice [][][]string

// Len returns the length of a Tuple.
func (t TupleSlice) Len() int { return len(t) }

// Swap switches the ith and jth elements in a Tuple.
func (t TupleSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// ReadTagged converts pre-tagged input into a TupleSlice suitable for training.
func ReadTagged(text, sep string) TupleSlice {
	lines := strings.Split(text, "\n")
	length := len(lines)
	t := make(TupleSlice, length)
	for i, sent := range lines {
		set := strings.Split(sent, " ")
		length = len(set)
		tokens := make([]string, length)
		tags := make([]string, length)
		for j, token := range set {
			parts := strings.Split(token, sep)
			tokens[j] = parts[0]
			tags[j] = parts[1]
		}
		t[i] = [][]string{tokens, tags}
	}
	return t
}

var none = regexp.MustCompile(`^(?:0|\*[\w?]\*|\*\-\d{1,3}|\*[A-Z]+\*\-\d{1,3}|\*)$`)
var keep = regexp.MustCompile(`^\-[A-Z]{3}\-$`)

// averagedPerceptron is a Averaged Perceptron classifier.
type averagedPerceptron struct {
	classes []string
	stamps  map[string]float64
	totals  map[string]float64
	tagMap  map[string]string
	weights map[string]map[string]float64

	// TODO: Training
	//
	// instances float64
}

// newAveragedPerceptron creates a new AveragedPerceptron model.
func newAveragedPerceptron(weights map[string]map[string]float64,
	tags map[string]string, classes []string) *averagedPerceptron {
	return &averagedPerceptron{
		totals: make(map[string]float64), stamps: make(map[string]float64),
		classes: classes, tagMap: tags, weights: weights}
}

// perceptronTagger is a port of Textblob's "fast and accurate" POS tagger.
// See https://github.com/sloria/textblob-aptagger for details.
type perceptronTagger struct {
	model *averagedPerceptron
}

// newPerceptronTagger creates a new PerceptronTagger and loads the built-in
// AveragedPerceptron model.
func newPerceptronTagger() *perceptronTagger {
	return &perceptronTagger{model: newAveragedPerceptron(wts, tags, classes)}
}

// tag takes a slice of words and returns a slice of tagged tokens.
func (pt *perceptronTagger) tag(tokens []*Token) []*Token {
	var tag string
	var found bool

	p1, p2 := "-START-", "-START2-"
	length := len(tokens) + 4
	context := make([]string, length)
	context[0] = p1
	context[1] = p2
	for i, t := range tokens {
		context[i+2] = normalize(t.Text)
	}
	context[length-2] = "-END-"
	context[length-1] = "-END2-"
	for i := 0; i < len(tokens); i++ {
		word := tokens[i].Text
		if word == "-" {
			tag = "-"
		} else if _, ok := emoticons[word]; ok {
			tag = "SYM"
		} else if strings.HasPrefix(word, "@") {
			// TODO: URLs and emails?
			tag = "NN"
		} else if none.MatchString(word) {
			tag = "-NONE-"
		} else if keep.MatchString(word) {
			tag = word
		} else if tag, found = pt.model.tagMap[word]; !found {
			tag = pt.model.predict(featurize(i, context, word, p1, p2))
		}
		tokens[i].Tag = tag
		p2 = p1
		p1 = tag
	}

	return tokens
}

func (m *averagedPerceptron) predict(features map[string]float64) string {
	var weights map[string]float64
	var found bool

	scores := make(map[string]float64)
	for feat, value := range features {
		if weights, found = m.weights[feat]; !found || value == 0 {
			continue
		}
		for label, weight := range weights {
			scores[label] += value * weight
		}
	}
	return max(scores)
}

func max(scores map[string]float64) string {
	var class string
	max := math.Inf(-1)
	for label, value := range scores {
		if value > max {
			max = value
			class = label
		}
	}
	return class
}

func featurize(i int, ctx []string, w, p1, p2 string) map[string]float64 {
	feats := make(map[string]float64)
	suf := min(len(w), 3)
	i = min(len(ctx)-2, i+2)
	iminus := min(len(ctx[i-1]), 3)
	iplus := min(len(ctx[i+1]), 3)
	feats = add([]string{"bias"}, feats)
	feats = add([]string{"i suffix", w[len(w)-suf:]}, feats)
	feats = add([]string{"i pref1", string(w[0])}, feats)
	feats = add([]string{"i-1 tag", p1}, feats)
	feats = add([]string{"i-2 tag", p2}, feats)
	feats = add([]string{"i tag+i-2 tag", p1, p2}, feats)
	feats = add([]string{"i word", ctx[i]}, feats)
	feats = add([]string{"i-1 tag+i word", p1, ctx[i]}, feats)
	feats = add([]string{"i-1 word", ctx[i-1]}, feats)
	feats = add([]string{"i-1 suffix", ctx[i-1][len(ctx[i-1])-iminus:]}, feats)
	feats = add([]string{"i-2 word", ctx[i-2]}, feats)
	feats = add([]string{"i+1 word", ctx[i+1]}, feats)
	feats = add([]string{"i+1 suffix", ctx[i+1][len(ctx[i+1])-iplus:]}, feats)
	feats = add([]string{"i+2 word", ctx[i+2]}, feats)
	return feats
}

func add(args []string, features map[string]float64) map[string]float64 {
	key := strings.Join(args, " ")
	features[key]++
	return features
}

func normalize(word string) string {
	if word == "" {
		return word
	}
	first := string(word[0])
	if strings.Contains(word, "-") && first != "-" {
		return "!HYPHEN"
	} else if _, err := strconv.Atoi(word); err == nil && len(word) == 4 {
		return "!YEAR"
	} else if _, err := strconv.Atoi(first); err == nil {
		return "!DIGITS"
	}
	return strings.ToLower(word)
}
