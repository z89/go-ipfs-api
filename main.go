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

func add(sh *shell.Shell, path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	cid, err := sh.Add(bufio.NewReader(file))
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	fmt.Printf("successfully added file: '%s' to IPFS w/ a CID of: %s\n", path, cid)

	return cid
}

func addDir(sh *shell.Shell, path string) string {
	cid, err := sh.AddDir(path)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	fmt.Printf("successfully added directory: '%s' to IPFS w/ a CID of: %s\n", path, cid)

	return cid
}

func ls(sh *shell.Shell, cid string) *shell.UnixLsObject {
	dir, err := sh.FileList(cid)
	if err != nil {
		panic(fmt.Errorf("yeet: %s", err))
	}

	return dir
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

	fmt.Printf("successfully wrote file to: '%s'\n", path)
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

	file := "./data/directory/picture.png"
	dir := "./data/directory"

	// add() - add any given file to IPFS
	fileCid := add(sh, file)

	// addDir() - add any given directory to IPFS
	dirCid := addDir(sh, dir)

	// cat() - get contents of any given file via it's CID from IPFS
	output := cat(sh, fileCid)

	fmt.Printf("\n")

	// write() - write any given file to any given dir
	write("./data/test.png", output, 0644)

	fmt.Printf("\n")

	// ls() - get contents of any given directory via it's CID from IPFS
	dirContents := ls(sh, dirCid)

	// print each Links object w/ type (*shell.UnixLsLink) from the listed directory
	fmt.Printf("directory %s:\n", dirContents.Hash)
	for _, v := range dirContents.Links {
		if v.Type == "File" {
			fmt.Printf(" - file: %s Name: %s Type: %s Size: %d bytes\n", v.Hash, v.Name, v.Type, v.Size)
		} else if v.Type == "Directory" {
			fmt.Printf(" - dir: %s Name: %s Type: %s Size: %d bytes\n", v.Hash, v.Name, v.Type, v.Size)
		}
	}

}
