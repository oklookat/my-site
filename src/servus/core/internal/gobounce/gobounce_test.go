package gobounce

import (
	"testing"
	"time"
)

func simpleTester(withSkip bool) int {
	var counter = 0
	var toDebounce = func() {
		counter++
	}
	var debouncer = New(100*time.Millisecond, withSkip)
	for i := 0; i < 10; i++ {
		debouncer(toDebounce)
		time.Sleep(50 * time.Millisecond)
	}
	return counter
}

func TestBasic(t *testing.T) {
	var counter = simpleTester(false)
	if counter != 5 {
		t.Fatalf("expected 5, got %v", counter)
	}
}

func TestBasicSkip(t *testing.T) {
	var counter = simpleTester(true)
	if counter != 7 {
		t.Fatalf("expected 7, got %v", counter)
	}
}
