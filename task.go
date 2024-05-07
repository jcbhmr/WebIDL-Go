//go:build ignore

package main

import (
	"log"
	"os"
	"os/exec"
)

func Dev() error {
	cmd := exec.Command("go", "run", "github.com/jcbhmr/go-bikeshed/cmd/bikeshed", "serve")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("$ %s", cmd.String())
	return cmd.Run()
}

func main() {
	log.SetFlags(0)
	var taskName string
	if len(os.Args) >= 2 {
		taskName = os.Args[1]
	} else {
		log.Fatal("no task")
	}
	task, ok := map[string]func() error{
		"dev":   Dev,
	}[taskName]
	if !ok {
		log.Fatal("no such task")
	}
	err := task()
	if err != nil {
		log.Fatal(err)
	}
}
