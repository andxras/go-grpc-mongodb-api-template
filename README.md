## MONGODB CRUD API

### setup

Download protoc binary from: https://developers.google.com/protocol-buffers/docs/downloads
for the required os and arch and place it in /proto dir


#### generate proto files

```
cd proto && protoc *.proto --go_out=plugins=grpc:.
```


#### start server & client

```
go run client/main.go
go run core/main.go
```
