package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/HACKERALERT/monocypher-go"
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

func read2(content *bufio.Reader) []byte {
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
	key := make([]byte, 32)
	nonce := make([]byte, 24)

	// write("./data/cipher", ciphertext, 0644)

	// plaintext2, authentic := monocypher.Unlock(ciphertext, nonce, key, mac)
	// if !authentic {
	// 	panic(fmt.Errorf("not authentic"))
	// }

	// fmt.Printf("plaintext: %s\n", plaintext2)
	// fmt.Printf("%s\n", read(bufio.NewReader('./data/directory/custom'))

	sh := shell.NewShell("localhost:5001")

	// cleartext picture
	file := "./data/directory/picture.png"

	// open picture with type *os.File
	fileTe, err := os.Open(file)
	if err != nil {
		panic(fmt.Errorf("error: %s", err))
	}

	// use the file to create a new reader with type *bufio.Reader
	t := bufio.NewReader(fileTe)

	// use *bufio.Reader & return a []byte
	p := read2(t)

	mac, ciphertext := monocypher.Lock(p, nonce, key)

	path := "./data/cipher"

	write(path, ciphertext, 0644)

	// dir := "./data/directory"

	// // add() - add any given file to IPFS
	fileCid := add(sh, path) // add cleartext
	add(sh, file)            // add ciphertext

	// // addDir() - add any given directory to IPFS
	// dirCid := addDir(sh, dir)

	// // cat() - get contents of any given file via it's CID from IPFS
	output := cat(sh, fileCid) // get plaintext from IPFS

	fmt.Printf("\n")

	// // // write() - write any given file to any given dir
	write("./data/encrypted-ipfs-pic.png", output, 0644) // save plaintext picture

	encryptedPicture := "./data/encrypted-ipfs-pic.png"

	encryptedFile, err := os.Open(encryptedPicture)
	if err != nil {
		panic(fmt.Errorf("error: %s", err))
	}

	encryptedByte := bufio.NewReader(encryptedFile)

	readout := read2(encryptedByte)

	plaintextPicture, authentic := monocypher.Unlock(readout, nonce, key, mac)
	if !authentic {
		panic(fmt.Errorf("not authentic"))
	}

	write("./data/decrypted-ipfs-pic.png", plaintextPicture, 0644)

	// fmt.Printf("\n")

	// // ls() - get contents of any given directory via it's CID from IPFS
	// dirContents := ls(sh, dirCid)

	// // print each Links object w/ type (*shell.UnixLsLink) from the listed directory
	// fmt.Printf("directory %s:\n", dirContents.Hash)
	// for _, v := range dirContents.Links {
	// 	if v.Type == "File" {
	// 		fmt.Printf(" - file: %s Name: %s Type: %s Size: %d bytes\n", v.Hash, v.Name, v.Type, v.Size)
	// 	} else if v.Type == "Directory" {
	// 		fmt.Printf(" - dir: %s Name: %s Type: %s Size: %d bytes\n", v.Hash, v.Name, v.Type, v.Size)
	// 	}
	// }

}
