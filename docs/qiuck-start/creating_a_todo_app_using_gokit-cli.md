Creating a TODO app using GoKit-CLI
===

> :warning: This article was published by kujtim Hoxha on October 15, 2017. The original article is [here](https://medium.com/@kujtimii.h/creating-a-todo-app-using-gokit-cli-20f066a58e1). Due to the long time, some contents in this article may no longer be completely accurate. We are working to update this article.  


Creating a microservices based application for a TODO app is probably an overkill, but I think it will serve as a good example of how you can use Go kit and GoKit-CLI to make it a much easier and faster process.  

## What should the TODO app do?
1. It should show existing todos.
2. It should allow users to add new todos.
3. It should allow users to mark todo’s as done.
4. It should allow users to delete todos.

## Let’s start with the basics.
To create the app I am going to use a great Go library that will make our lives easier called Go kit (to learn more about Go kit visit https://gokit.io/), besides that, I have also created a CLI tool that helps with some boilerplate code and will speed up the process.  
For the frontend I will use an existing Angular demo app and modify it to use our services.

You can find the original version at [jvandemo/angular2-todo-app](https://github.com/jvandemo/angular2-todo-app).  
And the modified one at [kujtimiihoxha/todo-demo](https://github.com/kujtimiihoxha/todo-demo)

I will not go in-depth with how the Angular app works, since it only servers as a means to visualize what we are doing.

I assume that you already have go installed and setup so start by getting the necessary libraries.

Install GoKit-CLI:
```shell
go get github.com/kujtimiihoxha/kit
```
The service will also use a package manager called `glide` so let’s install it also:
```shell
curl https://glide.sh/get | sh
```

More about glide: [Masterminds/glide](https://github.com/Masterminds/glide)

GoKit-CLI will only work if your project is inside your `$GOPATH` so create a project folder and `cd` into it with your terminal.

Inside that directory initiate `glide` :
```shell
glide init
```

Now install Go kit:
```shell
glide get github.com/go-kit/kit
```

## Create a new service using `kit`

```shell
kit new service todo
kit n s todo # using aliases
```
This will generate the main folder structure and the service interface:
```text
todo/
|---pkg/
|------service/
|----------service.go
```

### service.go
```go
package service

// TodoService describes the service.
type TodoService interface {
   // Add your methods here
   // e.x: Foo(ctx context.Context,s string)(rs string, err error)
}
```
In this interface we will define all of our service endpoints, first, let’s create a structure that will represent a `todo` .
```text
todo/
|---pkg/
|------io/
|----------io.go
```
### io.go
```go
package io

import "gopkg.in/mgo.v2/bson"

type Todo struct {
   Id       bson.ObjectId `json:"id" bson:"_id"`
   Title    string        `json:"title"  bson:"title"`
   Complete bool          `json:"complete" bson:"complete"`
}

func (t Todo) String() string {
   b, err := json.Marshal(t)
   if err != nil {
      return "unsupported value type"
   }
   return string(b)
}
```
We will be using mongo for this example that’s why we use a `bson.ObjectId` as the todo ID.

Now define our endpoints.

```go
import (
   "context"

   "github.com/kujtimiihoxha/todo-gokit-demo/todo/pkg/io"
)

// TodoService describes the service.
type TodoService interface {
   Get(ctx context.Context) (t []io.Todo, error error)
   Add(ctx context.Context, todo io.Todo) (t io.Todo, error error)
   SetComplete(ctx context.Context, id string) (error error)
   RemoveComplete(ctx context.Context, id string) (error error)
   Delete(ctx context.Context, id string) (error error)
}
```

This is all that `kit` needs to create the service.

```shell
kit g s todo -w --gorilla
```

`-w` generate some default service middleware.

`--gorilla` use [gorilla/mux](https://github.com/gorilla/mux) instead of the default http handler for the http transport.

`kit` will generate the following file structure

```text
todo/
|---cmd/
|------service/
|----------server.go          Wire the service.
|----------server_gen.go      Also wire the service.
|------main.go                Runs the service
|---pkg/
|------endpoints/
|----------endpoint.go        The endpoint logic.
|----------endpoint_gen.go    This will wire the endpoints.
|----------middleware.go      Endpoint middleware
|------http/
|----------handler.go         Transport logic encode/decode data.
|----------handler_gen.go     This will wire the transport.
|------io/
|----------io.go              The input output structs.
|------service/
|----------middleware.go      The service middleware.
|----------service.go         Business logic.
```

As it is right now the service can be run without any code change.

```shell
glide get github.com/gorilla/mux
glide get github.com/gorilla/handlers
```

then

```shell
go run todo/cmd/main.go
```

service starts

```text
ts=2017-10-14T16:51:40.595461383Z caller=service.go:78 tracer=none
ts=2017-10-14T16:51:40.595866216Z caller=service.go:100 transport=HTTP addr=:8081
ts=2017-10-14T16:51:40.595893015Z caller=service.go:134 transport=debug/HTTP addr=:8080
```

But the service wont do to much without implementing the business logic, so lets install the MongoDB driver and move on from there.

```shell
glide get gopkg.in/mgo.v2
```

Create our session:

```text
todo/
|---pkg/
|------db/
|----------db.go
```

### db.go
```go
package db

import (
   mgo "gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session
var mongo_conn_str = "mongodb:27017"

// Creates a new session if mgoSession is nil i.e there is no active mongo session.
//If there is an active mongo session it will return a Clone
func GetMongoSession() (*mgo.Session, error) {
   if mgoSession == nil {
      var err error
      mgoSession, err = mgo.Dial(mongo_conn_str)
      if err != nil {
         return nil, err
      }
   }
   return mgoSession.Clone(), nil
}
```

Now lets edit all the service endpoints:

### Get
`kit` will create the `Get` endpoint implementation in `todo/pkg/service/service.go` like:

```go
func (b *basicTodoService) Get(ctx context.Context) (t []io.Todo, error error) {
   // TODO implement the business logic of Get
   return t, error
}
```

Now lets implement it:

```go
func (b *basicTodoService) Get(ctx context.Context) (t []io.Todo, error error) {
   session, err := db.GetMongoSession()
   if err != nil {
      return t, err
   }
   defer session.Close()
   c := session.DB("todo_app").C("todos")
   error = c.Find(nil).All(&t)
   return t, error
}
```

We also need to update `todo/http/handler.go` so we expect a GET request and that the request is decoded properly:

change:

```go
// makeGetHandler creates the handler logic
func makeGetHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
   m.Methods("POST").Path("/get").Handler(
      handlers.CORS(
         handlers.AllowedMethods([]string{"POST"}),
         handlers.AllowedOrigins([]string{"*"}),
      )(http.NewServer(endpoints.GetEndpoint, decodeGetRequest, encodeGetResponse, options...)))
}
```

to:

```go
// makeGetHandler creates the handler logic
func makeGetHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
   m.Methods("GET").Path("/").Handler(
      handlers.CORS(
         handlers.AllowedMethods([]string{"GET"}),
         handlers.AllowedOrigins([]string{"*"}),
      )(http.NewServer(endpoints.GetEndpoint, decodeGetRequest, encodeGetResponse, options...)),
   )
}
```

and:

```go
// decodeGetResponse  is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetRequest(_ context.Context, r *http1.Request) (interface{}, error) {
   req := endpoint.GetRequest{}
   err := json.NewDecoder(r.Body).Decode(&req)
   return req, nil
}
```

to:

```go
// decodeGetResponse  is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetRequest(_ context.Context, r *http1.Request)(interface{}, error) {
   req := endpoint.GetRequest{}
   return req, nil
}
```

(editing the decoder is necessary because the request does not have any body)

### Add
Same as for the `Get` endpoint `kit` will implement an empty function for `Add` .  
Implement it like:

```go
func (b *basicTodoService) Add(ctx context.Context, todo io.Todo) (t io.Todo, error error) {
   todo.Id = bson.NewObjectId()
   session, err := db.GetMongoSession()
   if err != nil {
      return t, err
   }
   defer session.Close()
   c := session.DB("todo_app").C("todos")
   error = c.Insert(&todo)
   return todo, error
}
```

and update the handler:

```go
// makeAddHandler creates the handler logic
func makeAddHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
   m.Methods("POST", "OPTIONS").Path("/add").Handler(
      handlers.CORS(
         handlers.AllowedOrigins([]string{"*"}),
         handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
         handlers.AllowedMethods([]string{"POST"}),
      )(http.NewServer(endpoints.AddEndpoint, decodeAddRequest, encodeAddResponse, options...)))
}
```

So we accept `POST` and `OPTIONS` method because our request will be a `CORS` request.

We also need to add any additional `Header` we allow using `AllowHeaders`.

### SetComplete

Implementation:
```go
func (b *basicTodoService) SetComplete(ctx context.Context, id string) (error error) {
   session, err := db.GetMongoSession()
   if err != nil {
      return err
   }
   defer session.Close()
   c := session.DB("todo_app").C("todos")
   return c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"complete": true}})
}
```

Handler:
```go
// makeSetCompleteHandler creates the handler logic
func makeSetCompleteHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
   m.Methods("PUT", "OPTIONS").Path("/set-complete").Handler(
      handlers.CORS(
         handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
         handlers.AllowedMethods([]string{"PUT"}),
         handlers.AllowedOrigins([]string{"*"}),
      )(http.NewServer(endpoints.SetCompleteEndpoint, decodeSetCompleteRequest, encodeSetCompleteResponse, options...)))
}
```

### RemoveComplete
Implementation:
```go
func (b *basicTodoService) RemoveComplete(ctx context.Context, id string) (error error) {
   session, err := db.GetMongoSession()
   if err != nil {
      return err
   }
   defer session.Close()
   c := session.DB("todo_app").C("todos")
   return c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"complete": false}})
}
```
Handler:
```go
// makeRemoveCompleteHandler creates the handler logic
func makeRemoveCompleteHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
   m.Methods("PUT", "OPTIONS").Path("/remove-complete").Handler(
      handlers.CORS(
         handlers.AllowedMethods([]string{"PUT"}),
         handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
         handlers.AllowedOrigins([]string{"*"}),
      )(http.NewServer(endpoints.RemoveCompleteEndpoint, decodeRemoveCompleteRequest, encodeRemoveCompleteResponse, options...)))
}
```
### Delete
Implementation:
```go
func (b *basicTodoService) Delete(ctx context.Context, id string) (error error) {
   session, err := db.GetMongoSession()
   if err != nil {
      return err
   }
   defer session.Close()
   c := session.DB("todo_app").C("todos")
   return c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
```

Handler:

```go
// makeDeleteHandler creates the handler logic
func makeDeleteHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http1.ServerOption) {
   m.Methods("DELETE", "OPTIONS").Path("/delete/{id}").Handler(
      handlers.CORS(
         handlers.AllowedMethods([]string{"DELETE"}),
         handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
         handlers.AllowedOrigins([]string{"*"}),
      )(http1.NewServer(endpoints.DeleteEndpoint, decodeDeleteRequest, encodeDeleteResponse, options...)))
}
```

we also want to update the decoder for delete since we will give the id of the todo in the url `/delete/{id}`
```go
// decodeDeleteResponse  is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeDeleteRequest(_ context.Context, r *http1.Request) (interface{}, error) {
   vars := mux.Vars(r)
   id, ok := vars["id"]
   if !ok {
      return nil, errors.New("not a valid ID")
   }
   req := endpoint.DeleteRequest{
      Id: id,
   }
   return req, nil
}
```

## Running the service
Now the service is ready to be used. If you want to start the service its pretty easy `kit` has already created the `main` file and has wired your service together, you just have to run it:
```shell
go run todo/cmd/main.go
```

One thing to keep in mind is that you need to startup a mongo server and specify the address in `db.go`

## Using kit and docker
An easier way to startup your services is using `kit` and `docker` . Kit has a build in command that will generate your services `Dockerfile` and it will create a `docker-compose.yml` configuration file so you can use `docker-compose` and run your containers.

Run the generate docker command from the project folder.

```shell
kit g d --glide # we use the glide flag so that we use glide dependency manager.
```

This will generate:

### Dockerfile

```dockerfile
FROM golang

RUN mkdir -p /go/src/github.com/kujtimiihoxha/todo-gokit-demo

ADD ../.. /go/src/github.com/kujtimiihoxha/todo-gokit-demo

RUN curl https://glide.sh/get | sh
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

RUN cd /go/src/github.com/kujtimiihoxha/todo-gokit-demo && glide install

ENTRYPOINT  watcher -run github.com/kujtimiihoxha/todo-gokit-demo/todo/cmd -watch github.com/kujtimiihoxha/todo-gokit-demo/todo
```

This dockerfile is setup so that it will use `watcher` a nice little tool that will rebuild your services every time you change something in your source this makes development much easier and you don’t have to constantly restart your containers.

### Docker Compose Config

```yaml
version: "2"
services:
  todo:
    build:
      context: .
      dockerfile: todo/Dockerfile
    restart: always
    volumes:
    - .:/go/src/github.com/kujtimiihoxha/todo-gokit-demo/examples
    container_name: todo
    ports:
    - 8800:8081
  mongodb:
    command: mongod --smallfiles --logpath=/dev/null
    container_name: mongodb
    environment:
    - MONGO_DATA_DIR=/data/db
    - MONGO_LOG_DIR=/dev/null
    image: mongo:latest
    ports:
    - 27017:27017
```

`kit` will generate the service configurations for `todo` and I added `mongo` manually.

With all this setup you can run `docker-compose` up and your service will be up and running, you can now start the frontend by running `ng serve` inside `kujtimiihoxha/todo-demo` project.

After you run `docker-compose up` you should see that the container for your service will be build and your service will run like :

```text
Attaching to mongodb, todo
todo       | 2017/10/14 17:02:05 build started
todo       | Building github.com/kujtimiihoxha/todo-gokit-demo/examples/todo/cmd...
todo       | 2017/10/14 17:02:17 build completed
todo       | Running github.com/kujtimiihoxha/todo-gokit-demo/examples/todo/cmd...
todo       | ts=2017-10-14T17:02:17.239114284Z caller=service.go:78 tracer=none
todo       | ts=2017-10-14T17:02:17.239586923Z caller=service.go:100 transport=HTTP addr=:8081
todo       | ts=2017-10-14T17:02:17.239849812Z caller=service.go:134 transport=debug/HTTP addr=:8080
```

you can access your service at `http://localhost:8800` and the UI at `http://localhost:4200`.

One thing to keep in mind is that `dockerfile` and `docker-compose.yml` are not meant to be used in production they are only created to help with the development of the app.

The final app

![the todo app in action](the_todo_app_in_action.gif)

And while we update the todos in the UI the service will generate logs from the default middleware

```text
todo       | ts=2017-10-14T17:23:28.927066705Z caller=middleware.go:34 method=Add todo="{\"id\":\"\",\"title\":\"Todo\",\"complete\":false}" t="{\"id\":\"59e248107a74bb05eaaeffd4\",\"title\":\"Todo\",\"complete\":false}" error=null
todo       | ts=2017-10-14T17:23:28.927197649Z caller=middleware.go:33 method=Add transport_error=null took=2.686717ms
todo       | ts=2017-10-14T17:23:33.055004163Z caller=middleware.go:34 method=Add todo="{\"id\":\"\",\"title\":\"My Todo\",\"complete\":false}" t="{\"id\":\"59e248157a74bb05eaaeffd5\",\"title\":\"My Todo\",\"complete\":false}" error=null
todo       | ts=2017-10-14T17:23:33.055088163Z caller=middleware.go:33 method=Add transport_error=null took=611.675µs
todo       | ts=2017-10-14T17:23:33.999898829Z caller=middleware.go:40 method=SetComplete id=59e248157a74bb05eaaeffd5 error=null
todo       | ts=2017-10-14T17:23:34.000022283Z caller=middleware.go:33 method=SetComplete transport_error=null took=844.4µs
todo       | ts=2017-10-14T17:23:34.414960398Z caller=middleware.go:40 method=SetComplete id=59e248107a74bb05eaaeffd4 error=null
todo       | ts=2017-10-14T17:23:34.415056584Z caller=middleware.go:33 method=SetComplete transport_error=null took=919.311µs
todo       | ts=2017-10-14T17:23:35.363194608Z caller=middleware.go:52 method=Delete id=59e248107a74bb05eaaeffd4 error=null
todo       | ts=2017-10-14T17:23:35.363305377Z caller=middleware.go:33 method=Delete transport_error=null took=668.209µs
todo       | ts=2017-10-14T17:23:36.988528346Z caller=middleware.go:46 method=RemoveComplete id=59e248157a74bb05eaaeffd5 error=null
todo       | ts=2017-10-14T17:23:36.988574891Z caller=middleware.go:33 method=RemoveComplete transport_error=null took=644.037µs
todo       | ts=2017-10-14T17:23:37.872750792Z caller=middleware.go:52 method=Delete id=59e248157a74bb05eaaeffd5 error=null
todo       | ts=2017-10-14T17:23:37.872830897Z caller=middleware.go:33 method=Delete transport_error=null took=533.468µs
```

## Other features
A good thing about `kit` is that it can be used while you are developing and not only for one time project creation.

## Generate new service endpoints

A common thing that happens is that you want to add a new endpoint to your service that you did not think of in the beginning, `kit` makes it very easy to do that you just add the endpoint definition inside the service interface and rerun the generate command
```go
// TodoService describes the service.
type TodoService interface {
   Get(ctx context.Context) (t []io.Todo, error error)
   Add(ctx context.Context, todo io.Todo) (t io.Todo, error error)
   SetComplete(ctx context.Context, id string) (error error)
   RemoveComplete(ctx context.Context, id string) (error error)
   Delete(ctx context.Context, id string) (error error) 
   #Example we want to add a get by id method.
   GetById(ctx context.Context, id string) (t io.Todo, error error)
}
```

Then run:

```shell
kit g s todo --gorilla -w
```

`kit` will then create everything that is missing for the new endpoint and recreate the `_gen` files, this will not override any change you made in the non `_gen` files.

## Add new middleware
`kit` also supports generating new service middleware, the generator will create the boilerplate code for you but because `kit` can not possible know what parameters your middleware needs you will have to add them manually to your middleware and then add the middleware to your service.

Lets say we want to add an `auth` middleware to our service.

```shell
kit g m auth -s todo
```

this will generate the boilerplate code inside `todo/service/middleware.go`:

```go
type authMiddleware struct {
   next TodoService
}

// AuthMiddleware returns a TodoService Middleware.
func AuthMiddleware() Middleware {
   return func(next TodoService) TodoService {
      return &authMiddleware{next}
   }

}
func (a authMiddleware) Get(ctx context.Context) (t []io.Todo, error error) {
   // Implement your middleware logic here

   return a.next.Get(ctx)
}
func (a authMiddleware) Add(ctx context.Context, todo io.Todo) (t io.Todo, error error) {
   // Implement your middleware logic here

   return a.next.Add(ctx, todo)
}
func (a authMiddleware) SetComplete(ctx context.Context, id string) (error error) {
   // Implement your middleware logic here

   return a.next.SetComplete(ctx, id)
}
func (a authMiddleware) RemoveComplete(ctx context.Context, id string) (error error) {
   // Implement your middleware logic here

   return a.next.RemoveComplete(ctx, id)
}
func (a authMiddleware) Delete(ctx context.Context, id string) (error error) {
   // Implement your middleware logic here

   return a.next.Delete(ctx, id)
}
```

If you want to add an endpoint middleware just add the `-e` flag

```shell
kit g m auth -s todo -e
```

and that will generate the boilerplate code inside `todo/endpoints/middleware.go`

```go
// AuthMiddleware returns an endpoint middleware
func AuthMiddleware() endpoint.Middleware {
   return func(next endpoint.Endpoint) endpoint.Endpoint {
      return func(ctx context.Context, request interface{}) (response interface{}, err error) {
         // Add your middleware logic here
         return next(ctx, request)
      }
   }
}
```

Now to add your middleware to your service you will need to edit `todo/cmd/service/service.go#getServiceMiddleware` and `todo/cmd/service/service.go#getEndpointMiddleware` functions.

```go
func getServiceMiddleware(logger log.Logger) (mw []service.Middleware) {
   mw = []service.Middleware{}
   mw = addDefaultServiceMiddleware(logger, mw)
   // My auth middleware
   mw = append(mw, service.AuthMiddleware())
   return
}
func getEndpointMiddleware(logger log.Logger) (mw map[string][]endpoint1.Middleware) {
   mw = map[string][]endpoint1.Middleware{}
   duration := prometheus.NewSummaryFrom(prometheus1.SummaryOpts{
      Help:      "Request duration in seconds.",
      Name:      "request_duration_seconds",
      Namespace: "example",
      Subsystem: "todo",
   }, []string{"method", "success"})
   addDefaultEndpointMiddleware(logger, duration, mw)
   // My auth middleware
   addEndpointMiddlewareToAllMethods(mw,endpoint.AuthMiddleware())
   return
}
```

for adding the endpoint middleware I am using a small helper function inside `todo/cmd/service/service_gen.go` that will add the middleware to all the endpoints if you want to add the middleware to only specific endpoints you can do that by:

```go
func getEndpointMiddleware(logger log.Logger) (mw map[string][]endpoint1.Middleware) {
   mw = map[string][]endpoint1.Middleware{}
   duration := prometheus.NewSummaryFrom(prometheus1.SummaryOpts{
      Help:      "Request duration in seconds.",
      Name:      "request_duration_seconds",
      Namespace: "example",
      Subsystem: "todo",
   }, []string{"method", "success"})
   addDefaultEndpointMiddleware(logger, duration, mw)
   // My auth middleware only for the "Get" method
   mw["Get"] = append(mw["Get"], endpoint.AuthMiddleware())
   return
}
```

## Add GRPC transport

If you want to add the GRPC transport to your service you just rerun the generate command using `-t grpc`.

```shell
kit g s todo -w -t grpc
```

Since GRPC has some extra stuff you need to take care of `kit` will create your basic GRPC setup and will give you further instructions on how to complete the generation.

**You can find all the source code of the todo service in GitHub**  
[kujtimiihoxha/todo-gokit-demo](https://github.com/kujtimiihoxha/todo-gokit-demo)

**And you can find out more about GoKit-CLI here [kujtimiihoxha/kit](https://github.com/kujtimiihoxha/kit)**

