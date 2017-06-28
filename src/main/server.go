package main

import (
	"strings"
	"strconv"
	"time"
	"net"
	"log"
	"bufio"
	"bytes"
	"io"
	"net/url"
	"runtime"
)


type request struct {
	method string
	url string
	protocol string
}


func startServer()  {
	port := 80;

	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		panic("Failed start server: " + err.Error());
	}
	log.Print("Server started at " + strconv.Itoa(port) + " port")

	ch := make(chan net.Conn)

	runtime.GOMAXPROCS(*CPU_NUM)


	for i:=0; i < *CPU_NUM; i++ {
		println("Created worker...")
		go handleConnection(ch)
	}


	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accept new connection: %s", err)
			continue
		}
		ch <- conn
	}
}


func handleConnection(ch chan net.Conn) {
	for {
		conn := <-ch

		bufr := bufio.NewReader(conn)
		buf := make([]byte, 1024)

		var input bytes.Buffer

		for {
			readBytes, err := bufr.Read(buf)
			if err != nil {
				if (err != io.EOF) {
					log.Printf("handle connection error, err=%s", err)
				}
				break
			}

			input.Write(buf[:readBytes]) // Сохраняем полученные данные

			httpRequestEnd := "\r\n\r\n"
			readString := string(buf[:readBytes])
			if strings.Contains(readString, httpRequestEnd) {
				break // конец запроса
			}
		}

		response := parseInputData(input.String())
		conn.Write(response)
		conn.Close()
	}
}






func parseInputData(input string) ([]byte) {
	var infoLine= strings.Split(input, "\r\n")[0]

	var splitLine= strings.Split(infoLine, " ")

	var response bytes.Buffer
	var file File

	if (len(splitLine) < 3) {
		log.Println("Failed parsing 1st line user's request!")
		response.WriteString("400 Bad Request")
		response.WriteString("\r\n")

	} else {
		decoded_url, _ := url.QueryUnescape(splitLine[1])
		userRequest := request{
			method:   splitLine[0],
			url:      decoded_url,
			protocol: splitLine[2],
		}

		response.WriteString(userRequest.protocol)
		response.WriteString(" ")


		if !isMethodAllowed(userRequest.method) {
			response.WriteString("405 Method Not Allowed")
			response.WriteString("\r\n")
			response.WriteString("Allow: GET, HEAD")
			response.WriteString("\r\n")

		} else {
			head := strings.Compare(userRequest.method, "HEAD") == 0;
			file = GetFile(userRequest.url, head)

			switch file.status {
			case 200:
				response.WriteString("200 OK\r\n")
				response.WriteString("Content-Type: " +
					file.content_type + "\r\n")
				response.WriteString("Content-Length: " +
					strconv.Itoa(file.length) + "\r\n")
				break

			case 403:
				response.WriteString("403 Forbidden\r\n")
				break

			case 404:
				response.WriteString("404 File Not Found\r\n")
				break
			default:
				break
			}
		}
	}



	// Дописываем хедеры
	response.WriteString("Date: " + time.Now().String())
	response.WriteString("\r\n")
	response.WriteString("Server: Golang HTTP Server")
	response.WriteString("\r\n")
	response.WriteString("Connection: Close")
	response.WriteString("\r\n")
	response.WriteString("\r\n")

	if (file.status == 200) {
		response.Write(file.content)
	}

	return response.Bytes()
}


func isMethodAllowed(method string) (bool)  {
	return strings.Compare(method, "GET") == 0 || strings.Compare(method, "HEAD") == 0;
}

