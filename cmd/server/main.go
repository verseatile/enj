package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"grpc-start/todo"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var port = ":9000"

type length int64

const (
	sizeOfLength = 8
	dbPath       = "mydb.pb"
)

var endianness = binary.LittleEndian

type taskServer struct {
}

// Void - imitation void type
type Void struct{}

func main() {

	var tasks taskServer
	// start server
	srv := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute}))
	// <--- This fixes it!

	todo.RegisterTasksServer(srv, tasks)
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not listen to :9000: %v", err)
	} else {
		log.Printf("listening....\n")
	}
	log.Fatal(srv.Serve(l))

}

func (taskServer) Add(ctx context.Context, text *todo.Text) (*todo.Task, error) {
	task := &todo.Task{
		Text: text.Text,
		Done: false,
	}

	b, err := proto.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("could not encode task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %v", dbPath, err)
	}

	if err := binary.Write(f, endianness, length(len(b))); err != nil {
		return nil, fmt.Errorf("could not encode length of message: %v", err)
	}
	_, err = f.Write(b)
	if err != nil {
		return nil, fmt.Errorf("could not write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close file %s: %v", dbPath, err)
	}
	return task, nil
}

func (taskServer) List(ctx context.Context, void *todo.Void) (*todo.TaskList, error) {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	var tasks todo.TaskList
	for {
		if len(b) == 0 {
			return &tasks, nil
		} else if len(b) < sizeOfLength {
			return nil, fmt.Errorf("remaining odd %d bytes, what to do?", len(b))
		}

		var l length
		if err := binary.Read(bytes.NewReader(b[:sizeOfLength]), endianness, &l); err != nil {
			return nil, fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[sizeOfLength:]

		var task todo.Task
		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return nil, fmt.Errorf("could not read task: %v", err)
		}
		b = b[l:]
		tasks.Tasks = append(tasks.Tasks, &task)
	}
}

// flag.Parse()
// 	if flag.NArg() < 1 {
// 		fmt.Fprintln(os.Stdout, "missing subcommand: list or add")
// 		os.Exit(1)
// 	}

// 	var err error
// 	switch cmd := flag.Arg(0); cmd {
// 	case "list":
// 		err = list()
// 	case "add":
// 		err = add(strings.Join(flag.Args()[1:], " "))
// 	default:
// 		err = fmt.Errorf("unknown subcommand %s", cmd)
// 	}
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}
