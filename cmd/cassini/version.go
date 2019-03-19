package main

import (
	"context"
	"fmt"

	"github.com/QOSGroup/cassini/config"
)

/**
 * go build --ldflags "-X main.GitCommit=$(git rev-parse HEAD) -X main.Version=0.0.0 " -o ./cassini
 */

var (
	// Version of cassini
	Version = "0.0.5"
	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

var versioner = func(conf *config.Config) (context.CancelFunc, error) {

	fmt.Println("Version: \t", Version,
		"\nGitCommitID: \t", GitCommit)

	return nil, nil
}
