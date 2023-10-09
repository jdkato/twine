package summarize

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/jdkato/stransform/internal"
)

func TestReadability(t *testing.T) {
	tests := make([]testCase, 0)
	cases := internal.ReadDataFile(filepath.Join(testdata, "summarize.json"))

	err := json.Unmarshal(cases, &tests)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		d := NewDocument(test.Text)
		a := d.Assess()

		if !internal.EqualFloat(test.FleschKincaid, a.FleschKincaid) {
			t.Errorf("FleschKincaid: got %0.2f; expected %0.2f", a.FleschKincaid, test.FleschKincaid)
		}

		if !internal.EqualFloat(test.GunningFog, a.GunningFog) {
			t.Errorf("GunningFog: got %0.2f; expected %0.2f", a.GunningFog, test.GunningFog)
		}

		if !internal.EqualFloat(test.SMOG, a.SMOG) {
			t.Errorf("SMOG: got %0.2f; expected %0.2f", a.SMOG, test.SMOG)
		}

		if !internal.EqualFloat(test.ColemanLiau, a.ColemanLiau) {
			t.Errorf("ColemanLiau: got %0.2f; expected %0.2f", a.ColemanLiau, test.ColemanLiau)
		}
	}
}

func BenchmarkReadability(b *testing.B) {
	in := internal.ReadDataFile(filepath.Join(testdata, "sherlock.txt"))

	d := NewDocument(string(in))
	for n := 0; n < b.N; n++ {
		d.Assess()
	}
}
