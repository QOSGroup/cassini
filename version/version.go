package version

/**
 * go build --ldflags "-X version.GitCommit=$(git rev-parse HEAD) -X version.Version=0.0.0 " -o ./cassini
 */

var (
	// Version of cassini
	Version = "0.0.5"
	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)
