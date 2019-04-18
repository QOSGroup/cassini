#!/bin/sh

go build --ldflags "-X main.GitCommit=$(git rev-parse HEAD) -X main.Version=0.5.0 " -o ./cassini
