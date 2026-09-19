package main

type kind string

const (
	patch kind = "patch"
	minor kind = "minor"
	major kind = "major"
)

func (k kind) isValid() (result bool) {
	switch k {
	case patch, minor, major:
		return true
	}

	return false
}
