package main

import (
	"bufio"
	"fmt"
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

// func addDir(sh *shell.Shell, path string) string {
// 	cid, err := sh.AddDir(path)
// 	if err != nil {
// 		panic(fmt.Errorf("failed: %s", err))
// 	}

// 	fmt.Printf("successfully added directory: '%s' to IPFS w/ a CID of: %s\n", path, cid)

// 	return cid
// }

// func ls(sh *shell.Shell, cid string) *shell.UnixLsObject {
// 	dir, err := sh.FileList(cid)
// 	if err != nil {
// 		panic(fmt.Errorf("yeet: %s", err))
// 	}

// 	return dir
// }

func open(path string) *bufio.Reader {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("error: %s", err))
	}

	result := bufio.NewReader(file)

	return result
}

func read(reader *bufio.Reader) []byte {
	raw, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(fmt.Errorf("failed: %s", err))
	}

	return raw
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

func encrypt(plaintext []byte, nonce []byte, key []byte) ([]byte, []byte) {
	mac, ciphertext := monocypher.Lock(plaintext, nonce, key)

	fmt.Printf("mac exists: %s\n", mac)
	return ciphertext, mac
}

func decrypt(ciphertext []byte, nonce []byte, key []byte, mac []byte) []byte {
	decrypted, authentic := monocypher.Unlock((ciphertext), nonce, key, mac)
	if !authentic {
		panic(fmt.Errorf("not authentic"))
	}

	return decrypted
}

func main() {
	key := make([]byte, 32)
	nonce := make([]byte, 24)

	sh := shell.NewShell("localhost:5001")

	// plaintext picture inside two nested directories
	file := "./data/directory/picture.png"

	// read plaintext picture from unifs (flow: path string > reader *bufio.Reader > raw []byte)
	plaintextFile := read(open(file))

	// encrypt plaintext picture using key; nonce vars & output ciphertext; mac vars
	ciphertext, mac := encrypt(plaintextFile, nonce, key)

	// add ciphertext picture to IPFS
	ciphertextCID := add(sh, path)

	// cat (get contents of IPFS file) ciphertext picture from IPFS
	output := cat(sh, ciphertextCID)

	// write ciphertext picture fetched from IPFS to new file
	write("./data/encrypted-ipfs-pic.png", output, 0644)

	// define path of newly written ciphertext picture
	ciphertextPic := "./data/encrypted-ipfs-pic.png"

	// read ciphertext picture from unifs (flow: path string > reader *bufio.Reader > raw []byte)
	ciphertextPicBytes := read(open(ciphertextPic))

	// decrypt ciphertext picture using key; nonce; mac vars & output plaintext
	plaintextPic := decrypt(ciphertextPicBytes, nonce, key, mac)

	// write decrypted plaintext picture to a new file
	write("./data/decrypted-ipfs-pic.png", plaintextPic, 0644)
}
