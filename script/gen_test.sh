#!/bin/bash
serviceCode='package service

import "context"

type HelloService interface {
    Foo(ctx context.Context,s string)(rs string, err error)
    Bar(ctx context.Context,i int)(rs int, err error)
}'


if [ -d temp ]; then
    cd temp || exit
else
    mkdir temp && cd temp || exit
fi
# generate project
kit n s hello
echo "$serviceCode" > ./hello/pkg/service/service.go
kit g s hello
kit g s hello -t grpc
cd ./hello || exit
go mod tidy -compat=1.17
go build -v ./...
