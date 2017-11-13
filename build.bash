#!/bin/bash
echo "----------------------"
echo "executing tests"
echo "----------------------"
go test -cover ./...
echo "----------------------"
echo "building binary"
echo "----------------------"
go build -o asdf
echo "done"