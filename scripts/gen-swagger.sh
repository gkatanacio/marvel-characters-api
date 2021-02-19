#!/bin/sh -e

go get -u github.com/swaggo/swag/cmd/swag

$GOPATH/bin/swag init --dir ./cmd/api --parseDependency --parseInternal
