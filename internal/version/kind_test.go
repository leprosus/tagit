package version

import "testing"

func TestKindIsValid(t *testing.T) {
	t.Parallel()

	testList := map[Kind]bool{
		Patch:           true,
		Minor:           true,
		Major:           true,
		Kind("invalid"): false,
	}

	var got bool

	for input, want := range testList {
		got = input.IsValid()
		if got != want {
			t.Errorf(
				"Kind(%q).IsValid() = %t, want %t",
				input,
				got,
				want,
			)
		}
	}
}
