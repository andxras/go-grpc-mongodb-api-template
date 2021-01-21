## MONGODB CRUD API

### Setup

Download protoc binary from: https://developers.google.com/protocol-buffers/docs/downloads
for the required os and arch and place it in /proto dir


#### Generate proto files

```
cd proto && protoc *.proto --go_out=plugins=grpc:.
```


#### Start

```
go run client/main.go
```
