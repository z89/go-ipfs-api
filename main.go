package main

import (
	"fmt"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	sh := shell.NewShell("localhost:5001")

	str := "hello world"
	cid, err := sh.Add(strings.NewReader(str))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("successfully added string: '%s' to IPFS w/ a CID of: %s\n", str, cid)
}
