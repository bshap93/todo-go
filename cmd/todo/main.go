package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"pragprog.com/rggo/interacting/todo"
)

var todoFileName = ".todo.json"

func main() {
	// Parsing cmd line flags
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	delete := flag.Int("del", 0, "Item to be deleted")
	verbose := flag.Bool("verbose", false, "Output should be verbose")

	flag.Parse()
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}
	// same as before
	inst := &todo.Instance{}
	lst := &todo.List{}
	inst.ListInst = *lst
	inst.Verbose = false

	if err := lst.Get(todoFileName); err != nil {

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *verbose:
		inst.Verbose = true
	case *list:
		if inst.Verbose {
			inst.ListInst.VerbosePrint()
		} else {
			fmt.Print(lst)
		}

	case *delete > 0:
		if err := lst.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	case *complete > 0:
		if err := lst.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := lst.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		lst.Add(t)
		if err := lst.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}

}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}
	return s.Text(), nil
}
