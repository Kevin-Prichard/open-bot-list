package util

import (
	"testing"
)

func TestPTRLookup(t *testing.T) {
	got := LookupPTR("1.1.1.1")
	want := "one.one.one.one"

	if want != got {
		t.Errorf("LookupPTR() got = %+v, want %+v", got, want)
	}
}
