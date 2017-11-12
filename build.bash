#!/bin/bash
echo "----------------------"
echo "executing tests"
echo "----------------------"
go test -cover ./...
echo "----------------------"
echo "building binary"
echo "----------------------"
go build -o asdf github.com/moolen/asdf/cmd/asdf
echo "done"