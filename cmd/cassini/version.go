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
	Version = "0.0.6"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string

	// GoVersion is version info of golang
	GoVersion string

	// BuidDate is compile date and time
	BuidDate string
)

var versioner = func(conf *config.Config) (context.CancelFunc, error) {

	s := `cassini - the relay of cross-chain
version:	%s
revision:	%s
compile:	%s
go version:	%s
`

	fmt.Printf(s, Version, GitCommit, BuidDate, GoVersion)

	return nil, nil
}
