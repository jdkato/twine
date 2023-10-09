package strcase_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/jdkato/stransform/internal"
	"github.com/jdkato/stransform/strcase"
)

var testdata = filepath.Join("../testdata")

type testCase struct {
	Input  string
	Expect string
}

func TestAP(t *testing.T) {
	tests := make([]testCase, 0)
	cases := internal.ReadDataFile(filepath.Join(testdata, "AP.json"))

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		t.Error(err)
	}

	tc := strcase.NewTitleConverter(strcase.APStyle)
	for _, test := range tests {
		title := tc.Title(test.Input)
		if test.Expect != title {
			t.Fatalf("Got '%s'; expected '%s'", title, test.Expect)
		}
	}
}

func TestChicago(t *testing.T) {
	tests := make([]testCase, 0)
	cases := internal.ReadDataFile(filepath.Join(testdata, "Chicago.json"))

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		t.Error(err)
	}

	tc := strcase.NewTitleConverter(strcase.ChicagoStyle)
	for _, test := range tests {
		title := tc.Title(test.Input)
		if test.Expect != title {
			t.Fatalf("Got '%s'; expected '%s'", title, test.Expect)
		}
	}
}

func BenchmarkTitle(b *testing.B) {
	tests := make([]testCase, 0)
	cases := internal.ReadDataFile(filepath.Join(testdata, "title.json"))

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		b.Error(err)
	}

	tc := strcase.NewTitleConverter(strcase.APStyle)
	for n := 0; n < b.N; n++ {
		for _, test := range tests {
			_ = tc.Title(test.Input)
		}
	}
}
