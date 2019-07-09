# enj
grpc + protocol buffers

## Protocol Buffers
- message format like JSON, YAML, XML, etc.
- protocol buffers are defined by message types, and value types (just like key, value in json)

## Compile a .proto file
- after installing protoc (can use brew by doing brew install protobuf) run this command
- protoc --go_out=. file.proto
- if getting "--go_out: protoc-gen-go: Plugin failed with status code 1." you need to place the location of this inside of your PATH. can 'go get' it then point it there.

## To use locally
go install (directory of go project)
