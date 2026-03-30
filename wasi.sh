#!/bin/sh

outname=./bytes2hex_cli.wasm
mainpat=./cmd/bytes2hex/main.go

GOOS=wasip1 GOARCH=wasm go \
	build \
	-o "${outname}" \
	-ldflags="-s -w" \
	"${mainpat}"
