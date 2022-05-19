package utils

import (
	"testing"
	"time"
)

func debounce_simpleTester(withSkip bool) int {
	var counter = 0
	var toDebounce = func() {
		counter++
	}
	var debouncer = Debounce(100*time.Millisecond, withSkip)
	for i := 0; i < 10; i++ {
		debouncer(toDebounce)
		time.Sleep(50 * time.Millisecond)
	}
	return counter
}

func Test_DebounceBasic(t *testing.T) {
	var counter = debounce_simpleTester(false)
	if counter != 5 {
		t.Fatalf("expected 5, got %v", counter)
	}
}

func Test_DebounceBasic_Skip(t *testing.T) {
	var counter = debounce_simpleTester(true)
	if counter != 7 {
		t.Fatalf("expected 7, got %v", counter)
	}
}
