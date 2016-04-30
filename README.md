Simple example of sing grpc
======================


PREREQUISITES
-------------

- This requires Go 1.6
- Requires that [GOPATH is set](https://golang.org/doc/code.html#GOPATH)

```
$ go help gopath
$ # ensure the PATH contains $GOPATH/bin
$ export PATH=$PATH:$GOPATH/bin
```

INSTALL
-------

```
$ go get -u github.com/gengo/grpc-gateway/protoc-gen-grpc-gateway
$ go get -u github.com/gengo/grpc-gateway/protoc-gen-swagger
$ go get -u github.com/golang/protobuf/protoc-gen-go
$ go get -u github.com/olegsmetanin/golang-grpc-rest-gorm-example
```

Rebuilding the proto file
-------
```
$ # from this dir; invoke protoc
$  protoc -I ./api/ ./api/api.proto --go_out=plugins=grpc:./api/proto
$  protoc -I ./api/ ./api/api.proto --grpc-gateway_out=logtostderr=true:./api/proto

```

# Generate gRPC stub
$ protoc -I/usr/local/include -I./api \
 -I$GOPATH/src/github.com/gengo/grpc-gateway/third_party/googleapis \
 --go_out=Mgoogle/api/annotations.proto=github.com/gengo/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:./api/proto \
 ./api/api.proto

# Generate reverse-proxy
$ protoc -I/usr/local/include -I./api \
 -I$GOPATH/src/github.com/gengo/grpc-gateway/third_party/googleapis \
 --grpc-gateway_out=logtostderr=true:./api/proto \
 ./api/api.proto

# Generate swagger definitions
$ protoc -I/usr/local/include -I./api \
 -I$GOPATH/src/github.com/gengo/grpc-gateway/third_party/googleapis \
 --swagger_out=logtostderr=true:./api \
 ./api/api.proto


