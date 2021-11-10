package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	files "github.com/ipfs/go-ipfs-files"
)

// get a file or directory from unixfs
func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func main() {
	sh := shell.NewShell("localhost:5001")

	file, err := os.Open("./hello")
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	reader := bufio.NewReader(file)

	cid, err := sh.Add(reader)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	fmt.Printf("successfully added file: '%s' to IPFS w/ a CID of: %s\n", "./hello", cid)

	content, err := sh.Cat(cid)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	bodyBytes, err := ioutil.ReadAll(content)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	content.Close()

	fmt.Printf("%s\n", bytes.NewBuffer(bodyBytes))

	// fmt.Printf("contents contain: %s\n", new)
}
