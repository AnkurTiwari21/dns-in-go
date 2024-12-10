package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/app/models"
)

// Ensures gofmt doesn't remove the "net" import in stage 1 (feel free to remove this!)
var _ = net.ListenUDP

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		// response := []byte{}
		response := models.Message{
			Header: models.Header{
				PacketIdentifier:      1234,
				QueryIndicator:        1,
				OperationCode:         0,
				AuthoritativeAnswer:   0,
				Truncation:            0,
				RecursionDesired:      0,
				RecursionAvailable:    0,
				Reserved:              0,
				ResponseCode:          0,
				QuestionCount:         0,
				AnswerRecordCount:     0,
				AuthorityRecordCount:  0,
				AdditionalRecordCount: 0,
			},
		}
		responseMarshalled, _ := json.Marshal(response)
		_, err = udpConn.WriteToUDP(responseMarshalled, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
