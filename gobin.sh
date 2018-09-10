#! /usr/bash

if  command -v go-bindata > /dev/null; then
    echo "go-bindata found."
else
    echo "go-bindata not found. downloading go-bindata..."
    go get -u github.com/go-bindata/go-bindata/...
fi

echo "go-bindata running..."
go-bindata res
echo "go-bindata end"
