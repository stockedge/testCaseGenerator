package tcg

import (
	"reflect"
	"testing"
)

func TestParseFactorInfo(t *testing.T) {
	if !reflect.DeepEqual(ParseFactorInfo(`
A: A1,A2,A3
B: B1, B2
`), map[string][]string{
		"A": {"A1", "A2", "A3"},
		"B": {"B1", "B2"},
	}) {
		t.Fail()
	}
}
