# GoKit CLI  [![Build Status](https://github.com/maohieng/kit-generator/workflows/Go/badge.svg)](https://github.com/maohieng/kit-generator/actions)[![Go Report Card](https://goreportcard.com/badge/github.com/maohieng/kit-generator)](https://goreportcard.com/report/github.com/maohieng/kit-generator)[![Coverage Status](https://coveralls.io/repos/github/GrantZheng/kit/badge.svg?branch=master)](https://coveralls.io/github/GrantZheng/kit?branch=master)

translate to: English | [简体中文](./README_zh.md)  

I fork the project from [kit](https://github.com/kujtimiihoxha/kit) and plan to maintain it in the future. The kit tool is a great job, and deeply used in our team. Some features and bugs have been done and fixed, such as supporting go module,replacing some old dependencies and so on. I am very glad to receive recommend about it.

This project is a more advanced version of [gk](https://github.com/kujtimiihoxha/gk).
The goal of the gokit cli is to be a tool that you can use while you develop your microservices with `gokit`.

While `gk` did help you create your basic folder structure it was not really able to be used further on in your project.
This is what `GoKit Cli` is aiming to change.


# Prerequisites
`Go` is a requirement to be able to test your services.[gokit](https://github.com/go-kit/kit) is needed.To utilise generation of gRPC service code through kit generate service <SERVICE_NAME> -t grpc you will need to install the [grpc prequisites](https://grpc.io/docs/languages/go/quickstart/).
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

# Table of Content
- [Installation](#installation)
- [Usage](#usage)
- [Create a new service](#create-a-new-service)
- [Generate the service](#generate-the-service)
- [Generate the client library](#generate-the-client-library)
- [Generate new middlewares](#generate-new-middleware)
- [Enable docker integration](#enable-docker-integration)

# Installation
Before you install please read [prerequisites](#prerequisites)
```bash
# in the go1.17 or latest
go install github.com/maohieng/kit-generator@latest

# go version =< go1.16 
go install github.com/maohieng/kit-generator
# or
go get -u github.com/maohieng/kit-generator
```


# Usage
```bash
kit help
```

Also read this [medium story](docs/qiuck-start/creating_a_todo_app_using_gokit-cli.md)
# Create a new service
The kit tool use [Go Module](https://go.dev/doc/modules/managing-dependencies) to manage dependencies by default, please make sure your go version >= 1.3, or
GO111MODULE is set on. If you want to specify the module name, you should use the --module flag, otherwise, the module name in the go.mod file will be set as your project name.
```bash
kit new service --help
kit new service hello
kit n s hello # using aliases
```
or
```bash
kit new service hello --module github.com/{group name}/hello
kit n s hello -m github.com/{group name}/hello # using aliases
```

This will generate the initial folder structure, the go.mod file and the service interface

`service-name/pkg/service/service.go`
```go
package service

// HelloService describes the service.
type HelloService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
}
```
When you are generating the service and the client library, the module name in the go.mod file could be autodetected.

# Generate the service
```bash
kit g s hello
kit g s hello --dmw # to create the default middleware
kit g s hello -t grpc # specify the transport (default is http)
```
This command will do these things:
- Create the service boilerplate: `hello/pkg/service/service.go`
- Create the service middleware: `hello/pkg/service/middleware.go`
- Create the endpoint:  `hello/pkg/endpoint/endpoint.go` and `hello/pkg/endpoint/endpoint_gen.go`
- If using` --dmw` create the endpoint middleware: `hello/pkg/endpoint/middleware.go`
- Create the transport files e.x `http`: `service-name/pkg/http/handler.go`
- Create the service main file :boom:   
  `hello/cmd/service/service.go`  
  `hello/cmd/service/service_gen.go`   
  `hello/cmd/main.go`

:warning: **Notice** all the files that end with `_gen` will be regenerated when you add endpoints to your service and
you rerun `kit g s hello` :warning:

You can run the service by running:
```bash
go run hello/cmd/main.go
```

# Generate the client library
```bash
kit g c hello
```
This will generate the client library :sparkles: `http/client/http/http.go` that you can than use to call the service methods, you can use it like this:
```go
package main

import (
	"context"
	"fmt"

	client "hello/client/http"
	"github.com/go-kit/kit/transport/http"
)

func main() {
	svc, err := client.New("http://localhost:8081", map[string][]http.ClientOption{})
	if err != nil {
		panic(err)
	}

	r, err := svc.Foo(context.Background(), "hello")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", r)
}
```
# Generate new middleware
```bash
kit g m hi -s hello
kit g m hi -s hello -e # if you want to add endpoint middleware
```
The only thing left to do is add your middleware logic and wire the middleware with your service/endpoint.
# Enable docker integration

```bash
kit g d
```
This will add the individual service docker files and one `docker-compose.yml` file that will allow you to start
your services.
To start your services just run
```bash
docker-compose up
```

After you run `docker-compose up` your services will start up and any change you make to your code will automatically
rebuild and restart your service (only the service that is changed)
