package main

import (
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

func main() {
	connection, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	for {
		//Waiting for connection ...
		conn, err := connection.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//Connected to server, start receiving the file name and file size
		bufferFileName := make([]byte, 64)
		bufferFileSize := make([]byte, 10)

		conn.Read(bufferFileSize)
		fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

		conn.Read(bufferFileName)
		fileName := strings.Trim(string(bufferFileName), ":")

		newFile, err := os.Create(filepath.Join("/home/Files", fileName)) //you have to create a dir name Files under home dir

		if err != nil {
			panic(err)
		}
		defer newFile.Close()
		var receivedBytes int64

		for {
			if (fileSize - receivedBytes) < BUFFERSIZE {
				io.CopyN(newFile, conn, (fileSize - receivedBytes))
				conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
				break
			}
			io.CopyN(newFile, conn, BUFFERSIZE)
			receivedBytes += BUFFERSIZE
		}
		//Received file completely!
	}
}
