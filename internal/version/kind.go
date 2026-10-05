package version

type Kind string

const (
	Patch Kind = "patch"
	Minor Kind = "minor"
	Major Kind = "major"
)

func (k Kind) IsValid() (result bool) {
	switch k {
	case Patch, Minor, Major:
		return true
	}

	return false
}
