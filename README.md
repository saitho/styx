# Project Styx - POC

```
protoc --go_out=server/src --go_opt=paths=source_relative --go-grpc_out=server/src --go-grpc_opt=paths=source_relative ./proto/service-api.proto
```