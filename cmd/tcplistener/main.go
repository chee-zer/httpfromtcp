package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)


func main() {

	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Couldnt set up listener: ", err)
	}
	defer l.Close()

	for {
		con, err := l.Accept()
		if err != nil {
			log.Fatal("Couldnt accept connection: ", err)
		}
		log.Println("Connection is accepted successfully!")
		ch := GetLinesChannel(con)

		for line := range ch {
			fmt.Println(line)
		}
		log.Println("Connection has been closed.")
	}

	// lineChannel := GetLinesChannel()

	// for line := range lineChannel {
	// 	fmt.Printf("read: %s\n", line)
	// }
}


func GetLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	var currentLineContents string
	go func() {

		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
	
			for i := 0; i < len(parts) - 1; i++ {
				currentLineContents = currentLineContents + parts[i]
				ch <- currentLineContents
				currentLineContents = ""
			}
			currentLineContents = currentLineContents + parts[len(parts) - 1]
		}
		close(ch)
	} ()
	return ch
} 