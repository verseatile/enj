# Protocol Buffers
- message format like JSON, YAML, XML, etc.
- protocol buffers are defined by message types, and value types (just like key, value in json)

# Compile a .proto file
- after installing protoc (can use brew by doing brew install protobuf) run this command
- protoc --go_out=. file.proto
- if getting "--go_out: protoc-gen-go: Plugin failed with status code 1." you need to place the location of this inside of your PATH. can 'go get' it then point it there.

# To use locally
go install (directory of go project)

# Basic protoc to Go compilation
protoc -I . todo.proto --go_out=. 

# CLI Handling
```go
// take in commands
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stdout, "missing subcommand: list or add")
		os.Exit(1)
	}

	var err error
	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list()
	case "add":
		err = add(strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
    }
```

