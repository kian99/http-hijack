package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

const useTLS = true

func main() {
	http.HandleFunc("/foo", func(res http.ResponseWriter, req *http.Request) {
		conn, _, err := res.(http.Hijacker).Hijack()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		buffer := make([]byte, 1024)
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Server sending")
			_, err = conn.Write([]byte("Hijack server"))
			if err != nil {
				panic(err)
			}
			fmt.Println("Server receiving")
			n, err := conn.Read(buffer)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Server: %d bytes from client -> %s\n", n, string(buffer))
		}
	})

	go runClient()

	var err error
	if useTLS {
		err = http.ListenAndServeTLS(":8081", "./server.crt", "./server.key", nil)
	} else {
		err = http.ListenAndServe(":8081", nil)
	}

	if err != nil {
		panic(err)
	}
}

func loopSendReceive(conn net.Conn) {
	buffer := make([]byte, 1024)
	fmt.Println("Client: Enter routine to send and receive")
	for {
		time.Sleep(250 * time.Millisecond)
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Client: Received %d bytes -> %s\n", n, string(buffer[:n]))
		conn.Write([]byte("Hello world"))
	}
}

func runClient() {
	time.Sleep(1 * time.Second)
	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		panic(err)
	}

	dial, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	if useTLS {
		dial = tls.Client(dial, &tls.Config{InsecureSkipVerify: true})
	}

	err = req.Write(dial)
	if err != nil {
		panic(err)
	}
	loopSendReceive(dial)
}
