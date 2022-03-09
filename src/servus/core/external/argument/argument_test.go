package argument

import (
	"os"
	"testing"
)

func Test_Get(t *testing.T) {
	var one = "one"
	var zeroValue = "-zerovalue"
	var two = "two"
	var withValue = "-withvalue"
	var twoDots = "two:dots"
	var notPresented = "not-presented"
	var presented = "presented"
	var flags = []string{
		one,
		zeroValue,
		two,
		withValue,
		twoDots,
		notPresented,
		presented,
	}
	os.Args = []string{one, zeroValue + "=", two, withValue + "=1=2=3=4=5=6", twoDots, presented}
	// get.
	for _, flag := range flags {
		var val = Get(flag)
		//t.Logf("flag: %v", flag)
		if val == nil && flag == notPresented {
			continue
		} else if val == nil {
			t.Fatalf("flag %v failed", flag)
		}
		if val.Value != nil {
			t.Logf("value: %v", *val.Value)
		}
		// get values.
		if val.Value == nil && (flag == withValue || flag == zeroValue) {
			t.Fatalf("flag %v failed", withValue)
		}
	}
}
