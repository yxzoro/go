package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	port := ":18000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
	checkError(err)
	log.Println("Socket Server Started at :18000...")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		checkError(err)
		log.Printf("=> got a conection from: %s", conn.RemoteAddr().String())
		// handle a conn in a goroutine
		go handler(conn)
	}

}

func checkError(err error) {
	if err != nil {
		log.Printf("error: %s", err.Error())
	}
}

func handler(conn net.Conn) {
	log.Printf("=> Reading data from conn: %s", conn.RemoteAddr().String())

	// ************************************************************************************
	// use tmp as a buffer, append all bytes to total
	// you can see here:
	//	  just copy data bytes from network buffer to process buffer...
	// ************************************************************************************
	buf := make([]byte, 0, 1024) // total buffer, capcity=1024 bytes
	tmp := make([]byte, 10)      // tmp buffer, capcity=10 bytes
	for {
		log.Println("-> waiting for conn.Read()...")
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Printf("unexpected read error: %s", err)
			} else {
				log.Println("read EOF. byebye.")
				break
			}
		}
		log.Printf("got %d bytes:", n)
		log.Println(tmp)
		log.Printf("(bytes to string: %s)", tmp)
		buf = append(buf, tmp[:n]...)
	}
	log.Println("end...total bytes is:")
	log.Println(buf)

	daytime := time.Now().String()
	conn.Write([]byte("i have send you the time: " + daytime))
	conn.Close()

}