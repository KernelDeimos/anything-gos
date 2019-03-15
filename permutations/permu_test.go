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

func TestPermuteAbsolute(t *testing.T) {
	key := PermuSeq{
		true,
		true, false,
		true, false, true, false,
	}

	input := []interface{}{1, 2, 3, 4, 5, 6, 7, 8}

	expected := []interface{}{8, 7, 5, 6, 2, 1, 3, 4}

	output, err := key.Permutate(input)
	if err != nil {
		t.Error(err)
	}

	for i, v := range expected {
		if v != output[i] {
			t.Log(output)
			t.FailNow()
		}
	}

}

func TestPermuteRelative(t *testing.T) {
	key := PermuSeq{
		true,
		true, false,
		true, false, true, false,
	}

	input := []interface{}{1, 2, 3, 4, 5, 6, 7, 8}

	expected := []interface{}{6, 5, 7, 8, 3, 4, 2, 1}

	key, err := key.ToRelative()
	if err != nil {
		t.Error(err)
	}

	output, err := key.Permutate(input)
	if err != nil {
		t.Error(err)
	}

	for i, v := range expected {
		if v != output[i] {
			t.Log(output)
			t.FailNow()
		}
	}

}
