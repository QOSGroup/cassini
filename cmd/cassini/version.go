package main

import (
	"context"
	"fmt"

	"github.com/QOSGroup/cassini/commands"
)

/**
 * go build --ldflags "-X main.GitCommit=$(git rev-parse HEAD) -X main.Version=0.0.0 " -o ./cassini
 */

// nolint
var (
	Version   = "0.0.6"
	GitCommit string
	GoVersion string
	BuidDate  string
)

var versioner = func() (context.CancelFunc, error) {

	s := `cassini - %s
version:	%s
revision:	%s
compile:	%s
go version:	%s
`

	fmt.Printf(s, commands.ShortDescription,
		Version, GitCommit, BuidDate, GoVersion)

	return nil, nil
}
