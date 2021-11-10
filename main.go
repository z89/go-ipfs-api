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

func add(sh *shell.Shell, content *bufio.Reader) string {
	cid, err := sh.Add(content)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	fmt.Printf("successfully added file: '%s' to IPFS w/ a CID of: %s\n", "./hello", cid)

	return cid
}

func get(path string) *bufio.Reader {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	return bufio.NewReader(file)
}

func read(content io.ReadCloser) []byte {
	bodyBytes, err := ioutil.ReadAll(content)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	return bodyBytes
}

func write(path string, bytes []byte, perm fs.FileMode) {
	err := ioutil.WriteFile(path, bytes, perm)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}
}

func cat(sh *shell.Shell, cid string) []byte {
	result, err := sh.Cat(cid)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	bytes := read(result)

	result.Close()

	return bytes
}

func main() {
	sh := shell.NewShell("localhost:5001")

	path := "./data/directory/picture.png"

	// get() - get any given content from unixfs
	content := get(path)

	// add() - add any given content to IPFS via the shell API
	cid := add(sh, content)

	// cat() - display contents of any given content via it's CID via the shell API
	bytes := cat(sh, cid)

	fmt.Printf("%s\n", bytes)

	// write() - write any given content to any given dir
	write("./test.png", bytes, 0644)
}
