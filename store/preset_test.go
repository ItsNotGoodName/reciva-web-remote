package store

import (
	"testing"
)

type GetPresetTest struct {
	uri   string
	sid   int
	valid bool
}

func getPresetTests() []GetPresetTest {
	var getPresetTests = []GetPresetTest{
		{"/04.m3u", 0, false},
	}

	for _, v := range PresetURIs {
		getPresetTests = append(getPresetTests, GetPresetTest{v, 0, true})
	}
	return getPresetTests
}

func TestPreset(t *testing.T) {
	s := testStore(t)

	// GetPreset
	name := "GetPreset"
	var validGetPresetTests []GetPresetTest
	for _, gpt := range getPresetTests() {
		pt, ok := s.GetPreset(gpt.uri)
		if ok != gpt.valid {
			t.Errorf("%s(ok) = %t, want %t", name, ok, gpt.valid)
		}
		if ok {
			if pt.SID != gpt.sid {
				t.Errorf("%s(pt.SID) = %d, want %d", name, pt.SID, gpt.sid)
			}
			if pt.URI != gpt.uri {
				t.Errorf("%s(pt.URI) = %s, want %s", name, pt.URI, gpt.uri)
			}
			validGetPresetTests = append(validGetPresetTests, gpt)
		}
	}

	// GetPresets
	name = "GetPresets"
	if pts := s.GetPresets(); len(pts) != len(validGetPresetTests) {
		t.Fatalf("len(%s) = %d, want %d", name, len(pts), len(validGetPresetTests))
	}
}
