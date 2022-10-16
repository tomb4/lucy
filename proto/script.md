
# Party 
```
protoc --go_out=./party --go_opt=paths=source_relative \
--go-grpc_out=./party --go-grpc_opt=paths=source_relative \
party.proto
```

# Gateway
```
protoc --go_out=./MetaGateway --go_opt=paths=source_relative \
--go-grpc_out=./MetaGateway --go-grpc_opt=paths=source_relative \
MetaGateway.proto
```