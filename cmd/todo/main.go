package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"grpc-start/todo"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
)

func main() {
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
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// fmt.Println("TODO")

}

const dbPath = "mydb.pb"

func add(text string) error {
	task := &todo.Task{
		Text: text, Done: false}

	b, err := proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("could not encode task")
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", dbPath, err)
	}
	if err := gob.NewEncoder(f).Encode(int64(len(b))); err != nil {
		return fmt.Errorf("could not encode length of message: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could write task to file: %v", err)
	}

	// fmt.Println(proto.MarshalTextString(task))
	return nil
}

func list() error {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("could not find file: %s", dbPath)
	}

	for {
		if len(b) == 0 {
			return nil
		} else if len(b) < 4 {
			return fmt.Errorf("remaining odd %d bytes, what to do?", len(b))
		}

		var length int64
		if err := gob.NewDecoder(bytes.NewReader(b[:4])).Decode(&length); err != nil {
			return fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[4:]

		var task todo.Task
		if err := proto.Unmarshal(b[:length], &task); err != nil {
			return fmt.Errorf("could not read task: %v", err)
		}
		b = b[length:]

		if task.Done {
			fmt.Printf("🆗👍 ")
		} else {
			fmt.Printf("😨 ")
		}
		fmt.Printf("[Task Text]: %s \n", task.Text)
	}

	// return nil
}
