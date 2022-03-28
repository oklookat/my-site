package way

import (
	"reflect"
	"testing"
)

func TestShiftPath(t *testing.T) {
	var pathCases = map[string][2]string{
		"///1/2/3/4": {"1", "/2/3/4"},
		"/1/2/3/4//": {"1", "/2/3/4"},

		"1//2/3/4///": {"1", "/2/3/4"},
		"1/2//3/4/":   {"1", "/2/3/4"},
	}
	for key, thePath := range pathCases {
		var head, tail = shiftPath(key)
		var expectedHead = thePath[0]
		if head != expectedHead {
			t.Fatalf(`[head] expected: "%v" / have: "%v"`, expectedHead, head)
		}
		var expectedTail = thePath[1]
		if tail != expectedTail {
			t.Fatalf(`[tail] expected: "%v" / have: "%v"`, expectedTail, tail)
		}
	}
}

func TestRemoveDuplicatesFromSlice(t *testing.T) {
	var slice = []string{"1", "2", "3", "4", "1", "2", "3", "4"}
	var expectedSlice = []string{"1", "2", "3", "4"}
	var result = removeDuplicateValues(slice)
	var isSame = reflect.DeepEqual(result, expectedSlice)
	if !isSame {
		t.Fatal("expected same slices")
	}
}
