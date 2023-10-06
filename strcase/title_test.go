package strcase

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

var testdata = filepath.Join("../testdata")

type testCase struct {
	Input  string
	Expect string
}

func readDataFile(path string) ([]byte, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func TestAP(t *testing.T) {
	tests := make([]testCase, 0)

	cases, err := readDataFile(filepath.Join(testdata, "AP.json"))
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(cases, &tests)
	if err != nil {
		t.Error(err)
	}

	tc := NewTitleConverter(APStyle)
	for _, test := range tests {
		title := tc.Title(test.Input)
		if test.Expect != title {
			t.Fatalf("Got '%s'; expected '%s'", title, test.Expect)
		}
	}
}

func TestChicago(t *testing.T) {
	tests := make([]testCase, 0)

	cases, err := readDataFile(filepath.Join(testdata, "Chicago.json"))
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(cases, &tests)
	if err != nil {
		t.Error(err)
	}

	tc := NewTitleConverter(ChicagoStyle)
	for _, test := range tests {
		title := tc.Title(test.Input)
		if test.Expect != title {
			t.Fatalf("Got '%s'; expected '%s'", title, test.Expect)
		}
	}
}

func BenchmarkTitle(b *testing.B) {
	tests := make([]testCase, 0)

	cases, err := readDataFile(filepath.Join(testdata, "title.json"))
	if err != nil {
		b.Error(err)
	}

	err = json.Unmarshal(cases, &tests)
	if err != nil {
		b.Error(err)
	}

	tc := NewTitleConverter(APStyle)
	for n := 0; n < b.N; n++ {
		for _, test := range tests {
			_ = tc.Title(test.Input)
		}
	}
}
