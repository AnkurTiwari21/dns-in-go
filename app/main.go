package main

import (
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
		for _, b := range buf[:size] {
			fmt.Printf("%08b\n", b)
		}
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		response := models.Message{
			Header:   models.Header{},
			Question: models.Question{},
			Answer:   models.Answer{},
		}
		
		// response.Header.SetFlags(1, 0, 0, 0, 0, 0, 0, 0)                              //setting up flag
		headerBytes := response.Header.SetRemainingDataAndReturnBytes(buf[:size]) //sending remaining data and getting header bytes
		responseBytes := response.Bytes(headerBytes)

		// questionBytes := response.Question.SetAllDataAndReturnQuestionBytes("codecrafters.io", 1, 1)
		// responseBytes = append(responseBytes, questionBytes...) //appending question bytes

		// answerBytes := response.Answer.FillAnswerAndReturnBytes()
		// responseBytes = append(responseBytes, answerBytes...)

		fmt.Print("<------->")
		for _, b := range responseBytes {
			fmt.Printf("%08b ", b)
		}
		fmt.Print(len(responseBytes))
		fmt.Print("<>")
		_, err = udpConn.WriteToUDP(responseBytes, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
