//go:build ignore

package main

import (
	"log"
	"os"
	"os/exec"
)

func Setup() error {
	cmd := exec.Command("pipx", "install", "bikeshed")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("bikeshed", "update")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Dev() error {
	cmd := exec.Command("bikeshed", "serve")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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
		"setup": Setup,
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
