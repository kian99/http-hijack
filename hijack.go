package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

const useTLS = true

func main() {
	http.HandleFunc("/auth", func(res http.ResponseWriter, req *http.Request) {
		conn, _, err := res.(http.Hijacker).Hijack()
		if err != nil {
			panic(err)
		}
		conn.Write([]byte{})
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n")

		buffer := make([]byte, 1024)
		fmt.Println("Server : Enter routine")
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Server : I send")
			_, err = conn.Write([]byte("Hijack server"))
			if err != nil {
				panic(err)
			}
			fmt.Println("Server : I'm receiving")
			n, err := conn.Read(buffer)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Server : %d bytes from client : %s\n", n, string(buffer))
		}
	})

	if useTLS {
		go runTLSClient()
	} else {
		go runNonTLSClient()
	}

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

func loopSendReceive(conn net.Conn, reader io.Reader) {
	buffer := make([]byte, 1024)
	fmt.Println("Client : Enter client routine")
	for {
		time.Sleep(250 * time.Millisecond)
		n, err := reader.Read(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Receive %d bytes : %s\n", n, string(buffer))
		conn.Write([]byte("I am Leo"))
	}
}

func runNonTLSClient() {
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	initialHTTPreq := fmt.Sprintf("CONNECT localhost:8081/auth HTTP/1.1\r\nHost: localhost:8081\r\n\r\n")
	_, err = conn.Write([]byte(initialHTTPreq))
	if err != nil {
		panic(err)
	}

	loopSendReceive(conn, conn)
}

func runTLSClient() {
	time.Sleep(1 * time.Second)
	req, err := http.NewRequest("GET", "/auth", nil)
	if err != nil {
		panic(err)
	}

	dial, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	fmt.Println("Client : create TLS connection")
	tls_conn := tls.Client(dial, &tls.Config{InsecureSkipVerify: true})
	_ = tls_conn

	var conn *httputil.ClientConn
	fmt.Println("Client : create http connection from tls client")
	conn = httputil.NewClientConn(tls_conn, nil)

	fmt.Println("Client : do request through http connection")
	_, err = conn.Do(req)
	if err != httputil.ErrPersistEOF && err != nil {
		panic(err)
	}

	fmt.Println("Client : hijack https connection")
	connection, reader := conn.Hijack()

	buffer := make([]byte, 1024)
	fmt.Println("Client : Enter client routine")
	for {
		time.Sleep(250 * time.Millisecond)
		n, err := reader.Read(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Receive %d bytes : %s\n", n, string(buffer))
		connection.Write([]byte("I am Leo"))
	}
}
