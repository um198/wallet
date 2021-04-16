package main

import (
	"io"
	"log"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(wd)

	file, err := os.Open("data/readme.txt")
	if err != nil {
		log.Print(err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	log.Printf("%#v", file)

	buf := make([]byte, 4)
	content := make([]byte, 0)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}
		if err != nil {
			log.Print(err)
			return
		}

		content = append(content, buf[:read]...)

	}
	data := string(content)
	log.Print(data)

	file2, err := os.Create("data/message.txt")
	if err != nil {
		log.Print(err)
		return
	}

	defer func() {
		if cerr := file2.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	_, err = file2.Write([]byte("Hello from GO!"))
	if err != nil {
		log.Print(err)
		return
	}
	log.Print()
}

func closeF(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Print(err)
		return
	}

}
