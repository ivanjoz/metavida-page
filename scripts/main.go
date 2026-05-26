package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run . {generate_schemas}")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate_schemas":
		runSubpackage("./schemas")
	default:
		fmt.Fprintf(os.Stderr, "unknown script: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runSubpackage(pkg string) {
	cmd := exec.Command("go", "run", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running %s: %v\n", pkg, err)
		os.Exit(1)
	}
}
