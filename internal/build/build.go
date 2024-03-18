package build

var (
	sha     = "dev"
	version = "v0.0.0"
)

type Info struct {
	// SHA is the ref that built this binary
	// it is dev if it is not set.
	SHA string
	// Version is the semver tag
	// it is v0.0.0 if it is not set
	Version string
}

// Version returns the build info which is set via the build system.
func Version() Info {
	return Info{
		SHA:     sha,
		Version: version,
	}
}
