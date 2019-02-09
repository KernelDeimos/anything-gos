package permutations

import (
	"testing"
)

/*
0000001 -> 0000001
1000001 -> 1000100
1101010 -> 1011001


1101010
1011001
1101010
1011001
*/

func TestToRelative(t *testing.T) {
	input := PermuSeq{
		true,
		true, false,
		true, false, true, false,
	}
	expected := PermuSeq{
		true,
		false, true,
		true, false, false, true,
	}
	result, err := input.ToRelative()
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(expected) {
		t.Fatal("len != len")
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("index %d fails", i)
		}
	}
}
