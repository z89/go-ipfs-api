package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

func addFile(sh *shell.Shell, file io.Reader) string {
	cid, err := sh.Add(bufio.NewReader(file))
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	fmt.Printf("successfully added file to IPFS w/ a CID of: %s\n", cid)

	return cid
}

func open(path string) *bufio.Reader {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("error: %s", err))
	}

	result := bufio.NewReader(file)

	return result
}

func write(path string, bytes []byte, perm fs.FileMode) string {
	err := ioutil.WriteFile(path, bytes, perm)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	fmt.Printf("successfully wrote file to: '%s'\n", path)

	return path
}

func cat(sh *shell.Shell, cid string) []byte {
	result, err := sh.Cat(cid)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	bytes, err := ioutil.ReadAll(result)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	result.Close()

	return bytes
}

func main() {
	sh := shell.NewShell("localhost:5001")

	file := "./hello"

	// add (send file to IPFS instance and return CID)
	ciphertextCID := addFile(sh, open(file))

	// cat (get contents of a file via it's CID from an IPFS instance)
	output := cat(sh, ciphertextCID)

	// write (write a file to a unixfs)
	write(file+"-fromIPFS", output, 0644)
}
