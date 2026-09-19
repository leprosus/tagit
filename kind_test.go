package main

import "testing"

func TestKindIsValid(t *testing.T) {
	t.Parallel()

	testList := map[kind]bool{
		patch:           true,
		minor:           true,
		major:           true,
		kind("invalid"): false,
	}

	var got bool

	for input, want := range testList {
		got = input.isValid()
		if got != want {
			t.Errorf("kind(%q).isValid() = %t, want %t", input, got, want)
		}
	}
}
