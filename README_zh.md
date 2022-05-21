# GoKit CLI  [![Build Status](https://github.com/maohieng/kit-generator/workflows/Go/badge.svg)](https://github.com/maohieng/kit-generator/actions)[![Go Report Card](https://goreportcard.com/badge/github.com/maohieng/kit-generator)](https://goreportcard.com/report/github.com/maohieng/kit-generator)[![Coverage Status](https://coveralls.io/repos/github/GrantZheng/kit/badge.svg?branch=master)](https://coveralls.io/github/GrantZheng/kit?branch=master)

translate to: [English](./README.md) | 简体中文  


本项目fork自[kit](https://github.com/kujtimiihoxha/kit)，并计划在将来维护它。kit是一个优秀的工具，并在我们的团队中得到广泛应用。一些功能和bug已经完成并修复，例如支持go module，替换一些旧的依赖项等，欢迎各位开发者提供建议。

gokit cli 是一个可以在你开发微服务时帮助你的工具，它是[gk](https://github.com/kujtimiihoxha/gk)的高级版。虽然gk确实可以帮助您创建基本的文件夹结构，但它实际上无法在项目中进一步使用，gokit cli希望能够改变这一点。


# Prerequisites
你需要准备：
- [Go](https://go.dev)  
  Go是编译您服务的必要条件，所以你需要先搭建好一套Go环境
- [go-kit](https://github.com/go-kit/kit)  
  gokit cli生成的项目代码使用go-kit作为框架，因此你需要了解go-kit的基本概念以及框架的使用方法
- [Protocol Buffer](https://developers.google.cn/protocol-buffers) 和 [gRPC](https://grpc.io/docs/languages/go/quickstart/)   
  gokit cli使用`kit generate  service <SERVICE_NAME> -t grpc`来生成gRPC代码，所以你需要安装[Protocol Buffer]()和[gRPC需要的环境](https://grpc.io/docs/languages/go/quickstart/)

使用以下命令安装protocol编译器的Go插件
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

# Table of Content
- [安装](#installation)
- [使用方法](#usage)
- [创建一个新服务](#create-a-new-service)
- [生成服务代码](#generate-the-service)
- [生成client代码](#generate-the-client-library)
- [生成中间件](#generate-new-middleware)
- [使用Docker整合](#enable-docker-integration)

# Installation
在你安装之前，请先确保您的环境已经满足 [Prerequisites](#prerequisites)的要求。

```bash
# in the go1.17 or latest
go install github.com/maohieng/kit-generator@latest

# go version =< go1.16 
go install github.com/maohieng/kit-generator
# or
go get -u github.com/maohieng/kit-generator
```


# Usage
若要查看帮助，请使用：
```bash
kit help
```

或查看[这篇文章](docs/qiuck-start/creating_a_todo_app_using_gokit-cli.md)
# Create a new service
`gokit cli`默认使用[Go Module](https://go.dev/doc/modules/managing-dependencies)来管理依赖，请确保你的Go版本 >= 1.3, 或启用GO111MODULE. 如果你想指定module名, 你可以使用 `--module` 选项, 否则 `go.mod`里的模组名将被设置为项目名.
```bash
# 查看生成服务的帮助以及可用选项
kit new service --help

# 生成一个服务，其目录名为hello
kit new service hello

# 使用别名创建服务,等同于kit new service hello
kit n s hello
```
或
```bash
# 生成hello服务并设置module name
kit new service hello --module github.com/{group name}/hello

# 使用别名创建服务并设置module name
kit n s hello -m github.com/{group name}/hello # 
```

这将生成一个初始目录结构,一个`go.mod`文件和一个service 接口文件：

`service-name/pkg/service/service.go`
```go
package service

// HelloService describes the service.
type HelloService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
}
```
当你在生成service或client的时候, 可以自动检测`go.mod`文件中的module name.

# Generate the service
使用以下命令生成service代码：
```bash
# 为hello项目生成代码
kit g s hello

# 为hello项目生成代码并创建默认的中间件(middleware)
kit g s hello --dmw

# 生成代码,指定transport层使用gRPC(默认为http)
kit g s hello -t grpc
```
这个命令会做这些事:
- 创建一个service样板文件: `hello/pkg/service/service.go`
- 创建service中间件: `hello/pkg/service/middleware.go`
- 创建endpoint:  `hello/pkg/endpoint/endpoint.go` and `hello/pkg/endpoint/endpoint_gen.go`
- 如果使用` --dmw`选项，创建endpoint中间件: `hello/pkg/endpoint/middleware.go`
- 创建transport文件，例如`http`: `service-name/pkg/http/handler.go`
- 创建服务main文件:   
  `hello/cmd/service/service.go`  
  `hello/cmd/service/service_gen.go`   
  `hello/cmd/main.go`

:warning: **注意**：当你为你的服务添加endpoint并重新运行`kit g s hello`时， 所有以 `_gen`结尾的文件都会重新生成 :warning:

你可以使用下面的方法来运行你的服务:
```bash
go run hello/cmd/main.go
```

# Generate the client library
```bash
# 生成名为hello的client library项目
kit g c hello
```
这将生成一个client library :sparkles: 你可以使用`http/client/http/http.go`来调用service的方法 ,像这样:
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
要生成新的中间件，你可以：
```bash
# -s选项指定要创建中间件的service名称
kit g m hi -s hello
kit g m hi -s hello -e # 如果你想添加endpoint中间件
```
只需要添加中间件的逻辑，并把中间件和endpoint连接起来即可
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

