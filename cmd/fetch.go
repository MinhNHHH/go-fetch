package main

import (
	"os"

	f "github.com/minhnh/fetch/pkg/fetch"
)

func main() {
	var cmd []string
	if len(os.Args) < 2 {
		cmd = []string{""}
	} else {
		cmd = os.Args[1:]
	}
	f.HandleClient(cmd)
}
