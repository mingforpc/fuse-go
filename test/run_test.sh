#!/bin/sh
go test -cover -coverprofile coverage.cov -coverpkg=../fuse
go tool cover -html=./coverage.cov -o ./coverage.html